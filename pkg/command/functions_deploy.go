package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/altipla-consulting/chrono"
	"github.com/hashicorp/hcl"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/ernestoalejo/tfg-fn/pkg/minikube"
	"github.com/ernestoalejo/tfg-fn/pkg/spec"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

func init() {
	FunctionsCmd.AddCommand(FunctionsDeployCmd)
}

var FunctionsDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy function",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return errors.Trace(err)
		}
		app := new(spec.Application)
		if err := hcl.Unmarshal(content, app); err != nil {
			return errors.Trace(err)
		}

		dir, err := ioutil.TempDir("", "fndeploy")
		if err != nil {
			return errors.Trace(err)
		}
		defer os.RemoveAll(dir)
		logrus.WithFields(logrus.Fields{"directory": dir, "function": app.Function.Name}).Info("build production container")

		logrus.Info("copy runtime files")
		c := exec.Command("bash", "-c", fmt.Sprintf("cp runtime/nodejs/* %s", dir))
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return errors.Trace(err)
		}

		logrus.Info("copy app files")
		copyFiles := []string{}
		for _, file := range app.Function.Files {
			copyFiles = append(copyFiles, fmt.Sprintf("COPY %s /opt/%s", file, file))
			file = filepath.Join(filepath.Dir(args[0]), file)

			c = exec.Command("bash", "-c", fmt.Sprintf("cp %s %s", file, dir))
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return errors.Trace(err)
			}
		}

		logrus.Info("build Dockerfile")
		dockerfile, err := ioutil.ReadFile(filepath.Join(dir, "Dockerfile"))
		if err != nil {
			return errors.Trace(err)
		}
		dockerfile = []byte(fmt.Sprintf(string(dockerfile), strings.Join(copyFiles, "\n")))
		if err := ioutil.WriteFile(filepath.Join(dir, "Dockerfile"), dockerfile, 0600); err != nil {
			return errors.Trace(err)
		}

		logrus.Info("build container")
		registryAddress, err := minikube.ServiceIPPort("registry")
		if err != nil {
			return errors.Trace(err)
		}
		imageName := fmt.Sprintf("%s/%s", registryAddress, app.Function.Name)
		c = exec.Command("docker", "build", "-t", imageName, dir)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return errors.Trace(err)
		}

		logrus.Info("push container to registry")
		c = exec.Command("docker", "push", imageName)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return errors.Trace(err)
		}

		logrus.Info("register new function in api server")
		fn := &pb.Function{
			Name:      app.Function.Name,
			Call:      app.Function.Call,
			Trigger:   app.Function.Trigger,
			Method:    app.Function.Method,
			CreatedAt: chrono.DateTimeToProto(time.Now()),
		}
		if _, err := client.DeployFunction(context.Background(), &pb.DeployFunctionRequest{Function: fn}); err != nil {
			return errors.Trace(err)
		}

		logrus.Info("function deployed successfully")

		return nil
	},
}

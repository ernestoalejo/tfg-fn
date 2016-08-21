package command

import (
	"io/ioutil"
	"time"

	"github.com/altipla-consulting/chrono"
	"github.com/hashicorp/hcl"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

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

		return nil
	},
}

package minikube

import (
	"bytes"
	"net/url"
	"os/exec"
	"strings"

	"github.com/juju/errors"
)

func ServiceIPPort(name string) (string, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command("minikube", "service", "--url=true", name)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", errors.Trace(err)
	}

	u, err := url.Parse(strings.TrimSpace(buf.String()))
	if err != nil {
		return "", errors.Trace(err)
	}
	return u.Host, nil
}

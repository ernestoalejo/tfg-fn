package command

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "fnctl",
	Short: "Control the fn cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("use a subcommand")
	},
}

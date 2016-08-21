package command

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(FunctionsCmd)
}

var FunctionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Functions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("use a subcommand")
	},
}

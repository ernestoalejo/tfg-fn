package command

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	pb "github.com/ernestoalejo/tfg-fn/protos"
)

func init() {
	FunctionsCmd.AddCommand(FunctionsDeleteCmd)
}

var FunctionsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete function",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := client.DeleteFunction(context.Background(), &pb.DeleteFunctionRequest{Name: args[0]}); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}

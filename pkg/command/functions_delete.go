package command

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"

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
		logrus.WithFields(logrus.Fields{"function": args[0]}).Info("delete function")

		return nil
	},
}

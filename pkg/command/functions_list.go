package command

import (
	"time"

	"github.com/altipla-consulting/chrono"
	pb_empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	FunctionsCmd.AddCommand(FunctionsListCmd)
}

var FunctionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cluster functions",
	RunE: func(cmd *cobra.Command, args []string) error {
		reply, err := client.ListFunctions(context.Background(), new(pb_empty.Empty))
		if err != nil {
			return errors.Trace(err)
		}

		tabPrint([]string{
			"NAME",
			"DEPLOYED",
		})
		for _, fn := range reply.Functions {
			tabPrint([]string{
				fn.Name,
				time.Now().Sub(chrono.DateTimeFromProto(fn.CreatedAt)).String(),
			})
		}

		return nil
	},
}

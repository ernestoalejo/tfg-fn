package command

import (
	"fmt"
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
			"CALL",
			"TRIGGER",
			"DEPLOYED",
		})
		for _, fn := range reply.Functions {
			var trigger string
			switch fn.Trigger {
			case "http":
				trigger = fmt.Sprintf("http (method=%s)", fn.Method)
			}

			tabPrint([]string{
				fn.Name,
				fn.Call,
				trigger,
				time.Now().Sub(chrono.DateTimeFromProto(fn.CreatedAt)).String(),
			})
		}

		return nil
	},
}

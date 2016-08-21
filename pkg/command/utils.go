package command

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	pb "github.com/ernestoalejo/tfg-fn/protos"
)

var (
	client pb.FnClient
	writer = tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
)

func SetClient(c pb.FnClient) {
	client = c
}

func FlushOutput() {
	writer.Flush()
}

func tabPrint(fields []string) {
	fmt.Fprintf(writer, strings.Join(fields, "\t")+"\n")
}

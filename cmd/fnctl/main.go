package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/ernestoalejo/tfg-fn/pkg/command"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

func main() {
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to connect")
	}
	defer conn.Close()
	command.SetClient(pb.NewFnClient(conn))

	if err := command.RootCmd.Execute(); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to run command")
		os.Exit(-1)
	}
	command.FlushOutput()
}

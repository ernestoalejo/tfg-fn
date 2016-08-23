package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/ernestoalejo/tfg-fn/pkg/command"
	"github.com/ernestoalejo/tfg-fn/pkg/minikube"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

func main() {
	fnapi, err := minikube.ServiceIPPort("fnapi")
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to fetch fnapi address")
	}

	conn, err := grpc.Dial(fnapi, grpc.WithInsecure())
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "address": fnapi}).Fatal("failed to connect")
	}
	defer conn.Close()
	command.SetClient(pb.NewFnClient(conn))

	if err := command.RootCmd.Execute(); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to run command")
		os.Exit(-1)
	}
	command.FlushOutput()
}

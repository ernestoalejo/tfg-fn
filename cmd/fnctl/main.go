package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/ernestoalejo/tfg-fn/pkg/command"
	"github.com/ernestoalejo/tfg-fn/pkg/minikube"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

func main() {
	fnapi, err := minikube.ServiceIPPort("fnapi")
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to fetch fnapi address")
	}

	tlsCert, err := ioutil.ReadFile("certs/clients/fnctl.pem")
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Fatal("failed to load certificate")
	}
	tlsKey, err := ioutil.ReadFile("certs/clients/fnctl-key.pem")
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Fatal("failed to load certificate")
	}
	tlsCA, err := ioutil.ReadFile("certs/ca.pem")
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Fatal("failed to load certificate")
	}
	cert, err := tls.X509KeyPair(tlsCert, tlsKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Fatal("failed to load certificate")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(tlsCA)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})
	conn, err := grpc.Dial(fnapi, grpc.WithTransportCredentials(creds))
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

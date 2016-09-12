package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/husobee/vestigo"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	r "gopkg.in/dancannon/gorethink.v2"

	"github.com/ernestoalejo/tfg-fn/pkg/api"
	fnctx "github.com/ernestoalejo/tfg-fn/pkg/context"
	"github.com/ernestoalejo/tfg-fn/pkg/proxy"
	"github.com/ernestoalejo/tfg-fn/pkg/web"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

func main() {
	address := os.Getenv("DATABASE_ADDRESS")
	if address == "" {
		address = "localhost"
	}
	db, err := r.Connect(r.ConnectOpts{
		Address:  address,
		Database: "fn",
		Username: "fn",
		Password: "fnpass",
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to connect to RethinkDB")
	}

	tlsCert, err := ioutil.ReadFile("certs/server.pem")
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Fatal("failed to load certificate")
	}
	tlsKey, err := ioutil.ReadFile("certs/server-key.pem")
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
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	listener, err := net.Listen("tcp", ":50050")
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("failed to listen")
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor), grpc.Creds(creds))

	trace.AuthRequest = func(r *http.Request) (bool, bool) { return true, true }
	go func() {
		r := vestigo.NewRouter()
		proxy.NewServer(r, db)
		web.Register(r, db)

		logrus.Info("server listening in :8080 to HTTP connections")
		http.Handle("/", r)
		http.ListenAndServe(":8080", nil)
	}()

	go fnctx.BgProcessor(db)

	pb.RegisterFnServer(s, api.NewServer(db))

	logrus.Info("server listening in :50050 to GRPC connections")
	s.Serve(listener)
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, appErr := handler(ctx, req)
	if appErr != nil && grpc.Code(appErr) == codes.Unknown {
		logrus.WithFields(logrus.Fields{"err": appErr}).Error(appErr.Error())

		return resp, grpc.Errorf(codes.Unknown, "internal service error: %s", appErr.Error())
	}

	return resp, appErr
}

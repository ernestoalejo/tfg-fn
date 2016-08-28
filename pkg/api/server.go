package api

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/altipla-consulting/chrono"
	pb_empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"

	"github.com/ernestoalejo/tfg-fn/pkg/kubernetes"
	"github.com/ernestoalejo/tfg-fn/pkg/models"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

type Server struct {
	db    *r.Session
	token string
}

func NewServer(db *r.Session) *Server {
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err)
	}

	return &Server{db, string(token)}
}

func (s *Server) ListFunctions(ctx context.Context, in *pb_empty.Empty) (*pb.ListFunctionsReply, error) {
	cursor, err := r.Table(models.TableFunctions).Run(s.db)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer cursor.Close()

	functions := []*models.Function{}
	if err := cursor.All(&functions); err != nil {
		return nil, errors.Trace(err)
	}

	reply := new(pb.ListFunctionsReply)
	for _, fn := range functions {
		reply.Functions = append(reply.Functions, &pb.Function{
			Name:      fn.Name,
			Call:      fn.Call,
			Trigger:   fn.Trigger,
			Method:    fn.Method,
			CreatedAt: chrono.DateTimeToProto(fn.CreatedAt),
		})
	}

	return reply, nil
}

func (s *Server) DeployFunction(ctx context.Context, in *pb.DeployFunctionRequest) (*pb_empty.Empty, error) {
	fn := &models.Function{
		Name:      in.Function.Name,
		Call:      in.Function.Call,
		Trigger:   in.Function.Trigger,
		Method:    in.Function.Method,
		CreatedAt: time.Now(),
	}
	if _, err := r.Table(models.TableFunctions).Insert(fn).RunWrite(s.db); err != nil {
		return nil, errors.Trace(err)
	}

	client := kubernetes.NewClient(s.token)

	deployment := &kubernetes.Deployment{
		APIVersion: "extensions/v1beta1",
		Metadata:   &kubernetes.ObjectMeta{Name: fn.Name},
		Spec: &kubernetes.DeploymentSpec{
			RevisionHistoryLimit: 1,
			Strategy: &kubernetes.DeploymentStrategy{
				RollingUpdate: &kubernetes.RollingUpdateDeployment{
					MaxUnavailable: 0,
				},
			},
			Template: &kubernetes.PodTemplateSpec{
				Metadata: &kubernetes.ObjectMeta{
					Labels: map[string]string{
						"app": fn.Name,
					},
				},
				Spec: &kubernetes.PodSpec{
					Containers: []*kubernetes.Container{
						{
							Name:  fn.Name,
							Image: fmt.Sprintf("localhost:5000/%s", fn.Name),
							Ports: []*kubernetes.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 50050,
								},
							},
						},
					},
				},
			},
		},
	}
	if err := client.CreateDeployment(deployment); err != nil {
		return nil, errors.Trace(err)
	}

	return new(pb_empty.Empty), nil
}

func (s *Server) DeleteFunction(ctx context.Context, in *pb.DeleteFunctionRequest) (*pb_empty.Empty, error) {
	if _, err := r.Table(models.TableFunctions).Get(in.Name).Delete().RunWrite(s.db); err != nil {
		return nil, errors.Trace(err)
	}

	client := kubernetes.NewClient(s.token)

	if err := client.DeleteDeployment(in.Name); err != nil {
		return nil, errors.Trace(err)
	}

	return new(pb_empty.Empty), nil
}

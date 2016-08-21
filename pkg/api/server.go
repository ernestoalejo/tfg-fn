package api

import (
	"time"

	"github.com/altipla-consulting/chrono"
	pb_empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"

	"github.com/ernestoalejo/tfg-fn/pkg/models"
	pb "github.com/ernestoalejo/tfg-fn/protos"
)

type Server struct {
	db *r.Session
}

func NewServer(db *r.Session) *Server {
	return &Server{db}
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
			CreatedAt: chrono.DateTimeToProto(fn.CreatedAt),
		})
	}

	return reply, nil
}

func (s *Server) DeployFunction(ctx context.Context, in *pb.DeployFunctionRequest) (*pb_empty.Empty, error) {
	fn := &models.Function{
		Name:      in.Function.Name,
		CreatedAt: time.Now(),
	}
	if _, err := r.Table(models.TableFunctions).Insert(fn).RunWrite(s.db); err != nil {
		return nil, errors.Trace(err)
	}

	return new(pb_empty.Empty), nil
}

package proxy

import (
	"fmt"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/juju/errors"
	r "gopkg.in/dancannon/gorethink.v2"

	"github.com/ernestoalejo/tfg-fn/pkg/context"
)

type Server struct {
	db *r.Session
}

func NewServer(r *vestigo.Router, db *r.Session) {
	s := &Server{db}
	r.Get("/trigger/:name", s.Trigger)
}

func (s *Server) Trigger(w http.ResponseWriter, r *http.Request) {
	resp, err := context.CallFunction(r, vestigo.Param(r, "name"))
	if err != nil {
		fmt.Fprintln(w, errors.ErrorStack(err))
	} else {
		fmt.Fprintln(w, resp)
	}
}

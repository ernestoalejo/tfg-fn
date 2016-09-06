package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, errors.ErrorStack(err))
		return
	}
	values := url.Values{}
	for k, vals := range r.Form {
		if strings.HasPrefix(k, ":") {
			continue
		}
		values[k] = vals
	}
	resp, err := context.CallFunction(vestigo.Param(r, "name"), values)
	if err != nil {
		fmt.Fprintln(w, errors.ErrorStack(err))
	} else {
		fmt.Fprintln(w, resp)
	}
}

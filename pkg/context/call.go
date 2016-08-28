package context

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/juju/errors"
)

var (
	calls    = make(chan *Call, 100)
	backends = map[string]*Backend{}
)

type Call struct {
	Name     string
	Request  *Request
	Response chan string
	Error    chan error
}

type Request struct {
	Form url.Values `json:"form"`
}

func CallFunction(r *http.Request, name string) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", errors.Trace(err)
	}

	call := &Call{
		Name: name,
		Request: &Request{
			Form: r.Form,
		},
		Response: make(chan string, 1),
		Error:    make(chan error, 1),
	}
	calls <- call
	select {
	case resp := <-call.Response:
		return resp, nil
	case err := <-call.Error:
		return "", errors.Trace(err)
	}

	panic("should not reach here")
}

func BgProcessor() {
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err)
	}

	for call := range calls {
		if backends[call.Name] == nil {
			backends[call.Name] = NewBackend(string(token), call.Name)
		}
		backends[call.Name].Process <- call
	}
}

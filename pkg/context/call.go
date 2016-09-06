package context

import (
	"io/ioutil"
	"net/url"

	"github.com/ernestoalejo/tfg-fn/pkg/models"
	"github.com/juju/errors"
	r "gopkg.in/dancannon/gorethink.v2"
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

func CallFunction(name string, form url.Values) (string, error) {
	call := &Call{
		Name: name,
		Request: &Request{
			Form: form,
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

func BgProcessor(db *r.Session) {
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err)
	}

	for call := range calls {
		if backends[call.Name] == nil {
			cursor, err := r.Table(models.TableFunctions).Get(call.Name).Run(db)
			if err != nil {
				panic(err)
			}
			fn := new(models.Function)
			err = cursor.One(fn)
			cursor.Close()
			if err != nil {
				panic(err)
			}

			backends[call.Name] = NewBackend(string(token), call.Name, fn.Call)
		}
		backends[call.Name].Process <- call
	}
}

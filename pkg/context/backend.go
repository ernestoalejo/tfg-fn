package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/ernestoalejo/tfg-fn/pkg/kubernetes"
)

type Backend struct {
	Process chan *Call

	token   string
	name    string
	closers map[string]chan bool
}

func NewBackend(token, name string) *Backend {
	backend := &Backend{
		name:  name,
		token: token,
	}
	go backend.Controller()
	return backend
}

func (backend *Backend) Controller() {
	for {
		client := kubernetes.NewClient(backend.token)
		pods, err := client.GetPods()
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error("cannot get pods")
			time.Sleep(5 * time.Second)
			continue
		}

		names := []string{}
		for _, pod := range pods {
			app := pod.Metadata.Labels["app"]
			if app != backend.name {
				continue
			}
			names = append(names, pod.Metadata.Name)
		}
		logrus.WithFields(logrus.Fields{"pods": names, "function": backend.name}).Info("pods updated")

		active := map[string]bool{}
		for _, name := range names {
			active[name] = true
			if backend.closers[name] == nil {
				backend.closers[name] = make(chan bool, 1)
				go backend.Processor(name, backend.closers[name])
				logrus.WithFields(logrus.Fields{"pod": name, "function": backend.name}).Info("enable pod")
			}
		}
		for name := range backend.closers {
			if !active[name] {
				backend.closers[name] <- true
				delete(backend.closers, name)
				logrus.WithFields(logrus.Fields{"pod": name, "function": backend.name}).Info("disable pod")
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type Context struct {
	Request *Request `json:"request"`
}

func (backend *Backend) Processor(pod string, closer chan bool) {
	for {
		select {
		case <-closer:
			return

		case call := <-backend.Process:
			c := &Context{
				Request: call.Request,
			}
			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(c); err != nil {
				call.Error <- errors.Trace(err)
				continue
			}

			req, _ := http.NewRequest("POST", fmt.Sprintf("%s:50050", pod), buf)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				call.Error <- errors.Trace(err)
				continue
			}

			contents, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				call.Error <- errors.Trace(err)
				continue
			}

			call.Response <- string(contents)
		}
	}
}

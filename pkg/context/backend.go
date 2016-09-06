package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/ernestoalejo/tfg-fn/pkg/kubernetes"
)

type Backend struct {
	Process chan *Call

	latency chan int64

	// Read only info of the backend
	call  string
	token string
	name  string

	// Owned by the controller
	closers   map[string]chan bool
	lastClose time.Time

	// Owned by the latency controller, accesed by the controller too
	mu             *sync.Mutex
	averageLatency int64
}

func NewBackend(token, name, call string) *Backend {
	backend := &Backend{
		Process: make(chan *Call, 1000),
		call:    call,
		token:   token,
		name:    name,
		closers: map[string]chan bool{},
		latency: make(chan int64, 100),
		mu:      new(sync.Mutex),
	}
	go backend.Controller()
	go backend.LatencyController()
	return backend
}

func (backend *Backend) LatencyController() {
	var numerator int64
	var values int64
	for r := range backend.latency {
		values++
		numerator += r
		backend.mu.Lock()
		backend.averageLatency = numerator / values
		backend.mu.Unlock()
	}
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

		active := map[string]bool{}
		for _, name := range names {
			active[name] = true
			if backend.closers[name] == nil {
				backend.closers[name] = make(chan bool, 1)
				go backend.Processor(name, backend.closers[name])
			}
		}
		for name := range backend.closers {
			if !active[name] {
				backend.closers[name] <- true
				delete(backend.closers, name)
			}
		}

		backend.mu.Lock()
		averageLatency := backend.averageLatency
		backend.mu.Unlock()
		pending := int64(len(backend.Process))
		current := int64(len(backend.closers))
		logrus.WithFields(logrus.Fields{
			"pods":           names,
			"function":       backend.name,
			"pending":        pending,
			"current":        current,
			"averageLatency": averageLatency,
		}).Info("controller update")
		if (current == 0 && pending > 0) || (current > 0 && pending*averageLatency/current > 150) {
			logrus.WithFields(logrus.Fields{"function": backend.name, "desired": current + 1}).Info("scale up function")
			if err := client.ScaleDeployment(backend.name, current+1); err != nil {
				logrus.WithFields(logrus.Fields{"error": err}).Error("cannot scale deployment")
			}
		} else if current > 0 && pending*averageLatency/(current-1) < 100 && time.Now().Sub(backend.lastClose) > 50*time.Second {
			backend.lastClose = time.Now()
			logrus.WithFields(logrus.Fields{"function": backend.name, "desired": current - 1}).Info("scale down function")
			if err := client.ScaleDeployment(backend.name, current-1); err != nil {
				logrus.WithFields(logrus.Fields{"error": err}).Error("cannot scale deployment")
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type Context struct {
	Call    string   `json:"call"`
	Request *Request `json:"request"`
}

func (backend *Backend) Processor(podName string, closer chan bool) {
	client := kubernetes.NewClient(backend.token)
	var podIP string
	for {
		pod, err := client.GetPod(podName)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err, "pod": podName}).Error("cannot get pod")
			time.Sleep(5 * time.Second)
			continue
		}
		if pod.Status.Phase != "Running" {
			logrus.WithFields(logrus.Fields{"pod": podName}).Error("waiting pod")
			time.Sleep(5 * time.Second)
			continue
		}
		podIP = pod.Status.PodIP
		break
	}
	logrus.WithFields(logrus.Fields{"pod": podName, "function": backend.name}).Info("enable pod")

	for {
		select {
		case <-closer:
			logrus.WithFields(logrus.Fields{"pod": podName, "function": backend.name}).Info("disable pod")
			return

		case call := <-backend.Process:
			start := time.Now()
			c := &Context{
				Call:    backend.call,
				Request: call.Request,
			}
			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(c); err != nil {
				call.Error <- errors.Trace(err)
				continue
			}

			req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:50050", podIP), buf)
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
			backend.latency <- int64(time.Now().Sub(start) / time.Millisecond)
		}
	}
}

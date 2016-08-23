package kubernetes

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juju/errors"
)

type Client struct {
	token  string
	client *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *Client) call(method, url string, request interface{}) error {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(request); err != nil {
		return errors.Trace(err)
	}
	req, _ := http.NewRequest(method, fmt.Sprintf("kubernetes%s", url), buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("bad http status: %s", resp.Status)
	}

	return nil
}

type Deployment struct {
	APIVersion string          `json:"apiVersion"`
	Metadata   *ObjectMeta     `json:"metadata"`
	Spec       *DeploymentSpec `json:"spec"`
}

type ObjectMeta struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type DeploymentSpec struct {
	Replicas             int64               `json:"replicas"`
	RevisionHistoryLimit int64               `json:"revisionHistoryLimit"`
	Template             *PodTemplateSpec    `json:"template"`
	Strategy             *DeploymentStrategy `json:"strategy"`
}

type DeploymentStrategy struct {
	Type          string                   `json:"type"`
	RollingUpdate *RollingUpdateDeployment `json:"rollingUpdate"`
}

type RollingUpdateDeployment struct {
	MaxUnavailable int64 `json:"maxUnavailable"`
}

type PodTemplateSpec struct {
	Metadata *ObjectMeta `json:"metadata"`
	Spec     *PodSpec    `json:"spec"`
}

type PodSpec struct {
	Containers []*Container `json:"containers"`
}

type Container struct {
	Name  string           `json:"name"`
	Image string           `json:"image"`
	Ports []*ContainerPort `json:"ports"`
}

type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int64  `json:"containerPort"`
}

func (c *Client) CreateDeployment(deployment *Deployment) error {
	if err := c.call("POST", "/apis/extensions/v1beta1/namespaces/default/deployments", deployment); err != nil {
		return errors.Trace(err)
	}
	return nil
}

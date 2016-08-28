package kubernetes

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (c *Client) call(method, url string, request, response interface{}) error {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(request); err != nil {
		return errors.Trace(err)
	}
	req, _ := http.NewRequest(method, fmt.Sprintf("https://kubernetes%s", url), buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Trace(err)
		}
		return errors.Errorf("bad http status: %s: %s", resp.Status, content)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		if err == io.EOF {
			return nil
		}

		return errors.Trace(err)
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
	return errors.Trace(c.call("POST", "/apis/extensions/v1beta1/namespaces/default/deployments", deployment, nil))
}

func (c *Client) DeleteDeployment(name string) error {
	return errors.Trace(c.call("DELETE", fmt.Sprintf("/apis/extensions/v1beta1/namespaces/default/deployments/%s", name), nil, nil))
}

type Service struct {
	APIVersion string       `json:"apiVersion"`
	Metadata   *ObjectMeta  `json:"metadata"`
	Spec       *ServiceSpec `json:"spec"`
}

type ServiceSpec struct {
	Ports    []*ServicePort    `json:"ports"`
	Selector map[string]string `json:"selector"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int64  `json:"port"`
	TargetPort int64  `json:"targetPort"`
}

func (c *Client) CreateService(service *Service) error {
	return errors.Trace(c.call("POST", "/api/v1/namespaces/default/services", service, nil))
}

func (c *Client) DeleteService(name string) error {
	return errors.Trace(c.call("DELETE", fmt.Sprintf("/api/v1/namespaces/default/services/%s", name), nil, nil))
}

type PodList struct {
	Items []*Pod `json:"items"`
}

type Pod struct {
	Metadata *ObjectMeta `json:"metadata"`
}

func (c *Client) GetPods() ([]*Pod, error) {
	resp := new(PodList)
	if err := c.call("GET", "/api/v1/pods", nil, resp); err != nil {
		return nil, errors.Trace(err)
	}

	return resp.Items, nil
}

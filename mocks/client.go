package mocks

import (
	"fmt"
	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	v1beta1extensions "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

type Client struct {
	v1core.PodInterface
	v1core.ServiceInterface
	v1beta1extensions.DaemonSetInterface
	v1core.EndpointsInterface
	v1batch.JobInterface
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.PodInterface
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.ServiceInterface
}

func (c Client) DaemonSets(namespace string) v1beta1extensions.DaemonSetInterface {
	return c.DaemonSetInterface
}

func (c Client) Endpoints(namespace string) v1core.EndpointsInterface {
	return c.EndpointsInterface
}
func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.JobInterface
}

func (c Client) GetResourceName(kind, apiVersion string) (string, error) {
	if kind == "failKind" {
		return "failResourceName", fmt.Errorf("Expected Failure")
	}
	return "successResourceName", nil
}

func (c Client) CustomResource(apiVersion, namespace, resource, name string) (map[string]interface{}, error) {
	if name == "failName" {
		return nil, fmt.Errorf("Expected Failure")
	} else if name == "resolved" {
		return map[string]interface{}{
			"simple_key": "simple_value",
			"complex": map[string]interface{}{
				"key": map[string]interface{}{
					"with": map[string]interface{}{
						"layers": "complex_value",
					},
				},
			},
		}, nil
	}
	return map[string]interface{}{
		"key": "unexpected_value",
	}, nil
}

func NewClient() cli.ClientInterface {
	return Client{
		NewPClient(),
		NewSClient(),
		NewDSClient(),
		NewEClient(),
		NewJClient(),
	}
}

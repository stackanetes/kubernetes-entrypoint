package mocks

import (
	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

type Client struct {
	unversioned.PodInterface
	unversioned.ServiceInterface
	unversioned.DaemonSetInterface
	unversioned.EndpointsInterface
	unversioned.JobInterface
}

func (c Client) Pods(namespace string) unversioned.PodInterface {
	return c.PodInterface
}

func (c Client) Services(namespace string) unversioned.ServiceInterface {
	return c.ServiceInterface
}

func (c Client) DaemonSets(namespace string) unversioned.DaemonSetInterface {
	return c.DaemonSetInterface
}

func (c Client) Endpoints(namespace string) unversioned.EndpointsInterface {
	return c.EndpointsInterface
}
func (c Client) Jobs(namespace string) unversioned.JobInterface {
	return c.JobInterface
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

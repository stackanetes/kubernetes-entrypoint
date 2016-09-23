package client

import (
	restclient "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

type ClientInterface interface {
	Pods(string) unversioned.PodInterface
	Jobs(string) unversioned.JobInterface
	Endpoints(string) unversioned.EndpointsInterface
	DaemonSets(string) unversioned.DaemonSetInterface
	Services(string) unversioned.ServiceInterface
}

type Client struct {
	*unversioned.Client
}

func (c Client) Pods(namespace string) unversioned.PodInterface {
	return c.Client.Pods(namespace)
}

func (c Client) Jobs(namespace string) unversioned.JobInterface {
	return c.Client.Extensions().Jobs(namespace)
}

func (c Client) Endpoints(namespace string) unversioned.EndpointsInterface {
	return c.Client.Endpoints(namespace)
}
func (c Client) DaemonSets(namespace string) unversioned.DaemonSetInterface {
	return c.Client.Extensions().DaemonSets(namespace)
}

func (c Client) Services(namespace string) unversioned.ServiceInterface {
	return c.Client.Services(namespace)
}

func New(config *restclient.Config) (ClientInterface, error) {
	if config == nil {
		client, err := unversioned.NewInCluster()
		if err != nil {
			return nil, err
		}
		return Client{Client: client}, nil
	}

	client, err := unversioned.New(config)
	if err != nil {
		return nil, err
	}

	return Client{Client: client}, nil

}

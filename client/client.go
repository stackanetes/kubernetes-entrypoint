package client

import (
	"k8s.io/client-go/kubernetes"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	v1beta1extensions "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

type ClientInterface interface {
	Pods(string) v1core.PodInterface
	Jobs(string) v1batch.JobInterface
	Endpoints(string) v1core.EndpointsInterface
	DaemonSets(string) v1beta1extensions.DaemonSetInterface
	Services(string) v1core.ServiceInterface
}
type Client struct {
	*kubernetes.Clientset
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.Clientset.Core().Pods(namespace)
}

func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.Clientset.Batch().Jobs(namespace)
}

func (c Client) Endpoints(namespace string) v1core.EndpointsInterface {
	return c.Clientset.Core().Endpoints(namespace)
}
func (c Client) DaemonSets(namespace string) v1beta1extensions.DaemonSetInterface {
	return c.Clientset.Extensions().DaemonSets(namespace)
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.Clientset.Core().Services(namespace)
}

func New(config *rest.Config) (ClientInterface, error) {
	if config == nil {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		return Client{Clientset: clientset}, nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return Client{Clientset: clientset}, nil

}

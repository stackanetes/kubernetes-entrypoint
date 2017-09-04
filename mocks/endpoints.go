package mocks

import (
	"fmt"

	apicore "k8s.io/client-go/1.5/kubernetes/typed/core/v1"
	api "k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/watch"
	"k8s.io/client-go/1.5/rest"
)

type eClient struct {
}

const (
	MockEndpointError = "Mock endpoint didnt work"
)

func (e eClient) Get(name string) (*v1.Endpoints, error) {
	if name == FailingServiceName {
		return nil, fmt.Errorf(MockEndpointError)
	}

	subsets := []v1.EndpointSubset{}

	if name != EmptySubsetsServiceName {
		subsets = []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{IP: "127.0.0.1"},
				},
			},
		}
	}

	endpoint := &v1.Endpoints{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Subsets:    subsets,
	}

	return endpoint, nil
}
func (e eClient) Create(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Delete(name string, options *api.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}

func (e eClient) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (e eClient) List(options api.ListOptions) (*v1.EndpointsList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Update(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s eClient) UpdateStatus(ds *api.Endpoints) (*api.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	return nil
}

func (e eClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Endpoints, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewEClient() apicore.EndpointsInterface {
	return eClient{}
}

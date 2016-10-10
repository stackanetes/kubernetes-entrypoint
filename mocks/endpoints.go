package mocks

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type eClient struct {
}

func (e eClient) Get(name string) (*api.Endpoints, error) {
	if name != "lgtm" {
		return nil, fmt.Errorf("Mock endpoint didnt work")
	}
	endpoint := &api.Endpoints{
		ObjectMeta: api.ObjectMeta{Name: name},
		Subsets: []api.EndpointSubset{
			api.EndpointSubset{
				Addresses: []api.EndpointAddress{
					api.EndpointAddress{IP: "127.0.0.1"},
				},
			},
		},
	}
	return endpoint, nil
}
func (e eClient) Create(ds *api.Endpoints) (*api.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Delete(name string) error {
	return fmt.Errorf("Not implemented")
}
func (e eClient) List(options api.ListOptions) (*api.EndpointsList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Update(ds *api.Endpoints) (*api.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s eClient) UpdateStatus(ds *api.Endpoints) (*api.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) restclient.ResponseWrapper {
	return nil
}

func NewEClient() unversioned.EndpointsInterface {
	return eClient{}
}

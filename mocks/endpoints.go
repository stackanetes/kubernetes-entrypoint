package mocks

import (
	"fmt"

	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	api "k8s.io/client-go/pkg/api"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

type eClient struct {
}

func (e eClient) Get(name string) (*v1.Endpoints, error) {
	if name != "lgtm" {
		return nil, fmt.Errorf("Mock endpoint didnt work")
	}
	endpoint := &v1.Endpoints{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{IP: "127.0.0.1"},
				},
			},
		},
	}
	return endpoint, nil
}
func (e eClient) Create(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Delete(name string, options *v1.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}

func (e eClient) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (e eClient) List(options v1.ListOptions) (*v1.EndpointsList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Update(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s eClient) UpdateStatus(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Watch(options v1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	return nil
}

func (e eClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Endpoints, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewEClient() v1core.EndpointsInterface {
	return eClient{}
}

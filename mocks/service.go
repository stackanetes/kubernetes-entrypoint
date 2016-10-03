package mocks

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type sClient struct {
}

func (s sClient) Get(name string) (*api.Service, error) {
	if name != "lgtm" {
		return nil, fmt.Errorf("Mock service didnt work")
	}
	return &api.Service{
		ObjectMeta: api.ObjectMeta{Name: name},
	}, nil
}
func (s sClient) Create(ds *api.Service) (*api.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Delete(name string) error {
	return fmt.Errorf("Not implemented")
}
func (s sClient) List(options api.ListOptions) (*api.ServiceList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Update(ds *api.Service) (*api.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) UpdateStatus(ds *api.Service) (*api.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) restclient.ResponseWrapper {
	return nil
}

func NewSClient() unversioned.ServiceInterface {
	return sClient{}
}

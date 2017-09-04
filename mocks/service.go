package mocks

import (
	"fmt"

	v1core "k8s.io/client-go/1.5/kubernetes/typed/core/v1"
	api "k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/watch"
	"k8s.io/client-go/1.5/rest"
)

type sClient struct {
}

const (
	MockServiceError        = "Mock service didnt work"
	SucceedingServiceName   = "succeed"
	EmptySubsetsServiceName = "empty-subsets"
	FailingServiceName      = "fail"
)

func (s sClient) Get(name string) (*v1.Service, error) {
	if name == FailingServiceName {
		return nil, fmt.Errorf(MockServiceError)
	}
	return &v1.Service{
		ObjectMeta: v1.ObjectMeta{Name: name},
	}, nil
}
func (s sClient) Create(ds *v1.Service) (*v1.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Delete(name string, options *api.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}

func (s sClient) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (s sClient) List(options api.ListOptions) (*v1.ServiceList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Update(ds *v1.Service) (*v1.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) UpdateStatus(ds *v1.Service) (*v1.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	return nil
}

func (s sClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Service, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewSClient() v1core.ServiceInterface {
	return sClient{}
}

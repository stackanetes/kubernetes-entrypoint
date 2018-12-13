package mocks

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type sClient struct {
}

const (
	MockServiceError        = "Mock service didnt work"
	SucceedingServiceName   = "succeed"
	EmptySubsetsServiceName = "empty-subsets"
	FailingServiceName      = "fail"
)

func (s sClient) Get(name string, opts metav1.GetOptions) (*v1.Service, error) {
	if name == FailingServiceName {
		return nil, fmt.Errorf(MockServiceError)
	}
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}, nil
}
func (s sClient) Create(ds *v1.Service) (*v1.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Delete(name string, options *metav1.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}

func (s sClient) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (s sClient) List(options metav1.ListOptions) (*v1.ServiceList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Update(ds *v1.Service) (*v1.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) UpdateStatus(ds *v1.Service) (*v1.Service, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s sClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	return nil
}

func (s sClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Service, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewSClient() v1core.ServiceInterface {
	return sClient{}
}

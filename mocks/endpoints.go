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

type eClient struct {
}

const (
	MockEndpointError = "Mock endpoint didnt work"
)

func (e eClient) Get(name string, opts metav1.GetOptions) (*v1.Endpoints, error) {
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
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Subsets:    subsets,
	}

	return endpoint, nil
}
func (e eClient) Create(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Delete(name string, options *metav1.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}

func (e eClient) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (e eClient) List(options metav1.ListOptions) (*v1.EndpointsList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Update(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s eClient) UpdateStatus(ds *v1.Endpoints) (*v1.Endpoints, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (e eClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	return nil
}

func (e eClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Endpoints, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewEClient() v1core.EndpointsInterface {
	return eClient{}
}

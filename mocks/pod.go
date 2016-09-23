package mocks

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type pClient struct {
}

func (p pClient) Get(name string) (*api.Pod, error) {
	if name == "lgtm" {
		pod := new(api.Pod)
		container_one := api.ContainerStatus{
			Name:  "container_test",
			Ready: true,
		}
		pod.Status.ContainerStatuses = []api.ContainerStatus{container_one}
		return pod, nil
	}
	return nil, fmt.Errorf("Mock pod didnt work")
}
func (p pClient) Create(pod *api.Pod) (*api.Pod, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (p pClient) Delete(name string, options *api.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (p pClient) List(options api.ListOptions) (*api.PodList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (p pClient) Update(pod *api.Pod) (*api.Pod, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (p pClient) UpdateStatus(pod *api.Pod) (*api.Pod, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (p pClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (p pClient) GetLogs(name string, opts *api.PodLogOptions) *restclient.Request {
	return nil
}
func (p pClient) Bind(binding *api.Binding) error {
	return fmt.Errorf("Not implemented")
}
func NewPClient() unversioned.PodInterface {
	return pClient{}
}

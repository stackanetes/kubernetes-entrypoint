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
	return &api.Pod{
		ObjectMeta: api.ObjectMeta{Name: name},
		Status: api.PodStatus{
			ContainerStatuses: []api.ContainerStatus{
				api.ContainerStatus{
					Name:  "container_test",
					Ready: true,
				},
			},
			HostIP: "127.0.0.1",
		},
	}, nil

}
func (p pClient) Create(pod *api.Pod) (*api.Pod, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (p pClient) Delete(name string, options *api.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (p pClient) List(options api.ListOptions) (*api.PodList, error) {
	return &api.PodList{
		Items: []api.Pod{
			api.Pod{
				ObjectMeta: api.ObjectMeta{Name: "podList"},
				Status: api.PodStatus{
					HostIP: "127.0.01",
					Conditions: []api.PodCondition{
						api.PodCondition{
							Type:   api.PodReady,
							Status: "True",
						},
					},
					ContainerStatuses: []api.ContainerStatus{
						api.ContainerStatus{
							Name:  "container_test",
							Ready: true,
						},
					},
				},
			},
		},
	}, nil

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

package mocks

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	unv "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type dClient struct {
}

func (d dClient) Get(name string) (*extensions.DaemonSet, error) {
	if name != "lgtm" {
		return nil, fmt.Errorf("Mock daemonset didnt work")
	}
	ds := &extensions.DaemonSet{
		ObjectMeta: api.ObjectMeta{Name: name},
		Spec: extensions.DaemonSetSpec{
			Selector: &unv.LabelSelector{
				MatchLabels: map[string]string{"name": "test"},
			},
		},
	}
	return ds, nil
}
func (d dClient) Create(ds *extensions.DaemonSet) (*extensions.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Delete(name string) error {
	return fmt.Errorf("Not implemented")
}
func (d dClient) List(options api.ListOptions) (*extensions.DaemonSetList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Update(ds *extensions.DaemonSet) (*extensions.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) UpdateStatus(ds *extensions.DaemonSet) (*extensions.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewDSClient() unversioned.DaemonSetInterface {
	return dClient{}
}

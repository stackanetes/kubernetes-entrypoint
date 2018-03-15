package mocks

import (
	"fmt"

	v1beta1extensions "k8s.io/client-go/1.5/kubernetes/typed/extensions/v1beta1"
	api "k8s.io/client-go/1.5/pkg/api"
	v1 "k8s.io/client-go/1.5/pkg/api/v1"
	extensions "k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/1.5/pkg/watch"
)

type dClient struct {
}

const (
	SucceedingDaemonsetName         = "DAEMONSET_SUCCEED"
	FailingDaemonsetName            = "DAEMONSET_FAIL"
	CorrectNamespaceDaemonsetName   = "CORRECT_DAEMONSET_NAMESPACE"
	IncorrectNamespaceDaemonsetName = "INCORRECT_DAEMONSET_NAMESPACE"
	CorrectDaemonsetNamespace       = "CORRECT_DAEMONSET"

	FailingMatchLabelsDaemonsetName  = "DAEMONSET_INCORRECT_MATCH_LABELS"
	NotReadyMatchLabelsDaemonsetName = "DAEMONSET_NOT_READY_MATCH_LABELS"
)

func (d dClient) Get(name string) (*extensions.DaemonSet, error) {
	matchLabelName := MockContainerName

	if name == FailingDaemonsetName {
		return nil, fmt.Errorf("Mock daemonset didnt work")
	} else if name == FailingMatchLabelsDaemonsetName {
		matchLabelName = FailingMatchLabel
	} else if name == NotReadyMatchLabelsDaemonsetName {
		matchLabelName = SameHostNotReadyMatchLabel
	}

	ds := &extensions.DaemonSet{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Spec: extensions.DaemonSetSpec{
			Selector: &extensions.LabelSelector{
				MatchLabels: map[string]string{"name": matchLabelName},
			},
		},
	}

	if name == CorrectNamespaceDaemonsetName {
		ds.ObjectMeta.Namespace = CorrectDaemonsetNamespace
	} else if name == IncorrectNamespaceDaemonsetName {
		return nil, fmt.Errorf("Mock daemonset didnt work")
	}

	return ds, nil
}
func (d dClient) Create(ds *extensions.DaemonSet) (*extensions.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Delete(name string, options *api.DeleteOptions) error {
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

func (d dClient) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (d dClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *extensions.DaemonSet, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewDSClient() v1beta1extensions.DaemonSetInterface {
	return dClient{}
}

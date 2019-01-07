package mocks

import (
	"fmt"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1beta1extensions "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
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

func (d dClient) Get(name string, opts metav1.GetOptions) (*v1beta1.DaemonSet, error) {
	matchLabelName := MockContainerName

	if name == FailingDaemonsetName {
		return nil, fmt.Errorf("Mock daemonset didnt work")
	} else if name == FailingMatchLabelsDaemonsetName {
		matchLabelName = FailingMatchLabel
	} else if name == NotReadyMatchLabelsDaemonsetName {
		matchLabelName = SameHostNotReadyMatchLabel
	}

	ds := &v1beta1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: v1beta1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
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
func (d dClient) Create(ds *v1beta1.DaemonSet) (*v1beta1.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Delete(name string, options *metav1.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (d dClient) List(options metav1.ListOptions) (*v1beta1.DaemonSetList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Update(ds *v1beta1.DaemonSet) (*v1beta1.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) UpdateStatus(ds *v1beta1.DaemonSet) (*v1beta1.DaemonSet, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return fmt.Errorf("Not implemented")
}

func (d dClient) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d dClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.DaemonSet, err error) {
	return nil, fmt.Errorf("Not implemented")
}

func NewDSClient() v1beta1extensions.DaemonSetInterface {
	return dClient{}
}

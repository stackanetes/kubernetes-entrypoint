package mocks

import (
	"fmt"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
)

const (
	SucceedingJobName  = "succeed"
	FailingJobName     = "fail"
	SucceedingJobLabel = "succeed"
	FailingJobLabel    = "fail"
)

type jClient struct {
}

func (j jClient) Get(name string, opts metav1.GetOptions) (*v1.Job, error) {
	if name == SucceedingJobName {
		return &v1.Job{
			Status: v1.JobStatus{Succeeded: 1},
		}, nil
	}
	if name == FailingJobName {
		return &v1.Job{
			Status: v1.JobStatus{Succeeded: 0},
		}, nil
	}
	return nil, fmt.Errorf("Mock job didnt work")
}
func (j jClient) Create(job *v1.Job) (*v1.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Delete(name string, opts *metav1.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) List(options metav1.ListOptions) (*v1.JobList, error) {
	var jobs []v1.Job
	if options.LabelSelector == fmt.Sprintf("name=%s", SucceedingJobLabel) {
		jobs = []v1.Job{NewJob(1)}
	} else if options.LabelSelector == fmt.Sprintf("name=%s", FailingJobLabel) {
		jobs = []v1.Job{NewJob(1), NewJob(0)}
	} else {
		return nil, fmt.Errorf("Mock job didnt work")
	}
	return &v1.JobList{
		Items: jobs,
	}, nil
}

func (j jClient) Update(job *v1.Job) (*v1.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) UpdateStatus(job *v1.Job) (*v1.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Job, err error) {
	return nil, fmt.Errorf("Not implemented")
}
func NewJClient() v1batch.JobInterface {
	return jClient{}
}

func NewJob(succeeded int32) v1.Job {
	return v1.Job{
		Status: v1.JobStatus{Succeeded: succeeded},
	}
}

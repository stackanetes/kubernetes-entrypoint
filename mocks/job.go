package mocks

import (
	"fmt"

	v1batch "k8s.io/client-go/1.5/kubernetes/typed/batch/v1"
	api "k8s.io/client-go/1.5/pkg/api"
	batch "k8s.io/client-go/1.5/pkg/apis/batch/v1"
	"k8s.io/client-go/1.5/pkg/watch"
)

const (
	SucceedingJobName  = "succeed"
	FailingJobName     = "fail"
	SucceedingJobLabel = "succeed"
	FailingJobLabel    = "fail"
)

type jClient struct {
}

func (j jClient) Get(name string) (*batch.Job, error) {
	if name == SucceedingJobName {
		return &batch.Job{
			Status: batch.JobStatus{Succeeded: 1},
		}, nil
	}
	if name == FailingJobName {
		return &batch.Job{
			Status: batch.JobStatus{Succeeded: 0},
		}, nil
	}
	return nil, fmt.Errorf("Mock job didnt work")
}
func (j jClient) Create(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Delete(name string, opts *api.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) List(options api.ListOptions) (*batch.JobList, error) {
	var jobs []batch.Job
	if options.LabelSelector.String() == fmt.Sprintf("name=%s", SucceedingJobLabel) {
		jobs = []batch.Job{NewJob(1)}
	} else if options.LabelSelector.String() == fmt.Sprintf("name=%s", FailingJobLabel) {
		jobs = []batch.Job{NewJob(1), NewJob(0)}
	} else {
		return nil, fmt.Errorf("Mock job didnt work")
	}
	return &batch.JobList{
		Items: jobs,
	}, nil
}

func (j jClient) Update(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) UpdateStatus(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Watch(options api.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *batch.Job, err error) {
	return nil, fmt.Errorf("Not implemented")
}
func NewJClient() v1batch.JobInterface {
	return jClient{}
}

func NewJob(succeeded int32) batch.Job {
	return batch.Job{
		Status: batch.JobStatus{Succeeded: succeeded},
	}
}

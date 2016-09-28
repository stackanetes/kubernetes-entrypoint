package mocks

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type jClient struct {
}

func (j jClient) Get(name string) (*batch.Job, error) {
	if name == "lgtm" {
		job := new(batch.Job)
		job.Status.Succeeded = 1
		return job, nil
	}
	if name == "fail" {
		job := new(batch.Job)
		job.Status.Succeeded = 0
		return job, nil
	}
	return nil, fmt.Errorf("Mock job didnt work")
}
func (j jClient) Create(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Delete(name string, opts *api.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) List(options api.ListOptions) (*batch.JobList, error) {
	return nil, fmt.Errorf("Not implemented")
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

func NewJClient() unversioned.JobInterface {
	return jClient{}
}

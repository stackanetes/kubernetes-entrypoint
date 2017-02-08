package mocks

import (
	"fmt"

	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	api "k8s.io/client-go/pkg/api"
	v1 "k8s.io/client-go/pkg/api/v1"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/pkg/watch"
)

type jClient struct {
}

func (j jClient) Get(name string) (*batch.Job, error) {
	if name == "lgtm" {
		return &batch.Job{
			Status: batch.JobStatus{Succeeded: 1},
		}, nil
	}
	if name == "fail" {
		return &batch.Job{
			Status: batch.JobStatus{Succeeded: 0},
		}, nil
	}
	return nil, fmt.Errorf("Mock job didnt work")
}
func (j jClient) Create(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Delete(name string, opts *v1.DeleteOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return fmt.Errorf("Not implemented")
}
func (j jClient) List(options v1.ListOptions) (*batch.JobList, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Update(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) UpdateStatus(job *batch.Job) (*batch.Job, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Watch(options v1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (j jClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *batch.Job, err error) {
	return nil, fmt.Errorf("Not implemented")
}
func NewJClient() v1batch.JobInterface {
	return jClient{}
}

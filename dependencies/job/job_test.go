package job

import (
	//entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"fmt"
	mocks "github.com/stackanetes/kubernetes-entrypoint/mocks"
	//batch "k8s.io/kubernetes/pkg/apis/batch"
	//cl "k8s.io/kubernetes/pkg/client/unversioned"
	"testing"
)

func TestResolveNewJob(t *testing.T) {

	var entrypoint mocks.MockEntrypoint
	job := NewJob("test")
	status, err := job.IsResolved(entrypoint)
	fmt.Printf("%v, %v, %v", job, status, err)
}

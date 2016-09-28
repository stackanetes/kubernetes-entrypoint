package job

import (
	mocks "github.com/stackanetes/kubernetes-entrypoint/mocks"
	"testing"
)

func TestResolveNewJob(t *testing.T) {

	entrypoint := mocks.NewEntrypoint()
	job := NewJob("lgtm")
	status, err := job.IsResolved(entrypoint)
	if status != true {
		t.Errorf("Resolving job failed: %v", err)
	}
}

func TestFailResolveNewJob(t *testing.T) {

	entrypoint := mocks.NewEntrypoint()
	job := NewJob("fail")
	_, err := job.IsResolved(entrypoint)
	expectedError := "Job fail is not completed yet"
	if err.Error() != expectedError {
		t.Errorf("Something went wrong: %v", err)
	}
}

package service

import (
	"github.com/stackanetes/kubernetes-entrypoint/mocks"
	"testing"
)

func TestResolveService(t *testing.T) {
	entrypoint := mocks.NewEntrypoint()
	s := NewService("lgtm")
	_, err := s.IsResolved(entrypoint)
	if err != nil {
		t.Errorf("Checking condition fail with: %v", err)
	}

}

func TestResolveServiceFail(t *testing.T) {
	entrypoint := mocks.NewEntrypoint()
	s := NewService("fail")
	_, err := s.IsResolved(entrypoint)
	expectedError := "Mock endpoint didnt work"
	if err.Error() != expectedError {
		t.Errorf("Checking condition fail with: %v", err)
	}
}

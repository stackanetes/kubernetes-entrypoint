package container

import (
	"os"
	"testing"

	"github.com/stackanetes/kubernetes-entrypoint/mocks"
)

func TestResolveContainer(t *testing.T) {
	entrypoint := mocks.NewEntrypoint()
	c := NewContainer("container_test")
	os.Setenv("POD_NAME", "lgtm")
	_, err := c.IsResolved(entrypoint)
	if err != nil {
		t.Errorf("Resolving container failed: %v", err)
	}
}

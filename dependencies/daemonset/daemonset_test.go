package daemonset

import (
	"os"
	"testing"

	mocks "github.com/stackanetes/kubernetes-entrypoint/mocks"
)

func TestResolveDaemonset(t *testing.T) {
	entrypoint := mocks.NewEntrypoint()
	daemonset := NewDaemonset("lgtm")
	os.Setenv("POD_NAME", "podlist")
	status, err := daemonset.IsResolved(entrypoint)
	if err != nil {
		t.Errorf("Something went wrong status: %s : %v", status, err)
	}

}

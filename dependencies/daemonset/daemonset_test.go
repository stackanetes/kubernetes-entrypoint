package daemonset

import (
	"os"
	"testing"

	mocks "github.com/stackanetes/kubernetes-entrypoint/mocks"
)

func TestResolveDaemonset(t *testing.T) {
	entrypoint := mocks.NewEntrypoint()
	os.Setenv("POD_NAME", "podlist")
	daemonset, err := NewDaemonset("lgtm")
	if err != nil {
		t.Errorf("Cannot initialize daemonset: %v", err)
	}
	status, err := daemonset.IsResolved(entrypoint)
	if err != nil {
		t.Errorf("Something went wrong status: %s : %v", status, err)
	}

}

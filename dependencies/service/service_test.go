package service

import (
	"fmt"
	"os"
	"testing"

	entry "github.com/stackanetes/docker-entrypoint/dependencies"
	"github.com/stackanetes/docker-entrypoint/logger"
)

func init() {
	os.Setenv(fmt.Sprintf("%sSERVICE", entry.DependencyPrefix), "test")

}
func TestRegisterNewService(t *testing.T) {
	logger.Info.Printf("%v", os.Getenv(fmt.Sprintf("%sSERVICE", entry.DependencyPrefix)))
	if len(entry.Dependencies) != 1 {
		t.Errorf("Expecting len of dependencies to be 1 not %v", len(entry.Dependencies))
	}

	if entry.Dependencies[0].GetName() != "test" {
		t.Errorf("Expecting name to be test not %s", entry.Dependencies[0].GetName())
	}
}

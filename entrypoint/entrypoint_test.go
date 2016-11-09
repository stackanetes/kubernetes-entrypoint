package entrypoint

import (
	"testing"
)

type dummyResolver struct {
	name      string
	namespace string
}

func (d dummyResolver) IsResolved(entry EntrypointInterface) (bool, error) {
	return true, nil
}
func (d dummyResolver) GetName() (name string) {
	return d.name
}

func TestRegisterNewDependency(t *testing.T) {
	dummy := dummyResolver{name: "dummy"}
	Register(dummy)
	if len(dependencies) != 1 {
		t.Errorf("Expecting dependencies len to be 1 got %v", len(dependencies))
	}
}

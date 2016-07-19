package entrypoint

import "testing"

type dummyResolver struct {
}

func (d *dummyResolver) IsResolved(name string) (bool, error) {
	return true, nil
}
func TestRegisterNewDependency(t *testing.T) {
	dummy := new(dummyResolver)
	Register(dummy)
	if len(Dependencies) != 1 {
		t.Errorf("Expecting dependencies len to be 1 got %v", len(Dependencies))
	}
}

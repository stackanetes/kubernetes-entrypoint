package command

import "testing"

func TestExecuteCommandSuccess(t *testing.T) {
	successCommand := []string{"echo", "test"}
	err := Execute(successCommand)
	if err != nil {
		t.Errorf("Expecting: command to success not %v", err)
	}
}

func TestExecuteCommandFail(t *testing.T) {
	errorCommand := []string{"false"}
	err := Execute(errorCommand)
	if err == nil {
		t.Errorf("Expecting command to fail")
	}
}

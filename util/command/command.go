package command

import (
	"os"
	"os/exec"
)

func ExecuteCommand(command []string) error {
	path, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}
	cmd := exec.Cmd{
		Path:   path,
		Args:   command,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

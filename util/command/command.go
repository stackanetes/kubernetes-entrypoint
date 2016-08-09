package command

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/stackanetes/kubernetes-entrypoint/logger"
)

func ExecuteCommand(command []string) {
	path, err := exec.LookPath(command[0])
	if err != nil {
		logger.Error.Printf("Cannot find a binary %v : %v", command[0], err)
		os.Exit(1)
	}

	env := os.Environ()
	err = syscall.Exec(path, command, env)
	if err != nil {
		logger.Error.Print("Executing command %v failed: %v", command, err)
		os.Exit(1)
	}

}

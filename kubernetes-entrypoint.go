package main

import (
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"

	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/config"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/container"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/customresource"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/daemonset"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/job"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/pod"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/service"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/socket"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	command "github.com/stackanetes/kubernetes-entrypoint/util/command"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

func main() {
	var comm []string
	var entrypoint *entry.Entrypoint
	var err error
	if entrypoint, err = entry.New(nil); err != nil {
		logger.Error.Printf("Creating entrypoint failed: %v", err)
		os.Exit(1)
	}

	entrypoint.Resolve()

	if comm = env.SplitCommand(); len(comm) == 0 {
		// TODO(DTadrzak): we should consider other options to handle whether pod
		// is an init-container
		logger.Warning.Printf("COMMAND env is empty")
		os.Exit(0)
	}

	if err = command.Execute(comm); err != nil {
		logger.Error.Printf("Cannot execute command: %v", err)
		os.Exit(1)
	}
}

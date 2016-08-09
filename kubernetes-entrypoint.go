package main

import (
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/dependencies"

	"github.com/stackanetes/kubernetes-entrypoint/logger"
	comm "github.com/stackanetes/kubernetes-entrypoint/util/command"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	cl "k8s.io/kubernetes/pkg/client/unversioned"
	//Register resolvers
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/config"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/container"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/daemonset"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/job"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/service"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/socket"
)

func main() {
	var client *cl.Client
	var command []string
	var entrypoint *entry.Entrypoint
	var err error
	if entrypoint, err = entry.NewEntrypoint(client); err != nil {
		logger.Error.Printf("Creating entrypoint failed: %v", err)
		os.Exit(1)
	}
	entrypoint.Resolve()

	if command = env.SplitEnvToList("COMMAND", " "); len(command) == 0 {
		logger.Error.Printf("COMMAND env is empty")
		os.Exit(1)
	}
	if err = comm.ExecuteCommand(command); err != nil {
		logger.Error.Printf("Executing command failed: %v", err)
		os.Exit(1)
	}
}

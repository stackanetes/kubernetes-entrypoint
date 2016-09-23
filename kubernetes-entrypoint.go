package main

import (
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"

	"github.com/stackanetes/kubernetes-entrypoint/logger"
	command "github.com/stackanetes/kubernetes-entrypoint/util/command"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	//restclient "k8s.io/kubernetes/pkg/client/restclient"
	//Register resolvers
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/config"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/container"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/daemonset"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/job"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/service"
	_ "github.com/stackanetes/kubernetes-entrypoint/dependencies/socket"
)

func main() {
	//var client cli.ClientInterface
	var comm []string
	var entrypoint *entry.Entrypoint
	var err error
	if entrypoint, err = entry.New(nil); err != nil {
		logger.Error.Printf("Creating entrypoint failed: %v", err)
		os.Exit(1)
	}
	entrypoint.Resolve()

	if comm = env.SplitEnvToList("COMMAND", " "); len(comm) == 0 {
		logger.Error.Printf("COMMAND env is empty")
		os.Exit(1)

	}
	err = command.Execute(comm)
	if err != nil {
		logger.Error.Printf("Cannot execute command: %v", err)
		os.Exit(1)
	}
}

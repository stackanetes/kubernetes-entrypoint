// Copyright 2018 The kubernetes-entrypoint Authors
// Copyright 2026 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"slices"

	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/config"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/container"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/customresource"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/daemonset"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/deployment"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/job"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/pod"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/service"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/socket"
	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	command "opendev.org/airship/kubernetes-entrypoint/util/command"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

func getNoOpCommands() []string {
	return []string{"echo", "true"}
}

func main() {
	// Placeholder for version flag handling. Needed by smoke test in image build.
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		os.Exit(0)
	}

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

	if len(comm) == 1 && slices.Index(getNoOpCommands(), comm[0]) != -1 {
		logger.Info.Printf("Early exit thanks to COMMAND being set to \"%v\"", comm[0])
		os.Exit(0)
	}

	if err = command.Execute(comm); err != nil {
		logger.Error.Printf("Cannot execute command: %v", err)
		os.Exit(1)
	}
}

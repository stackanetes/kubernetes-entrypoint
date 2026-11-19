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

package env

import (
	"encoding/json"
	"os"
	"strings"

	"opendev.org/airship/kubernetes-entrypoint/logger"
)

const (
	Separator = ":"
)

type Dependency struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type PodDependency struct {
	Labels          map[string]string `json:"labels"`
	Namespace       string            `json:"namespace"`
	RequireSameNode bool              `json:"requireSameNode"`
}

type JobDependency struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Namespace string            `json:"namespace"`
}

type DeploymentDependency struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Annotations map[string]string `json:"annotations"`
}

type CustomResourceDependency struct {
	APIVersion string              `json:"apiVersion"`
	Name       string              `json:"name"`
	Namespace  string              `json:"namespace"`
	Kind       string              `json:"kind"`
	Fields     []map[string]string `json:"fields"`
}

func SplitCommand() []string {
	command := os.Getenv("COMMAND")
	if command == "" {
		return []string{}
	}
	commandList := strings.Split(command, " ")
	return commandList
}

// SplitEnvToDeps returns list of namespaces and names pairs
func SplitEnvToDeps(env string) (envList []Dependency) {
	separator := ","

	e := os.Getenv(env)
	if e == "" {
		return envList
	}

	envVars := strings.Split(e, separator)
	namespace := GetBaseNamespace()
	var dep Dependency
	for _, envVar := range envVars {
		if strings.Contains(envVar, Separator) {
			nameAfterSplit := strings.Split(envVar, Separator)
			if len(nameAfterSplit) != 2 {
				logger.Warning.Printf("Invalid format got %s, expected namespace:name", envVar)
				continue
			}
			if nameAfterSplit[0] == "" {
				logger.Warning.Printf("Invalid format, missing namespace %s", envVar)
				continue
			}
			dep = Dependency{Name: nameAfterSplit[1], Namespace: nameAfterSplit[0]}
		} else {
			dep = Dependency{Name: envVar, Namespace: namespace}
		}
		envList = append(envList, dep)
	}

	return envList
}

// SplitPodEnvToDeps returns list of PodDependency
func SplitPodEnvToDeps(env string) []PodDependency {
	deps := []PodDependency{}

	namespace := GetBaseNamespace()

	e := os.Getenv(env)
	if e == "" {
		return deps
	}

	err := json.Unmarshal([]byte(e), &deps)
	if err != nil {
		logger.Warning.Printf("Invalid format: %v", e)
		return []PodDependency{}
	}

	for i, dep := range deps {
		if dep.Namespace == "" {
			dep.Namespace = namespace
		}
		deps[i] = dep
	}

	return deps
}

// SplitJobEnvToDeps returns list of JobDependency
func SplitJobEnvToDeps(env, jsonEnv string) []JobDependency {
	deps := []JobDependency{}

	namespace := GetBaseNamespace()

	envVal := os.Getenv(env)
	jsonEnvVal := os.Getenv(jsonEnv)
	if jsonEnvVal != "" {
		if envVal != "" {
			logger.Warning.Printf("Ignoring %s since %s was specified", env, jsonEnv)
		}
		err := json.Unmarshal([]byte(jsonEnvVal), &deps)
		if err != nil {
			logger.Warning.Printf("Invalid format: %s", jsonEnvVal)
			return []JobDependency{}
		}

		valid := []JobDependency{}
		for _, dep := range deps {
			if dep.Namespace == "" {
				dep.Namespace = namespace
			}

			valid = append(valid, dep)
		}

		return valid
	}

	if envVal != "" {
		plainDeps := SplitEnvToDeps(env)

		deps = []JobDependency{}
		for _, dep := range plainDeps {
			deps = append(deps, JobDependency{Name: dep.Name, Namespace: dep.Namespace})
		}

		return deps
	}

	return deps
}

// SplitDeploymentEnvToDeps returns named Deployment dependencies, preferring JSON.
// Invalid configuration is logged and skipped, like other dependency parsers.
func SplitDeploymentEnvToDeps(env, jsonEnv string) []DeploymentDependency {
	deps := []DeploymentDependency{}
	jsonEnvVal := os.Getenv(jsonEnv)
	if jsonEnvVal == "" {
		for _, dep := range SplitEnvToDeps(env) {
			if strings.TrimSpace(dep.Name) == "" {
				logger.Warning.Printf("Invalid %s dependency: name is required", env)
				continue
			}
			deps = append(deps, DeploymentDependency{Name: dep.Name, Namespace: dep.Namespace})
		}
		return deps
	}

	if os.Getenv(env) != "" {
		logger.Warning.Printf("Ignoring %s since %s was specified", env, jsonEnv)
	}
	// Pointer values distinguish a JSON null from an expected empty string.
	var configured []struct {
		Name        string             `json:"name"`
		Namespace   json.RawMessage    `json:"namespace"`
		Annotations map[string]*string `json:"annotations"`
	}
	if err := json.Unmarshal([]byte(jsonEnvVal), &configured); err != nil {
		logger.Warning.Printf("Invalid %s: expected an array of named Deployment dependencies: %v", jsonEnv, err)
		return deps
	}
	if configured == nil {
		logger.Warning.Printf("Invalid %s: expected an array of named Deployment dependencies, not null", jsonEnv)
		return deps
	}

	for i, dep := range configured {
		namespace := GetBaseNamespace()
		if dep.Namespace != nil {
			if string(dep.Namespace) == "null" {
				logger.Warning.Printf("Invalid %s dependency %d: namespace must be a string, not null", jsonEnv, i+1)
				return []DeploymentDependency{}
			}
			if err := json.Unmarshal(dep.Namespace, &namespace); err != nil {
				logger.Warning.Printf("Invalid %s dependency %d: namespace must be a string: %v", jsonEnv, i+1, err)
				return []DeploymentDependency{}
			}
			if namespace == "" {
				namespace = GetBaseNamespace()
			}
		}
		annotations := map[string]string{}
		for key, value := range dep.Annotations {
			if value == nil {
				logger.Warning.Printf("Invalid %s dependency %d: annotation %q must be a string, not null", jsonEnv, i+1, key)
				return []DeploymentDependency{}
			}
			annotations[key] = *value
		}
		if strings.TrimSpace(dep.Name) == "" {
			logger.Warning.Printf("Invalid %s dependency %d: name is required", jsonEnv, i+1)
			continue
		}
		deps = append(deps, DeploymentDependency{
			Name: dep.Name, Namespace: namespace, Annotations: annotations,
		})
	}
	return deps
}

// GetBaseNamespace returns default namespace when user set empty one
func GetBaseNamespace() string {
	namespace, isSet := os.LookupEnv("NAMESPACE")
	if !isSet || namespace == "" {
		namespace = "default"
	}
	return namespace
}

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

type CronJobDependency struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
}

// SplitCronJobEnvToDeps returns named CronJob dependencies, preferring JSON.
// Invalid configuration is logged and skipped, like other dependency parsers.
func SplitCronJobEnvToDeps(env, jsonEnv string) []CronJobDependency {
	deps := []CronJobDependency{}
	jsonEnvVal := os.Getenv(jsonEnv)
	if jsonEnvVal == "" {
		for _, dep := range SplitEnvToDeps(env) {
			if strings.TrimSpace(dep.Name) == "" {
				logger.Warning.Printf("Invalid %s dependency: name is required", env)
				continue
			}
			deps = append(deps, CronJobDependency{Name: dep.Name, Namespace: dep.Namespace})
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
		Labels      map[string]*string `json:"labels"`
	}
	if err := json.Unmarshal([]byte(jsonEnvVal), &configured); err != nil {
		logger.Warning.Printf("Invalid %s: expected an array of named CronJob dependencies: %v", jsonEnv, err)
		return deps
	}
	if configured == nil {
		logger.Warning.Printf("Invalid %s: expected an array of named CronJob dependencies, not null", jsonEnv)
		return deps
	}

	for i, dep := range configured {
		namespace := GetBaseNamespace()
		if dep.Namespace != nil {
			if string(dep.Namespace) == "null" {
				logger.Warning.Printf("Invalid %s dependency %d: namespace must be a string, not null", jsonEnv, i+1)
				return []CronJobDependency{}
			}
			if err := json.Unmarshal(dep.Namespace, &namespace); err != nil {
				logger.Warning.Printf("Invalid %s dependency %d: namespace must be a string: %v", jsonEnv, i+1, err)
				return []CronJobDependency{}
			}
			if namespace == "" {
				namespace = GetBaseNamespace()
			}
		}
		annotations := map[string]string{}
		for key, value := range dep.Annotations {
			if value == nil {
				logger.Warning.Printf("Invalid %s dependency %d: annotation %q must be a string, not null", jsonEnv, i+1, key)
				return []CronJobDependency{}
			}
			annotations[key] = *value
		}
		labels := map[string]string{}
		for key, value := range dep.Labels {
			if value == nil {
				logger.Warning.Printf("Invalid %s dependency %d: label %q must be a string, not null", jsonEnv, i+1, key)
				return []CronJobDependency{}
			}
			labels[key] = *value
		}
		if strings.TrimSpace(dep.Name) == "" {
			logger.Warning.Printf("Invalid %s dependency %d: name is required", jsonEnv, i+1)
			continue
		}
		deps = append(deps, CronJobDependency{
			Name: dep.Name, Namespace: namespace, Annotations: annotations, Labels: labels,
		})
	}
	return deps
}

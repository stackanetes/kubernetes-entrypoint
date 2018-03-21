package env

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/stackanetes/kubernetes-entrypoint/logger"
)

const (
	Separator = ":"
)

type Dependency struct {
	Name      string
	Namespace string
}

type PodDependency struct {
	Labels          map[string]string
	Namespace       string
	RequireSameNode bool
}

func SplitCommand() []string {
	command := os.Getenv("COMMAND")
	if command == "" {
		return []string{}
	}
	commandList := strings.Split(command, " ")
	return commandList
}

//SplitEnvToDeps returns list of namespaces and names pairs
func SplitEnvToDeps(env string) (envList []Dependency) {
	separator := ","

	e := os.Getenv(env)
	if e == "" {
		return envList
	}

	envVars := strings.Split(e, separator)
	namespace := GetBaseNamespace()
	dep := Dependency{}
	for _, envVar := range envVars {
		if strings.Contains(envVar, Separator) {
			nameAfterSplit := strings.Split(envVar, Separator)
			if len(nameAfterSplit) != 2 {
				logger.Warning.Printf("Invalid format got %s, expected namespace:name", envVar)
				continue
			}
			if nameAfterSplit[0] == "" {
				logger.Warning.Printf("Invalid format, missing namespace", envVar)
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

//SplitPodEnvToDeps returns list of PodDependency
func SplitPodEnvToDeps(env string) []PodDependency {
	deps := []PodDependency{}

	namespace := GetBaseNamespace()

	e := os.Getenv(env)
	if e == "" {
		return deps
	}

	err := json.Unmarshal([]byte(e), &deps)
	if err != nil {
		logger.Warning.Printf("Invalid format: ", e)
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

//GetBaseNamespace returns default namespace when user set empty one
func GetBaseNamespace() string {
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	return namespace
}

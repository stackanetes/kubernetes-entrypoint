package env

import (
	"encoding/json"
	"os"
	"strconv"
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

type JobDependency struct {
	Name      string
	Labels    map[string]string
	Namespace string
}

type CustomResourceDependency struct {
	ApiVersion string
	Name       string
	Namespace  string
	Kind       string
	Fields     []map[string]string
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
				logger.Warning("Invalid format got %s, expected namespace:name", envVar)
				continue
			}
			if nameAfterSplit[0] == "" {
				logger.Warning("Invalid format, missing namespace %s", envVar)
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
		logger.Warning("Invalid format: %v", e)
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

//SplitJobEnvToDeps returns list of JobDependency
func SplitJobEnvToDeps(env string, jsonEnv string) []JobDependency {
	deps := []JobDependency{}

	namespace := GetBaseNamespace()

	envVal := os.Getenv(env)
	jsonEnvVal := os.Getenv(jsonEnv)
	if jsonEnvVal != "" {
		if envVal != "" {
			logger.Warning("Ignoring %s since %s was specified", env, jsonEnv)
		}
		err := json.Unmarshal([]byte(jsonEnvVal), &deps)
		if err != nil {
			logger.Warning("Invalid format: %s", jsonEnvVal)
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

func SplitCustomResourceEnvToDeps(jsonEnv string) []CustomResourceDependency {
	deps := []CustomResourceDependency{}
	namespace := GetBaseNamespace()
	jsonEnvVal := os.Getenv(jsonEnv)
	err := json.Unmarshal([]byte(jsonEnvVal), &deps)
	if err != nil {
		logger.Warning("Invalid format: %s", jsonEnvVal)
		return []CustomResourceDependency{}
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

//OutputJSON returns true if the OUTPUT_JSON environment variable is set
func OutputJSON() bool {
	return strings.EqualFold(os.Getenv("OUTPUT_JSON"), "true")
}

//GetEnvFloat gets the value of an environment variable as a float
func GetEnvFloat(name string, def float64) float64 {
	if strVal, ok := os.LookupEnv(name); ok {
		if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
			return floatVal
		}
	}
	return def
}

//Backoff returns the next slot in an exponential backoff
func Backoff(n, base, max float64) float64 {
	if n*base < max {
		return n * base
	}
	return max
}

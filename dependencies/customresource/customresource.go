package customresource

import (
	"fmt"
	"strings"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	env "github.com/stackanetes/kubernetes-entrypoint/util/env"
)

// A CustomResource represents the desired state of a CustomResourceDefinition
type CustomResource struct {
	APIVersion string
	Kind       string
	Name       string
	Namespace  string
	Fields     []map[string]string
}

func init() {
	crEnv := fmt.Sprintf("%sCUSTOM_RESOURCE", entry.DependencyPrefix)
	if crDeps := env.SplitCustomResourceEnvToDeps(crEnv); crDeps != nil {
		for _, dep := range crDeps {
			cr := NewCustomResource(&dep)
			entry.Register(cr)
		}
	}
}

// NewCustomResource creates a CustomResource from a dependecy
func NewCustomResource(dep *env.CustomResourceDependency) *CustomResource {
	return &CustomResource{
		APIVersion: dep.ApiVersion,
		Kind:       dep.Kind,
		Name:       dep.Name,
		Namespace:  dep.Namespace,
		Fields:     dep.Fields,
	}
}

// IsResolved will return true when the values for each key in cr.Fields is the same as the resource in the cluster
func (cr CustomResource) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	resourceName, err := entrypoint.Client().GetResourceName(cr.Kind, cr.APIVersion)
	if err != nil {
		return false, err
	}

	myCustomResource, err := entrypoint.Client().CustomResource(cr.APIVersion, cr.Namespace, resourceName, cr.Name)
	if err != nil {
		return false, err
	}

	for _, field := range cr.Fields {
		key := field["key"]
		expected := field["value"]

		// Extract the specified value from the resource
		actual := extractVal(myCustomResource, key)

		if actual != expected {
			return false, fmt.Errorf("Expected value of [%s] to be [%s], but got [%s]", key, expected, actual)
		}
	}

	return true, nil
}

func extractVal(customResource map[string]interface{}, key string) string {
	for i := strings.Index(key, "."); i != -1; i = strings.Index(key, ".") {
		first := key[:i]
		key = key[i+1:]
		customResource = customResource[first].(map[string]interface{})
	}
	return customResource[key].(string)
}

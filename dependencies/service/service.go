package service

import (
	"fmt"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FailingStatusFormat = "Service %v has no endpoints"

type Service struct {
	name      string
	namespace string
}

func init() {
	serviceEnv := fmt.Sprintf("%sSERVICE", entry.DependencyPrefix)
	if serviceDeps := env.SplitEnvToDeps(serviceEnv); serviceDeps != nil {
		if len(serviceDeps) > 0 {
			for _, dep := range serviceDeps {
				entry.Register(NewService(dep.Name, dep.Namespace))
			}
		}
	}
}

func NewService(name string, namespace string) Service {
	return Service{
		name:      name,
		namespace: namespace,
	}

}

func (s Service) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	e, err := entrypoint.Client().Endpoints(s.namespace).Get(s.name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	for _, subset := range e.Subsets {
		if len(subset.Addresses) > 0 {
			return true, nil
		}
	}
	return false, fmt.Errorf(FailingStatusFormat, s.name)
}

// GetDependency returns the details associated with this dependency
func (s Service) GetDependency() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "Service",
		"Details": s,
	}
}

func (s Service) String() string {
	return fmt.Sprintf("Service %s in namespace %s", s.name, s.namespace)
}

package service

import (
	"fmt"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
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
	e, err := entrypoint.Client().Endpoints(s.namespace).Get(s.name)
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

func (s Service) String() string {
	return fmt.Sprintf("Service %s in namespace %s", s.name, s.namespace)
}

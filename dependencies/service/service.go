package service

import (
	"fmt"
	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

type Service struct {
	name string
}

func init() {
	serviceEnv := fmt.Sprintf("%sSERVICE", entry.DependencyPrefix)
	if serviceDeps := env.SplitEnvToList(serviceEnv); len(serviceDeps) > 0 {
		for _, dep := range serviceDeps {
			entry.Register(NewService(dep))
		}
	}
}

func NewService(name string) Service {
	return Service{name: name}

}

func (s Service) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	e, err := entrypoint.Client().Endpoints(entrypoint.GetNamespace()).Get(s.GetName())
	if err != nil {
		return false, err
	}
	if len(e.Subsets) > 0 {
		return true, nil
	}
	return false, fmt.Errorf("Service %v has no endpoints", s.GetName())
}

func (s Service) GetName() string {
	return s.name
}

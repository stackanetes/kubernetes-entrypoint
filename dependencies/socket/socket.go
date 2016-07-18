package socket

import (
	"fmt"
	entry "github.com/stackanetes/docker-entrypoint/dependencies"
	"github.com/stackanetes/docker-entrypoint/util/env"
	"os"
)

type Socket struct {
	name string
}

func init() {
	socketEnv := fmt.Sprintf("%sSOCKET", entry.DependencyPrefix)
	if socketDeps := env.SplitEnvToList(socketEnv); len(socketDeps) > 0 {
		for _, dep := range socketDeps {
			entry.Register(NewSocket(dep))
		}
	}
}

func NewSocket(name string) Socket {
	return Socket{name: name}
}

func (s Socket) GetName() string {
	return s.name
}

func (s Socket) IsResolved(entrypoint *entry.Entrypoint) (bool, error) {
	_, err := os.Stat(s.GetName())
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("Socket %v doesn't exists", s.GetName())
	}
	if os.IsPermission(err) {
		return false, fmt.Errorf("I have no permission to %v", s.GetName())
	}
	return false, err
}

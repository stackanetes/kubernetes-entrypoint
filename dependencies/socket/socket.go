package socket

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

const (
	NonExistingErrorFormat = "Socket %v doesn't exists"
	NoPermsErrorFormat     = "I have no permission to %v"
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

func (s Socket) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	_, err := os.Stat(s.GetName())
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf(NonExistingErrorFormat, s.GetName())
	}
	if os.IsPermission(err) {
		return false, fmt.Errorf(NoPermsErrorFormat, s.GetName())
	}
	return false, err
}

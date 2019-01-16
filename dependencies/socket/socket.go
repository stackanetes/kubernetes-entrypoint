package socket

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

const (
	NonExistingErrorFormat = "%s doesn't exists"
	NoPermsErrorFormat     = "I have no permission to %s"
	NamespaceNotSupported  = "Socket doesn't accept namespace"
)

type Socket struct {
	name string
}

func init() {
	socketEnv := fmt.Sprintf("%sSOCKET", entry.DependencyPrefix)
	if util.ContainsSeparator(socketEnv, "Socket") {
		logger.Error(NamespaceNotSupported)
		os.Exit(1)
	}
	if socketDeps := env.SplitEnvToDeps(socketEnv); socketDeps != nil {
		if len(socketDeps) > 0 {
			for _, dep := range socketDeps {
				entry.Register(NewSocket(dep.Name))
			}
		}
	}
}

func NewSocket(name string) Socket {
	return Socket{name: name}
}

func (s Socket) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	_, err := os.Stat(s.name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf(NonExistingErrorFormat, s)
	}
	if os.IsPermission(err) {
		return false, fmt.Errorf(NoPermsErrorFormat, s)
	}
	return false, err
}

func (s Socket) String() string {
	return fmt.Sprintf("Socket %s", s.name)
}

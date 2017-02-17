package file

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

type File struct {
	name string
}

func init() {
	fileEnv := fmt.Sprintf("%sFILE", entry.DependencyPrefix)
	if fileDeps := env.SplitEnvToList(fileEnv); len(fileDeps) > 0 {
		for _, dep := range fileDeps {
			entry.Register(NewFile(dep))
		}
	}
}

func NewFile(name string) File {
	return File{name: name}
}

func (s File) GetName() string {
	return s.name
}

func (s File) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	_, err := os.Stat(s.GetName())
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("File %v doesn't exists", s.GetName())
	}
	if os.IsPermission(err) {
		return false, fmt.Errorf("I have no permission to %v", s.GetName())
	}
	return false, err
}

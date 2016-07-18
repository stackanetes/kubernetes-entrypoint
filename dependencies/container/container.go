package container

import (
	"fmt"
	entry "github.com/stackanetes/docker-entrypoint/dependencies"
	"github.com/stackanetes/docker-entrypoint/util/env"
	"os"
)

type Container struct {
	name string
}

func init() {
	containerEnv := fmt.Sprintf("%sCONTAINER", entry.DependencyPrefix)
	if containerDeps := env.SplitEnvToList(containerEnv); len(containerDeps) > 0 {
		for _, dep := range containerDeps {
			entry.Register(NewContainer(dep))
		}
	}
}

func NewContainer(name string) Container {
	return Container{name: name}

}

func (c Container) IsResolved(entrypoint *entry.Entrypoint) (bool, error) {
	myPodName := os.Getenv("POD_NAME")
	if myPodName == "" {
		return false, fmt.Errorf("Environment variable POD_NAME not set")
	}
	pod, err := entrypoint.Client.Pods(entrypoint.Namespace).Get(myPodName)
	if err != nil {
		return false, err
	}
	containers := pod.Status.ContainerStatuses
	for _, container := range containers {
		if container.Name == c.GetName() && container.State.Running != nil {
			return true, nil
		}
	}
	return false, nil
}

func (c Container) GetName() string {
	return c.name
}

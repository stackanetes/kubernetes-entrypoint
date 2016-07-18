package job

import (
	"fmt"

	entry "github.com/stackanetes/docker-entrypoint/dependencies"
	"github.com/stackanetes/docker-entrypoint/util/env"
)

type Job struct {
	name string
}

func init() {
	jobsEnv := fmt.Sprintf("%sJOBS", entry.DependencyPrefix)
	if jobsDeps := env.SplitEnvToList(jobsEnv); len(jobsDeps) > 0 {
		for _, dep := range jobsDeps {
			entry.Register(NewJob(dep))
		}
	}
}

func NewJob(name string) Job {
	return Job{name: name}

}

func (j Job) IsResolved(entrypoint *entry.Entrypoint) (bool, error) {
	job, err := entrypoint.Client.ExtensionsClient.Jobs(entrypoint.Namespace).Get(j.name)
	if err != nil {
		return false, err
	}
	if job.Status.Succeeded == 0 {
		return false, fmt.Errorf("Job %v is not completed yet", j.GetName())
	}
	return true, nil
}

func (j Job) GetName() string {
	return j.name
}

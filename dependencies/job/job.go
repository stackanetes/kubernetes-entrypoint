package job

import (
	"fmt"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

const FailingStatusFormat = "Job %v is not completed yet"

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

func (j Job) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	job, err := entrypoint.Client().Jobs(entrypoint.GetNamespace()).Get(j.GetName())
	if err != nil {
		return false, err
	}
	if job.Status.Succeeded == 0 {
		return false, fmt.Errorf(FailingStatusFormat, j.GetName())
	}
	return true, nil
}

func (j Job) GetName() string {
	return j.name
}

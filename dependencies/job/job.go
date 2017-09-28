package job

import (
	"fmt"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

const FailingStatusFormat = "Job %s is not completed yet"

type Job struct {
	name      string
	namespace string
}

func init() {
	jobsEnv := fmt.Sprintf("%sJOBS", entry.DependencyPrefix)
	if jobsDeps := env.SplitEnvToDeps(jobsEnv); jobsDeps != nil {
		if len(jobsDeps) > 0 {
			for _, dep := range jobsDeps {
				entry.Register(NewJob(dep.Name, dep.Namespace))
			}
		}
	}
}

func NewJob(name string, namespace string) Job {
	return Job{
		name:      name,
		namespace: namespace,
	}

}

func (j Job) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	job, err := entrypoint.Client().Jobs(j.namespace).Get(j.name)
	if err != nil {
		return false, err
	}
	if job.Status.Succeeded == 0 {
		return false, fmt.Errorf(FailingStatusFormat, j)
	}
	return true, nil
}

func (j Job) String() string {
	return fmt.Sprintf("Job %s in namespace %s", j.name, j.namespace)
}

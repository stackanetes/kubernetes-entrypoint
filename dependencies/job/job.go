package job

import (
	"fmt"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FailingStatusFormat = "Job %s is not completed yet"

type Job struct {
	name      string
	namespace string
	labels    map[string]string
}

func init() {
	jobsEnv := fmt.Sprintf("%sJOBS", entry.DependencyPrefix)
	jobsJsonEnv := fmt.Sprintf("%s%s", jobsEnv, entry.JsonSuffix)
	if jobsDeps := env.SplitJobEnvToDeps(jobsEnv, jobsJsonEnv); jobsDeps != nil {
		if len(jobsDeps) > 0 {
			for _, dep := range jobsDeps {
				job := NewJob(dep.Name, dep.Namespace, dep.Labels)
				if job != nil {
					entry.Register(*job)
				}
			}
		}
	}
}

func NewJob(name string, namespace string, labels map[string]string) *Job {
	if name != "" && labels != nil {
		logger.Warning("Cannot specify both name and labels for job depependency")
		return nil
	}
	return &Job{
		name:      name,
		namespace: namespace,
		labels:    labels,
	}
}

func (j Job) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	iface := entrypoint.Client().Jobs(j.namespace)
	var jobs []v1.Job

	if j.name != "" {
		job, err := iface.Get(j.name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		jobs = []v1.Job{*job}
	} else if j.labels != nil {
		labelSelector := &metav1.LabelSelector{MatchLabels: j.labels}
		label := metav1.FormatLabelSelector(labelSelector)
		opts := metav1.ListOptions{LabelSelector: label}
		jobList, err := iface.List(opts)
		if err != nil {
			return false, err
		}
		jobs = jobList.Items
	}
	if len(jobs) == 0 {
		return false, fmt.Errorf("No matching jobs found: %v", j)
	}

	for _, job := range jobs {
		if job.Status.Succeeded == 0 {
			return false, fmt.Errorf(FailingStatusFormat, j)
		}
	}
	return true, nil
}

// GetDependency returns the details associated with this dependency
func (j Job) GetDependency() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "Job",
		"Details": j,
	}
}

func (j Job) String() string {
	var prefix string
	if j.name != "" {
		prefix = fmt.Sprintf("Job %s", j.name)
	} else if j.labels != nil {
		prefix = fmt.Sprintf("Jobs with labels %s", j.labels)
	} else {
		prefix = "Jobs"
	}
	return fmt.Sprintf("%s in namespace %s", prefix, j.namespace)
}

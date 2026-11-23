// Copyright 2026 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cronjob

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

type CronJob struct {
	name        string
	namespace   string
	annotations map[string]string
	labels      map[string]string
}

func init() {
	registerCronJobs(entry.Register)
}

func registerCronJobs(register func(entry.Resolver)) {
	cronJobEnv := entry.DependencyPrefix + "CRONJOB"
	for _, dep := range env.SplitCronJobEnvToDeps(cronJobEnv, cronJobEnv+entry.JsonSuffix) {
		if cronJob := NewCronJob(dep.Name, dep.Namespace, dep.Annotations, dep.Labels); cronJob != nil {
			register(*cronJob)
		}
	}
}

func NewCronJob(name, namespace string, annotations, labels map[string]string) *CronJob {
	if strings.TrimSpace(name) == "" {
		logger.Warning.Printf("Cannot create CronJob dependency: name is required")
		return nil
	}
	return &CronJob{name: name, namespace: namespace, annotations: maps.Clone(annotations), labels: maps.Clone(labels)}
}

func (c CronJob) String() string {
	return fmt.Sprintf("CronJob %s in namespace %s", c.name, c.namespace)
}

func (c CronJob) IsResolved(ctx context.Context, entrypoint entry.EntrypointInterface) (bool, error) {
	iface := entrypoint.Client().CronJobs(c.namespace)
	snapshot, err := iface.Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("getting %s: %w", c, err)
	}
	if err := c.matchCronJob(snapshot); err != nil {
		return false, fmt.Errorf("%s: %w", c, err)
	}

	// Do not select Jobs using the target labels: that would hide old versions
	// and unmarked Jobs. CronJob status.active can also lag behind actual Jobs.
	jobs, err := entrypoint.Client().Jobs(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("listing Jobs for %s: %w", c, err)
	}
	if jobs == nil {
		return false, fmt.Errorf("%s: Jobs were not returned by the API", c)
	}
	for i := range jobs.Items {
		job := &jobs.Items[i]
		owner := metav1.GetControllerOf(job)
		if owner == nil || owner.Kind != "CronJob" || owner.Name != snapshot.Name ||
			schema.FromAPIVersionAndKind(owner.APIVersion, owner.Kind).Group != batchv1.GroupName {
			continue
		}
		// Also drain Jobs from a previous incarnation of this named CronJob.
		current := owner.UID == snapshot.UID && (len(c.annotations) > 0 || len(c.labels) > 0) &&
			matchMetadata("annotation", job.Annotations, c.annotations) == nil &&
			matchMetadata("label", job.Labels, c.labels) == nil
		if !current && !jobFinished(job) {
			return false, fmt.Errorf("%s: old-version Job %s is not finished", c, job.Name)
		}
	}

	// The two resources cannot be read atomically. Before accepting success,
	// ensure the CronJob was not replaced or updated while its Jobs were listed.
	latest, err := iface.Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("rechecking %s: %w", c, err)
	}
	if err := c.matchCronJob(latest); err != nil {
		return false, fmt.Errorf("%s: %w", c, err)
	}
	if latest.UID != snapshot.UID || latest.Generation != snapshot.Generation {
		return false, fmt.Errorf("%s: CronJob changed while checking Jobs", c)
	}
	return true, nil
}

func (c CronJob) matchCronJob(cronJob *batchv1.CronJob) error {
	if cronJob == nil {
		return errors.New("CronJob was not returned by the API")
	}
	if cronJob.DeletionTimestamp != nil {
		return errors.New("CronJob is being deleted")
	}
	if err := matchMetadata("annotation", cronJob.Annotations, c.annotations); err != nil {
		return err
	}
	return matchMetadata("label", cronJob.Labels, c.labels)
}

func matchMetadata(kind string, actual, expected map[string]string) error {
	for key, value := range expected {
		got, present := actual[key]
		if !present {
			return fmt.Errorf("pending target update: metadata %s %q is missing, expected %q", kind, key, value)
		}
		if got != value {
			return fmt.Errorf("pending target update: metadata %s %q is %q, expected %q", kind, key, got, value)
		}
	}
	return nil
}

// Zero active Pods does not prove completion: a Job may be pending, retrying,
// suspended, or still terminating. Only terminal conditions acknowledge it.
func jobFinished(job *batchv1.Job) bool {
	if job.Status.Active != 0 || (job.Status.Terminating != nil && *job.Status.Terminating != 0) {
		return false
	}
	for _, condition := range job.Status.Conditions {
		if condition.Status == corev1.ConditionTrue &&
			(condition.Type == batchv1.JobComplete || condition.Type == batchv1.JobFailed) {
			return true
		}
	}
	return false
}

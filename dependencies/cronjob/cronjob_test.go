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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	batchclient "k8s.io/client-go/kubernetes/typed/batch/v1"

	"opendev.org/airship/kubernetes-entrypoint/client"
	"opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/mocks"
)

const releaseAnnotation = "example.com/release-id"

// Embedded interfaces make accidental API calls other than Get/List fail.
type cronJobAPI struct {
	batchclient.CronJobInterface
	get func(context.Context, string, metav1.GetOptions) (*batchv1.CronJob, error)
}

func (a cronJobAPI) Get(ctx context.Context, name string, opts metav1.GetOptions) (*batchv1.CronJob, error) {
	return a.get(ctx, name, opts)
}

type jobAPI struct {
	batchclient.JobInterface
	list func(context.Context, metav1.ListOptions) (*batchv1.JobList, error)
}

func (a jobAPI) List(ctx context.Context, opts metav1.ListOptions) (*batchv1.JobList, error) {
	return a.list(ctx, opts)
}

type namespacedClient struct {
	client.ClientInterface
	cronJobs func(string) batchclient.CronJobInterface
	jobs     func(string) batchclient.JobInterface
}

func (c namespacedClient) CronJobs(namespace string) batchclient.CronJobInterface {
	return c.cronJobs(namespace)
}

func (c namespacedClient) Jobs(namespace string) batchclient.JobInterface {
	return c.jobs(namespace)
}

type cronJobEntrypoint struct {
	entrypoint.EntrypointInterface
	kubeClient client.ClientInterface
}

func (e cronJobEntrypoint) Client() client.ClientInterface {
	return e.kubeClient
}

func targetCronJob() *batchv1.CronJob {
	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cleanup", Namespace: "application", UID: "cronjob-uid", Generation: 7,
			Annotations: map[string]string{releaseAnnotation: "my-release:42"},
			Labels:      map[string]string{"version": "v42"},
		},
		Spec: batchv1.CronJobSpec{
			Schedule: "* * * * *",
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{releaseAnnotation: "my-release:42"},
					Labels:      map[string]string{"version": "v42"},
				},
			},
		},
	}
}

func oldJob() batchv1.Job {
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cleanup-old", Namespace: "application",
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(targetCronJob(), batchv1.SchemeGroupVersion.WithKind("CronJob"))},
			Annotations:     map[string]string{releaseAnnotation: "my-release:41"},
			Labels:          map[string]string{"version": "v41"},
		},
		Status: batchv1.JobStatus{Active: 1},
	}
}

func finishJob(job *batchv1.Job, condition batchv1.JobConditionType) {
	job.Status = batchv1.JobStatus{Conditions: []batchv1.JobCondition{{Type: condition, Status: corev1.ConditionTrue}}}
}

func targetDependency() *CronJob {
	return NewCronJob("cleanup", "application", map[string]string{releaseAnnotation: "my-release:42"}, map[string]string{"version": "v42"})
}

func snapshotEntrypoint(snapshot *batchv1.CronJob, jobs []batchv1.Job) cronJobEntrypoint {
	mockClient := mocks.NewClient()
	mockClient.CronJobInterface = cronJobAPI{get: func(context.Context, string, metav1.GetOptions) (*batchv1.CronJob, error) {
		return snapshot.DeepCopy(), nil
	}}
	mockClient.JobInterface = jobAPI{list: func(context.Context, metav1.ListOptions) (*batchv1.JobList, error) {
		return &batchv1.JobList{Items: jobs}, nil
	}}
	return cronJobEntrypoint{kubeClient: mockClient}
}

var _ = Describe("CronJob configuration", func() {
	It("constructs a named dependency and owns its metadata expectations", func() {
		annotations := map[string]string{releaseAnnotation: "42"}
		labels := map[string]string{"version": "v42"}
		dep := NewCronJob("cleanup", "application", annotations, labels)
		Expect(dep.name).To(Equal("cleanup"))
		Expect(dep.namespace).To(Equal("application"))
		Expect(dep.String()).To(Equal("CronJob cleanup in namespace application"))
		annotations[releaseAnnotation] = "43"
		labels["version"] = "v43"
		Expect(dep.annotations).To(Equal(map[string]string{releaseAnnotation: "42"}))
		Expect(dep.labels).To(Equal(map[string]string{"version": "v42"}))
		Expect(NewCronJob("", "application", nil, nil)).To(BeNil())
		Expect(NewCronJob(" ", "application", nil, nil)).To(BeNil())
	})

	DescribeTable("registers dependencies using the existing resolver contract",
		func(plain, json string, expected []CronJob) {
			GinkgoT().Setenv("NAMESPACE", "application")
			GinkgoT().Setenv("DEPENDENCY_CRONJOB", plain)
			GinkgoT().Setenv("DEPENDENCY_CRONJOB_JSON", json)
			var registered []CronJob
			registerCronJobs(func(resolver entrypoint.Resolver) {
				Expect(resolver).To(BeAssignableToTypeOf(CronJob{}))
				registered = append(registered, resolver.(CronJob))
			})
			Expect(registered).To(Equal(expected))
		},
		Entry("plain", "cleanup,other:backup", "", []CronJob{{name: "cleanup", namespace: "application"}, {name: "backup", namespace: "other"}}),
		Entry("JSON precedence", "ignored", `[{"name":"cleanup","annotations":{"release":"42"},"labels":{"version":"v42"}},{"name":"backup","namespace":"other"}]`, []CronJob{
			{name: "cleanup", namespace: "application", annotations: map[string]string{"release": "42"}, labels: map[string]string{"version": "v42"}},
			{name: "backup", namespace: "other", annotations: map[string]string{}, labels: map[string]string{}},
		}),
		Entry("unset", "", "", nil),
		Entry("invalid names skipped", "", `[{"labels":{"release":"42"}}]`, nil),
	)
})

var _ = Describe("CronJob update matching", func() {
	DescribeTable("matches every expected top-level annotation and label exactly",
		func(change func(*batchv1.CronJob), resolved bool, diagnostic string) {
			snapshot := targetCronJob()
			change(snapshot)
			result, err := targetDependency().IsResolved(context.Background(), snapshotEntrypoint(snapshot, nil))
			Expect(result).To(Equal(resolved))
			if resolved {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(ContainSubstring(diagnostic)))
			}
		},
		Entry("target already submitted, no transition required", func(c *batchv1.CronJob) {}, true, ""),
		Entry("extra metadata allowed", func(c *batchv1.CronJob) { c.Labels["extra"] = "ok"; c.Annotations["extra"] = "ok" }, true, ""),
		Entry("old annotation", func(c *batchv1.CronJob) { c.Annotations[releaseAnnotation] = "my-release:41" }, false, "pending target update"),
		Entry("newer annotation never silently accepted", func(c *batchv1.CronJob) { c.Annotations[releaseAnnotation] = "my-release:43" }, false, "pending target update"),
		Entry("missing annotation, matching template is insufficient", func(c *batchv1.CronJob) { c.Annotations = nil }, false, "annotation"),
		Entry("old label despite matching annotation", func(c *batchv1.CronJob) { c.Labels["version"] = "v41" }, false, "label"),
		Entry("newer label never silently accepted", func(c *batchv1.CronJob) { c.Labels["version"] = "v43" }, false, "label"),
		Entry("missing label, matching template is insufficient", func(c *batchv1.CronJob) { c.Labels = nil }, false, "is missing"),
		Entry("being deleted", func(c *batchv1.CronJob) { now := metav1.Now(); c.DeletionTimestamp = &now }, false, "being deleted"),
		Entry("suspended CronJob can still be updated and drained", func(c *batchv1.CronJob) { c.Spec.Suspend = new(true) }, true, ""),
		Entry("stale active references do not override the Job list", func(c *batchv1.CronJob) { c.Status.Active = []corev1.ObjectReference{{Name: "gone"}} }, true, ""),
	)

	DescribeTable("supports labels or annotations alone, including present empty values",
		func(annotations, labels map[string]string) {
			snapshot := targetCronJob()
			dep := NewCronJob("cleanup", "application", annotations, labels)
			resolved, err := dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, nil))
			Expect(resolved).To(BeFalse())
			Expect(err).To(MatchError(ContainSubstring("is missing")))
			snapshot.Annotations["empty"] = ""
			snapshot.Labels["empty"] = ""
			resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, nil))
			Expect(resolved).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		},
		Entry("annotation", map[string]string{"empty": ""}, nil),
		Entry("label", nil, map[string]string{"empty": ""}),
	)

	It("requires every key when multiple annotations and labels are expected", func() {
		snapshot := targetCronJob()
		dep := targetDependency()
		dep.annotations["another-annotation"] = "expected"
		dep.labels["another-label"] = "expected"
		resolved, err := dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, nil))
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring("another-annotation")))
		snapshot.Annotations["another-annotation"] = "expected"
		resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, nil))
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring("another-label")))
		snapshot.Labels["another-label"] = "expected"
		resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, nil))
		Expect(resolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("CronJob Job draining", func() {
	DescribeTable("waits only for unfinished old-version Jobs owned by the named CronJob",
		func(change func(*batchv1.Job), resolved bool) {
			job := oldJob()
			change(&job)
			result, err := targetDependency().IsResolved(context.Background(), snapshotEntrypoint(targetCronJob(), []batchv1.Job{job}))
			Expect(result).To(Equal(resolved))
			if resolved {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(ContainSubstring("old-version Job cleanup-old is not finished")))
			}
		},
		Entry("old active", func(j *batchv1.Job) {}, false),
		Entry("no metadata", func(j *batchv1.Job) { j.Annotations = nil; j.Labels = nil }, false),
		Entry("pending with zero active", func(j *batchv1.Job) { j.Status.Active = 0 }, false),
		Entry("retrying with failed Pods", func(j *batchv1.Job) { j.Status = batchv1.JobStatus{Failed: 2} }, false),
		Entry("partial success is not terminal", func(j *batchv1.Job) { j.Status = batchv1.JobStatus{Succeeded: 1} }, false),
		Entry("suspended old Job", func(j *batchv1.Job) { j.Spec.Suspend = new(true); j.Status.Active = 0 }, false),
		Entry("deleting old Job", func(j *batchv1.Job) { now := metav1.Now(); j.DeletionTimestamp = &now; j.Status.Active = 0 }, false),
		Entry("Complete True", func(j *batchv1.Job) { finishJob(j, batchv1.JobComplete) }, true),
		Entry("Failed True", func(j *batchv1.Job) { finishJob(j, batchv1.JobFailed) }, true),
		Entry("Complete False", func(j *batchv1.Job) {
			finishJob(j, batchv1.JobComplete)
			j.Status.Conditions[0].Status = corev1.ConditionFalse
		}, false),
		Entry("Failed Unknown", func(j *batchv1.Job) {
			finishJob(j, batchv1.JobFailed)
			j.Status.Conditions[0].Status = corev1.ConditionUnknown
		}, false),
		Entry("FailureTarget is not terminal", func(j *batchv1.Job) { finishJob(j, batchv1.JobFailureTarget) }, false),
		Entry("SuccessCriteriaMet is not terminal", func(j *batchv1.Job) { finishJob(j, batchv1.JobSuccessCriteriaMet) }, false),
		Entry("terminal but active Pods remain", func(j *batchv1.Job) { finishJob(j, batchv1.JobComplete); j.Status.Active = 1 }, false),
		Entry("terminal but terminating Pods remain", func(j *batchv1.Job) { finishJob(j, batchv1.JobFailed); j.Status.Terminating = new(int32(1)) }, false),
		Entry("terminal and zero terminating Pods", func(j *batchv1.Job) { finishJob(j, batchv1.JobComplete); j.Status.Terminating = new(int32(0)) }, true),
		Entry("current active Job", func(j *batchv1.Job) { j.Annotations = targetCronJob().Annotations; j.Labels = targetCronJob().Labels }, true),
		Entry("matching annotation alone is not enough", func(j *batchv1.Job) { j.Annotations = targetCronJob().Annotations }, false),
		Entry("matching label alone is not enough", func(j *batchv1.Job) { j.Labels = targetCronJob().Labels }, false),
		Entry("Pod template metadata cannot mark the Job", func(j *batchv1.Job) {
			j.Spec.Template.Annotations = targetCronJob().Annotations
			j.Spec.Template.Labels = targetCronJob().Labels
		}, false),
		Entry("another CronJob", func(j *batchv1.Job) { j.OwnerReferences[0].Name = "another" }, true),
		Entry("same name but another resource kind", func(j *batchv1.Job) { j.OwnerReferences[0].Kind = "Deployment" }, true),
		Entry("same kind but another API group", func(j *batchv1.Job) { j.OwnerReferences[0].APIVersion = "example.com/v1" }, true),
		Entry("manual Job without controller", func(j *batchv1.Job) { j.OwnerReferences = nil }, true),
		Entry("non-controlling reference", func(j *batchv1.Job) { j.OwnerReferences[0].Controller = new(false) }, true),
		Entry("previous CronJob incarnation", func(j *batchv1.Job) {
			j.OwnerReferences[0].UID = "previous-uid"
			j.Annotations = targetCronJob().Annotations
			j.Labels = targetCronJob().Labels
		}, false),
		Entry("completed previous incarnation", func(j *batchv1.Job) { j.OwnerReferences[0].UID = "previous-uid"; finishJob(j, batchv1.JobComplete) }, true),
		Entry("legacy batch API owner", func(j *batchv1.Job) { j.OwnerReferences[0].APIVersion = "batch/v1beta1" }, false),
	)

	DescribeTable("without version expectations drains all unfinished owned Jobs",
		func(annotations, labels map[string]string) {
			dep := NewCronJob("cleanup", "application", annotations, labels)
			job := oldJob()
			job.Annotations = targetCronJob().Annotations
			job.Labels = targetCronJob().Labels
			resolved, err := dep.IsResolved(context.Background(), snapshotEntrypoint(targetCronJob(), []batchv1.Job{job}))
			Expect(resolved).To(BeFalse())
			Expect(err).To(HaveOccurred())
			finishJob(&job, batchv1.JobComplete)
			resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(targetCronJob(), []batchv1.Job{job}))
			Expect(resolved).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		},
		Entry("nil maps", nil, nil),
		Entry("empty maps", map[string]string{}, map[string]string{}),
	)

	DescribeTable("allows current Jobs with a single configured marker",
		func(annotations, labels map[string]string) {
			dep := NewCronJob("cleanup", "application", annotations, labels)
			job := oldJob()
			job.Annotations = annotations
			job.Labels = labels
			resolved, err := dep.IsResolved(context.Background(), snapshotEntrypoint(targetCronJob(), []batchv1.Job{job}))
			Expect(resolved).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		},
		Entry("annotation", targetCronJob().Annotations, nil),
		Entry("label", nil, targetCronJob().Labels),
	)

	It("waits through release submission and drains every old Job without requiring a new Job to run", func() {
		snapshot := targetCronJob()
		snapshot.Annotations[releaseAnnotation] = "my-release:41"
		first, second := oldJob(), oldJob()
		second.Name = "cleanup-older"
		dep := targetDependency()
		resolved, err := dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, []batchv1.Job{first, second}))
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring("pending target update")))
		snapshot.Annotations = targetCronJob().Annotations
		resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, []batchv1.Job{first, second}))
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring(first.Name)))
		finishJob(&first, batchv1.JobComplete)
		resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, []batchv1.Job{first, second}))
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring(second.Name)))
		finishJob(&second, batchv1.JobFailed)
		resolved, err = dep.IsResolved(context.Background(), snapshotEntrypoint(snapshot, []batchv1.Job{first, second}))
		Expect(resolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("CronJob API handling", func() {
	It("uses the configured namespace and context and lists Jobs without a version selector", func() {
		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)
		gets, lists := 0, 0
		mockClient := mocks.NewClient()
		mockClient.CronJobInterface = cronJobAPI{get: func(actual context.Context, name string, opts metav1.GetOptions) (*batchv1.CronJob, error) {
			Expect(actual).To(BeIdenticalTo(ctx))
			Expect(name).To(Equal("cleanup"))
			Expect(opts).To(Equal(metav1.GetOptions{}))
			gets++
			return targetCronJob(), nil
		}}
		mockClient.JobInterface = jobAPI{list: func(actual context.Context, opts metav1.ListOptions) (*batchv1.JobList, error) {
			Expect(actual).To(BeIdenticalTo(ctx))
			Expect(opts).To(Equal(metav1.ListOptions{}))
			lists++
			return &batchv1.JobList{}, nil
		}}
		kubeClient := namespacedClient{
			cronJobs: func(namespace string) batchclient.CronJobInterface {
				Expect(namespace).To(Equal("application"))
				return mockClient.CronJobs(namespace)
			},
			jobs: func(namespace string) batchclient.JobInterface {
				Expect(namespace).To(Equal("application"))
				return mockClient.Jobs(namespace)
			},
		}
		resolved, err := targetDependency().IsResolved(ctx, cronJobEntrypoint{kubeClient: kubeClient})
		Expect(resolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
		Expect(gets).To(Equal(1))
		Expect(lists).To(Equal(1))
	})

	It("does not query Jobs before the CronJob matches the target update", func() {
		mockClient := mocks.NewClient()
		mockClient.JobInterface = nil
		mockClient.CronJobInterface = cronJobAPI{get: func(context.Context, string, metav1.GetOptions) (*batchv1.CronJob, error) {
			snapshot := targetCronJob()
			snapshot.Labels = nil
			return snapshot, nil
		}}
		resolved, err := targetDependency().IsResolved(context.Background(), cronJobEntrypoint{kubeClient: mockClient})
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring("pending target update")))
	})

	DescribeTable("preserves retryable API errors at every stage",
		func(stage string, apiError error) {
			mockClient := mocks.NewClient()
			mockClient.CronJobInterface = cronJobAPI{get: func(context.Context, string, metav1.GetOptions) (*batchv1.CronJob, error) {
				if stage == "get" {
					return nil, apiError
				}
				return targetCronJob(), nil
			}}
			mockClient.JobInterface = jobAPI{list: func(context.Context, metav1.ListOptions) (*batchv1.JobList, error) {
				if stage == "list" {
					return nil, apiError
				}
				return &batchv1.JobList{}, nil
			}}
			resolved, err := targetDependency().IsResolved(context.Background(), cronJobEntrypoint{kubeClient: mockClient})
			Expect(resolved).To(BeFalse())
			Expect(errors.Is(err, apiError)).To(BeTrue())
			resolved, err = targetDependency().IsResolved(context.Background(), snapshotEntrypoint(targetCronJob(), nil))
			Expect(resolved).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		},
		Entry("not yet created", "get", apierrors.NewNotFound(schema.GroupResource{Group: "batch", Resource: "cronjobs"}, "cleanup")),
		Entry("CronJob permission denied", "get", apierrors.NewForbidden(schema.GroupResource{Group: "batch", Resource: "cronjobs"}, "cleanup", errors.New("permission denied"))),
		Entry("Job list unavailable", "list", errors.New("API unavailable")),
		Entry("Job list context canceled", "list", context.Canceled),
	)

	DescribeTable("rejects missing API snapshots rather than accepting success",
		func(stage, diagnostic string) {
			mockClient := mocks.NewClient()
			mockClient.CronJobInterface = cronJobAPI{get: func(context.Context, string, metav1.GetOptions) (*batchv1.CronJob, error) {
				if stage == "get" {
					return nil, nil
				}
				return targetCronJob(), nil
			}}
			mockClient.JobInterface = jobAPI{list: func(context.Context, metav1.ListOptions) (*batchv1.JobList, error) {
				return nil, nil
			}}
			resolved, err := targetDependency().IsResolved(context.Background(), cronJobEntrypoint{kubeClient: mockClient})
			Expect(resolved).To(BeFalse())
			Expect(err).To(MatchError(ContainSubstring(diagnostic)))
		},
		Entry("CronJob", "get", "CronJob was not returned"),
		Entry("Job list", "list", "Jobs were not returned"),
	)
})

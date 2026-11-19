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

package deployment

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	appsclient "k8s.io/client-go/kubernetes/typed/apps/v1"

	"opendev.org/airship/kubernetes-entrypoint/client"
	"opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/mocks"
)

const releaseAnnotation = "example.com/release-id"

// Embedded interfaces ensure any accidental API call other than Get fails.
type deploymentAPI struct {
	appsclient.DeploymentInterface
	get func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error)
}

func (a deploymentAPI) Get(ctx context.Context, name string, options metav1.GetOptions) (*appsv1.Deployment, error) {
	return a.get(ctx, name, options)
}

type namespacedDeploymentClient struct {
	client.ClientInterface
	deployments func(string) appsclient.DeploymentInterface
}

func (c namespacedDeploymentClient) Deployments(namespace string) appsclient.DeploymentInterface {
	return c.deployments(namespace)
}

type deploymentEntrypoint struct {
	entrypoint.EntrypointInterface
	kubeClient client.ClientInterface
}

func (e deploymentEntrypoint) Client() client.ClientInterface {
	return e.kubeClient
}

func completeDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "api", Namespace: "application", Generation: 7,
			Annotations: map[string]string{releaseAnnotation: "my-release:42"},
		},
		Spec: appsv1.DeploymentSpec{Replicas: new(int32(3))},
		Status: appsv1.DeploymentStatus{
			ObservedGeneration: 7, UpdatedReplicas: 3, Replicas: 3, AvailableReplicas: 3, ReadyReplicas: 3,
		},
	}
}

var _ = Describe("Deployment snapshot evaluation", func() {
	var snapshot *appsv1.Deployment
	BeforeEach(func() {
		snapshot = completeDeployment()
	})

	DescribeTable("requires convergence of the currently observed specification",
		func(change func(*appsv1.Deployment), resolved bool, diagnostic string) {
			change(snapshot)
			actual, err := evaluateDeployment(snapshot, nil)
			Expect(actual).To(Equal(resolved))
			if resolved {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(ContainSubstring(diagnostic)))
			}
		},
		Entry("complete without conditions", func(d *appsv1.Deployment) {}, true, ""),
		Entry("observed generation can be higher", func(d *appsv1.Deployment) { d.Status.ObservedGeneration++ }, true, ""),
		Entry("stale successful status", func(d *appsv1.Deployment) { d.Generation++ }, false, "has not been observed"),
		Entry("not all replicas updated", func(d *appsv1.Deployment) { d.Status.UpdatedReplicas = 2 }, false, "rollout incomplete"),
		Entry("too many updated replicas", func(d *appsv1.Deployment) { d.Status.UpdatedReplicas = 4 }, false, "rollout incomplete"),
		Entry("old or surge replicas remain", func(d *appsv1.Deployment) { d.Status.Replicas = 4 }, false, "rollout incomplete"),
		Entry("total replicas too few", func(d *appsv1.Deployment) { d.Status.Replicas = 2 }, false, "rollout incomplete"),
		Entry("ready is insufficient until minReadySeconds", func(d *appsv1.Deployment) {
			d.Spec.MinReadySeconds = 30
			d.Status.AvailableReplicas = 2
		}, false, "rollout incomplete"),
		Entry("surplus available replicas", func(d *appsv1.Deployment) { d.Status.AvailableReplicas = 4 }, false, "rollout incomplete"),
		Entry("paused even if counts complete", func(d *appsv1.Deployment) { d.Spec.Paused = true }, false, "paused"),
		Entry("deleting even if counts complete", func(d *appsv1.Deployment) { now := metav1.Now(); d.DeletionTimestamp = &now }, false, "being deleted"),
		Entry("complete with Progressing=True", func(d *appsv1.Deployment) {
			d.Status.Conditions = []appsv1.DeploymentCondition{{Type: appsv1.DeploymentProgressing, Status: corev1.ConditionTrue}}
		}, true, ""),
		Entry("Available and previous success cannot bypass counts", func(d *appsv1.Deployment) {
			d.Status.UpdatedReplicas = 1
			d.Status.Conditions = []appsv1.DeploymentCondition{
				{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionTrue},
				{Type: appsv1.DeploymentProgressing, Status: corev1.ConditionTrue, Reason: "NewReplicaSetAvailable"},
			}
		}, false, "rollout incomplete"),
		Entry("zero replicas converged", func(d *appsv1.Deployment) {
			d.Spec.Replicas = new(int32(0))
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation}
		}, true, ""),
		Entry("zero replicas require observation", func(d *appsv1.Deployment) {
			d.Spec.Replicas = new(int32(0))
			d.Status = appsv1.DeploymentStatus{}
		}, false, "has not been observed"),
		Entry("zero replicas still have old replicas", func(d *appsv1.Deployment) {
			d.Spec.Replicas = new(int32(0))
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation, Replicas: 1}
		}, false, "rollout incomplete"),
		Entry("zero replicas still have updated replicas", func(d *appsv1.Deployment) {
			d.Spec.Replicas = new(int32(0))
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation, UpdatedReplicas: 1}
		}, false, "rollout incomplete"),
		Entry("zero replicas still have available replicas", func(d *appsv1.Deployment) {
			d.Spec.Replicas = new(int32(0))
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation, AvailableReplicas: 1}
		}, false, "rollout incomplete"),
		Entry("zero replicas still subject to pause guard", func(d *appsv1.Deployment) {
			d.Spec.Replicas = new(int32(0))
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation}
			d.Spec.Paused = true
		}, false, "paused"),
		Entry("nil replicas default to one", func(d *appsv1.Deployment) {
			d.Spec.Replicas = nil
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation, UpdatedReplicas: 1, Replicas: 1, AvailableReplicas: 1}
		}, true, ""),
		Entry("nil replicas do not accept empty status", func(d *appsv1.Deployment) {
			d.Spec.Replicas = nil
			d.Status = appsv1.DeploymentStatus{ObservedGeneration: d.Generation}
		}, false, "rollout incomplete"),
		Entry("terminating Pods need not disappear", func(d *appsv1.Deployment) { d.Status.TerminatingReplicas = new(int32(2)) }, true, ""),
		Entry("RollingUpdate uses the same predicate", func(d *appsv1.Deployment) { d.Spec.Strategy.Type = appsv1.RollingUpdateDeploymentStrategyType }, true, ""),
		Entry("Recreate uses the same predicate", func(d *appsv1.Deployment) { d.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType }, true, ""),
		Entry("Recreate cannot bypass incomplete counts", func(d *appsv1.Deployment) {
			d.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType
			d.Status.UpdatedReplicas = 0
		}, false, "rollout incomplete"),
	)

	DescribeTable("matches every expected top-level annotation exactly",
		func(actual, template, expected map[string]string, resolved bool, diagnostic string) {
			snapshot.Annotations = actual
			snapshot.Spec.Template.Annotations = template
			result, err := evaluateDeployment(snapshot, expected)
			Expect(result).To(Equal(resolved))
			if resolved {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(ContainSubstring(diagnostic)))
			}
		},
		Entry("no expectations", nil, nil, nil, true, ""),
		Entry("empty expectations", nil, nil, map[string]string{}, true, ""),
		Entry("all matching, actual extras allowed", map[string]string{"a": "1", "b": "2", "extra": "3"}, nil, map[string]string{"a": "1", "b": "2"}, true, ""),
		Entry("one match does not substitute for another", map[string]string{"a": "1", "b": "wrong"}, nil, map[string]string{"a": "1", "b": "2"}, false, `annotation "b" is "wrong", expected "2"`),
		Entry("missing empty value", nil, nil, map[string]string{"empty": ""}, false, `annotation "empty" is missing, expected ""`),
		Entry("present empty value", map[string]string{"empty": ""}, nil, map[string]string{"empty": ""}, true, ""),
		Entry("Pod template cannot substitute for missing metadata", nil, map[string]string{"release": "42"}, map[string]string{"release": "42"}, false, "is missing"),
		Entry("Pod template cannot substitute for wrong metadata", map[string]string{"release": "41"}, map[string]string{"release": "42"}, map[string]string{"release": "42"}, false, `is "41", expected "42"`),
		Entry("a newer marker is not silently accepted", map[string]string{"release": "43"}, nil, map[string]string{"release": "42"}, false, `is "43", expected "42"`),
	)

	It("reports only an incomplete rollout regardless of controller failures, then accepts recovery", func() {
		snapshot.Status.UpdatedReplicas = 2
		snapshot.Status.Conditions = []appsv1.DeploymentCondition{
			{Type: appsv1.DeploymentProgressing, Status: corev1.ConditionFalse, Reason: "ProgressDeadlineExceeded", Message: "replicas did not progress"},
			{Type: appsv1.DeploymentReplicaFailure, Status: corev1.ConditionTrue, Reason: "FailedCreate", Message: "quota exceeded"},
		}
		resolved, err := evaluateDeployment(snapshot, nil)
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError("rollout incomplete"))
		snapshot.Status.UpdatedReplicas = 3
		resolved, err = evaluateDeployment(snapshot, nil)
		Expect(resolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	It("evaluates deletion, annotations, observation and pause in order", func() {
		snapshot.Spec.Paused = true
		snapshot.Generation++
		expected := map[string]string{releaseAnnotation: "my-release:43"}
		now := metav1.Now()
		snapshot.DeletionTimestamp = &now
		_, err := evaluateDeployment(snapshot, expected)
		Expect(err).To(MatchError(ContainSubstring("being deleted")))
		snapshot.DeletionTimestamp = nil
		_, err = evaluateDeployment(snapshot, expected)
		Expect(err).To(MatchError(ContainSubstring("pending target update")))
		snapshot.Annotations = expected
		_, err = evaluateDeployment(snapshot, expected)
		Expect(err).To(MatchError(ContainSubstring("has not been observed")))
		snapshot.Status.ObservedGeneration = snapshot.Generation
		_, err = evaluateDeployment(snapshot, expected)
		Expect(err).To(MatchError(ContainSubstring("paused")))
	})
})

var _ = Describe("Deployment resolver", func() {
	It("constructs a named dependency and owns its annotation expectations", func() {
		annotations := map[string]string{releaseAnnotation: "42"}
		dep := NewDeployment("api", "application", annotations)
		Expect(dep.name).To(Equal("api"))
		Expect(dep.namespace).To(Equal("application"))
		Expect(dep.String()).To(Equal("Deployment api in namespace application"))
		annotations[releaseAnnotation] = "43"
		Expect(dep.annotations).To(Equal(map[string]string{releaseAnnotation: "42"}))
		Expect(NewDeployment("", "application", nil)).To(BeNil())
		Expect(NewDeployment(" ", "application", nil)).To(BeNil())
	})

	DescribeTable("registers both configuration modes using the existing resolver contract",
		func(plain, json string, expected []Deployment) {
			GinkgoT().Setenv("NAMESPACE", "application")
			GinkgoT().Setenv("DEPENDENCY_DEPLOYMENT", plain)
			GinkgoT().Setenv("DEPENDENCY_DEPLOYMENT_JSON", json)
			var registered []Deployment
			registerDeployments(func(resolver entrypoint.Resolver) {
				Expect(resolver).To(BeAssignableToTypeOf(Deployment{}))
				registered = append(registered, resolver.(Deployment))
			})
			Expect(registered).To(Equal(expected))
		},
		Entry("plain", "api,other:worker", "", []Deployment{{name: "api", namespace: "application"}, {name: "worker", namespace: "other"}}),
		Entry("JSON precedence with annotations", "ignored", `[{"name":"api","annotations":{"release":"42"}},{"name":"worker","namespace":"other"}]`, []Deployment{
			{name: "api", namespace: "application", annotations: map[string]string{"release": "42"}},
			{name: "worker", namespace: "other", annotations: map[string]string{}},
		}),
		Entry("unset", "", "", nil),
		Entry("invalid names skipped", "", `[{"annotations":{"release":"42"}}]`, nil),
	)

	It("fetches only the configured Deployment, once, with the supplied context", func() {
		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)
		calls := 0
		api := deploymentAPI{get: func(actual context.Context, name string, options metav1.GetOptions) (*appsv1.Deployment, error) {
			calls++
			Expect(actual).To(BeIdenticalTo(ctx))
			Expect(name).To(Equal("api"))
			Expect(options).To(Equal(metav1.GetOptions{}))
			return completeDeployment(), nil
		}}
		mockClient := mocks.NewClient()
		mockClient.DeploymentInterface = api
		kubeClient := namespacedDeploymentClient{ClientInterface: mockClient, deployments: func(namespace string) appsclient.DeploymentInterface {
			Expect(namespace).To(Equal("application"))
			return mockClient.Deployments(namespace)
		}}
		ep := deploymentEntrypoint{kubeClient: kubeClient}
		resolved, err := NewDeployment("api", "application", nil).IsResolved(ctx, ep)
		Expect(resolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
		Expect(calls).To(Equal(1))
	})

	DescribeTable("keeps missing resources and API failures retryable and preserves their cause",
		func(apiError error) {
			mockClient := mocks.NewClient()
			mockClient.DeploymentInterface = deploymentAPI{get: func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error) {
				return nil, apiError
			}}
			ep := deploymentEntrypoint{kubeClient: mockClient}
			dep := NewDeployment("api", "application", nil)
			resolved, err := dep.IsResolved(context.Background(), ep)
			Expect(resolved).To(BeFalse())
			Expect(err).To(MatchError(ContainSubstring("getting Deployment api in namespace application")))
			Expect(errors.Is(err, apiError)).To(BeTrue())
			mockClient.DeploymentInterface = deploymentAPI{get: func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error) {
				return completeDeployment(), nil
			}}
			resolved, err = dep.IsResolved(context.Background(), ep)
			Expect(resolved).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		},
		Entry("not yet created", apierrors.NewNotFound(schema.GroupResource{Group: "apps", Resource: "deployments"}, "api")),
		Entry("permission denied", apierrors.NewForbidden(schema.GroupResource{Group: "apps", Resource: "deployments"}, "api", errors.New("permission denied"))),
		Entry("API unavailable", errors.New("API unavailable")),
		Entry("context canceled", context.Canceled),
	)

	It("does not accept a missing snapshot even if a client returns no error", func() {
		resolved, err := evaluateDeployment(nil, nil)
		Expect(resolved).To(BeFalse())
		Expect(err).To(MatchError(ContainSubstring("not returned by the API")))
	})

	DescribeTable("handles release sequencing without sleeps or requiring a witnessed transition",
		func(annotationOnly, startAfterCompletion bool) {
			old := completeDeployment()
			old.Annotations[releaseAnnotation] = "my-release:41"
			submitted := old.DeepCopy()
			submitted.Annotations[releaseAnnotation] = "my-release:42"
			submitted.Generation++
			if !annotationOnly {
				submitted.Spec.Template.Spec.Containers = []corev1.Container{{Name: "api", Image: "api:v42"}}
			}
			observed := submitted.DeepCopy()
			observed.Status.ObservedGeneration = observed.Generation
			if !annotationOnly {
				observed.Status.UpdatedReplicas = 1
				observed.Status.Replicas = 4
				observed.Status.AvailableReplicas = 3
			}
			completed := observed.DeepCopy()
			completed.Status.UpdatedReplicas = 3
			completed.Status.Replicas = 3
			completed.Status.AvailableReplicas = 3
			snapshots := []*appsv1.Deployment{old, submitted, observed, completed}
			expected := []bool{false, false, annotationOnly, true}
			diagnostics := []string{"pending target update", "has not been observed", "rollout incomplete", ""}
			if startAfterCompletion {
				snapshots = []*appsv1.Deployment{completed}
				expected = []bool{true}
			}
			calls := 0
			mockClient := mocks.NewClient()
			mockClient.DeploymentInterface = deploymentAPI{get: func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error) {
				Expect(calls).To(BeNumerically("<", len(snapshots)))
				snapshot := snapshots[calls].DeepCopy()
				calls++
				return snapshot, nil
			}}
			ep := deploymentEntrypoint{kubeClient: mockClient}
			dep := NewDeployment("api", "application", map[string]string{releaseAnnotation: "my-release:42"})
			for i, want := range expected {
				resolved, err := dep.IsResolved(context.Background(), ep)
				Expect(resolved).To(Equal(want))
				if want {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(ContainSubstring(diagnostics[i])))
				}
			}
			Expect(calls).To(Equal(len(snapshots)))
			if annotationOnly {
				Expect(completed.Spec.Template).To(Equal(old.Spec.Template))
			}
		},
		Entry("template-changing upgrade", false, false),
		Entry("annotation-only upgrade", true, false),
		Entry("Job starts after template rollout completed", false, true),
		Entry("Job starts after annotation-only upgrade completed", true, true),
	)

	It("does not accept another update arriving while waiting", func() {
		pending := completeDeployment()
		pending.Status.UpdatedReplicas = 1
		superseded := completeDeployment()
		superseded.Generation++
		superseded.Annotations[releaseAnnotation] = "my-release:43"
		acknowledged := superseded.DeepCopy()
		acknowledged.Status.ObservedGeneration = acknowledged.Generation
		reusedMarker := superseded.DeepCopy()
		reusedMarker.Annotations[releaseAnnotation] = "my-release:42"
		snapshots := []*appsv1.Deployment{pending, superseded, acknowledged, reusedMarker}
		diagnostics := []string{"rollout incomplete", "pending target update", "pending target update", "has not been observed"}
		calls := 0
		mockClient := mocks.NewClient()
		mockClient.DeploymentInterface = deploymentAPI{get: func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error) {
			snapshot := snapshots[calls].DeepCopy()
			calls++
			return snapshot, nil
		}}
		dep := NewDeployment("api", "application", map[string]string{releaseAnnotation: "my-release:42"})
		for _, diagnostic := range diagnostics {
			resolved, err := dep.IsResolved(context.Background(), deploymentEntrypoint{kubeClient: mockClient})
			Expect(resolved).To(BeFalse())
			Expect(err).To(MatchError(ContainSubstring(diagnostic)))
		}
		Expect(calls).To(Equal(len(snapshots)))
	})
})

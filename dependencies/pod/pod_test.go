package pod

import (
	"fmt"
	"os"

	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	podEnvVariableValue = "podlist"
	podNamespace        = "test"
	requireSameNode     = true
)

var testEntrypoint entrypoint.EntrypointInterface
var testLabels = map[string]string{"foo": "bar"}

var _ = Describe("Pod", func() {

	BeforeEach(func() {
		err := os.Setenv(PodNameEnvVar, podEnvVariableValue)
		Expect(err).NotTo(HaveOccurred())

		testEntrypoint = mocks.NewEntrypoint()
	})

	It(fmt.Sprintf("checks failure of new pod creation without %s set", PodNameEnvVar), func() {
		os.Unsetenv(PodNameEnvVar)
		pod, err := NewPod(testLabels, podNamespace, requireSameNode)

		Expect(pod).To(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(PodNameNotSetErrorFormat, podNamespace)))
	})

	It(fmt.Sprintf("creates new pod with %s set and checks its name", PodNameEnvVar), func() {
		pod, err := NewPod(testLabels, podNamespace, requireSameNode)
		Expect(pod).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(pod.labels).To(Equal(testLabels))
	})

	It("is resolved via all pods matching labels ready on same host", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.SameHostReadyMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	It("is resolved via some pods matching labels ready on same host", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.SameHostSomeReadyMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	It("is not resolved via a pod matching labels not ready on same host", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.SameHostNotReadyMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("is not resolved via pod matching labels ready on different host", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.DifferentHostReadyMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("is resolved via pod matching labels ready on different host when requireSameNode=false", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.DifferentHostReadyMatchLabel}, podNamespace, false)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	It("is not resolved via pod matching labels not ready on different host when requireSameNode=false", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.DifferentHostNotReadyMatchLabel}, podNamespace, false)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("is not resolved via no pods matching labels", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.NoPodsMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("is not resolved when getting pods matching labels from api fails", func() {
		pod, _ := NewPod(map[string]string{"name": mocks.FailingMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It(fmt.Sprintf("is not resolved when getting current pod via %s value fails", PodNameEnvVar), func() {
		os.Setenv(PodNameEnvVar, mocks.PodNotPresent)
		pod, _ := NewPod(map[string]string{"name": mocks.SameHostReadyMatchLabel}, podNamespace, requireSameNode)

		isResolved, err := pod.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})
})

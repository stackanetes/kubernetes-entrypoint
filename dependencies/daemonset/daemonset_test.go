package daemonset

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
	daemonsetNamespace  = "test"
)

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Daemonset", func() {

	BeforeEach(func() {
		err := os.Setenv(PodNameEnvVar, podEnvVariableValue)
		Expect(err).NotTo(HaveOccurred())

		testEntrypoint = mocks.NewEntrypoint()
	})

	It(fmt.Sprintf("checks failure of new daemonset creation without %s set", PodNameEnvVar), func() {
		os.Unsetenv(PodNameEnvVar)
		daemonset, err := NewDaemonset(mocks.SucceedingDaemonsetName, daemonsetNamespace)

		Expect(daemonset).To(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(PodNameNotSetErrorFormat, mocks.SucceedingDaemonsetName, daemonsetNamespace)))
	})

	It(fmt.Sprintf("creates new daemonset with %s set and checks its name", PodNameEnvVar), func() {
		daemonset, err := NewDaemonset(mocks.SucceedingDaemonsetName, daemonsetNamespace)
		Expect(daemonset).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(daemonset.name).To(Equal(mocks.SucceedingDaemonsetName))
	})

	It("checks resolution of a succeeding daemonset", func() {
		daemonset, _ := NewDaemonset(mocks.SucceedingDaemonsetName, daemonsetNamespace)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks resolution failure of a daemonset with incorrect name", func() {
		daemonset, _ := NewDaemonset(mocks.FailingDaemonsetName, daemonsetNamespace)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("checks resolution failure of a daemonset with incorrect match labels", func() {
		daemonset, _ := NewDaemonset(mocks.IncorrectMatchLabelsDaemonsetName, daemonsetNamespace)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It(fmt.Sprintf("checks resolution failure of a daemonset with incorrect %s value", PodNameEnvVar), func() {
		// Set POD_NAME to value not present in the mocks
		os.Setenv(PodNameEnvVar, mocks.PodNotPresent)
		daemonset, _ := NewDaemonset(mocks.IncorrectMatchLabelsDaemonsetName, daemonsetNamespace)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("checks resolution failure of a daemonset with none of the pods with Ready status", func() {
		daemonset, _ := NewDaemonset(mocks.NotReadyMatchLabelsDaemonsetName, daemonsetNamespace)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("checks resolution of a correct daemonset namespace", func() {
		daemonset, err := NewDaemonset(mocks.CorrectNamespaceDaemonsetName, daemonsetNamespace)

		Expect(daemonset).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())

	})

	It("checks resolution of an incorrect daemonset namespace", func() {
		daemonset, err := NewDaemonset(mocks.IncorrectNamespaceDaemonsetName, daemonsetNamespace)

		Expect(daemonset).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(BeFalse())
		Expect(err).To(HaveOccurred())
	})

	It("resolve daemonset and entrypoint pod in different namespaces", func() {
		daemonset, err := NewDaemonset(mocks.CorrectNamespaceDaemonsetName, mocks.CorrectDaemonsetNamespace)
		Expect(err).NotTo(HaveOccurred())

		err = os.Setenv(PodNameEnvVar, "shouldwork")
		Expect(err).NotTo(HaveOccurred())

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(err).NotTo(HaveOccurred())
		Expect(isResolved).To(BeTrue())
		err = os.Unsetenv(PodNameEnvVar)
	})
})

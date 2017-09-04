package daemonset_test

import (
	"fmt"
	"os"

	. "github.com/stackanetes/kubernetes-entrypoint/dependencies/daemonset"
	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	podEnvVariableValue = "podlist"
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
		daemonset, err := NewDaemonset(mocks.SucceedingDaemonsetName)

		Expect(daemonset).To(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(PodNameNotSetErrorFormat, mocks.SucceedingDaemonsetName)))
	})

	It(fmt.Sprintf("creates new daemonset with %s set and checks its name", PodNameEnvVar), func() {
		daemonset, err := NewDaemonset(mocks.SucceedingDaemonsetName)
		Expect(daemonset).NotTo(Equal(nil))
		Expect(err).NotTo(HaveOccurred())

		Expect(daemonset.GetName()).To(Equal(mocks.SucceedingDaemonsetName))
	})

	It("checks resolution of a succeeding daemonset", func() {
		daemonset, _ := NewDaemonset(mocks.SucceedingDaemonsetName)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks resolution failure of a daemonset with incorrect name", func() {
		daemonset, _ := NewDaemonset(mocks.FailingDaemonsetName)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
	})

	It("checks resolution failure of a daemonset with incorrect match labels", func() {
		daemonset, _ := NewDaemonset(mocks.IncorrectMatchLabelsDaemonsetName)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
	})

	It(fmt.Sprintf("checks resolution failure of a daemonset with incorrect %s value", PodNameEnvVar), func() {
		// Set POD_NAME to value not present in the mocks
		os.Setenv(PodNameEnvVar, mocks.PodNotPresent)
		daemonset, _ := NewDaemonset(mocks.IncorrectMatchLabelsDaemonsetName)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
	})

	It("checks resolution failure of a daemonset with none of the pods with Ready status", func() {
		daemonset, _ := NewDaemonset(mocks.NotReadyMatchLabelsDaemonsetName)

		isResolved, err := daemonset.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
	})
})

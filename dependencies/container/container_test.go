package container_test

import (
	"fmt"
	"os"

	. "github.com/stackanetes/kubernetes-entrypoint/dependencies/container"
	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	podEnvVariableName = "POD_NAME"
)

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Container", func() {

	BeforeEach(func() {
		err := os.Setenv(podEnvVariableName, mocks.PodEnvVariableValue)
		Expect(err).NotTo(HaveOccurred())

		testEntrypoint = mocks.NewEntrypoint()
	})

	It("checks the name of a newly created container", func() {
		container := NewContainer(mocks.MockContainerName)

		Expect(container.GetName()).To(Equal(mocks.MockContainerName))
	})

	It(fmt.Sprintf("checks container resolution failure with %s not set", podEnvVariableName), func() {
		os.Unsetenv(podEnvVariableName)
		container := NewContainer(mocks.MockContainerName)

		isResolved, err := container.IsResolved(testEntrypoint)
		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(PodNameNotSetError))
	})

	It("checks resolution of a succeeding container", func() {
		container := NewContainer(mocks.MockContainerName)

		isResolved, err := container.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It(fmt.Sprintf("fails to resolve a mocked container for a given %s value", podEnvVariableName), func() {
		err := os.Setenv(podEnvVariableName, "INVALID_POD_LIST_VALUE")
		Expect(err).NotTo(HaveOccurred())

		container := NewContainer(mocks.PodNotPresent)
		Expect(container).NotTo(Equal(nil))

		var isResolved bool
		isResolved, err = container.IsResolved(testEntrypoint)
		Expect(isResolved).To(Equal(false))
		Expect(err).To(BeNil())
	})
})

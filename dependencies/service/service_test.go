package service_test

import (
	"fmt"

	. "github.com/stackanetes/kubernetes-entrypoint/dependencies/service"
	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testServiceName = "TEST_SERVICE"

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Service", func() {

	BeforeEach(func() {
		testEntrypoint = mocks.NewEntrypoint()
	})

	It("checks the name of a newly created service", func() {
		service := NewService(testServiceName)

		Expect(service.GetName()).To(Equal(testServiceName))
	})

	It("checks resolution of a succeeding service", func() {
		service := NewService(mocks.SucceedingServiceName)

		isResolved, err := service.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks resolution failure of a failing service", func() {
		service := NewService(mocks.FailingServiceName)

		isResolved, err := service.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(mocks.MockEndpointError))
	})

	It("checks resolution failure of a succeeding service with removed subsets", func() {
		service := NewService(mocks.EmptySubsetsServiceName)

		isResolved, err := service.IsResolved(testEntrypoint)
		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(fmt.Sprintf(FailingStatusFormat, service.GetName())))
	})
})

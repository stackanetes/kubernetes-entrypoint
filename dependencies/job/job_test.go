package job_test

import (
	"fmt"

	. "github.com/stackanetes/kubernetes-entrypoint/dependencies/job"
	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testJobName = "TEST_JOB"

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Job", func() {

	BeforeEach(func() {
		testEntrypoint = mocks.NewEntrypoint()
	})

	It("checks the name of a newly created job", func() {
		job := NewJob(testJobName)

		Expect(job.GetName()).To(Equal(testJobName))
	})

	It("checks resolution of a succeeding job", func() {
		job := NewJob(mocks.SucceedingJobName)

		isResolved, err := job.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks resolution failure of a failing job", func() {
		job := NewJob(mocks.FailingJobName)

		isResolved, err := job.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(fmt.Sprintf(FailingStatusFormat, job.GetName())))
	})
})

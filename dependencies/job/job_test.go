package job

import (
	"fmt"

	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testJobName = "TEST_JOB_NAME"
const testJobNamespace = "TEST_JOB_NAMESPACE"

var testLabels = map[string]string{
	"k1": "v1",
}

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Job", func() {

	BeforeEach(func() {
		testEntrypoint = mocks.NewEntrypoint()
	})

	It("constructor correctly assigns fields", func() {
		nameJob := NewJob(testJobName, testJobNamespace, nil)

		Expect(nameJob.name).To(Equal(testJobName))
		Expect(nameJob.namespace).To(Equal(testJobNamespace))

		labelsJob := NewJob("", testJobNamespace, testLabels)

		Expect(labelsJob.labels).To(Equal(testLabels))
	})

	It("constructor returns nil when both name and labels specified", func() {
		job := NewJob(testJobName, testJobNamespace, testLabels)

		Expect(job).To(BeNil())
	})

	It("checks resolution of a succeeding job by name", func() {
		job := NewJob(mocks.SucceedingJobName, mocks.SucceedingJobName, nil)

		isResolved, err := job.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks resolution failure of a failing job by name", func() {
		job := NewJob(mocks.FailingJobName, mocks.FailingJobName, nil)

		isResolved, err := job.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(fmt.Sprintf(FailingStatusFormat, job)))
	})

	It("checks resolution of a succeeding job by labels", func() {
		job := NewJob("", mocks.SucceedingJobName, map[string]string{"name": mocks.SucceedingJobLabel})

		isResolved, err := job.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks resolution failure of a failing job by labels", func() {
		job := NewJob("", mocks.FailingJobName, map[string]string{"name": mocks.FailingJobLabel})

		isResolved, err := job.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(fmt.Sprintf(FailingStatusFormat, job)))
	})

})

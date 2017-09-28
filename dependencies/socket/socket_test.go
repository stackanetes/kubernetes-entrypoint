package socket

import (
	"fmt"
	"os"

	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	existingSocketPath    = "/tmp/k8s-existing-socket"
	nonExistingSocketPath = "/tmp/k8s-nonexisting-socket"
	noPermsSocketPath     = "/root/k8s-no-permission-socket"
)

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("Socket", func() {

	BeforeEach(func() {
		testEntrypoint = mocks.NewEntrypoint()

		_, err := os.Create(existingSocketPath)
		Expect(err).NotTo(HaveOccurred())

	})

	It("checks the name of a newly created socket", func() {
		socket := NewSocket(existingSocketPath)

		Expect(socket.name).To(Equal(existingSocketPath))
	})

	It("resolves an existing socket socket", func() {
		socket := NewSocket(existingSocketPath)

		isResolved, err := socket.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

	It("fails on trying to resolve a nonexisting socket", func() {
		socket := NewSocket(nonExistingSocketPath)

		isResolved, err := socket.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf(NonExistingErrorFormat, socket)))
	})

	It("fails on trying to resolve a socket without permissions", func() {
		socket := NewSocket(noPermsSocketPath)

		isResolved, err := socket.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf(NoPermsErrorFormat, socket)))
	})
})

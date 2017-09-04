package socket_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSocket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Socket Suite")
}

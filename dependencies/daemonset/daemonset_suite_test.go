package daemonset_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDaemonset(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Daemonset Suite")
}

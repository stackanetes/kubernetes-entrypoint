package pod_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPod(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pod Suite")
}

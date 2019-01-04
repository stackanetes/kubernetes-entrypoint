package customresource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCustomResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Custom Resource Suite")
}

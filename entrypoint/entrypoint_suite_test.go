package entrypoint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEntrypoint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entrypoint Suite")
}

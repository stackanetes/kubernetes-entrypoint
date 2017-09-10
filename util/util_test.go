package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const failingNamespaceUtil = "foo:util"

var _ = Describe("Util", func() {

	It("fails on trying to resolve a socket with namespace", func() {
		contains := ContainsSeparator(failingNamespaceUtil, "Util")
		Expect(contains).To(Equal(true))
	})
})

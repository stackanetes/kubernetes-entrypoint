package client

import (
	"fmt"

	"k8s.io/client-go/kubernetes/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	nonExistingGroupVersion       = "nonexisting"
	nonExistingGroupVersionFormat = "GroupVersion \"%s\" not found"
)

var testClient *Client

func newTestClient() *Client {
	client := Client{}
	client.Clientset = fake.NewSimpleClientset()
	return &client
}

var _ = Describe("Client", func() {

	BeforeEach(func() {
		testClient = newTestClient()
	})

	It("asserts that tests are working", func() {
		Expect(true).To(BeTrue())
	})

	It("checks that GetResourceName fails when the GroupVersion doesn't exist", func() {
		name, err := testClient.GetResourceName("", nonExistingGroupVersion)

		Expect(name).To(BeEmpty())
		Expect(err.Error()).To(Equal(fmt.Sprintf(nonExistingGroupVersionFormat, nonExistingGroupVersion)))
	})
})

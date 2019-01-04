package customresource

import (
	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"
	env "github.com/stackanetes/kubernetes-entrypoint/util/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	expectedFailure              = "Expected Failure"
	testCustomResourceName       = "foo1"
	testCustomResourceNamespace  = "foospace1"
	testCustomResourceAPIVersion = "api1"
	testCustomResourceKind       = "kind1"
)

var testEntrypoint entrypoint.EntrypointInterface

var _ = Describe("CustomResource", func() {

	BeforeEach(func() {
		testEntrypoint = mocks.NewEntrypoint()
	})

	It("checks that constructor correctly assigns fields", func() {
		nameCustomResource := NewCustomResource(&env.CustomResourceDependency{
			ApiVersion: "api1",
			Kind:       "kind1",
			Namespace:  "foospace1",
			Name:       "foo1",
			Fields: []map[string]string{
				{
					"key":   "field1key1",
					"value": "field1val1",
				},
				{
					"key":   "field1key2",
					"value": "field1val2",
				},
			},
		})

		Expect(nameCustomResource.APIVersion).To(Equal(testCustomResourceAPIVersion))
		Expect(nameCustomResource.Kind).To(Equal(testCustomResourceKind))
		Expect(nameCustomResource.Name).To(Equal(testCustomResourceName))
		Expect(nameCustomResource.Namespace).To(Equal(testCustomResourceNamespace))
		Expect(nameCustomResource.Fields).To(Equal([]map[string]string{
			{
				"key":   "field1key1",
				"value": "field1val1",
			},
			{
				"key":   "field1key2",
				"value": "field1val2",
			},
		},
		))
	})

	It("checks that IsResolved returns true when customResource is resolved", func() {
		resolvedResource := CustomResource{
			Kind: "successKind",
			Name: "resolved",
			Fields: []map[string]string{
				{
					"key":   "simple_key",
					"value": "simple_value",
				},
				{
					"key":   "complex.key.with.layers",
					"value": "complex_value",
				},
			},
		}
		isResolved, err := resolvedResource.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).To(BeNil())
	})

	It("checks that IsResolved returns false when customResource is NOT resolved", func() {
		unresolvedResource := CustomResource{
			Kind: "successKind",
			Name: "unresolved",
			Fields: []map[string]string{
				{
					"key":   "key",
					"value": "expected_value",
				},
			},
		}
		isResolved, err := unresolvedResource.IsResolved(testEntrypoint)

		errMsg := "Expected value of [key] to be [expected_value], but got [unexpected_value]"
		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(errMsg))
	})

	It("checks resolution failure of a failing customResource on API requests", func() {
		failingResource := CustomResource{Kind: "failKind"}
		isResolved, err := failingResource.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(expectedFailure))

		successfulResource := CustomResource{
			Kind: "successKind",
			Name: "failName",
		}
		isResolved, err = successfulResource.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(false))
		Expect(err.Error()).To(Equal(expectedFailure))
	})

})

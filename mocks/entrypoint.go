package mocks

import (
	//	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	cli "github.com/stackanetes/kubernetes-entrypoint/client"
)

type MockEntrypoint struct {
	client    cli.ClientInterface
	namespace string
}

func (m MockEntrypoint) Resolve() {
}

func (m MockEntrypoint) Client() (client cli.ClientInterface) {
	return m.client
}
func (m MockEntrypoint) GetNamespace() (namespace string) {
	return m.namespace
}
func NewEntrypoint() MockEntrypoint {
	return MockEntrypoint{
		client:    NewClient(),
		namespace: "test",
	}
}

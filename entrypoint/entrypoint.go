package entrypoint

import (
	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	restclient "k8s.io/kubernetes/pkg/client/restclient"
	"os"
	"sync"
	"time"
)

var dependencies []Resolver // List containing all dependencies to be resolved
const (
	//DependencyPrefix is a prefix for env variables
	DependencyPrefix = "DEPENDENCY_"
	interval         = 2
)

type EntrypointInterface interface {
	Resolve()
	Client() cli.ClientInterface
	GetNamespace() string
}

// Entrypoint is a main struct which checks dependencies
type Entrypoint struct {
	client    cli.ClientInterface
	namespace string
}

//New is a constructor for entrypoint
func New(config *restclient.Config) (entry *Entrypoint, err error) {
	entry = new(Entrypoint)
	client, err := cli.New(config)
	if err != nil {
		return nil, err
	}
	entry.client = client
	if entry.namespace = os.Getenv("NAMESPACE"); entry.namespace == "" {
		logger.Warning.Print("NAMESPACE env not set, using default")
		entry.namespace = "default"
	}
	return entry, err
}

func (e Entrypoint) Client() (client cli.ClientInterface) {
	return e.client
}

func (e Entrypoint) GetNamespace() string {
	return e.namespace
}

//Resolve is a main loop which iterates through all dependencies and resolves them
func (e Entrypoint) Resolve() {
	var wg sync.WaitGroup
	for _, dep := range dependencies {
		wg.Add(1)
		go func(dep Resolver) {
			defer wg.Done()
			logger.Info.Printf("Resolving %s", dep.GetName())
			var err error
			status := false
			for status == false {
				if status, err = dep.IsResolved(e); err != nil {
					logger.Warning.Printf("Resolving dependency for %v failed: %v", dep.GetName(), err)
				}
				time.Sleep(interval * time.Second)
			}
			logger.Info.Printf("Dependency %v is resolved", dep.GetName())

		}(dep)
	}
	wg.Wait()

}

//Resolver is an interface which all dependencies should implement
type Resolver interface {
	IsResolved(entrypoint EntrypointInterface) (bool, error)
	GetName() string
}

//Register is a function which registers new dependencies
func Register(res Resolver) {
	if res == nil {
		panic("Entrypoint: could not register nil Resolver")
	}
	dependencies = append(dependencies, res)
}

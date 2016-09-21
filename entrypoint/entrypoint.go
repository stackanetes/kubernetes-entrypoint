package entrypoint

import (
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	//	"k8s.io/kubernetes/pkg/client/restclient"
	"fmt"
	cl "k8s.io/kubernetes/pkg/client/unversioned"
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

// Entrypoint is a main struct which check dependencies
type Entrypoint struct {
	Client    *cl.Client
	Namespace string
}

//New is a constructor for entrypoint
func New(client *cl.Client) (entry *Entrypoint, err error) {
	entry = new(Entrypoint)
	if entry.Client = client; client == nil {
		if entry.Client, err = cl.NewInCluster(); err != nil {
			err = fmt.Errorf("Error while creating k8s client: %s", err)
			return entry, err
		}
	}
	if entry.Namespace = os.Getenv("NAMESPACE"); entry.Namespace == "" {
		logger.Warning.Print("NAMESPACE env not set, using default")
		entry.Namespace = "default"
	}
	return entry, err
}

//Resolve is a main loop whic iterates through all dependencies and resolves them
func (e *Entrypoint) Resolve() {
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
	IsResolved(entrypoint *Entrypoint) (bool, error)
	GetName() string
}

//Register is a function which register new dependencies
func Register(res Resolver) {
	if res == nil {
		panic("Entrypoint: could not register nil Resolver")
	}
	dependencies = append(dependencies, res)
}

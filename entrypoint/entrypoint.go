package entrypoint

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"k8s.io/client-go/rest"
)

var dependencies []Resolver // List containing all dependencies to be resolved

const (
	//DependencyPrefix is a prefix for env variables
	DependencyPrefix      = "DEPENDENCY_"
	JsonSuffix            = "_JSON"
	resolverSleepInterval = 2
)

//Resolver is an interface which all dependencies should implement
type Resolver interface {
	IsResolved(entrypoint EntrypointInterface) (bool, error)
	GetDependency() map[string]interface{}
}

//Conflict is used to explain which dependencies are causing errors or a given dependency
type Conflict struct {
	Dependency map[string]interface{}
	resolved   bool
	Reason     string
	mux        sync.Mutex
}

//NewConflict creates a new Conflict object
func NewConflict(dep map[string]interface{}) *Conflict {
	return &Conflict{
		Dependency: dep,
		resolved:   false,
		Reason:     "Not yet checked",
		mux:        sync.Mutex{},
	}
}

//Update changes the state and cause of a conflict
func (c *Conflict) Update(resolved bool, reason string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.resolved = resolved
	c.Reason = reason
}

type EntrypointInterface interface {
	Resolve()
	Client() cli.ClientInterface
}

// Entrypoint is a main struct which checks dependencies
type Entrypoint struct {
	client    cli.ClientInterface
	namespace string
}

//Register is a function which registers new dependencies
func Register(res Resolver) {
	if res == nil {
		panic("Entrypoint: could not register nil Resolver")
	}
	dependencies = append(dependencies, res)
}

//New is a constructor for entrypoint
func New(config *rest.Config) (entry *Entrypoint, err error) {
	entry = new(Entrypoint)
	client, err := cli.New(config)
	if err != nil {
		return nil, err
	}
	entry.client = client
	return entry, err
}

func (e Entrypoint) Client() (client cli.ClientInterface) {
	return e.client
}

//Resolve is a main loop which iterates through all dependencies and resolves them
func (e Entrypoint) Resolve() {
	var wg sync.WaitGroup
	defer wg.Wait()
	conflicts := make(map[Resolver]*Conflict)
	for _, dep := range dependencies {
		wg.Add(1)
		conflicts[dep] = NewConflict(dep.GetDependency())
		go func(dep Resolver) {
			defer wg.Done()
			logger.Info("Resolving %v", dep)
			var err error
			status := false
			for status == false {
				if status, err = dep.IsResolved(e); err != nil {
					logger.Warning("Resolving dependency %s failed: %v .", dep, err)
					conflicts[dep].Update(false, err.Error())
				}
				if status == false {
					time.Sleep(resolverSleepInterval * time.Second)
				}
			}
			conflicts[dep].Update(true, "")
			logger.Info("Dependency %v is resolved.", dep)
		}(dep)
	}

	if !logger.OutputJSON {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			var unresolved []*Conflict
			for _, dep := range dependencies {
				if !conflicts[dep].resolved {
					unresolved = append(unresolved, conflicts[dep])
				}
			}

			if len(unresolved) == 0 {
				return
			}
			j, _ := json.Marshal(unresolved)
			fmt.Println(string(j))
			time.Sleep(resolverSleepInterval * time.Second)
		}
	}()
}

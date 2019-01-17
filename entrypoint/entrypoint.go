package entrypoint

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	"k8s.io/client-go/rest"
)

var dependencies []Resolver // List containing all dependencies to be resolved

const (
	//DependencyPrefix is a prefix for env variables
	DependencyPrefix            = "DEPENDENCY_"
	JsonSuffix                  = "_JSON"
	PrintSleepInterval          = 2.0     // The number of seconds to wait between each "JSON status report"
	DefaultInitialSleepInterval = 2000.0  // The default number of milliseconds to sleep before exponential backoff
	DefaultMaxSleepInterval     = 20000.0 // The default number of milliseconds to max out wait time
	DefaultBackoffRate          = 1.25    // The default rate at which to increase the exponential backoff. In other words, the base
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
	backoffRate := env.GetEnvFloat("BACKOFF_RATE", DefaultBackoffRate)
	maxSleepInterval := env.GetEnvFloat("MAX_SLEEP_INTERVAL", DefaultMaxSleepInterval)
	for _, dep := range dependencies {
		wg.Add(1)
		conflicts[dep] = NewConflict(dep.GetDependency())
		go func(dep Resolver) {
			defer wg.Done()
			logger.Info("Resolving %v", dep)
			sleepInterval := env.GetEnvFloat("INITIAL_SLEEP_INTERVAL", DefaultInitialSleepInterval)
			var err error
			status := false
			for status == false {
				if status, err = dep.IsResolved(e); err != nil {
					logger.Warning("Resolving dependency %s failed: %v - Trying again in %0.3f seconds", dep, err, sleepInterval/1000.0)
					conflicts[dep].Update(false, err.Error())
					time.Sleep(time.Duration(sleepInterval) * time.Millisecond)
					sleepInterval = env.Backoff(sleepInterval, backoffRate, maxSleepInterval)
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
			time.Sleep(PrintSleepInterval * time.Second)
		}
	}()
}

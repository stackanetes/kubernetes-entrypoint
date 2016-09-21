package daemonset

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type Daemonset struct {
	name string
}

func init() {
	daemonsetEnv := fmt.Sprintf("%sDAEMONSET", entry.DependencyPrefix)
	if daemonsetsDeps := env.SplitEnvToList(daemonsetEnv); daemonsetsDeps != nil {
		for _, dep := range daemonsetsDeps {
			entry.Register(NewDaemonSet(dep))
		}
	}
}

func NewDaemonSet(name string) Daemonset {
	return Daemonset{name: name}
}

func (d Daemonset) IsResolved(entrypoint *entry.Entrypoint) (bool, error) {
	daemonset, err := entrypoint.Client.ExtensionsClient.DaemonSets(entrypoint.Namespace).Get(d.name)
	if err != nil {
		return false, err
	}
	label := labels.SelectorFromSet(daemonset.Spec.Selector.MatchLabels)
	opts := api.ListOptions{LabelSelector: label}
	pods, err := entrypoint.Client.Pods(entrypoint.Namespace).List(opts)
	if err != nil {
		return false, err
	}
	myPodName := os.Getenv("POD_NAME")
	if myPodName == "" {
		logger.Error.Print("Environment variable POD_NAME not set")
		os.Exit(1)

	}
	myPod, err := entrypoint.Client.Pods(entrypoint.Namespace).Get(myPodName)
	if err != nil {
		logger.Error.Printf("Getting POD: %v failed : %v", myPodName, err)
		os.Exit(1)
	}
	myHost := myPod.Status.HostIP

	for _, pod := range pods.Items {
		if !isPodOnHost(&pod, myHost) {
			continue
		}
		if api.IsPodReady(&pod) {
			return true, nil
		}
		return false, fmt.Errorf("Pod %v of daemonset %v is not ready", pod.Name, d.GetName())

	}
	return true, nil
}

func (d Daemonset) GetName() string {
	return d.name
}

func isPodOnHost(pod *api.Pod, hostIP string) bool {
	if pod.Status.HostIP == hostIP {
		return true
	}
	return false
}

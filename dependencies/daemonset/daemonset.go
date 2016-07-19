package daemonset

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/docker-entrypoint/dependencies"
	"github.com/stackanetes/docker-entrypoint/logger"
	"github.com/stackanetes/docker-entrypoint/util/env"
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
			entry.Register(NewDaemonset(dep))
		}
	}
}

func NewDaemonset(name string) Daemonset {
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
		if !podReady(&pod) {
			return false, fmt.Errorf("Pod %v of daemonset %v is not ready", pod.Name, d.GetName())
		}

	}
	if !isPodOnHost(pods.Items, myHost) {
		return false, fmt.Errorf("Hostname mismatch: Daemonset %v is not on the same host as Pod %v", d.GetName(), myPodName)
	}
	return true, nil
}

func (d Daemonset) GetName() string {
	return d.name
}

func podReady(pod *api.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == api.PodReady && cond.Status == api.ConditionTrue {
			return true
		}
	}
	return false
}

func isPodOnHost(podList []api.Pod, hostIP string) bool {
	for _, pod := range podList {
		if pod.Status.HostIP == hostIP {
			return true
		}
	}
	return false
}

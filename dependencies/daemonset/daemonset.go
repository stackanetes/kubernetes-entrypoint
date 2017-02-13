package daemonset

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/labels"
)

const (
	PodNameEnvVar            = "POD_NAME"
	PodNameNotSetErrorFormat = "Env POD_NAME not set. Daemonset dependency %s will be ignored!"
)

type Daemonset struct {
	name    string
	podName string
}

func init() {
	daemonsetEnv := fmt.Sprintf("%sDAEMONSET", entry.DependencyPrefix)
	if daemonsetsDeps := env.SplitEnvToList(daemonsetEnv); daemonsetsDeps != nil {
		for _, dep := range daemonsetsDeps {
			daemonset, err := NewDaemonset(dep)
			if err != nil {
				logger.Error.Printf("Cannot initialize daemonset: %v", err)
				continue
			}
			entry.Register(daemonset)
		}
	}
}

func NewDaemonset(name string) (*Daemonset, error) {
	if os.Getenv(PodNameEnvVar) == "" {
		return nil, fmt.Errorf(PodNameNotSetErrorFormat, name)
	}
	return &Daemonset{
		name:    name,
		podName: os.Getenv(PodNameEnvVar),
	}, nil
}

func (d Daemonset) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	var myPodName string
	daemonset, err := entrypoint.Client().DaemonSets(entrypoint.GetNamespace()).Get(d.GetName())

	if err != nil {
		return false, err
	}

	label := labels.SelectorFromSet(daemonset.Spec.Selector.MatchLabels)
	opts := v1.ListOptions{LabelSelector: label.String()}
	pods, err := entrypoint.Client().Pods(entrypoint.GetNamespace()).List(opts)
	if err != nil {
		return false, err
	}

	myPod, err := entrypoint.Client().Pods(entrypoint.GetNamespace()).Get(d.podName)
	if err != nil {
		return false, fmt.Errorf("Getting POD: %v failed : %v", myPodName, err)
	}

	myHost := myPod.Status.HostIP

	for _, pod := range pods.Items {
		if !isPodOnHost(&pod, myHost) {
			continue
		}
		if isPodReady(pod) {
			return true, nil
		}
		return false, fmt.Errorf("Pod %v of daemonset %v is not ready", pod.Name, d.GetName())

	}
	return true, nil
}

func (d Daemonset) GetName() string {
	return d.name
}

func isPodOnHost(pod *v1.Pod, hostIP string) bool {
	if pod.Status.HostIP == hostIP {
		return true
	}
	return false
}

func isPodReady(pod v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady && condition.Status == "True" {
			return true
		}
	}
	return false
}

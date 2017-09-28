package daemonset

import (
	"fmt"
	"os"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
	api "k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/labels"
)

const (
	PodNameEnvVar            = "POD_NAME"
	PodNameNotSetErrorFormat = "Env POD_NAME not set. Daemonset dependency %s in namespace %s will be ignored!"
)

type Daemonset struct {
	name      string
	namespace string
	podName   string
}

func init() {
	daemonsetEnv := fmt.Sprintf("%sDAEMONSET", entry.DependencyPrefix)
	if daemonsetsDeps := env.SplitEnvToDeps(daemonsetEnv); daemonsetsDeps != nil {
		for _, dep := range daemonsetsDeps {
			daemonset, err := NewDaemonset(dep.Name, dep.Namespace)
			if err != nil {
				logger.Error.Printf("Cannot initialize daemonset: %v", err)
				continue
			}
			entry.Register(daemonset)
		}
	}
}

func NewDaemonset(name string, namespace string) (*Daemonset, error) {
	if os.Getenv(PodNameEnvVar) == "" {
		return nil, fmt.Errorf(PodNameNotSetErrorFormat, name, namespace)
	}
	return &Daemonset{
		name:      name,
		namespace: namespace,
		podName:   os.Getenv(PodNameEnvVar),
	}, nil
}

func (d Daemonset) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	var myPodName string
	daemonset, err := entrypoint.Client().DaemonSets(d.namespace).Get(d.name)

	if err != nil {
		return false, err
	}

	label := labels.SelectorFromSet(daemonset.Spec.Selector.MatchLabels)
	opts := api.ListOptions{LabelSelector: label}

	daemonsetPods, err := entrypoint.Client().Pods(d.namespace).List(opts)
	if err != nil {
		return false, err
	}

	myPod, err := entrypoint.Client().Pods(env.GetBaseNamespace()).Get(d.podName)
	if err != nil {
		return false, fmt.Errorf("Getting POD: %v failed : %v", myPodName, err)
	}

	myHost := myPod.Status.HostIP

	for _, pod := range daemonsetPods.Items {
		if !isPodOnHost(&pod, myHost) {
			continue
		}
		if isPodReady(pod) {
			return true, nil
		}
		return false, fmt.Errorf("Pod %v of daemonset %s is not ready", pod.Name, d)

	}
	return true, nil
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

func (d Daemonset) String() string {
	return fmt.Sprintf("Daemonset %s in namespace %s", d.name, d.namespace)
}

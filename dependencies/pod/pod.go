package pod

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
	PodNameNotSetErrorFormat = "Env POD_NAME not set. Pod dependency in namespace %s will be ignored!"
)

type Pod struct {
	namespace       string
	labels          map[string]string
	requireSameNode bool
	podName         string
}

func init() {
	podEnv := fmt.Sprintf("%sPOD", entry.DependencyPrefix)
	if podDeps := env.SplitPodEnvToDeps(podEnv); podDeps != nil {
		for _, dep := range podDeps {
			pod, err := NewPod(dep.Labels, dep.Namespace, dep.RequireSameNode)
			if err != nil {
				logger.Error.Printf("Cannot initialize pod: %v", err)
				continue
			}
			entry.Register(pod)
		}
	}
}

func NewPod(labels map[string]string, namespace string, requireSameNode bool) (*Pod, error) {
	if os.Getenv(PodNameEnvVar) == "" {
		return nil, fmt.Errorf(PodNameNotSetErrorFormat, namespace)
	}
	return &Pod{
		labels:          labels,
		namespace:       namespace,
		requireSameNode: requireSameNode,
		podName:         os.Getenv(PodNameEnvVar),
	}, nil
}

func (p Pod) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	myPod, err := entrypoint.Client().Pods(env.GetBaseNamespace()).Get(p.podName)
	if err != nil {
		return false, fmt.Errorf("Getting POD: %v failed : %v", p.podName, err)
	}
	myHost := myPod.Status.HostIP

	label := labels.SelectorFromSet(p.labels)
	opts := api.ListOptions{LabelSelector: label}

	matchingPodList, err := entrypoint.Client().Pods(p.namespace).List(opts)
	if err != nil {
		return false, err
	}

	matchingPods := matchingPodList.Items
	if len(matchingPods) == 0 {
		return false, fmt.Errorf("Found no pods matching labels: %v", p.labels)
	}

	podCount := 0
	for _, pod := range matchingPods {
		podCount++
		if p.requireSameNode && !isPodOnHost(&pod, myHost) {
			continue
		}
		if isPodReady(pod) {
			return true, nil
		}
	}
	onHostClause := ""
	if p.requireSameNode {
		onHostClause = " on host"
	}
	if podCount == 0 {
		return false, fmt.Errorf("Found no pods%v matching labels: %v", onHostClause, p.labels)
	} else {
		return false, fmt.Errorf("Found %v pods%v, but none ready, matching labels: %v", podCount, onHostClause, p.labels)
	}
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

func (p Pod) String() string {
	return fmt.Sprintf("Pod on same host with labels %v in namespace %s", p.labels, p.namespace)
}

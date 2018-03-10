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
	namespace string
	labels    map[string]string
	podName   string
}

func init() {
	podEnv := fmt.Sprintf("%sPOD", entry.DependencyPrefix)
	if podDeps := env.SplitPodEnvToDeps(podEnv); podDeps != nil {
		for _, dep := range podDeps {
			pod, err := NewPod(dep.Labels, dep.Namespace)
			if err != nil {
				logger.Error.Printf("Cannot initialize pod: %v", err)
				continue
			}
			entry.Register(pod)
		}
	}
}

func NewPod(labels map[string]string, namespace string) (*Pod, error) {
	if os.Getenv(PodNameEnvVar) == "" {
		return nil, fmt.Errorf(PodNameNotSetErrorFormat, namespace)
	}
	return &Pod{
		namespace: namespace,
		labels:    labels,
		podName:   os.Getenv(PodNameEnvVar),
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
		return false, fmt.Errorf("No pods found matching labels: %v", p.labels)
	}

	hostPodCount := 0
	for _, pod := range matchingPods {
		if !isPodOnHost(&pod, myHost) {
			continue
		}
		hostPodCount++
		if isPodReady(pod) {
			return true, nil
		}
	}
	if hostPodCount == 0 {
		return false, fmt.Errorf("Found no pods on host matching labels: %v", p.labels)
	} else {
		return false, fmt.Errorf("Found %v pods on host, but none ready, matching labels: %v", hostPodCount, p.labels)
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

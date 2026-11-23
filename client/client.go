// Copyright 2018 The kubernetes-entrypoint Authors
// Copyright 2026 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	discoveryv1client "k8s.io/client-go/kubernetes/typed/discovery/v1"
	"k8s.io/client-go/rest"

	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

type ClientInterface interface {
	Pods(string) v1core.PodInterface
	Jobs(string) v1batch.JobInterface
	CronJobs(string) v1batch.CronJobInterface
	EndpointSlices(string) discoveryv1client.EndpointSliceInterface
	DaemonSets(string) appsv1.DaemonSetInterface
	Services(string) v1core.ServiceInterface
	CustomResource(ctx context.Context, apiVersion, namespace, resource, name string) (*unstructured.Unstructured, error)
}
type Client struct {
	client        kubernetes.Interface
	dynamicClient dynamic.Interface
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.client.CoreV1().Pods(namespace)
}

func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.client.BatchV1().Jobs(namespace)
}

func (c Client) CronJobs(namespace string) v1batch.CronJobInterface {
	return c.client.BatchV1().CronJobs(namespace)
}

func (c Client) EndpointSlices(namespace string) discoveryv1client.EndpointSliceInterface {
	return c.client.DiscoveryV1().EndpointSlices(namespace)
}

func (c Client) DaemonSets(namespace string) appsv1.DaemonSetInterface {
	return c.client.AppsV1().DaemonSets(namespace)
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.client.CoreV1().Services(namespace)
}

func (c Client) CustomResource(
	ctx context.Context,
	apiVersion, kind, namespace, name string,
) (*unstructured.Unstructured, error) {

	apiResourceList, err := c.client.Discovery().ServerResourcesForGroupVersion(apiVersion)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(apiVersion, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf(`apiVersion [%s] must be "group/version"`,
			apiVersion)
	}
	group, version := parts[0], parts[1]

	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Kind == kind {
			gvr := schema.GroupVersionResource{
				Group:    group,
				Version:  version,
				Resource: apiResource.Name,
			}

			resourceClient := c.dynamicClient.Resource(gvr)
			if apiResource.Namespaced {
				if namespace == "" {
					namespace = env.GetBaseNamespace()
				}
				return resourceClient.Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
			}

			return resourceClient.Get(ctx, name, metav1.GetOptions{})
		}
	}
	return nil, fmt.Errorf("could not find resource with version %v, "+
		"kind %v, and name %v in namespace %v",
		apiVersion, kind, name, namespace)
}

func New(config *rest.Config) (ClientInterface, error) {
	if config == nil {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return Client{clientset, dynamicClient}, nil
}

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

package mocks

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	v1apps "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	discoveryv1client "k8s.io/client-go/kubernetes/typed/discovery/v1"
)

type Client struct {
	v1core.PodInterface
	v1core.ServiceInterface
	v1apps.DaemonSetInterface
	v1apps.DeploymentInterface
	discoveryv1client.EndpointSliceInterface
	v1batch.JobInterface

	FakeCustomResource *unstructured.Unstructured
	Err                error
}

func (c Client) Pods(namespace string) v1core.PodInterface {
	return c.PodInterface
}

func (c Client) Services(namespace string) v1core.ServiceInterface {
	return c.ServiceInterface
}

func (c Client) DaemonSets(namespace string) v1apps.DaemonSetInterface {
	return c.DaemonSetInterface
}

func (c Client) Deployments(namespace string) v1apps.DeploymentInterface {
	return c.DeploymentInterface
}

func (c Client) EndpointSlices(namespace string) discoveryv1client.EndpointSliceInterface {
	return c.EndpointSliceInterface
}

func (c Client) Jobs(namespace string) v1batch.JobInterface {
	return c.JobInterface
}

func (c Client) CustomResource(
	ctx context.Context,
	apiVersion, namespace, resource, name string,
) (*unstructured.Unstructured, error) {

	return c.FakeCustomResource, c.Err
}

func NewClient() *Client {
	return &Client{
		PodInterface:           NewPClient(),
		ServiceInterface:       NewSClient(),
		DaemonSetInterface:     NewDSClient(),
		EndpointSliceInterface: NewEClient(),
		JobInterface:           NewJClient(),
	}
}

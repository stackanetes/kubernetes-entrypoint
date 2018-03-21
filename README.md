# Kubernetes Entrypoint

[![Build Status](https://api.travis-ci.org/stackanetes/kubernetes-entrypoint.svg?branch=master "Build Status")](https://travis-ci.org/stackanetes/kubernetes-entrypoint)
[![Container Repository on Quay](https://quay.io/repository/stackanetes/kubernetes-entrypoint/status "Container Repository on Quay")](https://quay.io/repository/stackanetes/kubernetes-entrypoint)
[![Go Report Card](https://goreportcard.com/badge/stackanetes/kubernetes-entrypoint "Go Report Card")](https://goreportcard.com/report/stackanetes/kubernetes-entrypoint)


============

Kubernetes-entrypoint enables complex deployments on top of Kubernetes.

## Overview

Kubernetes-entrypoint is meant to be used as a container entrypoint, which means it has to bundled in the container.
Before launching the desired application, the entrypoint verifies and waits for all specified dependencies to be met.

The Kubernetes-entrypoint queries directly the Kubernetes API and each container is self-aware of its dependencies and their states.
Therefore, no centralized orchestration layer is required to manage deployments and scenarios such as failure recovery or pod migration become easy.

## Usage

Kubernetes-entrypoint reads the dependencies out of environment variables passed into a container.
There is only one required environment variable "COMMAND" which specifies a command (arguments delimited by whitespace) which has to be executed when all dependencies are resolved:

`COMMAND="sleep inf"`

Kubernetes-entrypoint introduces a wide variety of dependencies which can be used to better orchestrate once deployment.

## Latest features

Extending functionality of kubernetes-entrypoint by adding an ability to specify dependencies in different namespaces. The new format for writing dependencies is `namespace:name`, with the exception of pod dependencies which us json. To ensure backward compatibility if the `namespace:` is omitted, it behaves just like in previous versions so it assumes that dependecies are running at the same namespace as kubernetes-entrypoint. This feature is not implemented for container, config and socket dependency because in such cases the different namespace is irrelevant.

For instance:
`
DEPENDENCY_SERVICE=mysql:mariadb,keystone-api
`

The new entrypoint will resolve mariadb in mysql namespace and keystone-api in the same namespace as entrypoint was deployed in.

## Supported types of dependencies

All dependencies are passed as environement variables in format of `DEPENDENCY_<NAME>` delimited by colon. For dependencies to be effective please use [readiness probes](http://kubernetes.io/docs/user-guide/production-pods/#liveness-and-readiness-probes-aka-health-checks) for all containers.

### Service
Checks whether given kubernetes service has at least one endpoint.
Example:

`DEPENDENCY_SERVICE=mariadb,keystone-api`

### Container
Within a pod composed of multiple containers, it waits for the containers specified by their names to start.
This dependency requires a `POD_NAME` environement variable which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
Example:

`DEPENDENCY_CONTAINER=nova-libvirt,virtlogd`

### Daemonset
Checks if a specified daemonset is already running on the same host, this dependency requires a `POD_NAME`
env which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
The `POD_NAME` variable is mandatory and is used to resolve dependencies.
Example:

`DEPENDENCY_DAEMONSET=openvswitch-agent`

Simple example how to use downward API to get `POD_NAME` can be found [here](https://raw.githubusercontent.com/kubernetes/kubernetes.github.io/master/docs/user-guide/downward-api/dapi-pod.yaml).

### Job
Checks if a given job succeded at least once.
Example:

`DEPENDENCY_JOBS=nova-init,neutron-init`

### Config
This dependency performs a container level templating of configuration files. It can template an ip address `{{ .IP }}` and hostname `{{ .HOSTNAME }}`.
Templated config has to be stored in an arbitrary directory `/configmaps/<name_of_file>/<name_of_file>`.
This dependency requires `INTERFACE_NAME` environment variable to know which interface to use for obtain ip address.
Example:

`DEPENDENCY_CONFIG=/etc/nova/nova.conf`

The Kubernetes-entrypoint will look for the configuration file `/configmaps/nova.conf/nova.conf`, template
`{{ .IP }} and {{ .HOSTNAME }}` tags and save the file as `/etc/nova/nova.conf`.

### Socket
Checks whether a given file exists and container has rights to read it.
Example:

`DEPENDENCY_SOCKET=/var/run/openvswitch/ovs.socket`

### Pod
Checks if at least one pod matching the specified labels is already running, by
default anywhere in the cluster, or use `"requireSameNode": true` to require a
a pod on the same node.
In contrast to other dependencies, the syntax uses json in order to avoid inventing a new
format to specify labels and the parsing complexities that would come with that.
This dependency requires a `POD_NAME` env which can be easily passed through the
[downward api](http://kubernetes.io/docs/user-guide/downward-api/). The `POD_NAME` variable is mandatory and is used to resolve dependencies.
Example:

`DEPENDENCY_POD="[{\"namespace\": \"foo\", \"labels\": {\"k1\": \"v1\", \"k2\": \"v2\"}}, {\"labels\": {\"k1\": \"v1\", \"k2\": \"v2\"}, \"requireSameNode\": true}]"`

## Image

Build process for image is trigged after each commit.
Can be found [here](https://quay.io/repository/stackanetes/kubernetes-entrypoint?tab=tags), and pulled by executing:
`docker pull quay.io/stackanetes/kubernetes-entrypoint:v0.1.0`

## Examples

[Stackanetes](http://github.com/stackanetes/stackanetes) uses kubernetes-entrypoint to manage dependencies when deploying OpenStack on Kubernetes.

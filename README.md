# Kubernetes Entrypoint

[![Build Status](https://api.travis-ci.org/stackanetes/kubernetes-entrypoint.svg?branch=master "Build Status")](https://travis-ci.org/stackanetes/kubernetes-entrypoint)
[![Container Repository on Quay](https://quay.io/repository/stackanetes/kubernetes-entrypoint/status "Container Repository on Quay")](https://quay.io/repository/stackanetes/kubernetes-entrypoint)
[![Go Report Card](https://goreportcard.com/badge/stackanetes/kubernetes-entrypoint "Go Report Card")](https://goreportcard.com/report/stackanetes/kubernetes-entrypoint)


============

Kubernetes-entrypoint enables complex deployments on top of Kubernetes.

## Overview

Kubernetes-entrypoint needs to be deployed as a container entrypoint which means it has to bundled in container. Before launching desired application
it checks if all specified dependencies are met.

Kubernetes-entrypoint talks directly to Kubernetes API so there is no need for centralized orchestration layer to manage deployments.
Each container is self aware of its dependencies so thing like failure recovery or pods migration becomes really easy.

## Usage

Kubernetes-entrypoint uses environment variables passed into a container.
There is only one required environment variable "COMMAND" which specifies a command(arguments delimited by whitespace) which has to be executed when all dependencies are resolved:

```COMMAND="sleep inf"```

Kubernetes-entrypoint introduces a range variety of dependencies which can be used to better orchestrate once deployment.

## Dependencies

All dependencies are passed as environement variables in format of `DEPENDENCY_<NAME>` delimited by colon. For dependencies to be effective
please use [readiness probes](http://kubernetes.io/docs/user-guide/production-pods/#liveness-and-readiness-probes-aka-health-checks) for all containers.

- Service:
    checks whether given kubernetes service has at least one endpoint.  
    Example:

    ```DEPENDENCY_SERVICE=mariadb,keystone-api```

- Container:
    for multicontainer pod we can specify to wait with launching application on other container in the same pod,
    this dependency requires a `POD_NAME` env which can be easily passed as [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
    Example:

    ```DEPENDENCY_CONTAINER=nova-libvirt,virtlogd```

- Daemonset:
    checks if a specified daemonset is already running on the same host, this dependency requires a `POD_NAME`
    env which can be easily passed as [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
    Example:

    ```DEPENDENCY_DAEMONSET=openvswitch-agent```

- Job:
    checks if a given job succeded at least once.
    Example:

    ```DEPENDENCY_JOBS=nova-init,neutron-init```

- Config:
    this dependency performs a container level templating of configuration files. It can template an ip address `{{ .IP }}` and hostname `{{ .HOSTNAME }}`. 
    Templated config has to be stored in an arbitrary directory `/configmaps/<name_of_file>/<name_of_file>`.
    This dependency requires `INTERFACE_NAME` environment variable to know which interface to use for obtain ip address. 
    Example:

    ```DEPENDENCY_CONFIG=/etc/nova/nova.conf```

    kubernetes-entrypoint will look for config in /configmaps/nova.conf/nova.conf template `{{ .IP }} and {{ .HOSTNAME }}` and place it as
    ```/etc/nova/nova.conf```

- Socket:
    checks wheter a given file exists and container has rights to read it.
    Example:

    ```DEPENDENCY_SOCKET=/var/run/openvswitch/ovs.socket```

## Examples

[Stackanetes](http://github/stackanetes/stackanetes) uses kubernetes-entrypoint to manage dependencies when deploying OpenStack on Kubernetes.


# Kubernetes Entrypoint

============

Kubernetes-entrypoint enables complex deployments on top of Kubernetes.

## Overview

Kubernetes-entrypoint is meant to be used as a container entrypoint, which means it has to bundled in the container.
Before launching the desired application, the entrypoint verifies and waits for all specified dependencies to be met.

Kubernetes-entrypoint queries the Kubernetes API directly, and each container is self-aware of its dependencies and their states.
Therefore, no centralized orchestration layer is required to manage deployments, and scenarios (such as failure recovery or pod migration) become easy.

## Usage

Kubernetes-entrypoint reads the dependencies out of environment variables passed into a container.
There is only one required environment variable `COMMAND` which specifies a command (arguments delimited by whitespace) which has to be executed when all dependencies are resolved:

`COMMAND="sleep inf"`

## Latest features

Extending functionality of kubernetes-entrypoint by adding an ability to specify dependencies in different namespaces. The new format for writing dependencies is `namespace:name`, with the exception of pod dependencies (which use JSON).
In order to ensure backward compatibility, if `namespace` is omitted, it is assumed that dependencies are running in the same namespace as kubernetes-entrypoint, just like in previous versions.
This feature is not implemented for Container, Config or Socket dependency because the different namespace is irrelevant for those cases.

For instance:
`
DEPENDENCY_SERVICE=mysql:mariadb,keystone-api
`

The new entrypoint will resolve mariadb in the mysql namespace and keystone-api in the same namespace as kubernetes-entrypoint was deployed in.

## Supported types of dependencies

All dependencies are passed as environment variables with the format `DEPENDENCY_<NAME>`, delimited by a colon.
For dependencies to be effective please use [readiness probes](http://kubernetes.io/docs/user-guide/production-pods/#liveness-and-readiness-probes-aka-health-checks) for all containers.

### Service
Checks whether given kubernetes service has at least one endpoint.
Example:

`DEPENDENCY_SERVICE=mariadb,keystone-api`

### Container
Within a pod composed of multiple containers, kubernetes-entrypoint waits for the containers specified by their names to start.
This dependency requires a `POD_NAME` environment variable which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
Example:

`DEPENDENCY_CONTAINER=nova-libvirt,virtlogd`

### Daemonset
Checks if a specified daemonset is already running on the same host
This dependency requires a `POD_NAME` environment variable which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
The `POD_NAME` variable is mandatory and is used to resolve dependencies.
Example:

`DEPENDENCY_DAEMONSET=openvswitch-agent`

A simple example of how to use downward API to get `POD_NAME` can be found [here](https://raw.githubusercontent.com/kubernetes/kubernetes.github.io/master/docs/user-guide/downward-api/dapi-pod.yaml).

### CronJob

Waits for a named `batch/v1` CronJob to have the expected **top-level metadata
annotations and/or labels**, then waits until no unfinished Jobs from older
versions remain:

```sh
DEPENDENCY_CRONJOB_JSON='[{"name":"cleanup","namespace":"application","annotations":{"example.com/release-id":"my-release:42"}}]'
# Labels can be used instead of, or together with, annotations:
DEPENDENCY_CRONJOB_JSON='[{"name":"cleanup","labels":{"app.kubernetes.io/version":"v42"}}]'
```

`name` is required; `namespace` defaults to `NAMESPACE`, or `default` if unset.
Every expected annotation and label must exist and match exactly, including
empty-string values. Extra actual metadata is allowed. Labels are expectations
on the named CronJob, not selectors. A different or newer marker is not accepted.
Job-template metadata alone does not substitute for top-level CronJob metadata.

To identify the version of generated Jobs, copy the same expected markers into
`spec.jobTemplate.metadata` (**not** `spec.jobTemplate.spec.template.metadata`).
Jobs missing any expected marker, or having different values, are treated as old.
Jobs matching every expected marker may keep running without blocking the dependency.
Only Jobs controlled by this named CronJob in its namespace are considered;
unfinished Jobs from a previous incarnation of that CronJob also block. Unrelated
and manually created Jobs without a CronJob controller reference are ignored.

An old Job is finished when it has a `Complete=True` or `Failed=True` condition
and no active or reported terminating Pods. Pending, retrying, suspended and
terminating old Jobs still block, even with zero active Pods. Successful and
failed history entries do not block. The CronJob need not schedule a new Job:
an already-updated CronJob with no unfinished old Jobs resolves immediately.

Without version markers, the plain form (or JSON with empty metadata maps)
waits for the CronJob to exist and **all** its unfinished Jobs to drain; it cannot
distinguish an old release that has not yet been updated:

```sh
DEPENDENCY_CRONJOB=cleanup,other-namespace:backup
```

When both variables are nonempty, JSON takes precedence with a warning. Invalid
configuration is logged and may **skip dependencies rather than terminate the
process**, following existing registration behavior. Malformed JSON or invalid
field types skip the JSON configuration; entries missing a name are skipped.

#### Helm release-aware ordering

Render a shared release marker on the CronJob and its Job template:

```yaml
metadata:
  annotations:
    example.com/release-id: "{{ .Release.Name }}:{{ .Release.Revision }}"
spec:
  jobTemplate:
    metadata:
      annotations:
        example.com/release-id: "{{ .Release.Name }}:{{ .Release.Revision }}"
```

Use the same marker in the waiting container's environment:

```yaml
env:
  - name: DEPENDENCY_CRONJOB_JSON
    value: >-
      [{"name":"cleanup","annotations":{"example.com/release-id":"{{ .Release.Name }}:{{ .Release.Revision }}"}}]
```

Unlike Deployments, CronJobs have no `status.observedGeneration` acknowledgement.
This dependency checks the submitted metadata and the actual Job list, not
`status.active`. It rechecks the CronJob before accepting success, but polling
does not lock scheduling or prevent subsequent updates. Use unique version
markers and update the CronJob and Job-template markers together with its spec.
The waiting container needs RBAC permissions to `get` `batch/cronjobs` and
`list` `batch/jobs` in each dependency namespace.

### Job
Checks if a given job or set of jobs with matching name and/or labels succeeded at least once.
In order to use labels, `DEPENDENCY_JOBS_JSON` must be used.
DEPENDENCY_JOBS is supported as well for backward compatibility.
Examples:

`DEPENDENCY_JOBS_JSON='[{"namespace": "foo", "name": "nova-init"}, {"labels": {"initializes": "neutron"}}]'`
`DEPENDENCY_JOBS=nova-init,neutron-init'`

### Config
This dependency performs a container level templating of configuration files. It can template an ip address `{{ .IP }}` and hostname `{{ .HOSTNAME }}`.
Templated config has to be stored in an arbitrary directory `/configmaps/<name_of_file>/<name_of_file>`.
This dependency requires an `INTERFACE_NAME` environment variable to know which interface to use to obtain the ip address.
Example:

`DEPENDENCY_CONFIG=/etc/nova/nova.conf`

Kubernetes-entrypoint will look for the configuration file `/configmaps/nova.conf/nova.conf`, template the `{{ .IP }} and {{ .HOSTNAME }}` tags, and then save the file as `/etc/nova/nova.conf`.

### Socket
Checks whether a given file exists and that the container has rights to read it.
Example:

`DEPENDENCY_SOCKET=/var/run/openvswitch/ovs.socket`

### Pod
Checks if at least one pod matching the specified labels is already running, by default anywhere in the cluster, or use `"requireSameNode": true` to require a pod on the same node.
Labels are specified using JSON, as seen in the example below.
This dependency requires a `POD_NAME` env which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
The `POD_NAME` variable is mandatory and is used to resolve dependencies.
Example:

`DEPENDENCY_POD_JSON='[{"namespace": "foo", "labels": {"k1": "v1", "k2": "v2"}}, {"labels": {"k1": "v1", "k2": "v2"}, "requireSameNode": true}]'`

### Custom Resource
This dependency checks whether an arbitrary key on a given CustomResource
matches a desired value. The environment variable `DEPENDENCY_CUSTOM_RESOURCE`
dictates the specific object to watch, as well as the key and its desired
value.

The `key` field supports three formats:
1. **Dot notation** (e.g., `spec.arbitrary-key`) - for simple nested field access
2. **JSONPath with $ prefix** (e.g., `$.spec.arbitrary-key` or `$.items[0].status`) - for more complex queries including array indexing
3. **JSONPath with curly braces** (e.g., `{.spec.arbitrary-key}` or `{$.items[0].status}`) - alternative JSONPath syntax

**Important JSONPath restrictions:**
- Only single-value results are allowed (queries returning multiple values will fail)

**Value Comparison:**
- Values are compared using JSON encoding for consistency
- String values should be provided without quotes in the configuration (e.g., `"value": "ready"` not `"value": "\"ready\""`)
- Numeric values are compared as JSON-encoded numbers (e.g., `"value": "3.14"` for a float, `"value": "3"` for an integer)
- Boolean values should be provided as strings (e.g., `"value": "true"` or `"value": "false"`)
- `null` values are supported - use `"value": "null"` to match a null field

For example, suppose you have the following `DEPENDENCY_CUSTOM_RESOURCE` and CustomResource:
`DEPENDENCY_CUSTOM_RESOURCE='[{"apiVersion":"stable.example.com/v1","kind":"Foo","namespace":"default","name":"my-foo","fields":[{"key":"spec.arbitrary-key","value":"ready"}]}]'`
```
apiVersion: stable.example.com/v1
kind: Foo
metadata:
  name: my-foo
  namespace: default
spec:
  arbitrary-key: not-ready
```
Given the above, kubernetes-entrypoint will wait until the value of
`spec.arbitrary-key` has flipped from `not-ready` to `ready`.

**JSONPath Examples:**
```
# Simple field access
{"key": "$.status", "value": "ready"}

# Nested field access
{"key": "$.spec.template.status", "value": "ready"}

# Array indexing
{"key": "$.items[0].status", "value": "ready"}

# Numeric comparison
{"key": "$.spec.replicas", "value": "3"}

# Boolean comparison
{"key": "$.spec.enabled", "value": "true"}
```

Note also that `fields` is a list, meaning that multiple fields can be monitored.

## Image

Build process for image is triggered after each commit.
Can be found [here](https://quay.io/repository/airshipit/kubernetes-entrypoint?tab=tags), and pulled by executing:
`docker pull quay.io/airshipit/kubernetes-entrypoint:v0.1.0`

## Examples

[OpenStack Helm](https://github.com/openstack/openstack-helm-infra/blob/master/helm-toolkit/templates/snippets/_kubernetes_entrypoint_init_container.tpl) uses kubernetes-entrypoint to manage dependencies when deploying OpenStack on Kubernetes.

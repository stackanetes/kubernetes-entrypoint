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

package env

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"opendev.org/airship/kubernetes-entrypoint/logger"
)

const (
	name1            = "foo"
	name2            = "bar"
	defaultNamespace = "default"
	altNamespace1    = "fooNS"
	altNamespace2    = "barNS"
)

func TestSplitEnvToListWithComma(t *testing.T) {
	t.Setenv("TEST_LIST", name1+","+name2)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 2 {
		t.Errorf("Expected len to be 2 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0])
	}
	if list[1].Name != name2 {
		t.Errorf("Expected: %s got %s", name2, list[1])
	}
}

func TestSplitEnvToList(t *testing.T) {
	t.Setenv("TEST_LIST", name1)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 1 {
		t.Errorf("Expected len to be 1 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0])
	}
	if list[0].Namespace != defaultNamespace {
		t.Errorf("Expected: %s got %s", defaultNamespace, list[0].Namespace)
	}
}

func TestSplitEnvToListWithColon(t *testing.T) {
	t.Setenv("TEST_LIST", altNamespace1+":"+name1)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 1 {
		t.Errorf("Expected len to be 1 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0].Name)
	}
	if list[0].Namespace != altNamespace1 {
		t.Errorf("Expected: %s got %s", altNamespace1, list[0].Namespace)
	}
}

func TestSplitEnvToListWithTooManyColons(t *testing.T) {
	// TODO(howell): This should probably expect an error
	t.Setenv("TEST_LIST", "too:many:colons")
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 0 {
		t.Errorf("Expected list to be empty")
	}
}

func TestSplitEnvToListWithColonsAndCommas(t *testing.T) {
	t.Setenv("TEST_LIST", altNamespace1+":"+name1+","+altNamespace2+":"+name2)
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 2 {
		t.Errorf("Expected len to be 2 not %d", len(list))
	}
	if list[0].Name != name1 {
		t.Errorf("Expected: %s got %s", name1, list[0].Name)
	}
	if list[0].Namespace != altNamespace1 {
		t.Errorf("Expected: %s got %s", altNamespace1, list[0].Namespace)
	}
	if list[1].Name != name2 {
		t.Errorf("Expected: %s got %s", name2, list[0].Name)
	}
	if list[1].Namespace != altNamespace2 {
		t.Errorf("Expected: %s got %s", altNamespace2, list[0].Namespace)
	}
}

func TestSplitEnvToListWithMissingNamespace(t *testing.T) {
	t.Setenv("TEST_LIST", ":name")
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 0 {
		t.Errorf("Invalid format, missing namespace in pod")
	}
}

func TestSplitEmptyEnvWithColon(t *testing.T) {
	t.Setenv("TEST_LIST", "")
	list := SplitEnvToDeps("TEST_LIST")
	if len(list) != 0 {
		t.Errorf("Expected list to be empty")
	}
}

func TestSplitPodEnvToDepsSuccess(t *testing.T) {
	testListJSONVal := `[
  {
    "namespace": "` + name1 + `",
    "labels": {
      "k1": "v1",
      "k2": "v2"
    },
    "requireSameNode": true
  },
  {
    "labels": {
      "k1": "v1",
      "k2": "v2"
     }
  }
]`
	t.Setenv("NAMESPACE", "TEST_NAMESPACE")
	t.Setenv("TEST_LIST_JSON", testListJSONVal)
	actual := SplitPodEnvToDeps("TEST_LIST_JSON")
	expected := []PodDependency{
		{
			Namespace: name1,
			Labels: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
			RequireSameNode: true,
		},
		{
			Namespace: "TEST_NAMESPACE",
			Labels: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
			RequireSameNode: false,
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitPodEnvToDepsUnset(t *testing.T) {
	t.Setenv("TEST_LIST", "")
	actual := SplitPodEnvToDeps("TEST_LIST")
	if len(actual) != 0 {
		t.Errorf("Expected: no dependencies Got: %v", actual)
	}
}

func TestSplitPodEnvToDepsIgnoreInvalid(t *testing.T) {
	t.Setenv("TEST_LIST", `[{"invalid": json}`)
	actual := SplitPodEnvToDeps("TEST_LIST")
	if len(actual) != 0 {
		t.Errorf("Expected: ignore invalid dependencies Got: %v", actual)
	}
}

func TestSplitJobEnvToDepsJsonSuccess(t *testing.T) {
	testListJSONVal := `[
  {
    "namespace": "` + altNamespace1 + `",
    "labels": {
      "k1": "v1",
      "k2": "v2"
    }
  },
  {
    "name": "` + name1 + `"
  }
]`
	t.Setenv("NAMESPACE", "TEST_NAMESPACE")
	t.Setenv("TEST_LIST_JSON", testListJSONVal)
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	expected := []JobDependency{
		{
			Namespace: altNamespace1,
			Labels: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
		{
			Name:      name1,
			Namespace: "TEST_NAMESPACE",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitJobEnvToDepsPlainSuccess(t *testing.T) {
	t.Setenv("NAMESPACE", "TEST_NAMESPACE")
	t.Setenv("TEST_LIST", "plain")
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	expected := []JobDependency{
		{
			Name:      "plain",
			Namespace: "TEST_NAMESPACE",
		},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitJobEnvToDepsJsonPrecedence(t *testing.T) {
	t.Setenv("NAMESPACE", "TEST_NAMESPACE")
	t.Setenv("TEST_LIST_JSON", `[{"name": "json"}]`)
	t.Setenv("TEST_LIST", "plain")
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	expected := []JobDependency{
		{
			Name:      "json",
			Namespace: "TEST_NAMESPACE",
		},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitJobEnvToDepsUnset(t *testing.T) {
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	if len(actual) != 0 {
		t.Errorf("Expected: no dependencies Got: %v", actual)
	}
}

func TestSplitJobEnvToDepsIgnoreInvalid(t *testing.T) {
	t.Setenv("TEST_LIST_JSON", `[{"invalid": json}`)
	actual := SplitJobEnvToDeps("TEST_LIST", "TEST_LIST_JSON")
	if len(actual) != 0 {
		t.Errorf("Expected: ignore invalid dependencies Got: %v", actual)
	}
}

func TestSplitCommandUnset(t *testing.T) {
	defer os.Unsetenv("COMMAND")
	list := SplitCommand()
	if len(list) > 0 {
		t.Errorf("Expected len to be 0, got %d", len(list))
	}
}

func TestSplitCommandEmpty(t *testing.T) {
	t.Setenv("COMMAND", "")
	list := SplitCommand()
	if len(list) > 0 {
		t.Errorf("Expected len to be 0, got %v", len(list))
	}
}

func TestSplitCommand(t *testing.T) {
	t.Setenv("COMMAND", "echo test")
	list := SplitCommand()
	if len(list) != 2 {
		t.Errorf("Expected two elements, got %v", len(list))
	}
	if list[0] != "echo" {
		t.Errorf("Expected echo, got %s", list[0])
	}
	if list[1] != "test" {
		t.Errorf("Expected test, got %s", list[1])
	}
}

func TestGetBaseNamespaceEmpty(t *testing.T) {
	t.Setenv("NAMESPACE", "")
	getBaseNamespace := GetBaseNamespace()
	if getBaseNamespace != defaultNamespace {
		t.Errorf("Expected namespace to be %s, got %s", defaultNamespace, getBaseNamespace)
	}
}

func TestGetBaseNamespace(t *testing.T) {
	t.Setenv("NAMESPACE", "foo")
	getBaseNamespace := GetBaseNamespace()
	if getBaseNamespace != "foo" {
		t.Errorf("Expected namespace to be foo, got %v", getBaseNamespace)
	}
}

func TestSplitDeploymentEnvToDeps(t *testing.T) {
	tests := []struct {
		name           string
		plain          string
		json           string
		namespace      string
		unsetNamespace bool
		want           []DeploymentDependency
		warning        string
	}{
		{name: "unset", want: []DeploymentDependency{}},
		{name: "one plain", plain: "api", want: []DeploymentDependency{{Name: "api", Namespace: "default"}}},
		{name: "unset namespace", plain: "api", unsetNamespace: true, want: []DeploymentDependency{{Name: "api", Namespace: "default"}}},
		{name: "multiple plain", plain: "api,worker", namespace: "application", want: []DeploymentDependency{
			{Name: "api", Namespace: "application"}, {Name: "worker", Namespace: "application"},
		}},
		{name: "cross namespace", plain: "api,other:worker", namespace: "application", want: []DeploymentDependency{
			{Name: "api", Namespace: "application"}, {Name: "worker", Namespace: "other"},
		}},
		{name: "JSON annotations and namespace", json: `[{"name":"api","annotations":{"release":"42","empty":""}},{"name":"worker","namespace":"other","annotations":{}}]`, namespace: "application", want: []DeploymentDependency{
			{Name: "api", Namespace: "application", Annotations: map[string]string{"release": "42", "empty": ""}},
			{Name: "worker", Namespace: "other", Annotations: map[string]string{}},
		}},
		{name: "JSON missing annotations", json: `[{"name":"api"}]`, want: []DeploymentDependency{{Name: "api", Namespace: "default", Annotations: map[string]string{}}}},
		{name: "JSON empty namespace", json: `[{"name":"api","namespace":""}]`, namespace: "application", want: []DeploymentDependency{{Name: "api", Namespace: "application", Annotations: map[string]string{}}}},
		{name: "JSON null annotations", json: `[{"name":"api","annotations":null}]`, unsetNamespace: true, want: []DeploymentDependency{{Name: "api", Namespace: "default", Annotations: map[string]string{}}}},
		{name: "JSON precedence", plain: "old", json: `[{"name":"api"}]`, want: []DeploymentDependency{{Name: "api", Namespace: "default", Annotations: map[string]string{}}}, warning: "Ignoring TEST_DEPLOYMENT since TEST_DEPLOYMENT_JSON was specified"},
		{name: "empty JSON array takes precedence", plain: "old", json: `[]`, want: []DeploymentDependency{}, warning: "Ignoring TEST_DEPLOYMENT"},
		{name: "malformed JSON", json: `[{"name":`, want: []DeploymentDependency{}, warning: "Invalid TEST_DEPLOYMENT_JSON"},
		{name: "invalid JSON does not fall back to plain", plain: "old", json: `invalid`, want: []DeploymentDependency{}, warning: "Invalid TEST_DEPLOYMENT_JSON"},
		{name: "JSON object instead of array", json: `{"name":"api"}`, want: []DeploymentDependency{}, warning: "expected an array"},
		{name: "JSON null instead of array", json: `null`, want: []DeploymentDependency{}, warning: "not null"},
		{name: "invalid name type", json: `[{"name":42}]`, want: []DeploymentDependency{}, warning: "cannot unmarshal number"},
		{name: "invalid namespace type", json: `[{"name":"api","namespace":true}]`, want: []DeploymentDependency{}, warning: "cannot unmarshal bool"},
		{name: "null namespace", json: `[{"name":"api","namespace":null}]`, want: []DeploymentDependency{}, warning: "namespace must be a string, not null"},
		{name: "invalid annotation map type", json: `[{"name":"api","annotations":[]}]`, want: []DeploymentDependency{}, warning: "cannot unmarshal array"},
		{name: "invalid annotation value type", json: `[{"name":"api","annotations":{"release":42}}]`, want: []DeploymentDependency{}, warning: "cannot unmarshal number"},
		{name: "null annotation value", json: `[{"name":"api","annotations":{"release":null}}]`, want: []DeploymentDependency{}, warning: `annotation "release" must be a string, not null`},
		{name: "invalid annotation value discards entire JSON configuration", json: `[{"name":"valid"},{"name":"api","annotations":{"release":null}}]`, want: []DeploymentDependency{}, warning: `dependency 2: annotation "release" must be a string, not null`},
		{name: "missing name", json: `[{"namespace":"application"}]`, want: []DeploymentDependency{}, warning: "dependency 1: name is required"},
		{name: "label selection unsupported", json: `[{"labels":{"app":"api"}}]`, want: []DeploymentDependency{}, warning: "name is required"},
		{name: "empty name", json: `[{"name":""}]`, want: []DeploymentDependency{}, warning: "name is required"},
		{name: "null name", json: `[{"name":null}]`, want: []DeploymentDependency{}, warning: "name is required"},
		{name: "whitespace name", json: `[{"name":" "}]`, want: []DeploymentDependency{}, warning: "name is required"},
		{name: "skip missing name only", json: `[{"name":""},{"name":"api"}]`, want: []DeploymentDependency{{Name: "api", Namespace: "default", Annotations: map[string]string{}}}, warning: "dependency 1: name is required"},
		{name: "plain missing name", plain: "application:", want: []DeploymentDependency{}, warning: "name is required"},
		{name: "plain empty entry", plain: ",api", want: []DeploymentDependency{{Name: "api", Namespace: "default"}}, warning: "name is required"},
		{name: "plain invalid namespace", plain: ":api", want: []DeploymentDependency{}, warning: "missing namespace"},
		{name: "plain too many colons", plain: "a:b:c", want: []DeploymentDependency{}, warning: "expected namespace:name"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("TEST_DEPLOYMENT", test.plain)
			t.Setenv("TEST_DEPLOYMENT_JSON", test.json)
			t.Setenv("NAMESPACE", test.namespace)
			if test.unsetNamespace {
				if err := os.Unsetenv("NAMESPACE"); err != nil {
					t.Fatal(err)
				}
			}
			var warnings bytes.Buffer
			previous := logger.Warning.Writer()
			logger.Warning.SetOutput(&warnings)
			t.Cleanup(func() { logger.Warning.SetOutput(previous) })
			actual := SplitDeploymentEnvToDeps("TEST_DEPLOYMENT", "TEST_DEPLOYMENT_JSON")
			if !reflect.DeepEqual(actual, test.want) {
				t.Errorf("Expected: %#v Got: %#v", test.want, actual)
			}
			if test.warning == "" && warnings.Len() != 0 {
				t.Errorf("Unexpected warning: %s", warnings.String())
			}
			if !strings.Contains(warnings.String(), test.warning) {
				t.Errorf("Expected warning %q Got: %s", test.warning, warnings.String())
			}
		})
	}
}

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

func TestSplitCronJobEnvToDeps(t *testing.T) {
	emptyMetadata := func(namespace string) []CronJobDependency {
		return []CronJobDependency{{Name: "cleanup", Namespace: namespace, Annotations: map[string]string{}, Labels: map[string]string{}}}
	}
	tests := []struct {
		name           string
		plain          string
		json           string
		namespace      string
		unsetNamespace bool
		want           []CronJobDependency
		warning        string
	}{
		{name: "unset", want: []CronJobDependency{}},
		{name: "one plain", plain: "cleanup", want: []CronJobDependency{{Name: "cleanup", Namespace: "default"}}},
		{name: "unset namespace", plain: "cleanup", unsetNamespace: true, want: []CronJobDependency{{Name: "cleanup", Namespace: "default"}}},
		{name: "cross namespace", plain: "cleanup,other:backup", namespace: "application", want: []CronJobDependency{
			{Name: "cleanup", Namespace: "application"}, {Name: "backup", Namespace: "other"},
		}},
		{name: "JSON metadata and namespace", json: `[{"name":"cleanup","annotations":{"release":"my-release:42","empty":""},"labels":{"version":"v42"}},{"name":"backup","namespace":"other"}]`, namespace: "application", want: []CronJobDependency{
			{Name: "cleanup", Namespace: "application", Annotations: map[string]string{"release": "my-release:42", "empty": ""}, Labels: map[string]string{"version": "v42"}},
			{Name: "backup", Namespace: "other", Annotations: map[string]string{}, Labels: map[string]string{}},
		}},
		{name: "JSON no expectations", json: `[{"name":"cleanup"}]`, want: emptyMetadata("default")},
		{name: "JSON empty namespace", json: `[{"name":"cleanup","namespace":""}]`, namespace: "application", want: emptyMetadata("application")},
		{name: "JSON null maps", json: `[{"name":"cleanup","annotations":null,"labels":null}]`, want: emptyMetadata("default")},
		{name: "JSON empty label value", json: `[{"name":"cleanup","labels":{"empty":""}}]`, want: []CronJobDependency{
			{Name: "cleanup", Namespace: "default", Annotations: map[string]string{}, Labels: map[string]string{"empty": ""}},
		}},
		{name: "JSON precedence", plain: "old", json: `[{"name":"cleanup"}]`, want: emptyMetadata("default"), warning: "Ignoring TEST_CRONJOB since TEST_CRONJOB_JSON was specified"},
		{name: "empty JSON array takes precedence", plain: "old", json: `[]`, want: []CronJobDependency{}, warning: "Ignoring TEST_CRONJOB"},
		{name: "malformed JSON", json: `[{"name":`, want: []CronJobDependency{}, warning: "Invalid TEST_CRONJOB_JSON"},
		{name: "invalid JSON does not fall back to plain", plain: "old", json: `invalid`, want: []CronJobDependency{}, warning: "Invalid TEST_CRONJOB_JSON"},
		{name: "object instead of array", json: `{"name":"cleanup"}`, want: []CronJobDependency{}, warning: "expected an array"},
		{name: "null instead of array", json: `null`, want: []CronJobDependency{}, warning: "not null"},
		{name: "invalid name type", json: `[{"name":42}]`, want: []CronJobDependency{}, warning: "cannot unmarshal number"},
		{name: "invalid namespace type", json: `[{"name":"cleanup","namespace":true}]`, want: []CronJobDependency{}, warning: "cannot unmarshal bool"},
		{name: "null namespace", json: `[{"name":"cleanup","namespace":null}]`, want: []CronJobDependency{}, warning: "namespace must be a string, not null"},
		{name: "invalid annotation map", json: `[{"name":"cleanup","annotations":[]}]`, want: []CronJobDependency{}, warning: "cannot unmarshal array"},
		{name: "invalid annotation value", json: `[{"name":"cleanup","annotations":{"release":42}}]`, want: []CronJobDependency{}, warning: "cannot unmarshal number"},
		{name: "null annotation value", json: `[{"name":"cleanup","annotations":{"release":null}}]`, want: []CronJobDependency{}, warning: `annotation "release" must be a string, not null`},
		{name: "invalid label map", json: `[{"name":"cleanup","labels":[]}]`, want: []CronJobDependency{}, warning: "cannot unmarshal array"},
		{name: "invalid label value", json: `[{"name":"cleanup","labels":{"version":42}}]`, want: []CronJobDependency{}, warning: "cannot unmarshal number"},
		{name: "null label value", json: `[{"name":"cleanup","labels":{"version":null}}]`, want: []CronJobDependency{}, warning: `label "version" must be a string, not null`},
		{name: "invalid annotation discards entire JSON", json: `[{"name":"valid"},{"name":"cleanup","annotations":{"release":null}}]`, want: []CronJobDependency{}, warning: "dependency 2"},
		{name: "invalid label discards entire JSON", json: `[{"name":"valid"},{"name":"cleanup","labels":{"version":null}}]`, want: []CronJobDependency{}, warning: "dependency 2"},
		{name: "missing name", json: `[{"namespace":"application"}]`, want: []CronJobDependency{}, warning: "name is required"},
		{name: "label selection unsupported", json: `[{"labels":{"app":"cleanup"}}]`, want: []CronJobDependency{}, warning: "name is required"},
		{name: "empty name", json: `[{"name":""}]`, want: []CronJobDependency{}, warning: "name is required"},
		{name: "null name", json: `[{"name":null}]`, want: []CronJobDependency{}, warning: "name is required"},
		{name: "whitespace name", json: `[{"name":" "}]`, want: []CronJobDependency{}, warning: "name is required"},
		{name: "skip missing name only", json: `[{"name":""},{"name":"cleanup"}]`, want: emptyMetadata("default"), warning: "dependency 1: name is required"},
		{name: "plain missing name", plain: "application:", want: []CronJobDependency{}, warning: "name is required"},
		{name: "plain empty entry", plain: ",cleanup", want: []CronJobDependency{{Name: "cleanup", Namespace: "default"}}, warning: "name is required"},
		{name: "plain missing namespace", plain: ":cleanup", want: []CronJobDependency{}, warning: "missing namespace"},
		{name: "plain too many colons", plain: "a:b:c", want: []CronJobDependency{}, warning: "expected namespace:name"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("TEST_CRONJOB", test.plain)
			t.Setenv("TEST_CRONJOB_JSON", test.json)
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
			actual := SplitCronJobEnvToDeps("TEST_CRONJOB", "TEST_CRONJOB_JSON")
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

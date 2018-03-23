package env

import (
	"os"
	"reflect"
	"testing"
)

func TestSplitEnvToListWithColon(t *testing.T) {

	defer os.Unsetenv("TEST_LIST")

	os.Setenv("TEST_LIST", "foo,bar")
	list := SplitEnvToDeps("TEST_LIST")
	if list == nil {
		t.Errorf("Expected: not nil")
	}
	if list[0].Name != "foo" {
		t.Errorf("Expected: foo got %s", list[0])
	}
	if list[1].Name != "bar" {
		t.Errorf("Expected: bar got %s", list[1])
	}

	os.Setenv("TEST_LIST", "foo1")
	list1 := SplitEnvToDeps("TEST_LIST")
	if list1 == nil {
		t.Errorf("Expected: not nil")
	}
	if len(list1) != 1 {
		t.Errorf("Expected len to be 1 not %i", len(list1))
	}
	if list1[0].Name != "foo1" {
		t.Errorf("Expected: foo1 got %s", list1[0])
	}

	os.Setenv("TEST_LIST", "foo:foo")
	list2 := SplitEnvToDeps("TEST_LIST")
	if list2[0].Name != "foo" {
		t.Errorf("Expected: foo got %s", list2[0].Name)
	}
	if list2[0].Namespace != "foo" {
		t.Errorf("Expected: foo got %s", list2[0].Namespace)
	}

	os.Setenv("TEST_LIST", "bar")
	list3 := SplitEnvToDeps("TEST_LIST")
	if list3[0].Name != "bar" {
		t.Errorf("Expected: bar got %s", list3[0].Name)
	}
	if list3[0].Namespace != "default" {
		t.Errorf("Expected: default got %s", list3[0].Namespace)
	}

	os.Setenv("TEST_LIST", "foo:foo1:foo2")
	list4 := SplitEnvToDeps("TEST_LIST")
	if len(list4) != 0 {
		t.Errorf("Expected list to be empty")
	}

	os.Setenv("TEST_LIST", "foo:foo1:foo2,bar")
	list5 := SplitEnvToDeps("TEST_LIST")
	if list5[0].Namespace != "default" {
		t.Errorf("Expected: default got %s", list5[0].Namespace)
	}
	if list5[0].Name != "bar" {
		t.Errorf("Expected: bar got %s", list5[0].Name)
	}

	os.Setenv("TEST_LIST", "foo:foo1:foo2,bar:foo")
	list6 := SplitEnvToDeps("TEST_LIST")
	if list6[0].Namespace != "bar" {
		t.Errorf("Expected: bar got %s", list6[0].Namespace)
	}
	if list6[0].Name != "foo" {
		t.Errorf("Expected: foo got %s", list6[0].Name)
	}

	os.Setenv("TEST_LIST", ":foo")
	list7 := SplitEnvToDeps("TEST_LIST")
	if len(list7) != 0 {
		t.Errorf("Invalid format, missing namespace in pod")
	}
}

func TestSplitEmptyEnvWithColon(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "")
	list := SplitEnvToDeps("TEST_LIST")
	if list != nil {
		t.Errorf("Expected nil got %v", list)
	}
}

func TestSplitPodEnvToDepsSuccess(t *testing.T) {
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", `TEST_NAMESPACE`)
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", `[{"namespace": "foo", "labels": {"k1": "v1", "k2": "v2"}, "requireSameNode": true}, {"labels": {"k1": "v1", "k2": "v2"}}]`)
	actual := SplitPodEnvToDeps("TEST_LIST")
	expected := []PodDependency{
		PodDependency{Namespace: "foo", Labels: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}, RequireSameNode: true},
		PodDependency{Namespace: "TEST_NAMESPACE", Labels: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}, RequireSameNode: false},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v Got: %v", expected, actual)
	}
}

func TestSplitPodEnvToDepsUnset(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", "")
	actual := SplitPodEnvToDeps("TEST_LIST")
	if len(actual) != 0 {
		t.Errorf("Expected: no dependencies Got: %v", actual)
	}
}

func TestSplitPodEnvToDepsIgnoreInvalid(t *testing.T) {
	defer os.Unsetenv("TEST_LIST")
	os.Setenv("TEST_LIST", `[{"invalid": json}`)
	actual := SplitPodEnvToDeps("TEST_LIST")
	if len(actual) != 0 {
		t.Errorf("Expected: ignore invalid dependencies Got: %v", actual)
	}
}

func TestSplitCommand(t *testing.T) {
	defer os.Unsetenv("COMMAND")
	list2 := SplitCommand()
	if len(list2) > 0 {
		t.Errorf("Expected len to be 0, got %v", len(list2))
	}
	os.Setenv("COMMAND", "echo test")
	list := SplitCommand()
	if list == nil {
		t.Errorf("Expected slice, got nil")
		return
	}
	if len(list) != 2 {
		t.Errorf("Expected two elements, got %v", len(list))
	}
	if list[0] != "echo" {
		t.Errorf("Expected echo, got %s", list[0])
	}
	if list[1] != "test" {
		t.Errorf("Expected test, got %s", list[1])
	}

	os.Setenv("COMMAND", "")
	list1 := SplitCommand()
	if len(list1) > 0 {
		t.Errorf("Expected len to be 0, got %v", len(list1))
	}

}

func TestGetBaseNamespace(t *testing.T) {
	defer os.Unsetenv("NAMESPACE")
	os.Setenv("NAMESPACE", "")
	getBaseNamespace := GetBaseNamespace()
	if getBaseNamespace != "default" {
		t.Errorf("Expected namespace to be default, got %v", getBaseNamespace)
	}
	os.Setenv("NAMESPACE", "foo")
	getBaseNamespace = GetBaseNamespace()
	if getBaseNamespace != "foo" {
		t.Errorf("Expected namespace to be foo, got %v", getBaseNamespace)
	}
	os.Setenv("NAMESPACE", "default")
	getBaseNamespace = GetBaseNamespace()
	if getBaseNamespace != "default" {
		t.Errorf("Expected namespace to be default, got %v", getBaseNamespace)
	}
}

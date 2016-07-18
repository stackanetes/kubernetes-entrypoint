package env

import (
	"os"
	"testing"
)

func TestSplitEnvToListWithColon(t *testing.T) {
	os.Setenv("TEST_LIST", "foo,bar")
	list := SplitEnvToList("TEST_LIST")
	if list == nil {
		t.Errorf("Expected: not nil")
	}
	if list[0] != "foo" {
		t.Errorf("Expected: foo got %s", list[0])
	}
	if list[1] != "bar" {
		t.Errorf("Expected: bar got %s", list[1])
	}

	os.Setenv("TEST_LIST", "foo1")
	list1 := SplitEnvToList("TEST_LIST")
	if list1 == nil {
		t.Errorf("Expected: not nil")
	}
	if len(list1) != 1 {
		t.Errorf("Expected len to be 1 not %i", len(list1))
	}
	if list1[0] != "foo1" {
		t.Errorf("Expected: foo1 got %s", list1[0])
	}

}

func TestSplitEnvToListWithSpace(t *testing.T) {
	os.Setenv("TEST_LIST", "foo bar")
	list := SplitEnvToList("TEST_LIST", " ")
	if list == nil {
		t.Errorf("Expected: not nil")
	}
	if list[0] != "foo" {
		t.Errorf("Expected: foo got %s", list[0])
	}
	if list[1] != "bar" {
		t.Errorf("Expected: bar got %s", list[1])
	}

	os.Setenv("TEST_LIST", "foo1")
	list1 := SplitEnvToList("TEST_LIST", " ")
	if list1 == nil {
		t.Errorf("Expected: not nil")
	}
	if len(list1) != 1 {
		t.Errorf("Expected len to be 1 not %i", len(list1))
	}
	if list1[0] != "foo1" {
		t.Errorf("Expected: foo1 got %s", list1[0])
	}

}

func TestSplitEmptyEnvWithColon(t *testing.T) {
	os.Setenv("TEST_LIST", "")
	list := SplitEnvToList("TEST_LIST")
	if list != nil {
		t.Errorf("Expected nil got %v", list)
	}
}

func TestSplitEmptyEnvWithSpace(t *testing.T) {
	os.Setenv("TEST_LIST", "")
	list := SplitEnvToList("TEST_LIST", " ")
	if list != nil {
		t.Errorf("Expected nil got %v", list)
	}

}

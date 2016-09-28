package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	mocks "github.com/stackanetes/kubernetes-entrypoint/mocks"
)

func prepareEnv() (err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	ifaceName := ifaces[0].Name
	os.Setenv("INTERFACE_NAME", ifaceName)
	return nil

}

func createTemplate(template string) (err error) {
	configContent := []byte("LGTM {{ .HOSTNAME }}\n")
	if err = createDirectory(template); err != nil {
		return fmt.Errorf("Couldn't create directory in tmp: %v", err)
	}

	if err = ioutil.WriteFile(template, configContent, 0644); err != nil {
		return err
	}
	return

}
func TestIsResolved(t *testing.T) {
	name := "/tmp/lgtm"
	template := "/tmp/templates/lgtm/lgtm"
	hostname, err := os.Hostname()
	if err != nil {
		t.Errorf("couldn't get hostname", err)
	}

	entry := mocks.NewEntrypoint()

	err = prepareEnv()
	if err != nil {
		t.Errorf("Something went wrong: %v", err)
	}

	if err = createTemplate(template); err != nil {
		t.Errorf("Couldn't create %s template: %v", template, err)
	}

	config := NewConfig(name, "/tmp/templates")

	if _, err := config.IsResolved(entry); err != nil {
		t.Errorf("Something went wrong: %v", err)
	}

	result, err := ioutil.ReadFile(name)
	if err != nil {
		t.Errorf("Something went wrong file %s: doesnt exist: %v", name, err)
	}

	expectedFile := fmt.Sprintf("LGTM %s", hostname)

	same := strings.Compare(strings.TrimRight(string(result[:]), "\n"), expectedFile)
	if same != 0 {
		t.Errorf("Expected: %s got: %s: same %v", expectedFile, string(result[:]), same)
	}
}

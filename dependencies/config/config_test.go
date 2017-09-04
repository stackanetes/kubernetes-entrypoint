package config_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	. "github.com/stackanetes/kubernetes-entrypoint/dependencies/config"
	"github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testDir        = "/tmp"
	interfaceName  = "INTERFACE_NAME"
	testConfigName = "KUBERNETES_ENTRYPOINT_TEST_CONFIG"

	testConfigContentsFormat = "TEST_CONFIG %s\n"

	// configPath       = "/tmp/lgtm"
	templatePrefix = "/tmp/templates"
)

var testEntrypoint entrypoint.EntrypointInterface
var testConfigContents string
var testConfigPath string
var testTemplatePath string
var hostname string

// var testClient cli.ClientInterface

func init() {
	var err error
	testConfigContents = fmt.Sprintf(testConfigContentsFormat, "{{ .HOSTNAME }}")

	testTemplatePath = fmt.Sprintf("%s/%s/%s", templatePrefix, testConfigName, testConfigName)
	testConfigPath = fmt.Sprintf("%s/%s", testDir, testConfigName)

	hostname, err = os.Hostname()

	if err != nil {
		fmt.Errorf("Could not get hostname", err)
	}
}

func setupOsEnvironment() (err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	ifaceName := ifaces[0].Name
	return os.Setenv(interfaceName, ifaceName)
}

func teardownOsEnvironment() (err error) {
	return os.Unsetenv(interfaceName)
}

func setupConfigTemplate(templatePath string) (err error) {
	configContent := []byte(testConfigContents)
	if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
		return err
	}

	if err = ioutil.WriteFile(templatePath, configContent, 0644); err != nil {
		return err
	}

	return
}

func teardownConfigTemplate(templatePath string) (err error) {
	if err := os.RemoveAll(templatePath); err != nil {
		return err
	}

	return
}

var _ = Describe("Config", func() {

	BeforeEach(func() {
		err := setupOsEnvironment()
		Expect(err).NotTo(HaveOccurred())

		err = setupConfigTemplate(testTemplatePath)
		Expect(err).NotTo(HaveOccurred())

		testEntrypoint = mocks.NewEntrypoint()
	})

	AfterEach(func() {
		err := teardownOsEnvironment()
		Expect(err).NotTo(HaveOccurred())

		err = teardownConfigTemplate(testTemplatePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("creates new config from file", func() {
		config, err := NewConfig(testConfigPath, templatePrefix)

		Expect(config).NotTo(Equal(nil))
		Expect(err).NotTo(HaveOccurred())
	})

	It("checks the name of a newly created config file", func() {
		config, _ := NewConfig(testConfigPath, templatePrefix)

		Expect(config.GetName()).To(Equal(testConfigPath))
	})

	It("checks the format of a newly created config file", func() {
		config, _ := NewConfig(testConfigPath, templatePrefix)
		config.IsResolved(testEntrypoint)

		result, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testDir, testConfigName))
		Expect(err).NotTo(HaveOccurred())

		expectedFile := fmt.Sprintf(testConfigContentsFormat, hostname)

		readConfig := string(result[:])
		Expect(readConfig).To(BeEquivalentTo(expectedFile))
	})

	It("checks resolution of a config", func() {
		config, _ := NewConfig(testConfigPath, templatePrefix)

		isResolved, err := config.IsResolved(testEntrypoint)

		Expect(isResolved).To(Equal(true))
		Expect(err).NotTo(HaveOccurred())
	})

})

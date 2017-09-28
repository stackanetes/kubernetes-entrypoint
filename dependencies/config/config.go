package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/util"
	"github.com/stackanetes/kubernetes-entrypoint/util/env"
)

const (
	configmapDirPrefix    = "/configmaps"
	NamespaceNotSupported = "Config doesn't accept namespace"
)

type configParams struct {
	HOSTNAME  string
	IP        string
	IP_ERLANG string
}

type Config struct {
	name   string
	params configParams
	prefix string
}

func init() {
	configEnv := fmt.Sprintf("%sCONFIG", entry.DependencyPrefix)
	if util.ContainsSeparator(configEnv, "Config") {
		logger.Error.Printf(NamespaceNotSupported)
		os.Exit(1)
	}
	if configDeps := env.SplitEnvToDeps(configEnv); len(configDeps) > 0 {
		for _, dep := range configDeps {
			config, err := NewConfig(dep.Name, configmapDirPrefix)
			if err != nil {
				logger.Error.Printf("Cannot initialize config dep: %v", err)
			}
			entry.Register(config)
		}
	}
}

func NewConfig(name string, prefix string) (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Cannot determine hostname: %v", err)
	}

	ip, err := util.GetIp()
	if err != nil {
		return nil, fmt.Errorf("Cannot get ip address: %v", err)
	}

	return &Config{
		name: name,
		params: configParams{
			IP:        ip,
			IP_ERLANG: strings.Replace(ip, ".", ",", -1),
			HOSTNAME:  hostname},
		prefix: prefix,
	}, nil
}

func (c Config) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	//Create directory to ensure it exists
	if err := createDirectory(c.name); err != nil {
		return false, fmt.Errorf("Couldn't create directory: %v", err)
	}
	if err := c.createAndTemplateConfig(); err != nil {
		return false, fmt.Errorf("Cannot template %s: %v", c.name, err)
	}
	return true, nil

}

func (c Config) createAndTemplateConfig() (err error) {
	config, err := os.Create(c.name)
	if err != nil {
		return err
	}
	file := filepath.Base(c.name)

	temp := template.Must(template.New(file).ParseFiles(getSrcConfig(c.prefix, file)))
	if err = temp.Execute(config, c.params); err != nil {
		return err
	}
	return
}

func getSrcConfig(prefix string, config string) (srcConfig string) {
	return fmt.Sprintf("%s/%s/%s", prefix, config, config)
}

func createDirectory(file string) error {
	return os.MkdirAll(filepath.Dir(file), 0755)
}

func (c Config) String() string {
	return fmt.Sprintf("Config %s", c.name)
}

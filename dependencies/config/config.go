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

const configmapDirPrefix = "/configmaps"

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
	if configDeps := env.SplitEnvToList(configEnv); len(configDeps) > 0 {
		for _, dep := range configDeps {
			config, err := NewConfig(dep, configmapDirPrefix)
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
	if err := createDirectory(c.GetName()); err != nil {
		return false, fmt.Errorf("Couldn't create directory: %v", err)
	}
	if err := createAndTemplateConfig(c.GetName(), c.params, c.prefix); err != nil {
		return false, fmt.Errorf("Cannot template %s: %v", c.GetName(), err)
	}
	return true, nil

}

func createAndTemplateConfig(name string, params configParams, prefix string) (err error) {
	config, err := os.Create(name)
	if err != nil {
		return err
	}
	file := filepath.Base(name)
	temp := template.Must(template.New(file).ParseFiles(getSrcConfig(prefix, file)))
	if err = temp.Execute(config, params); err != nil {
		return err
	}

	return
}
func (c Config) GetName() string {
	return c.name
}

func getSrcConfig(prefix string, config string) (srcConfig string) {
	return fmt.Sprintf("%s/%s/%s", prefix, config, config)
}

func createDirectory(file string) error {
	return os.MkdirAll(filepath.Dir(file), 0755)
}

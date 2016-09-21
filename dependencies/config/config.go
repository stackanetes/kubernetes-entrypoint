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

type configParams struct {
	HOSTNAME  string
	IP        string
	IP_ERLANG string
}

type Config struct {
	name   string
	params configParams
}

func init() {
	configEnv := fmt.Sprintf("%sCONFIG", entry.DependencyPrefix)
	if configDeps := env.SplitEnvToList(configEnv); len(configDeps) > 0 {
		for _, dep := range configDeps {
			entry.Register(NewConfig(dep))
		}
	}
}

func NewConfig(name string) Config {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error.Printf("Cannot determine hostname: %v", err)
		os.Exit(1)
	}

	ip, err := util.GetIp()
	if err != nil {
		logger.Error.Printf("Cannot get ip address: %v", err)
		os.Exit(1)
	}

	return Config{
		name: name,
		params: configParams{
			IP:        ip,
			IP_ERLANG: strings.Replace(ip, ".", ",", -1),
			HOSTNAME:  hostname},
	}
}

func (c Config) IsResolved(entrypoint *entry.Entrypoint) (bool, error) {
	//Create directory to ensure it exists
	err := createDirectory(c.GetName())
	if err != nil {
		return false, fmt.Errorf("Couldn't create directory: %v", err)
	}
	err = createAndTemplateConfig(c.GetName(), c.params)
	if err != nil {
		return false, fmt.Errorf("Cannot template %s: %v", c.GetName(), err)
	}
	return true, nil

}

func createAndTemplateConfig(name string, params configParams) (err error) {
	config, err := os.Create(name)
	if err != nil {
		return
	}
	file := filepath.Base(name)
	temp := template.Must(template.New(file).ParseFiles(fmt.Sprintf("/configmaps/%s/%s", file, file)))
	if err = temp.Execute(config, params); err != nil {
		return err
	}

	return
}
func (c Config) GetName() string {
	return c.name
}

func createDirectory(file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0644)
	if err != nil {
		return err
	}
	return nil
}

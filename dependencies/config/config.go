package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	entry "github.com/stackanetes/kubernetes-entrypoint/entrypoint"
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
			entry.Register(NewConfig(dep, configmapDirPrefix))
		}
	}
}

func NewConfig(name string, prefix string) Config {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Cannot determine hostname: %v", err))
	}

	ip, err := util.GetIp()
	if err != nil {
		panic(fmt.Sprintf("Cannot get ip address: %v", err))
	}

	return Config{
		name: name,
		params: configParams{
			IP:        ip,
			IP_ERLANG: strings.Replace(ip, ".", ",", -1),
			HOSTNAME:  hostname},
		prefix: prefix,
	}
}

func (c Config) IsResolved(entrypoint entry.EntrypointInterface) (bool, error) {
	var err error
	//Create directory to ensure it exists
	if err = createDirectory(c.GetName()); err != nil {
		return false, fmt.Errorf("Couldn't create directory: %v", err)
	}
	if err = createAndTemplateConfig(c.GetName(), c.params, c.prefix); err != nil {
		return false, fmt.Errorf("Cannot template %s: %v", c.GetName(), err)
	}
	return true, nil

}

func createAndTemplateConfig(name string, params configParams, prefix string) (err error) {
	config, err := os.Create(name)
	if err != nil {
		return
	}
	file := filepath.Base(name)
	fmt.Sprintf("SrcConfig: %s", getSrcConfig(prefix, file))
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
	srcConfig = fmt.Sprintf("%s/%s/%s", prefix, config, config)
	return
}

func createDirectory(file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}
	return nil
}

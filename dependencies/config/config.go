package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	entry "github.com/stackanetes/docker-entrypoint/dependencies"
	"github.com/stackanetes/docker-entrypoint/logger"
	"github.com/stackanetes/docker-entrypoint/util/env"
)

type Config struct {
	name   string
	params struct {
		iface     string
		HOSTNAME  string
		IP        string
		IP_ERLANG string
	}
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
	var config Config
	config.name = name
	iface := os.Getenv("INTERFACE_NAME")
	if iface == "" {
		logger.Error.Print("Environment variable INTERFACE_NAME not set")
		os.Exit(1)
	}
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		logger.Error.Print("Environment variable HOSTNAME not set")
	}
	config.params.HOSTNAME = hostname
	config.params.iface = iface
	i, err := net.InterfaceByName(iface)
	if err != nil {
		logger.Error.Printf("Cannot get iface: %v", err)
		os.Exit(1)
	}

	address, err := i.Addrs()
	if err != nil || len(address) == 0 {
		logger.Error.Printf("Cannot get ip: %v", err)
		os.Exit(1)
	}
	config.params.IP = strings.Split(address[0].String(), "/")[0]
	config.params.IP_ERLANG = strings.Replace(config.params.IP, ".", ",", -1)

	return config
}

func (c Config) IsResolved(entrypoint *entry.Entrypoint) (bool, error) {
	logger.Info.Print(c.GetName())
	err := CreateDirectory(c.GetName())
	if err != nil {
		return false, fmt.Errorf("Couldn't create directory: %v", err)
	}
	config, err := os.Create(c.GetName())
	if err != nil {
		return false, fmt.Errorf("Couldn't touch file %v: %v", c.GetName(), err)
	}
	file := filepath.Base(c.GetName())
	temp := template.Must(template.New(file).ParseFiles(fmt.Sprintf("/configmaps/%s/%s", file, file)))
	if err = temp.Execute(config, c.params); err != nil {
		return false, err
	}
	return true, nil

}

func (c Config) GetName() string {
	return c.name
}

func CreateDirectory(file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0644)
	if err != nil {
		return err
	}
	return nil
}

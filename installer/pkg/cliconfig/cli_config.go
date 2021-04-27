package cliconfig

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path"
)

type CliConfigItem struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type CFCliConfig struct {
	Contexts       map[string]CliConfigItem `yaml:"contexts"`
	CurrentContext string                   `yaml:"current-context"`
}

type cliConfig struct {
}

type CliConfig interface {
	GetCurrentConfig() (*CliConfigItem, error)
}

func NewCliConfig() CliConfig {
	return &cliConfig{}
}

func (cf *cliConfig) GetCurrentConfig() (*CliConfigItem, error) {
	currentUser, err := user.Current()

	if err != nil {
		return nil, err
	}

	configPath := path.Join(currentUser.HomeDir, ".cfconfig")

	data, err := ioutil.ReadFile(configPath)

	var config CFCliConfig

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}
	result := config.Contexts[config.CurrentContext]
	return &result, nil
}

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

type CliConfig struct {
	Contexts       map[string]CliConfigItem `yaml:"contexts"`
	CurrentContext string                   `yaml:"current-context"`
}

func GetCurrentConfig() (*CliConfigItem, error) {
	currentUser, err := user.Current()

	if err != nil {
		return nil, err
	}

	configPath := path.Join(currentUser.HomeDir, ".cfconfig")

	data, err := ioutil.ReadFile(configPath)

	var config CliConfig

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

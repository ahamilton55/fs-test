package utils

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	DefaultConfigFile = ".deploy/config.toml"
)

// Info from config file
type Config struct {
	Service      string
	PackagerType string
	PackagerArgs string
	DeployerType string
	DeployerArgs string
}

type CommandConfig struct {
	Service string
	Env     string
	Profile string
	Region  string

	Params map[string]string
	Config Config
}

func getConfigfile(configFile string) string {
	if configFile == "" {
		configFile = DefaultConfigFile
	}

	return configFile
}

// Reads info from config file
func ReadConfig(configFile string) (Config, error) {
	var config Config

	configFile = getConfigfile(configFile)
	if _, err := os.Stat(configFile); err != nil {
		return config, fmt.Errorf("Config file is missing: %s", configFile)
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return config, err
	}

	return config, nil
}

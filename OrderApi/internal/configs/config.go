package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC    GRPCConfig    `yaml:"grpc"`
	Clients ClientsConfig `yaml:"clients"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ClientsConfig struct {
	User ClientConfig `yaml:"user"`
	Spot ClientConfig `yaml:"spot"`
}

type ClientConfig struct {
	Addr string `yaml:"addr"`
}

func Load(
	configPath string,
) (
	*Config,
	error,
) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

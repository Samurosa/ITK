package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC     GRPCConfig `yaml:"grpc"`
	Postgres Postgres   `yaml:"postgres"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Postgres struct {
	Link       string `yaml:"link"`
	MaxRetries int    `yaml:"max_retries"`
}

func Load(
	configPath string,
) (*Config, error) {
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

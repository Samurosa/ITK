package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env       string        `yaml:"env"`
	Token_ttl time.Duration `yaml:"token_ttl"`
	GRPC      GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load(
	configPath string,
) *Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("error read config: ", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("error unmarshal config: ", err)
	}
	return &cfg
}

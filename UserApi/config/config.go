package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TokenTTl TokensTTL  `yaml:"token_ttl"`
	GRPC     GRPCConfig `yaml:"grpc"`
	Env      string     `yaml:"env"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type TokensTTL struct {
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
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

package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TokenTTl TokensTTL  `yaml:"token_ttl"`
	GRPC     GRPCConfig `yaml:"grpc"`
	Env      string     `yaml:"env"`
	Postgres Postgres   `yaml:"postgres"`
	Redis    Redis      `yaml:"redis"`
	Limiter  Limiter    `yaml:"limiter"`
}

type TokensTTL struct {
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Postgres struct {
	Link       string `yaml:"link"`
	MaxRetries int    `yaml:"max_retries"`
}

type Redis struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

type Limiter struct {
	Limit int           `yaml:"limit"`
	Timer time.Duration `yaml:"timer"`
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

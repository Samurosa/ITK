package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWTSecretPath string     `yaml:"jwt_secret_path"`
	GRPC          GRPCConfig `yaml:"grpc"`
	Postgres      Postgres   `yaml:"postgres"`
	Redis         Redis      `yaml:"redis"`
}

type Redis struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Postgres struct {
	Link           string `yaml:"link"`
	MaxRetries     int    `yaml:"max_retries"`
	MaxConnections int32  `yaml:"max_connections"`
	MinConnections int32  `yaml:"min_connections"`
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

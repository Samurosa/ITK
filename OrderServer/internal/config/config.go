package config

import (
	"os"
	"time"

	sharedredis "github.com/Samurosa/exchange-common/shared/redis"
	"gopkg.in/yaml.v3"
)

type Config struct {
	JWTSecretPath string     `yaml:"jwt_secret_path"`
	GRPC          GRPCConfig `yaml:"grpc"`
	Spot          Spot       `yaml:"spot"`
	Postgres      Postgres   `yaml:"postgres"`
	Redis         Redis      `yaml:"redis"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Spot struct {
	GRPCAddr string `yaml:"grpc_addr"`
}

type Postgres struct {
	Link           string `yaml:"link"`
	MaxRetries     int    `yaml:"max_retries"`
	MaxConnections int32  `yaml:"max_connections"`
	MinConnections int32  `yaml:"min_connections"`
}

type Redis = sharedredis.Config

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

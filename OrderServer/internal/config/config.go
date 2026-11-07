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
	Kafka         Kafka      `yaml:"kafka"`
}

type Kafka struct {
	Brokers  []string      `yaml:"brokers"`
	ClientID string        `yaml:"client_id"`
	Topics   KafkaTopics   `yaml:"topics"`
	Producer KafkaProducer `yaml:"producer"`
	Consumer KafkaConsumer `yaml:"consumer"`
	Security KafkaSecurity `yaml:"security"`
}

type KafkaTopics struct {
	Commands string `yaml:"commands"`
	Events   string `yaml:"events"`
}

type KafkaProducer struct {
	RequiredAcks string        `yaml:"required_acks"`
	Idempotent   bool          `yaml:"idempotent"`
	Compression  string        `yaml:"compression"`
	BatchTimeout time.Duration `yaml:"batch_timeout"`
	RetryMax     int           `yaml:"retry_max"`
}

type KafkaConsumer struct {
	GroupID     string `yaml:"group_id"`
	AutoCommit  bool   `yaml:"auto_commit"`
	StartOffset string `yaml:"start_offset"`
}

type KafkaSecurity struct {
	Protocol      string `yaml:"protocol"`
	SASLMechanism string `yaml:"sasl_mechanism"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
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

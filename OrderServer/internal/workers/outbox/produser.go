package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ITK_Code/m/v2/internal/config"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer interface {
	Publish(
		ctx context.Context,
		topic string,
		key string,
		payload []byte,
	) error
}

type Kafka struct {
	client *kgo.Client
}

func NewClient(cfg config.Kafka) (*Kafka, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers are not configured")
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(cfg.Producer.RetryMax),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.ProducerLinger(cfg.Producer.BatchTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := client.Ping(pingCtx); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping kafka: %w", err)
	}

	return &Kafka{
		client: client,
	}, nil
}

func (k *Kafka) Publish(
	ctx context.Context,
	topic string,
	key string,
	payload []byte,
) error {
	if k.client == nil {
		return errors.New("kafka client is not initialized")
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}

	if err := k.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("publish kafka message to topic %s: %w", topic, err)
	}

	return nil
}

func (k *Kafka) Close() {
	if k.client != nil {
		k.client.Close()
	}
}

package ports

import (
	"ITK_Code/m/v2/internal/core/models"
	"context"
	"time"
)

type Storage interface {
	Claim(
		ctx context.Context,
		workerID string,
		batchSize int,
		leaseDuration time.Duration,
	) ([]models.OutboxMessage, error)

	MarkPublished(
		ctx context.Context,
		messageID string,
		workerID string,
	) error

	MarkFailed(
		ctx context.Context,
		messageID string,
		workerID string,
		messageErr string,
		retryAfter time.Duration,
	) error
}

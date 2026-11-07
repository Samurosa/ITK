package outbox

import (
	"context"
	"time"

	"ITK_Code/m/v2/internal/workers/outbox/ports"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Worker struct {
	storage  ports.Storage
	producer Producer
	log      *zap.Logger

	workerID      string
	batchSize     int
	pollInterval  time.Duration
	leaseDuration time.Duration
}

func NewWorker(
	storage ports.Storage,
	producer Producer,
	log *zap.Logger,
) *Worker {
	return &Worker{
		storage:       storage,
		producer:      producer,
		log:           log.Named("outbox_worker"),
		workerID:      uuid.NewString(),
		batchSize:     100,
		pollInterval:  time.Second,
		leaseDuration: 30 * time.Second,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	w.log.Info("outbox worker started",
		zap.String("worker_id", w.workerID),
	)

	for {
		if ctx.Err() != nil {
			w.log.Info("outbox worker stopped")
			return
		}

		if err := w.processBatch(ctx); err != nil {
			if ctx.Err() != nil {
				break
			}

			w.log.Error("outbox batch processing failed", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			w.log.Info("outbox worker stopped")
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) error {
	messages, err := w.storage.Claim(
		ctx,
		w.workerID,
		w.batchSize,
		w.leaseDuration,
	)
	if err != nil {
		return err
	}

	for _, message := range messages {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := w.producer.Publish(
			ctx,
			"order.commands",
			message.OrderID,
			message.Payload,
		)

		if err != nil {
			w.log.Error(
				"failed to publish outbox message",
				zap.String("message_id", message.ID),
				zap.String("order_id", message.OrderID),
				zap.Error(err),
			)

			retryAfter := w.retryDelay(message.Attempts)

			if markErr := w.storage.MarkFailed(
				ctx,
				message.ID,
				w.workerID,
				err.Error(),
				retryAfter,
			); markErr != nil {
				w.log.Error(
					"failed to update outbox retry state",
					zap.String("message_id", message.ID),
					zap.Error(markErr),
				)
			}

			continue
		}

		if err := w.storage.MarkPublished(
			ctx,
			message.ID,
			w.workerID,
		); err != nil {
			w.log.Error(
				"failed to mark outbox message as published",
				zap.String("message_id", message.ID),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (w *Worker) retryDelay(attempts int) time.Duration {
	return time.Second * time.Duration(1<<min(attempts, 6))
}

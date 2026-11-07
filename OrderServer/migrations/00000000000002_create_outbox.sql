-- +goose Up

CREATE TABLE order_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    message_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_until TIMESTAMPTZ,
    locked_by UUID,

    CONSTRAINT chk_order_outbox_attempts
        CHECK (attempts >= 0)
);

CREATE INDEX idx_order_outbox_available
    ON order_outbox (next_attempt_at, created_at)
    WHERE published_at IS NULL;

CREATE INDEX idx_order_outbox_pending
    ON order_outbox (created_at, id)
    WHERE published_at IS NULL;

CREATE INDEX idx_order_outbox_order_id
    ON order_outbox (order_id);


-- +goose Down

DROP TABLE IF EXISTS order_outbox;
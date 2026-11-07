-- +goose Up

CREATE TABLE order_inbox (
    event_id UUID PRIMARY KEY,

    order_id UUID NOT NULL
        REFERENCES orders(id) ON DELETE RESTRICT,

    sequence BIGINT NOT NULL,
    event_type VARCHAR(50) NOT NULL,

    payload JSONB NOT NULL,

    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_order_inbox_sequence
        UNIQUE (order_id, sequence),

    CONSTRAINT chk_order_inbox_sequence
        CHECK (sequence > 0)
);

CREATE INDEX idx_order_inbox_order_id
    ON order_inbox (order_id);

-- +goose Down

DROP TABLE IF EXISTS order_inbox;
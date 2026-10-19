-- +goose Up
CREATE TABLE order_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL
    REFERENCES orders(id) ON DELETE RESTRICT,

    version BIGINT NOT NULL,

    event_type VARCHAR(50) NOT NULL,

    order_status VARCHAR(50) NOT NULL,
    filled_quantity NUMERIC(30, 18) NOT NULL,

    payload JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_order_outbox_version
    UNIQUE (order_id, version),

    CONSTRAINT chk_order_outbox_version
    CHECK (version > 0)
);


CREATE INDEX idx_order_outbox_order_version
    ON order_outbox (order_id, version);

-- +goose Down
DROP TABLE IF EXISTS order_outbox;


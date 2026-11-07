-- +goose Up

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    idempotency_key VARCHAR(50) NOT NULL,
    last_engine_sequence BIGINT NOT NULL DEFAULT 0,

    user_id UUID NOT NULL,
    spot_id UUID NOT NULL,

    order_side VARCHAR(50) NOT NULL,
    order_status VARCHAR(50) NOT NULL,

    price NUMERIC(30, 18) NOT NULL,
    price_currency VARCHAR(10) NOT NULL,

    quantity NUMERIC(30, 18) NOT NULL,
    quantity_currency VARCHAR(10) NOT NULL,
    filled_quantity NUMERIC(30, 18) NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_orders_user_idempotency
        UNIQUE (user_id, idempotency_key),

    CONSTRAINT chk_orders_price_positive
        CHECK (price > 0),

    CONSTRAINT chk_orders_quantity_positive
        CHECK (quantity > 0),

    CONSTRAINT chk_orders_filled_quantity
        CHECK (filled_quantity >= 0 AND filled_quantity <= quantity),

    CONSTRAINT chk_orders_last_engine_sequence
        CHECK (last_engine_sequence >= 0),

    CONSTRAINT chk_orders_side
        CHECK (
            order_side IN (
                           'ORDER_SIDE_BUY',
                           'ORDER_SIDE_SELL'
                )
            ),

    CONSTRAINT chk_orders_status
        CHECK (
            order_status IN (
                             'ORDER_STATUS_NEW',
                             'ORDER_STATUS_OPEN',
                             'ORDER_STATUS_PARTIALLY_FILLED',
                             'ORDER_STATUS_FILLED',
                             'ORDER_STATUS_CANCELED',
                             'ORDER_STATUS_REJECTED'
                )
            )
);

CREATE INDEX idx_orders_user_created_id
    ON orders (user_id, created_at DESC, id DESC);

-- +goose Down

DROP TABLE IF EXISTS orders;
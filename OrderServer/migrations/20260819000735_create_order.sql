-- +goose Up
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    idempotency_key VARCHAR(50) NOT NULL,

    user_id UUID NOT NULL,
    spot_id UUID NOT NULL,

    order_side VARCHAR(50) NOT NULL,
    order_status VARCHAR(50) NOT NULL,

    price NUMERIC(30, 18) NOT NULL,
    quantity NUMERIC(30, 18) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_orders_user_idempotency
        UNIQUE (user_id, idempotency_key),

    CONSTRAINT chk_orders_price_positive
        CHECK (price > 0),

    CONSTRAINT chk_orders_quantity_positive
        CHECK (quantity > 0),

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

CREATE INDEX idx_orders_created_at_id
    ON orders (created_at DESC, id DESC);

-- +goose Down
DROP TABLE orders;

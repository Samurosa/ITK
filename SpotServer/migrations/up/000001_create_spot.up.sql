CREATE TABLE spot (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    symbol VARCHAR(255) UNIQUE NOT NULL,

    base_asset VARCHAR(255) NOT NULL,
    quote_asset VARCHAR(255) NOT NULL,

    price_precision INT NOT NULL,
    quantity_precision INT NOT NULL,

    min_order_size NUMERIC(30,18) NOT NULL,
    max_order_size NUMERIC(30,18) NOT NULL,

    allowed_roles TEXT[] NOT NULL,

    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),

    status VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    disabled_at TIMESTAMPTZ,

    CONSTRAINT chk_spot_different_assets
        CHECK (base_asset <> quote_asset),

    CONSTRAINT chk_spot_min_order_positive
        CHECK (min_order_size > 0),

    CONSTRAINT chk_spot_max_order_positive
        CHECK (max_order_size > 0),

    CONSTRAINT chk_spot_order_size_range
        CHECK (min_order_size <= max_order_size),

    CONSTRAINT chk_spot_price_precision
        CHECK (price_precision BETWEEN 0 AND 18),

    CONSTRAINT chk_spot_quantity_precision
        CHECK (quantity_precision BETWEEN 0 AND 18),

    CONSTRAINT chk_spot_status
        CHECK (status IN ('SPOT_STATUS_ACTIVE', 'SPOT_STATUS_DISABLED')),

    CONSTRAINT uq_spot_asset_pair
        UNIQUE (base_asset, quote_asset)
);
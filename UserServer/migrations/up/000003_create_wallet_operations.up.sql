CREATE TABLE wallet_operations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(50) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    asset VARCHAR(20) NOT NULL,
    amount NUMERIC(30, 18) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
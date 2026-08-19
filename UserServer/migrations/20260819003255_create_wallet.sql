-- +goose Up
CREATE TABLE balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset VARCHAR(20) NOT NULL,
    available NUMERIC(30,18) NOT NULL DEFAULT 0,
    locked NUMERIC(30,18) NOT NULL DEFAULT 0,
    UNIQUE(user_id, asset)
);

-- +goose Down
DROP TABLE balances

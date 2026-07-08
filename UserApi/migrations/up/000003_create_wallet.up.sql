CREATE TABLE balances (
    user_id UUID NOT NULL,

    asset VARCHAR(20) NOT NULL,

    available NUMERIC(38,18) NOT NULL DEFAULT 0,
    locked NUMERIC(38,18) NOT NULL DEFAULT 0,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_balances_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,

    CONSTRAINT unique_user_asset
        UNIQUE(user_id, asset)
);
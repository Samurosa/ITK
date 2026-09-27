-- +goose Up
ALTER TABLE orders
    ADD COLUMN price_currency VARCHAR(10) NOT NULL DEFAULT '',
    ADD COLUMN quantity_currency VARCHAR(10) NOT NULL DEFAULT '',
    ADD COLUMN filled_quantity NUMERIC(30, 18) NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE orders
    DROP COLUMN filled_quantity,
    DROP COLUMN quantity_currency,
    DROP COLUMN price_currency;

-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN price DECIMAL(10,2) NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN price;

-- +goose StatementEnd

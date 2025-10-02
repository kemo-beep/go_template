-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN damny VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN damny;

-- +goose StatementEnd

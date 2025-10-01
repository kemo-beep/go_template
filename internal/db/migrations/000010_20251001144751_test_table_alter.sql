-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN cool VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN cool;

-- +goose StatementEnd

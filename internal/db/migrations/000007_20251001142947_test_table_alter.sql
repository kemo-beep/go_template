-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN col22 VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN col22;

-- +goose StatementEnd

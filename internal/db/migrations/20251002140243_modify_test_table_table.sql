-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN location VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN location;

-- +goose StatementEnd

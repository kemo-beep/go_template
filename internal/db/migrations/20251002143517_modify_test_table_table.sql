-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN sammy VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN sammy;

-- +goose StatementEnd

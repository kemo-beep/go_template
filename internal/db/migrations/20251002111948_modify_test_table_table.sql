-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN preferences VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN preferences;

-- +goose StatementEnd

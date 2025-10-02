-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN preferrence VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN preferrence;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN prefrence VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN prefrence;

-- +goose StatementEnd

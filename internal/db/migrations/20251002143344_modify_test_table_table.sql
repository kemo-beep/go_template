-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN zoom VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN zoom;

-- +goose StatementEnd

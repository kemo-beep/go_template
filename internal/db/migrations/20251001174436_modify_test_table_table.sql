-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN family_name VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN family_name;

-- +goose StatementEnd

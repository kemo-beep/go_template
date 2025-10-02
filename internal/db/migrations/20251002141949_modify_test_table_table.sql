-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN wow VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN wow;

-- +goose StatementEnd

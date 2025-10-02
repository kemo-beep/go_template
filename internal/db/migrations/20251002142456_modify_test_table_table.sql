-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN damn VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN damn;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ADD COLUMN switch VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN switch;

-- +goose StatementEnd

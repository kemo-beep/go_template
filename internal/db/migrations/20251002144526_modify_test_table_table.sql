-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table DROP COLUMN zoom;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- TODO: Restore dropped column zoom

-- +goose StatementEnd

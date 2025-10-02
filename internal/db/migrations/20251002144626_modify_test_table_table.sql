-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table ALTER COLUMN wowedit TYPE VARCHAR;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- TODO: Restore original column definition for wowedit

-- +goose StatementEnd

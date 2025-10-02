-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table RENAME COLUMN wow TO wowedit;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table RENAME COLUMN wowedit TO wow;

-- +goose StatementEnd

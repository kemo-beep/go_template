-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table RENAME COLUMN cooledit TO name;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table RENAME COLUMN name TO cooledit;

-- +goose StatementEnd

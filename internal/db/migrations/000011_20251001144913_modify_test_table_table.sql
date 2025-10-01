-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_table RENAME COLUMN cool TO cooledit;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_table RENAME COLUMN cooledit TO cool;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE test ADD COLUMN col1 VARCHAR(255);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test DROP COLUMN col1;

-- +goose StatementEnd

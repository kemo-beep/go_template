-- +goose Down
DROP INDEX IF EXISTS idx_files_is_public;
DROP INDEX IF EXISTS idx_files_deleted_at;
DROP INDEX IF EXISTS idx_files_r2_key;
DROP INDEX IF EXISTS idx_files_user_id;
DROP TABLE IF EXISTS files;

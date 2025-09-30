-- +goose Down
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

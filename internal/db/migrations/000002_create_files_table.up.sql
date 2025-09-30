-- +goose Up
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    r2_key VARCHAR(500) UNIQUE NOT NULL,
    r2_url VARCHAR(1000) NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_r2_key ON files(r2_key);
CREATE INDEX idx_files_deleted_at ON files(deleted_at);
CREATE INDEX idx_files_is_public ON files(is_public);

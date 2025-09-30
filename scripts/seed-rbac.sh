#!/bin/bash

# Seed RBAC data (roles and permissions)
# This script populates the database with default roles and permissions

set -e

# Database configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5433}"
DB_USER="${DB_USER:-app}"
DB_PASS="${DB_PASS:-secret}"
DB_NAME="${DB_NAME:-myapp}"

echo "🌱 Seeding RBAC data..."

# Function to execute SQL
execute_sql() {
    PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "$1"
}

# Create default roles
echo "Creating default roles..."
execute_sql "INSERT INTO roles (name, description) VALUES 
    ('admin', 'Full system administrator') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO roles (name, description) VALUES 
    ('user', 'Standard user role') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO roles (name, description) VALUES 
    ('moderator', 'Content moderator') ON CONFLICT (name) DO NOTHING;"

# Create default permissions
echo "Creating default permissions..."

# User permissions
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('users.read', 'Read user data', 'users', 'read') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('users.write', 'Create and update users', 'users', 'write') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('users.delete', 'Delete users', 'users', 'delete') ON CONFLICT (name) DO NOTHING;"

# File permissions
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('files.read', 'Read file data', 'files', 'read') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('files.write', 'Upload and update files', 'files', 'write') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('files.delete', 'Delete files', 'files', 'delete') ON CONFLICT (name) DO NOTHING;"

# Role permissions
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('roles.read', 'Read role data', 'roles', 'read') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('roles.write', 'Create and update roles', 'roles', 'write') ON CONFLICT (name) DO NOTHING;"

# Database permissions (for DB management UI)
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('database.read', 'View database schema and data', 'database', 'read') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('database.write', 'Modify database schema and data', 'database', 'write') ON CONFLICT (name) DO NOTHING;"
execute_sql "INSERT INTO permissions (name, description, resource, action) VALUES 
    ('database.delete', 'Delete database tables and data', 'database', 'delete') ON CONFLICT (name) DO NOTHING;"

# Assign permissions to roles
echo "Assigning permissions to roles..."

# Get role IDs
ADMIN_ROLE_ID=$(PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT id FROM roles WHERE name = 'admin';")
USER_ROLE_ID=$(PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT id FROM roles WHERE name = 'user';")
MOD_ROLE_ID=$(PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT id FROM roles WHERE name = 'moderator';")

# Admin gets all permissions
execute_sql "INSERT INTO role_permissions (role_id, permission_id) 
    SELECT $ADMIN_ROLE_ID, id FROM permissions ON CONFLICT DO NOTHING;"

# User gets basic read/write permissions
execute_sql "INSERT INTO role_permissions (role_id, permission_id) 
    SELECT $USER_ROLE_ID, id FROM permissions 
    WHERE name IN ('users.read', 'files.read', 'files.write') 
    ON CONFLICT DO NOTHING;"

# Moderator gets user and file management permissions
execute_sql "INSERT INTO role_permissions (role_id, permission_id) 
    SELECT $MOD_ROLE_ID, id FROM permissions 
    WHERE name IN ('users.read', 'users.write', 'files.read', 'files.write', 'files.delete') 
    ON CONFLICT DO NOTHING;"

echo "✅ RBAC data seeded successfully!"
echo ""
echo "Default roles created:"
echo "  - admin (full access)"
echo "  - user (basic access)"
echo "  - moderator (content management)"
echo ""
echo "To assign admin role to a user, run:"
echo "  INSERT INTO user_roles (user_id, role_id) VALUES (<user_id>, $ADMIN_ROLE_ID);"

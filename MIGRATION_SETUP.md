# Database Migration System Setup

This document explains how to set up the database migration system with Google Scripts integration.

## Overview

The migration system provides:
- **Automatic Migration File Generation**: Creates organized migration files using Google Apps Script
- **Database Migration Execution**: Executes migrations safely with rollback support
- **Migration History Tracking**: Tracks all migration attempts and their status
- **Google Drive Integration**: Stores migration files in Google Drive for version control

## Setup Instructions

### 1. Google Apps Script Setup

1. **Create a new Google Apps Script project**:
   - Go to [script.google.com](https://script.google.com)
   - Click "New Project"
   - Replace the default code with the content from `scripts/google-apps-script/migration-generator.js`

2. **Create a Google Drive folder for migrations**:
   - Go to [drive.google.com](https://drive.google.com)
   - Create a new folder called "Database Migrations"
   - Copy the folder ID from the URL
   - Update `CONFIG.MIGRATIONS_FOLDER_ID` in the script

3. **Deploy the script as a web app**:
   - In the Apps Script editor, click "Deploy" > "New deployment"
   - Choose "Web app" as the type
   - Set "Execute as" to "Me"
   - Set "Who has access" to "Anyone"
   - Click "Deploy"
   - Copy the web app URL

4. **Get an access token**:
   - Go to [console.developers.google.com](https://console.developers.google.com)
   - Create a new project or select existing one
   - Enable the Google Drive API
   - Create credentials (OAuth 2.0 Client ID)
   - Use the client ID and secret to get an access token

### 2. Backend Configuration

Update your `config/config.yaml` file:

```yaml
google_scripts:
  url: "https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec"
  access_token: "YOUR_ACCESS_TOKEN"
  project_id: "your-project-id"
  migrations_dir: "./migrations"
```

### 3. Database Migration

Run the database migration to create the migrations table:

```bash
# Run the migration
go run cmd/server/main.go migrate up

# Or if using the Makefile
make migrate-up
```

### 4. Frontend Integration

The frontend is already configured to use the new migration API. The alter table functionality will:

1. **Create Migration**: When user saves changes, creates a migration record
2. **Generate Files**: Calls Google Scripts to generate migration files
3. **Execute Migration**: Runs the migration against the database
4. **Track Status**: Updates migration status throughout the process

## API Endpoints

### Migration Management

- `POST /api/v1/migrations` - Create a new migration
- `GET /api/v1/migrations` - List all migrations
- `GET /api/v1/migrations/:id` - Get specific migration
- `POST /api/v1/migrations/:id/execute` - Execute migration
- `POST /api/v1/migrations/:id/rollback` - Rollback migration
- `GET /api/v1/migrations/history?table_name=table` - Get migration history for table

### Migration Status

- `GET /api/v1/migration-status/:id` - Get migration status

## Migration File Structure

Generated migration files follow this structure:

```
migrations/
├── migration_0001_2024-01-15T10-30-00Z_users.sql
├── migration_0001_2024-01-15T10-30-00Z_users_rollback.sql
└── migration_0001_2024-01-15T10-30-00Z_users_metadata.json
```

### File Contents

**Migration File** (`migration_XXXX_*.sql`):
```sql
-- =============================================
-- Migration: migration-id
-- Table: users
-- Version: 0001
-- Created: 2024-01-15T10:30:00Z
-- Created By: admin
-- =============================================

-- Migration Description:
-- This migration modifies the table structure for users
-- Changes: 2 column modifications

-- =============================================
-- Migration SQL
-- =============================================

ALTER TABLE users ADD COLUMN new_field VARCHAR(255);
ALTER TABLE users DROP COLUMN old_field;

-- =============================================
-- End of Migration
-- =============================================
```

**Rollback File** (`migration_XXXX_*_rollback.sql`):
```sql
-- =============================================
-- Rollback Migration: migration-id
-- Table: users
-- Version: 0001
-- Created: 2024-01-15T10:30:00Z
-- Created By: admin
-- =============================================

-- Rollback Description:
-- This rollback reverses the changes made to users
-- Original Changes: 2 column modifications

-- =============================================
-- Rollback SQL
-- =============================================

ALTER TABLE users DROP COLUMN new_field;
ALTER TABLE users ADD COLUMN old_field VARCHAR(255);

-- =============================================
-- End of Rollback
-- =============================================
```

**Metadata File** (`migration_XXXX_*_metadata.json`):
```json
{
  "migration_id": "migration-id",
  "table_name": "users",
  "version": "0001",
  "created_at": "2024-01-15T10:30:00Z",
  "created_by": "admin",
  "changes": [
    {
      "action": "add",
      "column_name": "new_field",
      "type": "VARCHAR(255)",
      "nullable": true
    },
    {
      "action": "drop",
      "column_name": "old_field"
    }
  ],
  "files": {
    "migration_file_id": "file-id-1",
    "rollback_file_id": "file-id-2"
  },
  "status": "generated",
  "project": "Your Project Name"
}
```

## Security Considerations

1. **Access Tokens**: Store access tokens securely and rotate them regularly
2. **Permissions**: Use least-privilege access for Google Drive and Apps Script
3. **Validation**: Always validate migration SQL before execution
4. **Backup**: Ensure database backups before running migrations
5. **Testing**: Test migrations in development environment first

## Troubleshooting

### Common Issues

1. **Google Scripts API Error**:
   - Check if the script URL is correct
   - Verify access token has proper permissions
   - Ensure the script is deployed as a web app

2. **Migration Execution Failed**:
   - Check database connection
   - Verify SQL syntax
   - Check for table locks or constraints

3. **File Generation Failed**:
   - Verify Google Drive folder permissions
   - Check if the folder ID is correct
   - Ensure access token has Drive API permissions

### Debug Mode

Enable debug logging by setting the log level to `debug` in your config:

```yaml
logging:
  level: "debug"
  format: "json"
```

## Monitoring

The system provides several monitoring capabilities:

1. **Migration Status**: Real-time status updates
2. **Error Tracking**: Detailed error messages for failed migrations
3. **History**: Complete migration history with timestamps
4. **File Tracking**: Links to generated migration files

## Best Practices

1. **Always test migrations** in development first
2. **Use transactions** for complex migrations
3. **Keep migrations small** and focused
4. **Document changes** in migration descriptions
5. **Regular backups** before production migrations
6. **Monitor migration status** during execution
7. **Have rollback plans** ready for critical changes

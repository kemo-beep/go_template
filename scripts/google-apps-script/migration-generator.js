/**
 * Google Apps Script for generating database migration files
 * This script receives migration data and creates organized migration files
 */

// Configuration
const CONFIG = {
    MIGRATIONS_FOLDER_ID: 'YOUR_GOOGLE_DRIVE_FOLDER_ID', // Replace with actual folder ID
    PROJECT_NAME: 'Your Project Name',
    MIGRATION_PREFIX: 'migration',
    VERSION_PREFIX: 'v'
};

/**
 * Main function to handle POST requests
 * @param {Object} e - The event object containing the request data
 * @return {Object} Response object
 */
function doPost(e) {
    try {
        // Parse the request data
        const data = JSON.parse(e.postData.contents);

        // Validate required fields
        if (!data.migration_id || !data.table_name || !data.sql_query) {
            return createResponse(400, 'Missing required fields', null);
        }

        // Generate migration file
        const result = generateMigrationFile(data);

        return createResponse(200, 'Migration file generated successfully', result);

    } catch (error) {
        console.error('Error in doPost:', error);
        return createResponse(500, 'Internal server error', { error: error.toString() });
    }
}

/**
 * Generate migration file and save to Google Drive
 * @param {Object} data - Migration data
 * @return {Object} Result object with file information
 */
function generateMigrationFile(data) {
    try {
        // Create timestamp for versioning
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const version = getNextVersion();

        // Generate file name
        const fileName = `${CONFIG.MIGRATION_PREFIX}_${version}_${timestamp}_${data.table_name}.sql`;

        // Generate file content
        const fileContent = generateFileContent(data, version);

        // Create the file in Google Drive
        const file = DriveApp.createFile(fileName, fileContent, MimeType.PLAIN_TEXT);

        // Move to migrations folder
        const folder = DriveApp.getFolderById(CONFIG.MIGRATIONS_FOLDER_ID);
        file.moveTo(folder);

        // Create rollback file
        const rollbackFileName = `${CONFIG.MIGRATION_PREFIX}_${version}_${timestamp}_${data.table_name}_rollback.sql`;
        const rollbackContent = generateRollbackContent(data, version);
        const rollbackFile = DriveApp.createFile(rollbackFileName, rollbackContent, MimeType.PLAIN_TEXT);
        rollbackFile.moveTo(folder);

        // Create metadata file
        const metadataFileName = `${CONFIG.MIGRATION_PREFIX}_${version}_${timestamp}_${data.table_name}_metadata.json`;
        const metadataContent = generateMetadataContent(data, version, file.getId(), rollbackFile.getId());
        const metadataFile = DriveApp.createFile(metadataFileName, metadataContent, MimeType.PLAIN_TEXT);
        metadataFile.moveTo(folder);

        return {
            migration_id: data.migration_id,
            version: version,
            files: {
                migration: {
                    name: fileName,
                    id: file.getId(),
                    url: file.getUrl()
                },
                rollback: {
                    name: rollbackFileName,
                    id: rollbackFile.getId(),
                    url: rollbackFile.getUrl()
                },
                metadata: {
                    name: metadataFileName,
                    id: metadataFile.getId(),
                    url: metadataFile.getUrl()
                }
            },
            created_at: new Date().toISOString()
        };

    } catch (error) {
        console.error('Error generating migration file:', error);
        throw error;
    }
}

/**
 * Generate the main migration file content
 * @param {Object} data - Migration data
 * @param {string} version - Migration version
 * @return {string} File content
 */
function generateFileContent(data, version) {
    const header = `-- =============================================
-- Migration: ${data.migration_id}
-- Table: ${data.table_name}
-- Version: ${version}
-- Created: ${new Date().toISOString()}
-- Created By: ${data.created_by}
-- =============================================

-- Migration Description:
-- This migration modifies the table structure for ${data.table_name}
-- Changes: ${data.changes ? data.changes.length : 0} column modifications

-- =============================================
-- Migration SQL
-- =============================================

`;

    const footer = `
-- =============================================
-- End of Migration
-- =============================================
`;

    return header + data.sql_query + footer;
}

/**
 * Generate rollback file content
 * @param {Object} data - Migration data
 * @param {string} version - Migration version
 * @return {string} Rollback content
 */
function generateRollbackContent(data, version) {
    const header = `-- =============================================
-- Rollback Migration: ${data.migration_id}
-- Table: ${data.table_name}
-- Version: ${version}
-- Created: ${new Date().toISOString()}
-- Created By: ${data.created_by}
-- =============================================

-- Rollback Description:
-- This rollback reverses the changes made to ${data.table_name}
-- Original Changes: ${data.changes ? data.changes.length : 0} column modifications

-- =============================================
-- Rollback SQL
-- =============================================

`;

    const footer = `
-- =============================================
-- End of Rollback
-- =============================================
`;

    return header + (data.rollback_sql || '-- No rollback SQL available') + footer;
}

/**
 * Generate metadata file content
 * @param {Object} data - Migration data
 * @param {string} version - Migration version
 * @param {string} migrationFileId - Migration file ID
 * @param {string} rollbackFileId - Rollback file ID
 * @return {string} Metadata content
 */
function generateMetadataContent(data, version, migrationFileId, rollbackFileId) {
    const metadata = {
        migration_id: data.migration_id,
        table_name: data.table_name,
        version: version,
        created_at: new Date().toISOString(),
        created_by: data.created_by,
        changes: data.changes || [],
        files: {
            migration_file_id: migrationFileId,
            rollback_file_id: rollbackFileId
        },
        status: 'generated',
        project: CONFIG.PROJECT_NAME
    };

    return JSON.stringify(metadata, null, 2);
}

/**
 * Get the next version number
 * @return {string} Next version number
 */
function getNextVersion() {
    try {
        // Get existing files in the migrations folder
        const folder = DriveApp.getFolderById(CONFIG.MIGRATIONS_FOLDER_ID);
        const files = folder.getFiles();

        let maxVersion = 0;

        while (files.hasNext()) {
            const file = files.next();
            const fileName = file.getName();

            // Check if it's a migration file
            if (fileName.startsWith(CONFIG.MIGRATION_PREFIX)) {
                const versionMatch = fileName.match(new RegExp(`${CONFIG.MIGRATION_PREFIX}_(\\d+)_`));
                if (versionMatch) {
                    const version = parseInt(versionMatch[1]);
                    if (version > maxVersion) {
                        maxVersion = version;
                    }
                }
            }
        }

        return (maxVersion + 1).toString().padStart(4, '0');

    } catch (error) {
        console.error('Error getting next version:', error);
        // Return timestamp-based version as fallback
        return Date.now().toString();
    }
}

/**
 * Create a standardized response
 * @param {number} statusCode - HTTP status code
 * @param {string} message - Response message
 * @param {Object} data - Response data
 * @return {Object} Response object
 */
function createResponse(statusCode, message, data) {
    return {
        statusCode: statusCode,
        message: message,
        data: data,
        timestamp: new Date().toISOString()
    };
}

/**
 * Test function for development
 */
function testMigrationGeneration() {
    const testData = {
        migration_id: 'test-migration-123',
        table_name: 'users',
        sql_query: 'ALTER TABLE users ADD COLUMN new_field VARCHAR(255);',
        rollback_sql: 'ALTER TABLE users DROP COLUMN new_field;',
        changes: [
            {
                action: 'add',
                column_name: 'new_field',
                type: 'VARCHAR(255)',
                nullable: true
            }
        ],
        created_by: 'test-user',
        created_at: new Date().toISOString()
    };

    const result = generateMigrationFile(testData);
    console.log('Test result:', result);
    return result;
}

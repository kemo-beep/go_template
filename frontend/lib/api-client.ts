import axios from 'axios';

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Handle auth errors
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Clear auth token and redirect to login
            localStorage.removeItem('auth_token');
            // Only redirect if we're not already on the login page
            if (window.location.pathname !== '/login') {
                window.location.href = '/login';
            }
        }
        return Promise.reject(error);
    }
);

// API Types
export interface User {
    id: string;
    email: string;
    name: string;
    is_active: boolean;
    is_admin: boolean;
    created_at: string;
    updated_at: string;
}

export interface FileRecord {
    id: number;
    file_name: string;
    file_size: number;
    file_type: string;
    r2_url: string;
    is_public: boolean;
    created_at: string;
}

export interface ApiLog {
    timestamp: string;
    method: string;
    path: string;
    status: number;
    latency: string;
    client_ip: string;
    user_agent: string;
    error?: string;
}

export interface Metric {
    timestamp: string;
    requests_per_minute: number;
    error_rate: number;
    avg_latency: number;
}

export interface Migration {
    id: string;
    table_name: string;
    sql_query: string;
    rollback_sql?: string;
    status: 'pending' | 'running' | 'completed' | 'failed' | 'rolled_back';
    error_message?: string;
    created_at: string;
    completed_at?: string;
    created_by: string;
}

export interface ColumnChange {
    action: 'add' | 'modify' | 'drop' | 'rename';
    column_name: string;
    new_name?: string;
    type?: string;
    nullable?: boolean;
    default_value?: string;
    is_primary_key?: boolean;
    is_foreign_key?: boolean;
    references?: string;
}

// API Functions
export const api = {
    // Auth
    login: (email: string, password: string) =>
        apiClient.post('/api/v1/auth/login', { email, password }),

    // Users
    getUsers: () => apiClient.get<User[]>('/api/v1/admin/users'),
    getUser: (id: string) => apiClient.get<User>(`/api/v1/admin/users/${id}`),
    createUser: (data: Partial<User>) => apiClient.post('/api/v1/admin/users', data),
    updateUser: (id: string, data: Partial<User>) => apiClient.put(`/api/v1/admin/users/${id}`, data),
    deleteUser: (id: string) => apiClient.delete(`/api/v1/admin/users/${id}`),
    resetPassword: (id: string, newPassword: string) =>
        apiClient.post(`/api/v1/admin/users/${id}/reset-password`, { password: newPassword }),

    // Files
    getFiles: () => apiClient.get('/api/v1/files'),
    getFile: (id: string) => apiClient.get<FileRecord>(`/api/v1/files/${id}`),
    uploadFile: (data: { file: File; path: string }) => {
        const formData = new FormData();
        formData.append('file', data.file);
        formData.append('path', data.path);
        return apiClient.post('/api/v1/files/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        });
    },
    deleteFile: (id: string) => apiClient.delete(`/api/v1/files/${id}`),
    getDownloadUrl: (id: string) => apiClient.get(`/api/v1/files/${id}/download`),
    generateSignedUrl: (id: string) => apiClient.get(`/api/v1/files/${id}/download`),

    // Database
    getTables: () => apiClient.get('/api/v1/admin/database/tables'),
    getTableData: (table: string, page = 1, limit = 100) =>
        apiClient.get(`/api/v1/admin/database/tables/${table}/data`, { params: { page, limit } }),
    executeQuery: (query: string) =>
        apiClient.post('/api/v1/admin/database/query', { query }),
    getSchema: () => apiClient.get('/api/v1/admin/database/schema'),
    getTableSchema: (table: string) => apiClient.get(`/api/v1/admin/database/tables/${table}/schema`),

    // Logs
    getLogs: (params?: { limit?: number; level?: string; since?: string }) =>
        apiClient.get<ApiLog[]>('/api/v1/admin/logs', { params }),
    getMetrics: (params?: { from?: string; to?: string; interval?: string }) =>
        apiClient.get<Metric[]>('/api/v1/admin/metrics', { params }),
    getSentryErrors: () => apiClient.get('/api/v1/admin/sentry-errors'),

    // Developer Tools
    runMigration: (direction: 'up' | 'down', version?: string) =>
        apiClient.post('/api/v1/admin/migrations/run', { direction, version }),
    getGooseMigrations: () => apiClient.get('/api/v1/admin/migrations'),
    getFeatureFlags: () => apiClient.get('/api/v1/admin/feature-flags'),
    createFeatureFlag: (name: string, enabled: boolean) =>
        apiClient.post('/api/v1/admin/feature-flags', { name, enabled }),
    updateFeatureFlag: (name: string, enabled: boolean) =>
        apiClient.put(`/api/v1/admin/feature-flags/${name}`, { enabled }),
    toggleFeatureFlag: (name: string, enabled: boolean) =>
        apiClient.put(`/api/v1/admin/feature-flags/${name}`, { enabled }),
    getBackgroundJobs: () => apiClient.get('/api/v1/admin/jobs'),
    runBackgroundJob: (job: string, params?: any) =>
        apiClient.post('/api/v1/admin/jobs/run', { job, params }),

    // Settings
    getSettings: () => apiClient.get('/api/v1/admin/settings'),
    updateSettings: (settings: any) => apiClient.put('/api/v1/admin/settings', settings),

    // Realtime
    getRealtimePresence: () => apiClient.get('/api/v1/realtime/presence'),
    getRealtimeStats: () => apiClient.get('/api/v1/realtime/stats'),

    // Migrations
    createMigration: (data: { table_name: string; changes: ColumnChange[]; requested_by: string }) =>
        apiClient.post<Migration>('/api/v1/migrations', data),
    getMigrations: (params?: { limit?: number; offset?: number }) =>
        apiClient.get<{ migrations: Migration[]; limit: number; offset: number }>('/api/v1/migrations', { params }),
    getMigration: (id: string) => apiClient.get<Migration>(`/api/v1/migrations/${id}`),
    getMigrationHistory: (tableName: string, params?: { limit?: number; offset?: number }) =>
        apiClient.get<{ migrations: Migration[]; table_name: string; limit: number; offset: number }>('/api/v1/migrations/history', {
            params: { table_name: tableName, ...params }
        }),
    getMigrationFile: (id: string) => apiClient.get<{ id: string; table_name: string; sql_query: string; rollback_sql?: string; status: string }>(`/api/v1/migrations/${id}/file`),
    validateMigration: (id: string) => apiClient.get<{ valid: boolean; warnings: string[]; errors: string[]; migration_id: string; table_name: string; status: string }>(`/api/v1/migrations/${id}/validate`),
    executeMigration: (id: string) => apiClient.post(`/api/v1/migrations/${id}/execute`, {}),
    rollbackMigration: (id: string) => apiClient.post(`/api/v1/migrations/${id}/rollback`, {}),
    getMigrationStatus: (id: string) => apiClient.get<{ id: string; status: string; error_message?: string; created_at: string; completed_at?: string }>(`/api/v1/migration-status/${id}`),

    // Table data operations
    insertTableRow: (tableName: string, row: Record<string, any>) =>
        apiClient.post(`/api/v1/admin/database/tables/${tableName}/rows`, row),
    updateTableRow: (tableName: string, pkValue: any, row: Record<string, any>) =>
        apiClient.put(`/api/v1/admin/database/tables/${tableName}/rows/${pkValue}`, row),
    deleteTableRow: (tableName: string, pkValue: any) =>
        apiClient.delete(`/api/v1/admin/database/tables/${tableName}/rows/${pkValue}`),

    // Auth validation
    validateToken: () => apiClient.get('/api/v1/auth/validate'),
};

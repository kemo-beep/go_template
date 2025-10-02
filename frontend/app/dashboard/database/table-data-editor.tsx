'use client';

import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { toast } from 'sonner';
import { Pencil, Trash2, Plus, Save, X, ChevronLeft, ChevronRight, BookOpen, RefreshCw, Loader2, Copy, ExternalLink, Code, Database, Type, Shield, Zap, Check, AlertCircle, Edit2 } from 'lucide-react';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger, SheetFooter } from '@/components/ui/sheet';
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import { api, ColumnChange } from '@/lib/api-client';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { cn } from '@/lib/utils';
import { ColumnEditorSheet } from '@/components/column-editor-sheet';

const COLUMN_TYPES = [
    'VARCHAR',
    'TEXT',
    'INTEGER',
    'BIGINT',
    'SMALLINT',
    'DECIMAL',
    'NUMERIC',
    'FLOAT',
    'DOUBLE',
    'BOOLEAN',
    'DATE',
    'TIME',
    'TIMESTAMP',
    'JSON',
    'JSONB',
    'UUID',
    'SERIAL',
    'BIGSERIAL',
];

interface TableDataEditorProps {
    tableName: string;
    onRefresh?: () => void;
}

interface ColumnInfo {
    name: string;
    type: string;
    nullable: boolean;
    default_value?: string | null;
    is_primary_key: boolean;
    is_foreign_key?: boolean;
    references?: string;
    references_table?: string;
    references_column?: string;
    on_delete?: 'CASCADE' | 'SET NULL' | 'RESTRICT' | 'NO ACTION';
    on_update?: 'CASCADE' | 'SET NULL' | 'RESTRICT' | 'NO ACTION';
    is_editing?: boolean;
    is_new?: boolean;
    original_name?: string;
    constraints?: {
        unique?: boolean;
        check?: string;
        not_null?: boolean;
    };
    indexes?: {
        name?: string;
        type?: 'BTREE' | 'HASH' | 'GIN' | 'GIST';
        unique?: boolean;
    };
}

interface RowData {
    [key: string]: any;
}

export default function TableDataEditor({ tableName, onRefresh }: TableDataEditorProps) {
    const [page, setPage] = useState(1);
    const [pageSize] = useState(50);
    const [editingRow, setEditingRow] = useState<RowData | null>(null);
    const [isAddingRow, setIsAddingRow] = useState(false);
    const [newRow, setNewRow] = useState<RowData>({});
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
    const [rowToDelete, setRowToDelete] = useState<RowData | null>(null);
    const [columns, setColumns] = useState<ColumnInfo[]>([]);
    const [originalColumns, setOriginalColumns] = useState<ColumnInfo[]>([]);
    const [hasChanges, setHasChanges] = useState(false);
    const [migrationStatus, setMigrationStatus] = useState<'idle' | 'running' | 'completed' | 'error'>('idle');
    const [migrationId, setMigrationId] = useState<string | null>(null);
    const [apiDocsOpen, setApiDocsOpen] = useState(false);
    const [columnEditorOpen, setColumnEditorOpen] = useState(false);
    const [editingColumnIndex, setEditingColumnIndex] = useState<number | null>(null);

    const queryClient = useQueryClient();

    // Reset page when table changes
    useEffect(() => {
        setPage(1);
        setEditingRow(null);
        setIsAddingRow(false);
    }, [tableName]);

    // Fetch table schema
    const { data: schemaData } = useQuery({
        queryKey: ['tableSchema', tableName],
        queryFn: () => api.getTableSchema(tableName),
        enabled: !!tableName,
    });

    // Initialize columns when schema is loaded
    useEffect(() => {
        if (schemaData?.data?.data?.columns) {
            const schemaColumns = schemaData.data.data.columns.map((col: Record<string, unknown>) => ({
                name: col.name as string,
                type: col.type as string,
                nullable: col.nullable as boolean,
                default_value: col.default_value as string | null,
                is_primary_key: (col.is_primary_key as boolean) || false,
                is_foreign_key: (col.is_foreign_key as boolean) || false,
                references: (col.references as string) || '',
                references_table: (col.references_table as string) || '',
                references_column: (col.references_column as string) || '',
                on_delete: (col.on_delete as 'CASCADE' | 'SET NULL' | 'RESTRICT' | 'NO ACTION') || 'NO ACTION',
                on_update: (col.on_update as 'CASCADE' | 'SET NULL' | 'RESTRICT' | 'NO ACTION') || 'NO ACTION',
                is_editing: false,
                is_new: false,
                original_name: col.name as string,
                constraints: {
                    unique: (col.unique as boolean) || false,
                    check: (col.check as string) || '',
                    not_null: !(col.nullable as boolean),
                },
                indexes: {
                    name: (col.index_name as string) || '',
                    type: (col.index_type as 'BTREE' | 'HASH' | 'GIN' | 'GIST') || 'BTREE',
                    unique: (col.index_unique as boolean) || false,
                },
            }));
            setColumns(schemaColumns);
            setOriginalColumns(JSON.parse(JSON.stringify(schemaColumns)));
        }
    }, [schemaData]);

    const primaryKey = columns.find(col => col.is_primary_key)?.name;

    // Check for changes
    useEffect(() => {
        const hasModifications = JSON.stringify(columns) !== JSON.stringify(originalColumns);
        setHasChanges(hasModifications);
    }, [columns, originalColumns]);

    // Fetch table data
    const { data: tableData, isLoading } = useQuery({
        queryKey: ['tableData', tableName, page],
        queryFn: () => api.getTableData(tableName, page, pageSize),
        enabled: !!tableName,
    });

    const rows: RowData[] = tableData?.data?.data?.data || tableData?.data?.rows || [];
    const totalRows = tableData?.data?.data?.total || tableData?.data?.total || 0;
    const totalPages = Math.ceil(totalRows / pageSize);

    // Initialize new row with default values (excluding auto-generated columns)
    useEffect(() => {
        if (isAddingRow && columns.length > 0) {
            const initialRow: RowData = {};
            columns.forEach(col => {
                // Skip auto-generated columns (primary keys with default nextval, created_at, updated_at)
                const isAutoGenerated = col.is_primary_key ||
                    col.default_value?.includes('nextval') ||
                    ['created_at', 'updated_at', 'deleted_at'].includes(col.name.toLowerCase());

                if (isAutoGenerated) {
                    return; // Don't include in new row form
                }

                if (col.default_value && !col.default_value.includes('nextval')) {
                    initialRow[col.name] = col.default_value;
                } else if (col.nullable) {
                    initialRow[col.name] = null;
                } else {
                    initialRow[col.name] = '';
                }
            });
            setNewRow(initialRow);
        }
    }, [isAddingRow, columns]);

    // Add row mutation
    const addRowMutation = useMutation({
        mutationFn: (row: RowData) => api.insertTableRow(tableName, row),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['tableData', tableName] });
            toast.success('Row added successfully');
            setIsAddingRow(false);
            setNewRow({});
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to add row');
        },
    });

    // Update row mutation
    const updateRowMutation = useMutation({
        mutationFn: ({ row, pkValue }: { row: RowData; pkValue: any }) =>
            api.updateTableRow(tableName, pkValue, row),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['tableData', tableName] });
            toast.success('Row updated successfully');
            setEditingRow(null);
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to update row');
        },
    });

    // Delete row mutation
    const deleteRowMutation = useMutation({
        mutationFn: (pkValue: any) => api.deleteTableRow(tableName, pkValue),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['tableData', tableName] });
            toast.success('Row deleted successfully');
            setDeleteConfirmOpen(false);
            setRowToDelete(null);
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to delete row');
        },
    });

    // Column management functions
    const addNewColumn = () => {
        setEditingColumnIndex(null);
        setColumnEditorOpen(true);
    };

    const startEditing = (index: number) => {
        setEditingColumnIndex(index);
        setColumnEditorOpen(true);
    };

    const handleColumnSave = (columnData: ColumnInfo) => {
        if (editingColumnIndex === null) {
            // Adding new column
            const newColumn: ColumnInfo = {
                ...columnData,
                is_editing: false,
                is_new: true,
                original_name: columnData.name,
            };
            setColumns(prev => [...prev, newColumn]);
            toast.success('Column added successfully');
        } else {
            // Editing existing column
            setColumns(prev => prev.map((col, i) =>
                i === editingColumnIndex ? { ...columnData, is_editing: false } : col
            ));
            toast.success('Column updated successfully');
        }
        setColumnEditorOpen(false);
        setEditingColumnIndex(null);
    };

    const handleColumnCancel = () => {
        setColumnEditorOpen(false);
        setEditingColumnIndex(null);
    };


    const deleteColumn = (index: number) => {
        const column = columns[index];
        if (column.is_primary_key) {
            toast.error('Cannot delete primary key column');
            return;
        }

        setColumns(prev => prev.filter((_, i) => i !== index));
        toast.success('Column deleted');
    };

    const updateColumn = (index: number, field: keyof ColumnInfo, value: unknown) => {
        setColumns(prev => prev.map((col, i) =>
            i === index ? { ...col, [field]: value } : col
        ));
    };

    // Migration mutation
    const migrationMutation = useMutation({
        mutationFn: async () => {
            setMigrationStatus('running');

            // Convert columns to changes format
            const changes: ColumnChange[] = [];

            // Find new columns
            const newColumns = columns.filter(col => col.is_new);
            newColumns.forEach(col => {
                changes.push({
                    action: 'add' as const,
                    column_name: col.name,
                    type: col.type,
                    nullable: col.nullable,
                    default_value: col.default_value || undefined,
                    is_primary_key: col.is_primary_key,
                    is_foreign_key: col.is_foreign_key,
                    references: col.references || '',
                    references_table: col.references_table || '',
                    references_column: col.references_column || '',
                    on_delete: col.on_delete || 'NO ACTION',
                    on_update: col.on_update || 'NO ACTION',
                    constraints: col.constraints,
                    indexes: col.indexes,
                });
            });

            // Find renamed columns
            const renamedColumns = columns.filter(col =>
                col.original_name && col.name !== col.original_name && !col.is_new
            );
            renamedColumns.forEach(col => {
                changes.push({
                    action: 'rename' as const,
                    column_name: col.original_name!,
                    new_name: col.name,
                });
            });

            // Find modified columns
            const modifiedColumns = columns.filter(col =>
                !col.is_new &&
                col.original_name &&
                col.name === col.original_name &&
                originalColumns.some(orig => orig.name === col.name)
            );
            modifiedColumns.forEach(col => {
                const originalCol = originalColumns.find(orig => orig.name === col.name);
                if (originalCol && (
                    originalCol.type !== col.type ||
                    originalCol.nullable !== col.nullable ||
                    originalCol.default_value !== col.default_value ||
                    originalCol.is_primary_key !== col.is_primary_key ||
                    originalCol.is_foreign_key !== col.is_foreign_key ||
                    originalCol.references_table !== col.references_table ||
                    originalCol.references_column !== col.references_column ||
                    originalCol.on_delete !== col.on_delete ||
                    originalCol.on_update !== col.on_update ||
                    JSON.stringify(originalCol.constraints) !== JSON.stringify(col.constraints) ||
                    JSON.stringify(originalCol.indexes) !== JSON.stringify(col.indexes)
                )) {
                    changes.push({
                        action: 'modify' as const,
                        column_name: col.name,
                        type: col.type,
                        nullable: col.nullable,
                        default_value: col.default_value || undefined,
                        is_primary_key: col.is_primary_key,
                        is_foreign_key: col.is_foreign_key,
                        references: col.references || '',
                        references_table: col.references_table || '',
                        references_column: col.references_column || '',
                        on_delete: col.on_delete || 'NO ACTION',
                        on_update: col.on_update || 'NO ACTION',
                        constraints: col.constraints,
                        indexes: col.indexes,
                    });
                }
            });

            // Find deleted columns
            const deletedColumns = originalColumns.filter(origCol =>
                !columns.some(col => col.original_name === origCol.name || col.name === origCol.name)
            );
            deletedColumns.forEach(col => {
                changes.push({
                    action: 'drop' as const,
                    column_name: col.name,
                });
            });

            console.log('Changes to migrate:', changes);

            if (changes.length === 0) {
                throw new Error('No changes to migrate');
            }

            // Create migration
            const response = await api.createMigration({
                table_name: tableName,
                changes: changes,
                requested_by: 'current-user',
            });

            console.log('Migration response:', response);
            const migrationId = response.data?.id;
            console.log('Migration ID:', migrationId);

            if (!migrationId) {
                console.error('Full response structure:', JSON.stringify(response, null, 2));
                throw new Error('Failed to get migration ID from response');
            }

            setMigrationId(migrationId);

            // Execute migration
            await api.executeMigration(migrationId);

            return response;
        },
        onSuccess: () => {
            setMigrationStatus('completed');
            toast.success('Migration completed successfully');
            queryClient.invalidateQueries({ queryKey: ['tableSchema', tableName] });
            queryClient.invalidateQueries({ queryKey: ['tableData', tableName] });
            onRefresh?.();
        },
        onError: (error: unknown) => {
            setMigrationStatus('error');
            const errorMessage = error instanceof Error ? error.message : 'Migration failed';
            toast.error(errorMessage);
        },
    });

    const handleSave = () => {
        if (!hasChanges) {
            toast.info('No changes to save');
            return;
        }
        migrationMutation.mutate();
    };

    const handleEdit = (row: RowData) => {
        setEditingRow({ ...row });
    };

    const handleSaveEdit = () => {
        if (editingRow && primaryKey) {
            updateRowMutation.mutate({
                row: editingRow,
                pkValue: editingRow[primaryKey],
            });
        }
    };

    const handleCancelEdit = () => {
        setEditingRow(null);
    };

    const handleAddRow = () => {
        setIsAddingRow(true);
    };

    const handleSaveNewRow = () => {
        addRowMutation.mutate(newRow);
    };

    const handleCancelAddRow = () => {
        setIsAddingRow(false);
        setNewRow({});
    };

    const handleDeleteClick = (row: RowData) => {
        setRowToDelete(row);
        setDeleteConfirmOpen(true);
    };

    const handleConfirmDelete = () => {
        if (rowToDelete && primaryKey) {
            deleteRowMutation.mutate(rowToDelete[primaryKey]);
        }
    };


    const handleRefresh = () => {
        queryClient.invalidateQueries({ queryKey: ['tableData', tableName] });
        queryClient.invalidateQueries({ queryKey: ['tableSchema', tableName] });
        onRefresh?.();
        toast.success('Table refreshed');
    };

    const copyToClipboard = async (text: string) => {
        try {
            await navigator.clipboard.writeText(text);
            toast.success('Copied to clipboard');
        } catch (error) {
            toast.error('Failed to copy to clipboard');
        }
    };

    const renderCellValue = (value: any) => {
        if (value === null) {
            return <Badge variant="secondary">NULL</Badge>;
        }
        if (typeof value === 'boolean') {
            return <Badge variant={value ? 'default' : 'outline'}>{value.toString()}</Badge>;
        }
        if (typeof value === 'object') {
            return JSON.stringify(value);
        }
        return String(value);
    };

    const renderEditableCell = (row: RowData, column: ColumnInfo, isNew: boolean = false) => {
        const value = row[column.name];
        const isEditing = primaryKey && editingRow?.[primaryKey] === row[primaryKey];
        const displayRow = isNew ? newRow : (isEditing ? editingRow : row);

        // Check if column is auto-generated
        const isAutoGenerated = column.is_primary_key ||
            column.default_value?.includes('nextval') ||
            ['created_at', 'updated_at', 'deleted_at'].includes(column.name.toLowerCase());

        // For new rows, don't show auto-generated columns
        if (isNew && isAutoGenerated) {
            return <span className="text-muted-foreground text-xs">Auto</span>;
        }

        // For existing rows, show auto-generated columns as read-only
        if (isAutoGenerated && !isNew) {
            return <span className="text-muted-foreground">{renderCellValue(value)}</span>;
        }

        // Safety check for displayRow
        if (!displayRow) {
            return <span className="text-muted-foreground">-</span>;
        }

        return (
            <Input
                type={column.type.includes('INT') ? 'number' : 'text'}
                value={displayRow[column.name] ?? ''}
                onChange={(e) => {
                    const newValue = e.target.value;
                    if (isNew) {
                        setNewRow(prev => ({ ...prev, [column.name]: newValue }));
                    } else {
                        setEditingRow(prev => prev ? { ...prev, [column.name]: newValue } : null);
                    }
                }}
                className="h-8"
                placeholder={column.nullable ? 'NULL' : ''}
            />
        );
    };

    if (!tableName) {
        return (
            <div className="flex items-center justify-center h-64 text-muted-foreground">
                <p>Select a table to view and edit data</p>
            </div>
        );
    }

    if (!columns || columns.length === 0) {
        return (
            <div className="flex items-center justify-center h-64">
                <p className="text-muted-foreground">Loading table schema...</p>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-lg font-semibold">{tableName}</h3>
                    <p className="text-sm text-muted-foreground">
                        {totalRows} rows total
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={handleRefresh}
                                    disabled={isLoading}
                                >
                                    <RefreshCw className={cn("h-4 w-4", isLoading && "animate-spin")} />
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>Refresh table data</p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>

                    <Sheet open={apiDocsOpen} onOpenChange={setApiDocsOpen}>
                        <SheetTrigger asChild>
                            <Button variant="outline" size="sm">
                                <BookOpen className="h-4 w-4" />
                            </Button>
                        </SheetTrigger>
                        <SheetContent className="w-[900px] sm:max-w-[900px] overflow-y-auto px-6">
                            <SheetHeader>
                                <SheetTitle className="flex items-center gap-2">
                                    <Database className="h-5 w-5" />
                                    API Documentation - {tableName}
                                </SheetTitle>
                                <SheetDescription>
                                    Complete API reference and documentation for the {tableName} table
                                </SheetDescription>
                            </SheetHeader>

                            <div className="mt-6 space-y-6">
                                {/* Table Overview */}
                                <div className="p-4 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-950/20 dark:to-indigo-950/20 rounded-lg border">
                                    <div className="flex items-center gap-2 mb-2">
                                        <Shield className="h-4 w-4 text-blue-600" />
                                        <h3 className="font-semibold text-blue-900 dark:text-blue-100">Table Overview</h3>
                                    </div>
                                    <div className="grid grid-cols-2 gap-4 text-sm">
                                        <div>
                                            <span className="text-muted-foreground">Table Name:</span>
                                            <span className="ml-2 font-mono">{tableName}</span>
                                        </div>
                                        <div>
                                            <span className="text-muted-foreground">Total Rows:</span>
                                            <span className="ml-2 font-semibold">{totalRows.toLocaleString()}</span>
                                        </div>
                                        <div>
                                            <span className="text-muted-foreground">Columns:</span>
                                            <span className="ml-2 font-semibold">{columns.length}</span>
                                        </div>
                                        <div>
                                            <span className="text-muted-foreground">Base URL:</span>
                                            <span className="ml-2 font-mono text-xs">/api/v1/tables/{tableName}</span>
                                        </div>
                                    </div>
                                </div>

                                {/* API Endpoints */}
                                <div className="space-y-4">
                                    <h3 className="text-lg font-semibold flex items-center gap-2">
                                        <Zap className="h-5 w-5" />
                                        API Endpoints
                                    </h3>

                                    <Accordion type="multiple" className="w-full">
                                        {/* GET All Rows */}
                                        <AccordionItem value="get-all" className="border rounded-lg">
                                            <AccordionTrigger className="px-4 py-3 hover:no-underline">
                                                <div className="flex items-center gap-3">
                                                    <Badge variant="outline" className="bg-green-100 text-green-800 border-green-200">
                                                        GET
                                                    </Badge>
                                                    <code className="text-sm font-mono">/api/v1/tables/{tableName}</code>
                                                    <span className="text-sm text-muted-foreground">Get all rows</span>
                                                </div>
                                            </AccordionTrigger>
                                            <AccordionContent className="px-4 pb-4">
                                                <div className="space-y-4">
                                                    <div>
                                                        <h4 className="font-medium mb-2">Description</h4>
                                                        <p className="text-sm text-muted-foreground">
                                                            Retrieve all rows from the {tableName} table with optional pagination and filtering.
                                                        </p>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Query Parameters</h4>
                                                        <div className="space-y-2 text-sm">
                                                            <div className="flex justify-between items-center p-2 bg-muted rounded">
                                                                <span className="font-mono">page</span>
                                                                <span className="text-muted-foreground">integer (optional)</span>
                                                            </div>
                                                            <div className="flex justify-between items-center p-2 bg-muted rounded">
                                                                <span className="font-mono">limit</span>
                                                                <span className="text-muted-foreground">integer (optional, default: 50)</span>
                                                            </div>
                                                            <div className="flex justify-between items-center p-2 bg-muted rounded">
                                                                <span className="font-mono">sort</span>
                                                                <span className="text-muted-foreground">string (optional)</span>
                                                            </div>
                                                        </div>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Example Request</h4>
                                                        <div className="relative">
                                                            <pre className="bg-gray-900 text-gray-100 p-3 rounded-lg text-xs overflow-x-auto">
                                                                <code>{`curl -X GET "http://localhost:8080/api/v1/tables/${tableName}?page=1&limit=10" \\
  -H "Authorization: Bearer YOUR_TOKEN"`}</code>
                                                            </pre>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                className="absolute top-2 right-2 h-6 w-6 p-0"
                                                                onClick={() => copyToClipboard(`curl -X GET "http://localhost:8080/api/v1/tables/${tableName}?page=1&limit=10" \\\n  -H "Authorization: Bearer YOUR_TOKEN"`)}
                                                            >
                                                                <Copy className="h-3 w-3" />
                                                            </Button>
                                                        </div>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Response</h4>
                                                        <pre className="bg-gray-900 text-gray-100 p-3 rounded-lg text-xs overflow-x-auto">
                                                            <code>{`{
  "data": {
    "data": [...],
    "total": ${totalRows},
    "page": 1,
    "limit": 10
  },
  "success": true
}`}</code>
                                                        </pre>
                                                    </div>
                                                </div>
                                            </AccordionContent>
                                        </AccordionItem>

                                        {/* POST Create Row */}
                                        <AccordionItem value="post-create" className="border rounded-lg">
                                            <AccordionTrigger className="px-4 py-3 hover:no-underline">
                                                <div className="flex items-center gap-3">
                                                    <Badge variant="outline" className="bg-blue-100 text-blue-800 border-blue-200">
                                                        POST
                                                    </Badge>
                                                    <code className="text-sm font-mono">/api/v1/tables/{tableName}/rows</code>
                                                    <span className="text-sm text-muted-foreground">Create new row</span>
                                                </div>
                                            </AccordionTrigger>
                                            <AccordionContent className="px-4 pb-4">
                                                <div className="space-y-4">
                                                    <div>
                                                        <h4 className="font-medium mb-2">Description</h4>
                                                        <p className="text-sm text-muted-foreground">
                                                            Create a new row in the {tableName} table.
                                                        </p>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Request Body</h4>
                                                        <div className="space-y-2 text-sm">
                                                            {columns.filter(col => !col.is_primary_key && !col.default_value?.includes('nextval')).map((col) => (
                                                                <div key={col.name} className="flex justify-between items-center p-2 bg-muted rounded">
                                                                    <span className="font-mono">{col.name}</span>
                                                                    <div className="flex items-center gap-2">
                                                                        <span className="text-muted-foreground">{col.type}</span>
                                                                        {!col.nullable && <Badge variant="destructive" className="text-xs">required</Badge>}
                                                                    </div>
                                                                </div>
                                                            ))}
                                                        </div>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Example Request</h4>
                                                        <div className="relative">
                                                            <pre className="bg-gray-900 text-gray-100 p-3 rounded-lg text-xs overflow-x-auto">
                                                                <code>{`curl -X POST "http://localhost:8080/api/v1/tables/${tableName}/rows" \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{
    ${columns.filter(col => !col.is_primary_key && !col.default_value?.includes('nextval')).slice(0, 3).map(col => `"${col.name}": "value"`).join(',\n    ')}
  }'`}</code>
                                                            </pre>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                className="absolute top-2 right-2 h-6 w-6 p-0"
                                                                onClick={() => copyToClipboard(`curl -X POST "http://localhost:8080/api/v1/tables/${tableName}/rows" \\\n  -H "Content-Type: application/json" \\\n  -H "Authorization: Bearer YOUR_TOKEN" \\\n  -d '{\n    ${columns.filter(col => !col.is_primary_key && !col.default_value?.includes('nextval')).slice(0, 3).map(col => `"${col.name}": "value"`).join(',\n    ')}\n  }'`)}
                                                            >
                                                                <Copy className="h-3 w-3" />
                                                            </Button>
                                                        </div>
                                                    </div>
                                                </div>
                                            </AccordionContent>
                                        </AccordionItem>

                                        {/* PUT Update Row */}
                                        <AccordionItem value="put-update" className="border rounded-lg">
                                            <AccordionTrigger className="px-4 py-3 hover:no-underline">
                                                <div className="flex items-center gap-3">
                                                    <Badge variant="outline" className="bg-yellow-100 text-yellow-800 border-yellow-200">
                                                        PUT
                                                    </Badge>
                                                    <code className="text-sm font-mono">/api/v1/tables/{tableName}/rows/:id</code>
                                                    <span className="text-sm text-muted-foreground">Update row</span>
                                                </div>
                                            </AccordionTrigger>
                                            <AccordionContent className="px-4 pb-4">
                                                <div className="space-y-4">
                                                    <div>
                                                        <h4 className="font-medium mb-2">Description</h4>
                                                        <p className="text-sm text-muted-foreground">
                                                            Update an existing row in the {tableName} table by ID.
                                                        </p>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Path Parameters</h4>
                                                        <div className="space-y-2 text-sm">
                                                            <div className="flex justify-between items-center p-2 bg-muted rounded">
                                                                <span className="font-mono">id</span>
                                                                <span className="text-muted-foreground">string (required) - Row ID</span>
                                                            </div>
                                                        </div>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Example Request</h4>
                                                        <div className="relative">
                                                            <pre className="bg-gray-900 text-gray-100 p-3 rounded-lg text-xs overflow-x-auto">
                                                                <code>{`curl -X PUT "http://localhost:8080/api/v1/tables/${tableName}/rows/123" \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{
    ${columns.filter(col => !col.is_primary_key && !col.default_value?.includes('nextval')).slice(0, 3).map(col => `"${col.name}": "updated_value"`).join(',\n    ')}
  }'`}</code>
                                                            </pre>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                className="absolute top-2 right-2 h-6 w-6 p-0"
                                                                onClick={() => copyToClipboard(`curl -X PUT "http://localhost:8080/api/v1/tables/${tableName}/rows/123" \\\n  -H "Content-Type: application/json" \\\n  -H "Authorization: Bearer YOUR_TOKEN" \\\n  -d '{\n    ${columns.filter(col => !col.is_primary_key && !col.default_value?.includes('nextval')).slice(0, 3).map(col => `"${col.name}": "updated_value"`).join(',\n    ')}\n  }'`)}
                                                            >
                                                                <Copy className="h-3 w-3" />
                                                            </Button>
                                                        </div>
                                                    </div>
                                                </div>
                                            </AccordionContent>
                                        </AccordionItem>

                                        {/* DELETE Row */}
                                        <AccordionItem value="delete-row" className="border rounded-lg">
                                            <AccordionTrigger className="px-4 py-3 hover:no-underline">
                                                <div className="flex items-center gap-3">
                                                    <Badge variant="outline" className="bg-red-100 text-red-800 border-red-200">
                                                        DELETE
                                                    </Badge>
                                                    <code className="text-sm font-mono">/api/v1/tables/{tableName}/rows/:id</code>
                                                    <span className="text-sm text-muted-foreground">Delete row</span>
                                                </div>
                                            </AccordionTrigger>
                                            <AccordionContent className="px-4 pb-4">
                                                <div className="space-y-4">
                                                    <div>
                                                        <h4 className="font-medium mb-2">Description</h4>
                                                        <p className="text-sm text-muted-foreground">
                                                            Delete a row from the {tableName} table by ID.
                                                        </p>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Path Parameters</h4>
                                                        <div className="space-y-2 text-sm">
                                                            <div className="flex justify-between items-center p-2 bg-muted rounded">
                                                                <span className="font-mono">id</span>
                                                                <span className="text-muted-foreground">string (required) - Row ID</span>
                                                            </div>
                                                        </div>
                                                    </div>

                                                    <div>
                                                        <h4 className="font-medium mb-2">Example Request</h4>
                                                        <div className="relative">
                                                            <pre className="bg-gray-900 text-gray-100 p-3 rounded-lg text-xs overflow-x-auto">
                                                                <code>{`curl -X DELETE "http://localhost:8080/api/v1/tables/${tableName}/rows/123" \\
  -H "Authorization: Bearer YOUR_TOKEN"`}</code>
                                                            </pre>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                className="absolute top-2 right-2 h-6 w-6 p-0"
                                                                onClick={() => copyToClipboard(`curl -X DELETE "http://localhost:8080/api/v1/tables/${tableName}/rows/123" \\\n  -H "Authorization: Bearer YOUR_TOKEN"`)}
                                                            >
                                                                <Copy className="h-3 w-3" />
                                                            </Button>
                                                        </div>
                                                    </div>
                                                </div>
                                            </AccordionContent>
                                        </AccordionItem>
                                    </Accordion>
                                </div>

                                {/* Schema Information */}
                                <div className="space-y-4">
                                    <h3 className="text-lg font-semibold flex items-center gap-2">
                                        <Type className="h-5 w-5" />
                                        Table Schema
                                    </h3>

                                    <div className="border rounded-lg overflow-hidden">
                                        <Table>
                                            <TableHeader>
                                                <TableRow>
                                                    <TableHead>Column Name</TableHead>
                                                    <TableHead>Type</TableHead>
                                                    <TableHead>Nullable</TableHead>
                                                    <TableHead>Default</TableHead>
                                                    <TableHead>Constraints</TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {columns.map((col) => (
                                                    <TableRow key={col.name}>
                                                        <TableCell className="font-mono">{col.name}</TableCell>
                                                        <TableCell>
                                                            <Badge variant="outline" className="font-mono text-xs">
                                                                {col.type}
                                                            </Badge>
                                                        </TableCell>
                                                        <TableCell>
                                                            <Badge variant={col.nullable ? "secondary" : "destructive"} className="text-xs">
                                                                {col.nullable ? "Yes" : "No"}
                                                            </Badge>
                                                        </TableCell>
                                                        <TableCell className="text-sm text-muted-foreground">
                                                            {col.default_value || "NULL"}
                                                        </TableCell>
                                                        <TableCell>
                                                            <div className="flex gap-1">
                                                                {col.is_primary_key && (
                                                                    <Badge variant="default" className="text-xs">PK</Badge>
                                                                )}
                                                                {col.is_foreign_key && (
                                                                    <Badge variant="outline" className="text-xs">FK</Badge>
                                                                )}
                                                            </div>
                                                        </TableCell>
                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>
                            </div>
                        </SheetContent>
                    </Sheet>

                    <Button
                        variant="outline"
                        onClick={addNewColumn}
                        disabled={migrationStatus === 'running'}
                    >
                        <Plus className="h-4 w-4 mr-2" />
                        Add Column
                    </Button>

                    {hasChanges && (
                        <Button
                            onClick={handleSave}
                            disabled={migrationStatus === 'running'}
                        >
                            {migrationStatus === 'running' ? (
                                <>
                                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                                    Running Migration...
                                </>
                            ) : (
                                <>
                                    <Save className="h-4 w-4 mr-2" />
                                    Save Changes
                                </>
                            )}
                        </Button>
                    )}

                    <Button onClick={handleAddRow} disabled={isAddingRow}>
                        <Plus className="h-4 w-4 mr-2" />
                        Add Row
                    </Button>
                </div>
            </div>

            {/* Migration Status */}
            {migrationStatus !== 'idle' && (
                <div className={cn(
                    "p-3 rounded-lg border flex items-center gap-2",
                    migrationStatus === 'running' && "bg-blue-50 border-blue-200",
                    migrationStatus === 'completed' && "bg-green-50 border-green-200",
                    migrationStatus === 'error' && "bg-red-50 border-red-200"
                )}>
                    {migrationStatus === 'running' && <Loader2 className="h-4 w-4 animate-spin" />}
                    {migrationStatus === 'completed' && <Check className="h-4 w-4 text-green-600" />}
                    {migrationStatus === 'error' && <AlertCircle className="h-4 w-4 text-red-600" />}
                    <span className="text-sm font-medium">
                        {migrationStatus === 'running' && 'Migration in progress...'}
                        {migrationStatus === 'completed' && 'Migration completed successfully'}
                        {migrationStatus === 'error' && 'Migration failed'}
                    </span>
                </div>
            )}

            <Card>
                <CardContent className="p-0">
                    <div className="overflow-x-auto">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    {columns.map((column, index) => (
                                        <TableHead key={index} className="w-[200px]">
                                            {column.is_editing ? (
                                                <Input
                                                    value={column.name}
                                                    onChange={(e) => updateColumn(index, 'name', e.target.value)}
                                                    placeholder="Column name"
                                                    className="h-8"
                                                />
                                            ) : (
                                                <div className="flex items-center gap-2">
                                                    <span className="font-medium">{column.name}</span>
                                                    {column.is_primary_key && (
                                                        <Badge variant="outline" className="text-xs">PK</Badge>
                                                    )}
                                                    {column.is_new && <Badge variant="secondary" className="text-xs">New</Badge>}
                                                </div>
                                            )}
                                        </TableHead>
                                    ))}
                                    <TableHead className="w-24">Actions</TableHead>
                                </TableRow>
                                <TableRow>
                                    {columns.map((column, index) => (
                                        <TableHead key={`type-${index}`} className="w-[200px]">
                                            <div className="space-y-1">
                                                <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                                                    {column.type}
                                                </code>
                                                <div className="flex gap-1">
                                                    <Badge variant={column.nullable ? 'secondary' : 'outline'} className="text-xs">
                                                        {column.nullable ? 'Nullable' : 'Not Null'}
                                                    </Badge>
                                                </div>
                                                <span className="text-xs text-gray-500">
                                                    {column.default_value || 'No default'}
                                                </span>
                                            </div>
                                        </TableHead>
                                    ))}
                                    <TableHead className="w-24">Type & Properties</TableHead>
                                </TableRow>
                                <TableRow>
                                    {columns.map((column, index) => (
                                        <TableHead key={`constraints-${index}`} className="w-[200px]">
                                            <div className="space-y-1">
                                                <div className="flex gap-1 flex-wrap">
                                                    {column.is_primary_key && (
                                                        <Badge variant="default" className="text-xs">PK</Badge>
                                                    )}
                                                    {column.is_foreign_key && (
                                                        <Badge variant="outline" className="text-xs">FK</Badge>
                                                    )}
                                                    {column.constraints?.unique && (
                                                        <Badge variant="secondary" className="text-xs">UNIQUE</Badge>
                                                    )}
                                                    {column.constraints?.check && (
                                                        <Badge variant="outline" className="text-xs">CHECK</Badge>
                                                    )}
                                                    {column.indexes?.name && (
                                                        <Badge variant="outline" className="text-xs">INDEX</Badge>
                                                    )}
                                                </div>
                                                {column.is_foreign_key && column.references_table && (
                                                    <div className="text-xs text-muted-foreground">
                                                        → {column.references_table}.{column.references_column}
                                                    </div>
                                                )}
                                                {column.constraints?.check && (
                                                    <div className="text-xs text-muted-foreground">
                                                        {column.constraints.check}
                                                    </div>
                                                )}
                                            </div>
                                        </TableHead>
                                    ))}
                                    <TableHead className="w-24">Keys & Constraints</TableHead>
                                </TableRow>
                                <TableRow>
                                    {columns.map((column, index) => (
                                        <TableHead key={`actions-${index}`} className="w-[200px]">
                                            <div className="flex gap-1">
                                                <Button
                                                    size="sm"
                                                    variant="ghost"
                                                    onClick={() => startEditing(index)}
                                                    className="h-6 w-6 p-0"
                                                >
                                                    <Edit2 className="h-3 w-3" />
                                                </Button>
                                                {!column.is_primary_key && (
                                                    <Button
                                                        size="sm"
                                                        variant="ghost"
                                                        onClick={() => deleteColumn(index)}
                                                        className="h-6 w-6 p-0 text-red-600 hover:text-red-700"
                                                    >
                                                        <Trash2 className="h-3 w-3" />
                                                    </Button>
                                                )}
                                            </div>
                                        </TableHead>
                                    ))}
                                    <TableHead className="w-24">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {isAddingRow && (
                                    <TableRow className="bg-muted/50">
                                        {columns.map((column) => (
                                            <TableCell key={column.name}>
                                                {renderEditableCell(newRow, column, true)}
                                            </TableCell>
                                        ))}
                                        <TableCell>
                                            <div className="flex items-center gap-1">
                                                <Button
                                                    size="sm"
                                                    variant="ghost"
                                                    onClick={handleSaveNewRow}
                                                    disabled={addRowMutation.isPending}
                                                >
                                                    <Save className="h-4 w-4" />
                                                </Button>
                                                <Button
                                                    size="sm"
                                                    variant="ghost"
                                                    onClick={handleCancelAddRow}
                                                >
                                                    <X className="h-4 w-4" />
                                                </Button>
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                )}
                                {isLoading ? (
                                    <TableRow>
                                        <TableCell colSpan={columns.length + 1} className="text-center">
                                            Loading...
                                        </TableCell>
                                    </TableRow>
                                ) : rows.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={columns.length + 1} className="text-center text-muted-foreground">
                                            No data available
                                        </TableCell>
                                    </TableRow>
                                ) : (
                                    rows.map((row, idx) => {
                                        const isEditing = editingRow?.[primaryKey!] === row[primaryKey!];
                                        return (
                                            <TableRow key={idx} className={isEditing ? 'bg-muted/50' : ''}>
                                                {columns.map((column) => (
                                                    <TableCell key={column.name}>
                                                        {isEditing ? (
                                                            renderEditableCell(row, column)
                                                        ) : (
                                                            renderCellValue(row[column.name])
                                                        )}
                                                    </TableCell>
                                                ))}
                                                <TableCell>
                                                    {isEditing ? (
                                                        <div className="flex items-center gap-1">
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                onClick={handleSaveEdit}
                                                                disabled={updateRowMutation.isPending}
                                                            >
                                                                <Save className="h-4 w-4" />
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                onClick={handleCancelEdit}
                                                            >
                                                                <X className="h-4 w-4" />
                                                            </Button>
                                                        </div>
                                                    ) : (
                                                        <div className="flex items-center gap-1">
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                onClick={() => handleEdit(row)}
                                                            >
                                                                <Pencil className="h-4 w-4" />
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                onClick={() => handleDeleteClick(row)}
                                                            >
                                                                <Trash2 className="h-4 w-4 text-destructive" />
                                                            </Button>
                                                        </div>
                                                    )}
                                                </TableCell>
                                            </TableRow>
                                        );
                                    })
                                )}
                            </TableBody>
                        </Table>
                    </div>
                </CardContent>
            </Card>

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex items-center justify-between">
                    <p className="text-sm text-muted-foreground">
                        Page {page} of {totalPages}
                    </p>
                    <div className="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage(p => Math.max(1, p - 1))}
                            disabled={page === 1}
                        >
                            <ChevronLeft className="h-4 w-4" />
                            Previous
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                            disabled={page === totalPages}
                        >
                            Next
                            <ChevronRight className="h-4 w-4" />
                        </Button>
                    </div>
                </div>
            )}

            {/* Delete Confirmation Dialog */}
            <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Confirm Delete</DialogTitle>
                        <DialogDescription>
                            Are you sure you want to delete this row? This action cannot be undone.
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setDeleteConfirmOpen(false)}>
                            Cancel
                        </Button>
                        <Button
                            variant="destructive"
                            onClick={handleConfirmDelete}
                            disabled={deleteRowMutation.isPending}
                        >
                            Delete
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Column Editor Sheet */}
            <ColumnEditorSheet
                isOpen={columnEditorOpen}
                onClose={handleColumnCancel}
                onSave={handleColumnSave}
                column={editingColumnIndex !== null ? columns[editingColumnIndex] : null}
                isEditing={editingColumnIndex !== null}
                existingColumns={columns}
            />

        </div>
    );
}


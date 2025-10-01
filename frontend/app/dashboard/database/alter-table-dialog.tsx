'use client';

import { useState, useEffect } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import {
    Plus,
    Trash2,
    Edit2,
    Save,
    X,
    Check,
    AlertCircle,
    Database,
    ArrowLeft,
    Loader2,
} from 'lucide-react';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';

interface Column {
    name: string;
    type: string;
    nullable: boolean;
    default_value: string | null;
    is_primary_key: boolean;
    is_foreign_key: boolean;
    references?: string;
    is_editing?: boolean;
    is_new?: boolean;
    original_name?: string;
}

interface AlterTableDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    tableName: string;
}

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


export function AlterTableDialog({ open, onOpenChange, tableName }: AlterTableDialogProps) {
    const [columns, setColumns] = useState<Column[]>([]);
    const [originalColumns, setOriginalColumns] = useState<Column[]>([]);
    const [hasChanges, setHasChanges] = useState(false);
    const [migrationStatus, setMigrationStatus] = useState<'idle' | 'running' | 'completed' | 'error'>('idle');
    const [migrationId, setMigrationId] = useState<string | null>(null);

    // Fetch table schema
    const { data: schema, isLoading: schemaLoading } = useQuery({
        queryKey: ['table-schema', tableName],
        queryFn: () => api.getTableSchema(tableName),
        enabled: open && !!tableName,
    });

    // Initialize columns when schema is loaded
    useEffect(() => {
        if (schema?.data?.data?.columns) {
            const schemaColumns = schema.data.data.columns.map((col: Record<string, unknown>) => ({
                name: col.name as string,
                type: col.type as string,
                nullable: col.nullable as boolean,
                default_value: col.default_value as string | null,
                is_primary_key: (col.is_primary_key as boolean) || false,
                is_foreign_key: (col.is_foreign_key as boolean) || false,
                references: (col.references as string) || '',
                is_editing: false,
                is_new: false,
                original_name: col.name as string,
            }));
            setColumns(schemaColumns);
            setOriginalColumns(JSON.parse(JSON.stringify(schemaColumns)));
        }
    }, [schema]);

    // Check for changes
    useEffect(() => {
        const hasModifications = JSON.stringify(columns) !== JSON.stringify(originalColumns);
        setHasChanges(hasModifications);
    }, [columns, originalColumns]);

    const addNewColumn = () => {
        const newColumn: Column = {
            name: '',
            type: 'VARCHAR',
            nullable: true,
            default_value: null,
            is_primary_key: false,
            is_foreign_key: false,
            references: '',
            is_editing: true,
            is_new: true,
        };
        setColumns(prev => [...prev, newColumn]);
    };

    const startEditing = (index: number) => {
        setColumns(prev => prev.map((col, i) =>
            i === index ? { ...col, is_editing: true } : col
        ));
    };

    const cancelEditing = (index: number) => {
        const column = columns[index];
        if (column.is_new) {
            // Remove new column if canceling
            setColumns(prev => prev.filter((_, i) => i !== index));
        } else {
            // Restore original values
            const originalColumn = originalColumns[index];
            setColumns(prev => prev.map((col, i) =>
                i === index ? { ...originalColumn, is_editing: false } : col
            ));
        }
    };

    const saveColumn = (index: number) => {
        const column = columns[index];

        // Validation
        if (!column.name.trim()) {
            toast.error('Column name is required');
            return;
        }

        // Check for duplicate names
        const duplicateIndex = columns.findIndex((col, i) =>
            i !== index && col.name.toLowerCase() === column.name.toLowerCase()
        );
        if (duplicateIndex !== -1) {
            toast.error('Column name must be unique');
            return;
        }

        setColumns(prev => prev.map((col, i) =>
            i === index ? { ...col, is_editing: false } : col
        ));
        toast.success('Column saved');
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

    const updateColumn = (index: number, field: keyof Column, value: unknown) => {
        setColumns(prev => prev.map((col, i) =>
            i === index ? { ...col, [field]: value } : col
        ));
    };

    const generateMigration = () => {
        const changes = [];

        // Find new columns
        const newColumns = columns.filter(col => col.is_new);
        if (newColumns.length > 0) {
            changes.push(`-- Adding new columns`);
            newColumns.forEach(col => {
                const typeWithLength = col.type === 'VARCHAR' && col.default_value
                    ? `${col.type}(${col.default_value})`
                    : col.type;
                changes.push(`ALTER TABLE ${tableName} ADD COLUMN ${col.name} ${typeWithLength}${col.nullable ? '' : ' NOT NULL'};`);
            });
        }

        // Find modified columns
        const modifiedColumns = columns.filter((col, index) =>
            !col.is_new &&
            col.original_name &&
            JSON.stringify(col) !== JSON.stringify(originalColumns[index])
        );

        if (modifiedColumns.length > 0) {
            changes.push(`-- Modifying existing columns`);
            modifiedColumns.forEach(col => {
                if (col.name !== col.original_name) {
                    changes.push(`ALTER TABLE ${tableName} RENAME COLUMN ${col.original_name} TO ${col.name};`);
                }
                // Add more modification logic as needed
            });
        }

        // Find deleted columns
        const deletedColumns = originalColumns.filter(origCol =>
            !columns.some(col => col.original_name === origCol.name)
        );

        if (deletedColumns.length > 0) {
            changes.push(`-- Dropping columns`);
            deletedColumns.forEach(col => {
                changes.push(`ALTER TABLE ${tableName} DROP COLUMN ${col.name};`);
            });
        }

        return changes.join('\n');
    };

    const migrationMutation = useMutation({
        mutationFn: async (migrationSQL: string) => {
            setMigrationStatus('running');
            const response = await api.executeQuery(migrationSQL);
            setMigrationId(response.data?.migration_id || Date.now().toString());
            return response;
        },
        onSuccess: () => {
            setMigrationStatus('completed');
            toast.success('Migration completed successfully');
            // Refresh schema
            window.location.reload();
        },
        onError: (error: unknown) => {
            setMigrationStatus('error');
            const errorMessage = error instanceof Error ? error.message : 'Migration failed';
            toast.error(errorMessage);
        },
    });

    const rollbackMutation = useMutation({
        mutationFn: async () => {
            if (!migrationId) throw new Error('No migration to rollback');
            // This would call a rollback API endpoint
            toast.info('Rollback functionality coming soon');
        },
        onError: (error: unknown) => {
            const errorMessage = error instanceof Error ? error.message : 'Rollback failed';
            toast.error(errorMessage);
        },
    });

    const handleSave = () => {
        const migrationSQL = generateMigration();
        if (!migrationSQL.trim()) {
            toast.info('No changes to save');
            return;
        }
        migrationMutation.mutate(migrationSQL);
    };

    const handleRollback = () => {
        if (migrationStatus === 'completed' && migrationId) {
            rollbackMutation.mutate();
        }
    };

    const resetChanges = () => {
        setColumns(JSON.parse(JSON.stringify(originalColumns)));
        toast.info('Changes reset');
    };

    if (schemaLoading) {
        return (
            <Dialog open={open} onOpenChange={onOpenChange}>
                <DialogContent className="max-w-4xl max-h-[80vh] overflow-hidden">
                    <div className="flex items-center justify-center py-8">
                        <Loader2 className="h-8 w-8 animate-spin" />
                        <span className="ml-2">Loading table schema...</span>
                    </div>
                </DialogContent>
            </Dialog>
        );
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-4xl max-h-[80vh] overflow-hidden flex flex-col">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Database className="h-5 w-5" />
                        Alter Table: {tableName}
                    </DialogTitle>
                    <DialogDescription>
                        Modify table structure by adding, editing, or removing columns. Changes will be applied via migration.
                    </DialogDescription>
                </DialogHeader>

                <div className="flex-1 overflow-y-auto">
                    <div className="space-y-4">
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

                        {/* Columns Table */}
                        <div className="border rounded-lg">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead className="w-[200px]">Column Name</TableHead>
                                        <TableHead className="w-[150px]">Type</TableHead>
                                        <TableHead className="w-[100px]">Nullable</TableHead>
                                        <TableHead className="w-[150px]">Default Value</TableHead>
                                        <TableHead className="w-[100px]">Constraints</TableHead>
                                        <TableHead className="w-[100px]">Actions</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {columns.map((column, index) => (
                                        <TableRow key={index} className={cn(
                                            column.is_new && "bg-blue-50",
                                            column.is_editing && "bg-yellow-50"
                                        )}>
                                            <TableCell>
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
                                                        {column.is_new && <Badge variant="secondary" className="text-xs">New</Badge>}
                                                    </div>
                                                )}
                                            </TableCell>
                                            <TableCell>
                                                {column.is_editing ? (
                                                    <Select
                                                        value={column.type}
                                                        onValueChange={(value) => updateColumn(index, 'type', value)}
                                                    >
                                                        <SelectTrigger className="h-8">
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                        <SelectContent>
                                                            {COLUMN_TYPES.map(type => (
                                                                <SelectItem key={type} value={type}>
                                                                    {type}
                                                                </SelectItem>
                                                            ))}
                                                        </SelectContent>
                                                    </Select>
                                                ) : (
                                                    <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                                                        {column.type}
                                                    </code>
                                                )}
                                            </TableCell>
                                            <TableCell>
                                                {column.is_editing ? (
                                                    <Select
                                                        value={column.nullable ? 'true' : 'false'}
                                                        onValueChange={(value) => updateColumn(index, 'nullable', value === 'true')}
                                                    >
                                                        <SelectTrigger className="h-8">
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                        <SelectContent>
                                                            <SelectItem value="true">Yes</SelectItem>
                                                            <SelectItem value="false">No</SelectItem>
                                                        </SelectContent>
                                                    </Select>
                                                ) : (
                                                    <Badge variant={column.nullable ? 'secondary' : 'outline'}>
                                                        {column.nullable ? 'Yes' : 'No'}
                                                    </Badge>
                                                )}
                                            </TableCell>
                                            <TableCell>
                                                {column.is_editing ? (
                                                    <Input
                                                        value={column.default_value || ''}
                                                        onChange={(e) => updateColumn(index, 'default_value', e.target.value || null)}
                                                        placeholder="Default value"
                                                        className="h-8"
                                                    />
                                                ) : (
                                                    <span className="text-sm text-gray-500">
                                                        {column.default_value || 'NULL'}
                                                    </span>
                                                )}
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex gap-1">
                                                    {column.is_primary_key && (
                                                        <Badge variant="default" className="text-xs">PK</Badge>
                                                    )}
                                                    {column.is_foreign_key && (
                                                        <Badge variant="outline" className="text-xs">FK</Badge>
                                                    )}
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex gap-1">
                                                    {column.is_editing ? (
                                                        <>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                onClick={() => saveColumn(index)}
                                                                className="h-6 w-6 p-0"
                                                            >
                                                                <Check className="h-3 w-3" />
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                variant="ghost"
                                                                onClick={() => cancelEditing(index)}
                                                                className="h-6 w-6 p-0"
                                                            >
                                                                <X className="h-3 w-3" />
                                                            </Button>
                                                        </>
                                                    ) : (
                                                        <>
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
                                                        </>
                                                    )}
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>

                        {/* Add Column Button */}
                        <Button
                            onClick={addNewColumn}
                            variant="outline"
                            className="w-full"
                        >
                            <Plus className="h-4 w-4 mr-2" />
                            Add New Column
                        </Button>
                    </div>
                </div>

                <DialogFooter className="flex justify-between">
                    <div className="flex gap-2">
                        {hasChanges && (
                            <Button
                                onClick={resetChanges}
                                variant="outline"
                                disabled={migrationStatus === 'running'}
                            >
                                <ArrowLeft className="h-4 w-4 mr-2" />
                                Reset Changes
                            </Button>
                        )}
                        {migrationStatus === 'completed' && migrationId && (
                            <Button
                                onClick={handleRollback}
                                variant="outline"
                                disabled={rollbackMutation.isPending}
                            >
                                Rollback Migration
                            </Button>
                        )}
                    </div>
                    <div className="flex gap-2">
                        <Button
                            onClick={() => onOpenChange(false)}
                            variant="outline"
                        >
                            Cancel
                        </Button>
                        <Button
                            onClick={handleSave}
                            disabled={!hasChanges || migrationStatus === 'running'}
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
                    </div>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

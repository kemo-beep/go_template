'use client';

import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetFooter } from '@/components/ui/sheet';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Plus, Trash2, Loader2, Database, Key, Shield, Save, X } from 'lucide-react';
import { toast } from 'sonner';

interface Column {
    name: string;
    type: string;
    length?: number;
    not_null: boolean;
    primary_key: boolean;
    unique: boolean;
    default_value?: string;
}

interface CreateTableSheetProps {
    isOpen: boolean;
    onClose: () => void;
}

const DATA_TYPES = [
    'VARCHAR', 'INTEGER', 'BIGINT', 'SMALLINT', 'BOOLEAN',
    'TEXT', 'TIMESTAMP', 'DATE', 'TIME', 'DECIMAL',
    'NUMERIC', 'REAL', 'DOUBLE PRECISION', 'SERIAL', 'BIGSERIAL',
    'JSON', 'JSONB', 'UUID', 'BYTEA'
];

export function CreateTableSheet({ isOpen, onClose }: CreateTableSheetProps) {
    const [tableName, setTableName] = useState('');
    const [columns, setColumns] = useState<Column[]>([
        { name: 'id', type: 'SERIAL', not_null: true, primary_key: true, unique: false }
    ]);

    const queryClient = useQueryClient();

    const createTableMutation = useMutation({
        mutationFn: async (data: { table_name: string; columns: Column[] }) => {
            const response = await fetch('http://localhost:8080/api/v1/admin/database/tables', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
                },
                body: JSON.stringify(data),
            });
            if (!response.ok) throw new Error('Failed to create table');
            return response.json();
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['database-tables'] });
            toast.success('Table created successfully');
            onClose();
            setTableName('');
            setColumns([{ name: 'id', type: 'SERIAL', not_null: true, primary_key: true, unique: false }]);
        },
        onError: (error: any) => {
            toast.error(error.message || 'Failed to create table');
        },
    });

    const addColumn = () => {
        setColumns([...columns, {
            name: '',
            type: 'VARCHAR',
            not_null: false,
            primary_key: false,
            unique: false,
        }]);
    };

    const removeColumn = (index: number) => {
        if (columns.length > 1) {
            setColumns(columns.filter((_, i) => i !== index));
        }
    };

    const updateColumn = (index: number, field: keyof Column, value: any) => {
        const newColumns = [...columns];
        (newColumns[index] as any)[field] = value;
        setColumns(newColumns);
    };

    const handleCreate = () => {
        if (!tableName.trim()) {
            toast.error('Table name is required');
            return;
        }

        if (columns.length === 0) {
            toast.error('At least one column is required');
            return;
        }

        const invalidColumn = columns.find(col => !col.name.trim());
        if (invalidColumn) {
            toast.error('All columns must have a name');
            return;
        }

        // Check for duplicate column names
        const columnNames = columns.map(col => col.name.toLowerCase());
        const duplicateNames = columnNames.filter((name, index) => columnNames.indexOf(name) !== index);
        if (duplicateNames.length > 0) {
            toast.error('Column names must be unique');
            return;
        }

        // Check for multiple primary keys
        const primaryKeyCount = columns.filter(col => col.primary_key).length;
        if (primaryKeyCount > 1) {
            toast.error('Only one primary key is allowed');
            return;
        }

        createTableMutation.mutate({
            table_name: tableName,
            columns: columns,
        });
    };

    const handleCancel = () => {
        onClose();
    };

    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent className="w-[700px] sm:max-w-[700px] overflow-y-auto px-6">
                <SheetHeader>
                    <SheetTitle className="flex items-center gap-2">
                        <Database className="h-5 w-5" />
                        Create New Table
                    </SheetTitle>
                    <SheetDescription>
                        Define your table structure with columns and constraints
                    </SheetDescription>
                </SheetHeader>

                <div className="space-y-4 max-h-[calc(100vh-200px)] overflow-y-auto">
                    {/* Table Name */}
                    <div className="space-y-2">
                        <h3 className="flex items-center gap-2">
                            <Database className="h-4 w-4" />
                            Table Information
                        </h3>
                        <div className="space-y-2 px-2">
                            <Label htmlFor="table-name">Table Name *</Label>
                            <Input
                                id="table-name"
                                placeholder="e.g., users, products, orders"
                                value={tableName}
                                onChange={(e) => setTableName(e.target.value)}
                            />
                        </div>
                    </div>

                    <Separator />

                    {/* Columns */}
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <h3 className="flex items-center gap-2">
                                <Key className="h-4 w-4" />
                                Columns ({columns.length})
                            </h3>
                            <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                onClick={addColumn}
                            >
                                <Plus className="h-4 w-4 mr-2" />
                                Insert Column
                            </Button>
                        </div>

                        <div className="border rounded-lg overflow-hidden">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead className="w-[120px]">Name</TableHead>
                                        <TableHead className="w-[120px]">Type</TableHead>
                                        <TableHead className="w-[80px]">Length</TableHead>
                                        <TableHead className="w-[120px]">Default</TableHead>
                                        <TableHead className="w-[80px]">NOT NULL</TableHead>
                                        <TableHead className="w-[80px]">Primary Key</TableHead>
                                        <TableHead className="w-[80px]">Unique</TableHead>
                                        <TableHead className="w-[60px]">Actions</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {columns.map((column, index) => (
                                        <TableRow key={index} className="hover:bg-muted/50">
                                            <TableCell>
                                                <Input
                                                    placeholder="column_name"
                                                    value={column.name}
                                                    onChange={(e) => updateColumn(index, 'name', e.target.value)}
                                                    className="h-8"
                                                />
                                            </TableCell>
                                            <TableCell>
                                                <Select
                                                    value={column.type}
                                                    onValueChange={(value) => updateColumn(index, 'type', value)}
                                                >
                                                    <SelectTrigger className="h-8">
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {DATA_TYPES.map((type) => (
                                                            <SelectItem key={type} value={type}>
                                                                {type}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            </TableCell>
                                            <TableCell>
                                                {(column.type === 'VARCHAR' || column.type === 'CHAR') ? (
                                                    <Input
                                                        type="number"
                                                        placeholder="255"
                                                        value={column.length || ''}
                                                        onChange={(e) => updateColumn(index, 'length', parseInt(e.target.value))}
                                                        className="h-8"
                                                    />
                                                ) : (
                                                    <span className="text-xs text-muted-foreground">-</span>
                                                )}
                                            </TableCell>
                                            <TableCell>
                                                <Input
                                                    placeholder="default"
                                                    value={column.default_value || ''}
                                                    onChange={(e) => updateColumn(index, 'default_value', e.target.value || undefined)}
                                                    className="h-8"
                                                />
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex items-center justify-center">
                                                    <Switch
                                                        id={`not-null-${index}`}
                                                        checked={column.not_null}
                                                        onCheckedChange={(checked) => updateColumn(index, 'not_null', checked)}
                                                    />
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex items-center justify-center">
                                                    <Switch
                                                        id={`primary-${index}`}
                                                        checked={column.primary_key}
                                                        onCheckedChange={(checked) => updateColumn(index, 'primary_key', checked)}
                                                    />
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex items-center justify-center">
                                                    <Switch
                                                        id={`unique-${index}`}
                                                        checked={column.unique}
                                                        onCheckedChange={(checked) => updateColumn(index, 'unique', checked)}
                                                    />
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                {columns.length > 1 && (
                                                    <Button
                                                        type="button"
                                                        variant="ghost"
                                                        size="sm"
                                                        onClick={() => removeColumn(index)}
                                                        className="h-6 w-6 p-0 text-red-500 hover:text-red-700"
                                                    >
                                                        <Trash2 className="h-3 w-3" />
                                                    </Button>
                                                )}
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    </div>
                </div>

                <SheetFooter className="flex gap-2">
                    <div className="flex items-center gap-2">
                        <Button variant="outline" onClick={handleCancel}>
                            <X className="h-4 w-4 mr-2" />
                            Cancel
                        </Button>
                        <Button onClick={handleCreate} disabled={createTableMutation.isPending}>
                            {createTableMutation.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating...
                                </>
                            ) : (
                                <>
                                    <Save className="h-4 w-4 mr-2" />
                                    Create Table
                                </>
                            )}
                        </Button>
                    </div>
                </SheetFooter>
            </SheetContent>
        </Sheet>
    );
}

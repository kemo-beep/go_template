'use client';

import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Plus, Trash2, Loader2 } from 'lucide-react';
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

interface CreateTableDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const DATA_TYPES = [
    'VARCHAR', 'INTEGER', 'BIGINT', 'SMALLINT', 'BOOLEAN',
    'TEXT', 'TIMESTAMP', 'DATE', 'TIME', 'DECIMAL',
    'NUMERIC', 'REAL', 'DOUBLE PRECISION', 'SERIAL', 'BIGSERIAL',
    'JSON', 'JSONB', 'UUID', 'BYTEA'
];

export function CreateTableDialog({ open, onOpenChange }: CreateTableDialogProps) {
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
            onOpenChange(false);
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
        setColumns(columns.filter((_, i) => i !== index));
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

        createTableMutation.mutate({
            table_name: tableName,
            columns: columns,
        });
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Create New Table</DialogTitle>
                    <DialogDescription>
                        Define your table structure with columns and constraints
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-6 py-4">
                    {/* Table Name */}
                    <div className="space-y-2">
                        <Label htmlFor="table-name">Table Name</Label>
                        <Input
                            id="table-name"
                            placeholder="e.g., users, products, orders"
                            value={tableName}
                            onChange={(e) => setTableName(e.target.value)}
                        />
                    </div>

                    {/* Columns */}
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <Label>Columns</Label>
                            <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                onClick={addColumn}
                            >
                                <Plus className="h-4 w-4 mr-2" />
                                Add Column
                            </Button>
                        </div>

                        <div className="space-y-4 border rounded-lg p-4">
                            {columns.map((column, index) => (
                                <div key={index} className="space-y-3 pb-4 border-b last:border-b-0 last:pb-0">
                                    <div className="flex items-center justify-between">
                                        <span className="text-sm font-medium">Column {index + 1}</span>
                                        {index > 0 && (
                                            <Button
                                                type="button"
                                                variant="ghost"
                                                size="sm"
                                                onClick={() => removeColumn(index)}
                                            >
                                                <Trash2 className="h-4 w-4 text-red-500" />
                                            </Button>
                                        )}
                                    </div>

                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="space-y-2">
                                            <Label>Name</Label>
                                            <Input
                                                placeholder="column_name"
                                                value={column.name}
                                                onChange={(e) => updateColumn(index, 'name', e.target.value)}
                                            />
                                        </div>

                                        <div className="space-y-2">
                                            <Label>Type</Label>
                                            <Select
                                                value={column.type}
                                                onValueChange={(value) => updateColumn(index, 'type', value)}
                                            >
                                                <SelectTrigger>
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
                                        </div>

                                        {(column.type === 'VARCHAR' || column.type === 'CHAR') && (
                                            <div className="space-y-2">
                                                <Label>Length</Label>
                                                <Input
                                                    type="number"
                                                    placeholder="255"
                                                    value={column.length || ''}
                                                    onChange={(e) => updateColumn(index, 'length', parseInt(e.target.value))}
                                                />
                                            </div>
                                        )}

                                        <div className="space-y-2">
                                            <Label>Default Value (optional)</Label>
                                            <Input
                                                placeholder="e.g., NOW(), 0, 'default'"
                                                value={column.default_value || ''}
                                                onChange={(e) => updateColumn(index, 'default_value', e.target.value || undefined)}
                                            />
                                        </div>
                                    </div>

                                    <div className="flex gap-6">
                                        <div className="flex items-center space-x-2">
                                            <Switch
                                                id={`not-null-${index}`}
                                                checked={column.not_null}
                                                onCheckedChange={(checked) => updateColumn(index, 'not_null', checked)}
                                            />
                                            <Label htmlFor={`not-null-${index}`} className="text-sm">
                                                NOT NULL
                                            </Label>
                                        </div>

                                        <div className="flex items-center space-x-2">
                                            <Switch
                                                id={`primary-${index}`}
                                                checked={column.primary_key}
                                                onCheckedChange={(checked) => updateColumn(index, 'primary_key', checked)}
                                            />
                                            <Label htmlFor={`primary-${index}`} className="text-sm">
                                                Primary Key
                                            </Label>
                                        </div>

                                        <div className="flex items-center space-x-2">
                                            <Switch
                                                id={`unique-${index}`}
                                                checked={column.unique}
                                                onCheckedChange={(checked) => updateColumn(index, 'unique', checked)}
                                            />
                                            <Label htmlFor={`unique-${index}`} className="text-sm">
                                                Unique
                                            </Label>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Cancel
                    </Button>
                    <Button onClick={handleCreate} disabled={createTableMutation.isPending}>
                        {createTableMutation.isPending ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                Creating...
                            </>
                        ) : (
                            'Create Table'
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

'use client';

import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Textarea } from '@/components/ui/textarea';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetFooter } from '@/components/ui/sheet';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { toast } from 'sonner';
import { Save, X, Key, Link, Shield, Database, Loader2 } from 'lucide-react';
import { api } from '@/lib/api-client';

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

interface ColumnEditorSheetProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (column: ColumnInfo) => void;
    column?: ColumnInfo | null;
    isEditing: boolean;
    existingColumns: ColumnInfo[];
}

export function ColumnEditorSheet({
    isOpen,
    onClose,
    onSave,
    column,
    isEditing,
    existingColumns
}: ColumnEditorSheetProps) {
    const [formData, setFormData] = useState<ColumnInfo>({
        name: '',
        type: 'VARCHAR',
        nullable: true,
        default_value: null,
        is_primary_key: false,
        is_foreign_key: false,
        references: '',
        references_table: '',
        references_column: '',
        on_delete: 'NO ACTION',
        on_update: 'NO ACTION',
        is_editing: false,
        is_new: false,
        constraints: {
            unique: false,
            check: '',
            not_null: false,
        },
        indexes: {
            name: '',
            type: 'BTREE',
            unique: false,
        },
    });

    // Fetch all tables for foreign key selection
    const { data: tablesData, isLoading: tablesLoading } = useQuery({
        queryKey: ['all-tables'],
        queryFn: () => api.getAllTables(),
        enabled: formData.is_foreign_key,
    });

    // Fetch columns for selected table using the existing getTableSchema API
    const { data: columnsData, isLoading: columnsLoading } = useQuery({
        queryKey: ['table-schema', formData.references_table],
        queryFn: () => api.getTableSchema(formData.references_table!),
        enabled: formData.is_foreign_key && !!formData.references_table,
    });

    // Debug logging
    useEffect(() => {
        if (tablesData) {
            console.log('Tables data:', tablesData);
            console.log('Tables array:', tablesData?.data?.data);
        }
    }, [tablesData]);

    useEffect(() => {
        if (columnsData) {
            console.log('Columns data:', columnsData);
            console.log('Columns array:', columnsData?.data?.data?.columns);
        }
    }, [columnsData]);

    // Initialize form data when column changes
    useEffect(() => {
        if (column) {
            setFormData({
                ...column,
                constraints: column.constraints || {
                    unique: false,
                    check: '',
                    not_null: false,
                },
                indexes: column.indexes || {
                    name: '',
                    type: 'BTREE',
                    unique: false,
                },
            });
        } else {
            // Reset form for new column
            setFormData({
                name: '',
                type: 'VARCHAR',
                nullable: true,
                default_value: null,
                is_primary_key: false,
                is_foreign_key: false,
                references: '',
                references_table: '',
                references_column: '',
                on_delete: 'NO ACTION',
                on_update: 'NO ACTION',
                is_editing: false,
                is_new: true,
                constraints: {
                    unique: false,
                    check: '',
                    not_null: false,
                },
                indexes: {
                    name: '',
                    type: 'BTREE',
                    unique: false,
                },
            });
        }
    }, [column]);

    const updateField = (field: keyof ColumnInfo, value: any) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const updateConstraints = (field: keyof NonNullable<ColumnInfo['constraints']>, value: any) => {
        setFormData(prev => ({
            ...prev,
            constraints: {
                ...prev.constraints,
                [field]: value
            }
        }));
    };

    const updateIndexes = (field: keyof NonNullable<ColumnInfo['indexes']>, value: any) => {
        setFormData(prev => ({
            ...prev,
            indexes: {
                ...prev.indexes,
                [field]: value
            }
        }));
    };

    const handleTableSelection = (tableName: string) => {
        setFormData(prev => ({
            ...prev,
            references_table: tableName,
            references_column: '', // Reset column selection when table changes
        }));
    };

    const handleSave = () => {
        // Validation
        if (!formData.name.trim()) {
            toast.error('Column name is required');
            return;
        }

        // Check for duplicate names (excluding current column if editing)
        const duplicateIndex = existingColumns.findIndex((col, i) =>
            col.name.toLowerCase() === formData.name.toLowerCase() &&
            (!isEditing || col.original_name !== formData.original_name)
        );
        if (duplicateIndex !== -1) {
            toast.error('Column name must be unique');
            return;
        }

        // Check for multiple primary keys
        if (formData.is_primary_key) {
            const existingPrimaryKey = existingColumns.find(col =>
                col.is_primary_key && (!isEditing || col.original_name !== formData.original_name)
            );
            if (existingPrimaryKey) {
                toast.error('Only one primary key is allowed per table');
                return;
            }
        }

        // Validate foreign key references
        if (formData.is_foreign_key) {
            if (!formData.references_table || !formData.references_column) {
                toast.error('Foreign key must specify both table and column references');
                return;
            }

            // Check if the referenced table and column exist
            if (tablesLoading || columnsLoading) {
                toast.error('Please wait for table and column data to load');
                return;
            }

            if (!tablesData?.data?.data?.find((table: any) => table.name === formData.references_table)) {
                toast.error('Selected table does not exist');
                return;
            }

            if (!columnsData?.data?.data?.columns?.find((col: any) => col.name === formData.references_column)) {
                toast.error('Selected column does not exist in the referenced table');
                return;
            }
        }

        // Validate check constraint
        if (formData.constraints?.check && formData.constraints.check.trim()) {
            const checkExpr = formData.constraints.check.trim();
            if (!checkExpr.match(/^[a-zA-Z_][a-zA-Z0-9_]*\s*[<>=!]+.*$/)) {
                toast.error('Check constraint must be a valid SQL expression (e.g., "age &gt; 0")');
                return;
            }
        }

        onSave(formData);
        onClose();
    };

    const handleCancel = () => {
        onClose();
    };

    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent className="w-[600px] sm:max-w-[600px] overflow-y-auto px-6">
                <SheetHeader>
                    <SheetTitle className="flex items-center gap-2">
                        <Database className="h-5 w-5" />
                        {isEditing ? 'Edit Column' : 'Add New Column'}
                    </SheetTitle>
                    <SheetDescription>
                        {isEditing
                            ? 'Modify the column properties and constraints'
                            : 'Define a new column with its properties, constraints, and indexes'
                        }
                    </SheetDescription>
                </SheetHeader>

                <div className="space-y-4  max-h-[calc(100vh-200px)] overflow-y-auto">
                    {/* Basic Properties */}
                    <div className="space-y-2">
                        <h3 className="flex items-center gap-2">
                            <Database className="h-4 w-4" />
                            Properties
                        </h3>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="column-name">Column Name *</Label>
                                <Input
                                    id="column-name"
                                    value={formData.name}
                                    onChange={(e) => updateField('name', e.target.value)}
                                    placeholder="Enter column name"
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="column-type">Data Type *</Label>
                                <Select
                                    value={formData.type}
                                    onValueChange={(value) => updateField('type', value)}
                                >
                                    <SelectTrigger>
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
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="default-value">Default Value</Label>
                                <Input
                                    id="default-value"
                                    value={formData.default_value || ''}
                                    onChange={(e) => updateField('default_value', e.target.value || null)}
                                    placeholder="Enter default value"
                                />
                            </div>

                            <div className="space-y-2">
                                <Label>Nullable</Label>
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="nullable"
                                        checked={formData.nullable}
                                        onCheckedChange={(checked) => updateField('nullable', checked)}
                                    />
                                    <Label htmlFor="nullable" className="text-sm">
                                        Allow NULL values
                                    </Label>
                                </div>
                            </div>
                        </div>
                    </div>

                    <Separator />

                    {/* Keys */}
                    <div className="py-2  flex items-center gap-2 justify-between">
                        <h3 className=" font-semibold flex items-center gap-2">
                            <Key className="h-4 w-4" />
                            Keys
                        </h3>

                        <div className="flex items-center gap-2">
                            {/* Primary Key */}
                            <div className="flex items-center space-x-2">
                                <Checkbox
                                    id="primary-key"
                                    checked={formData.is_primary_key}
                                    onCheckedChange={(checked) => updateField('is_primary_key', checked)}
                                />
                                <Label htmlFor="primary-key" className="text-sm font-medium">
                                    Primary Key
                                </Label>
                                {formData.is_primary_key && (
                                    <Badge variant="default" className="text-xs">PK</Badge>
                                )}
                            </div>

                            {/* Foreign Key */}
                            <div className="space-y-3">
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="foreign-key"
                                        checked={formData.is_foreign_key || false}
                                        onCheckedChange={(checked) => updateField('is_foreign_key', checked)}
                                    />
                                    <Label htmlFor="foreign-key" className="text-sm font-medium">
                                        Foreign Key
                                    </Label>
                                    {formData.is_foreign_key && (
                                        <Badge variant="outline" className="text-xs">FK</Badge>
                                    )}
                                </div>

                                {formData.is_foreign_key && (
                                    <div className="ml-6 space-y-3 p-3 bg-muted/50 rounded-lg">
                                        <div className="grid grid-cols-2 gap-3">
                                            <div className="space-y-2">
                                                <Label htmlFor="ref-table">Referenced Table *</Label>
                                                <Select
                                                    value={formData.references_table || ''}
                                                    onValueChange={handleTableSelection}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue placeholder="Select a table" />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {tablesLoading ? (
                                                            <div className="flex items-center justify-center p-2">
                                                                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                                                                <span className="text-sm text-muted-foreground">Loading tables...</span>
                                                            </div>
                                                        ) : tablesData?.data?.data && tablesData.data.data.length > 0 ? (
                                                            tablesData.data.data.map((table: any) => (
                                                                <SelectItem key={table.name} value={table.name}>
                                                                    <div className="flex items-center gap-2">
                                                                        <Database className="h-3 w-3" />
                                                                        {table.name}
                                                                    </div>
                                                                </SelectItem>
                                                            ))
                                                        ) : (
                                                            <div className="p-2 text-sm text-muted-foreground">
                                                                No tables found
                                                            </div>
                                                        )}
                                                    </SelectContent>
                                                </Select>
                                            </div>
                                            <div className="space-y-2">
                                                <Label htmlFor="ref-column">Referenced Column *</Label>
                                                <Select
                                                    value={formData.references_column || ''}
                                                    onValueChange={(value) => updateField('references_column', value)}
                                                    disabled={!formData.references_table}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue placeholder={
                                                            !formData.references_table
                                                                ? "Select a table first"
                                                                : "Select a column"
                                                        } />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {columnsLoading ? (
                                                            <div className="flex items-center justify-center p-2">
                                                                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                                                                <span className="text-sm text-muted-foreground">Loading columns...</span>
                                                            </div>
                                                        ) : columnsData?.data?.data?.columns && columnsData.data.data.columns.length > 0 ? (
                                                            columnsData.data.data.columns.map((col: any) => (
                                                                <SelectItem key={col.name} value={col.name}>
                                                                    <div className="flex items-center gap-2">
                                                                        <Key className="h-3 w-3" />
                                                                        <span>{col.name}</span>
                                                                        <span className="text-xs text-muted-foreground">({col.type})</span>
                                                                    </div>
                                                                </SelectItem>
                                                            ))
                                                        ) : formData.references_table ? (
                                                            <div className="p-2 text-sm text-muted-foreground">
                                                                No columns found in {formData.references_table}
                                                                {columnsData && (
                                                                    <div className="mt-1 text-xs">
                                                                        Debug: {JSON.stringify(columnsData, null, 2)}
                                                                    </div>
                                                                )}
                                                            </div>
                                                        ) : (
                                                            <div className="p-2 text-sm text-muted-foreground">
                                                                Select a table first
                                                            </div>
                                                        )}
                                                    </SelectContent>
                                                </Select>
                                            </div>
                                        </div>

                                        <div className="grid grid-cols-2 gap-3">
                                            <div className="space-y-2">
                                                <Label htmlFor="on-delete">ON DELETE</Label>
                                                <Select
                                                    value={formData.on_delete || 'NO ACTION'}
                                                    onValueChange={(value) => updateField('on_delete', value)}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value="CASCADE">CASCADE</SelectItem>
                                                        <SelectItem value="SET NULL">SET NULL</SelectItem>
                                                        <SelectItem value="RESTRICT">RESTRICT</SelectItem>
                                                        <SelectItem value="NO ACTION">NO ACTION</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </div>
                                            <div className="space-y-2">
                                                <Label htmlFor="on-update">ON UPDATE</Label>
                                                <Select
                                                    value={formData.on_update || 'NO ACTION'}
                                                    onValueChange={(value) => updateField('on_update', value)}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value="CASCADE">CASCADE</SelectItem>
                                                        <SelectItem value="SET NULL">SET NULL</SelectItem>
                                                        <SelectItem value="RESTRICT">RESTRICT</SelectItem>
                                                        <SelectItem value="NO ACTION">NO ACTION</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </div>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>

                    <Separator />

                    {/* Constraints */}
                    <div className="space-y-2">
                        <div className="flex items-center gap-2 justify-between">
                            <h3 className="font-semibold flex items-center gap-2">
                                <Shield className="h-4 w-4" />
                                Constraints
                            </h3>

                            <div className="flex items-center gap-2">
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="unique-constraint"
                                        checked={formData.constraints?.unique || false}
                                        onCheckedChange={(checked) => updateConstraints('unique', checked)}
                                    />
                                    <Label htmlFor="unique-constraint" className="text-sm font-medium">
                                        Unique
                                    </Label>
                                    {formData.constraints?.unique && (
                                        <Badge variant="secondary" className="text-xs">UNIQUE</Badge>
                                    )}
                                </div>

                                <div className="">
                                    {/* <Label htmlFor="check-constraint" className="text-xs">Check Constraint</Label> */}
                                    <Input
                                        className="max-w-40"
                                        id="check-constraint"
                                        value={formData.constraints?.check || ''}
                                        onChange={(e) => updateConstraints('check', e.target.value)}
                                        placeholder="e.g., age &gt; 0"
                                    />

                                </div>
                            </div>

                        </div>
                        <p className="text-xs text-muted-foreground">
                            Enter a SQL expression for validation (e.g., &quot;age &gt; 0&quot;, &quot;status IN ('active', 'inactive')&quot;)
                        </p>

                    </div>

                    <Separator />

                    {/* Indexes */}
                    <div className="space-y-2">
                        <h3 className=" flex items-center gap-2">
                            <Database className="h-4 w-4" />
                            Indexes
                        </h3>

                        <div className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="index-name">Index Name</Label>
                                <Input
                                    id="index-name"
                                    value={formData.indexes?.name || ''}
                                    onChange={(e) => updateIndexes('name', e.target.value)}
                                    placeholder="Enter index name (optional)"
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-2">
                                    <Label htmlFor="index-type">Index Type</Label>
                                    <Select
                                        value={formData.indexes?.type || 'BTREE'}
                                        onValueChange={(value) => updateIndexes('type', value)}
                                    >
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="BTREE">BTREE</SelectItem>
                                            <SelectItem value="HASH">HASH</SelectItem>
                                            <SelectItem value="GIN">GIN</SelectItem>
                                            <SelectItem value="GIST">GIST</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>

                                <div className="space-y-2">
                                    <Label>Unique Index</Label>
                                    <div className="flex items-center space-x-2">
                                        <Checkbox
                                            id="unique-index"
                                            checked={formData.indexes?.unique || false}
                                            onCheckedChange={(checked) => updateIndexes('unique', checked)}
                                        />
                                        <Label htmlFor="unique-index" className="text-sm">
                                            Create unique index
                                        </Label>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <SheetFooter className="flex gap-2">
                    <div className="flex items-center gap-2">
                        <Button variant="outline" onClick={handleCancel}>
                            <X className="h-4 w-4 mr-2" />
                            Cancel
                        </Button>
                        <Button onClick={handleSave}>
                            <Save className="h-4 w-4 mr-2" />
                            {isEditing ? 'Update Column' : 'Insert Column'}
                        </Button>
                    </div>
                </SheetFooter>
            </SheetContent>
        </Sheet>
    );
}

'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog';
import {
    History,
    Clock,
    CheckCircle,
    XCircle,
    Loader2,
    ArrowLeft,
    Eye,
    AlertTriangle,
} from 'lucide-react';

interface Migration {
    id: string;
    table_name: string;
    sql_query: string;
    status: 'pending' | 'running' | 'completed' | 'failed' | 'rolled_back';
    created_at: string;
    completed_at?: string;
    error_message?: string;
    rollback_sql?: string;
}

export function MigrationHistory() {

    const { data: migrations, isLoading } = useQuery({
        queryKey: ['migrations'],
        queryFn: () => api.getMigrations({ limit: 50, offset: 0 }),
        refetchInterval: 5000, // Refresh every 5 seconds
    });

    const migrationList = migrations?.data?.migrations || [];

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'pending':
                return <Clock className="h-4 w-4 text-yellow-500" />;
            case 'running':
                return <Loader2 className="h-4 w-4 animate-spin text-blue-500" />;
            case 'completed':
                return <CheckCircle className="h-4 w-4 text-green-500" />;
            case 'failed':
                return <XCircle className="h-4 w-4 text-red-500" />;
            case 'rolled_back':
                return <ArrowLeft className="h-4 w-4 text-gray-500" />;
            default:
                return <AlertTriangle className="h-4 w-4 text-gray-500" />;
        }
    };

    const getStatusBadge = (status: string) => {
        const variants = {
            pending: 'secondary',
            running: 'default',
            completed: 'default',
            failed: 'destructive',
            rolled_back: 'outline',
        } as const;

        return (
            <Badge variant={variants[status as keyof typeof variants] || 'secondary'}>
                {status.replace('_', ' ').toUpperCase()}
            </Badge>
        );
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString();
    };

    const handleRollback = async (migration: Migration) => {
        if (!migration.rollback_sql) {
            alert('No rollback SQL available for this migration');
            return;
        }

        if (window.confirm('Are you sure you want to rollback this migration?')) {
            try {
                await api.rollbackMigration(migration.id);
                alert('Migration rolled back successfully');
            } catch (error: unknown) {
                const errorMessage = error instanceof Error ? error.message : 'Unknown error';
                alert(`Rollback failed: ${errorMessage}`);
            }
        }
    };

    return (
        <div className="space-y-4">
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <History className="h-5 w-5" />
                        Migration History
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <div className="flex items-center justify-center py-8">
                            <Loader2 className="h-6 w-6 animate-spin" />
                            <span className="ml-2">Loading migrations...</span>
                        </div>
                    ) : migrationList.length === 0 ? (
                        <div className="text-center py-8 text-gray-500">
                            <History className="h-12 w-12 mx-auto mb-4 opacity-50" />
                            <p>No migrations found</p>
                        </div>
                    ) : (
                        <div className="space-y-4">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>Status</TableHead>
                                        <TableHead>Table</TableHead>
                                        <TableHead>Created</TableHead>
                                        <TableHead>Completed</TableHead>
                                        <TableHead>Actions</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {migrationList.map((migration: Migration) => (
                                        <TableRow key={migration.id}>
                                            <TableCell>
                                                <div className="flex items-center gap-2">
                                                    {getStatusIcon(migration.status)}
                                                    {getStatusBadge(migration.status)}
                                                </div>
                                            </TableCell>
                                            <TableCell className="font-medium">
                                                {migration.table_name}
                                            </TableCell>
                                            <TableCell>
                                                {formatDate(migration.created_at)}
                                            </TableCell>
                                            <TableCell>
                                                {migration.completed_at
                                                    ? formatDate(migration.completed_at)
                                                    : '-'
                                                }
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex gap-2">
                                                    <Dialog>
                                                        <DialogTrigger asChild>
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                            >
                                                                <Eye className="h-4 w-4" />
                                                            </Button>
                                                        </DialogTrigger>
                                                        <DialogContent className="max-w-2xl">
                                                            <DialogHeader>
                                                                <DialogTitle>
                                                                    Migration Details
                                                                </DialogTitle>
                                                                <DialogDescription>
                                                                    Migration ID: {migration.id}
                                                                </DialogDescription>
                                                            </DialogHeader>
                                                            <div className="space-y-4">
                                                                <div>
                                                                    <h4 className="font-medium mb-2">SQL Query</h4>
                                                                    <pre className="bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                                                                        {migration.sql_query}
                                                                    </pre>
                                                                </div>
                                                                {migration.error_message && (
                                                                    <div>
                                                                        <h4 className="font-medium mb-2 text-red-600">Error Message</h4>
                                                                        <pre className="bg-red-50 p-3 rounded text-sm text-red-700">
                                                                            {migration.error_message}
                                                                        </pre>
                                                                    </div>
                                                                )}
                                                                {migration.rollback_sql && (
                                                                    <div>
                                                                        <h4 className="font-medium mb-2">Rollback SQL</h4>
                                                                        <pre className="bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                                                                            {migration.rollback_sql}
                                                                        </pre>
                                                                    </div>
                                                                )}
                                                            </div>
                                                        </DialogContent>
                                                    </Dialog>

                                                    {migration.status === 'completed' && migration.rollback_sql && (
                                                        <Button
                                                            variant="outline"
                                                            size="sm"
                                                            onClick={() => handleRollback(migration)}
                                                        >
                                                            <ArrowLeft className="h-4 w-4 mr-1" />
                                                            Rollback
                                                        </Button>
                                                    )}
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}

'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog';
import {
    Database,
    Play,
    Terminal,
    Flag,
    Settings,
    ChevronUp,
    ChevronDown,
    Plus,
} from 'lucide-react';
import { api } from '@/lib/api-client';
import { toast } from 'sonner';

export default function DevToolsPage() {
    const queryClient = useQueryClient();
    const [selectedJob, setSelectedJob] = useState('');
    const [newFlagName, setNewFlagName] = useState('');
    const [newFlagValue, setNewFlagValue] = useState(false);

    const { data: migrations } = useQuery({
        queryKey: ['migrations'],
        queryFn: () => api.getMigrations(),
    });

    const { data: jobs } = useQuery({
        queryKey: ['background-jobs'],
        queryFn: () => api.getBackgroundJobs(),
    });

    const { data: featureFlags } = useQuery({
        queryKey: ['feature-flags'],
        queryFn: () => api.getFeatureFlags(),
    });

    const runMigrationMutation = useMutation({
        mutationFn: (direction: 'up' | 'down') => api.runMigration(direction),
        onSuccess: (data, direction) => {
            toast.success(`Migration ${direction} executed successfully`);
            queryClient.invalidateQueries({ queryKey: ['migrations'] });
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Migration failed');
        },
    });

    const runJobMutation = useMutation({
        mutationFn: (jobName: string) => api.runBackgroundJob(jobName),
        onSuccess: () => {
            toast.success('Job triggered successfully');
            queryClient.invalidateQueries({ queryKey: ['background-jobs'] });
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Job execution failed');
        },
    });

    const toggleFeatureFlagMutation = useMutation({
        mutationFn: ({ flag, enabled }: { flag: string; enabled: boolean }) =>
            api.toggleFeatureFlag(flag, enabled),
        onSuccess: () => {
            toast.success('Feature flag updated');
            queryClient.invalidateQueries({ queryKey: ['feature-flags'] });
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Update failed');
        },
    });

    const createFeatureFlagMutation = useMutation({
        mutationFn: ({ name, enabled }: { name: string; enabled: boolean }) =>
            api.createFeatureFlag(name, enabled),
        onSuccess: () => {
            toast.success('Feature flag created');
            queryClient.invalidateQueries({ queryKey: ['feature-flags'] });
            setNewFlagName('');
            setNewFlagValue(false);
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Creation failed');
        },
    });

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Developer Tools</h1>
                <p className="text-gray-500 mt-1">
                    Manage migrations, background jobs, and feature flags
                </p>
            </div>

            <Tabs defaultValue="migrations" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="migrations">
                        <Database className="h-4 w-4 mr-2" />
                        Migrations
                    </TabsTrigger>
                    <TabsTrigger value="jobs">
                        <Terminal className="h-4 w-4 mr-2" />
                        Background Jobs
                    </TabsTrigger>
                    <TabsTrigger value="flags">
                        <Flag className="h-4 w-4 mr-2" />
                        Feature Flags
                    </TabsTrigger>
                </TabsList>

                {/* Migrations Tab */}
                <TabsContent value="migrations" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <CardTitle>Database Migrations (Goose)</CardTitle>
                                <div className="flex gap-2">
                                    <Button
                                        onClick={() => runMigrationMutation.mutate('up')}
                                        disabled={runMigrationMutation.isPending}
                                    >
                                        <ChevronUp className="h-4 w-4 mr-2" />
                                        Migrate Up
                                    </Button>
                                    <Button
                                        variant="outline"
                                        onClick={() => runMigrationMutation.mutate('down')}
                                        disabled={runMigrationMutation.isPending}
                                    >
                                        <ChevronDown className="h-4 w-4 mr-2" />
                                        Migrate Down
                                    </Button>
                                </div>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>Version</TableHead>
                                        <TableHead>Description</TableHead>
                                        <TableHead>Applied At</TableHead>
                                        <TableHead>Status</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {migrations?.data?.map((migration: any) => (
                                        <TableRow key={migration.version}>
                                            <TableCell className="font-mono font-medium">
                                                {migration.version}
                                            </TableCell>
                                            <TableCell>{migration.description}</TableCell>
                                            <TableCell className="text-sm text-gray-500">
                                                {migration.applied_at
                                                    ? new Date(migration.applied_at).toLocaleString()
                                                    : '-'}
                                            </TableCell>
                                            <TableCell>
                                                <Badge variant={migration.is_applied ? 'default' : 'secondary'}>
                                                    {migration.is_applied ? 'Applied' : 'Pending'}
                                                </Badge>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle className="text-sm">Migration Commands</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-2">
                            <div className="font-mono text-sm bg-gray-100 p-3 rounded">
                                <p className="text-gray-600"># Run all pending migrations</p>
                                <p className="text-blue-600">make migrate-up</p>
                            </div>
                            <div className="font-mono text-sm bg-gray-100 p-3 rounded">
                                <p className="text-gray-600"># Rollback last migration</p>
                                <p className="text-blue-600">make migrate-down</p>
                            </div>
                            <div className="font-mono text-sm bg-gray-100 p-3 rounded">
                                <p className="text-gray-600"># Create new migration</p>
                                <p className="text-blue-600">make migrate-create name=add_users_table</p>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Background Jobs Tab */}
                <TabsContent value="jobs" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Available Background Jobs</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-3">
                                {jobs?.data?.map((job: any) => (
                                    <div
                                        key={job.name}
                                        className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50"
                                    >
                                        <div className="flex-1">
                                            <h4 className="font-medium">{job.name}</h4>
                                            <p className="text-sm text-gray-500 mt-1">{job.description}</p>
                                            <div className="flex items-center gap-4 mt-2">
                                                <span className="text-xs text-gray-500">
                                                    Last run: {job.last_run ? new Date(job.last_run).toLocaleString() : 'Never'}
                                                </span>
                                                {job.schedule && (
                                                    <Badge variant="outline" className="text-xs">
                                                        {job.schedule}
                                                    </Badge>
                                                )}
                                            </div>
                                        </div>
                                        <Button
                                            onClick={() => runJobMutation.mutate(job.name)}
                                            disabled={runJobMutation.isPending}
                                        >
                                            <Play className="h-4 w-4 mr-2" />
                                            Run Now
                                        </Button>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>

                    <div className="grid md:grid-cols-3 gap-4">
                        <Card>
                            <CardHeader className="pb-3">
                                <CardTitle className="text-sm font-medium text-gray-500">
                                    Total Jobs
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{jobs?.data?.length || 0}</div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="pb-3">
                                <CardTitle className="text-sm font-medium text-gray-500">
                                    Scheduled
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">
                                    {jobs?.data?.filter((j: any) => j.schedule).length || 0}
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="pb-3">
                                <CardTitle className="text-sm font-medium text-gray-500">
                                    Running
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">
                                    {jobs?.data?.filter((j: any) => j.status === 'running').length || 0}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                {/* Feature Flags Tab */}
                <TabsContent value="flags" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <CardTitle>Feature Flags</CardTitle>
                                <Dialog>
                                    <DialogTrigger asChild>
                                        <Button size="sm">
                                            <Plus className="h-4 w-4 mr-2" />
                                            Add Flag
                                        </Button>
                                    </DialogTrigger>
                                    <DialogContent>
                                        <DialogHeader>
                                            <DialogTitle>Create Feature Flag</DialogTitle>
                                        </DialogHeader>
                                        <div className="space-y-4 py-4">
                                            <div>
                                                <label className="text-sm font-medium mb-2 block">Flag Name</label>
                                                <Input
                                                    value={newFlagName}
                                                    onChange={(e) => setNewFlagName(e.target.value)}
                                                    placeholder="new_feature_enabled"
                                                />
                                            </div>
                                            <div className="flex items-center justify-between">
                                                <label className="text-sm font-medium">Enabled by default</label>
                                                <Switch
                                                    checked={newFlagValue}
                                                    onCheckedChange={setNewFlagValue}
                                                />
                                            </div>
                                            <Button
                                                onClick={() =>
                                                    createFeatureFlagMutation.mutate({
                                                        name: newFlagName,
                                                        enabled: newFlagValue,
                                                    })
                                                }
                                                disabled={!newFlagName || createFeatureFlagMutation.isPending}
                                                className="w-full"
                                            >
                                                Create Flag
                                            </Button>
                                        </div>
                                    </DialogContent>
                                </Dialog>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>Flag Name</TableHead>
                                        <TableHead>Description</TableHead>
                                        <TableHead>Environment</TableHead>
                                        <TableHead>Status</TableHead>
                                        <TableHead>Toggle</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {featureFlags?.data?.map((flag: any) => (
                                        <TableRow key={flag.name}>
                                            <TableCell className="font-mono font-medium">{flag.name}</TableCell>
                                            <TableCell className="text-sm text-gray-600">
                                                {flag.description}
                                            </TableCell>
                                            <TableCell>
                                                <Badge variant="outline">{flag.environment || 'all'}</Badge>
                                            </TableCell>
                                            <TableCell>
                                                <Badge variant={flag.enabled ? 'default' : 'secondary'}>
                                                    {flag.enabled ? 'Enabled' : 'Disabled'}
                                                </Badge>
                                            </TableCell>
                                            <TableCell>
                                                <Switch
                                                    checked={flag.enabled}
                                                    onCheckedChange={(checked) =>
                                                        toggleFeatureFlagMutation.mutate({
                                                            flag: flag.name,
                                                            enabled: checked,
                                                        })
                                                    }
                                                />
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}

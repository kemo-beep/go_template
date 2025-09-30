'use client';

import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Database, FolderOpen, Activity } from 'lucide-react';
import { api } from '@/lib/api-client';

export default function DashboardPage() {
    const { data: users } = useQuery({
        queryKey: ['users'],
        queryFn: () => api.getUsers(),
        retry: false, // Don't retry if admin API is not available
    });

    const { data: files } = useQuery({
        queryKey: ['files'],
        queryFn: () => api.getFiles(),
    });

    const stats = [
        {
            name: 'Total Users',
            value: users?.data?.length || 0,
            icon: Users,
            color: 'text-blue-600',
            bgColor: 'bg-blue-50',
        },
        {
            name: 'Database Tables',
            value: '3',
            icon: Database,
            color: 'text-green-600',
            bgColor: 'bg-green-50',
        },
        {
            name: 'Stored Files',
            value: files?.data?.length || 0,
            icon: FolderOpen,
            color: 'text-purple-600',
            bgColor: 'bg-purple-50',
        },
        {
            name: 'API Status',
            value: 'Healthy',
            icon: Activity,
            color: 'text-emerald-600',
            bgColor: 'bg-emerald-50',
        },
    ];

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
                <p className="text-gray-500 mt-1">
                    Overview of your backend system
                </p>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {stats.map((stat) => (
                    <Card key={stat.name}>
                        <CardHeader className="flex flex-row items-center justify-between pb-2">
                            <CardTitle className="text-sm font-medium text-gray-600">
                                {stat.name}
                            </CardTitle>
                            <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                                <stat.icon className={`h-4 w-4 ${stat.color}`} />
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{stat.value}</div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            <div className="grid gap-4 md:grid-cols-2">
                <Card>
                    <CardHeader>
                        <CardTitle>Quick Actions</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <p className="text-sm text-gray-600">
                            Use the sidebar to navigate to different sections of the admin console.
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>System Status</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2">
                            <div className="flex items-center justify-between">
                                <span className="text-sm">API Server</span>
                                <span className="text-sm font-medium text-green-600">Online</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-sm">Database</span>
                                <span className="text-sm font-medium text-green-600">Connected</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-sm">Storage (R2)</span>
                                <span className="text-sm font-medium text-green-600">Available</span>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}

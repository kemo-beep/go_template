'use client';

import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Users,
    Database,
    Activity,
    Wifi,
    Shield,
    FileText,
    TrendingUp,
    Clock,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { format } from 'date-fns';

export default function OverviewPage() {
    // Fetch various statistics
    const { data: dbStats } = useQuery({
        queryKey: ['database-stats'],
        queryFn: async () => {
            const response = await fetch('http://localhost:8080/api/v1/admin/database/stats', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('auth_token')}` },
            });
            if (!response.ok) throw new Error('Failed to fetch database stats');
            return response.json();
        },
    });

    const { data: realtimeStats } = useQuery({
        queryKey: ['realtime-stats'],
        queryFn: async () => {
            const response = await fetch('http://localhost:8080/api/v1/realtime/stats', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('auth_token')}` },
            });
            if (!response.ok) throw new Error('Failed to fetch realtime stats');
            return response.json();
        },
        refetchInterval: 5000,
    });

    const { data: usersData } = useQuery({
        queryKey: ['admin-users-count'],
        queryFn: async () => {
            const response = await fetch('http://localhost:8080/api/v1/admin/users?limit=1', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('auth_token')}` },
            });
            if (!response.ok) throw new Error('Failed to fetch users');
            return response.json();
        },
    });

    const { data: rolesData } = useQuery({
        queryKey: ['roles-count'],
        queryFn: async () => {
            const response = await fetch('http://localhost:8080/api/v1/admin/roles', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('auth_token')}` },
            });
            if (!response.ok) throw new Error('Failed to fetch roles');
            return response.json();
        },
    });

    const dbStatsData = dbStats?.data || {};
    const realtimeStatsData = realtimeStats?.data || {};
    const totalUsers = usersData?.data?.total || 0;
    const totalRoles = rolesData?.data?.total || 0;

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Admin Overview</h1>
                <p className="text-gray-500 mt-1">
                    Complete system statistics and real-time monitoring
                </p>
            </div>

            {/* Main Stats */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Users</CardTitle>
                        <Users className="h-4 w-4 text-blue-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{totalUsers}</div>
                        <p className="text-xs text-muted-foreground mt-1">
                            Registered accounts
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Database Tables</CardTitle>
                        <Database className="h-4 w-4 text-green-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{dbStatsData.table_count || 0}</div>
                        <p className="text-xs text-muted-foreground mt-1">
                            {dbStatsData.database_size || 'N/A'}
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Online Now</CardTitle>
                        <Wifi className="h-4 w-4 text-purple-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{realtimeStatsData.online_users || 0}</div>
                        <p className="text-xs text-muted-foreground mt-1">
                            {realtimeStatsData.total_clients || 0} total connections
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Roles</CardTitle>
                        <Shield className="h-4 w-4 text-orange-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{totalRoles}</div>
                        <p className="text-xs text-muted-foreground mt-1">
                            Permission groups
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* System Health */}
            <div className="grid gap-4 md:grid-cols-2">
                <Card>
                    <CardHeader>
                        <CardTitle className="text-sm">System Health</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Activity className="h-4 w-4 text-green-600" />
                                    <span className="text-sm">API Server</span>
                                </div>
                                <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                                    Healthy
                                </Badge>
                            </div>

                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Database className="h-4 w-4 text-green-600" />
                                    <span className="text-sm">Database</span>
                                </div>
                                <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                                    Connected
                                </Badge>
                            </div>

                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Wifi className="h-4 w-4 text-green-600" />
                                    <span className="text-sm">WebSocket</span>
                                </div>
                                <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                                    Active
                                </Badge>
                            </div>

                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Shield className="h-4 w-4 text-green-600" />
                                    <span className="text-sm">Authentication</span>
                                </div>
                                <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                                    Secure
                                </Badge>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle className="text-sm">Database Statistics</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <span className="text-sm text-gray-600">Total Rows</span>
                                <span className="font-bold">{dbStatsData.total_rows || 0}</span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm text-gray-600">Database Size</span>
                                <span className="font-bold">{dbStatsData.database_size || 'N/A'}</span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm text-gray-600">Table Count</span>
                                <span className="font-bold">{dbStatsData.table_count || 0}</span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm text-gray-600">Last Backup</span>
                                <span className="text-sm text-gray-500">Not configured</span>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Real-time Activity */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-sm">Real-time Channels</CardTitle>
                </CardHeader>
                <CardContent>
                    {realtimeStatsData.rooms && Object.keys(realtimeStatsData.rooms).length > 0 ? (
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            {Object.entries(realtimeStatsData.rooms).map(([room, count]: [string, any]) => (
                                <div key={room} className="border rounded-lg p-4 bg-gray-50">
                                    <div className="text-xs text-gray-500 mb-1">#{room}</div>
                                    <div className="text-2xl font-bold">{count}</div>
                                    <div className="text-xs text-gray-500">active</div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="text-center text-gray-500 py-8">
                            <Wifi className="h-8 w-8 mx-auto mb-2 text-gray-400" />
                            <p className="text-sm">No active channels</p>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-sm">Quick Actions</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <a
                            href="/dashboard/users"
                            className="border rounded-lg p-4 hover:bg-gray-50 transition-colors cursor-pointer"
                        >
                            <Users className="h-6 w-6 text-blue-600 mb-2" />
                            <div className="font-medium text-sm">Manage Users</div>
                            <div className="text-xs text-gray-500">View and edit users</div>
                        </a>

                        <a
                            href="/dashboard/database"
                            className="border rounded-lg p-4 hover:bg-gray-50 transition-colors cursor-pointer"
                        >
                            <Database className="h-6 w-6 text-green-600 mb-2" />
                            <div className="font-medium text-sm">Database</div>
                            <div className="text-xs text-gray-500">Explore tables</div>
                        </a>

                        <a
                            href="/dashboard/realtime"
                            className="border rounded-lg p-4 hover:bg-gray-50 transition-colors cursor-pointer"
                        >
                            <Activity className="h-6 w-6 text-purple-600 mb-2" />
                            <div className="font-medium text-sm">Real-time</div>
                            <div className="text-xs text-gray-500">Live monitoring</div>
                        </a>

                        <a
                            href="/dashboard/settings"
                            className="border rounded-lg p-4 hover:bg-gray-50 transition-colors cursor-pointer"
                        >
                            <Shield className="h-6 w-6 text-orange-600 mb-2" />
                            <div className="font-medium text-sm">Settings</div>
                            <div className="text-xs text-gray-500">System config</div>
                        </a>
                    </div>
                </CardContent>
            </Card>

            {/* Recent Activity */}
            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle className="text-sm">Recent Activity</CardTitle>
                        <Clock className="h-4 w-4 text-gray-400" />
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        <div className="flex items-start gap-3 text-sm">
                            <div className="h-2 w-2 rounded-full bg-blue-500 mt-2"></div>
                            <div className="flex-1">
                                <div className="font-medium">System started</div>
                                <div className="text-xs text-gray-500">Server is running on port 8080</div>
                            </div>
                            <div className="text-xs text-gray-400">Just now</div>
                        </div>

                        <div className="flex items-start gap-3 text-sm">
                            <div className="h-2 w-2 rounded-full bg-green-500 mt-2"></div>
                            <div className="flex-1">
                                <div className="font-medium">Database connected</div>
                                <div className="text-xs text-gray-500">PostgreSQL connection established</div>
                            </div>
                            <div className="text-xs text-gray-400">1 min ago</div>
                        </div>

                        <div className="flex items-start gap-3 text-sm">
                            <div className="h-2 w-2 rounded-full bg-purple-500 mt-2"></div>
                            <div className="flex-1">
                                <div className="font-medium">WebSocket server active</div>
                                <div className="text-xs text-gray-500">Real-time features enabled</div>
                            </div>
                            <div className="text-xs text-gray-400">2 min ago</div>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}

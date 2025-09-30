'use client';

import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    Activity,
    TrendingUp,
    Clock,
    AlertTriangle,
    Database,
    Zap,
    Users,
    Server,
} from 'lucide-react';
import { api } from '@/lib/api-client';
import { useState } from 'react';

export default function MetricsPage() {
    const [timeRange, setTimeRange] = useState('1h');

    const { data: metrics, isLoading } = useQuery({
        queryKey: ['metrics', timeRange],
        queryFn: () => api.getMetrics(),
        refetchInterval: 10000, // Auto-refresh every 10 seconds
    });

    const metricCards = [
        {
            title: 'Requests/min',
            value: metrics?.data?.[0]?.requests_per_minute || 0,
            icon: Activity,
            change: '+12%',
            positive: true,
        },
        {
            title: 'Avg Response Time',
            value: `${metrics?.data?.[0]?.avg_latency || 0}ms`,
            icon: Clock,
            change: '-5%',
            positive: true,
        },
        {
            title: 'Error Rate',
            value: `${metrics?.data?.[0]?.error_rate || 0}%`,
            icon: AlertTriangle,
            change: '+2%',
            positive: false,
        },
        {
            title: 'Success Rate',
            value: `${100 - (metrics?.data?.[0]?.error_rate || 0)}%`,
            icon: TrendingUp,
            change: '+0.2%',
            positive: true,
        },
    ];

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Metrics Dashboard</h1>
                    <p className="text-gray-500 mt-1">
                        Real-time application performance metrics
                    </p>
                </div>
                <Select value={timeRange} onValueChange={setTimeRange}>
                    <SelectTrigger className="w-40">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="5m">Last 5 min</SelectItem>
                        <SelectItem value="15m">Last 15 min</SelectItem>
                        <SelectItem value="1h">Last 1 hour</SelectItem>
                        <SelectItem value="6h">Last 6 hours</SelectItem>
                        <SelectItem value="24h">Last 24 hours</SelectItem>
                        <SelectItem value="7d">Last 7 days</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            {/* Metric Cards */}
            <div className="grid md:grid-cols-4 gap-4">
                {metricCards.map((metric) => {
                    const Icon = metric.icon;
                    return (
                        <Card key={metric.title}>
                            <CardHeader className="pb-3">
                                <div className="flex items-center justify-between">
                                    <CardTitle className="text-sm font-medium text-gray-500">
                                        {metric.title}
                                    </CardTitle>
                                    <Icon className="h-4 w-4 text-gray-400" />
                                </div>
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{metric.value}</div>
                                <Badge
                                    variant={metric.positive ? 'default' : 'secondary'}
                                    className="mt-2"
                                >
                                    {metric.change} from last hour
                                </Badge>
                            </CardContent>
                        </Card>
                    );
                })}
            </div>

            {/* Request Distribution */}
            <div className="grid md:grid-cols-2 gap-4">
                <Card>
                    <CardHeader>
                        <CardTitle>Request Distribution by Method</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {[
                                { method: 'GET', count: 1245, percent: 62, color: 'bg-blue-500' },
                                { method: 'POST', count: 487, percent: 24, color: 'bg-green-500' },
                                { method: 'PUT', count: 183, percent: 9, color: 'bg-yellow-500' },
                                { method: 'DELETE', count: 98, percent: 5, color: 'bg-red-500' },
                            ].map((item) => (
                                <div key={item.method}>
                                    <div className="flex items-center justify-between mb-1">
                                        <span className="text-sm font-medium">{item.method}</span>
                                        <span className="text-sm text-gray-500">
                                            {item.count} ({item.percent}%)
                                        </span>
                                    </div>
                                    <div className="w-full bg-gray-200 rounded-full h-2">
                                        <div
                                            className={`${item.color} h-2 rounded-full`}
                                            style={{ width: `${item.percent}%` }}
                                        />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>Response Status Codes</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {[
                                { status: '2xx Success', count: 1876, percent: 93, color: 'bg-green-500' },
                                { status: '4xx Client Error', count: 98, percent: 5, color: 'bg-yellow-500' },
                                { status: '5xx Server Error', count: 39, percent: 2, color: 'bg-red-500' },
                            ].map((item) => (
                                <div key={item.status}>
                                    <div className="flex items-center justify-between mb-1">
                                        <span className="text-sm font-medium">{item.status}</span>
                                        <span className="text-sm text-gray-500">
                                            {item.count} ({item.percent}%)
                                        </span>
                                    </div>
                                    <div className="w-full bg-gray-200 rounded-full h-2">
                                        <div
                                            className={`${item.color} h-2 rounded-full`}
                                            style={{ width: `${item.percent}%` }}
                                        />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Top Endpoints */}
            <Card>
                <CardHeader>
                    <CardTitle>Top API Endpoints</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        {[
                            { path: '/api/v1/users', requests: 4521, avg_time: '45ms', errors: 12 },
                            { path: '/api/v1/auth/login', requests: 3214, avg_time: '78ms', errors: 5 },
                            { path: '/api/v1/posts', requests: 2876, avg_time: '62ms', errors: 8 },
                            { path: '/api/v1/comments', requests: 1923, avg_time: '34ms', errors: 3 },
                            { path: '/api/v1/files/upload', requests: 1456, avg_time: '234ms', errors: 23 },
                        ].map((endpoint) => (
                            <div
                                key={endpoint.path}
                                className="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50"
                            >
                                <div className="flex-1">
                                    <code className="text-sm font-medium">{endpoint.path}</code>
                                    <div className="flex items-center gap-4 mt-1">
                                        <span className="text-xs text-gray-500">
                                            {endpoint.requests.toLocaleString()} requests
                                        </span>
                                        <span className="text-xs text-gray-500">
                                            Avg: {endpoint.avg_time}
                                        </span>
                                        <span className="text-xs text-red-600">
                                            {endpoint.errors} errors
                                        </span>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>

            {/* System Resources */}
            <Card>
                <CardHeader>
                    <CardTitle>System Resources</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid md:grid-cols-3 gap-6">
                        <div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm font-medium">CPU Usage</span>
                                <span className="text-sm text-gray-500">34%</span>
                            </div>
                            <div className="w-full bg-gray-200 rounded-full h-3">
                                <div className="bg-blue-500 h-3 rounded-full" style={{ width: '34%' }} />
                            </div>
                        </div>
                        <div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm font-medium">Memory Usage</span>
                                <span className="text-sm text-gray-500">58%</span>
                            </div>
                            <div className="w-full bg-gray-200 rounded-full h-3">
                                <div className="bg-green-500 h-3 rounded-full" style={{ width: '58%' }} />
                            </div>
                        </div>
                        <div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm font-medium">Disk Usage</span>
                                <span className="text-sm text-gray-500">42%</span>
                            </div>
                            <div className="w-full bg-gray-200 rounded-full h-3">
                                <div className="bg-yellow-500 h-3 rounded-full" style={{ width: '42%' }} />
                            </div>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}

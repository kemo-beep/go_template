"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { RefreshCw, Database, Zap, CheckCircle, AlertCircle } from "lucide-react";

interface AutoRegistryStatus {
    registered_apis: number;
    apis: Record<string, boolean>;
    watcher_running: boolean;
    last_checksum: string;
}

interface RegisteredAPIs {
    apis: string[];
    count: number;
}

export default function AutoRegistryStatus() {
    const [status, setStatus] = useState<AutoRegistryStatus | null>(null);
    const [registeredAPIs, setRegisteredAPIs] = useState<RegisteredAPIs | null>(null);
    const [loading, setLoading] = useState(true);
    const [regenerating, setRegenerating] = useState(false);

    const fetchStatus = async () => {
        try {
            const response = await fetch('/api/v1/admin/auto-registry/status', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                },
            });

            if (response.ok) {
                const data = await response.json();
                setStatus(data.data);
            }
        } catch (error) {
            console.error('Failed to fetch auto-registry status:', error);
        }
    };

    const fetchRegisteredAPIs = async () => {
        try {
            const response = await fetch('/api/v1/admin/auto-registry/apis', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                },
            });

            if (response.ok) {
                const data = await response.json();
                setRegisteredAPIs(data.data);
            }
        } catch (error) {
            console.error('Failed to fetch registered APIs:', error);
        }
    };

    const handleRegenerate = async () => {
        setRegenerating(true);
        try {
            const response = await fetch('/api/v1/admin/auto-registry/regenerate', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                },
            });

            if (response.ok) {
                // Refresh status after regeneration
                await fetchStatus();
                await fetchRegisteredAPIs();
            }
        } catch (error) {
            console.error('Failed to regenerate APIs:', error);
        } finally {
            setRegenerating(false);
        }
    };

    useEffect(() => {
        const loadData = async () => {
            setLoading(true);
            await Promise.all([fetchStatus(), fetchRegisteredAPIs()]);
            setLoading(false);
        };

        loadData();

        // Refresh every 30 seconds
        const interval = setInterval(loadData, 30000);
        return () => clearInterval(interval);
    }, []);

    if (loading) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <RefreshCw className="h-5 w-5 animate-spin" />
                        Auto-Registry Status
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p>Loading...</p>
                </CardContent>
            </Card>
        );
    }

    return (
        <div className="space-y-6">
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Database className="h-5 w-5" />
                        Auto-Registry Status
                    </CardTitle>
                    <CardDescription>
                        Monitor the automatic API generation and registration system
                    </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div className="flex items-center gap-2">
                            <span className="text-sm font-medium">Registered APIs:</span>
                            <Badge variant="secondary">{status?.registered_apis || 0}</Badge>
                        </div>

                        <div className="flex items-center gap-2">
                            <span className="text-sm font-medium">Schema Watcher:</span>
                            {status?.watcher_running ? (
                                <Badge variant="default" className="bg-green-100 text-green-800">
                                    <CheckCircle className="h-3 w-3 mr-1" />
                                    Running
                                </Badge>
                            ) : (
                                <Badge variant="destructive">
                                    <AlertCircle className="h-3 w-3 mr-1" />
                                    Stopped
                                </Badge>
                            )}
                        </div>

                        <div className="flex items-center gap-2">
                            <span className="text-sm font-medium">Last Checksum:</span>
                            <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                                {status?.last_checksum?.substring(0, 8) || 'N/A'}...
                            </code>
                        </div>
                    </div>

                    <div className="flex gap-2">
                        <Button
                            onClick={() => {
                                fetchStatus();
                                fetchRegisteredAPIs();
                            }}
                            variant="outline"
                            size="sm"
                        >
                            <RefreshCw className="h-4 w-4 mr-2" />
                            Refresh
                        </Button>

                        <Button
                            onClick={handleRegenerate}
                            disabled={regenerating}
                            size="sm"
                        >
                            {regenerating ? (
                                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                            ) : (
                                <Zap className="h-4 w-4 mr-2" />
                            )}
                            {regenerating ? 'Regenerating...' : 'Regenerate APIs'}
                        </Button>
                    </div>
                </CardContent>
            </Card>

            {registeredAPIs && (
                <Card>
                    <CardHeader>
                        <CardTitle>Registered APIs</CardTitle>
                        <CardDescription>
                            List of currently registered auto-generated APIs
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
                            {registeredAPIs.apis.map((api) => (
                                <Badge key={api} variant="outline" className="justify-start">
                                    {api}
                                </Badge>
                            ))}
                        </div>
                        {registeredAPIs.apis.length === 0 && (
                            <p className="text-sm text-gray-500">No APIs registered yet</p>
                        )}
                    </CardContent>
                </Card>
            )}
        </div>
    );
}

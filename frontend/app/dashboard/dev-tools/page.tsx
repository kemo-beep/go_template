'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import {
    Code,
    Database,
    Zap,
    Settings,
    Download,
    Upload,
    Play,
    Pause,
    RotateCcw,
    Eye,
    EyeOff,
    Copy,
    CheckCircle,
    AlertCircle,
    Clock,
    Activity,
    Terminal,
    FileText,
    Wrench,
    RefreshCw
} from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

interface HealthCheck {
    service: string;
    status: 'healthy' | 'unhealthy' | 'checking';
    responseTime?: number;
    lastChecked: string;
    error?: string;
}

interface CodeGenerator {
    language: string;
    framework: string;
    template: string;
}

export default function DevToolsPage() {
    const [healthChecks, setHealthChecks] = useState<HealthCheck[]>([]);
    const [isRunningChecks, setIsRunningChecks] = useState(false);
    const [codeGenerator, setCodeGenerator] = useState<CodeGenerator>({
        language: 'typescript',
        framework: 'react',
        template: 'api-client'
    });
    const [generatedCode, setGeneratedCode] = useState('');
    const [autoRefresh, setAutoRefresh] = useState(false);
    const [refreshInterval, setRefreshInterval] = useState(30);
    const { toast } = useToast();

    const services = [
        { name: 'API Server', endpoint: '/healthz' },
        { name: 'Database', endpoint: '/api/v1/health/database' },
        { name: 'Redis Cache', endpoint: '/api/v1/health/redis' },
        { name: 'File Storage', endpoint: '/api/v1/health/storage' },
    ];

    useEffect(() => {
        runHealthChecks();

        if (autoRefresh) {
            const interval = setInterval(runHealthChecks, refreshInterval * 1000);
            return () => clearInterval(interval);
        }
    }, [autoRefresh, refreshInterval]);

    const runHealthChecks = async () => {
        setIsRunningChecks(true);
        const checks: HealthCheck[] = [];

        for (const service of services) {
            const startTime = Date.now();
            try {
                const response = await fetch(service.endpoint, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json' }
                });

                const responseTime = Date.now() - startTime;

                checks.push({
                    service: service.name,
                    status: response.ok ? 'healthy' : 'unhealthy',
                    responseTime,
                    lastChecked: new Date().toISOString(),
                    error: response.ok ? undefined : `HTTP ${response.status}`
                });
            } catch (error) {
                checks.push({
                    service: service.name,
                    status: 'unhealthy',
                    lastChecked: new Date().toISOString(),
                    error: error instanceof Error ? error.message : 'Unknown error'
                });
            }
        }

        setHealthChecks(checks);
        setIsRunningChecks(false);
    };

    const generateCode = async () => {
        try {
            const response = await fetch('/api/v1/generate-code', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(codeGenerator)
            });

            if (response.ok) {
                const code = await response.text();
                setGeneratedCode(code);
                toast({
                    title: "Code Generated",
                    description: `${codeGenerator.language} code generated successfully`,
                    variant: "success"
                });
            } else {
                throw new Error('Failed to generate code');
            }
        } catch (error) {
            toast({
                title: "Error",
                description: "Failed to generate code",
                variant: "destructive"
            });
        }
    };

    const copyToClipboard = (text: string, label: string) => {
        navigator.clipboard.writeText(text);
        toast({
            title: "Copied!",
            description: `${label} copied to clipboard`,
        });
    };

    const exportLogs = () => {
        const logs = healthChecks.map(check => ({
            service: check.service,
            status: check.status,
            responseTime: check.responseTime,
            timestamp: check.lastChecked,
            error: check.error
        }));

        const blob = new Blob([JSON.stringify(logs, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `health-checks-${new Date().toISOString().split('T')[0]}.json`;
        a.click();
        URL.revokeObjectURL(url);

        toast({
            title: "Exported",
            description: "Health check logs exported successfully",
        });
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'healthy':
                return <CheckCircle className="h-4 w-4 text-green-500" />;
            case 'unhealthy':
                return <AlertCircle className="h-4 w-4 text-red-500" />;
            case 'checking':
                return <RefreshCw className="h-4 w-4 text-blue-500 animate-spin" />;
            default:
                return <Clock className="h-4 w-4 text-gray-500" />;
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'healthy':
                return 'bg-green-100 text-green-800 border-green-200';
            case 'unhealthy':
                return 'bg-red-100 text-red-800 border-red-200';
            case 'checking':
                return 'bg-blue-100 text-blue-800 border-blue-200';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-200';
        }
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Developer Tools</h1>
                    <p className="text-muted-foreground">
                        Advanced tools for development and debugging
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={exportLogs}>
                        <Download className="h-4 w-4 mr-2" />
                        Export Logs
                    </Button>
                    <Button onClick={runHealthChecks} disabled={isRunningChecks}>
                        {isRunningChecks ? (
                            <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                        ) : (
                            <Play className="h-4 w-4 mr-2" />
                        )}
                        {isRunningChecks ? 'Running...' : 'Run Checks'}
                    </Button>
                </div>
            </div>

            <Tabs defaultValue="health" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="health">Health Monitoring</TabsTrigger>
                    <TabsTrigger value="codegen">Code Generator</TabsTrigger>
                    <TabsTrigger value="debug">Debug Tools</TabsTrigger>
                    <TabsTrigger value="settings">Settings</TabsTrigger>
                </TabsList>

                {/* Health Monitoring */}
                <TabsContent value="health" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <div>
                                    <CardTitle className="flex items-center gap-2">
                                        <Activity className="h-5 w-5" />
                                        System Health
                                    </CardTitle>
                                    <CardDescription>
                                        Monitor the health of all system components
                                    </CardDescription>
                                </div>
                                <div className="flex items-center gap-4">
                                    <div className="flex items-center gap-2">
                                        <Switch
                                            id="auto-refresh"
                                            checked={autoRefresh}
                                            onCheckedChange={setAutoRefresh}
                                        />
                                        <Label htmlFor="auto-refresh">Auto Refresh</Label>
                                    </div>
                                    {autoRefresh && (
                                        <Select value={refreshInterval.toString()} onValueChange={(value) => setRefreshInterval(Number(value))}>
                                            <SelectTrigger className="w-32">
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="10">10s</SelectItem>
                                                <SelectItem value="30">30s</SelectItem>
                                                <SelectItem value="60">1m</SelectItem>
                                                <SelectItem value="300">5m</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    )}
                                </div>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="grid gap-4">
                                {healthChecks.map((check, index) => (
                                    <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                                        <div className="flex items-center gap-3">
                                            {getStatusIcon(check.status)}
                                            <div>
                                                <div className="font-medium">{check.service}</div>
                                                <div className="text-sm text-muted-foreground">
                                                    Last checked: {new Date(check.lastChecked).toLocaleTimeString()}
                                                </div>
                                                {check.error && (
                                                    <div className="text-sm text-red-600 mt-1">
                                                        Error: {check.error}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-2">
                                            {check.responseTime && (
                                                <Badge variant="outline">
                                                    {check.responseTime}ms
                                                </Badge>
                                            )}
                                            <Badge className={getStatusColor(check.status)}>
                                                {check.status}
                                            </Badge>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Code Generator */}
                <TabsContent value="codegen" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Code className="h-5 w-5" />
                                Code Generator
                            </CardTitle>
                            <CardDescription>
                                Generate client code for your API
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <Label htmlFor="language">Language</Label>
                                    <Select value={codeGenerator.language} onValueChange={(value) => setCodeGenerator(prev => ({ ...prev, language: value }))}>
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="typescript">TypeScript</SelectItem>
                                            <SelectItem value="javascript">JavaScript</SelectItem>
                                            <SelectItem value="python">Python</SelectItem>
                                            <SelectItem value="go">Go</SelectItem>
                                            <SelectItem value="java">Java</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div>
                                    <Label htmlFor="framework">Framework</Label>
                                    <Select value={codeGenerator.framework} onValueChange={(value) => setCodeGenerator(prev => ({ ...prev, framework: value }))}>
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="react">React</SelectItem>
                                            <SelectItem value="vue">Vue.js</SelectItem>
                                            <SelectItem value="angular">Angular</SelectItem>
                                            <SelectItem value="vanilla">Vanilla JS</SelectItem>
                                            <SelectItem value="axios">Axios</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div>
                                    <Label htmlFor="template">Template</Label>
                                    <Select value={codeGenerator.template} onValueChange={(value) => setCodeGenerator(prev => ({ ...prev, template: value }))}>
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="api-client">API Client</SelectItem>
                                            <SelectItem value="hooks">React Hooks</SelectItem>
                                            <SelectItem value="service">Service Class</SelectItem>
                                            <SelectItem value="types">Type Definitions</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <Button onClick={generateCode} className="w-full">
                                <Code className="h-4 w-4 mr-2" />
                                Generate Code
                            </Button>
                            {generatedCode && (
                                <div className="space-y-2">
                                    <div className="flex items-center justify-between">
                                        <Label>Generated Code</Label>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={() => copyToClipboard(generatedCode, 'Generated code')}
                                        >
                                            <Copy className="h-4 w-4 mr-2" />
                                            Copy
                                        </Button>
                                    </div>
                                    <Textarea
                                        value={generatedCode}
                                        readOnly
                                        className="font-mono text-sm min-h-[300px]"
                                    />
                                </div>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Debug Tools */}
                <TabsContent value="debug" className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Terminal className="h-5 w-5" />
                                    API Testing
                                </CardTitle>
                                <CardDescription>
                                    Test API endpoints with custom requests
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div>
                                    <Label htmlFor="endpoint">Endpoint</Label>
                                    <Input placeholder="/api/v1/users" />
                                </div>
                                <div>
                                    <Label htmlFor="method">Method</Label>
                                    <Select>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select method" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="GET">GET</SelectItem>
                                            <SelectItem value="POST">POST</SelectItem>
                                            <SelectItem value="PUT">PUT</SelectItem>
                                            <SelectItem value="DELETE">DELETE</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div>
                                    <Label htmlFor="body">Request Body</Label>
                                    <Textarea placeholder='{"key": "value"}' />
                                </div>
                                <Button className="w-full">
                                    <Play className="h-4 w-4 mr-2" />
                                    Send Request
                                </Button>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Database className="h-5 w-5" />
                                    Database Tools
                                </CardTitle>
                                <CardDescription>
                                    Database management and monitoring
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="space-y-2">
                                    <Button variant="outline" className="w-full justify-start">
                                        <Database className="h-4 w-4 mr-2" />
                                        View Tables
                                    </Button>
                                    <Button variant="outline" className="w-full justify-start">
                                        <FileText className="h-4 w-4 mr-2" />
                                        Export Schema
                                    </Button>
                                    <Button variant="outline" className="w-full justify-start">
                                        <RefreshCw className="h-4 w-4 mr-2" />
                                        Run Migrations
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                {/* Settings */}
                <TabsContent value="settings" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Settings className="h-5 w-5" />
                                Developer Settings
                            </CardTitle>
                            <CardDescription>
                                Configure development environment settings
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <div>
                                        <Label htmlFor="debug-mode">Debug Mode</Label>
                                        <p className="text-sm text-muted-foreground">
                                            Enable detailed logging and debugging information
                                        </p>
                                    </div>
                                    <Switch id="debug-mode" />
                                </div>
                                <div className="flex items-center justify-between">
                                    <div>
                                        <Label htmlFor="auto-save">Auto Save</Label>
                                        <p className="text-sm text-muted-foreground">
                                            Automatically save changes to files
                                        </p>
                                    </div>
                                    <Switch id="auto-save" defaultChecked />
                                </div>
                                <div className="flex items-center justify-between">
                                    <div>
                                        <Label htmlFor="hot-reload">Hot Reload</Label>
                                        <p className="text-sm text-muted-foreground">
                                            Automatically reload on file changes
                                        </p>
                                    </div>
                                    <Switch id="hot-reload" defaultChecked />
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
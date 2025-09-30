'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    Settings as SettingsIcon,
    Shield,
    Database,
    Bell,
    Mail,
    Save,
} from 'lucide-react';
import { api } from '@/lib/api-client';
import { toast } from 'sonner';

export default function SettingsPage() {
    const queryClient = useQueryClient();

    const [appName, setAppName] = useState('Go Mobile Backend');
    const [appUrl, setAppUrl] = useState('https://api.example.com');
    const [jwtExpiry, setJwtExpiry] = useState('24h');
    const [maxUploadSize, setMaxUploadSize] = useState('10MB');

    const [emailNotifications, setEmailNotifications] = useState(true);
    const [slackNotifications, setSlackNotifications] = useState(false);
    const [maintenanceMode, setMaintenanceMode] = useState(false);
    const [rateLimitEnabled, setRateLimitEnabled] = useState(true);
    const [rateLimitRequests, setRateLimitRequests] = useState('100');

    const { data: settings } = useQuery({
        queryKey: ['settings'],
        queryFn: () => api.getSettings(),
    });

    const updateSettingsMutation = useMutation({
        mutationFn: (data: any) => api.updateSettings(data),
        onSuccess: () => {
            toast.success('Settings updated successfully');
            queryClient.invalidateQueries({ queryKey: ['settings'] });
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Update failed');
        },
    });

    const handleSaveGeneral = () => {
        updateSettingsMutation.mutate({
            app_name: appName,
            app_url: appUrl,
            jwt_expiry: jwtExpiry,
            max_upload_size: maxUploadSize,
        });
    };

    const handleSaveNotifications = () => {
        updateSettingsMutation.mutate({
            email_notifications: emailNotifications,
            slack_notifications: slackNotifications,
        });
    };

    const handleSaveSecurity = () => {
        updateSettingsMutation.mutate({
            maintenance_mode: maintenanceMode,
            rate_limit_enabled: rateLimitEnabled,
            rate_limit_requests: parseInt(rateLimitRequests),
        });
    };

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
                <p className="text-gray-500 mt-1">
                    Manage application configuration and preferences
                </p>
            </div>

            <Tabs defaultValue="general" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="general">
                        <SettingsIcon className="h-4 w-4 mr-2" />
                        General
                    </TabsTrigger>
                    <TabsTrigger value="security">
                        <Shield className="h-4 w-4 mr-2" />
                        Security
                    </TabsTrigger>
                    <TabsTrigger value="database">
                        <Database className="h-4 w-4 mr-2" />
                        Database
                    </TabsTrigger>
                    <TabsTrigger value="notifications">
                        <Bell className="h-4 w-4 mr-2" />
                        Notifications
                    </TabsTrigger>
                </TabsList>

                {/* General Settings */}
                <TabsContent value="general" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Application Settings</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid md:grid-cols-2 gap-4">
                                <div>
                                    <Label htmlFor="app-name">Application Name</Label>
                                    <Input
                                        id="app-name"
                                        value={appName}
                                        onChange={(e) => setAppName(e.target.value)}
                                        placeholder="My Application"
                                    />
                                </div>
                                <div>
                                    <Label htmlFor="app-url">Application URL</Label>
                                    <Input
                                        id="app-url"
                                        value={appUrl}
                                        onChange={(e) => setAppUrl(e.target.value)}
                                        placeholder="https://api.example.com"
                                    />
                                </div>
                                <div>
                                    <Label htmlFor="jwt-expiry">JWT Token Expiry</Label>
                                    <Select value={jwtExpiry} onValueChange={setJwtExpiry}>
                                        <SelectTrigger id="jwt-expiry">
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="1h">1 Hour</SelectItem>
                                            <SelectItem value="6h">6 Hours</SelectItem>
                                            <SelectItem value="24h">24 Hours</SelectItem>
                                            <SelectItem value="7d">7 Days</SelectItem>
                                            <SelectItem value="30d">30 Days</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div>
                                    <Label htmlFor="max-upload">Max Upload Size</Label>
                                    <Select value={maxUploadSize} onValueChange={setMaxUploadSize}>
                                        <SelectTrigger id="max-upload">
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="1MB">1 MB</SelectItem>
                                            <SelectItem value="5MB">5 MB</SelectItem>
                                            <SelectItem value="10MB">10 MB</SelectItem>
                                            <SelectItem value="50MB">50 MB</SelectItem>
                                            <SelectItem value="100MB">100 MB</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <Button onClick={handleSaveGeneral} disabled={updateSettingsMutation.isPending}>
                                <Save className="h-4 w-4 mr-2" />
                                Save Changes
                            </Button>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Security Settings */}
                <TabsContent value="security" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Security & Access Control</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="flex items-center justify-between">
                                <div>
                                    <Label htmlFor="maintenance-mode" className="text-base">
                                        Maintenance Mode
                                    </Label>
                                    <p className="text-sm text-gray-500">
                                        Temporarily disable API access for maintenance
                                    </p>
                                </div>
                                <Switch
                                    id="maintenance-mode"
                                    checked={maintenanceMode}
                                    onCheckedChange={setMaintenanceMode}
                                />
                            </div>

                            <div className="flex items-center justify-between">
                                <div>
                                    <Label htmlFor="rate-limit" className="text-base">
                                        Rate Limiting
                                    </Label>
                                    <p className="text-sm text-gray-500">
                                        Limit API requests per IP address
                                    </p>
                                </div>
                                <Switch
                                    id="rate-limit"
                                    checked={rateLimitEnabled}
                                    onCheckedChange={setRateLimitEnabled}
                                />
                            </div>

                            {rateLimitEnabled && (
                                <div>
                                    <Label htmlFor="rate-limit-requests">Requests per minute</Label>
                                    <Input
                                        id="rate-limit-requests"
                                        type="number"
                                        value={rateLimitRequests}
                                        onChange={(e) => setRateLimitRequests(e.target.value)}
                                        className="w-32"
                                    />
                                </div>
                            )}

                            <Button onClick={handleSaveSecurity} disabled={updateSettingsMutation.isPending}>
                                <Save className="h-4 w-4 mr-2" />
                                Save Security Settings
                            </Button>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>API Keys</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-3">
                            <div className="flex items-center justify-between p-3 border rounded-lg">
                                <div>
                                    <p className="font-medium">Production API Key</p>
                                    <code className="text-sm text-gray-500">sk_prod_••••••••••••••••</code>
                                </div>
                                <Button variant="outline" size="sm">Regenerate</Button>
                            </div>
                            <div className="flex items-center justify-between p-3 border rounded-lg">
                                <div>
                                    <p className="font-medium">Development API Key</p>
                                    <code className="text-sm text-gray-500">sk_dev_••••••••••••••••</code>
                                </div>
                                <Button variant="outline" size="sm">Regenerate</Button>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Database Settings */}
                <TabsContent value="database" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Database Configuration</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid md:grid-cols-2 gap-4">
                                <div>
                                    <Label>Database Host</Label>
                                    <Input value="localhost:5432" disabled />
                                </div>
                                <div>
                                    <Label>Database Name</Label>
                                    <Input value="go_mobile_backend" disabled />
                                </div>
                                <div>
                                    <Label>Max Connections</Label>
                                    <Input type="number" defaultValue="100" />
                                </div>
                                <div>
                                    <Label>Connection Timeout</Label>
                                    <Select defaultValue="30s">
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="10s">10 seconds</SelectItem>
                                            <SelectItem value="30s">30 seconds</SelectItem>
                                            <SelectItem value="60s">60 seconds</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                            <Button disabled={updateSettingsMutation.isPending}>
                                <Save className="h-4 w-4 mr-2" />
                                Save Database Settings
                            </Button>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>Connection Pool</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="grid md:grid-cols-3 gap-4">
                                <div>
                                    <p className="text-sm text-gray-500">Active Connections</p>
                                    <p className="text-2xl font-bold">24</p>
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Idle Connections</p>
                                    <p className="text-2xl font-bold">6</p>
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Max Connections</p>
                                    <p className="text-2xl font-bold">100</p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Notifications Settings */}
                <TabsContent value="notifications" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Notification Preferences</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="flex items-center justify-between">
                                <div>
                                    <Label htmlFor="email-notif" className="text-base">
                                        Email Notifications
                                    </Label>
                                    <p className="text-sm text-gray-500">
                                        Receive alerts and updates via email
                                    </p>
                                </div>
                                <Switch
                                    id="email-notif"
                                    checked={emailNotifications}
                                    onCheckedChange={setEmailNotifications}
                                />
                            </div>

                            <div className="flex items-center justify-between">
                                <div>
                                    <Label htmlFor="slack-notif" className="text-base">
                                        Slack Notifications
                                    </Label>
                                    <p className="text-sm text-gray-500">
                                        Send alerts to your Slack workspace
                                    </p>
                                </div>
                                <Switch
                                    id="slack-notif"
                                    checked={slackNotifications}
                                    onCheckedChange={setSlackNotifications}
                                />
                            </div>

                            {slackNotifications && (
                                <div>
                                    <Label htmlFor="slack-webhook">Slack Webhook URL</Label>
                                    <Input
                                        id="slack-webhook"
                                        placeholder="https://hooks.slack.com/services/..."
                                        className="mt-1"
                                    />
                                </div>
                            )}

                            <Button onClick={handleSaveNotifications} disabled={updateSettingsMutation.isPending}>
                                <Save className="h-4 w-4 mr-2" />
                                Save Notification Settings
                            </Button>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>Alert Triggers</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            {[
                                'High error rate (>5%)',
                                'API response time >1s',
                                'Database connection errors',
                                'Failed authentication attempts',
                                'Server CPU usage >80%',
                            ].map((trigger, idx) => (
                                <div key={idx} className="flex items-center justify-between p-3 border rounded-lg">
                                    <span className="text-sm">{trigger}</span>
                                    <Switch defaultChecked={idx < 3} />
                                </div>
                            ))}
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}

'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
    Activity,
    Users,
    Wifi,
    WifiOff,
    Send,
    Globe,
    Radio,
} from 'lucide-react';
import { toast } from 'sonner';
import { format } from 'date-fns';
import { api } from '@/lib/api-client';

interface Message {
    type: string;
    channel?: string;
    event: string;
    payload: Record<string, unknown>;
    user_id?: number;
    timestamp: string;
}

interface PresenceInfo {
    user_id: number;
    username: string;
    status: string;
    last_seen: string;
    device_info?: string;
}

export default function RealtimePage() {
    const [connected, setConnected] = useState(false);
    const [messages, setMessages] = useState<Message[]>([]);
    const [messageInput, setMessageInput] = useState('');
    const [channel, setChannel] = useState('general');
    const [reconnectAttempts, setReconnectAttempts] = useState(0);
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    // Fetch presence data
    const { data: presenceData, refetch: refetchPresence } = useQuery({
        queryKey: ['presence'],
        queryFn: () => api.getRealtimePresence(),
        refetchInterval: 5000, // Refresh every 5 seconds
    });

    // Fetch real-time stats
    const { data: statsData } = useQuery({
        queryKey: ['realtime-stats'],
        queryFn: () => api.getRealtimeStats(),
        refetchInterval: 3000,
    });

    const connectWebSocket = useCallback(() => {
        const token = localStorage.getItem('auth_token');
        if (!token) {
            toast.error('No authentication token found');
            return;
        }

        const wsUrl = process.env.NEXT_PUBLIC_API_URL?.replace('http', 'ws') || 'ws://localhost:8080';
        const ws = new WebSocket(`${wsUrl}/api/v1/realtime/ws?token=${encodeURIComponent(token)}&channel=${channel}`);

        ws.onopen = () => {
            setConnected(true);
            setReconnectAttempts(0);
            toast.success('Connected to real-time server');

            // Subscribe to general channel
            ws.send(JSON.stringify({
                type: 'subscribe',
                channel: 'general',
            }));
        };

        ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                setMessages((prev) => [...prev, message].slice(-50)); // Keep last 50 messages

                if (message.type === 'presence') {
                    refetchPresence();
                }
            } catch (error) {
                console.error('Failed to parse message:', error);
            }
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            console.error('WebSocket readyState:', ws.readyState);
            console.error('WebSocket url:', ws.url);
            toast.error('WebSocket connection error. Is the backend server running?');
        };

        ws.onclose = (event) => {
            setConnected(false);
            console.log('WebSocket closed:', event.code, event.reason);

            if (event.code !== 1000) { // Not a normal closure
                toast.info('Disconnected from real-time server. Attempting to reconnect...');

                // Exponential backoff for reconnection
                const maxAttempts = 5;
                if (reconnectAttempts < maxAttempts) {
                    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000); // Max 30 seconds

                    reconnectTimeoutRef.current = setTimeout(() => {
                        setReconnectAttempts(prev => prev + 1);
                        connectWebSocket();
                    }, delay);
                } else {
                    toast.error('Failed to reconnect after multiple attempts. Please refresh the page.');
                }
            } else {
                toast.info('Disconnected from real-time server');
            }
        };

        wsRef.current = ws;
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        connectWebSocket();

        return () => {
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            if (wsRef.current) {
                wsRef.current.close(1000); // Normal closure
            }
        };
    }, []); // Remove connectWebSocket from dependencies to prevent reconnection loop

    const sendMessage = () => {
        if (!messageInput.trim() || !wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
            return;
        }

        const message = {
            type: 'broadcast',
            channel: 'general',
            event: 'user_message',
            payload: {
                text: messageInput,
                timestamp: new Date().toISOString(),
            },
        };

        wsRef.current.send(JSON.stringify(message));
        setMessageInput('');
    };

    const subscribeToChannel = (newChannel: string) => {
        if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
            return;
        }

        wsRef.current.send(JSON.stringify({
            type: 'subscribe',
            channel: newChannel,
        }));

        setChannel(newChannel);
        toast.success(`Subscribed to ${newChannel}`);
    };

    const presence: PresenceInfo[] = presenceData?.data || {};
    const onlineUsers = Object.values(presence).filter((p: PresenceInfo) => p.status === 'online');
    const stats = statsData?.data || {};

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Real-time Dashboard</h1>
                    <p className="text-gray-500 mt-1">
                        WebSocket connections, presence tracking, and live updates
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    {connected ? (
                        <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                            <Wifi className="h-3 w-3 mr-1" />
                            Connected
                        </Badge>
                    ) : (
                        <Badge variant="destructive">
                            <WifiOff className="h-3 w-3 mr-1" />
                            Disconnected
                        </Badge>
                    )}
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid gap-4 md:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Connections</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats.total_clients || 0}</div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Online Users</CardTitle>
                        <Users className="h-4 w-4 text-green-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats.online_users || 0}</div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Channels</CardTitle>
                        <Radio className="h-4 w-4 text-blue-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats.total_rooms || 0}</div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Current Channel</CardTitle>
                        <Globe className="h-4 w-4 text-purple-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-lg font-bold">{channel}</div>
                    </CardContent>
                </Card>
            </div>

            <div className="grid md:grid-cols-3 gap-6">
                {/* Online Users */}
                <Card className="md:col-span-1">
                    <CardHeader>
                        <CardTitle className="text-sm">Online Users ({onlineUsers.length})</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2 max-h-[400px] overflow-y-auto">
                            {onlineUsers.map((user: PresenceInfo) => (
                                <div key={user.user_id} className="flex items-center gap-3 p-2 rounded-lg bg-gray-50">
                                    <div className="h-2 w-2 rounded-full bg-green-500"></div>
                                    <div className="flex-1">
                                        <div className="font-medium text-sm">{user.username || `User ${user.user_id}`}</div>
                                        <div className="text-xs text-gray-500">
                                            {user.last_seen ? format(new Date(user.last_seen), 'HH:mm:ss') : 'Online'}
                                        </div>
                                    </div>
                                    <Badge variant="outline" className="text-xs">
                                        {user.status}
                                    </Badge>
                                </div>
                            ))}
                            {onlineUsers.length === 0 && (
                                <div className="text-center text-gray-500 py-8 text-sm">
                                    No users currently online
                                </div>
                            )}
                        </div>
                    </CardContent>
                </Card>

                {/* Messages Feed */}
                <Card className="md:col-span-2">
                    <CardHeader>
                        <CardTitle className="text-sm">Live Messages</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {/* Messages List */}
                            <div className="space-y-2 max-h-[300px] overflow-y-auto border rounded-lg p-4 bg-gray-50">
                                {messages.length === 0 ? (
                                    <div className="text-center text-gray-500 py-8 text-sm">
                                        <Activity className="h-8 w-8 mx-auto mb-2 text-gray-400" />
                                        Waiting for messages...
                                    </div>
                                ) : (
                                    messages.map((msg, index) => (
                                        <div key={index} className="text-sm p-2 rounded bg-white border">
                                            <div className="flex items-center gap-2 mb-1">
                                                <Badge variant="outline" className="text-xs">
                                                    {msg.event}
                                                </Badge>
                                                {msg.channel && (
                                                    <span className="text-xs text-gray-500">#{msg.channel}</span>
                                                )}
                                                <span className="text-xs text-gray-400 ml-auto">
                                                    {msg.timestamp ? format(new Date(msg.timestamp), 'HH:mm:ss') : '--:--:--'}
                                                </span>
                                            </div>
                                            <div className="text-xs text-gray-600">
                                                {JSON.stringify(msg.payload, null, 2)}
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>

                            {/* Send Message */}
                            <div className="flex gap-2">
                                <Input
                                    placeholder="Type a message..."
                                    value={messageInput}
                                    onChange={(e) => setMessageInput(e.target.value)}
                                    onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                                    disabled={!connected}
                                />
                                <Button
                                    onClick={sendMessage}
                                    disabled={!connected || !messageInput.trim()}
                                >
                                    <Send className="h-4 w-4" />
                                </Button>
                            </div>

                            {/* Channel Selector */}
                            <div className="flex gap-2 flex-wrap">
                                <span className="text-sm text-gray-500">Channels:</span>
                                {['general', 'db:users', 'db:files', 'db:*'].map((ch) => (
                                    <Button
                                        key={ch}
                                        variant={channel === ch ? 'default' : 'outline'}
                                        size="sm"
                                        onClick={() => subscribeToChannel(ch)}
                                        disabled={!connected}
                                    >
                                        #{ch}
                                    </Button>
                                ))}
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Channel Stats */}
            {stats.rooms && Object.keys(stats.rooms).length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle className="text-sm">Channel Statistics</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            {Object.entries(stats.rooms as Record<string, number>).map(([room, count]) => (
                                <div key={room} className="border rounded-lg p-3">
                                    <div className="text-xs text-gray-500 mb-1">#{room}</div>
                                    <div className="text-2xl font-bold">{count}</div>
                                    <div className="text-xs text-gray-500">connections</div>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}

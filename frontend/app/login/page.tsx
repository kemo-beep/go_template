'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation } from '@tanstack/react-query';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Shield, Loader2 } from 'lucide-react';
import { api } from '@/lib/api-client';
import { useAuthStore } from '@/lib/store';
import { toast } from 'sonner';

export default function LoginPage() {
    const router = useRouter();
    const { login } = useAuthStore();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const loginMutation = useMutation({
        mutationFn: () => api.login(email, password),
        onSuccess: (response) => {
            // Backend returns: { success: true, message: "...", data: { access_token, refresh_token, user, expires_in } }
            const { data: authData } = response.data;
            const { access_token, user } = authData;
            login(user, access_token);
            toast.success('Login successful!');
            router.push('/dashboard');
        },
        onError: (error: any) => {
            const errorMessage = error.response?.data?.message || error.response?.data?.error || 'Login failed. Please check your credentials.';
            toast.error(errorMessage);
            console.error('Login error:', error.response?.data);
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!email || !password) {
            toast.error('Please enter both email and password');
            return;
        }
        loginMutation.mutate();
    };

    return (
        <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white flex items-center justify-center p-4">
            <div className="w-full max-w-md">
                {/* Logo */}
                <div className="flex justify-center mb-8">
                    <div className="p-4 bg-blue-600 rounded-2xl">
                        <Shield className="h-12 w-12 text-white" />
                    </div>
                </div>

                {/* Login Card */}
                <Card>
                    <CardHeader className="space-y-1">
                        <CardTitle className="text-2xl font-bold text-center">
                            Admin Console
                        </CardTitle>
                        <CardDescription className="text-center">
                            Sign in to access the developer dashboard
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="space-y-2">
                                <label htmlFor="email" className="text-sm font-medium">
                                    Email
                                </label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="admin@example.com"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    disabled={loginMutation.isPending}
                                    autoComplete="email"
                                    required
                                />
                            </div>

                            <div className="space-y-2">
                                <label htmlFor="password" className="text-sm font-medium">
                                    Password
                                </label>
                                <Input
                                    id="password"
                                    type="password"
                                    placeholder="••••••••"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    disabled={loginMutation.isPending}
                                    autoComplete="current-password"
                                    required
                                />
                            </div>

                            <Button
                                type="submit"
                                className="w-full"
                                disabled={loginMutation.isPending}
                            >
                                {loginMutation.isPending ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        Signing in...
                                    </>
                                ) : (
                                    'Sign In'
                                )}
                            </Button>
                        </form>

                        {/* Demo Credentials */}
                        <div className="mt-6 p-4 bg-blue-50 rounded-lg">
                            <p className="text-sm font-medium text-blue-900 mb-2">
                                📝 First time? Create an admin account:
                            </p>
                            <div className="text-xs text-blue-700 space-y-1">
                                <p>1. Make sure your backend is running (port 8080)</p>
                                <p>2. Run: <code className="bg-blue-100 px-1 rounded">bash scripts/create-admin.sh</code></p>
                                <p>3. Then login with:</p>
                                <p className="mt-2"><strong>Email:</strong> admin@example.com</p>
                                <p><strong>Password:</strong> Admin123!</p>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Footer */}
                <p className="text-center text-sm text-gray-600 mt-6">
                    Go Backend Template Admin Console
                </p>
            </div>
        </div>
    );
}

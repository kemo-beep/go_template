'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/store';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Shield, Database, Users, Activity } from 'lucide-react';

export default function Home() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();

  useEffect(() => {
    if (isAuthenticated) {
      router.push('/dashboard');
    }
  }, [isAuthenticated, router]);

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white">
      <div className="container mx-auto px-4 py-16">
        {/* Hero Section */}
        <div className="text-center space-y-4 mb-16">
          <div className="flex justify-center mb-6">
            <div className="p-4 bg-blue-600 rounded-2xl">
              <Shield className="h-12 w-12 text-white" />
            </div>
          </div>
          <h1 className="text-5xl font-bold tracking-tight">Admin Console</h1>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            A powerful developer console for managing your Go backend infrastructure
          </p>
          <div className="flex gap-4 justify-center mt-8">
            <Button size="lg" onClick={() => router.push('/login')}>
              Get Started
            </Button>
            <Button size="lg" variant="outline" onClick={() => router.push('/dashboard')}>
              View Dashboard
            </Button>
          </div>
        </div>

        {/* Features */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-16">
          <Card>
            <CardHeader>
              <div className="p-3 bg-blue-100 rounded-lg w-fit mb-3">
                <Users className="h-6 w-6 text-blue-600" />
              </div>
              <CardTitle>User Management</CardTitle>
              <CardDescription>
                Create, update, and manage user accounts with full CRUD operations
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ul className="text-sm text-gray-600 space-y-2">
                <li>• Create and delete users</li>
                <li>• Reset passwords</li>
                <li>• Manage roles and permissions</li>
                <li>• View JWT sessions</li>
              </ul>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <div className="p-3 bg-green-100 rounded-lg w-fit mb-3">
                <Database className="h-6 w-6 text-green-600" />
              </div>
              <CardTitle>Database Explorer</CardTitle>
              <CardDescription>
                Browse and query your PostgreSQL database
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ul className="text-sm text-gray-600 space-y-2">
                <li>• List tables and rows</li>
                <li>• Execute custom SQL queries</li>
                <li>• View schema relationships</li>
                <li>• Export data</li>
              </ul>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <div className="p-3 bg-purple-100 rounded-lg w-fit mb-3">
                <Activity className="h-6 w-6 text-purple-600" />
              </div>
              <CardTitle>Monitoring & Logs</CardTitle>
              <CardDescription>
                Monitor your API performance and view logs
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ul className="text-sm text-gray-600 space-y-2">
                <li>• API request logs</li>
                <li>• Real-time metrics</li>
                <li>• Error tracking</li>
                <li>• Performance analytics</li>
              </ul>
            </CardContent>
          </Card>
        </div>

        {/* Tech Stack */}
        <div className="bg-white rounded-lg p-8 shadow-sm border">
          <h2 className="text-2xl font-bold mb-6 text-center">Built With Modern Technologies</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6 text-center">
            <div>
              <div className="font-semibold">Next.js 14</div>
              <div className="text-sm text-gray-600">React Framework</div>
            </div>
            <div>
              <div className="font-semibold">Shadcn UI</div>
              <div className="text-sm text-gray-600">Component Library</div>
            </div>
            <div>
              <div className="font-semibold">React Query</div>
              <div className="text-sm text-gray-600">Data Fetching</div>
            </div>
            <div>
              <div className="font-semibold">Zustand</div>
              <div className="text-sm text-gray-600">State Management</div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="text-center mt-16 text-gray-600">
          <p>© 2025 Go Backend Template. Built with ❤️ for developers.</p>
        </div>
      </div>
    </div>
  );
}
'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
    Users,
    Database,
    FolderOpen,
    Activity,
    Settings,
    Menu,
    X,
    LogOut,
    Shield,
    BarChart3,
    Wrench,
    LayoutDashboard,
    Wifi,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useAuthStore } from '@/lib/store';
import { cn } from '@/lib/utils';

const navigation = [
    { name: 'Overview', href: '/dashboard/overview', icon: LayoutDashboard },
    { name: 'Users', href: '/dashboard/users', icon: Users },
    { name: 'Database', href: '/dashboard/database', icon: Database },
    { name: 'Real-time', href: '/dashboard/realtime', icon: Wifi },
    { name: 'Storage', href: '/dashboard/storage', icon: FolderOpen },
    { name: 'Logs', href: '/dashboard/logs', icon: Activity },
    { name: 'Metrics', href: '/dashboard/metrics', icon: BarChart3 },
    { name: 'Dev Tools', href: '/dashboard/dev-tools', icon: Wrench },
    { name: 'Settings', href: '/dashboard/settings', icon: Settings },
];

export function DashboardLayout({ children }: { children: React.ReactNode }) {
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isHovered, setIsHovered] = useState(false);
    const pathname = usePathname();
    const { user, logout } = useAuthStore();

    const isExpanded = sidebarOpen || isHovered;

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Sidebar */}
            <aside
                className={cn(
                    'fixed inset-y-0 left-0 z-50 bg-white border-r border-gray-200 transition-all duration-300 group',
                    isExpanded ? 'w-64' : 'w-16'
                )}
                onMouseEnter={() => setIsHovered(true)}
                onMouseLeave={() => setIsHovered(false)}
            >
                <div className="flex h-full flex-col">
                    {/* Logo */}
                    <div className="flex h-16 items-center gap-2 border-b px-6">
                        <Shield className="h-6 w-6 text-blue-600 flex-shrink-0" />
                        <span className={cn(
                            "font-semibold text-lg transition-opacity duration-300",
                            isExpanded ? "opacity-100" : "opacity-0 w-0 overflow-hidden"
                        )}>
                            Admin Console
                        </span>
                    </div>

                    {/* Navigation */}
                    <nav className="flex-1 space-y-1 px-3 py-4">
                        {navigation.map((item) => {
                            const isActive = pathname === item.href;
                            return (
                                <Link
                                    key={item.name}
                                    href={item.href}
                                    className={cn(
                                        'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-300',
                                        isActive
                                            ? 'bg-blue-50 text-blue-600'
                                            : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                                    )}
                                    title={!isExpanded ? item.name : undefined}
                                >
                                    <item.icon className="h-5 w-5 flex-shrink-0" />
                                    <span className={cn(
                                        "transition-opacity duration-300",
                                        isExpanded ? "opacity-100" : "opacity-0 w-0 overflow-hidden"
                                    )}>
                                        {item.name}
                                    </span>
                                </Link>
                            );
                        })}
                    </nav>

                    {/* User section */}
                    <div className="border-t p-4">
                        <div className="flex items-center gap-3">
                            <Avatar>
                                <AvatarFallback>
                                    {user?.name?.substring(0, 2).toUpperCase() || 'AD'}
                                </AvatarFallback>
                            </Avatar>
                            <div className={cn(
                                "flex-1 min-w-0 transition-opacity duration-300",
                                isExpanded ? "opacity-100" : "opacity-0 w-0 overflow-hidden"
                            )}>
                                <p className="text-sm font-medium truncate">{user?.name || 'Admin'}</p>
                                <p className="text-xs text-gray-500 truncate">{user?.email}</p>
                            </div>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={logout}
                                title="Logout"
                                className={cn(
                                    "transition-opacity duration-300",
                                    isExpanded ? "opacity-100" : "opacity-0 w-0 overflow-hidden"
                                )}
                            >
                                <LogOut className="h-4 w-4" />
                            </Button>
                        </div>
                    </div>
                </div>
            </aside>

            {/* Main content */}
            <div
                className={cn(
                    'transition-all duration-300',
                    isExpanded ? 'ml-64' : 'ml-16'
                )}
            >
                {/* Top bar */}
                <header className="sticky top-0 z-40 flex h-16 items-center gap-4 border-b bg-white px-6">
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setSidebarOpen(!sidebarOpen)}
                        title={sidebarOpen ? "Collapse sidebar" : "Expand sidebar"}
                    >
                        {sidebarOpen ? (
                            <X className="h-5 w-5" />
                        ) : (
                            <Menu className="h-5 w-5" />
                        )}
                    </Button>
                    <div className="flex-1" />
                </header>

                {/* Page content */}
                <main className="p-6">{children}</main>
            </div>
        </div>
    );
}

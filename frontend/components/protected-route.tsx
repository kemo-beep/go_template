'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/store';
import { Loader2 } from 'lucide-react';

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const { isAuthenticated, isHydrated } = useAuthStore();

    useEffect(() => {
        // Only redirect to login if we're hydrated and not authenticated
        if (isHydrated && !isAuthenticated) {
            router.push('/login');
        }
    }, [isAuthenticated, isHydrated, router]);

    // Show loading while hydrating or if not authenticated after hydration
    if (!isHydrated || !isAuthenticated) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-600" />
                    <p className="text-gray-600">
                        {!isHydrated ? 'Loading...' : 'Checking authentication...'}
                    </p>
                </div>
            </div>
        );
    }

    return <>{children}</>;
}

'use client';

import { useEffect } from 'react';
import { X, CheckCircle, AlertCircle, Info, AlertTriangle } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface ToastProps {
    id: string;
    title: string;
    description?: string;
    variant?: 'default' | 'destructive' | 'success' | 'warning';
    duration?: number;
    onDismiss: (id: string) => void;
}

const variantStyles = {
    default: 'bg-background border-border text-foreground',
    destructive: 'bg-destructive text-destructive-foreground border-destructive',
    success: 'bg-green-500 text-white border-green-500',
    warning: 'bg-yellow-500 text-white border-yellow-500',
};

const variantIcons = {
    default: Info,
    destructive: AlertCircle,
    success: CheckCircle,
    warning: AlertTriangle,
};

export function Toast({ id, title, description, variant = 'default', onDismiss }: ToastProps) {
    useEffect(() => {
        const timer = setTimeout(() => {
            onDismiss(id);
        }, 5000);

        return () => clearTimeout(timer);
    }, [id, onDismiss]);

    const Icon = variantIcons[variant];

    return (
        <div
            className={cn(
                'pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg border shadow-lg',
                variantStyles[variant]
            )}
        >
            <div className="p-4">
                <div className="flex items-start">
                    <div className="flex-shrink-0">
                        <Icon className="h-5 w-5" />
                    </div>
                    <div className="ml-3 w-0 flex-1">
                        <p className="text-sm font-medium">{title}</p>
                        {description && (
                            <p className="mt-1 text-sm opacity-90">{description}</p>
                        )}
                    </div>
                    <div className="ml-4 flex flex-shrink-0">
                        <button
                            className="inline-flex rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2"
                            onClick={() => onDismiss(id)}
                        >
                            <span className="sr-only">Close</span>
                            <X className="h-5 w-5" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export function ToastContainer({ toasts, onDismiss }: { toasts: any[]; onDismiss: (id: string) => void }) {
    if (toasts.length === 0) return null;

    return (
        <div className="fixed top-4 right-4 z-50 space-y-2">
            {toasts.map((toast) => (
                <Toast
                    key={toast.id}
                    {...toast}
                    onDismiss={onDismiss}
                />
            ))}
        </div>
    );
}

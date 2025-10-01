import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from './api-client';

interface User {
    id: string;
    email: string;
    name: string;
    is_admin: boolean;
}

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isHydrated: boolean;
    login: (user: User, token: string) => void;
    logout: () => void;
    setUser: (user: User) => void;
    setHydrated: (hydrated: boolean) => void;
    validateToken: () => Promise<boolean>;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set, get) => ({
            user: null,
            token: null,
            isAuthenticated: false,
            isHydrated: false,
            login: (user, token) => {
                localStorage.setItem('auth_token', token);
                set({ user, token, isAuthenticated: true });
            },
            logout: () => {
                localStorage.removeItem('auth_token');
                set({ user: null, token: null, isAuthenticated: false });
            },
            setUser: (user) => set({ user }),
            setHydrated: (hydrated) => set({ isHydrated: hydrated }),
            validateToken: async () => {
                try {
                    const response = await api.validateToken();
                    if (response.data.success) {
                        const userData = response.data.data;
                        set({
                            user: {
                                id: userData.user_id.toString(),
                                email: userData.user_email as string,
                                name: userData.user_email as string, // We don't have name in validate response
                                is_admin: userData.is_admin as boolean,
                            },
                            isAuthenticated: true,
                        });
                        return true;
                    }
                } catch (error) {
                    console.error('Token validation failed:', error);
                    set({ user: null, token: null, isAuthenticated: false });
                    localStorage.removeItem('auth_token');
                }
                return false;
            },
        }),
        {
            name: 'auth-storage',
            onRehydrateStorage: () => (state) => {
                // Check if there's a token in localStorage and set authenticated state
                const token = localStorage.getItem('auth_token');
                if (token && state) {
                    state.isAuthenticated = true;
                }
                state?.setHydrated(true);
            },
        }
    )
);

interface UIState {
    sidebarOpen: boolean;
    theme: 'light' | 'dark';
    toggleSidebar: () => void;
    setSidebarOpen: (open: boolean) => void;
    setTheme: (theme: 'light' | 'dark') => void;
}

export const useUIStore = create<UIState>((set) => ({
    sidebarOpen: true,
    theme: 'light',
    toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
    setSidebarOpen: (open) => set({ sidebarOpen: open }),
    setTheme: (theme) => set({ theme }),
}));

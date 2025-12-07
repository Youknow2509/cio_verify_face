import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface SystemAdmin {
    id: string;
    email: string;
    full_name: string;
    role: 'super_admin' | 'admin';
    avatar?: string;
}

interface AuthState {
    user: SystemAdmin | null;
    accessToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
    setAuth: (user: SystemAdmin, accessToken: string, refreshToken: string) => void;
    setUser: (user: SystemAdmin) => void;
    clearAuth: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            accessToken: null,
            refreshToken: null,
            isAuthenticated: false,
            setAuth: (user, accessToken, refreshToken) =>
                set({
                    user,
                    accessToken,
                    refreshToken,
                    isAuthenticated: true,
                }),
            setUser: (user) => set({ user }),
            clearAuth: () =>
                set({
                    user: null,
                    accessToken: null,
                    refreshToken: null,
                    isAuthenticated: false,
                }),
        }),
        {
            name: 'system-admin-auth',
            partialize: (state) => ({
                accessToken: state.accessToken,
                refreshToken: state.refreshToken,
            }),
        }
    )
);

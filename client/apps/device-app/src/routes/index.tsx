import { createBrowserRouter, Navigate } from 'react-router-dom';
import { DeviceLayout } from '@/components/layouts/DeviceLayout';
import { TokenAuthPage } from '@/features/auth/TokenAuthPage';
import { AttendancePage } from '@/features/attendance/AttendancePage';
import { ProtectedRoute } from '@/components/ProtectedRoute';

export const router = createBrowserRouter([
    {
        path: '/',
        element: <Navigate to="/token-auth" replace />,
    },
    {
        path: '/token-auth',
        element: (
            <DeviceLayout>
                <TokenAuthPage />
            </DeviceLayout>
        ),
    },
    {
        path: '/attendance',
        element: (
            <ProtectedRoute>
                <AttendancePage />
            </ProtectedRoute>
        ),
    },
]);

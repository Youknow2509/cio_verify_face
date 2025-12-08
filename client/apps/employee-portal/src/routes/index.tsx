import { createBrowserRouter, Navigate } from 'react-router-dom';
import { LoginPage } from '@/features/auth/LoginPage';
import { DashboardPage } from '@/features/dashboard/DashboardPage';
import { AttendancePage } from '@/features/attendance/AttendancePage';
import { DailySummaryPage } from '@/features/attendance/DailySummaryPage';
import { ShiftsPage } from '@/features/shifts/ShiftsPage';
import { ProfilePage } from '@/features/profile/ProfilePage';
import { ExportPage } from '@/features/attendance/ExportPage';
import { MainLayout } from '@/components/layouts/MainLayout';
import { ProtectedRoute } from '@/components/ProtectedRoute';

export const router = createBrowserRouter([
    {
        path: '/',
        element: <Navigate to="/dashboard" replace />,
    },
    {
        path: '/login',
        element: <LoginPage />,
    },
    {
        path: '/',
        element: (
            <ProtectedRoute>
                <MainLayout />
            </ProtectedRoute>
        ),
        children: [
            {
                path: 'dashboard',
                element: <DashboardPage />,
            },
            {
                path: 'attendance',
                element: <AttendancePage />,
            },
            {
                path: 'daily-summary',
                element: <DailySummaryPage />,
            },
            {
                path: 'shifts',
                element: <ShiftsPage />,
            },
            {
                path: 'profile',
                element: <ProfilePage />,
            },
            {
                path: 'export',
                element: <ExportPage />,
            },
        ],
    },
]);

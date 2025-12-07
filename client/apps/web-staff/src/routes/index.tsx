import { createBrowserRouter, Navigate } from 'react-router-dom';
import { StaffLayout } from '@/components/layouts/StaffLayout';
import { AuthLayout } from '@/components/layouts/AuthLayout';
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { LoginPage } from '@/features/auth/LoginPage';
import { DashboardPage } from '@/features/dashboard/DashboardPage';
import { MonthlyReportPage } from '@/features/reports/MonthlyReportPage';
import { FaceUpdatePage } from '@/features/face-update/FaceUpdatePage';
import { ProfilePage } from '@/features/profile/ProfilePage';

export const router = createBrowserRouter([
    {
        path: '/',
        element: <Navigate to="/dashboard" replace />,
    },
    {
        path: '/login',
        element: (
            <AuthLayout>
                <LoginPage />
            </AuthLayout>
        ),
    },
    {
        element: <ProtectedRoute />,
        children: [
            {
                element: <StaffLayout />,
                children: [
                    {
                        path: '/dashboard',
                        element: <DashboardPage />,
                    },
                    {
                        path: '/reports/monthly',
                        element: <MonthlyReportPage />,
                    },
                    {
                        path: '/face-update',
                        element: <FaceUpdatePage />,
                    },
                    {
                        path: '/profile',
                        element: <ProfilePage />,
                    },
                ],
            },
        ],
    },
]);

import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AdminLayout } from '@/components/layouts/AdminLayout';
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { LoginPage } from '@/features/auth/LoginPage';
import { DashboardPage } from '@/features/dashboard/DashboardPage';
import { CompanyListPage } from '@/features/companies/CompanyListPage';
import { CompanyDetailPage } from '@/features/companies/CompanyDetailPage';
import { MonitoringPage } from '@/features/monitoring/MonitoringPage';
import { SettingsPage } from '@/features/settings/SettingsPage';

import { AuditLogPage } from '@/features/audit/AuditLogPage';
import { ServicePlansPage } from '@/features/settings/ServicePlansPage';
import { SecuritySettingsPage } from '@/features/settings/SecuritySettingsPage';
import { NotificationSettingsPage } from '@/features/settings/NotificationSettingsPage';
import { StorageSettingsPage } from '@/features/settings/StorageSettingsPage';

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
        element: <ProtectedRoute />,
        children: [
            {
                element: <AdminLayout />,
                children: [
                    { path: '/dashboard', element: <DashboardPage /> },
                    { path: '/companies', element: <CompanyListPage /> },
                    { path: '/companies/:id', element: <CompanyDetailPage /> },
                    { path: '/monitoring', element: <MonitoringPage /> },
                    { path: '/audit-log', element: <AuditLogPage /> },
                    { path: '/settings', element: <SettingsPage /> },
                    { path: '/settings/plans', element: <ServicePlansPage /> },
                    { path: '/settings/security', element: <SecuritySettingsPage /> },
                    { path: '/settings/notifications', element: <NotificationSettingsPage /> },
                    { path: '/settings/storage', element: <StorageSettingsPage /> },
                ],
            },
        ],
    },
]);


import { createBrowserRouter, Navigate } from 'react-router-dom';
import { MainLayout } from '@/components/layouts/MainLayout';
import { AuthLayout } from '@/components/layouts/AuthLayout';
import { LoginPage } from '@/features/auth/LoginPage';
import { DashboardPage } from '@/features/dashboard/DashboardPage';
import { EmployeeListPage } from '@/features/employees/EmployeeListPage';
import { EmployeeFormPage } from '@/features/employees/EmployeeFormPage';
import { EmployeeFaceDataPage } from '@/features/employees/EmployeeFaceDataPage';
import { DeviceListPage } from '@/features/devices/DeviceListPage';
import { DeviceFormPage } from '@/features/devices/DeviceFormPage';
import { DeviceConfigPage } from '@/features/devices/DeviceConfigPage';
import { ShiftListPage } from '@/features/shifts/ShiftListPage';
import { ShiftFormPage } from '@/features/shifts/ShiftFormPage';
import { DailyReportPage } from '@/features/reports/DailyReportPage';
import { SummaryReportPage } from '@/features/reports/SummaryReportPage';
import { SettingsPage } from '@/features/settings/SettingsPage';
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { DeviceAddPage } from '@/features/devices/DeviceAdd';

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
        path: 'employees',
        element: <EmployeeListPage />,
      },
      {
        path: 'employees/add',
        element: <EmployeeFormPage />,
      },
      {
        path: 'employees/:id/edit',
        element: <EmployeeFormPage />,
      },
      {
        path: 'employees/:id/face-data',
        element: <EmployeeFaceDataPage />,
      },
      {
        path: 'devices',
        element: <DeviceListPage />,
      },
      {
        path: 'devices/add',
        element: <DeviceAddPage />,
      },
      {
        path: 'devices/:id/edit',
        element: <DeviceFormPage />,
      },
      {
        path: 'devices/:id/config',
        element: <DeviceConfigPage />,
      },
      {
        path: 'shifts',
        element: <ShiftListPage />,
      },
      {
        path: 'shifts/add',
        element: <ShiftFormPage />,
      },
      {
        path: 'shifts/:id/edit',
        element: <ShiftFormPage />,
      },
      {
        path: 'reports/daily',
        element: <DailyReportPage />,
      },
      {
        path: 'reports/summary',
        element: <SummaryReportPage />,
      },
      {
        path: 'settings',
        element: <SettingsPage />,
      },
    ],
  },
]);

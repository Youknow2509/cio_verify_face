// src/app/routes.tsx

import { lazy, Suspense } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { Layout } from '@/app/layouts/Layout/Layout';
const Login = lazy(() => import('@/features/auth/Login'));

const Dashboard = lazy(() => import('@/features/dashboard/Dashboard'));
const EmployeesPage = lazy(() => import('@/features/employees/EmployeesPage'));
const EmployeeDetail = lazy(() => import('@/features/employees/EmployeeDetail'));
const EmployeeEditPage = lazy(() => import('@/features/employees/EmployeeEditPage'));
const Devices = lazy(() => import('@/features/devices/Devices'));
const Attendance = lazy(() => import('@/features/attendance/Attendance'));
const Reports = lazy(() => import('@/features/reports/Reports'));
const Shifts = lazy(() => import('@/features/shifts/Shifts'));
const Settings = lazy(() => import('@/features/settings/Settings'));
const AccountProfilePage = lazy(() => import('@/features/account/AccountProfilePage'));
const ChangePasswordPage = lazy(() => import('@/features/account/ChangePasswordPage'));
const FaceRegistration = lazy(() => import('@/features/face-registration/FaceRegistration'));

function LoadingFallback() {
  return (
    <div style={{ 
      display: 'flex', 
      alignItems: 'center', 
      justifyContent: 'center', 
      height: '200px',
      fontSize: '14px',
      color: 'var(--color-text-secondary)'
    }}>
      Đang tải...
    </div>
  );
}

export function AppRoutes() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <Routes>
        <Route path="login" element={<Login />} />
        <Route element={<Layout />}>
          <Route index element={<Navigate to="/login" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="employees" element={<EmployeesPage />} />
          <Route path="employees/:id" element={<EmployeeDetail />} />
          <Route path="employees/:id/edit" element={<EmployeeEditPage />} />
          <Route path="devices" element={<Devices />} />
          <Route path="attendance" element={<Attendance />} />
          <Route path="reports" element={<Reports />} />
          <Route path="shifts" element={<Shifts />} />
          <Route path="face-registration" element={<FaceRegistration />} />
          <Route path="settings" element={<Settings />} />
          <Route path="account" element={<AccountProfilePage />} />
          <Route path="account/password" element={<ChangePasswordPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Suspense>
  );
}
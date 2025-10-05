// src/app/routes.tsx

import { Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from '../layouts/Layout/Layout';

// Lazy load pages for better performance
import { lazy, Suspense } from 'react';

const Dashboard = lazy(() => import('../pages/Dashboard/Dashboard'));
const Employees = lazy(() => import('../pages/Employees/Employees'));
const EmployeeDetail = lazy(() => import('../pages/Employees/EmployeeDetail'));
const Devices = lazy(() => import('../pages/Devices/Devices'));
const Attendance = lazy(() => import('../pages/Attendance/Attendance'));
const Reports = lazy(() => import('../pages/Reports/Reports'));
const Shifts = lazy(() => import('../pages/Shifts/Shifts'));
const Settings = lazy(() => import('../pages/Settings/Settings'));

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
    <Layout>
      <Suspense fallback={<LoadingFallback />}>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/employees" element={<Employees />} />
          <Route path="/employees/:id" element={<EmployeeDetail />} />
          <Route path="/devices" element={<Devices />} />
          <Route path="/attendance" element={<Attendance />} />
          <Route path="/reports" element={<Reports />} />
          <Route path="/shifts" element={<Shifts />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </Suspense>
    </Layout>
  );
}
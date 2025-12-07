import axios from 'axios';

// API Gateway will route to appropriate services
const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || '/api/v1';

export const api = axios.create({
    baseURL: API_BASE_URL,
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor - add system admin token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('system_admin_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Response interceptor
api.interceptors.response.use(
    (response) => response.data,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('system_admin_token');
            window.location.href = '/login';
        }
        return Promise.reject(error.response?.data || error);
    }
);

// ============================================
// AUTH SERVICE (service_auth)
// ============================================
export const authApi = {
    login: (credentials: { email: string; password: string }) =>
        api.post('/auth/login/admin', credentials),
    refresh: () => api.post('/auth/refresh'),
    logout: () => api.post('/auth/logout'),
    getMe: () => api.get('/auth/me'),
};

// ============================================
// IDENTITY SERVICE - Companies (service_identity)
// ============================================
export const companiesApi = {
    getAll: () => api.get('/companies'),
    create: (data: any) => api.post('/companies', data),
    getById: (id: string) => api.get(`/companies/${id}`),
    update: (id: string, data: any) => api.put(`/companies/${id}`, data),
    delete: (id: string) => api.delete(`/companies/${id}`),
};

// ============================================
// IDENTITY SERVICE - Users (service_identity)
// ============================================
export const usersApi = {
    getAll: () => api.get('/users'),
    create: (data: any) => api.post('/users', data),
    getById: (id: string) => api.get(`/users/${id}`),
    update: (id: string, data: any) => api.put(`/users/${id}`, data),
    delete: (id: string) => api.delete(`/users/${id}`),
};

// ============================================
// DEVICE SERVICE (service_device)
// ============================================
export const devicesApi = {
    getAll: () => api.get('/device'),
    create: (data: any) => api.post('/device', data),
    getById: (id: string) => api.get(`/device/${id}`),
    update: (id: string, data: any) => api.put(`/device/${id}`, data),
    delete: (id: string) => api.delete(`/device/${id}`),
    updateStatus: (data: any) => api.post('/device/status', data),
};

// ============================================
// ANALYTIC SERVICE - Reports & Audit (service_analytic)
// ============================================
export const reportsApi = {
    getDaily: () => api.get('/reports/daily'),
    getSummary: () => api.get('/reports/summary'),
    getAuditLogs: () => api.get('/audit-logs'),
    getAuditLogsRange: (params: { start_date: string; end_date: string }) =>
        api.get('/audit-logs/range', { params }),
};

// ============================================
// WS DELIVERY SERVICE - Health (service_ws_delivery)
// ============================================
export const healthApi = {
    getHealth: () => api.get('/health'),
    getHealthDetails: () => api.get('/api/health/details'),
};

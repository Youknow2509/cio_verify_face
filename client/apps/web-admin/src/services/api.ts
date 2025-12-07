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

// Request interceptor - add auth token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Response interceptor - handle errors
api.interceptors.response.use(
    (response) => response.data,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('access_token');
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
// ANALYTIC SERVICE (service_analytic)
// ============================================
export const reportsApi = {
    getDaily: () => api.get('/reports/daily'),
    getSummary: () => api.get('/reports/summary'),
    exportReport: (params: any) => api.post('/reports/export', params),
    downloadReport: (filename: string) => api.get(`/reports/download/${filename}`),

    // Attendance records
    getAttendanceRecords: () => api.get('/attendance-records'),
    getAttendanceRecordsRange: (params: { start_date: string; end_date: string }) =>
        api.get('/attendance-records/range', { params }),
    getAttendanceByEmployee: (employeeId: string) =>
        api.get(`/attendance-records/employee/${employeeId}`),

    // Daily summaries
    getDailySummaries: () => api.get('/daily-summaries'),
    getDailySummariesByUser: (employeeId: string) =>
        api.get(`/daily-summaries/user/${employeeId}`),

    // Audit logs
    getAuditLogs: () => api.get('/audit-logs'),
    getAuditLogsRange: (params: { start_date: string; end_date: string }) =>
        api.get('/audit-logs/range', { params }),
    createAuditLog: (data: any) => api.post('/audit-logs', data),

    // Company reports
    getCompanyDailyStatus: () => api.get('/company/daily-attendance-status'),
    getCompanyStatusRange: (params: { start_date: string; end_date: string }) =>
        api.get('/company/attendance-status/range', { params }),
    getCompanyMonthlySummary: () => api.get('/company/monthly-summary'),
    exportCompanyDailyStatus: (params: any) => api.post('/company/export-daily-status', params),
    exportCompanyMonthlySummary: (params: any) => api.post('/company/export-monthly-summary', params),
};

// ============================================
// IDENTITY SERVICE (service_identity)
// ============================================
export const companiesApi = {
    getAll: () => api.get('/companies'),
    create: (data: any) => api.post('/companies', data),
    getById: (id: string) => api.get(`/companies/${id}`),
    update: (id: string, data: any) => api.put(`/companies/${id}`, data),
    delete: (id: string) => api.delete(`/companies/${id}`),
};

export const usersApi = {
    getAll: () => api.get('/users'),
    create: (data: any) => api.post('/users', data),
    getById: (id: string) => api.get(`/users/${id}`),
    update: (id: string, data: any) => api.put(`/users/${id}`, data),
    delete: (id: string) => api.delete(`/users/${id}`),

    // Face data
    uploadFaceData: (userId: string, data: FormData) =>
        api.post(`/users/${userId}/face-data/upload`, data, {
            headers: { 'Content-Type': 'multipart/form-data' },
        }),
    getFaceData: (userId: string) => api.get(`/users/${userId}/face-data`),
    deleteFaceData: (userId: string, faceId: string) =>
        api.delete(`/users/${userId}/face-data/${faceId}`),
    setPrimaryFace: (userId: string, faceId: string) =>
        api.put(`/users/${userId}/face-data/${faceId}/primary`),
};

// ============================================
// DEVICE SERVICE (service_device)
// ============================================
export const devicesApi = {
    getAll: () => api.get('/device'),
    create: (data: any) => api.post('/device', data),
    getById: (id: string) => api.get(`/device/${id}`),
    getToken: (id: string) => api.get(`/device/token/${id}`),
    refreshToken: (id: string) => api.post(`/device/token/refresh/${id}`),
    update: (id: string, data: any) => api.put(`/device/${id}`, data),
    delete: (id: string) => api.delete(`/device/${id}`),
    updateLocation: (data: any) => api.post('/device/location', data),
    updateName: (data: any) => api.post('/device/name', data),
    updateInfo: (data: any) => api.post('/device/info', data),
    updateStatus: (data: any) => api.post('/device/status', data),
};

// ============================================
// WORKFORCE SERVICE (service_workforce)
// ============================================
export const shiftsApi = {
    getAll: () => api.get('/shift'),
    create: (data: any) => api.post('/shift', data),
    getById: (id: string) => api.get(`/shift/${id}`),
    edit: (data: any) => api.post('/shift/edit', data),
    delete: (id: string) => api.delete(`/shift/${id}`),
    updateStatus: (data: any) => api.post('/shift/status', data),
};

export const employeeShiftsApi = {
    get: (data: any) => api.post('/employee/shift', data),
    editEffective: (data: any) => api.post('/employee/shift/edit/effective', data),
    enable: (data: any) => api.post('/employee/shift/enable', data),
    disable: (data: any) => api.post('/employee/shift/disable', data),
    delete: (data: any) => api.post('/employee/shift/delete', data),
    add: (data: any) => api.post('/employee/shift/add', data),
    addList: (data: any) => api.post('/employee/shift/add/list', data),
    getNotIn: (data: any) => api.post('/employee/shift/not_in', data),
    getIn: (data: any) => api.post('/employee/shift/in', data),
};

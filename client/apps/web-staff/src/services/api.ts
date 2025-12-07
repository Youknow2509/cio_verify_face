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
        const token = localStorage.getItem('staff_access_token');
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
            localStorage.removeItem('staff_access_token');
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
        api.post('/auth/login', credentials),
    refresh: () => api.post('/auth/refresh'),
    logout: () => api.post('/auth/logout'),
    getMe: () => api.get('/auth/me'),
};

// ============================================
// ANALYTIC SERVICE - Employee Self-Service (service_analytic)
// ============================================
export const employeeReportsApi = {
    // My attendance
    getMyAttendanceRecords: () => api.get('/employee/my-attendance-records'),
    getMyAttendanceRecordsRange: (params: { start_date: string; end_date: string }) =>
        api.get('/employee/my-attendance-records/range', { params }),

    // Daily summaries
    getMyDailySummaries: () => api.get('/employee/my-daily-summaries'),
    getMyDailySummary: (date: string) => api.get(`/employee/my-daily-summary/${date}`),

    // Stats
    getMyStats: () => api.get('/employee/my-stats'),
    getMyDailyStatus: () => api.get('/employee/my-daily-status'),
    getMyStatusRange: (params: { start_date: string; end_date: string }) =>
        api.get('/employee/my-status/range', { params }),
    getMyMonthlySummary: () => api.get('/employee/my-monthly-summary'),

    // Export
    exportMyDailyStatus: (params: any) => api.post('/employee/export-daily-status', params),
    exportMyMonthlySummary: (params: any) => api.post('/employee/export-monthly-summary', params),
};

// ============================================
// PROFILE UPDATE SERVICE (service_profile_update)
// ============================================
export const profileApi = {
    // Face update requests
    createUpdateRequest: (data: any) => api.post('/profile-update/requests', data),
    getMyRequests: () => api.get('/profile-update/requests/me'),
    validateToken: () => api.get('/profile-update/token/validate'),
    uploadFace: (data: FormData) =>
        api.post('/profile-update/face', data, {
            headers: { 'Content-Type': 'multipart/form-data' },
        }),
    resetPassword: (data: any) => api.post('/password/reset', data),
};

// ============================================
// IDENTITY SERVICE - Get own user data (service_identity)
// ============================================
export const userApi = {
    getMe: () => api.get('/auth/me'),
    getFaceData: (userId: string) => api.get(`/users/${userId}/face-data`),
};

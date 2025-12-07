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

// Request interceptor - add device token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('device_token');
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
            localStorage.removeItem('device_token');
            window.location.href = '/token-auth';
        }
        return Promise.reject(error.response?.data || error);
    }
);

// ============================================
// AUTH SERVICE - Device Auth (service_auth)
// ============================================
export const deviceAuthApi = {
    authenticate: (token: string) => api.post('/auth/device', { token }),
    deactivate: () => api.delete('/auth/device'),
};

// ============================================
// AI SERVICE - Face Verification (service_ai)
// ============================================
export const faceApi = {
    verify: (data: FormData) =>
        api.post('/face/verify/upload', data, {
            headers: { 'Content-Type': 'multipart/form-data' },
        }),
    enroll: (data: FormData) =>
        api.post('/face/enroll/upload', data, {
            headers: { 'Content-Type': 'multipart/form-data' },
        }),
};

// ============================================
// ATTENDANCE SERVICE (service_attendance)
// ============================================
export const attendanceApi = {
    checkIn: (data: any) => api.post('/attendance/', data),
    createRecord: (data: any) => api.post('/attendance/records', data),
    getDailySummary: (data: any) => api.post('/attendance/records/summary/daily', data),
    getEmployeeRecords: (data: any) => api.post('/attendance/records/employee', data),
    getEmployeeDailySummary: (data: any) => api.post('/attendance/records/employee/summary/daily', data),
};

// ============================================
// DEVICE SERVICE (service_device)
// ============================================
export const deviceApi = {
    getInfo: () => api.get('/device'),
    updateLocation: (data: any) => api.post('/device/location', data),
    updateStatus: (data: any) => api.post('/device/status', data),
};

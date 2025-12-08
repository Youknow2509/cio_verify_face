import axios from 'axios';

const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080';

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

// Response interceptor - handle token refresh
api.interceptors.response.use(
    (response) => response.data,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            const refreshToken = localStorage.getItem('refresh_token');
            const accessToken = localStorage.getItem('access_token');
            
            if (refreshToken && accessToken) {
                try {
                    const response = await axios.post(
                        `${API_BASE_URL}/api/v1/auth/refresh`,
                        {
                            access_token: accessToken,
                            refresh_token: refreshToken,
                        }
                    );

                    const { access_token, refresh_token } = response.data.data;
                    localStorage.setItem('access_token', access_token);
                    localStorage.setItem('refresh_token', refresh_token);
                    
                    originalRequest.headers.Authorization = `Bearer ${access_token}`;
                    return api(originalRequest);
                } catch (refreshError) {
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    window.location.href = '/login';
                    return Promise.reject(refreshError);
                }
            } else {
                localStorage.removeItem('access_token');
                localStorage.removeItem('refresh_token');
                window.location.href = '/login';
            }
        }

        return Promise.reject(error.response?.data || error);
    }
);

// Auth API
export const authApi = {
    login: (credentials: { username: string; password: string }) =>
        api.post('/api/v1/auth/login', credentials),
    refresh: (data: { access_token: string; refresh_token: string }) =>
        api.post('/api/v1/auth/refresh', data),
    getMe: () => api.get('/api/v1/auth/me'),
};

// Profile Update API
export const profileUpdateApi = {
    createRequest: (data: { reason: string }) =>
        api.post('/api/v1/profile-update/requests', data),
    getMyRequest: () =>
        api.get('/api/v1/profile-update/requests/me'),
    updateFace: (formData: FormData) =>
        api.post('/api/v1/profile-update/face', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        }),
};

// Shift API
export const shiftApi = {
    getEmployeeShifts: (params: { page: number; size: number }) =>
        api.get('/api/v1/shift/employee', { params }),
};

// Attendance API
export const attendanceApi = {
    getMyAttendanceRecords: (params: { year_month: string }) =>
        api.get('/api/v1/employee/my-attendance-records', { params }),
    getMySummaries: (params: { month: string }) =>
        api.get('/api/v1/employee/my-daily-summaries', { params }),
    exportMonthlySummary: (data: { email: string; format: string; month: string }) =>
        api.post('/api/v1/employee/export-monthly-summary', data),
};

import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

export const API_BASE_URL =
    (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8000';

export const createApiClient = (
    baseURL: string = API_BASE_URL
): AxiosInstance => {
    const client = axios.create({
        baseURL,
        headers: {
            'Content-Type': 'application/json',
        },
    });

    // Request interceptor to add auth token
    client.interceptors.request.use(
        (config) => {
            const token = localStorage.getItem('access_token');
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
            return config;
        },
        (error) => Promise.reject(error)
    );

    // Response interceptor to handle token refresh
    client.interceptors.response.use(
        (response) => response,
        async (error) => {
            const originalRequest = error.config;

            if (error.response?.status === 401 && !originalRequest._retry) {
                originalRequest._retry = true;

                const refreshToken = localStorage.getItem('refresh_token');
                if (refreshToken) {
                    try {
                        const access_token_old =
                            localStorage.getItem('access_token');
                        const headerAuth = `Bearer ${access_token_old}`;
                        const response = await axios.post(
                            `${baseURL}/api/v1/auth/refresh`,
                            {
                                access_token: access_token_old,
                                refresh_token: refreshToken,
                            },
                            {
                                headers: {
                                    Authorization: headerAuth,
                                },
                            }
                        );

                        const { access_token } = response.data;
                        if (!access_token) {
                            localStorage.setItem('access_token', access_token);
                        }
                        originalRequest.headers.Authorization = `Bearer ${access_token}`;
                        return client(originalRequest);
                    } catch (refreshError) {
                        localStorage.removeItem('access_token');
                        localStorage.removeItem('refresh_token');
                        window.location.href = '/login';
                        return Promise.reject(refreshError);
                    }
                }
            }

            return Promise.reject(error);
        }
    );

    return client;
};

export const apiClient = createApiClient();

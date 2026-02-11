import axios, { type InternalAxiosRequestConfig, type AxiosResponse, type AxiosError } from 'axios';
import router from '@/router';

const client = axios.create({
    baseURL: (import.meta.env.VITE_API_URL || 'http://localhost:8080') + '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor to add auth token
client.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error: AxiosError) => {
        return Promise.reject(error);
    }
);

// Response interceptor to unwrap backend { message, data } envelope and handle 401s
client.interceptors.response.use(
    (response: AxiosResponse) => {
        if (response.data && typeof response.data === 'object' && 'data' in response.data) {
            response.data = response.data.data;
        }
        return response;
    },
    (error: AxiosError) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            router.push('/login');
        }
        return Promise.reject(error);
    }
);

export default client;

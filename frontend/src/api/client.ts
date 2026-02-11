import axios, { type InternalAxiosRequestConfig, type AxiosResponse, type AxiosError } from 'axios';
import router from '@/router';

const client = axios.create({
    baseURL: (import.meta.env.VITE_API_URL || 'http://localhost:8080') + '/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true, // Send cookies with requests
});

// Request interceptor to add CSRF token (cookies are sent automatically with withCredentials: true)
client.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // Add CSRF token to state-changing requests
        const csrfToken = window.__CSRF_TOKEN__;
        const method = config.method?.toLowerCase();
        if (csrfToken && method && ['post', 'put', 'patch', 'delete'].includes(method)) {
            config.headers['X-CSRF-Token'] = csrfToken;
        }
        return config;
    },
    (error: AxiosError) => {
        return Promise.reject(error);
    }
);

// Response interceptor to unwrap backend { message, data } envelope and handle 401s and CSRF errors
client.interceptors.response.use(
    (response: AxiosResponse) => {
        if (response.data && typeof response.data === 'object' && 'data' in response.data) {
            response.data = response.data.data;
        }
        return response;
    },
    async (error: AxiosError) => {
        if (error.response?.status === 401) {
            // Clear user data (token is in HttpOnly cookie, cleared by backend)
            localStorage.removeItem('user');
            router.push('/login');
        }

        // Handle CSRF token errors
        if (error.response?.status === 403) {
            const errorData = error.response.data as { error?: string };
            if (errorData?.error?.includes('CSRF')) {
                try {
                    // Token expired or invalid - refetch
                    const response = await client.get('/csrf-token');
                    window.__CSRF_TOKEN__ = response.data.csrf_token;

                    // Retry original request
                    const originalRequest = error.config;
                    if (originalRequest) {
                        originalRequest.headers['X-CSRF-Token'] = window.__CSRF_TOKEN__;
                        return client(originalRequest);
                    }
                } catch (retryError) {
                    return Promise.reject(retryError);
                }
            }
        }

        return Promise.reject(error);
    }
);

export default client;

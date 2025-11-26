import { ref, type Ref } from 'vue';
import type { AppError } from '@/types/errors';

interface UseErrorHandlerReturn {
    error: Ref<AppError | null>;
    setError: (error: AppError | null) => void;
    clearError: () => void;
    getUserMessage: (error: AppError) => string;
}

export function useErrorHandler(): UseErrorHandlerReturn {
    const error = ref<AppError | null>(null);

    const setError = (err: AppError | null) => {
        error.value = err;

        // Log error for debugging
        if (err) {
            console.error('[Error Handler]', err);
        }
    };

    const clearError = () => {
        error.value = null;
    };

    const getUserMessage = (err: AppError): string => {
        // Map error codes to user-friendly messages
        const messageMap: Record<string, string> = {
            'API_ERROR': 'Something went wrong. Please try again.',
            'NETWORK_ERROR': 'Unable to connect. Please check your internet connection.',
            'VALIDATION_ERROR': err.message, // Use the specific validation message
            'AUTH_ERROR': 'Authentication failed. Please log in again.',
        };

        return messageMap[err.code] || err.message || 'An unexpected error occurred.';
    };

    return {
        error,
        setError,
        clearError,
        getUserMessage,
    };
}

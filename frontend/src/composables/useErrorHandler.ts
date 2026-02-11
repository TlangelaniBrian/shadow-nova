import { ref, type Ref } from 'vue';
import type { AppError } from '@/types/errors';
import { ERROR_MESSAGES } from '@/types/errors';

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
    };

    const clearError = () => {
        error.value = null;
    };

    const getUserMessage = (err: AppError): string => {
        if (err.code === 'VALIDATION_ERROR') {
            return err.message;
        }
        return ERROR_MESSAGES[err.code] || err.message || 'An unexpected error occurred.';
    };

    return {
        error,
        setError,
        clearError,
        getUserMessage,
    };
}

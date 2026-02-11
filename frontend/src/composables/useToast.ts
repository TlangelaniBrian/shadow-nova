import { toast } from 'vue-sonner';
import type { AppError } from '@/types/errors';
import { ERROR_MESSAGES } from '@/types/errors';

interface UseToastReturn {
    success: (message: string, description?: string) => void;
    error: (message: string, description?: string) => void;
    info: (message: string, description?: string) => void;
    warning: (message: string, description?: string) => void;
    showError: (error: AppError) => void;
}

export function useToast(): UseToastReturn {
    const success = (message: string, description?: string) => {
        toast.success(message, {
            description,
            duration: 4000,
        });
    };

    const error = (message: string, description?: string) => {
        toast.error(message, {
            description,
            duration: 5000,
        });
    };

    const info = (message: string, description?: string) => {
        toast.info(message, {
            description,
            duration: 4000,
        });
    };

    const warning = (message: string, description?: string) => {
        toast.warning(message, {
            description,
            duration: 4000,
        });
    };

    const showError = (appError: AppError) => {
        const message = appError.code === 'VALIDATION_ERROR'
            ? appError.message
            : (ERROR_MESSAGES[appError.code] || appError.message || 'An unexpected error occurred.');

        toast.error(message, {
            duration: 5000,
        });
    };

    return {
        success,
        error,
        info,
        warning,
        showError,
    };
}

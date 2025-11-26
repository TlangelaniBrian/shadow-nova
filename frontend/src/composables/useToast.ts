import { toast } from 'vue-sonner';
import type { AppError } from '@/types/errors';

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
        const errorMessages: Record<string, string> = {
            'API_ERROR': 'Something went wrong. Please try again.',
            'NETWORK_ERROR': 'Unable to connect. Please check your internet connection.',
            'VALIDATION_ERROR': appError.message,
            'AUTH_ERROR': 'Authentication failed. Please log in again.',
        };

        const message = errorMessages[appError.code] || appError.message || 'An unexpected error occurred.';

        toast.error(message, {
            description: appError.details ? String(appError.details) : undefined,
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

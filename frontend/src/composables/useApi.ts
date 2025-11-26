import { ref, type Ref } from 'vue';
import type { Result, AppError } from '@/types/errors';
import { transformAxiosError, failure } from '@/types/errors';

interface UseApiOptions<T> {
    immediate?: boolean;
    onSuccess?: (data: T) => void;
    onError?: (error: AppError) => void;
}

interface UseApiReturn<T> {
    data: Ref<T | null>;
    error: Ref<AppError | null>;
    isLoading: Ref<boolean>;
    execute: (...args: unknown[]) => Promise<Result<T>>;
}

export function useApi<T>(
    apiFunction: (...args: unknown[]) => Promise<T>,
    options: UseApiOptions<T> = {}
): UseApiReturn<T> {
    const data = ref<T | null>(null);
    const error = ref<AppError | null>(null);
    const isLoading = ref(false);

    const execute = async (...args: unknown[]): Promise<Result<T>> => {
        isLoading.value = true;
        error.value = null;

        try {
            const result = await apiFunction(...args);
            data.value = result;
            isLoading.value = false;

            if (options.onSuccess) {
                options.onSuccess(result);
            }

            return { data: result, error: null };
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            data.value = null;
            isLoading.value = false;

            if (options.onError) {
                options.onError(appError);
            }

            return failure(appError);
        }
    };

    // Execute immediately if requested
    if (options.immediate) {
        execute();
    }

    return {
        data: data as Ref<T | null>,
        error,
        isLoading,
        execute,
    };
}

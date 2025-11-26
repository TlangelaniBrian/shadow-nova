import { ref, computed, type Ref } from 'vue';
import { useUserStore } from '@/stores/user';
import { authApi, type User, type AuthResponse } from '@/api/auth';
import type { Result, AppError } from '@/types/errors';
import { transformAxiosError, success, failure } from '@/types/errors';
import { useToast } from '@/composables/useToast';

interface UseAuthReturn {
    user: Ref<User | null>;
    token: Ref<string | null>;
    isAuthenticated: Ref<boolean>;
    isLoading: Ref<boolean>;
    error: Ref<AppError | null>;
    login: (googleToken: string) => Promise<Result<AuthResponse>>;
    handleGoogleCallback: (code: string) => Promise<Result<AuthResponse>>;
    linkGitHub: (code: string) => Promise<Result<void>>;
    logout: () => void;
}

export function useAuth(): UseAuthReturn {
    const store = useUserStore();
    const isLoading = ref(false);
    const error = ref<AppError | null>(null);
    const toast = useToast();

    const login = async (googleToken: string): Promise<Result<AuthResponse>> => {
        isLoading.value = true;
        error.value = null;

        try {
            const response = await authApi.loginWithGoogle(googleToken);
            store.setSession(response.data.token, response.data.user);
            isLoading.value = false;
            toast.success('Welcome back!', `Logged in as ${response.data.user.name}`);
            return success(response.data);
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            isLoading.value = false;
            toast.showError(appError);
            return failure(appError);
        }
    };

    const handleGoogleCallback = async (code: string): Promise<Result<AuthResponse>> => {
        isLoading.value = true;
        error.value = null;

        try {
            const response = await authApi.handleGoogleCallback(code);
            store.setSession(response.data.token, response.data.user);
            isLoading.value = false;
            toast.success('Welcome!', 'Successfully logged in with Google');
            return success(response.data);
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            isLoading.value = false;
            toast.showError(appError);
            return failure(appError);
        }
    };

    const linkGitHub = async (code: string): Promise<Result<void>> => {
        isLoading.value = true;
        error.value = null;

        try {
            await authApi.linkGitHub(code);
            isLoading.value = false;
            toast.success('GitHub Linked!', 'Your GitHub account has been successfully linked');
            return success(undefined as void);
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            isLoading.value = false;
            toast.showError(appError);
            return failure(appError);
        }
    };

    const logout = () => {
        store.logout();
        error.value = null;
        toast.info('Logged out', 'You have been successfully logged out');
    };

    return {
        user: computed(() => store.user),
        token: computed(() => store.token),
        isAuthenticated: computed(() => store.isAuthenticated),
        isLoading,
        error,
        login,
        handleGoogleCallback,
        linkGitHub,
        logout,
    };
}

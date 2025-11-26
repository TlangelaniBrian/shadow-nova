// Error types for the application

export interface AppError {
    code: string;
    message: string;
    details?: unknown;
    timestamp: Date;
}

export class ApiError implements AppError {
    code: string;
    message: string;
    details?: unknown;
    timestamp: Date;
    statusCode?: number;

    constructor(message: string, statusCode?: number, details?: unknown) {
        this.code = 'API_ERROR';
        this.message = message;
        this.statusCode = statusCode;
        this.details = details;
        this.timestamp = new Date();
    }
}

export class ValidationError implements AppError {
    code: string;
    message: string;
    details?: unknown;
    timestamp: Date;
    field?: string;

    constructor(message: string, field?: string, details?: unknown) {
        this.code = 'VALIDATION_ERROR';
        this.message = message;
        this.field = field;
        this.details = details;
        this.timestamp = new Date();
    }
}

export class NetworkError implements AppError {
    code: string;
    message: string;
    details?: unknown;
    timestamp: Date;

    constructor(message: string = 'Network request failed', details?: unknown) {
        this.code = 'NETWORK_ERROR';
        this.message = message;
        this.details = details;
        this.timestamp = new Date();
    }
}

export class AuthenticationError implements AppError {
    code: string;
    message: string;
    details?: unknown;
    timestamp: Date;

    constructor(message: string = 'Authentication failed', details?: unknown) {
        this.code = 'AUTH_ERROR';
        this.message = message;
        this.details = details;
        this.timestamp = new Date();
    }
}

// Result type for operations that can fail
export type Result<T> = {
    data: T | null;
    error: AppError | null;
};

// Helper to create success result
export function success<T>(data: T): Result<T> {
    return { data, error: null };
}

// Helper to create error result
export function failure<T>(error: AppError): Result<T> {
    return { data: null, error };
}

// Transform axios error to AppError
export function transformAxiosError(error: unknown): AppError {
    if (typeof error === 'object' && error !== null && 'response' in error) {
        const axiosError = error as { response?: { status?: number; data?: { message?: string } }; message?: string };
        const status = axiosError.response?.status;
        const message = axiosError.response?.data?.message || axiosError.message || 'An error occurred';

        if (status === 401 || status === 403) {
            return new AuthenticationError(message, axiosError);
        }

        if (status === 400) {
            return new ValidationError(message, undefined, axiosError);
        }

        return new ApiError(message, status, axiosError);
    }

    if (error instanceof Error) {
        return new NetworkError(error.message, error);
    }

    return new NetworkError('An unknown error occurred', error);
}

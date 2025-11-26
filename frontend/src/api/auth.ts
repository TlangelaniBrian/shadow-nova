import client from './client';

export interface User {
    id: string;
    email: string;
    name: string;
    picture?: string;
    github_username?: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}

export const authApi = {
    loginWithGoogle(token: string) {
        return client.post<AuthResponse>('/auth/google', { token });
    },

    handleGoogleCallback(code: string) {
        return client.get<AuthResponse>(`/auth/google/callback?code=${code}`);
    },

    linkGitHub(code: string) {
        return client.get(`/auth/github/callback?code=${code}`);
    },
};

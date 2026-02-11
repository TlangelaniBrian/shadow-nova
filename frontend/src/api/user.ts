import client from './client';

export interface UserProfile {
    id: string;
    email: string;
    username?: string;
    name: string;
    picture?: string;
    github_username?: string;
    role?: string;
}

export interface UpdateProfileData {
    username?: string;
    email?: string;
}

export interface UpdatePasswordData {
    current_password: string;
    new_password: string;
}

export const userApi = {
    getUserProfile() {
        return client.get<UserProfile>('/user/profile');
    },

    updateUserProfile(data: UpdateProfileData) {
        return client.patch<UserProfile>('/user/profile', data);
    },

    updatePassword(data: UpdatePasswordData) {
        return client.put('/user/password', data);
    },
};

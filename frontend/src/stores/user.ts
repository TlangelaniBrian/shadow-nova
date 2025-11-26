import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authApi, type User } from '@/api/auth';

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('token'));
  const isAuthenticated = ref(!!token.value);

  async function loginWithGoogle(googleToken: string) {
    try {
      const response = await authApi.loginWithGoogle(googleToken);
      setSession(response.data.token, response.data.user);
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  }

  async function handleGoogleCallback(code: string) {
    try {
      const response = await authApi.handleGoogleCallback(code);
      setSession(response.data.token, response.data.user);
    } catch (error) {
      console.error('Google callback failed:', error);
      throw error;
    }
  }

  async function linkGitHub(code: string) {
    try {
      await authApi.linkGitHub(code);
      // Optionally refresh user profile if it contains linked accounts info
    } catch (error) {
      console.error('GitHub linking failed:', error);
      throw error;
    }
  }

  function setSession(newToken: string, newUser: User) {
    token.value = newToken;
    user.value = newUser;
    isAuthenticated.value = true;
    localStorage.setItem('token', newToken);
    // Store user info if needed, or fetch it on app load
    localStorage.setItem('user', JSON.stringify(newUser));
  }

  function logout() {
    token.value = null;
    user.value = null;
    isAuthenticated.value = false;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  // Initialize from local storage
  const storedUser = localStorage.getItem('user');
  if (storedUser) {
    user.value = JSON.parse(storedUser);
  }

  return {
    user,
    token,
    isAuthenticated,
    loginWithGoogle,
    handleGoogleCallback,
    linkGitHub,
    logout,
    setSession // Export setSession for use by composables
  };
});

import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authApi, type User } from '@/api/auth';

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null);
  const isAuthenticated = ref(false);

  async function loginWithGoogle(googleToken: string) {
    try {
      const response = await authApi.loginWithGoogle(googleToken);
      setSession(response.data.user);
    } catch (error) {
      throw error;
    }
  }

  async function handleGoogleCallback(code: string) {
    try {
      const response = await authApi.handleGoogleCallback(code);
      setSession(response.data.user);
    } catch (error) {
      throw error;
    }
  }

  async function linkGitHub(code: string) {
    try {
      await authApi.linkGitHub(code);
      // Optionally refresh user profile if it contains linked accounts info
    } catch (error) {
      throw error;
    }
  }

  function setSession(newUser: User) {
    user.value = newUser;
    isAuthenticated.value = true;
    // Store user info for UI state (token is in HttpOnly cookie)
    localStorage.setItem('user', JSON.stringify(newUser));
  }

  async function logout() {
    try {
      // Call logout API to clear the cookie
      await authApi.logout();
    } catch (error) {
      // Continue with local cleanup even if API call fails
      console.error('Logout API call failed:', error);
    } finally {
      user.value = null;
      isAuthenticated.value = false;
      localStorage.removeItem('user');
    }
  }

  // Initialize from local storage
  const storedUser = localStorage.getItem('user');
  if (storedUser) {
    user.value = JSON.parse(storedUser);
    isAuthenticated.value = true;
  }

  return {
    user,
    isAuthenticated,
    loginWithGoogle,
    handleGoogleCallback,
    linkGitHub,
    logout,
    setSession // Export setSession for use by composables
  };
});

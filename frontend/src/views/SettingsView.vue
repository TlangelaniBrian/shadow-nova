<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Settings</h2>
      <p class="text-gray-400 mt-1">Manage your account preferences</p>
    </div>

    <!-- Profile Settings Card -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <h3 class="text-lg font-bold text-gray-900 mb-6">Profile Information</h3>

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
          <input v-model="username" type="text"
            class="w-full px-4 py-2 border border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent outline-none transition-all" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
          <input v-model="email" type="email"
            class="w-full px-4 py-2 border border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent outline-none transition-all" />
        </div>

        <button @click="updateProfile" :disabled="savingProfile"
          class="px-6 py-2 bg-purple-600 text-white rounded-xl hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
          {{ savingProfile ? 'Saving...' : 'Save Changes' }}
        </button>
      </div>
    </div>

    <!-- Password Change Card -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <h3 class="text-lg font-bold text-gray-900 mb-6">Change Password</h3>

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Current Password</label>
          <input v-model="currentPassword" type="password"
            class="w-full px-4 py-2 border border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent outline-none transition-all" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">New Password</label>
          <input v-model="newPassword" type="password"
            class="w-full px-4 py-2 border border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent outline-none transition-all" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Confirm New Password</label>
          <input v-model="confirmPassword" type="password"
            class="w-full px-4 py-2 border border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent outline-none transition-all" />
        </div>

        <button @click="changePassword" :disabled="savingPassword"
          class="px-6 py-2 bg-purple-600 text-white rounded-xl hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
          {{ savingPassword ? 'Updating...' : 'Change Password' }}
        </button>
      </div>
    </div>

    <!-- Notification Preferences Card -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <h3 class="text-lg font-bold text-gray-900 mb-6">Notification Preferences</h3>

      <div class="space-y-4">
        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-gray-700">Email Notifications</span>
          <input v-model="emailNotifications" type="checkbox"
            class="w-5 h-5 text-purple-600 rounded border-gray-300 focus:ring-2 focus:ring-purple-500" />
        </label>

        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-gray-700">Dark Mode</span>
          <input v-model="darkMode" type="checkbox"
            class="w-5 h-5 text-purple-600 rounded border-gray-300 focus:ring-2 focus:ring-purple-500" />
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useUserStore } from '@/stores/user';
import { userApi } from '@/api/user';
import { useToast } from '@/composables/useToast';

const userStore = useUserStore();
const toast = useToast();

const savingProfile = ref(false);
const savingPassword = ref(false);

const username = ref('');
const email = ref('');
const currentPassword = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const emailNotifications = ref(true);
const darkMode = ref(false);

onMounted(async () => {
    try {
        const response = await userApi.getUserProfile();
        username.value = response.data.username || '';
        email.value = response.data.email || '';
    } catch (error: any) {
        toast.error('Failed to load profile', error.response?.data?.error || 'Please try again');
    }
});

async function updateProfile() {
    if (!username.value.trim() || !email.value.trim()) {
        toast.error('Validation Error', 'Username and email are required');
        return;
    }

    savingProfile.value = true;
    try {
        const response = await userApi.updateUserProfile({
            username: username.value,
            email: email.value
        });

        // Update user store with new data
        if (userStore.user) {
            userStore.user.email = response.data.email;
            userStore.user.name = response.data.username || response.data.name;
            localStorage.setItem('user', JSON.stringify(userStore.user));
        }

        toast.success('Profile updated successfully');
    } catch (error: any) {
        toast.error('Failed to update profile', error.response?.data?.error || 'Please try again');
    } finally {
        savingProfile.value = false;
    }
}

async function changePassword() {
    if (!currentPassword.value || !newPassword.value || !confirmPassword.value) {
        toast.error('Validation Error', 'All password fields are required');
        return;
    }

    if (newPassword.value !== confirmPassword.value) {
        toast.error('Passwords do not match', 'New password and confirmation must match');
        return;
    }

    if (newPassword.value.length < 8) {
        toast.error('Weak Password', 'Password must be at least 8 characters long');
        return;
    }

    savingPassword.value = true;
    try {
        await userApi.updatePassword({
            current_password: currentPassword.value,
            new_password: newPassword.value
        });

        currentPassword.value = '';
        newPassword.value = '';
        confirmPassword.value = '';

        toast.success('Password changed successfully');
    } catch (error: any) {
        toast.error('Failed to change password', error.response?.data?.error || 'Please check your current password');
    } finally {
        savingPassword.value = false;
    }
}
</script>

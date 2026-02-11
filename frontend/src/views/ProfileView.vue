<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Profile Settings</h2>
        <p class="text-gray-400 mt-1">Manage your account and integrations</p>
      </div>
    </div>

    <!-- Profile Card Component -->
    <ProfileCard :user="profileData" :courses-completed="12" />

    <!-- Profile Edit Section -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <h3 class="text-xl font-bold text-gray-900 mb-6">Edit Profile</h3>

      <form @submit.prevent="updateProfile" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
          <input
            v-model="editForm.username"
            type="text"
            class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            placeholder="Enter username"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
          <input
            v-model="editForm.email"
            type="email"
            class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            placeholder="Enter email"
          />
        </div>

        <button
          type="submit"
          :disabled="isUpdating"
          class="px-6 py-2 bg-purple-600 text-white rounded-xl hover:bg-purple-700 disabled:opacity-50 transition-colors"
        >
          {{ isUpdating ? 'Updating...' : 'Update Profile' }}
        </button>
      </form>
    </div>

    <!-- Password Change Section -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <h3 class="text-xl font-bold text-gray-900 mb-6">Change Password</h3>

      <form @submit.prevent="changePassword" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Current Password</label>
          <input
            v-model="passwordForm.currentPassword"
            type="password"
            class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            placeholder="Enter current password"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">New Password</label>
          <input
            v-model="passwordForm.newPassword"
            type="password"
            class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            placeholder="Enter new password (min 8 characters)"
          />
        </div>

        <button
          type="submit"
          :disabled="isChangingPassword"
          class="px-6 py-2 bg-purple-600 text-white rounded-xl hover:bg-purple-700 disabled:opacity-50 transition-colors"
        >
          {{ isChangingPassword ? 'Changing...' : 'Change Password' }}
        </button>
      </form>
    </div>

    <!-- GitHub Integration Component -->
    <GitHubIntegration
      :is-loading="isLoading"
      :is-connecting="isConnecting"
      :is-connected="githubLinked"
      :username="githubUsername"
      :stats="githubStats"
      @connect="linkGitHub"
      @disconnect="unlinkGitHub"
    />

    <!-- Account Settings Component -->
    <AccountSettings
      :email-notifications="emailNotifications"
      :dark-mode="darkMode"
      @update:email-notifications="emailNotifications = $event"
      @update:dark-mode="darkMode = $event"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect, onMounted } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useToast } from '@/composables/useToast';
import client from '@/api/client';
import { disconnectGitHub } from '@/api/github';
import { userApi } from '@/api/user';
import ProfileCard from '@/components/profile/ProfileCard.vue';
import GitHubIntegration from '@/components/profile/GitHubIntegration.vue';
import AccountSettings from '@/components/profile/AccountSettings.vue';

const { user } = useAuth();
const toast = useToast();

const profileData = ref<any>(null);
const isLoading = ref(false);
const isConnecting = ref(false);
const isUpdating = ref(false);
const isChangingPassword = ref(false);

const githubLinked = ref(false);
const githubUsername = ref('');
const githubStats = ref({
  repos: 0,
  contributions: 0,
  followers: 0,
});

const emailNotifications = ref(true);
const darkMode = ref(false);

const editForm = ref({
  username: '',
  email: '',
});

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
});

const fetchProfile = async () => {
  try {
    const response = await userApi.getUserProfile();
    profileData.value = response.data;

    editForm.value.username = response.data.username || '';
    editForm.value.email = response.data.email || '';
  } catch (error: any) {
    console.error('[Profile] Error fetching profile:', error);
    const errorMessage = error.response?.data?.error || error.message || 'Failed to load profile';
    toast.error('Profile Load Failed', errorMessage);
  }
};

const updateProfile = async () => {
  if (isUpdating.value) return;

  const updates: any = {};
  if (editForm.value.username !== profileData.value?.username) {
    updates.username = editForm.value.username;
  }
  if (editForm.value.email !== profileData.value?.email) {
    updates.email = editForm.value.email;
  }

  if (Object.keys(updates).length === 0) {
    toast.info('No Changes', 'No profile updates to save');
    return;
  }

  isUpdating.value = true;

  try {
    const response = await userApi.updateUserProfile(updates);
    profileData.value = response.data;

    if (user.value) {
      user.value.username = response.data.username;
      user.value.email = response.data.email;
      localStorage.setItem('user', JSON.stringify(user.value));
    }

    toast.success('Profile Updated', 'Your profile has been updated successfully');
  } catch (error: any) {
    console.error('[Profile] Update error:', error);
    const errorMessage = error.response?.data?.error || error.message || 'Failed to update profile';
    toast.error('Update Failed', errorMessage);
  } finally {
    isUpdating.value = false;
  }
};

const changePassword = async () => {
  if (isChangingPassword.value) return;

  if (!passwordForm.value.currentPassword || !passwordForm.value.newPassword) {
    toast.error('Validation Error', 'Both current and new password are required');
    return;
  }

  if (passwordForm.value.newPassword.length < 8) {
    toast.error('Validation Error', 'New password must be at least 8 characters');
    return;
  }

  isChangingPassword.value = true;

  try {
    await userApi.updatePassword({
      current_password: passwordForm.value.currentPassword,
      new_password: passwordForm.value.newPassword,
    });

    passwordForm.value.currentPassword = '';
    passwordForm.value.newPassword = '';

    toast.success('Password Changed', 'Your password has been updated successfully');
  } catch (error: any) {
    console.error('[Profile] Password change error:', error);
    const errorMessage = error.response?.data?.error || error.message || 'Failed to change password';
    toast.error('Password Change Failed', errorMessage);
  } finally {
    isChangingPassword.value = false;
  }
};

const linkGitHub = async () => {
  if (isConnecting.value) return;

  isConnecting.value = true;

  try {
    const response = await client.get('/auth/github/connect');

    if (response.data.url) {
      toast.info('Redirecting to GitHub', 'Please authorize the application');
      window.location.href = response.data.url;
    } else {
      isConnecting.value = false;
      toast.error('Connection Failed', 'No redirect URL received from server');
    }
  } catch (error: any) {
    console.error('[GitHub] Error:', error);
    isConnecting.value = false;
    const errorMessage = error.response?.data?.error || error.message || 'Unknown error occurred';
    toast.error('GitHub Connection Failed', errorMessage);
  }
};

const unlinkGitHub = async () => {
  try {
    await disconnectGitHub();
    githubLinked.value = false;
    githubUsername.value = '';
    githubStats.value = {
      repos: 0,
      contributions: 0,
      followers: 0,
    };

    if (user.value) {
      user.value.github_username = undefined;
      localStorage.setItem('user', JSON.stringify(user.value));
    }

    toast.success('GitHub Disconnected', 'Your GitHub account has been unlinked');
  } catch (error: any) {
    console.error('[GitHub] Disconnect error:', error);
    const errorMessage = error.response?.data?.error || error.message || 'Failed to disconnect';
    toast.error('Disconnect Failed', errorMessage);
  }
};

onMounted(() => {
  fetchProfile();
});

watchEffect(() => {
  if (user.value?.github_username) {
    githubLinked.value = true;
    githubUsername.value = user.value.github_username;
    githubStats.value = {
      repos: 24,
      contributions: 1247,
      followers: 89,
    };
  } else {
    githubLinked.value = false;
    githubUsername.value = '';
  }
});
</script>

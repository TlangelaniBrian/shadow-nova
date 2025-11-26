<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Profile Settings</h2>
        <p class="text-gray-400 mt-1">Manage your account and integrations</p>
      </div>
    </div>

    <!-- Profile Card -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <div class="flex flex-col md:flex-row items-center md:items-start gap-6 text-center md:text-left">
        <!-- Avatar -->
        <div class="w-24 h-24 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white text-3xl font-bold shrink-0 shadow-lg">
          {{ user?.name?.charAt(0) || 'U' }}
        </div>
        
        <!-- User Info -->
        <div class="flex-1 w-full">
          <h3 class="text-2xl font-bold text-gray-900">{{ user?.name || 'User' }}</h3>
          <p class="text-gray-400 mt-1">{{ user?.email }}</p>
          
          <div class="flex flex-wrap justify-center md:justify-start gap-3 mt-4">
            <div class="px-4 py-2 bg-purple-50 rounded-xl">
              <p class="text-xs text-gray-500">Member since</p>
              <p class="text-sm font-bold text-gray-900">{{ memberSince }}</p>
            </div>
            <div class="px-4 py-2 bg-blue-50 rounded-xl">
              <p class="text-xs text-gray-500">Courses Completed</p>
              <p class="text-sm font-bold text-gray-900">12</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- GitHub Integration Card -->
    <div class="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-12 h-12 rounded-2xl bg-gray-900 flex items-center justify-center shrink-0">
          <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-bold text-gray-900">GitHub Integration</h3>
          <p class="text-sm text-gray-400">Connect your GitHub account to track projects</p>
        </div>
      </div>

      <div v-if="isLoading" class="flex items-center justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600"></div>
      </div>

      <div v-else-if="githubLinked" class="space-y-4">
        <!-- Connected State -->
        <div class="flex flex-col sm:flex-row items-center justify-between p-4 bg-green-50 border border-green-200 rounded-2xl gap-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-green-100 flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
            <div>
              <p class="font-bold text-gray-900">Connected</p>
              <p class="text-sm text-gray-600">@{{ githubUsername }}</p>
            </div>
          </div>
          <button 
            @click="unlinkGitHub"
            class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-50 rounded-xl transition-colors"
          >
            Disconnect
          </button>
        </div>

        <!-- GitHub Stats -->
        <div class="grid grid-cols-3 gap-3 md:gap-4 mt-4">
          <div class="p-4 bg-gray-50 rounded-2xl text-center">
            <p class="text-2xl font-bold text-gray-900">{{ githubStats.repos }}</p>
            <p class="text-xs text-gray-500 mt-1">Repositories</p>
          </div>
          <div class="p-4 bg-gray-50 rounded-2xl text-center">
            <p class="text-2xl font-bold text-gray-900">{{ githubStats.contributions }}</p>
            <p class="text-xs text-gray-500 mt-1">Contributions</p>
          </div>
          <div class="p-4 bg-gray-50 rounded-2xl text-center">
            <p class="text-2xl font-bold text-gray-900">{{ githubStats.followers }}</p>
            <p class="text-xs text-gray-500 mt-1">Followers</p>
          </div>
        </div>
      </div>

      <div v-else class="space-y-4">
        <!-- Not Connected State -->
        <div class="p-6 bg-gray-50 border border-gray-200 rounded-2xl text-center">
          <p class="text-gray-600 mb-4">Connect your GitHub account to:</p>
          <ul class="text-sm text-gray-500 space-y-2 mb-6">
            <li>✓ Track your project submissions</li>
            <li>✓ Sync your repositories</li>
            <li>✓ Showcase your work</li>
          </ul>
          <button 
            @click="linkGitHub"
            :disabled="isConnecting"
            class="w-full bg-gray-900 text-white font-bold py-3 rounded-xl hover:bg-gray-800 transition-colors flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <!-- Spinner when connecting -->
            <svg v-if="isConnecting" class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <!-- GitHub icon when not connecting -->
            <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            {{ isConnecting ? 'Connecting...' : 'Connect GitHub Account' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Account Settings -->
    <div class="bg-white rounded-3xl p-8 border border-gray-100 shadow-sm">
      <h3 class="text-lg font-bold text-gray-900 mb-6">Account Settings</h3>
      
      <div class="space-y-4">
        <div class="flex items-center justify-between py-3 border-b border-gray-100">
          <div>
            <p class="font-medium text-gray-900">Email Notifications</p>
            <p class="text-sm text-gray-400">Receive updates about your progress</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="emailNotifications" class="sr-only peer">
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
          </label>
        </div>

        <div class="flex items-center justify-between py-3">
          <div>
            <p class="font-medium text-gray-900">Dark Mode</p>
            <p class="text-sm text-gray-400">Switch to dark theme</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="darkMode" class="sr-only peer">
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watchEffect } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useToast } from '@/composables/useToast';
import client from '@/api/client';

const { user } = useAuth();
const toast = useToast();

const isLoading = ref(false);
const isConnecting = ref(false);
const githubLinked = ref(false);
const githubUsername = ref('');
const githubStats = ref({
  repos: 0,
  contributions: 0,
  followers: 0,
});

const emailNotifications = ref(true);
const darkMode = ref(false);

const memberSince = computed(() => {
  const date = new Date();
  return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
});

const linkGitHub = async () => {
  if (isConnecting.value) return; // Prevent double-clicks
  
  isConnecting.value = true;
  console.log('[GitHub] Starting connection...');
  
  // Debug: Check if token exists
  const token = localStorage.getItem('token');
  console.log('[GitHub] Token exists:', !!token);
  console.log('[GitHub] Token value:', token ? `${token.substring(0, 20)}...` : 'null');
  
  try {
    console.log('[GitHub] Calling /auth/github/connect...');
    const response = await client.get('/auth/github/connect');
    console.log('[GitHub] Response:', response.data);
    
    if (response.data.url) {
      console.log('[GitHub] Redirecting to:', response.data.url);
      toast.info('Redirecting to GitHub', 'Please authorize the application');
      window.location.href = response.data.url;
    } else {
      console.error('[GitHub] No URL in response');
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

const unlinkGitHub = () => {
  // TODO: Call API to unlink GitHub
  githubLinked.value = false;
  githubUsername.value = '';
  toast.info('GitHub Disconnected', 'Your GitHub account has been unlinked');
};

// Reactively update UI when user state changes (e.g. after linking)
watchEffect(() => {
  if (user.value?.github_username) {
    githubLinked.value = true;
    githubUsername.value = user.value.github_username;
    // Mock stats - would come from GitHub API
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

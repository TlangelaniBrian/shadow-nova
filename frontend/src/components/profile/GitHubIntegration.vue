<template>
  <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 md:p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
    <div class="flex items-center gap-3 mb-6">
      <div class="w-12 h-12 rounded-2xl bg-gray-900 dark:bg-gray-700 flex items-center justify-center shrink-0">
        <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
        </svg>
      </div>
      <div>
        <h3 class="text-lg font-bold text-gray-900 dark:text-gray-100">GitHub Integration</h3>
        <p class="text-sm text-gray-400 dark:text-gray-500">Connect your GitHub account to track projects</p>
      </div>
    </div>

    <div v-if="isLoading" class="flex items-center justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600"></div>
    </div>

    <div v-else-if="isConnected" class="space-y-4">
      <!-- Connected State -->
      <div class="flex flex-col sm:flex-row items-center justify-between p-4 bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-700 rounded-2xl gap-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-green-100 dark:bg-green-900/50 flex items-center justify-center shrink-0">
            <svg class="w-5 h-5 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>
          <div>
            <p class="font-bold text-gray-900 dark:text-gray-100">Connected</p>
            <p class="text-sm text-gray-600 dark:text-gray-400">@{{ username }}</p>
          </div>
        </div>
        <button
          @click="handleDisconnect"
          class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-xl transition-colors"
        >
          Disconnect
        </button>
      </div>

      <!-- GitHub Stats -->
      <div class="grid grid-cols-3 gap-3 md:gap-4 mt-4">
        <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-2xl text-center">
          <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ stats.repos }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Repositories</p>
        </div>
        <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-2xl text-center">
          <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ stats.contributions }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Contributions</p>
        </div>
        <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-2xl text-center">
          <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ stats.followers }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Followers</p>
        </div>
      </div>
    </div>

    <div v-else class="space-y-4">
      <!-- Not Connected State -->
      <div class="p-6 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-2xl text-center">
        <p class="text-gray-600 dark:text-gray-400 mb-4">Connect your GitHub account to:</p>
        <ul class="text-sm text-gray-500 dark:text-gray-400 space-y-2 mb-6">
          <li>✓ Track your project submissions</li>
          <li>✓ Sync your repositories</li>
          <li>✓ Showcase your work</li>
        </ul>
        <button
          @click="handleConnect"
          :disabled="isConnecting"
          class="w-full bg-gray-900 dark:bg-gray-700 text-white font-bold py-3 rounded-xl hover:bg-gray-800 dark:hover:bg-gray-600 transition-colors flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
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
</template>

<script setup lang="ts">
interface GitHubStats {
  repos: number;
  contributions: number;
  followers: number;
}

interface Props {
  isLoading?: boolean;
  isConnecting?: boolean;
  isConnected?: boolean;
  username?: string;
  stats?: GitHubStats;
}

withDefaults(defineProps<Props>(), {
  isLoading: false,
  isConnecting: false,
  isConnected: false,
  username: '',
  stats: () => ({ repos: 0, contributions: 0, followers: 0 }),
});

interface Emits {
  (e: 'connect'): void;
  (e: 'disconnect'): void;
}

const emit = defineEmits<Emits>();

const handleConnect = () => {
  emit('connect');
};

const handleDisconnect = () => {
  emit('disconnect');
};
</script>

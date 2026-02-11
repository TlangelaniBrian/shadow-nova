<template>
  <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 md:p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
    <div class="flex flex-col md:flex-row items-center md:items-start gap-6 text-center md:text-left">
      <!-- Avatar -->
      <div class="w-24 h-24 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white text-3xl font-bold shrink-0 shadow-lg">
        {{ userInitial }}
      </div>

      <!-- User Info -->
      <div class="flex-1 w-full">
        <h3 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ userName }}</h3>
        <p class="text-gray-400 dark:text-gray-500 mt-1">{{ user?.email }}</p>

        <div class="flex flex-wrap justify-center md:justify-start gap-3 mt-4">
          <div class="px-4 py-2 bg-purple-50 dark:bg-purple-900/30 rounded-xl">
            <p class="text-xs text-gray-500 dark:text-gray-400">Member since</p>
            <p class="text-sm font-bold text-gray-900 dark:text-gray-100">{{ memberSince }}</p>
          </div>
          <div class="px-4 py-2 bg-blue-50 dark:bg-blue-900/30 rounded-xl">
            <p class="text-xs text-gray-500 dark:text-gray-400">Courses Completed</p>
            <p class="text-sm font-bold text-gray-900 dark:text-gray-100">{{ coursesCompleted }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface User {
  name?: string;
  email?: string;
  created_at?: string;
}

interface Props {
  user?: User;
  coursesCompleted?: number;
}

const props = withDefaults(defineProps<Props>(), {
  coursesCompleted: 12,
});

const userInitial = computed(() => {
  return props.user?.name?.charAt(0) || 'U';
});

const userName = computed(() => {
  return props.user?.name || 'User';
});

const memberSince = computed(() => {
  const date = new Date();
  return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
});
</script>

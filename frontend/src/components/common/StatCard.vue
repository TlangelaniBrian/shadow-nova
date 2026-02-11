<template>
  <div
    class="bg-white dark:bg-gray-800 rounded-3xl p-6 border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition-all"
  >
    <div class="flex items-center justify-between mb-4">
      <div
        class="w-12 h-12 rounded-2xl flex items-center justify-center text-2xl"
        :class="colorClass"
      >
        {{ icon }}
      </div>
      <span
        v-if="trend"
        class="text-green-500 dark:text-green-400 text-sm font-medium bg-green-50 dark:bg-green-900/30 px-2 py-1 rounded-full"
      >
        {{ trend }}
      </span>
    </div>
    <p class="text-gray-400 dark:text-gray-500 text-sm mb-1">{{ label }}</p>
    <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ value }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  label: string;
  value: string | number;
  trend?: string;
  icon: string;
  color?: 'purple' | 'blue' | 'green' | 'orange' | 'pink' | 'red';
}

const props = withDefaults(defineProps<Props>(), {
  color: 'purple',
});

const colorClass = computed(() => {
  const colors = {
    purple: 'bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400',
    blue: 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400',
    green: 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400',
    orange: 'bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400',
    pink: 'bg-pink-100 dark:bg-pink-900/30 text-pink-600 dark:text-pink-400',
    red: 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400',
  };
  return colors[props.color];
});
</script>

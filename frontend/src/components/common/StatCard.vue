<template>
  <div
    class="bg-white rounded-3xl p-6 border border-gray-100 shadow-sm hover:shadow-md transition-all"
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
        class="text-green-500 text-sm font-medium bg-green-50 px-2 py-1 rounded-full"
      >
        {{ trend }}
      </span>
    </div>
    <p class="text-gray-400 text-sm mb-1">{{ label }}</p>
    <p class="text-2xl font-bold text-gray-900">{{ value }}</p>
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
    purple: 'bg-purple-100 text-purple-600',
    blue: 'bg-blue-100 text-blue-600',
    green: 'bg-green-100 text-green-600',
    orange: 'bg-orange-100 text-orange-600',
    pink: 'bg-pink-100 text-pink-600',
    red: 'bg-red-100 text-red-600',
  };
  return colors[props.color];
});
</script>

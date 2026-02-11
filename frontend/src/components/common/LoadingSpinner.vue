<template>
  <div class="flex items-center justify-center" :class="containerSizeClass">
    <div
      class="animate-spin rounded-full border-b-2"
      :class="[borderColorClass, spinnerSizeClass]"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  size?: 'sm' | 'md' | 'lg';
  color?: 'purple' | 'blue' | 'green' | 'red' | 'white';
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'purple',
});

const containerSizeClass = computed(() => {
  const sizes = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };
  return sizes[props.size];
});

const spinnerSizeClass = computed(() => {
  const sizes = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };
  return sizes[props.size];
});

const borderColorClass = computed(() => {
  const colors = {
    purple: 'border-purple-600',
    blue: 'border-blue-600',
    green: 'border-green-600',
    red: 'border-red-600',
    white: 'border-white',
  };
  return colors[props.color];
});
</script>

<template>
  <div class="mt-8 space-y-3">
    <div
      v-for="feature in features"
      :key="feature.text"
      class="flex items-center gap-3 text-white/80 text-sm"
    >
      <svg
        class="w-5 h-5 text-green-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M5 13l4 4L19 7"
        />
      </svg>
      <span>{{ feature.text }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Feature {
  text: string;
}

const features: Feature[] = [
  { text: 'Structured learning paths' },
  { text: 'Hands-on projects with real code' },
  { text: 'Progress tracking & achievements' },
];
</script>

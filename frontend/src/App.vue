<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import { computed } from 'vue'

const route = useRoute()
const showLayout = computed(() => route.path !== '/login' && route.name !== 'login')
</script>

<template>
  <AppLayout v-if="showLayout">
    <RouterView />
  </AppLayout>
  <RouterView v-else />
</template>

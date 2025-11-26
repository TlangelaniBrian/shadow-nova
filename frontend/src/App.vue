<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import { Toaster } from 'vue-sonner'
import { computed } from 'vue'

const route = useRoute()
const showLayout = computed(() => route.path !== '/login' && route.name !== 'login')
</script>

<template>
  <Toaster position="top-right" richColors />
  <AppLayout v-if="showLayout">
    <RouterView />
  </AppLayout>
  <RouterView v-else />
</template>

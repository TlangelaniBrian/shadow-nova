<script setup lang="ts">
import { RouterView, useRoute, useRouter } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import RouteLoading from '@/components/RouteLoading.vue'
import { Toaster } from 'vue-sonner'
import { computed, ref } from 'vue'

const route = useRoute()
const router = useRouter()
const isLoading = ref(false)

const showLayout = computed(() => route.path !== '/login' && route.name !== 'login')

router.beforeEach(() => {
  isLoading.value = true
})

router.afterEach(() => {
  isLoading.value = false
})
</script>

<template>
  <Toaster position="top-right" richColors />
  <RouteLoading v-if="isLoading" />
  <AppLayout v-else-if="showLayout">
    <RouterView />
  </AppLayout>
  <RouterView v-else />
</template>

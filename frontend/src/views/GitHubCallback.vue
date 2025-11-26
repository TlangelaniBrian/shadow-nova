<template>
  <div class="flex items-center justify-center min-h-screen bg-slate-900 text-white">
    <div class="text-center">
      <h2 class="text-2xl font-bold mb-4">Linking GitHub Account...</h2>
      <p v-if="errorMessage" class="text-red-500">{{ errorMessage }}</p>
      <div v-else-if="isLoading" class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
      <p v-else class="text-green-500">Successfully linked!</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuth } from '@/composables/useAuth';
import { useErrorHandler } from '@/composables/useErrorHandler';

const route = useRoute();
const router = useRouter();
const { linkGitHub, isLoading } = useAuth();
const { getUserMessage } = useErrorHandler();
const errorMessage = ref('');

onMounted(async () => {
  const code = route.query.code as string;
  
  if (!code) {
    errorMessage.value = 'No code provided';
    setTimeout(() => router.push('/profile'), 2000);
    return;
  }

  const result = await linkGitHub(code);
  
  if (result.error) {
    errorMessage.value = getUserMessage(result.error);
    setTimeout(() => router.push('/profile'), 3000);
  } else {
    // Success - redirect to profile
    setTimeout(() => router.push('/profile'), 1000);
  }
});
</script>

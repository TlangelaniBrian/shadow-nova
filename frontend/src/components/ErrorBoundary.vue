<template>
  <div v-if="hasError" class="error-boundary">
    <slot name="error" :error="errorInfo" :reset="reset">
      <div class="error-fallback">
        <div class="error-icon">⚠️</div>
        <h2 class="error-title">Something went wrong</h2>
        <p class="error-message">{{ errorMessage }}</p>
        <button @click="reset" class="error-reset-btn">
          Try Again
        </button>
      </div>
    </slot>
  </div>
  <slot v-else></slot>
</template>

<script setup lang="ts">
import { ref, onErrorCaptured, type Ref } from 'vue';
import type { AppError } from '@/types/errors';
import { transformAxiosError } from '@/types/errors';

interface Props {
  fallback?: string;
  onError?: (error: AppError) => void;
}

const props = withDefaults(defineProps<Props>(), {
  fallback: 'An error occurred',
});

const hasError = ref(false);
const errorInfo: Ref<AppError | null> = ref(null);
const errorMessage = ref('');

onErrorCaptured((err: unknown) => {
  hasError.value = true;
  const appError = transformAxiosError(err);
  errorInfo.value = appError;
  errorMessage.value = appError.message || props.fallback;
  
  // Call custom error handler if provided
  if (props.onError) {
    props.onError(appError);
  }
  
  // Log error
  console.error('[ErrorBoundary]', appError);
  
  // Prevent error from propagating
  return false;
});

const reset = () => {
  hasError.value = false;
  errorInfo.value = null;
  errorMessage.value = '';
};
</script>

<style scoped>
.error-boundary {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-fallback {
  text-align: center;
  padding: 2rem;
  max-width: 500px;
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.error-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.error-message {
  color: #6b7280;
  margin-bottom: 1.5rem;
}

.error-reset-btn {
  background-color: #3b82f6;
  color: white;
  padding: 0.5rem 1.5rem;
  border-radius: 0.5rem;
  border: none;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.error-reset-btn:hover {
  background-color: #2563eb;
}
</style>

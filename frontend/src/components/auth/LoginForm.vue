<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-gray-200 mb-1">Email</label>
      <input
        v-model="localEmail"
        type="email"
        required
        class="w-full px-4 py-2 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500"
        placeholder="name@example.com"
      />
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-200 mb-1">Password</label>
      <input
        v-model="localPassword"
        type="password"
        required
        class="w-full px-4 py-2 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500"
        placeholder="••••••••"
      />
    </div>
    <button
      type="submit"
      :disabled="isLoading"
      class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-purple-600 hover:bg-purple-700 transition-all duration-200 rounded-lg text-white font-medium shadow-lg hover:shadow-purple-500/30 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span v-if="isLoading">Signing in...</span>
      <span v-else>Sign in with Email</span>
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

interface Props {
  email?: string;
  password?: string;
  isLoading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  email: '',
  password: '',
  isLoading: false,
});

interface Emits {
  (e: 'submit', credentials: { email: string; password: string }): void;
}

const emit = defineEmits<Emits>();

const localEmail = ref(props.email);
const localPassword = ref(props.password);

watch(() => props.email, (newEmail) => {
  localEmail.value = newEmail;
});

watch(() => props.password, (newPassword) => {
  localPassword.value = newPassword;
});

const handleSubmit = () => {
  emit('submit', {
    email: localEmail.value,
    password: localPassword.value,
  });
};
</script>

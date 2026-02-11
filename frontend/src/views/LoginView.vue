<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthProviders from '@/components/auth/AuthProviders.vue'
import LoginForm from '@/components/auth/LoginForm.vue'
import LoginFeatures from '@/components/auth/LoginFeatures.vue'

const router = useRouter()
const isLoading = ref(false)

const user = localStorage.getItem('user')
if (user) {
  router.push('/dashboard')
}

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const handleLogin = async (credentials: { email: string; password: string }) => {
  isLoading.value = true
  try {
    const res = await fetch(`${apiUrl}/api/v1/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(credentials),
    })

    if (!res.ok) {
      throw new Error('Invalid credentials')
    }

    const data = await res.json()

    localStorage.setItem('user', JSON.stringify({
      id: data.data.id,
      email: data.data.email,
      username: data.data.username,
      role: data.data.role
    }))

    router.push('/dashboard')

  } catch (error) {
    console.error(error)
    alert('Login failed: Invalid credentials')
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div
    class="min-h-screen bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-800 flex items-center justify-center p-4"
  >
    <!-- Background Pattern -->
    <div
      class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4wNSI+PHBhdGggZD0iTTM2IDEzNGgyLTJNMjggMTloMi0ybTAgMTVoMi0ybTAgMTVoMi0ybTAgMTVoMi0yTTEzIDM0aDItMm0wIDE1aDItMm0wIDE1aDItMm0xNS0zMGgyLTJtMCAxNWgyLTJtMCAxNWgyLTJNNDMgMzRoMi0ybTAgMTVoMi0ybTAgMTVoMi0ybTAtNDVoMi0ybTAgMTVoMi0yIi8+PC9nPjwvZz48L3N2Zz4=')] opacity-20"
    ></div>

    <!-- Login Card -->
    <div class="relative w-full max-w-md">
      <!-- Glowing Effect -->
      <div
        class="absolute -inset-1 bg-gradient-to-r from-purple-600 to-pink-600 rounded-2xl blur-xl opacity-30"
      ></div>

      <!-- Card -->
      <div
        class="relative bg-white/10 backdrop-blur-xl rounded-2xl p-8 shadow-2xl border border-white/20"
      >
        <!-- Logo & Title -->
        <div class="text-center mb-8">
          <div
            class="inline-block p-3 bg-gradient-to-br from-purple-500 to-pink-500 rounded-2xl mb-4"
          >
            <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          </div>
          <h1 class="text-4xl font-bold text-white mb-2">Shadow Nova</h1>
          <p class="text-purple-200">Master the Stack. Build Your Future.</p>
        </div>

        <!-- Description -->
        <div class="mb-8 text-center">
          <p class="text-white/80 text-sm">
            The ultimate platform for junior developers to learn Vue, Go, and Cloud technologies
            through hands-on projects.
          </p>
        </div>

        <!-- Auth Providers -->
        <div class="space-y-4">
          <AuthProviders />

          <!-- Divider -->
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/20"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-transparent text-purple-200">Or continue with</span>
            </div>
          </div>

          <!-- Login Form -->
          <LoginForm :is-loading="isLoading" @submit="handleLogin" />
        </div>

        <!-- Features -->
        <LoginFeatures />

        <!-- Footer -->
        <div class="mt-8 text-center text-xs text-white/60">
          <p>By signing in, you agree to our Terms of Service and Privacy Policy</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>

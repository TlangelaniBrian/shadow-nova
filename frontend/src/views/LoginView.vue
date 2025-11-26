<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import GoogleSignIn from '@/components/GoogleSignIn.vue'

const router = useRouter()
const isLoading = ref(false)

// Check if user is already logged in
const token = localStorage.getItem('token')
if (token) {
  router.push('/dashboard')
}

const email = ref('')
const password = ref('')
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const handleLogin = async () => {
  isLoading.value = true
  try {
    const res = await fetch(`${apiUrl}/api/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    })

    if (!res.ok) {
      throw new Error('Invalid credentials')
    }

    const data = await res.json()
    
    localStorage.setItem('token', data.data.token)
    localStorage.setItem('user', JSON.stringify({
      email: data.data.email,
      username: data.data.username
    }))
    
    // Check feature flag (same as Google login)
        // Note: Ideally we should use the injected unleash client here too
    router.push('/dashboard') // Default to dashboard for now
    
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

        <!-- Google Sign In -->
        <div class="space-y-4">
          <GoogleSignIn />

          <!-- Divider -->
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/20"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-transparent text-purple-200">Or continue with</span>
            </div>
          </div>

          <!-- Email Login Form -->
          <form @submit.prevent="handleLogin" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-200 mb-1">Email</label>
              <input
                v-model="email"
                type="email"
                required
                class="w-full px-4 py-2 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="name@example.com"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-200 mb-1">Password</label>
              <input
                v-model="password"
                type="password"
                required
                class="w-full px-4 py-2 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="••••••••"
              />
            </div>
            <button
              type="submit"
              :disabled="isLoading"
              class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-purple-600 hover:bg-purple-700 transition-all duration-200 rounded-lg text-white font-medium shadow-lg hover:shadow-purple-500/30"
            >
              <span v-if="isLoading">Signing in...</span>
              <span v-else>Sign in with Email</span>
            </button>
          </form>
        </div>

        <!-- Features -->
        <div class="mt-8 space-y-3">
          <div class="flex items-center gap-3 text-white/80 text-sm">
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
            <span>Structured learning paths</span>
          </div>
          <div class="flex items-center gap-3 text-white/80 text-sm">
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
            <span>Hands-on projects with real code</span>
          </div>
          <div class="flex items-center gap-3 text-white/80 text-sm">
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
            <span>Progress tracking & achievements</span>
          </div>
        </div>

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

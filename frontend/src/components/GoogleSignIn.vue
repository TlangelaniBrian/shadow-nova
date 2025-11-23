<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { inject } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'

const router = useRouter()
const unleashClient = inject('unleash') as any
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:3000'
const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID

const isLoading = ref(false)
const error = ref('')

// Declare google on window
declare global {
  interface Window {
    google: any
    handleGoogleResponse: (response: any) => void
  }
}

// Initialize Google Sign-In
const loadGoogleScript = () => {
  const script = document.createElement('script')
  script.src = 'https://accounts.google.com/gsi/client'
  script.async = true
  script.defer = true
  document.head.appendChild(script)

  script.onload = () => {
    if (window.google) {
      window.google.accounts.id.initialize({
        client_id: clientId,
        callback: handleGoogleResponse,
      })
    }
  }
}

// Handle Google Sign-In response
const handleGoogleResponse = async (response: any) => {
  isLoading.value = true
  error.value = ''

  try {
    // Send Google ID token to backend for verification
    const res = await fetch(`${apiUrl}/api/auth/google/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id_token: response.credential,
      }),
    })

    if (!res.ok) {
      throw new Error('Failed to verify Google token')
    }

    const data = await res.json()

    // Store JWT token and user info
    if (typeof window !== 'undefined') {
      window.localStorage.setItem('auth_token', data.token)
      window.localStorage.setItem('user', JSON.stringify(data.user))
    }

    // Check feature flag for landing page
    const unleash = window.localStorage.getItem('unleash') // Fallback if inject doesn't work in async
    // Better way: use the global property or inject
    
    // Redirect based on feature flag
    // We need to access the unleash instance. Since we can't easily inject inside the async function callback context if not setup,
    // let's try to get it from the component setup.
    
    if (unleashClient && unleashClient.isEnabled('enable-landing')) {
      router.push('/home')
    } else {
      router.push('/dashboard')
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to sign in with Google'
  } finally {
    isLoading.value = false
  }
}

// Trigger Google Sign-In
const signInWithGoogle = () => {
  if (window.google) {
    window.google.accounts.id.prompt()
  }
}

// Load Google script on component mount
onMounted(() => {
  if (typeof window !== 'undefined') {
    // Expose handler to window for Google library to call
    window.handleGoogleResponse = handleGoogleResponse
    loadGoogleScript()
  }
})
</script>

<template>
  <div class="w-full">
    <Button
      @click="signInWithGoogle"
      :disabled="isLoading"
      class="w-full flex items-center justify-center gap-3 px-6 py-3 bg-white hover:bg-gray-50 text-gray-900 font-medium rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl"
      variant="outline"
    >
      <svg class="w-5 h-5" viewBox="0 0 24 24">
        <path
          fill="#4285F4"
          d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
        />
        <path
          fill="#34A853"
          d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        />
        <path
          fill="#FBBC05"
          d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
        />
        <path
          fill="#EA4335"
          d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        />
      </svg>
      <span v-if="isLoading" class="text-gray-700">Signing in...</span>
      <span v-else class="text-gray-900 font-semibold">Sign in with Google</span>
    </Button>

    <p v-if="error" class="text-red-400 text-sm mt-3 text-center">{{ error }}</p>

    <!-- Google One Tap container -->
    <div
      id="g_id_onload"
      :data-client_id="clientId"
      data-callback="handleGoogleResponse"
      data-auto_prompt="false"
    ></div>
  </div>
</template>

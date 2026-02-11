import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import unleashPlugin from './plugins/unleash'
import apiClient from './api/client'
// Uncomment after installing: npm install pinia-plugin-persistedstate
// import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

// Clean up old token storage (tokens are now in HttpOnly cookies)
if (typeof window !== 'undefined') {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('token')
}

const app = createApp(App)

// Fetch CSRF token before mounting app
async function initApp() {
  try {
    const response = await apiClient.get('/csrf-token')
    const csrfToken = response.data.csrf_token

    // Store in memory (not localStorage - that's XSS vulnerable)
    window.__CSRF_TOKEN__ = csrfToken
  } catch (error) {
    console.error('Failed to fetch CSRF token:', error)
  }

  // Mount app
  const pinia = createPinia()
  // Uncomment after installing pinia-plugin-persistedstate
  // pinia.use(piniaPluginPersistedstate)
  app.use(pinia)
  app.use(router)
  app.use(unleashPlugin)
  app.mount('#app')
}

initApp()

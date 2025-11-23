import { UnleashClient } from 'unleash-proxy-client'
import type { App } from 'vue'

export default {
  install: (app: App) => {
    const unleash = new UnleashClient({
      url: import.meta.env.VITE_UNLEASH_URL || 'http://localhost:4242/api/frontend',
      clientKey:
        import.meta.env.VITE_UNLEASH_TOKEN || 'default:development.unleash-insecure-frontend-token',
      appName: 'shadow-nova-frontend',
    })

    unleash.start()

    app.provide('unleash', unleash)
    app.config.globalProperties.$unleash = unleash
  },
}

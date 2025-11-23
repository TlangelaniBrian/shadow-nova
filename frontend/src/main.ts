import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import unleashPlugin from './plugins/unleash'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(unleashPlugin)
app.mount('#app')

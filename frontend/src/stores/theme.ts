import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
    const isDarkMode = ref(false)

    // Load from localStorage on init
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme === 'dark') {
        isDarkMode.value = true
    } else if (savedTheme === null) {
        // Check system preference
        isDarkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    }

    // Apply theme on init
    applyTheme()

    // Watch for changes and apply
    watch(isDarkMode, () => {
        applyTheme()
        localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light')
    })

    function applyTheme() {
        if (isDarkMode.value) {
            document.documentElement.classList.add('dark')
        } else {
            document.documentElement.classList.remove('dark')
        }
    }

    function toggleDarkMode() {
        isDarkMode.value = !isDarkMode.value
    }

    return {
        isDarkMode,
        toggleDarkMode
    }
})

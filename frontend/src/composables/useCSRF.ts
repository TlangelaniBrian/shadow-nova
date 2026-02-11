import { ref } from 'vue'
import apiClient from '@/api/client'

const csrfToken = ref<string>('')

export function useCSRF() {
    async function refreshToken() {
        try {
            const response = await apiClient.get('/csrf-token')
            csrfToken.value = response.data.csrf_token
            window.__CSRF_TOKEN__ = csrfToken.value
        } catch (error) {
            console.error('Failed to refresh CSRF token:', error)
            throw error
        }
    }

    return {
        csrfToken,
        refreshToken
    }
}

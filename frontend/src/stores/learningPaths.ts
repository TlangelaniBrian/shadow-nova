import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { pathsApi, type LearningPath } from '@/api/paths'
import type { Result } from '@/types/errors'
import { transformAxiosError, success, failure } from '@/types/errors'

export const useLearningPathsStore = defineStore('learningPaths', () => {
    const paths = ref<LearningPath[]>([])
    const currentPath = ref<LearningPath | null>(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    const pathsByDifficulty = computed(() => {
        const grouped: Record<string, LearningPath[]> = {}
        paths.value.forEach(path => {
            if (!grouped[path.difficulty]) {
                grouped[path.difficulty] = []
            }
            grouped[path.difficulty].push(path)
        })
        return grouped
    })

    async function fetchPaths(): Promise<Result<LearningPath[]>> {
        loading.value = true
        error.value = null
        try {
            const response = await pathsApi.getLearningPaths()
            paths.value = response.data
            return success(response.data)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    async function fetchPath(id: string): Promise<Result<LearningPath>> {
        loading.value = true
        error.value = null
        try {
            const response = await pathsApi.getLearningPath(id)
            currentPath.value = response.data
            return success(response.data)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    return {
        paths,
        currentPath,
        loading,
        error,
        pathsByDifficulty,
        fetchPaths,
        fetchPath
    }
}, {
    // Enable after installing pinia-plugin-persistedstate
    // persist: {
    //     paths: ['paths', 'currentPath']
    // }
})

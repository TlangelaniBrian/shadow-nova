import { defineStore } from 'pinia'
import { ref } from 'vue'
import { progressApi, type UserProgress, type UserStats, type UpdateProgressRequest } from '@/api/progress'
import type { Result } from '@/types/errors'
import { transformAxiosError, success, failure } from '@/types/errors'

export const useProgressStore = defineStore('progress', () => {
    const progressMap = ref<Map<number, boolean>>(new Map())
    const stats = ref<UserStats | null>(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    async function updateProgress(lessonId: number, completed: boolean): Promise<Result<void>> {
        loading.value = true
        error.value = null
        try {
            const data: UpdateProgressRequest = { lesson_id: lessonId, completed }
            await progressApi.updateProgress(data)
            progressMap.value.set(lessonId, completed)
            await fetchStats()
            return success(undefined as void)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    async function fetchStats(): Promise<Result<UserStats>> {
        try {
            const response = await progressApi.getStats()
            stats.value = response.data
            return success(response.data)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        }
    }

    async function fetchPathProgress(pathId: string): Promise<Result<UserProgress[]>> {
        loading.value = true
        error.value = null
        try {
            const response = await progressApi.getPathProgress(pathId)
            response.data.forEach((item: UserProgress) => {
                progressMap.value.set(item.lesson_id, item.completed)
            })
            return success(response.data)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    function isLessonCompleted(lessonId: number): boolean {
        return progressMap.value.get(lessonId) || false
    }

    return {
        progressMap,
        stats,
        loading,
        error,
        updateProgress,
        fetchStats,
        fetchPathProgress,
        isLessonCompleted
    }
}, {
    // Enable after installing pinia-plugin-persistedstate
    // persist: {
    //     paths: ['progressMap', 'stats']
    // }
})

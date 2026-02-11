import client from './client'

export interface UserProgress {
    id: number
    user_id: number
    lesson_id: number
    completed: boolean
    completed_at: string
    created_at: string
}

export interface UserStats {
    courses_completed: number
    hours_learned: number
    rank: number
    current_streak: number
    total_xp: number
}

export interface PathProgress {
    path_id: string
    total_lessons: number
    completed_lessons: number
    percentage: number
}

export interface UpdateProgressRequest {
    lesson_id: number
    completed: boolean
}

export const progressApi = {
    updateProgress(data: UpdateProgressRequest) {
        return client.post<void>('/progress', data)
    },

    getStats() {
        return client.get<UserStats>('/progress/stats')
    },

    getPathProgress(pathId: string) {
        return client.get<UserProgress[]>(`/progress/paths/${pathId}`)
    }
}

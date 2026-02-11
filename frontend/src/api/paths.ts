import client from './client'

export interface Module {
    id: number
    path_id: string
    title: string
    description: string
    order_index: number
    created_at: string
    lessons?: Lesson[]
}

export interface Lesson {
    id: number
    module_id: number
    title: string
    content_type: string
    content_url?: string
    content_body?: string
    duration_minutes: number
    order_index: number
    created_at: string
}

export interface LearningPath {
    id: string
    title: string
    description: string
    difficulty: string
    created_at: string
    modules?: Module[]
}

export const pathsApi = {
    getLearningPaths() {
        return client.get<LearningPath[]>('/learning-paths')
    },

    getLearningPath(id: string) {
        return client.get<LearningPath>(`/learning-paths/${id}`)
    }
}

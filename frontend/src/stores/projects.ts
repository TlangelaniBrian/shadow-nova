import { defineStore } from 'pinia'
import { ref } from 'vue'
import { projectsApi, type Project, type SubmitProjectRequest } from '@/api/projects'
import type { Result } from '@/types/errors'
import { transformAxiosError, success, failure } from '@/types/errors'

export interface ProjectSubmission {
    id: number
    user_id: number
    project_id: string
    github_repo_url: string
    pr_url?: string
    demo_url?: string
    status: string
    feedback?: string
    submitted_at: string
}

export const useProjectsStore = defineStore('projects', () => {
    const projects = ref<Project[]>([])
    const submissions = ref<ProjectSubmission[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)

    async function fetchProjects(): Promise<Result<Project[]>> {
        loading.value = true
        error.value = null
        try {
            const response = await projectsApi.getProjects()
            projects.value = response.data
            return success(response.data)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    async function submitProject(data: SubmitProjectRequest): Promise<Result<void>> {
        loading.value = true
        error.value = null
        try {
            await projectsApi.submitProject(data)
            return success(undefined as void)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    return {
        projects,
        submissions,
        loading,
        error,
        fetchProjects,
        submitProject
    }
}, {
    // Enable after installing pinia-plugin-persistedstate
    // persist: {
    //     paths: ['projects']
    // }
})

import { ref, type Ref } from 'vue';
import { projectsApi, type Project, type CreateProjectRequest, type SubmitProjectRequest } from '@/api/projects';
import type { Result, AppError } from '@/types/errors';
import { transformAxiosError, success, failure } from '@/types/errors';

interface UseProjectsReturn {
    projects: Ref<Project[]>;
    isLoading: Ref<boolean>;
    error: Ref<AppError | null>;
    fetchProjects: () => Promise<Result<Project[]>>;
    createProject: (data: CreateProjectRequest) => Promise<Result<Project>>;  // Changed from void to Project
    submitProject: (data: SubmitProjectRequest) => Promise<Result<void>>;
}

export function useProjects(): UseProjectsReturn {
    const projects = ref<Project[]>([]);
    const isLoading = ref(false);
    const error = ref<AppError | null>(null);

    const fetchProjects = async (): Promise<Result<Project[]>> => {
        isLoading.value = true;
        error.value = null;

        try {
            const response = await projectsApi.getProjects();
            // Backend returns: { message: string, data: Project[] }
            // Axios gives us response.data = { message, data }
            // So we need response.data.data to get the actual projects
            const projectsData = (response.data as any).data || response.data;
            projects.value = projectsData;
            isLoading.value = false;
            return success(projectsData);
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            isLoading.value = false;
            return failure(appError);
        }
    };

    const createProject = async (data: CreateProjectRequest): Promise<Result<Project>> => {
        isLoading.value = true;
        error.value = null;

        try {
            const response = await projectsApi.createProject(data);
            // Unwrap backend response
            const newProject = (response.data as any).data || response.data;
            projects.value.push(newProject);
            isLoading.value = false;
            return success(newProject);
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            isLoading.value = false;
            return failure(appError);
        }
    };

    const submitProject = async (data: SubmitProjectRequest): Promise<Result<void>> => {
        isLoading.value = true;
        error.value = null;

        try {
            await projectsApi.submitProject(data);
            isLoading.value = false;
            return success(undefined as void);
        } catch (err) {
            const appError = transformAxiosError(err);
            error.value = appError;
            isLoading.value = false;
            return failure(appError);
        }
    };

    return {
        projects,
        isLoading,
        error,
        fetchProjects,
        createProject,
        submitProject,
    };
}

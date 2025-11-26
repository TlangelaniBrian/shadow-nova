import client from './client';

export interface Project {
    id: string;
    title: string;
    description: string;
    difficulty: string;
    tech_stack: string[];
    icon?: string;
    github_repo?: string;
    created_at?: string;
}

export interface CreateProjectRequest {
    id: string;
    title: string;
    description: string;
    difficulty: string;
    technologies: string[];
}

export interface SubmitProjectRequest {
    project_id: string;
    github_repo_url: string;
    pr_url?: string;
    demo_url?: string;
}

export const projectsApi = {
    getProjects() {
        return client.get<Project[]>('/projects');
    },

    createProject(data: CreateProjectRequest) {
        return client.post('/projects', data);
    },

    submitProject(data: SubmitProjectRequest) {
        return client.post('/projects/submit', data);
    },
};

import { Layout, Server, Globe } from 'lucide-vue-next'

export const projects = [
  {
    id: 'portfolio-site',
    title: 'Personal Portfolio',
    description: 'Build a responsive portfolio website to showcase your skills and projects.',
    techStack: ['Vue.js', 'Tailwind CSS'],
    difficulty: 'Beginner',
    icon: Layout,
  },
  {
    id: 'task-manager',
    title: 'Task Manager API',
    description: 'Create a RESTful API for a task management application using Go and PostgreSQL.',
    techStack: ['Go', 'Chi', 'PostgreSQL'],
    difficulty: 'Intermediate',
    icon: Server,
  },
  {
    id: 'weather-dashboard',
    title: 'Weather Dashboard',
    description: 'Build a weather dashboard that fetches data from a public API.',
    techStack: ['Vue.js', 'Axios', 'Chart.js'],
    difficulty: 'Intermediate',
    icon: Globe,
  },
]

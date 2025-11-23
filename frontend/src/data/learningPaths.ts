import { Code2, Database, Server, Cloud } from 'lucide-vue-next'

export const learningPaths = [
  {
    id: 'frontend-vue',
    title: 'Frontend Mastery with Vue.js',
    description: 'Master modern frontend development using Vue 3, TypeScript, and Tailwind CSS.',
    icon: Code2,
    modules: [
      { title: 'Vue 3 Fundamentals', description: 'Components, Props, Events, and Reactivity.' },
      { title: 'State Management with Pinia', description: 'Managing global state efficiently.' },
      { title: 'Routing with Vue Router', description: 'Building Single Page Applications.' },
      { title: 'Styling with Tailwind CSS', description: 'Utility-first CSS framework.' },
    ],
  },
  {
    id: 'backend-go',
    title: 'Backend Engineering with Go',
    description: 'Build robust, high-performance backend services using Go (Golang).',
    icon: Server,
    modules: [
      { title: 'Go Syntax & Primitives', description: 'Variables, Loops, and Data Structures.' },
      { title: 'Concurrency in Go', description: 'Goroutines and Channels.' },
      { title: 'Building APIs with Chi', description: 'RESTful API development.' },
      { title: 'Database Integration', description: 'Working with PostgreSQL and pgx.' },
    ],
  },
  {
    id: 'devops-cloud',
    title: 'DevOps & Cloud Engineering',
    description: 'Deploy and scale applications using Docker, AWS, and CI/CD pipelines.',
    icon: Cloud,
    modules: [
      { title: 'Containerization with Docker', description: 'Building and running containers.' },
      { title: 'AWS Essentials', description: 'EC2, S3, and RDS.' },
      { title: 'CI/CD Pipelines', description: 'Automating deployment workflows.' },
      { title: 'Infrastructure as Code', description: 'Terraform basics.' },
    ],
  },
]

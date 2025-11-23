import { BookOpen, Video, FileText, Github } from 'lucide-vue-next'

export const resources = [
  {
    id: 'vue-docs',
    title: 'Vue.js Documentation',
    description: 'The official Vue.js documentation. The best place to start.',
    type: 'Documentation',
    url: 'https://vuejs.org',
    icon: FileText,
  },
  {
    id: 'go-tour',
    title: 'A Tour of Go',
    description: 'Interactive tour of Go syntax and features.',
    type: 'Interactive',
    url: 'https://go.dev/tour',
    icon: BookOpen,
  },
  {
    id: 'docker-101',
    title: 'Docker for Beginners',
    description: 'Video tutorial series on Docker basics.',
    type: 'Video',
    url: 'https://youtube.com',
    icon: Video,
  },
  {
    id: 'awesome-go',
    title: 'Awesome Go',
    description: 'A curated list of awesome Go frameworks, libraries and software.',
    type: 'Repository',
    url: 'https://github.com/avelino/awesome-go',
    icon: Github,
  },
]

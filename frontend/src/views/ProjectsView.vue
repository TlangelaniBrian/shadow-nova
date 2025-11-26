<script setup lang="ts">
import { onMounted } from 'vue';
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import AppLayout from '@/layouts/AppLayout.vue'
import { ArrowRight } from 'lucide-vue-next'
import { useProjects } from '@/composables/useProjects';
import { useToast } from '@/composables/useToast';

const { projects, isLoading, error, fetchProjects } = useProjects();
const toast = useToast();

const getDifficultyColor = (difficulty: string) => {
  switch (difficulty.toLowerCase()) {
    case 'beginner':
      return 'bg-green-500/10 text-green-500 hover:bg-green-500/20'
    case 'intermediate':
      return 'bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20'
    case 'advanced':
      return 'bg-red-500/10 text-red-500 hover:bg-red-500/20'
    default:
      return 'secondary'
  }
}

onMounted(async () => {
  const result = await fetchProjects();
  if (result.error) {
    toast.showError(result.error);
  }
});
</script>

<template>
  <AppLayout>
    <div class="container mx-auto py-12 max-w-5xl">
      <div class="text-center mb-12 space-y-4">
        <h1 class="text-4xl font-bold">Hands-On Projects</h1>
        <p class="text-xl text-muted-foreground">
          Build real-world applications to reinforce your learning.
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-20">
        <div class="mb-4">
          <svg class="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <p class="text-red-500 font-medium mb-2">Failed to load projects</p>
          <p class="text-sm text-gray-500 mb-4">{{ error.message || 'Unable to connect to the server' }}</p>
        </div>
        <Button @click="fetchProjects">Try Again</Button>
      </div>

      <!-- Projects Grid -->
      <div v-else-if="projects.length > 0" class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <Card
          v-for="project in projects"
          :key="project.id"
          class="flex flex-col hover:shadow-md transition-shadow"
        >
          <CardHeader>
            <div class="flex justify-between items-start mb-4">
              <div class="h-10 w-10 rounded-lg bg-secondary flex items-center justify-center">
                <span class="text-2xl">{{ project.icon || '📦' }}</span>
              </div>
              <Badge :class="getDifficultyColor(project.difficulty)">
                {{ project.difficulty }}
              </Badge>
            </div>
            <CardTitle class="text-xl">{{ project.title }}</CardTitle>
            <CardDescription>{{ project.description }}</CardDescription>
          </CardHeader>
          <CardContent class="flex-1">
            <div class="flex flex-wrap gap-2 mt-2">
              <Badge
                v-for="tech in project.tech_stack"
                :key="tech"
                variant="secondary"
                class="text-xs"
              >
                {{ tech }}
              </Badge>
            </div>
          </CardContent>
          <CardFooter>
            <Button class="w-full" variant="outline">
              View Details <ArrowRight class="ml-2 h-4 w-4" />
            </Button>
          </CardFooter>
        </Card>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-20">
        <p class="text-gray-500 mb-4">No projects available yet</p>
        <p class="text-sm text-gray-400">Check back soon for new projects!</p>
      </div>
    </div>
  </AppLayout>
</template>

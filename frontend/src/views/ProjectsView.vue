<script setup lang="ts">
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
import { projects } from '@/data/projects'
import { ArrowRight } from 'lucide-vue-next'

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

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <Card
          v-for="project in projects"
          :key="project.id"
          class="flex flex-col hover:shadow-md transition-shadow"
        >
          <CardHeader>
            <div class="flex justify-between items-start mb-4">
              <div class="h-10 w-10 rounded-lg bg-secondary flex items-center justify-center">
                <component :is="project.icon" class="h-5 w-5" />
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
                v-for="tech in project.techStack"
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
    </div>
  </AppLayout>
</template>

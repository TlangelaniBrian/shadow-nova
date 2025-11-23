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
import { PlayCircle, Clock, BarChart } from 'lucide-vue-next'

const guides = [
  {
    id: 'vue-setup',
    title: 'Setting up a Vue 3 Project',
    description:
      'Step-by-step guide to initializing a Vue 3 project with Vite, TypeScript, and Tailwind CSS.',
    duration: '10 min',
    level: 'Beginner',
    steps: 5,
  },
  {
    id: 'go-api',
    title: 'Building a REST API with Go',
    description: 'Learn how to create a robust REST API using Go, Chi router, and PostgreSQL.',
    duration: '45 min',
    level: 'Intermediate',
    steps: 12,
  },
  {
    id: 'docker-deploy',
    title: 'Dockerizing Your Application',
    description:
      'A comprehensive guide to creating multi-stage Docker builds for your full-stack app.',
    duration: '30 min',
    level: 'Advanced',
    steps: 8,
  },
]
</script>

<template>
  <AppLayout>
    <div class="container mx-auto py-12 max-w-5xl">
      <div class="text-center mb-12 space-y-4">
        <h1 class="text-4xl font-bold">Interactive Guides</h1>
        <p class="text-xl text-muted-foreground">
          Follow step-by-step tutorials to master new skills.
        </p>
      </div>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <Card
          v-for="guide in guides"
          :key="guide.id"
          class="flex flex-col hover:shadow-md transition-shadow"
        >
          <CardHeader>
            <div class="flex justify-between items-start mb-4">
              <Badge variant="outline">{{ guide.level }}</Badge>
              <div class="flex items-center text-sm text-muted-foreground">
                <Clock class="h-4 w-4 mr-1" /> {{ guide.duration }}
              </div>
            </div>
            <CardTitle class="text-xl">{{ guide.title }}</CardTitle>
            <CardDescription>{{ guide.description }}</CardDescription>
          </CardHeader>
          <CardContent class="flex-1">
            <div class="flex items-center text-sm text-muted-foreground">
              <BarChart class="h-4 w-4 mr-2" /> {{ guide.steps }} Steps
            </div>
          </CardContent>
          <CardFooter>
            <Button class="w-full group">
              Start Guide
              <PlayCircle class="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </AppLayout>
</template>

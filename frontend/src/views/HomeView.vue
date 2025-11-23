<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import AppLayout from '@/layouts/AppLayout.vue'
import { ArrowRight, Code2, Layers, Trophy } from 'lucide-vue-next'

const features = [
  {
    title: 'Structured Learning Paths',
    description:
      'Follow curated roadmaps for Frontend (Vue), Backend (Go), and DevOps (AWS/Docker).',
    icon: Layers,
  },
  {
    title: 'Hands-On Projects',
    description: 'Build real-world applications like this one to build your portfolio.',
    icon: Code2,
  },
  {
    title: 'Progress Tracking',
    description: 'Earn badges and track your milestones as you master new skills.',
    icon: Trophy,
  },
]

const techStack = [
  'Vue.js',
  'TypeScript',
  'Tailwind CSS',
  'Go (Golang)',
  'PostgreSQL',
  'Docker',
  'AWS',
]
</script>

<template>
  <AppLayout>
    <!-- Hero Section -->
    <section class="py-20 md:py-32 text-center space-y-6">
      <Badge variant="secondary" class="mb-4">v1.0 Beta Now Live</Badge>
      <h1 class="text-4xl md:text-6xl font-bold tracking-tight">
        Master the Stack. <br class="hidden md:block" />
        <span class="text-primary">Build Your Future.</span>
      </h1>
      <p class="text-xl text-muted-foreground max-w-2xl mx-auto">
        The ultimate platform for junior developers to learn Vue, Go, and Cloud technologies through
        hands-on projects.
      </p>
      <div class="flex justify-center gap-4 pt-4">
        <Button size="lg"> Start Learning <ArrowRight class="ml-2 h-4 w-4" /> </Button>
        <Button variant="outline" size="lg"> View Projects </Button>
      </div>
    </section>

    <!-- Tech Stack -->
    <section class="py-12 border-y bg-muted/50">
      <div class="container mx-auto px-4 text-center">
        <p class="text-sm font-semibold text-muted-foreground mb-6 uppercase tracking-wider">
          Powered by Modern Tech
        </p>
        <div class="flex flex-wrap justify-center gap-2">
          <Badge v-for="tech in techStack" :key="tech" variant="outline" class="text-md py-1 px-4">
            {{ tech }}
          </Badge>
        </div>
      </div>
    </section>

    <!-- Features Grid -->
    <section class="py-20 container mx-auto px-4">
      <h2 class="text-3xl font-bold text-center mb-12">Everything you need to level up</h2>
      <div class="grid md:grid-cols-3 gap-8">
        <Card
          v-for="feature in features"
          :key="feature.title"
          class="bg-card hover:shadow-lg transition-shadow"
        >
          <CardHeader>
            <component :is="feature.icon" class="h-10 w-10 text-primary mb-4" />
            <CardTitle>{{ feature.title }}</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription class="text-base">
              {{ feature.description }}
            </CardDescription>
          </CardContent>
        </Card>
      </div>
    </section>
  </AppLayout>
</template>

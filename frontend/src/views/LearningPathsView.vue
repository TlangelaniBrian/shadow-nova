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
import { learningPaths } from '@/data/learningPaths'
import { ArrowRight } from 'lucide-vue-next'
</script>

<template>
  <AppLayout>
    <div class="container mx-auto py-12 max-w-5xl">
      <div class="text-center mb-12 space-y-4">
        <h1 class="text-4xl font-bold">Choose Your Path</h1>
        <p class="text-xl text-muted-foreground">
          Select a technology track to start your journey.
        </p>
      </div>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <Card
          v-for="path in learningPaths"
          :key="path.id"
          class="flex flex-col hover:border-primary/50 transition-colors"
        >
          <CardHeader>
            <div class="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
              <component :is="path.icon" class="h-6 w-6 text-primary" />
            </div>
            <CardTitle class="text-xl">{{ path.title }}</CardTitle>
            <CardDescription>{{ path.description }}</CardDescription>
          </CardHeader>
          <CardContent class="flex-1">
            <div class="space-y-2">
              <div
                v-for="(module, index) in path.modules.slice(0, 3)"
                :key="index"
                class="flex items-center text-sm text-muted-foreground"
              >
                <div class="h-1.5 w-1.5 rounded-full bg-primary mr-2"></div>
                {{ module.title }}
              </div>
              <div v-if="path.modules.length > 3" class="text-xs text-muted-foreground pl-3.5">
                + {{ path.modules.length - 3 }} more modules
              </div>
            </div>
          </CardContent>
          <CardFooter>
            <Button class="w-full group" as-child>
              <RouterLink :to="`/learning-paths/${path.id}`">
                Start Path
                <ArrowRight class="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
              </RouterLink>
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </AppLayout>
</template>

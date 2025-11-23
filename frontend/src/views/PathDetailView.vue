<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { learningPaths } from '@/data/learningPaths'
import AppLayout from '@/layouts/AppLayout.vue'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import {
  Accordion,
  AccordionItem,
  AccordionTrigger,
  AccordionContent,
} from '@/components/ui/accordion'
import { Badge } from '@/components/ui/badge'
import { CheckCircle2, Circle, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const pathId = route.params.id as string

const path = computed(() => learningPaths.find(p => p.id === pathId))

// Mock progress state
const completedModules = [0, 1] // First two modules completed
const progress = computed(() => {
  if (!path.value) return 0
  return (completedModules.length / path.value.modules.length) * 100
})
</script>

<template>
  <AppLayout>
    <div v-if="path" class="container mx-auto py-8 max-w-4xl">
      <Button variant="ghost" class="mb-6 pl-0 hover:pl-2 transition-all" as-child>
        <RouterLink to="/learning-paths">
          <ArrowLeft class="mr-2 h-4 w-4" /> Back to Paths
        </RouterLink>
      </Button>

      <div class="flex items-start justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold mb-2">{{ path.title }}</h1>
          <p class="text-muted-foreground">{{ path.description }}</p>
        </div>
        <div class="text-right">
          <Badge variant="outline" class="mb-2">Intermediate</Badge>
        </div>
      </div>

      <div class="bg-card border rounded-lg p-6 mb-8">
        <div class="flex justify-between text-sm mb-2">
          <span class="font-medium">Your Progress</span>
          <span class="text-muted-foreground">{{ Math.round(progress) }}% Completed</span>
        </div>
        <Progress :model-value="progress" class="h-2" />
      </div>

      <h2 class="text-2xl font-bold mb-6">Course Modules</h2>

      <Accordion type="single" collapsible class="w-full space-y-4">
        <AccordionItem
          v-for="(module, index) in path.modules"
          :key="index"
          :value="`item-${index}`"
          class="border rounded-lg px-4"
        >
          <AccordionTrigger class="hover:no-underline py-4">
            <div class="flex items-center gap-4 text-left">
              <div
                :class="completedModules.includes(index) ? 'text-primary' : 'text-muted-foreground'"
              >
                <CheckCircle2 v-if="completedModules.includes(index)" class="h-5 w-5" />
                <Circle v-else class="h-5 w-5" />
              </div>
              <div>
                <div class="font-semibold">{{ module.title }}</div>
                <div class="text-sm text-muted-foreground font-normal">
                  {{ module.description }}
                </div>
              </div>
            </div>
          </AccordionTrigger>
          <AccordionContent class="pl-9 pb-4 text-muted-foreground">
            <p class="mb-4">
              In this module, you will learn the core concepts of {{ module.title }}. Complete the
              hands-on exercises to verify your understanding.
            </p>
            <Button size="sm" :variant="completedModules.includes(index) ? 'outline' : 'default'">
              {{ completedModules.includes(index) ? 'Review Module' : 'Start Module' }}
            </Button>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>

    <div v-else class="container mx-auto py-20 text-center">
      <h1 class="text-2xl font-bold mb-4">Path Not Found</h1>
      <Button as-child>
        <RouterLink to="/learning-paths">Return to Learning Paths</RouterLink>
      </Button>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useLearningPathsStore } from '@/stores/learningPaths'
import { useProgressStore } from '@/stores/progress'
import { useToast } from '@/composables/useToast'
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
const pathsStore = useLearningPathsStore()
const progressStore = useProgressStore()
const toast = useToast()

const pathId = route.params.id as string

onMounted(async () => {
  const pathResult = await pathsStore.fetchPath(pathId)
  if (pathResult.error) {
    toast.showError(pathResult.error)
    return
  }

  const progressResult = await progressStore.fetchPathProgress(pathId)
  if (progressResult.error) {
    toast.showError(progressResult.error)
  }
})

const path = computed(() => pathsStore.currentPath)

const progress = computed(() => {
  if (!path.value || !path.value.modules) return 0

  const totalLessons = path.value.modules.reduce((sum, module) => {
    return sum + (module.lessons?.length || 0)
  }, 0)

  if (totalLessons === 0) return 0

  let completedLessons = 0
  path.value.modules.forEach(module => {
    module.lessons?.forEach(lesson => {
      if (progressStore.isLessonCompleted(lesson.id)) {
        completedLessons++
      }
    })
  })

  return (completedLessons / totalLessons) * 100
})

function isModuleCompleted(moduleIndex: number): boolean {
  const module = path.value?.modules?.[moduleIndex]
  if (!module || !module.lessons) return false

  return module.lessons.every(lesson => progressStore.isLessonCompleted(lesson.id))
}
</script>

<template>
  <AppLayout>
    <!-- Loading State -->
    <div v-if="pathsStore.loading" class="container mx-auto py-20">
      <div class="flex items-center justify-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="pathsStore.error" class="container mx-auto py-20 text-center">
      <div class="mb-4">
        <svg class="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <p class="text-red-500 font-medium mb-2">Failed to load learning path</p>
        <p class="text-sm text-gray-500 mb-4">{{ pathsStore.error }}</p>
      </div>
      <Button @click="pathsStore.fetchPath(pathId)">Try Again</Button>
    </div>

    <!-- Path Content -->
    <div v-else-if="path" class="container mx-auto py-8 max-w-4xl">
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
          <Badge variant="outline" class="mb-2">{{ path.difficulty }}</Badge>
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
                :class="isModuleCompleted(index) ? 'text-primary' : 'text-muted-foreground'"
              >
                <CheckCircle2 v-if="isModuleCompleted(index)" class="h-5 w-5" />
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
            <div v-if="module.lessons" class="space-y-2 mb-4">
              <div v-for="lesson in module.lessons" :key="lesson.id" class="flex items-center gap-2 text-sm">
                <CheckCircle2
                  v-if="progressStore.isLessonCompleted(lesson.id)"
                  class="h-4 w-4 text-primary"
                />
                <Circle v-else class="h-4 w-4" />
                <span>{{ lesson.title }}</span>
                <span class="text-xs text-muted-foreground ml-auto">{{ lesson.duration_minutes }}min</span>
              </div>
            </div>
            <Button size="sm" :variant="isModuleCompleted(index) ? 'outline' : 'default'">
              {{ isModuleCompleted(index) ? 'Review Module' : 'Start Module' }}
            </Button>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>

    <!-- Not Found State -->
    <div v-else class="container mx-auto py-20 text-center">
      <h1 class="text-2xl font-bold mb-4">Path Not Found</h1>
      <Button as-child>
        <RouterLink to="/learning-paths">Return to Learning Paths</RouterLink>
      </Button>
    </div>
  </AppLayout>
</template>

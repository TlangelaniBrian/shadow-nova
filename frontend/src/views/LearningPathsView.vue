<script setup lang="ts">
import { onMounted } from 'vue'
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import AppLayout from '@/layouts/AppLayout.vue'
import { useLearningPathsStore } from '@/stores/learningPaths'
import { useToast } from '@/composables/useToast'
import { ArrowRight } from 'lucide-vue-next'

const pathsStore = useLearningPathsStore()
const toast = useToast()

onMounted(async () => {
  const result = await pathsStore.fetchPaths()
  if (result.error) {
    toast.showError(result.error)
  }
})
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

      <!-- Loading State -->
      <div v-if="pathsStore.loading" class="flex items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="pathsStore.error" class="text-center py-20">
        <div class="mb-4">
          <svg class="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <p class="text-red-500 font-medium mb-2">Failed to load learning paths</p>
          <p class="text-sm text-gray-500 mb-4">{{ pathsStore.error }}</p>
        </div>
        <Button @click="pathsStore.fetchPaths">Try Again</Button>
      </div>

      <!-- Paths Grid -->
      <div v-else-if="pathsStore.paths.length > 0" class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <Card
          v-for="path in pathsStore.paths"
          :key="path.id"
          class="flex flex-col hover:border-primary/50 transition-colors"
        >
          <CardHeader>
            <div class="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
              <span class="text-2xl">📚</span>
            </div>
            <CardTitle class="text-xl">{{ path.title }}</CardTitle>
            <CardDescription>{{ path.description }}</CardDescription>
          </CardHeader>
          <CardContent class="flex-1">
            <div class="space-y-2">
              <div
                v-for="(module, index) in (path.modules || []).slice(0, 3)"
                :key="index"
                class="flex items-center text-sm text-muted-foreground"
              >
                <div class="h-1.5 w-1.5 rounded-full bg-primary mr-2"></div>
                {{ module.title }}
              </div>
              <div v-if="path.modules && path.modules.length > 3" class="text-xs text-muted-foreground pl-3.5">
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

      <!-- Empty State -->
      <div v-else class="text-center py-20">
        <p class="text-gray-500 mb-4">No learning paths available yet</p>
        <p class="text-sm text-gray-400">Check back soon for new paths!</p>
      </div>
    </div>
  </AppLayout>
</template>

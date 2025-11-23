<script setup lang="ts">
import { Card, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import AppLayout from '@/layouts/AppLayout.vue'
import { resources } from '@/data/resources'
import { ExternalLink } from 'lucide-vue-next'

const getTypeColor = (type: string) => {
  switch (type.toLowerCase()) {
    case 'documentation':
      return 'secondary'
    case 'video':
      return 'destructive'
    case 'interactive':
      return 'default'
    default:
      return 'outline'
  }
}
</script>

<template>
  <AppLayout>
    <div class="container mx-auto py-12 max-w-5xl">
      <div class="text-center mb-12 space-y-4">
        <h1 class="text-4xl font-bold">Resource Library</h1>
        <p class="text-xl text-muted-foreground">Curated tutorials, documentation, and tools.</p>
      </div>

      <div class="grid md:grid-cols-2 gap-6">
        <Card
          v-for="resource in resources"
          :key="resource.id"
          class="hover:bg-muted/50 transition-colors"
        >
          <CardHeader>
            <div class="flex justify-between items-start mb-2">
              <component :is="resource.icon" class="h-6 w-6 text-primary" />
              <Badge :variant="getTypeColor(resource.type)">{{ resource.type }}</Badge>
            </div>
            <CardTitle>{{ resource.title }}</CardTitle>
            <CardDescription>{{ resource.description }}</CardDescription>
          </CardHeader>
          <CardFooter>
            <Button variant="ghost" class="w-full justify-between group" as-child>
              <a :href="resource.url" target="_blank" rel="noopener noreferrer">
                Open Resource
                <ExternalLink class="h-4 w-4 group-hover:translate-x-1 transition-transform" />
              </a>
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </AppLayout>
</template>

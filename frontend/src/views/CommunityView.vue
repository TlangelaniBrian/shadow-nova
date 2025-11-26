<script setup lang="ts">
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import AppLayout from '@/layouts/AppLayout.vue'
import { MessageSquare, ThumbsUp, Share2 } from 'lucide-vue-next'

const posts = [
  {
    id: 1,
    author: 'Sarah Chen',
    avatar: 'https://github.com/shadcn.png',
    initials: 'SC',
    title: 'Best practices for Vue 3 Composition API?',
    content:
      'I am migrating a large app from Options API to Composition API. Any tips on structuring composables?',
    tags: ['Vue.js', 'Best Practices'],
    likes: 24,
    comments: 12,
    time: '2h ago',
  },
  {
    id: 2,
    author: 'Mike Ross',
    avatar: '',
    initials: 'MR',
    title: 'Deploying Go app to AWS Lambda',
    content:
      'Has anyone tried deploying a Chi router app to Lambda using the new Go runtime? I am hitting some cold start issues.',
    tags: ['Go', 'AWS', 'Serverless'],
    likes: 15,
    comments: 8,
    time: '5h ago',
  },
  {
    id: 3,
    author: 'Alex T.',
    avatar: '',
    initials: 'AT',
    title: 'Shadow Nova is awesome!',
    content:
      'Just finished the "Task Manager API" project. The step-by-step guide was super helpful. Highly recommend!',
    tags: ['Feedback', 'Success Story'],
    likes: 156,
    comments: 42,
    time: '1d ago',
  },
]
</script>

<template>
  <AppLayout>
    <div class="container mx-auto py-12 max-w-4xl">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-3xl font-bold">Community Forum</h1>
          <p class="text-muted-foreground">Connect, ask questions, and share your journey.</p>
        </div>
        <Button>New Post</Button>
      </div>

      <div class="space-y-6">
        <Card v-for="post in posts" :key="post.id">
          <CardHeader>
            <div class="flex items-center gap-4 mb-2">
              <Avatar>
                <AvatarImage :src="post.avatar" v-if="post.avatar" />
                <AvatarFallback>{{ post.initials }}</AvatarFallback>
              </Avatar>
              <div>
                <div class="font-semibold text-sm">{{ post.author }}</div>
                <div class="text-xs text-muted-foreground">{{ post.time }}</div>
              </div>
            </div>
            <CardTitle class="text-xl">{{ post.title }}</CardTitle>
            <div class="flex gap-2 mt-2">
              <Badge v-for="tag in post.tags" :key="tag" variant="secondary" class="text-xs">
                {{ tag }}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <p class="text-muted-foreground">{{ post.content }}</p>
          </CardContent>
          <CardFooter class="border-t pt-4 flex gap-6 text-muted-foreground text-sm">
            <button class="flex items-center gap-2 hover:text-primary transition-colors">
              <ThumbsUp class="h-4 w-4" /> {{ post.likes }}
            </button>
            <button class="flex items-center gap-2 hover:text-primary transition-colors">
              <MessageSquare class="h-4 w-4" /> {{ post.comments }}
            </button>
            <button class="flex items-center gap-2 hover:text-primary transition-colors ml-auto">
              <Share2 class="h-4 w-4" /> Share
            </button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </AppLayout>
</template>

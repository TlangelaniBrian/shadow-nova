<script setup lang="ts">
import { 
  LayoutDashboard, 
  BookOpen, 
  Rocket, 
  Users, 
  User, 
  Settings 
} from 'lucide-vue-next'
import { useRoute } from 'vue-router'

const route = useRoute()

const menuItems = [
  {
    category: 'Platform',
    items: [
      { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
      { name: 'Learning Paths', icon: BookOpen, path: '/paths' },
      { name: 'Projects', icon: Rocket, path: '/projects' },
      { name: 'Community', icon: Users, path: '/community' },
    ]
  },
  {
    category: 'Account',
    items: [
      { name: 'Profile', icon: User, path: '/profile' },
      { name: 'Settings', icon: Settings, path: '/settings' },
    ]
  }
]
</script>

<template>
  <aside class="w-64 bg-white h-screen fixed left-0 top-0 border-r border-gray-100 flex flex-col z-40">
    <!-- Logo -->
    <div class="p-8">
      <h1 class="text-2xl font-bold text-gray-800">Shadow Nova</h1>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto px-4">
      <div v-for="(section, idx) in menuItems" :key="idx" class="mb-8">
        <h3 class="text-gray-400 text-sm font-medium mb-4 px-4">{{ section.category }}</h3>
        <ul class="space-y-1">
          <li v-for="item in section.items" :key="item.name">
            <RouterLink 
              :to="item.path"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{ 'text-gray-900 font-medium': route.path === item.path }"
            >
              <component :is="item.icon" class="w-5 h-5" />
              <span>{{ item.name }}</span>
            </RouterLink>
          </li>
        </ul>
      </div>
    </nav>
  </aside>
</template>

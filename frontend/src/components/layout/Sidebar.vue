<script setup lang="ts">
import { 
  LayoutDashboard, 
  BookOpen, 
  Rocket, 
  Users, 
  User, 
  Settings,
  X
} from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { storeToRefs } from 'pinia'

const route = useRoute()
const uiStore = useUIStore()
const { isSidebarOpen } = storeToRefs(uiStore)
const { closeSidebar } = uiStore

const menuItems = [
  {
    category: 'Platform',
    items: [
      { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
      { name: 'Learning Paths', icon: BookOpen, path: '/paths' },
      { name: 'Projects', icon: Rocket, path: '/projects' },
      { name: 'Community', icon: Users, path: '/community' },
      { name: 'Profile', icon: User, path: '/profile' },
    ]
  },
  {
    category: 'Account',
    items: [
      { name: 'Settings', icon: Settings, path: '/settings' },
    ]
  }
]
</script>

<template>
  <!-- Mobile Overlay -->
  <div 
    v-if="isSidebarOpen" 
    class="fixed inset-0 bg-black/50 z-40 md:hidden backdrop-blur-sm transition-opacity"
    @click="closeSidebar"
  ></div>

  <!-- Sidebar -->
  <aside 
    class="fixed top-0 left-0 h-screen bg-white border-r border-gray-100 flex flex-col z-50 transition-transform duration-300 ease-in-out w-64"
    :class="[
      isSidebarOpen ? 'translate-x-0' : '-translate-x-full',
      'md:translate-x-0'
    ]"
  >
    <!-- Logo -->
    <div class="p-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-800">Shadow Nova</h1>
      <!-- Mobile Close Button -->
      <button 
        class="md:hidden p-2 text-gray-500 hover:bg-gray-100 rounded-lg"
        @click="closeSidebar"
      >
        <X class="w-5 h-5" />
      </button>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto px-4 py-4">
      <div v-for="(section, idx) in menuItems" :key="idx" class="mb-8">
        <h3 class="text-gray-400 text-xs font-semibold uppercase tracking-wider mb-4 px-4">{{ section.category }}</h3>
        <ul class="space-y-1">
          <li v-for="item in section.items" :key="item.name">
            <RouterLink 
              :to="item.path"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{ 'text-gray-900 font-medium bg-gray-50': route.path === item.path }"
              @click="closeSidebar"
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

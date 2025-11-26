<script setup lang="ts">
import { ref } from 'vue'
import { Search, Bell, ChevronDown, Menu } from 'lucide-vue-next'
import { useUIStore } from '@/stores/ui'

const uiStore = useUIStore()
const { toggleSidebar } = uiStore

const isLive = ref(true)
const user = ref({
  name: "Anddy's Makeover",
  id: "1234567",
  avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Anddy"
})
</script>

<template>
  <header class="h-20 bg-white/50 backdrop-blur-sm fixed top-0 left-0 md:left-64 right-0 lg:right-80 z-30 px-4 md:px-8 flex items-center justify-between transition-all duration-300 ease-in-out">
    <!-- Left Section -->
    <div class="flex items-center gap-4 flex-1 max-w-md">
      <!-- Mobile Menu Button -->
      <button 
        class="md:hidden p-2 -ml-2 text-gray-500 hover:bg-gray-100 rounded-lg"
        @click="toggleSidebar"
      >
        <Menu class="w-6 h-6" />
      </button>

      <!-- Search -->
      <div class="flex-1 relative hidden sm:block">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input 
          type="text" 
          placeholder="Search" 
          class="w-full pl-12 pr-4 py-3 bg-gray-50 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-purple-100"
        >
      </div>
    </div>

    <!-- Right Actions -->
    <div class="flex items-center gap-6">
      <!-- Live Toggle -->
      <div class="flex items-center gap-3">
        <button 
          class="w-12 h-6 rounded-full relative transition-colors duration-300"
          :class="isLive ? 'bg-green-100' : 'bg-gray-200'"
          @click="isLive = !isLive"
        >
          <div 
            class="absolute top-1 w-4 h-4 rounded-full transition-all duration-300"
            :class="[
              isLive ? 'left-7 bg-green-500' : 'left-1 bg-gray-400'
            ]"
          ></div>
        </button>
        <span class="text-sm text-gray-500">Live</span>
      </div>

      <!-- Notifications -->
      <button class="w-10 h-10 rounded-full bg-gray-50 flex items-center justify-center relative hover:bg-gray-100 transition-colors">
        <Bell class="w-5 h-5 text-gray-500" />
        <span class="absolute top-2 right-2 w-2 h-2 bg-purple-500 rounded-full border-2 border-white"></span>
      </button>

      <!-- User Profile -->
      <div class="flex items-center gap-3 pl-4 border-l border-gray-100">
        <img :src="user.avatar" alt="User" class="w-10 h-10 rounded-xl object-cover bg-gray-100" />
        <div class="hidden lg:block">
          <p class="text-sm font-bold text-gray-900">{{ user.name }}</p>
          <p class="text-xs text-gray-400">ID: {{ user.id }}</p>
        </div>
        <ChevronDown class="w-4 h-4 text-gray-400" />
      </div>
    </div>
  </header>
</template>

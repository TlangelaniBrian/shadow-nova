<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Search, Bell, ChevronDown, Menu, LogOut, User as UserIcon, Settings, Moon, Sun } from 'lucide-vue-next'
import { useUIStore } from '@/stores/ui'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'

const uiStore = useUIStore()
const { toggleSidebar } = uiStore
const userStore = useUserStore()
const themeStore = useThemeStore()
const router = useRouter()

const isLive = ref(true)
const showUserMenu = ref(false)

const user = computed(() => ({
  name: userStore.user?.name || 'Developer',
  id: userStore.user?.id || '',
  avatar: userStore.user?.picture || `https://api.dicebear.com/7.x/avataaars/svg?seed=${userStore.user?.name || 'default'}`,
}))

async function handleLogout() {
  try {
    await userStore.logout()
    showUserMenu.value = false
    router.push('/login')
    toast.success('Logged out successfully')
  } catch (error) {
    toast.error('Logout failed')
  }
}

// Close menu when clicking outside
function handleClickOutside(e: MouseEvent) {
  const target = e.target as Element
  if (showUserMenu.value && !target.closest('.user-menu-container')) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <header class="h-20 bg-white/50 dark:bg-gray-800/50 backdrop-blur-sm fixed top-0 left-0 md:left-64 right-0 lg:right-80 z-30 px-4 md:px-8 flex items-center justify-between transition-all duration-300 ease-in-out border-b border-gray-100 dark:border-gray-700">
    <!-- Left Section -->
    <div class="flex items-center gap-4 flex-1 max-w-md">
      <!-- Mobile Menu Button -->
      <button
        class="md:hidden p-2 -ml-2 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg"
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
          class="w-full pl-12 pr-4 py-3 bg-gray-50 dark:bg-gray-900 dark:text-gray-100 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-purple-100 dark:focus:ring-purple-900"
        >
      </div>
    </div>

    <!-- Right Actions -->
    <div class="flex items-center gap-6">
      <!-- Dark Mode Toggle -->
      <button @click="themeStore.toggleDarkMode()"
        class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        :title="themeStore.isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
        <Moon v-if="!themeStore.isDarkMode" class="w-5 h-5 text-gray-600" />
        <Sun v-else class="w-5 h-5 text-yellow-500" />
      </button>

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
        <span class="text-sm text-gray-500 dark:text-gray-400">Live</span>
      </div>

      <!-- Notifications -->
      <button class="w-10 h-10 rounded-full bg-gray-50 dark:bg-gray-900 flex items-center justify-center relative hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
        <Bell class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        <span class="absolute top-2 right-2 w-2 h-2 bg-purple-500 rounded-full border-2 border-white dark:border-gray-800"></span>
      </button>

      <!-- User Profile / Login -->
      <div v-if="!userStore.isAuthenticated" class="flex items-center gap-3 pl-4 border-l border-gray-100 dark:border-gray-700">
        <router-link
          to="/login"
          class="px-6 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-full hover:from-purple-700 hover:to-pink-700 transition-all duration-200 font-medium text-sm"
        >
          Login
        </router-link>
      </div>

      <div v-else class="relative user-menu-container pl-4 border-l border-gray-100 dark:border-gray-700">
        <button
          @click="showUserMenu = !showUserMenu"
          class="flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-800 rounded-xl p-2 -m-2 transition-colors"
        >
          <img :src="user.avatar" alt="User" class="w-10 h-10 rounded-xl object-cover bg-gray-100 dark:bg-gray-800" />
          <div class="hidden lg:block text-left">
            <p class="text-sm font-bold text-gray-900 dark:text-gray-100">{{ user.name }}</p>
            <p class="text-xs text-gray-400">ID: {{ user.id }}</p>
          </div>
          <ChevronDown
            class="w-4 h-4 text-gray-400 transition-transform duration-200"
            :class="{ 'rotate-180': showUserMenu }"
          />
        </button>

        <!-- Dropdown Menu -->
        <transition
          enter-active-class="transition ease-out duration-100"
          enter-from-class="transform opacity-0 scale-95"
          enter-to-class="transform opacity-100 scale-100"
          leave-active-class="transition ease-in duration-75"
          leave-from-class="transform opacity-100 scale-100"
          leave-to-class="transform opacity-0 scale-95"
        >
          <div
            v-if="showUserMenu"
            class="absolute right-0 mt-2 w-56 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-100 dark:border-gray-700 py-2 z-50"
          >
            <router-link
              to="/profile"
              @click="showUserMenu = false"
              class="flex items-center gap-3 px-4 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-purple-50 dark:hover:bg-gray-700 transition-colors"
            >
              <UserIcon class="w-4 h-4" />
              <span class="text-sm">Profile</span>
            </router-link>

            <router-link
              to="/settings"
              @click="showUserMenu = false"
              class="flex items-center gap-3 px-4 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-purple-50 dark:hover:bg-gray-700 transition-colors"
            >
              <Settings class="w-4 h-4" />
              <span class="text-sm">Settings</span>
            </router-link>

            <hr class="my-2 border-gray-100 dark:border-gray-700" />

            <button
              @click="handleLogout"
              class="w-full flex items-center gap-3 px-4 py-2.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
            >
              <LogOut class="w-4 h-4" />
              <span class="text-sm">Logout</span>
            </button>
          </div>
        </transition>
      </div>
    </div>
  </header>
</template>

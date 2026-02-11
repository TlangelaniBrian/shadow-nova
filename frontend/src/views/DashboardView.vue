<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { useProgressStore } from '@/stores/progress'
import { useToast } from '@/composables/useToast'

const userStore = useUserStore()
const progressStore = useProgressStore()
const toast = useToast()

const user = computed(() => userStore.user)

onMounted(async () => {
  const result = await progressStore.fetchStats()
  if (result.error) {
    toast.showError(result.error)
  }
})

const stats = computed(() => {
  if (!progressStore.stats) {
    return [
      { label: 'Courses Completed', value: '0', trend: '+0%', icon: '📚', color: 'bg-purple-100 text-purple-600' },
      { label: 'Projects Built', value: '0', trend: '+0%', icon: '🚀', color: 'bg-blue-100 text-blue-600' },
      { label: 'Hours Learned', value: '0', trend: '+0%', icon: '⏱️', color: 'bg-green-100 text-green-600' },
      { label: 'Rank', value: '#0', trend: '+0', icon: '🏆', color: 'bg-orange-100 text-orange-600' },
    ]
  }

  return [
    {
      label: 'Courses Completed',
      value: progressStore.stats.courses_completed.toString(),
      trend: '+15%',
      icon: '📚',
      color: 'bg-purple-100 text-purple-600'
    },
    {
      label: 'XP Earned',
      value: progressStore.stats.total_xp.toString(),
      trend: '+23%',
      icon: '🚀',
      color: 'bg-blue-100 text-blue-600'
    },
    {
      label: 'Hours Learned',
      value: progressStore.stats.hours_learned.toString(),
      trend: '+8%',
      icon: '⏱️',
      color: 'bg-green-100 text-green-600'
    },
    {
      label: 'Rank',
      value: `#${progressStore.stats.rank}`,
      trend: `+${progressStore.stats.current_streak} streak`,
      icon: '🏆',
      color: 'bg-orange-100 text-orange-600'
    },
  ]
})

const currentPath = ref({
  name: 'Full Stack Development',
  progress: 65,
  modules: 12,
  completed: 8,
})

</script>

<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Welcome back, {{ user?.name?.split(' ')[0] || 'Developer' }}! 👋
        </h2>
        <p class="text-gray-400 dark:text-gray-500 mt-1">Here's your learning progress today</p>
      </div>
    </div>

    <!-- Loading State for Stats -->
    <div v-if="progressStore.loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
    </div>

    <!-- Stats Grid (Bursa Style) -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-white dark:bg-gray-800 rounded-3xl p-6 border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition-all"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 rounded-2xl flex items-center justify-center text-2xl" :class="stat.color">
            {{ stat.icon }}
          </div>
          <span class="text-green-500 dark:text-green-400 text-sm font-medium bg-green-50 dark:bg-green-900/30 px-2 py-1 rounded-full">{{ stat.trend }}</span>
        </div>
        <p class="text-gray-400 dark:text-gray-500 text-sm mb-1">{{ stat.label }}</p>
        <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ stat.value }}</p>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Current Learning Path (Styled like Bursa Main Chart) -->
      <div class="lg:col-span-2 bg-white dark:bg-gray-800 rounded-3xl p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-bold text-gray-900 dark:text-gray-100">Current Learning Path</h3>
          <button class="text-purple-600 dark:text-purple-400 text-sm font-medium hover:text-purple-700 dark:hover:text-purple-300">View Details</button>
        </div>
        
        <div class="bg-gradient-to-br from-purple-600 to-indigo-600 rounded-2xl p-8 text-white relative overflow-hidden">
          <!-- Decorative circles -->
          <div class="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -translate-y-1/2 translate-x-1/2 blur-3xl"></div>
          <div class="absolute bottom-0 left-0 w-32 h-32 bg-white/10 rounded-full translate-y-1/2 -translate-x-1/2 blur-2xl"></div>

          <div class="relative z-10">
            <div class="flex items-start justify-between mb-8">
              <div>
                <span class="bg-white/20 text-white text-xs px-3 py-1 rounded-full backdrop-blur-sm">In Progress</span>
                <h4 class="text-2xl font-bold mt-4 mb-2">{{ currentPath.name }}</h4>
                <p class="text-white/80">{{ currentPath.completed }} of {{ currentPath.modules }} modules completed</p>
              </div>
              <div class="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center backdrop-blur-sm">
                <span class="text-2xl">💻</span>
              </div>
            </div>

            <!-- Progress Bar -->
            <div class="space-y-2">
              <div class="flex justify-between text-sm font-medium">
                <span>Progress</span>
                <span>{{ currentPath.progress }}%</span>
              </div>
              <div class="h-3 bg-black/20 rounded-full overflow-hidden backdrop-blur-sm">
                <div 
                  class="h-full bg-white rounded-full transition-all duration-1000 ease-out"
                  :style="{ width: `${currentPath.progress}%` }"
                ></div>
              </div>
            </div>

            <button class="mt-8 w-full bg-white text-purple-600 font-bold py-3 rounded-xl hover:bg-purple-50 transition-colors">
              Continue Learning
            </button>
          </div>
        </div>
      </div>

      <!-- Recommended/Next Steps (Styled like Bursa Side Cards) -->
      <div class="bg-white dark:bg-gray-800 rounded-3xl p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
        <h3 class="text-lg font-bold text-gray-900 dark:text-gray-100 mb-6">Next Up</h3>

        <div class="space-y-4">
          <div class="p-4 rounded-2xl bg-gray-50 dark:bg-gray-900 border border-gray-100 dark:border-gray-700 hover:border-purple-200 dark:hover:border-purple-700 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600 dark:text-orange-400 group-hover:scale-110 transition-transform">
                ⚡
              </div>
              <div>
                <h5 class="font-bold text-gray-900 dark:text-gray-100 text-sm">Go Concurrency</h5>
                <p class="text-xs text-gray-400 dark:text-gray-500">Advanced Module</p>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-gray-50 dark:bg-gray-900 border border-gray-100 dark:border-gray-700 hover:border-purple-200 dark:hover:border-purple-700 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform">
                🐳
              </div>
              <div>
                <h5 class="font-bold text-gray-900 dark:text-gray-100 text-sm">Docker Mastery</h5>
                <p class="text-xs text-gray-400 dark:text-gray-500">DevOps Path</p>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-gray-50 dark:bg-gray-900 border border-gray-100 dark:border-gray-700 hover:border-purple-200 dark:hover:border-purple-700 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-pink-100 dark:bg-pink-900/30 flex items-center justify-center text-pink-600 dark:text-pink-400 group-hover:scale-110 transition-transform">
                🎨
              </div>
              <div>
                <h5 class="font-bold text-gray-900 dark:text-gray-100 text-sm">UI Design Systems</h5>
                <p class="text-xs text-gray-400 dark:text-gray-500">Frontend Path</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

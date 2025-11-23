<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref<any>(null)

// Restoring Shadow Nova Stats Data
const stats = ref([
  { label: 'Courses Completed', value: '12', trend: '+15%', icon: '📚', color: 'bg-purple-100 text-purple-600' },
  { label: 'Projects Built', value: '8', trend: '+23%', icon: '🚀', color: 'bg-blue-100 text-blue-600' },
  { label: 'Hours Learned', value: '147', trend: '+8%', icon: '⏱️', color: 'bg-green-100 text-green-600' },
  { label: 'Rank', value: '#42', trend: '+5', icon: '🏆', color: 'bg-orange-100 text-orange-600' },
])

const currentPath = ref({
  name: 'Full Stack Development',
  progress: 65,
  modules: 12,
  completed: 8,
})

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    user.value = JSON.parse(userData)
  }
})
</script>

<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">
          Welcome back, {{ user?.name?.split(' ')[0] || 'Developer' }}! 👋
        </h2>
        <p class="text-gray-400 mt-1">Here's your learning progress today</p>
      </div>
    </div>

    <!-- Stats Grid (Bursa Style) -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-white rounded-3xl p-6 border border-gray-100 shadow-sm hover:shadow-md transition-all"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 rounded-2xl flex items-center justify-center text-2xl" :class="stat.color">
            {{ stat.icon }}
          </div>
          <span class="text-green-500 text-sm font-medium bg-green-50 px-2 py-1 rounded-full">{{ stat.trend }}</span>
        </div>
        <p class="text-gray-400 text-sm mb-1">{{ stat.label }}</p>
        <p class="text-2xl font-bold text-gray-900">{{ stat.value }}</p>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Current Learning Path (Styled like Bursa Main Chart) -->
      <div class="lg:col-span-2 bg-white rounded-3xl p-8 border border-gray-100 shadow-sm">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-bold text-gray-900">Current Learning Path</h3>
          <button class="text-purple-600 text-sm font-medium hover:text-purple-700">View Details</button>
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
      <div class="bg-white rounded-3xl p-8 border border-gray-100 shadow-sm">
        <h3 class="text-lg font-bold text-gray-900 mb-6">Next Up</h3>
        
        <div class="space-y-4">
          <div class="p-4 rounded-2xl bg-gray-50 border border-gray-100 hover:border-purple-200 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-orange-100 flex items-center justify-center text-orange-600 group-hover:scale-110 transition-transform">
                ⚡
              </div>
              <div>
                <h5 class="font-bold text-gray-900 text-sm">Go Concurrency</h5>
                <p class="text-xs text-gray-400">Advanced Module</p>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-gray-50 border border-gray-100 hover:border-purple-200 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-blue-100 flex items-center justify-center text-blue-600 group-hover:scale-110 transition-transform">
                🐳
              </div>
              <div>
                <h5 class="font-bold text-gray-900 text-sm">Docker Mastery</h5>
                <p class="text-xs text-gray-400">DevOps Path</p>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-gray-50 border border-gray-100 hover:border-purple-200 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-pink-100 flex items-center justify-center text-pink-600 group-hover:scale-110 transition-transform">
                🎨
              </div>
              <div>
                <h5 class="font-bold text-gray-900 text-sm">UI Design Systems</h5>
                <p class="text-xs text-gray-400">Frontend Path</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

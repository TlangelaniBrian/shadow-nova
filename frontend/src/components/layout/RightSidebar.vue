<script setup lang="ts">
import { ref } from 'vue'
import { ChevronRight, Trophy, Code, BookOpen } from 'lucide-vue-next'

const activeTab = ref('Activity')

const activities = [
  {
    id: 1,
    icon: BookOpen,
    title: 'Completed Module',
    desc: 'Vue Composition API',
    time: '2h ago',
    color: 'bg-purple-100 text-purple-600'
  },
  {
    id: 2,
    icon: Code,
    title: 'Started Project',
    desc: 'Go Microservices',
    time: '5h ago',
    color: 'bg-blue-100 text-blue-600'
  },
  {
    id: 3,
    icon: Trophy,
    title: 'Earned Badge',
    desc: 'Backend Master',
    time: '1d ago',
    color: 'bg-yellow-100 text-yellow-600'
  },
  {
    id: 4,
    icon: BookOpen,
    title: 'Completed Quiz',
    desc: 'TypeScript Basics',
    time: '2d ago',
    color: 'bg-green-100 text-green-600'
  }
]
</script>

<template>
  <aside class="w-80 bg-white h-screen fixed right-0 top-0 border-l border-gray-100 flex flex-col z-40 p-6">
    <!-- Tabs -->
    <div class="flex items-center gap-8 mb-8">
      <button 
        v-for="tab in ['Activity', 'Community']" 
        :key="tab"
        class="text-lg font-medium transition-colors relative pb-2"
        :class="activeTab === tab ? 'text-gray-900' : 'text-gray-400'"
        @click="activeTab = tab"
      >
        {{ tab }}
        <div 
          v-if="activeTab === tab" 
          class="absolute bottom-0 left-0 w-full h-1 bg-gray-900 rounded-full"
        ></div>
      </button>
    </div>

    <!-- Activity List -->
    <div class="flex-1 overflow-y-auto space-y-4 pr-2">
      <div 
        v-for="item in activities" 
        :key="item.id"
        class="p-4 bg-white rounded-2xl shadow-[0_2px_10px_rgba(0,0,0,0.02)] border border-gray-50 hover:shadow-md transition-shadow cursor-pointer group"
      >
        <div class="flex items-start gap-4">
          <div 
            class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
            :class="item.color"
          >
            <component :is="item.icon" class="w-5 h-5" />
          </div>
          
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-1">
              <h4 class="text-sm font-bold text-gray-900">{{ item.title }}</h4>
              <span class="text-[10px] text-gray-400">{{ item.time }}</span>
            </div>
            <p class="text-xs text-gray-500 truncate">{{ item.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>

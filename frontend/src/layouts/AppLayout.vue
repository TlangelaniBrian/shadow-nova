<script setup lang="ts">
import Sidebar from '@/components/layout/Sidebar.vue'
import Header from '@/components/layout/Header.vue'
import RightSidebar from '@/components/layout/RightSidebar.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 font-sans flex flex-col md:flex-row">
    <!-- Mobile Header/Nav Trigger could go here or in Header component -->
    
    <Sidebar />
    
    <div class="flex-1 flex flex-col min-h-screen transition-all duration-300 ease-in-out"
         :class="['md:pl-64', 'lg:pr-80']">
      
      <Header />

      <!-- Main Content Area -->
      <main class="flex-1 pt-20 px-4 md:px-8 pb-8 w-full max-w-[100vw] overflow-x-hidden">
        <div class="max-w-7xl mx-auto w-full">
          <slot />
        </div>
      </main>
    </div>

    <RightSidebar class="hidden lg:block" />
  </div>
</template>

<style>
/* Global font adjustment if needed */
body {
  font-family: 'Inter', sans-serif;
}
</style>

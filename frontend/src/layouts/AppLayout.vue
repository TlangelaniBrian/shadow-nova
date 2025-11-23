<script setup lang="ts">
import Sidebar from '@/components/layout/Sidebar.vue'
import Header from '@/components/layout/Header.vue'
import RightSidebar from '@/components/layout/RightSidebar.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 font-sans">
    <Sidebar />
    <Header />
    <RightSidebar />

    <!-- Main Content Area -->
    <main class="pl-64 pr-80 pt-20 min-h-screen">
      <div class="p-8 max-w-7xl mx-auto">
        <slot />
      </div>
    </main>
  </div>
</template>

<style>
/* Global font adjustment if needed */
body {
  font-family: 'Inter', sans-serif;
}
</style>

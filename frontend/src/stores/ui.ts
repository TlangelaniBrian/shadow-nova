import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useUIStore = defineStore('ui', () => {
    const isSidebarOpen = ref(false);
    const isRightSidebarOpen = ref(false);

    function toggleSidebar() {
        isSidebarOpen.value = !isSidebarOpen.value;
    }

    function closeSidebar() {
        isSidebarOpen.value = false;
    }

    function toggleRightSidebar() {
        isRightSidebarOpen.value = !isRightSidebarOpen.value;
    }

    return {
        isSidebarOpen,
        isRightSidebarOpen,
        toggleSidebar,
        closeSidebar,
        toggleRightSidebar,
    };
});

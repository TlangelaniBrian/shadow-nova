<template>
  <div class="bg-white rounded-3xl p-8 border border-gray-100 shadow-sm">
    <h3 class="text-lg font-bold text-gray-900 mb-6">Account Settings</h3>

    <div class="space-y-4">
      <div class="flex items-center justify-between py-3 border-b border-gray-100">
        <div>
          <p class="font-medium text-gray-900">Email Notifications</p>
          <p class="text-sm text-gray-400">Receive updates about your progress</p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            :checked="emailNotifications"
            @change="handleEmailNotificationsChange"
            class="sr-only peer"
          >
          <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
        </label>
      </div>

      <div class="flex items-center justify-between py-3">
        <div>
          <p class="font-medium text-gray-900">Dark Mode</p>
          <p class="text-sm text-gray-400">Switch to dark theme</p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            :checked="darkMode"
            @change="handleDarkModeChange"
            class="sr-only peer"
          >
          <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  emailNotifications?: boolean;
  darkMode?: boolean;
}

withDefaults(defineProps<Props>(), {
  emailNotifications: true,
  darkMode: false,
});

interface Emits {
  (e: 'update:emailNotifications', value: boolean): void;
  (e: 'update:darkMode', value: boolean): void;
}

const emit = defineEmits<Emits>();

const handleEmailNotificationsChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:emailNotifications', target.checked);
};

const handleDarkModeChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:darkMode', target.checked);
};
</script>

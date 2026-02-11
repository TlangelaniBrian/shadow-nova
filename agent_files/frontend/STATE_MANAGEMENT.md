# State Management Guide

This document describes the state management architecture for Shadow Nova frontend using Pinia.

## Overview

Shadow Nova uses [Pinia](https://pinia.vuejs.org/) as its state management library. Pinia provides:
- Type safety with TypeScript
- Composition API syntax
- Devtools integration
- Plugin system for persistence

## Store Architecture

### Available Stores

1. **User Store** (`stores/user.ts`)
   - Manages authentication state
   - Handles user profile data
   - Manages login/logout flows

2. **Projects Store** (`stores/projects.ts`)
   - Manages project listings
   - Handles project submissions
   - Tracks submission status

3. **Learning Paths Store** (`stores/learningPaths.ts`)
   - Manages learning path data
   - Organizes paths by difficulty
   - Handles current path navigation

4. **Progress Store** (`stores/progress.ts`)
   - Tracks lesson completion
   - Manages user statistics
   - Calculates path progress

5. **UI Store** (`stores/ui.ts`)
   - Manages UI state (sidebars, modals)
   - Handles responsive behavior

## When to Use Stores vs Composables

### Use Stores When:
- Data needs to be shared across multiple components
- State needs to persist across route changes
- You need centralized state management
- State requires complex computed properties
- Data comes from API endpoints

**Example:** User authentication, project listings, learning path data

### Use Composables When:
- Logic is component-specific
- No state sharing is needed
- Providing utility functions
- Creating reusable UI logic
- Handling side effects

**Example:** Toast notifications, form validation, error handling

## Store Pattern

All stores follow a consistent pattern:

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Result } from '@/types/errors'
import { transformAxiosError, success, failure } from '@/types/errors'

export const useExampleStore = defineStore('example', () => {
    // State
    const data = ref<DataType[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)

    // Computed/Getters
    const computedValue = computed(() => {
        // Complex calculations
    })

    // Actions
    async function fetchData(): Promise<Result<DataType[]>> {
        loading.value = true
        error.value = null
        try {
            const response = await api.getData()
            data.value = response.data
            return success(response.data)
        } catch (err) {
            const appError = transformAxiosError(err)
            error.value = appError.message
            return failure(appError)
        } finally {
            loading.value = false
        }
    }

    return {
        // State
        data,
        loading,
        error,
        // Computed
        computedValue,
        // Actions
        fetchData
    }
}, {
    // Persistence config (optional)
    // persist: {
    //     paths: ['data']
    // }
})
```

## Error Handling

All stores use the `Result<T>` type for async operations:

```typescript
// In component
const store = useExampleStore()
const result = await store.fetchData()

if (result.error) {
    toast.showError(result.error)
} else {
    // Success - data is in store.data
}
```

This pattern ensures:
- Type-safe error handling
- Consistent error messages
- No thrown exceptions to catch
- Explicit success/failure paths

## Using Stores in Components

```vue
<script setup lang="ts">
import { onMounted } from 'vue'
import { useProjectsStore } from '@/stores/projects'
import { useToast } from '@/composables/useToast'

const projectsStore = useProjectsStore()
const toast = useToast()

onMounted(async () => {
    const result = await projectsStore.fetchProjects()
    if (result.error) {
        toast.showError(result.error)
    }
})
</script>

<template>
    <div v-if="projectsStore.loading">Loading...</div>
    <div v-else-if="projectsStore.error">Error: {{ projectsStore.error }}</div>
    <div v-else>
        <div v-for="project in projectsStore.projects" :key="project.id">
            {{ project.title }}
        </div>
    </div>
</template>
```

## Persistence Strategy

State persistence is handled by `pinia-plugin-persistedstate`.

### Installation

```bash
npm install pinia-plugin-persistedstate
```

Then uncomment the plugin configuration in `main.ts`:

```typescript
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
```

### Configuring Persistence

In each store, uncomment the persist configuration:

```typescript
export const useProjectsStore = defineStore('projects', () => {
    // ... store implementation
}, {
    persist: {
        paths: ['projects'] // Only persist specific state
    }
})
```

### What to Persist

**Persist:**
- User preferences
- Cached data that's expensive to fetch
- Progress tracking
- UI state (theme, sidebar state)

**Don't Persist:**
- Loading states
- Error messages
- Temporary UI state
- Sensitive data (use HttpOnly cookies instead)

## Testing Stores

```typescript
import { setActivePinia, createPinia } from 'pinia'
import { useProjectsStore } from '@/stores/projects'

describe('Projects Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('fetches projects successfully', async () => {
        const store = useProjectsStore()
        const result = await store.fetchProjects()

        expect(result.error).toBeNull()
        expect(store.projects.length).toBeGreaterThan(0)
    })

    it('handles fetch errors', async () => {
        const store = useProjectsStore()
        // Mock API to fail
        const result = await store.fetchProjects()

        expect(result.error).toBeTruthy()
        expect(store.error).toBeTruthy()
    })
})
```

## Migration from Composables

When migrating from composables to stores:

1. Create the store following the pattern above
2. Move state and logic to the store
3. Update components to use the store
4. Remove the old composable file
5. Update imports across the codebase

Example migration:

```typescript
// Before (composable)
import { useProjects } from '@/composables/useProjects'
const { projects, loading, error, fetchProjects } = useProjects()

// After (store)
import { useProjectsStore } from '@/stores/projects'
const projectsStore = useProjectsStore()
// Access: projectsStore.projects, projectsStore.loading, etc.
```

## Best Practices

1. **Single Responsibility**: Each store manages one domain
2. **Type Safety**: Always define TypeScript interfaces
3. **Error Handling**: Use Result<T> pattern consistently
4. **Loading States**: Track loading for better UX
5. **Computed Properties**: Use for derived state
6. **Actions Return Results**: Always return Result<T> from async actions
7. **Clear Error Messages**: Set human-readable error messages
8. **Selective Persistence**: Only persist what's necessary
9. **Reset on Logout**: Clear sensitive data on logout
10. **Consistent Naming**: Use clear, descriptive names

## Common Patterns

### Optimistic Updates

```typescript
async function updateItem(id: number, data: UpdateData) {
    // Optimistically update UI
    const original = items.value.find(i => i.id === id)
    const index = items.value.findIndex(i => i.id === id)
    items.value[index] = { ...original, ...data }

    try {
        await api.update(id, data)
        return success(undefined)
    } catch (err) {
        // Revert on error
        items.value[index] = original
        return failure(transformAxiosError(err))
    }
}
```

### Polling for Updates

```typescript
let pollInterval: number | null = null

function startPolling(intervalMs: number = 5000) {
    pollInterval = window.setInterval(() => {
        fetchData()
    }, intervalMs)
}

function stopPolling() {
    if (pollInterval) {
        clearInterval(pollInterval)
        pollInterval = null
    }
}
```

### Pagination

```typescript
const currentPage = ref(1)
const pageSize = ref(20)
const totalPages = ref(0)

const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    return items.value.slice(start, start + pageSize.value)
})
```

## Resources

- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vue 3 Composition API](https://vuejs.org/api/composition-api-setup.html)
- [TypeScript with Vue](https://vuejs.org/guide/typescript/overview.html)
- [Pinia Plugin Persistedstate](https://github.com/prazdevs/pinia-plugin-persistedstate)

# Store Setup Completion Checklist

## ✅ Completed Tasks

### 1. Created New Stores
- ✅ `src/stores/projects.ts` - Projects and submissions management
- ✅ `src/stores/learningPaths.ts` - Learning paths with difficulty grouping
- ✅ `src/stores/progress.ts` - User progress tracking and stats

### 2. Created API Modules
- ✅ `src/api/paths.ts` - Learning paths API client
- ✅ `src/api/progress.ts` - Progress tracking API client

### 3. Updated Views
- ✅ `src/views/ProjectsView.vue` - Now uses store instead of composable

### 4. Enhanced Main Configuration
- ✅ `src/main.ts` - Added comments for pinia-plugin-persistedstate

### 5. Added Persistence Configuration
- ✅ All stores configured with commented persistence settings
- ✅ Ready to enable after plugin installation

### 6. Documentation
- ✅ `STATE_MANAGEMENT.md` - Comprehensive state management guide

## 🔄 Remaining Tasks

### 1. Install Persistence Plugin
```bash
cd frontend
npm install pinia-plugin-persistedstate
```

### 2. Enable Persistence
After installation, uncomment in `src/main.ts`:
```typescript
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
// ...
pinia.use(piniaPluginPersistedstate)
```

Then uncomment persistence config in each store:
```typescript
// In stores/projects.ts, stores/learningPaths.ts, stores/progress.ts
persist: {
    paths: ['projects'] // Adjust per store
}
```

### 3. Update Remaining Views
The following views need to be updated to use stores:

**LearningPathsView.vue**
- Currently uses static data from `@/data/learningPaths`
- Should use `useLearningPathsStore()`
```typescript
import { useLearningPathsStore } from '@/stores/learningPaths'
const pathsStore = useLearningPathsStore()
onMounted(() => pathsStore.fetchPaths())
```

**PathDetailView.vue**
- Currently uses static data from `@/data/learningPaths`
- Should use `useLearningPathsStore()` and `useProgressStore()`
```typescript
import { useLearningPathsStore } from '@/stores/learningPaths'
import { useProgressStore } from '@/stores/progress'

const pathsStore = useLearningPathsStore()
const progressStore = useProgressStore()

onMounted(async () => {
    await pathsStore.fetchPath(route.params.id as string)
    await progressStore.fetchPathProgress(route.params.id as string)
})
```

**DashboardView.vue**
- May need `useProgressStore()` for stats display

### 4. Remove Unused Composables
After confirming stores work correctly, consider removing:
- `src/composables/useProjects.ts` (replaced by projects store)

**Keep these composables:**
- `useAuth.ts` - Auth logic (could migrate to user store if desired)
- `useToast.ts` - UI utility
- `useErrorHandler.ts` - Error handling utility
- `useCSRF.ts` - Security utility
- `useApi.ts` - API wrapper utility

### 5. Backend API Routes to Verify
Ensure these routes exist in the backend:

**Projects:**
- ✅ GET `/api/projects` - List projects
- ✅ POST `/api/projects/submit` - Submit project

**Learning Paths:**
- ⚠️ GET `/api/learning-paths` - List paths (verify route exists)
- ⚠️ GET `/api/learning-paths/:id` - Get path details

**Progress:**
- ⚠️ POST `/api/progress` - Update progress
- ⚠️ GET `/api/progress/stats` - Get user stats
- ⚠️ GET `/api/progress/paths/:id` - Get path progress

### 6. Testing Checklist
- [ ] Test projects store fetches data correctly
- [ ] Test learning paths store fetches data correctly
- [ ] Test progress store updates correctly
- [ ] Test error handling in all stores
- [ ] Test persistence after page reload
- [ ] Test loading states display correctly
- [ ] Test error states display correctly

### 7. Type Definitions
All TypeScript interfaces are defined in:
- `src/api/projects.ts` - Project types
- `src/api/paths.ts` - Learning path types
- `src/api/progress.ts` - Progress types
- `src/types/errors.ts` - Error types

## 📋 Store Usage Examples

### Using Projects Store
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

async function handleSubmit(data) {
    const result = await projectsStore.submitProject(data)
    if (result.error) {
        toast.showError(result.error)
    } else {
        toast.showSuccess('Project submitted successfully!')
    }
}
</script>

<template>
    <div v-if="projectsStore.loading">Loading...</div>
    <div v-else-if="projectsStore.error">{{ projectsStore.error }}</div>
    <div v-else>
        <div v-for="project in projectsStore.projects" :key="project.id">
            {{ project.title }}
        </div>
    </div>
</template>
```

### Using Learning Paths Store
```vue
<script setup lang="ts">
import { onMounted } from 'vue'
import { useLearningPathsStore } from '@/stores/learningPaths'

const pathsStore = useLearningPathsStore()

onMounted(() => pathsStore.fetchPaths())
</script>

<template>
    <div v-for="(paths, difficulty) in pathsStore.pathsByDifficulty" :key="difficulty">
        <h2>{{ difficulty }}</h2>
        <div v-for="path in paths" :key="path.id">
            {{ path.title }}
        </div>
    </div>
</template>
```

### Using Progress Store
```vue
<script setup lang="ts">
import { useProgressStore } from '@/stores/progress'

const progressStore = useProgressStore()

async function markComplete(lessonId: number) {
    const result = await progressStore.updateProgress(lessonId, true)
    if (result.error) {
        // Handle error
    }
}
</script>

<template>
    <div v-if="progressStore.stats">
        <p>XP: {{ progressStore.stats.total_xp }}</p>
        <p>Streak: {{ progressStore.stats.current_streak }}</p>
    </div>
</template>
```

## 🎯 Benefits of This Architecture

1. **Type Safety**: Full TypeScript support with proper types
2. **Error Handling**: Consistent Result<T> pattern
3. **Centralized State**: Single source of truth for each domain
4. **Persistence Ready**: Easy to enable localStorage persistence
5. **Testable**: Stores can be tested independently
6. **DevTools**: Full Pinia devtools integration
7. **Composable**: Easy to reuse across components
8. **Maintainable**: Clear separation of concerns

## 📚 Documentation

See `STATE_MANAGEMENT.md` for:
- Complete architecture overview
- Store vs Composable guidelines
- Testing strategies
- Migration patterns
- Best practices
- Common patterns (polling, pagination, etc.)

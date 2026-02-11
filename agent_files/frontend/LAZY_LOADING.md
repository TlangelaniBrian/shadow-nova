# Lazy Loading Implementation Guide

## Overview

This document explains the lazy loading implementation for Shadow Nova's Vue.js frontend, which significantly reduces initial bundle size and improves application load time.

## What is Lazy Loading?

Lazy loading (also known as code splitting) is a technique where route components are loaded on-demand rather than bundled together in a single JavaScript file. When a user navigates to a route, only then is the component code downloaded.

## Benefits

### Bundle Size Reduction
- **Initial bundle**: Only includes core Vue, router, and shared dependencies
- **Route chunks**: Each view is split into separate chunks (typically 10-50KB each)
- **On-demand loading**: Components load only when the user navigates to them

### Performance Improvements
- **Faster initial load**: Smaller initial bundle means faster Time to Interactive (TTI)
- **Better caching**: Individual chunks can be cached separately
- **Progressive loading**: Users only download what they need
- **Improved Core Web Vitals**: Better FCP (First Contentful Paint) and LCP (Largest Contentful Paint)

### Before Lazy Loading
```
main-bundle.js: 500-800KB (all views included)
```

### After Lazy Loading
```
vendor-vue.js: 150KB (Vue core)
vendor-ui.js: 100KB (UI components)
login.js: 25KB (Login view)
dashboard.js: 45KB (Dashboard view)
... (other routes loaded on-demand)
```

## Implementation Details

### 1. Router Configuration

**Before (Static Imports):**
```typescript
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'

const routes = [
  {
    path: '/login',
    component: LoginView
  }
]
```

**After (Dynamic Imports):**
```typescript
const routes = [
  {
    path: '/login',
    component: () => import(/* webpackChunkName: "login" */ '../views/LoginView.vue')
  }
]
```

### 2. Chunk Naming

The `webpackChunkName` comment gives meaningful names to chunks:
```typescript
// Creates "login-[hash].js" instead of "chunk-[hash].js"
component: () => import(/* webpackChunkName: "login" */ '../views/LoginView.vue')
```

### 3. Loading States

A loading component displays during route transitions:
- `RouteLoading.vue`: Spinner component
- `App.vue`: Manages loading state with router guards

### 4. Vendor Chunking

Vite configuration splits vendor dependencies:
```typescript
manualChunks: {
  'vendor-vue': ['vue', 'vue-router', 'pinia'],
  'vendor-ui': ['radix-vue', 'lucide-vue-next'],
  'vendor-utils': ['axios', 'jwt-decode'],
}
```

## Bundle Analysis

Run the bundle analyzer to see chunk sizes:

```bash
./scripts/analyze-bundle.sh
```

This script will:
1. Build the production bundle
2. Display individual chunk sizes
3. Show total bundle size
4. Count number of chunks created

## How Route Chunks are Created

### Automatic Code Splitting

Vite automatically creates separate chunks for each dynamic import:

```typescript
// Each import creates a separate chunk
import('../views/DashboardView.vue')  // -> dashboard-[hash].js
import('../views/ProfileView.vue')     // -> profile-[hash].js
```

### Chunk Dependencies

Each route chunk includes:
- Component code
- Component-specific dependencies
- Imported child components (if not lazy loaded)

Shared dependencies are moved to vendor chunks to avoid duplication.

## Best Practices

### 1. Lazy Load All Routes

Every route should use dynamic imports:
```typescript
// Good
component: () => import('../views/DashboardView.vue')

// Bad - increases initial bundle
import DashboardView from '../views/DashboardView.vue'
component: DashboardView
```

### 2. Strategic Vendor Chunking

Group related vendor dependencies:
```typescript
manualChunks: {
  'vendor-vue': ['vue', 'vue-router', 'pinia'],  // Core framework
  'vendor-ui': ['radix-vue', 'lucide-vue-next'], // UI libraries
  'vendor-utils': ['axios', 'jwt-decode'],        // Utilities
}
```

### 3. Preload Critical Routes

For routes users are likely to visit soon:
```typescript
// In a component
import { useRouter } from 'vue-router'

const router = useRouter()

// Preload dashboard after login
router.beforeResolve((to) => {
  if (to.name === 'login') {
    import('../views/DashboardView.vue') // Preload in background
  }
})
```

### 4. Keep Component Imports Dynamic

Only lazy load at the route level:
```vue
<!-- Good - child components are bundled with parent -->
<script setup>
import UserCard from '@/components/UserCard.vue'
</script>

<!-- Unnecessary - don't lazy load small components -->
<script setup>
const UserCard = defineAsyncComponent(() =>
  import('@/components/UserCard.vue')
)
</script>
```

### 5. Monitor Chunk Sizes

Keep route chunks under 100KB:
- If a route chunk is too large, consider splitting complex child components
- Move heavy dependencies to separate dynamic imports
- Use virtual scrolling for long lists

## Preloading Critical Routes

Preload routes users are likely to visit:

```typescript
// After successful login, preload dashboard
router.afterEach((to) => {
  if (to.name === 'login') {
    // Preload likely next routes
    setTimeout(() => {
      import('../views/DashboardView.vue')
      import('../views/ProfileView.vue')
    }, 1000)
  }
})
```

## Troubleshooting

### Slow Route Loads

**Symptoms:**
- Long delay when navigating to a route
- Blank screen during navigation

**Solutions:**
1. Check network tab - chunk might be large
2. Implement route preloading
3. Add loading states
4. Optimize component dependencies

### Chunk Size Warnings

**Symptoms:**
```
(!) Some chunks are larger than 500 KiB
```

**Solutions:**
1. Split large components
2. Move heavy libraries to dynamic imports
3. Check for duplicate dependencies
4. Increase `chunkSizeWarningLimit` if intentional

### Chunks Not Created

**Symptoms:**
- All code still in one bundle
- No separate route chunks

**Solutions:**
1. Verify dynamic imports: `() => import()`
2. Check build output in `dist/assets/`
3. Ensure production build: `pnpm run build`
4. Clear Vite cache: `rm -rf node_modules/.vite`

### Loading State Flashing

**Symptoms:**
- Loading spinner flashes briefly
- Poor user experience

**Solutions:**
1. Add minimum loading time:
```typescript
const MIN_LOADING_TIME = 300

router.beforeEach(() => {
  isLoading.value = true
  loadingStartTime = Date.now()
})

router.afterEach(() => {
  const elapsed = Date.now() - loadingStartTime
  const delay = Math.max(0, MIN_LOADING_TIME - elapsed)

  setTimeout(() => {
    isLoading.value = false
  }, delay)
})
```

2. Use transition delays:
```vue
<Transition name="fade" mode="out-in">
  <RouteLoading v-if="isLoading" />
</Transition>
```

## Performance Metrics

### Key Metrics to Track

1. **First Contentful Paint (FCP)**: Time until first content renders
2. **Largest Contentful Paint (LCP)**: Time until main content renders
3. **Time to Interactive (TTI)**: Time until page is fully interactive
4. **Total Bundle Size**: Sum of all JavaScript and CSS
5. **Initial Bundle Size**: Size of code loaded on first visit

### Target Metrics

- Initial bundle: < 200KB (gzipped)
- FCP: < 1.8s
- LCP: < 2.5s
- TTI: < 3.8s

### Measuring Impact

Use Lighthouse to compare before/after:
```bash
# Install Lighthouse CLI
npm install -g lighthouse

# Run audit
lighthouse http://localhost:5173 --view
```

## Related Files

- `/frontend/src/router/index.ts` - Route definitions with lazy loading
- `/frontend/src/App.vue` - Loading state management
- `/frontend/src/components/RouteLoading.vue` - Loading component
- `/frontend/vite.config.ts` - Build and chunking configuration
- `/scripts/analyze-bundle.sh` - Bundle analysis script

## Additional Resources

- [Vue Router Lazy Loading](https://router.vuejs.org/guide/advanced/lazy-loading.html)
- [Vite Build Optimization](https://vitejs.dev/guide/build.html)
- [Web.dev Code Splitting](https://web.dev/reduce-javascript-payloads-with-code-splitting/)
- [Rollup Manual Chunks](https://rollupjs.org/configuration-options/#output-manualchunks)

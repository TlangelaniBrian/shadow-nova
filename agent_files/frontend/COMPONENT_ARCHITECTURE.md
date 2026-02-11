# Component Architecture

This document outlines the component structure, best practices, and guidelines for building and organizing Vue components in Shadow Nova.

## Component Hierarchy

```
frontend/src/components/
├── auth/                    # Authentication-related components
│   ├── LoginForm.vue       # Email/password login form
│   ├── AuthProviders.vue   # Social auth providers (Google, GitHub)
│   └── LoginFeatures.vue   # Login page feature highlights
│
├── profile/                 # User profile components
│   ├── ProfileCard.vue     # User info and stats display
│   ├── GitHubIntegration.vue # GitHub connection management
│   └── AccountSettings.vue  # User account settings toggles
│
├── common/                  # Reusable shared components
│   ├── LoadingSpinner.vue  # Configurable loading spinner
│   └── StatCard.vue        # Statistics display card
│
├── layout/                  # Layout components
│   └── ...
│
└── ui/                      # Base UI components
    └── ...
```

## Component Organization Principles

### 1. Single Responsibility
Each component should have one clear purpose:
- **Good**: `ProfileCard` displays user info
- **Bad**: `ProfileCard` displays user info AND handles GitHub integration AND manages settings

### 2. Component Size Guidelines
- **Small components**: 20-50 lines (utility components like `LoadingSpinner`)
- **Medium components**: 50-150 lines (feature components like `ProfileCard`)
- **Large components**: 150-300 lines (container/view components)
- **Too large**: 300+ lines (split into smaller components)

### 3. Directory Structure
Organize components by feature/domain:
```
components/
  ├── feature-name/         # Group related components
  │   ├── ComponentA.vue
  │   └── ComponentB.vue
  └── common/               # Shared across features
      └── SharedComponent.vue
```

## When to Extract Components

Extract a component when:

1. **Repeated UI patterns** - Used in 2+ places
2. **Clear responsibility** - Has a distinct, single purpose
3. **Too much complexity** - Parent component exceeds 200 lines
4. **Independent logic** - Can function independently
5. **Testability** - Needs isolated testing

### Example: ProfileView Decomposition

**Before** (238 lines):
```vue
<template>
  <div>
    <!-- Profile card inline -->
    <div class="profile-card">...</div>
    <!-- GitHub integration inline -->
    <div class="github">...</div>
    <!-- Settings inline -->
    <div class="settings">...</div>
  </div>
</template>
```

**After** (80 lines):
```vue
<template>
  <div>
    <ProfileCard :user="user" />
    <GitHubIntegration @connect="linkGitHub" />
    <AccountSettings v-model:dark-mode="darkMode" />
  </div>
</template>
```

## Props vs Events Guidelines

### Props (Parent → Child)
Use props to pass data down:
```typescript
interface Props {
  user?: User;           // Optional object
  isLoading?: boolean;   // Optional boolean with default
  stats: GitHubStats;    // Required object
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
});
```

### Events (Child → Parent)
Use events to communicate up:
```typescript
interface Emits {
  (e: 'connect'): void;                           // Simple event
  (e: 'update:value', value: string): void;       // v-model pattern
  (e: 'submit', data: FormData): void;            // Event with payload
}

const emit = defineEmits<Emits>();
```

### Two-way Binding (v-model)
For form controls and toggles:
```vue
<!-- Parent -->
<AccountSettings v-model:dark-mode="darkMode" />

<!-- Child -->
<script setup lang="ts">
interface Props {
  darkMode?: boolean;
}
interface Emits {
  (e: 'update:darkMode', value: boolean): void;
}
</script>
```

## TypeScript Best Practices

### 1. Define Clear Interfaces
```typescript
interface User {
  name?: string;
  email?: string;
  created_at?: string;
}

interface GitHubStats {
  repos: number;
  contributions: number;
  followers: number;
}
```

### 2. Type Props Strictly
```typescript
// Good
interface Props {
  user?: User;
  count: number;
}

// Bad
interface Props {
  user?: any;        // Avoid 'any'
  count?: number;    // Should count be optional?
}
```

### 3. Type Events with Payloads
```typescript
interface Emits {
  (e: 'submit', credentials: { email: string; password: string }): void;
}
```

## Component Communication Patterns

### Pattern 1: Props Down, Events Up
```vue
<!-- Parent -->
<GitHubIntegration
  :is-connected="githubLinked"
  @connect="linkGitHub"
  @disconnect="unlinkGitHub"
/>

<!-- Child -->
<script setup lang="ts">
defineProps<{ isConnected: boolean }>();
const emit = defineEmits<{
  (e: 'connect'): void;
  (e: 'disconnect'): void;
}>();
</script>
```

### Pattern 2: v-model for Two-Way Binding
```vue
<!-- Parent -->
<AccountSettings v-model:email-notifications="emailNotifications" />

<!-- Child emits update:emailNotifications -->
```

### Pattern 3: Composables for Shared Logic
```typescript
// composables/useAuth.ts
export function useAuth() {
  const user = ref<User | null>(null);
  return { user };
}

// In component
const { user } = useAuth();
```

## Examples of Good Component Design

### 1. StatCard - Single Responsibility
```vue
<StatCard
  label="Courses Completed"
  value="12"
  trend="+15%"
  icon="📚"
  color="purple"
/>
```
- Single purpose: display a stat
- Highly reusable
- Props control all behavior
- No internal state

### 2. LoginForm - Clear Interface
```vue
<LoginForm
  :is-loading="isLoading"
  @submit="handleLogin"
/>
```
- Manages its own form state
- Emits structured data
- Parent controls loading state
- Validation built-in

### 3. GitHubIntegration - Controlled Component
```vue
<GitHubIntegration
  :is-loading="isLoading"
  :is-connecting="isConnecting"
  :is-connected="githubLinked"
  :username="githubUsername"
  :stats="githubStats"
  @connect="linkGitHub"
  @disconnect="unlinkGitHub"
/>
```
- Parent controls all state
- Component is purely presentational
- Clear event handlers
- TypeScript ensures type safety

## Anti-Patterns to Avoid

### 1. God Components
```vue
<!-- BAD: 500-line component doing everything -->
<template>
  <div>
    <!-- profile, settings, github, notifications, etc. -->
  </div>
</template>
```

### 2. Prop Drilling
```vue
<!-- BAD: Passing props through many levels -->
<Parent :user="user">
  <Child1 :user="user">
    <Child2 :user="user">
      <Child3 :user="user" />
```
**Solution**: Use composables or provide/inject

### 3. Tight Coupling
```vue
<!-- BAD: Child knows too much about parent -->
<script>
import { useParentStore } from '@/stores/parent'
const parentStore = useParentStore() // Tight coupling
</script>
```
**Solution**: Pass data via props, communicate via events

### 4. Missing TypeScript
```vue
<!-- BAD: No types -->
<script>
defineProps(['user', 'data', 'stuff'])
</script>
```
**Solution**: Always use TypeScript interfaces

## Refactoring Checklist

When refactoring a large component:

- [ ] Identify distinct UI sections
- [ ] Extract repeated patterns
- [ ] Define clear prop interfaces
- [ ] Implement proper event handling
- [ ] Add TypeScript types
- [ ] Test all interactions still work
- [ ] Update parent component imports
- [ ] Document any breaking changes

## Resources

- [Vue 3 Component Composition](https://vuejs.org/guide/components/registration.html)
- [TypeScript with Vue](https://vuejs.org/guide/typescript/overview.html)
- [Composables Pattern](https://vuejs.org/guide/reusability/composables.html)

---

**Last Updated**: February 2026

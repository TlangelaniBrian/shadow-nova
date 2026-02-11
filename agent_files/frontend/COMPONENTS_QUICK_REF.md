# Components Quick Reference

Quick lookup for all custom components in Shadow Nova frontend.

## Auth Components (`@/components/auth`)

### AuthProviders
```vue
<AuthProviders />
```
Social authentication buttons (Google, GitHub, etc.)

### LoginForm
```vue
<LoginForm
  :is-loading="false"
  @submit="handleLogin"
/>
```
Email/password login form. Emits credentials on submit.

### LoginFeatures
```vue
<LoginFeatures />
```
Static feature highlights for login page.

---

## Profile Components (`@/components/profile`)

### ProfileCard
```vue
<ProfileCard
  :user="user"
  :courses-completed="12"
/>
```
Displays user avatar, name, email, and stats.

### GitHubIntegration
```vue
<GitHubIntegration
  :is-loading="false"
  :is-connecting="false"
  :is-connected="true"
  :username="'johndoe'"
  :stats="{ repos: 24, contributions: 1247, followers: 89 }"
  @connect="handleConnect"
  @disconnect="handleDisconnect"
/>
```
GitHub connection management with stats display.

### AccountSettings
```vue
<AccountSettings
  v-model:email-notifications="emailNotifications"
  v-model:dark-mode="darkMode"
/>
```
User settings toggles (email notifications, dark mode).

---

## Common Components (`@/components/common`)

### LoadingSpinner
```vue
<LoadingSpinner size="md" color="purple" />
```
Configurable loading spinner.
- **Sizes**: `sm`, `md`, `lg`
- **Colors**: `purple`, `blue`, `green`, `red`, `white`

### StatCard
```vue
<StatCard
  label="Courses Completed"
  :value="12"
  trend="+15%"
  icon="📚"
  color="purple"
/>
```
Statistics display card with icon and optional trend.
- **Colors**: `purple`, `blue`, `green`, `orange`, `pink`, `red`

---

## Import Patterns

### Direct Import
```typescript
import ProfileCard from '@/components/profile/ProfileCard.vue'
```

### Barrel Import
```typescript
import { ProfileCard, GitHubIntegration } from '@/components/profile'
```

---

## Component Sizing Guidelines

- **Small**: 20-50 lines (utilities like LoadingSpinner)
- **Medium**: 50-150 lines (features like ProfileCard)
- **Large**: 150-300 lines (views and containers)
- **Too Large**: 300+ lines (consider decomposing)

---

## Common Patterns

### Props + Events
```vue
<!-- Parent -->
<MyComponent
  :data="myData"
  @action="handleAction"
/>

<!-- Child -->
<script setup lang="ts">
interface Props {
  data: MyData;
}
interface Emits {
  (e: 'action', payload: string): void;
}
const props = defineProps<Props>();
const emit = defineEmits<Emits>();
</script>
```

### v-model Pattern
```vue
<!-- Parent -->
<MyComponent v-model:value="myValue" />

<!-- Child -->
<script setup lang="ts">
interface Props {
  value: string;
}
interface Emits {
  (e: 'update:value', value: string): void;
}
</script>
```

---

## TypeScript Interfaces

### User
```typescript
interface User {
  name?: string;
  email?: string;
  github_username?: string;
  created_at?: string;
}
```

### GitHubStats
```typescript
interface GitHubStats {
  repos: number;
  contributions: number;
  followers: number;
}
```

---

## See Also

- [COMPONENT_ARCHITECTURE.md](./COMPONENT_ARCHITECTURE.md) - Detailed architecture guide
- [COMPONENT_MIGRATION_GUIDE.md](./COMPONENT_MIGRATION_GUIDE.md) - Migration documentation

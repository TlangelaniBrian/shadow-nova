# Shadow Nova Components

This directory contains all reusable Vue components for the Shadow Nova application.

## Directory Structure

```
components/
├── auth/                    # Authentication components
│   ├── index.ts            # Barrel exports
│   ├── AuthProviders.vue   # Social auth providers
│   ├── LoginForm.vue       # Email/password login
│   └── LoginFeatures.vue   # Login page features
│
├── profile/                 # User profile components
│   ├── index.ts            # Barrel exports
│   ├── ProfileCard.vue     # User info display
│   ├── GitHubIntegration.vue # GitHub connection
│   └── AccountSettings.vue # Settings toggles
│
├── common/                  # Shared utility components
│   ├── index.ts            # Barrel exports
│   ├── LoadingSpinner.vue  # Loading indicator
│   └── StatCard.vue        # Statistics card
│
├── layout/                  # Layout components
│   └── ... (existing)
│
└── ui/                      # Base UI components
    └── ... (existing)
```

## Component Categories

### Auth Components
Components related to authentication and user onboarding.

- **AuthProviders** - Social authentication buttons (Google, GitHub)
- **LoginForm** - Email and password login form
- **LoginFeatures** - Feature highlights for login/signup pages

### Profile Components
Components for user profile management.

- **ProfileCard** - Displays user avatar, name, email, and basic stats
- **GitHubIntegration** - Manages GitHub account connection and displays stats
- **AccountSettings** - User preference toggles (notifications, theme, etc.)

### Common Components
Reusable utility components used across the application.

- **LoadingSpinner** - Configurable loading spinner (size, color)
- **StatCard** - Statistics display card with icon and trend

### Layout Components
Components that define page structure and navigation.

- See `layout/README.md` for details

### UI Components
Base-level UI components (buttons, inputs, cards, etc.).

- See `ui/README.md` for details

## Usage Guidelines

### Importing Components

#### Direct Import
```typescript
import ProfileCard from '@/components/profile/ProfileCard.vue'
```

#### Barrel Import
```typescript
import { ProfileCard, GitHubIntegration } from '@/components/profile'
```

### Component Sizing

- **Utility Components**: 20-50 lines
  - LoadingSpinner, icons, small widgets
- **Feature Components**: 50-150 lines
  - ProfileCard, LoginForm, StatCard
- **Container Components**: 150-300 lines
  - Complex integrations like GitHubIntegration

### When to Create a Component

Create a new component when:
1. The UI pattern is used in multiple places
2. The component has a clear, single responsibility
3. The parent component is becoming too complex (200+ lines)
4. The logic can be tested independently
5. The component is reusable in different contexts

### Component Design Principles

1. **Single Responsibility** - Each component does one thing well
2. **Props Down, Events Up** - Unidirectional data flow
3. **TypeScript First** - Strict typing for props and events
4. **Composable Logic** - Use composables for shared logic
5. **Accessible** - Follow WCAG guidelines

## Quick Reference

### Props Pattern
```vue
<script setup lang="ts">
interface Props {
  required: string;
  optional?: number;
}

const props = withDefaults(defineProps<Props>(), {
  optional: 0,
});
</script>
```

### Events Pattern
```vue
<script setup lang="ts">
interface Emits {
  (e: 'action', payload: string): void;
}

const emit = defineEmits<Emits>();
</script>
```

### v-model Pattern
```vue
<script setup lang="ts">
interface Props {
  modelValue: string;
}

interface Emits {
  (e: 'update:modelValue', value: string): void;
}
</script>
```

## Documentation

For more detailed information:

- **[COMPONENT_ARCHITECTURE.md](../../COMPONENT_ARCHITECTURE.md)** - Architecture and best practices
- **[COMPONENTS_QUICK_REF.md](../../COMPONENTS_QUICK_REF.md)** - Quick reference guide
- **[COMPONENT_HIERARCHY.md](../../COMPONENT_HIERARCHY.md)** - Visual component tree

## Testing

### Unit Testing
```bash
npm run test:unit
```

### Component Testing
```bash
npm run test:component
```

### E2E Testing
```bash
npm run test:e2e
```

## Contributing

When adding a new component:

1. Choose the appropriate directory (auth, profile, common, etc.)
2. Create the component file with TypeScript
3. Define clear prop and event interfaces
4. Add the component to the directory's `index.ts`
5. Update this README if adding a new category
6. Add unit tests
7. Update relevant documentation

## Component Standards

### File Naming
- PascalCase: `ProfileCard.vue`, `LoginForm.vue`
- Match the component name exactly

### Script Setup
```vue
<script setup lang="ts">
// Always use TypeScript
// Always use script setup
</script>
```

### Template Style
```vue
<template>
  <!-- Use semantic HTML -->
  <!-- Keep templates focused and readable -->
  <!-- Extract complex logic to computed properties -->
</template>
```

### Styling
```vue
<style scoped>
/* Use scoped styles */
/* Prefer Tailwind classes */
/* Only use <style> for complex/dynamic styles */
</style>
```

## Component Checklist

Before committing a new component:

- [ ] TypeScript interfaces defined
- [ ] Props have defaults where appropriate
- [ ] Events are typed
- [ ] Component is under 150 lines (or has good reason to be larger)
- [ ] Added to directory's index.ts
- [ ] Documentation strings added
- [ ] Unit tests written
- [ ] Tested in isolation
- [ ] Tested in parent component
- [ ] Responsive design verified
- [ ] Accessibility checked

---

**Last Updated**: February 12, 2026
**Components**: 20+ components across 5 categories

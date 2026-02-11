# Component Migration Guide

This guide documents the component decomposition performed in February 2026 and provides a template for future refactoring work.

## Summary of Changes

### ProfileView (238 → 108 lines, -54% reduction)

**Extracted Components:**
1. **ProfileCard.vue** (59 lines)
   - User avatar and display name
   - User email
   - Member stats (member since, courses completed)

2. **GitHubIntegration.vue** (145 lines)
   - GitHub connection status
   - Connect/disconnect functionality
   - GitHub stats display (repos, contributions, followers)

3. **AccountSettings.vue** (67 lines)
   - Email notifications toggle
   - Dark mode toggle
   - Future settings options

**Benefits:**
- Improved maintainability
- Better testability (each component can be tested in isolation)
- Reusable components (GitHubIntegration could be used elsewhere)
- Clearer separation of concerns

### LoginView (230 → 146 lines, -36% reduction)

**Extracted Components:**
1. **LoginForm.vue** (55 lines)
   - Email and password input fields
   - Form validation
   - Submit handling with typed events

2. **AuthProviders.vue** (8 lines)
   - Google OAuth button
   - Future: GitHub, Microsoft, etc.

3. **LoginFeatures.vue** (28 lines)
   - Feature highlights list
   - Marketing content
   - Static display

**Benefits:**
- Login form is now reusable (could be used in a modal, different page, etc.)
- Auth providers can be easily extended
- Features list can be updated without touching form logic

### Common Components Created

1. **LoadingSpinner.vue** (58 lines)
   - Configurable size (sm, md, lg)
   - Configurable color (purple, blue, green, red, white)
   - Reusable across entire app

2. **StatCard.vue** (48 lines)
   - Display statistics with icon and trend
   - Configurable colors
   - Used in Dashboard and could be used elsewhere

## File Structure

```
frontend/src/
├── components/
│   ├── auth/
│   │   ├── index.ts              # Barrel export
│   │   ├── AuthProviders.vue     # Social auth buttons
│   │   ├── LoginForm.vue         # Email/password form
│   │   └── LoginFeatures.vue     # Feature highlights
│   │
│   ├── profile/
│   │   ├── index.ts              # Barrel export
│   │   ├── ProfileCard.vue       # User info card
│   │   ├── GitHubIntegration.vue # GitHub connection
│   │   └── AccountSettings.vue   # Settings toggles
│   │
│   └── common/
│       ├── index.ts              # Barrel export
│       ├── LoadingSpinner.vue    # Loading indicator
│       └── StatCard.vue          # Stat display card
│
└── views/
    ├── ProfileView.vue           # Now 108 lines (was 238)
    └── LoginView.vue             # Now 146 lines (was 230)
```

## Migration Steps Followed

### 1. Analysis Phase
- Identified large components (200+ lines)
- Analyzed component responsibilities
- Found repeated UI patterns
- Identified clear boundaries

### 2. Planning Phase
- Decided on component boundaries
- Defined prop interfaces
- Planned event structure
- Considered TypeScript types

### 3. Extraction Phase
- Created new component directories
- Extracted template sections
- Moved related logic to new components
- Defined TypeScript interfaces

### 4. Integration Phase
- Updated parent components to use new children
- Wired up props and events
- Tested all interactions
- Verified TypeScript compilation

### 5. Documentation Phase
- Created COMPONENT_ARCHITECTURE.md
- Documented patterns and best practices
- Created this migration guide

## Import Patterns

### Before (Direct Imports)
```typescript
import ProfileCard from '@/components/profile/ProfileCard.vue'
import GitHubIntegration from '@/components/profile/GitHubIntegration.vue'
import AccountSettings from '@/components/profile/AccountSettings.vue'
```

### After (Barrel Imports - Optional)
```typescript
import { ProfileCard, GitHubIntegration, AccountSettings } from '@/components/profile'
```

Both patterns work, choose based on preference.

## Component API Examples

### ProfileCard
```vue
<ProfileCard
  :user="user"
  :courses-completed="12"
/>
```

**Props:**
- `user?: User` - User object with name, email
- `coursesCompleted?: number` - Number of courses (default: 12)

**No Events**

### GitHubIntegration
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

**Props:**
- `isLoading?: boolean` - Loading state
- `isConnecting?: boolean` - Connecting state
- `isConnected?: boolean` - Connected state
- `username?: string` - GitHub username
- `stats?: GitHubStats` - Stats object

**Events:**
- `@connect` - User clicked connect
- `@disconnect` - User clicked disconnect

### AccountSettings
```vue
<AccountSettings
  :email-notifications="emailNotifications"
  :dark-mode="darkMode"
  @update:email-notifications="emailNotifications = $event"
  @update:dark-mode="darkMode = $event"
/>
```

**Props:**
- `emailNotifications?: boolean` - Email notification setting
- `darkMode?: boolean` - Dark mode setting

**Events:**
- `@update:emailNotifications` - Setting changed
- `@update:darkMode` - Setting changed

### LoginForm
```vue
<LoginForm
  :is-loading="isLoading"
  @submit="handleLogin"
/>
```

**Props:**
- `email?: string` - Pre-filled email
- `password?: string` - Pre-filled password
- `isLoading?: boolean` - Loading state

**Events:**
- `@submit` - Form submitted with `{ email, password }`

### LoadingSpinner
```vue
<LoadingSpinner size="md" color="purple" />
```

**Props:**
- `size?: 'sm' | 'md' | 'lg'` - Spinner size (default: 'md')
- `color?: 'purple' | 'blue' | 'green' | 'red' | 'white'` - Color (default: 'purple')

### StatCard
```vue
<StatCard
  label="Courses Completed"
  value="12"
  trend="+15%"
  icon="📚"
  color="purple"
/>
```

**Props:**
- `label: string` - Stat label
- `value: string | number` - Stat value
- `trend?: string` - Trend indicator
- `icon: string` - Emoji or icon
- `color?: 'purple' | 'blue' | 'green' | 'orange' | 'pink' | 'red'` - Color scheme

## Testing Checklist

After component extraction:

- [ ] All imports resolve correctly
- [ ] TypeScript compiles without errors
- [ ] All interactions work (buttons, forms, toggles)
- [ ] Props are passed correctly
- [ ] Events are emitted and handled
- [ ] Styles are preserved
- [ ] Responsive layout still works
- [ ] No console errors
- [ ] Visual regression test passes

## Future Refactoring Candidates

Components that may benefit from decomposition:

1. **DashboardView.vue** (147 lines)
   - Could extract stat cards into reusable pattern
   - Already relatively clean

2. Any view over 200 lines should be considered for decomposition

## Lessons Learned

1. **Start with clear boundaries** - Define what each component should do before extracting
2. **Type everything** - TypeScript interfaces prevent integration errors
3. **Test incrementally** - Test after each extraction, don't extract everything at once
4. **Document as you go** - Write down the component API while creating it
5. **Consider reusability** - Make components generic enough to be reused

## Resources

- [COMPONENT_ARCHITECTURE.md](./COMPONENT_ARCHITECTURE.md) - Architecture guidelines
- [Vue 3 Component Guide](https://vuejs.org/guide/components/registration.html)
- [TypeScript with Vue](https://vuejs.org/guide/typescript/overview.html)

---

**Migration Date**: February 12, 2026
**Migrated By**: Claude Sonnet 4.5
**Files Changed**: 13 files created, 2 files updated

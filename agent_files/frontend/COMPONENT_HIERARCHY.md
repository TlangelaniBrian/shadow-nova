# Component Hierarchy Diagram

Visual representation of component relationships in Shadow Nova frontend.

## ProfileView Component Tree

```
ProfileView.vue (108 lines)
│
├── Header Section (inline)
│   └── "Profile Settings" title
│
├── ProfileCard.vue (59 lines)
│   ├── Props: user, coursesCompleted
│   └── Displays:
│       ├── User avatar (gradient circle)
│       ├── User name and email
│       └── Stats badges (member since, courses)
│
├── GitHubIntegration.vue (145 lines)
│   ├── Props: isLoading, isConnecting, isConnected, username, stats
│   ├── Events: @connect, @disconnect
│   └── Displays:
│       ├── GitHub icon and title
│       ├── Loading spinner (when loading)
│       ├── Connected state:
│       │   ├── Username display
│       │   ├── Disconnect button
│       │   └── Stats grid (repos, contributions, followers)
│       └── Not connected state:
│           ├── Benefits list
│           └── Connect button
│
└── AccountSettings.vue (67 lines)
    ├── Props: emailNotifications, darkMode
    ├── Events: @update:emailNotifications, @update:darkMode
    └── Displays:
        ├── Email notifications toggle
        └── Dark mode toggle
```

## LoginView Component Tree

```
LoginView.vue (146 lines)
│
├── Background Pattern (inline)
│   └── SVG pattern overlay
│
├── Logo & Title Section (inline)
│   ├── Logo icon
│   ├── "Shadow Nova" title
│   └── Tagline
│
├── Description (inline)
│   └── Platform description text
│
├── AuthProviders.vue (8 lines)
│   └── GoogleSignIn.vue
│       └── Google OAuth button
│
├── Divider (inline)
│   └── "Or continue with" text
│
├── LoginForm.vue (55 lines)
│   ├── Props: email, password, isLoading
│   ├── Events: @submit(credentials)
│   └── Displays:
│       ├── Email input field
│       ├── Password input field
│       └── Submit button
│
├── LoginFeatures.vue (28 lines)
│   └── Displays:
│       ├── ✓ Structured learning paths
│       ├── ✓ Hands-on projects
│       └── ✓ Progress tracking
│
└── Footer (inline)
    └── Terms of Service text
```

## Common Components (Shared)

```
LoadingSpinner.vue (58 lines)
├── Props: size ('sm'|'md'|'lg'), color
└── Renders: Animated spinning circle

StatCard.vue (48 lines)
├── Props: label, value, trend, icon, color
└── Renders:
    ├── Icon badge (colored)
    ├── Trend indicator (optional)
    ├── Label text
    └── Value (large)
```

## Component Dependencies Graph

```
views/ProfileView.vue
    ├── components/profile/ProfileCard.vue
    ├── components/profile/GitHubIntegration.vue
    └── components/profile/AccountSettings.vue

views/LoginView.vue
    ├── components/auth/AuthProviders.vue
    │   └── components/GoogleSignIn.vue
    ├── components/auth/LoginForm.vue
    └── components/auth/LoginFeatures.vue

views/DashboardView.vue
    └── (could use) components/common/StatCard.vue

(any view)
    └── (could use) components/common/LoadingSpinner.vue
```

## Data Flow Diagrams

### ProfileView Data Flow

```
ProfileView State
    ├── user (from useAuth composable)
    │   └── → ProfileCard :user
    │
    ├── githubLinked, githubUsername, githubStats
    │   └── → GitHubIntegration :isConnected, :username, :stats
    │
    └── emailNotifications, darkMode
        └── ⟷ AccountSettings v-model:emailNotifications, v-model:darkMode

Events:
    GitHubIntegration @connect    → linkGitHub()    → API call
    GitHubIntegration @disconnect → unlinkGitHub()  → API call
    AccountSettings   @update:*   → update state    → (future: save to API)
```

### LoginView Data Flow

```
LoginView State
    └── isLoading
        └── → LoginForm :isLoading

Events:
    LoginForm @submit(credentials)
        └── handleLogin(credentials)
            └── fetch('/api/v1/login')
                └── localStorage.setItem('user')
                    └── router.push('/dashboard')
```

## Component Size Comparison

### Before Decomposition
```
ProfileView.vue  ██████████████████████████████████████████████ 238 lines
LoginView.vue    ████████████████████████████████████████████   230 lines
```

### After Decomposition
```
ProfileView.vue        ███████████████████████ 108 lines
LoginView.vue          ███████████████████████████████ 146 lines

ProfileCard.vue        ████████████ 59 lines
GitHubIntegration.vue  ███████████████████████████████ 145 lines
AccountSettings.vue    ███████████████ 67 lines

LoginForm.vue          ████████████ 55 lines
AuthProviders.vue      ██ 8 lines
LoginFeatures.vue      ██████ 28 lines

LoadingSpinner.vue     ████████████ 58 lines
StatCard.vue           ██████████ 48 lines
```

## Component Reusability Matrix

| Component | Used In | Can Be Reused In |
|-----------|---------|------------------|
| ProfileCard | ProfileView | User dropdown, Team page |
| GitHubIntegration | ProfileView | Onboarding, Settings |
| AccountSettings | ProfileView | Settings drawer |
| LoginForm | LoginView | Modal, Registration |
| AuthProviders | LoginView | Registration, Settings |
| LoginFeatures | LoginView | Landing page, About |
| LoadingSpinner | GitHubIntegration | Any async operation |
| StatCard | (future) DashboardView | Analytics, Reports |

## Component Communication Patterns

### Pattern 1: Props Down (Parent → Child)
```
ProfileView
    ↓ :user
ProfileCard
```

### Pattern 2: Events Up (Child → Parent)
```
GitHubIntegration
    ↑ @connect
ProfileView
```

### Pattern 3: v-model (Two-Way Binding)
```
AccountSettings
    ⟷ v-model:darkMode
ProfileView
```

### Pattern 4: Composable Sharing
```
ProfileView ← useAuth() → DashboardView
                ↓
              user ref
```

## File Size Distribution

```
Component Type       | Count | Avg Lines | Total Lines
---------------------|-------|-----------|------------
View Components      |   2   |    127    |    254
Profile Components   |   3   |     90    |    271
Auth Components      |   3   |     30    |     91
Common Components    |   2   |     53    |    106
---------------------|-------|-----------|------------
TOTAL                |  10   |     72    |    722
```

## Import Statement Examples

### ProfileView Imports
```typescript
import { ref, watchEffect } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useToast } from '@/composables/useToast';
import client from '@/api/client';
import ProfileCard from '@/components/profile/ProfileCard.vue';
import GitHubIntegration from '@/components/profile/GitHubIntegration.vue';
import AccountSettings from '@/components/profile/AccountSettings.vue';
```

### LoginView Imports
```typescript
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthProviders from '@/components/auth/AuthProviders.vue'
import LoginForm from '@/components/auth/LoginForm.vue'
import LoginFeatures from '@/components/auth/LoginFeatures.vue'
```

---

**Visual Key:**
- `█` = 10 lines of code
- `→` = Data flow (props)
- `←` = Data flow (reverse)
- `↑` = Event emission
- `↓` = Prop passing
- `⟷` = Two-way binding (v-model)

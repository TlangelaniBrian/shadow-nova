# Dark Mode Implementation

This document describes how dark mode is implemented in Shadow Nova and provides guidelines for maintaining and extending it.

## Architecture

Shadow Nova uses a Pinia store-based dark mode system with class-based Tailwind CSS dark mode support.

### Key Components

1. **Theme Store** (`src/stores/theme.ts`)
   - Manages dark mode state
   - Persists preference to localStorage
   - Applies theme by adding/removing 'dark' class on document root
   - Detects system preference on first load

2. **Tailwind Configuration** (`tailwind.config.js`)
   - Uses `darkMode: 'class'` strategy
   - Dark mode classes are prefixed with `dark:`

3. **UI Controls**
   - Header toggle button (Moon/Sun icon)
   - Settings page checkbox
   - Both controls sync automatically through the store

## How It Works

### Initialization Flow

1. On app load, the theme store checks localStorage for saved preference
2. If no preference exists, it checks system preference using `window.matchMedia('(prefers-color-scheme: dark)')`
3. The theme is immediately applied to avoid flash of unstyled content
4. A Vue watcher persists any changes to localStorage

### Theme Application

The `applyTheme()` function toggles the `dark` class on `document.documentElement`:

```typescript
function applyTheme() {
    if (isDarkMode.value) {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }
}
```

This allows all Tailwind `dark:` utilities to activate/deactivate globally.

## Using Dark Mode in Components

### Basic Pattern

Always provide both light and dark variants for visual properties:

```vue
<div class="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">
  <!-- content -->
</div>
```

### Common Patterns

#### Background Colors
```vue
<!-- Cards -->
class="bg-white dark:bg-gray-800"

<!-- Main background -->
class="bg-gray-50 dark:bg-gray-900"

<!-- Hover states -->
class="hover:bg-gray-100 dark:hover:bg-gray-800"
```

#### Text Colors
```vue
<!-- Primary text -->
class="text-gray-900 dark:text-gray-100"

<!-- Secondary text -->
class="text-gray-600 dark:text-gray-400"

<!-- Muted text -->
class="text-gray-400 dark:text-gray-500"
```

#### Borders
```vue
<!-- Standard borders -->
class="border-gray-100 dark:border-gray-700"

<!-- Subtle borders -->
class="border-gray-50 dark:border-gray-800"
```

#### Form Inputs
```vue
<input
  class="bg-white dark:bg-gray-900
         border-gray-200 dark:border-gray-700
         text-gray-900 dark:text-gray-100
         focus:ring-purple-500"
/>
```

#### Shadows
```vue
<!-- Light shadow in light mode, darker in dark mode -->
class="shadow-sm dark:shadow-[0_2px_10px_rgba(0,0,0,0.3)]"
```

### Component Checklist

When adding dark mode to a component:

- [ ] Background colors
- [ ] Text colors (headings, body, muted)
- [ ] Border colors
- [ ] Hover states
- [ ] Active/focus states
- [ ] Icons (if using colored icons)
- [ ] Shadows
- [ ] Form inputs and controls

## Best Practices

### 1. Consistency

Use consistent color mappings across the app:
- Light bg-white → Dark bg-gray-800
- Light bg-gray-50 → Dark bg-gray-900
- Light text-gray-900 → Dark text-gray-100
- Light text-gray-600 → Dark text-gray-400

### 2. Contrast

Ensure sufficient contrast in both modes:
- Test with browser DevTools accessibility checker
- Verify text is readable on all backgrounds
- Check focus indicators are visible

### 3. Color Usage

- Don't just invert colors - dark mode is not a simple inverse
- Maintain brand colors (purple gradient) consistently
- Use lighter shades for dark mode backgrounds
- Use darker shades for dark mode text

### 4. Testing

Test dark mode on:
- All major pages/views
- Interactive states (hover, focus, active)
- Form validations and errors
- Loading states
- Empty states

### 5. System Preference

The app respects system preference on first load. Users can override this preference, and their choice is saved.

## Accessing the Theme Store

In any component:

```vue
<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

// Check current theme
if (themeStore.isDarkMode) {
  // dark mode specific logic
}

// Toggle theme
themeStore.toggleDarkMode()
</script>
```

## Troubleshooting

### Dark mode not persisting
- Check localStorage in DevTools
- Verify theme store is initialized before components render

### Flash of wrong theme
- Theme should be applied in store initialization before first render
- Check that `applyTheme()` is called synchronously on init

### Styles not updating
- Verify Tailwind config has `darkMode: 'class'`
- Check that dark: variants are included in your purge/content config
- Ensure you're using dark: prefix, not @media queries

### Missing dark mode classes
- Add dark: variants to new components
- Check that the 'dark' class is on document.documentElement
- Inspect element in DevTools to see applied classes

## Future Enhancements

Potential improvements:
- Add transition animations between theme switches
- Support for custom theme colors (not just light/dark)
- More granular theme controls (accent colors, etc.)
- Theme-aware image/SVG loading
- Auto-toggle based on time of day

## Resources

- [Tailwind CSS Dark Mode](https://tailwindcss.com/docs/dark-mode)
- [Pinia Store Documentation](https://pinia.vuejs.org/)
- [Web.dev: prefers-color-scheme](https://web.dev/prefers-color-scheme/)

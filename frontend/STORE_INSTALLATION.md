# Store Installation Guide

## Quick Start

### 1. Install Dependencies

```bash
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/frontend
npm install pinia-plugin-persistedstate
```

### 2. Enable Persistence Plugin

Edit `src/main.ts` and uncomment these lines:

```typescript
// Line 8 - Uncomment:
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

// Lines 31-33 - Replace:
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)
```

### 3. Enable Store Persistence

Uncomment persistence configuration in each store:

**`src/stores/projects.ts`** (lines 65-69):
```typescript
}, {
    persist: {
        paths: ['projects']
    }
})
```

**`src/stores/learningPaths.ts`** (lines 58-62):
```typescript
}, {
    persist: {
        paths: ['paths', 'currentPath']
    }
})
```

**`src/stores/progress.ts`** (lines 69-73):
```typescript
}, {
    persist: {
        paths: ['progressMap', 'stats']
    }
})
```

### 4. Test the Implementation

Start the dev server:
```bash
npm run dev
```

Visit the projects page and verify:
- Projects load correctly
- Loading states display
- Error handling works
- State persists on page reload

## Verification Checklist

- [ ] `pinia-plugin-persistedstate` installed
- [ ] Plugin enabled in `main.ts`
- [ ] Persistence enabled in all three stores
- [ ] Projects page loads correctly
- [ ] No TypeScript errors
- [ ] Dev server runs without errors
- [ ] State persists after page reload

## Troubleshooting

### TypeScript Errors
If you see TypeScript errors about missing types:
```bash
npm install --save-dev @types/node
```

### API Connection Issues
Verify your `.env` file has the correct API URL:
```
VITE_API_URL=http://localhost:8080
```

### Store Not Persisting
Check browser console for:
- localStorage permissions
- Any Pinia plugin errors
- Verify persistence config is uncommented

## Next Steps

After installation is complete:

1. Update `LearningPathsView.vue` to use the learning paths store
2. Update `PathDetailView.vue` to use both learning paths and progress stores
3. Consider removing the now-unused `useProjects` composable
4. Test all views thoroughly
5. Add unit tests for stores

See `STORE_SETUP_COMPLETE.md` for detailed next steps.

import { createRouter, createWebHistory, type RouteLocationNormalized, type NavigationGuardNext } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import(/* webpackChunkName: "login" */ '../views/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/auth/github/callback',
      name: 'github-callback',
      component: () => import(/* webpackChunkName: "github-callback" */ '../views/GitHubCallback.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import(/* webpackChunkName: "profile" */ '../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import(/* webpackChunkName: "dashboard" */ '../views/DashboardView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/home',
      name: 'home',
      component: () => import(/* webpackChunkName: "home" */ '../views/HomeView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/paths',
      name: 'learning-paths',
      component: () => import(/* webpackChunkName: "learning-paths" */ '../views/LearningPathsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/paths/:id',
      name: 'path-detail',
      component: () => import(/* webpackChunkName: "path-detail" */ '../views/PathDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import(/* webpackChunkName: "projects" */ '../views/ProjectsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/resources',
      name: 'resources',
      component: () => import(/* webpackChunkName: "resources" */ '../views/ResourcesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/community',
      name: 'community',
      component: () => import(/* webpackChunkName: "community" */ '../views/CommunityView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/guides',
      name: 'guides',
      component: () => import(/* webpackChunkName: "guides" */ '../views/GuidesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login',
    },
  ],
})

// Navigation guard for authentication
router.beforeEach((to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) => {
  // Check if user object exists (populated by API call)
  const userStr = localStorage.getItem('user')
  const hasUser = !!userStr
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // Token validation is now done by backend on each request
  // Frontend only checks if user object exists for UI state
  if (requiresAuth && !hasUser) {
    next('/login')
  } else if (to.path === '/login' && hasUser) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router

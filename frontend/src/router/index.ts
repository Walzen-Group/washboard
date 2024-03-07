/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router/auto'
import { setupLayouts } from 'virtual:generated-layouts'
import { isAuthorized } from '@/api/lib'

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  extendRoutes: setupLayouts,
})

router.beforeEach(async (to, from) => {
  const authorized = await isAuthorized();
  if (!authorized && to.name !== '/login') {
    return {
      name: '/login',
      query: { redirect: to.path }
    };
  }
})

export default router

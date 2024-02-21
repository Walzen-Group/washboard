/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router/auto'
import { setupLayouts } from 'virtual:generated-layouts'

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  extendRoutes: setupLayouts,
})

router.onError((error: any, to: any) => {
  if (error.message.includes('Failed to fetch dynamically imported module') || error.message.includes("Importing a module script failed")) {
    window.location = to.fullPath
  }
})

export default router

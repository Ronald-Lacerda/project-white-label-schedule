export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuth()

  if (auth.isAuthenticated.value) return

  // Tenta restaurar sessão via refresh token
  const restored = await auth.refresh()
  if (restored) return

  return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
})

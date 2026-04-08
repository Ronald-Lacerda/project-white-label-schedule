interface User {
  id: string
  name: string
  email: string
  role: string
  establishment_id: string
}

interface AuthState {
  user: User | null
  accessToken: string | null
}

interface AuthResponse {
  data: {
    access_token: string
    refresh_token: string
    user: User
  }
}

interface RegisterPayload {
  owner_name: string
  establishment_name: string
  email: string
  password: string
  slug: string
  contact_phone: string | null
}

const state = reactive<AuthState>({
  user: null,
  accessToken: null,
})

export function useAuth() {
  const config = useRuntimeConfig()
  const router = useRouter()

  const isAuthenticated = computed(() => !!state.accessToken)
  const user = computed(() => state.user)

  function persistSession(payload: AuthResponse['data']) {
    state.accessToken = payload.access_token
    state.user = payload.user
    if (import.meta.client) {
      localStorage.setItem('refresh_token', payload.refresh_token)
    }
  }

  async function login(email: string, password: string): Promise<void> {
    const res = await $fetch<AuthResponse>(
      `${config.public.apiBaseUrl}/api/v1/auth/login`,
      { method: 'POST', body: { email, password } },
    )
    persistSession(res.data)
  }

  async function register(payload: RegisterPayload): Promise<void> {
    const res = await $fetch<AuthResponse>(
      `${config.public.apiBaseUrl}/api/v1/auth/register`,
      { method: 'POST', body: payload },
    )
    persistSession(res.data)
  }

  async function logout(): Promise<void> {
    const refreshToken = import.meta.client ? localStorage.getItem('refresh_token') : null
    try {
      if (state.accessToken) {
        await $fetch(`${config.public.apiBaseUrl}/api/v1/auth/logout`, {
          method: 'POST',
          headers: authHeaders(),
          body: { refresh_token: refreshToken ?? '' },
        })
      }
    } finally {
      state.accessToken = null
      state.user = null
      if (import.meta.client) localStorage.removeItem('refresh_token')
      await router.push('/login')
    }
  }

  async function refresh(): Promise<boolean> {
    if (!import.meta.client) return false
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) return false

    try {
      const res = await $fetch<AuthResponse>(
        `${config.public.apiBaseUrl}/api/v1/auth/refresh`,
        { method: 'POST', body: { refresh_token: refreshToken } },
      )
      persistSession(res.data)
      return true
    } catch {
      state.accessToken = null
      state.user = null
      localStorage.removeItem('refresh_token')
      return false
    }
  }

  function authHeaders(): Record<string, string> {
    if (!state.accessToken) return {}
    return { Authorization: `Bearer ${state.accessToken}` }
  }

  return { isAuthenticated, user, login, register, logout, refresh, authHeaders }
}

interface ProfessionalRef {
  id: string
  name: string
  google_calendar_id: string | null
}

interface GoogleStatus {
  connected: boolean
  professionals: ProfessionalRef[]
}

export function useGoogleCalendar() {
  const api = useApi()

  const status = ref<GoogleStatus | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  function normalizeStatus(data: GoogleStatus | null | undefined): GoogleStatus {
    return {
      connected: Boolean(data?.connected),
      professionals: Array.isArray(data?.professionals) ? data.professionals : [],
    }
  }

  async function fetchStatus() {
    loading.value = true
    error.value = null
    try {
      status.value = normalizeStatus(await api.get<GoogleStatus | null>('/api/v1/google/status'))
    } catch (e: any) {
      status.value = normalizeStatus(null)
      error.value = e?.data?.error?.message ?? 'Erro ao carregar status.'
    } finally {
      loading.value = false
    }
  }

  async function getAuthUrl(): Promise<string> {
    const data = await api.get<{ url: string }>('/api/v1/google/auth-url')
    return data.url
  }

  async function disconnect() {
    await api.del('/api/v1/google/disconnect')
    status.value = normalizeStatus({
      connected: false,
      professionals: status.value?.professionals ?? [],
    })
  }

  return {
    status,
    loading,
    error,
    fetchStatus,
    getAuthUrl,
    disconnect,
  }
}

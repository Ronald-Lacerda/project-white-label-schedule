export interface Establishment {
  id: string
  name: string
  slug: string
  timezone: string
  contact_email: string | null
  contact_phone: string | null
  min_advance_cancel_hours: number
  active: boolean
  google_calendar_connected: boolean
}

export interface WhitelabelConfig {
  establishment_id: string
  logo_url: string | null
  primary_color: string
  secondary_color: string | null
  custom_css: string | null
}

export function useEstablishment() {
  const api = useApi()

  const establishment = ref<Establishment | null>(null)
  const whitelabel = ref<WhitelabelConfig | null>(null)
  const loading = ref(false)
  const error = ref('')

  async function fetchEstablishment() {
    loading.value = true
    error.value = ''
    try {
      establishment.value = await api.get<Establishment>('/api/v1/establishment')
    } catch {
      error.value = 'Erro ao carregar os dados do estabelecimento.'
    } finally {
      loading.value = false
    }
  }

  async function updateEstablishment(data: Partial<Establishment>) {
    const updated = await api.put<Establishment>('/api/v1/establishment', data)
    establishment.value = updated
    return updated
  }

  async function fetchWhitelabel() {
    whitelabel.value = await api.get<WhitelabelConfig>('/api/v1/whitelabel')
  }

  async function updateWhitelabel(data: Partial<WhitelabelConfig>) {
    const updated = await api.put<WhitelabelConfig>('/api/v1/whitelabel', data)
    whitelabel.value = updated
    return updated
  }

  async function uploadLogo(file: File) {
    const config = useRuntimeConfig()
    const auth = useAuth()
    const formData = new FormData()
    formData.append('logo', file)

    const res = await $fetch<{ data: { logo_url: string } }>(
      `${config.public.apiBaseUrl}/api/v1/whitelabel/logo`,
      {
        method: 'POST',
        body: formData,
        headers: auth.authHeaders(),
      },
    )

    whitelabel.value = whitelabel.value
      ? { ...whitelabel.value, logo_url: res.data.logo_url }
      : null

    return res.data.logo_url
  }

  return {
    establishment,
    whitelabel,
    loading,
    error,
    fetchEstablishment,
    updateEstablishment,
    fetchWhitelabel,
    updateWhitelabel,
    uploadLogo,
  }
}

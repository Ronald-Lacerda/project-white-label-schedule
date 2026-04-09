export interface Service {
  id: string
  name: string
  description: string | null
  duration_minutes: number
  price_cents: number | null
  active: boolean
  display_order: number
  created_at: string
}

export function useServices() {
  const api = useApi()

  const services = ref<Service[]>([])
  const loading = ref(false)
  const error = ref('')

  function normalizeServices(data: Service[] | null | undefined): Service[] {
    return Array.isArray(data) ? data : []
  }

  async function fetchServices() {
    loading.value = true
    error.value = ''
    try {
      services.value = normalizeServices(await api.get<Service[] | null>('/api/v1/services'))
    } catch {
      services.value = []
      error.value = 'Erro ao carregar serviços.'
    } finally {
      loading.value = false
    }
  }

  const activeServices = computed(() => services.value.filter(service => service.active))

  async function createService(data: {
    name: string
    description?: string
    duration_minutes: number
    price_cents?: number
    display_order?: number
  }) {
    const s = await api.post<Service>('/api/v1/services', data)
    services.value.push(s)
    return s
  }

  async function updateService(id: string, data: Partial<Service>) {
    const s = await api.put<Service>(`/api/v1/services/${id}`, data)
    const idx = services.value.findIndex(x => x.id === id)
    if (idx !== -1) services.value[idx] = s
    return s
  }

  async function deleteService(id: string) {
    await api.del(`/api/v1/services/${id}`)
    services.value = services.value.filter(s => s.id !== id)
  }

  function formatPrice(cents: number | null): string {
    if (cents == null) return '-'
    return (cents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
  }

  function formatDuration(minutes: number): string {
    if (minutes < 60) return `${minutes} min`
    const h = Math.floor(minutes / 60)
    const m = minutes % 60
    return m === 0 ? `${h}h` : `${h}h ${m}min`
  }

  return {
    services,
    activeServices,
    loading,
    error,
    fetchServices,
    createService,
    updateService,
    deleteService,
    formatPrice,
    formatDuration,
  }
}

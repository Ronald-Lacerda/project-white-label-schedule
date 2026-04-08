export interface Professional {
  id: string
  name: string
  avatar_url: string | null
  phone: string | null
  display_order: number
  active: boolean
  created_at: string
}

export interface ProfessionalHour {
  id: string
  day_of_week: number
  start_time: string
  end_time: string
  is_unavailable: boolean
}

export function useProfessionals() {
  const api = useApi()

  const professionals = ref<Professional[]>([])
  const loading = ref(false)
  const error = ref('')

  function normalizeProfessionals(data: Professional[] | null | undefined): Professional[] {
    return Array.isArray(data) ? data : []
  }

  async function fetchProfessionals() {
    loading.value = true
    error.value = ''
    try {
      professionals.value = normalizeProfessionals(await api.get<Professional[] | null>('/api/v1/professionals'))
    } catch {
      professionals.value = []
      error.value = 'Erro ao carregar profissionais.'
    } finally {
      loading.value = false
    }
  }

  async function createProfessional(data: { name: string; phone?: string; display_order?: number }) {
    const p = await api.post<Professional>('/api/v1/professionals', data)
    professionals.value.push(p)
    return p
  }

  async function updateProfessional(id: string, data: Partial<Professional>) {
    const p = await api.put<Professional>(`/api/v1/professionals/${id}`, data)
    const idx = professionals.value.findIndex(x => x.id === id)
    if (idx !== -1) professionals.value[idx] = p
    return p
  }

  async function deleteProfessional(id: string) {
    await api.del(`/api/v1/professionals/${id}`)
    professionals.value = professionals.value.filter(p => p.id !== id)
  }

  async function updateHours(id: string, hours: Omit<ProfessionalHour, 'id'>[]) {
    return api.put<ProfessionalHour[]>(`/api/v1/professionals/${id}/hours`, { hours })
  }

  async function updateServices(id: string, serviceIds: string[]) {
    await api.put(`/api/v1/professionals/${id}/services`, { service_ids: serviceIds })
  }

  return {
    professionals,
    loading,
    error,
    fetchProfessionals,
    createProfessional,
    updateProfessional,
    deleteProfessional,
    updateHours,
    updateServices,
  }
}

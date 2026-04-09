export interface ManagerAppointment {
  id: string
  professional_id: string
  professional_name: string
  service_id: string
  service_name: string
  duration_minutes: number
  client_name: string
  client_email?: string
  client_phone: string
  client_birth_date?: string
  starts_at: string
  ends_at: string
  status: 'confirmed' | 'completed' | 'no_show' | 'cancelled'
  source: string
  notes?: string
  created_at: string
}

export interface ManagerBlockedPeriod {
  id: string
  professional_id: string
  professional_name: string
  starts_at: string
  ends_at: string
  reason?: string
}

export interface AppointmentFilters {
  date_from?: string
  date_to?: string
  professional_id?: string
  status?: string
  page?: number
  per_page?: number
}

export interface AppointmentsMeta {
  page: number
  per_page: number
  total: number
}

export interface BlockedPeriodInput {
  professional_id: string
  starts_at: string
  ends_at: string
  reason?: string
}


export function useAppointments() {
  const api = useApi()

  const appointments = ref<ManagerAppointment[]>([])
  const meta = ref<AppointmentsMeta>({ page: 1, per_page: 20, total: 0 })
  const blockedPeriods = ref<ManagerBlockedPeriod[]>([])
  const loading = ref(false)
  const error = ref('')

  function buildQuery(filters: AppointmentFilters): string {
    const params = new URLSearchParams()
    if (filters.date_from) params.set('date_from', filters.date_from)
    if (filters.date_to) params.set('date_to', filters.date_to)
    if (filters.professional_id) params.set('professional_id', filters.professional_id)
    if (filters.status) params.set('status', filters.status)
    if (filters.page) params.set('page', String(filters.page))
    if (filters.per_page) params.set('per_page', String(filters.per_page))
    const qs = params.toString()
    return qs ? `?${qs}` : ''
  }

  async function fetchAppointments(filters: AppointmentFilters = {}) {
    loading.value = true
    error.value = ''
    try {
      const response = await api.list<ManagerAppointment>(`/api/v1/appointments${buildQuery(filters)}`)
      appointments.value = response.data ?? []
      meta.value = response.meta ?? { page: 1, per_page: 20, total: 0 }
    } catch {
      appointments.value = []
      error.value = 'Erro ao carregar agendamentos.'
    } finally {
      loading.value = false
    }
  }

  async function fetchAppointment(id: string): Promise<ManagerAppointment> {
    return api.get<ManagerAppointment>(`/api/v1/appointments/${id}`)
  }

  async function updateStatus(id: string, status: string): Promise<ManagerAppointment> {
    const updated = await api.patch<ManagerAppointment>(`/api/v1/appointments/${id}/status`, { status })
    const idx = appointments.value.findIndex(a => a.id === id)
    if (idx !== -1) appointments.value[idx] = updated
    return updated
  }

  async function fetchBlockedPeriods(professionalId?: string, dateFrom?: string, dateTo?: string) {
    const params = new URLSearchParams()
    if (professionalId) params.set('professional_id', professionalId)
    if (dateFrom) params.set('date_from', dateFrom)
    if (dateTo) params.set('date_to', dateTo)
    const qs = params.toString()
    const path = `/api/v1/blocked-periods${qs ? `?${qs}` : ''}`
    blockedPeriods.value = (await api.get<ManagerBlockedPeriod[]>(path)) ?? []
  }

  async function createBlockedPeriod(input: BlockedPeriodInput): Promise<ManagerBlockedPeriod> {
    const period = await api.post<ManagerBlockedPeriod>('/api/v1/blocked-periods', input)
    blockedPeriods.value.push(period)
    return period
  }

  async function deleteBlockedPeriod(id: string) {
    await api.del(`/api/v1/blocked-periods/${id}`)
    blockedPeriods.value = blockedPeriods.value.filter(p => p.id !== id)
  }

  return {
    appointments,
    meta,
    blockedPeriods,
    loading,
    error,
    fetchAppointments,
    fetchAppointment,
    updateStatus,
    fetchBlockedPeriods,
    createBlockedPeriod,
    deleteBlockedPeriod,
  }
}

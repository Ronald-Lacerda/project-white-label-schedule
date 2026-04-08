export interface BusinessHour {
  id: string
  day_of_week: number
  open_time: string
  close_time: string
  is_closed: boolean
}

const DAY_NAMES = ['Domingo', 'Segunda', 'Terça', 'Quarta', 'Quinta', 'Sexta', 'Sábado']

export function useBusinessHours() {
  const api = useApi()

  const hours = ref<BusinessHour[]>([])
  const loading = ref(false)
  const error = ref('')

  // Inicializa os 7 dias com valores padrão se não existirem
  function initDefaults(data: BusinessHour[]): BusinessHour[] {
    return Array.from({ length: 7 }, (_, i) => {
      const existing = data.find(h => h.day_of_week === i)
      return existing ?? {
        id: '',
        day_of_week: i,
        open_time: '08:00:00',
        close_time: '18:00:00',
        is_closed: i === 0, // domingo fechado por padrão
      }
    })
  }

  async function fetchHours() {
    loading.value = true
    error.value = ''
    try {
      const data = await api.get<BusinessHour[]>('/api/v1/establishment/business-hours')
      hours.value = initDefaults(data)
    } catch {
      error.value = 'Erro ao carregar horarios.'
    } finally {
      loading.value = false
    }
  }

  async function saveHours() {
    error.value = ''
    try {
      const updated = await api.put<BusinessHour[]>(
        '/api/v1/establishment/business-hours',
        { hours: hours.value },
      )
      hours.value = initDefaults(updated)
      return updated
    } catch {
      error.value = 'Erro ao salvar horarios.'
      throw new Error(error.value)
    }
  }

  function dayName(dayOfWeek: number): string {
    return DAY_NAMES[dayOfWeek] ?? ''
  }

  return { hours, loading, error, fetchHours, saveHours, dayName }
}

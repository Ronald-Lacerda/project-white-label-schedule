<template>
  <div>
    <div v-if="pending" class="py-8">
      <AppSurface tone="default">
        <div class="ds-empty-state">
          <AppStatusPill tone="info">Carregando</AppStatusPill>
          <p class="text-sm">Montando a experiencia de agendamento.</p>
        </div>
      </AppSurface>
    </div>

    <div v-else-if="error || !establishment" class="py-8">
      <AppSurface tone="default">
        <div class="ds-empty-state">
          <AppStatusPill tone="danger">Nao encontrado</AppStatusPill>
          <p class="text-sm">Estabelecimento nao encontrado.</p>
        </div>
      </AppSurface>
    </div>

    <div v-else class="space-y-6">
      <PublicHeader
        :title="establishment.name"
        subtitle="Escolha o servico, selecione o profissional e reserve um horario disponivel sem criar conta."
        :logo-url="whitelabel?.logo_url"
        :primary-color="whitelabel?.primary_color"
        :secondary-color="whitelabel?.secondary_color"
        :active-step="activeStep"
      />

      <AppSurface tone="default" padding="lg">
        <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p class="ds-kicker">Consulta rapida</p>
            <h2 class="mt-2 text-2xl font-semibold" style="color: var(--color-text);">Consultar ou cancelar agendamento</h2>
          </div>
          <p class="text-sm" style="color: var(--color-text-muted);">Use o codigo e o telefone informados no agendamento.</p>
        </div>

        <form class="mt-5 grid gap-3 sm:grid-cols-[1fr_1fr_auto]" @submit.prevent="fetchAppointmentLookup">
          <input v-model="lookupForm.id" type="text" class="ds-input" placeholder="Codigo do agendamento" />
          <input v-model="lookupForm.phone" type="tel" class="ds-input" placeholder="Telefone" />
          <AppButton type="submit" variant="primary" :disabled="lookupLoading">
            {{ lookupLoading ? 'Buscando...' : 'Consultar' }}
          </AppButton>
        </form>

        <div v-if="lookupError" class="mt-4 rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
          {{ lookupError }}
        </div>

        <div v-if="lookupResult" class="mt-5 rounded-[1.5rem] border p-5" style="border-color: var(--color-border); background: var(--color-surface-muted);">
          <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div class="grid gap-3 text-sm sm:grid-cols-2" style="color: var(--color-text-muted);">
              <p><span class="font-semibold" style="color: var(--color-text);">Codigo:</span> {{ lookupResult.id }}</p>
              <p><span class="font-semibold" style="color: var(--color-text);">Profissional:</span> {{ professionalName(lookupResult.professional_id) }}</p>
              <p><span class="font-semibold" style="color: var(--color-text);">Data:</span> {{ formatDateLong(lookupResult.starts_at) }}</p>
              <p><span class="font-semibold" style="color: var(--color-text);">Horario:</span> {{ formatTime(lookupResult.starts_at) }}</p>
            </div>

            <div class="flex flex-col items-start gap-3">
              <AppStatusPill :tone="statusTone(lookupResult.status)">
                {{ humanStatus(lookupResult.status) }}
              </AppStatusPill>

              <div v-if="lookupResult.can_cancel" class="flex flex-wrap gap-2">
                <AppButton variant="danger" :disabled="cancelLoading" @click="cancelLookupAppointment">
                  {{ cancelLoading ? 'Cancelando...' : 'Cancelar agendamento' }}
                </AppButton>

                <AppButton variant="secondary" @click="toggleReschedule">
                  {{ rescheduleOpen ? 'Fechar reagendamento' : 'Reagendar' }}
                </AppButton>
              </div>

              <p v-else class="max-w-xs text-sm" style="color: var(--color-text-muted);">
                O cancelamento online exige antecedencia minima de {{ lookupResult.min_advance_cancel_hours }}h.
              </p>
            </div>
          </div>

          <div v-if="rescheduleOpen && lookupResult.can_cancel" class="mt-5 border-t pt-5" style="border-color: var(--color-border);">
            <p class="text-sm font-semibold" style="color: var(--color-text);">Escolha um novo horario</p>
            <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
              Vamos procurar disponibilidade para o mesmo servico e profissional.
            </p>

            <div class="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
              <AppSelectableCard
                v-for="date in dateOptions"
                :key="`reschedule-${date.value}`"
                :title="date.label"
                :description="date.weekday"
                :selected="rescheduleDate === date.value"
                @select="rescheduleDate = date.value"
              />
            </div>

            <div v-if="rescheduleError" class="mt-4 rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
              {{ rescheduleError }}
            </div>

            <div class="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
              <AppSelectableCard
                v-for="slot in rescheduleSlotsForDate"
                :key="`reschedule-slot-${slot.professional_id}-${slot.starts_at}`"
                :title="formatTime(slot.starts_at)"
                :description="professionalName(slot.professional_id)"
                :selected="rescheduleSlotKey === `${slot.professional_id}-${slot.starts_at}`"
                @select="selectRescheduleSlot(slot)"
              />
            </div>

            <div v-if="rescheduleSlotsForDate.length === 0" class="mt-4 rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
              Nenhum horario disponivel para a data escolhida.
            </div>

            <AppButton class="mt-4" variant="primary" :disabled="!rescheduleSelectedSlot || rescheduleLoading" @click="submitReschedule">
              {{ rescheduleLoading ? 'Reagendando...' : 'Confirmar novo horario' }}
            </AppButton>
          </div>
        </div>
      </AppSurface>

      <AppSurface tone="default" padding="lg">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="ds-kicker">Passo 1</p>
            <h2 class="mt-2 text-2xl font-semibold" style="color: var(--color-text);">Escolha o servico</h2>
          </div>
          <AppStatusPill v-if="servicesLoading" tone="info">Carregando</AppStatusPill>
        </div>

        <div v-if="services.length === 0" class="mt-4 rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
          Nenhum servico disponivel no momento.
        </div>

        <div v-else class="mt-4 space-y-3">
          <AppSelectableCard
            v-for="service in services"
            :key="service.id"
            :title="service.name"
            :description="service.description || ''"
            :selected="selectedService?.id === service.id"
            @select="selectService(service)"
          >
            <template #meta>
              <p class="text-sm font-semibold">{{ formatDuration(service.duration_minutes) }}</p>
              <p class="mt-1 text-sm" :style="metaStyle(selectedService?.id === service.id)">
                {{ formatPrice(service.price_cents) }}
              </p>
            </template>
          </AppSelectableCard>
        </div>
      </AppSurface>

      <AppSurface tone="default" padding="lg">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="ds-kicker">Passo 2</p>
            <h2 class="mt-2 text-2xl font-semibold" style="color: var(--color-text);">Escolha o profissional</h2>
          </div>
          <AppStatusPill v-if="professionalsLoading" tone="info">Carregando</AppStatusPill>
        </div>

        <p v-if="!selectedService" class="mt-4 text-sm" style="color: var(--color-text-muted);">
          Selecione um servico para ver os profissionais disponiveis.
        </p>

        <div v-else class="mt-4 space-y-3">
          <AppSelectableCard
            title="Qualquer profissional disponivel"
            description="Mostra os melhores horarios sem limitar a um nome especifico."
            :selected="selectedProfessionalId === ''"
            @select="selectedProfessionalId = ''"
          />

          <AppSelectableCard
            v-for="professional in professionals"
            :key="professional.id"
            :title="professional.name"
            :description="selectedProfessionalId === professional.id ? 'Selecionado para este agendamento.' : 'Disponivel para este servico.'"
            :selected="selectedProfessionalId === professional.id"
            @select="selectedProfessionalId = professional.id"
          />

          <div
            v-if="selectedService && !professionalsLoading && professionals.length === 0"
            class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm"
            style="border-color: var(--color-border); color: var(--color-text-muted);"
          >
            Nenhum profissional foi vinculado a este servico ainda.
          </div>
        </div>
      </AppSurface>

      <AppSurface tone="default" padding="lg">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="ds-kicker">Passo 3</p>
            <h2 class="mt-2 text-2xl font-semibold" style="color: var(--color-text);">Escolha o horario</h2>
          </div>
          <AppStatusPill v-if="availabilityLoading" tone="info">Buscando horarios</AppStatusPill>
        </div>

        <p v-if="!selectedService" class="mt-4 text-sm" style="color: var(--color-text-muted);">
          O calendario de horarios aparece depois da escolha do servico.
        </p>

        <div v-else class="mt-4 space-y-5">
          <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
            <AppSelectableCard
              v-for="date in dateOptions"
              :key="date.value"
              :title="date.label"
              :description="date.weekday"
              :selected="selectedDate === date.value"
              @select="selectedDate = date.value"
            />
          </div>

          <div v-if="availabilityError" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(184, 107, 22, 0.24); background: var(--color-warning-soft); color: var(--color-warning);">
            {{ availabilityError }}
          </div>

          <div v-if="slotsForSelectedDate.length === 0 && !availabilityLoading" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
            Nenhum horario disponivel para a data selecionada.
          </div>

          <div v-else class="grid grid-cols-2 gap-2 sm:grid-cols-4">
            <AppSelectableCard
              v-for="slot in slotsForSelectedDate"
              :key="`${slot.professional_id}-${slot.starts_at}`"
              :title="formatTime(slot.starts_at)"
              :description="professionalName(slot.professional_id)"
              :selected="selectedSlotKey === `${slot.professional_id}-${slot.starts_at}`"
              @select="selectSlot(slot)"
            />
          </div>
        </div>
      </AppSurface>

      <BookingSummary
        v-if="selectedService && selectedSlot"
        :service-name="selectedService.name"
        :professional-name="professionalName(selectedSlot.professional_id)"
        :date-label="formatDateLong(selectedSlot.starts_at)"
        :time-label="formatTime(selectedSlot.starts_at)"
      />

      <AppSurface v-if="selectedService && selectedSlot && !bookingResult" tone="default" padding="lg">
        <div>
          <p class="ds-kicker">Passo 4</p>
          <h2 class="mt-2 text-2xl font-semibold" style="color: var(--color-text);">Seus dados</h2>
        </div>

        <form class="mt-5 space-y-4" @submit.prevent="submitAppointment">
          <div>
            <label class="ds-label">Nome</label>
            <input v-model="bookingForm.client_name" type="text" required class="ds-input" placeholder="Seu nome completo" />
          </div>
          <div>
            <label class="ds-label">Telefone</label>
            <input v-model="bookingForm.client_phone" type="tel" required class="ds-input" placeholder="(11) 99999-9999" />
          </div>

          <div v-if="bookingError" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
            {{ bookingError }}
          </div>

          <AppButton type="submit" variant="primary" block :disabled="bookingLoading">
            {{ bookingLoading ? 'Confirmando agendamento...' : 'Confirmar agendamento' }}
          </AppButton>
        </form>
      </AppSurface>

      <AppSurface v-if="bookingResult" tone="brand" padding="lg">
        <p class="ds-kicker">Passo concluido</p>
        <h2 class="mt-2 font-display text-4xl font-semibold tracking-tight" style="color: var(--color-brand-text);">Agendamento confirmado</h2>
        <div class="mt-5 grid gap-3 text-sm md:grid-cols-2">
          <p><span class="font-semibold">Codigo:</span> {{ bookingResult.id }}</p>
          <p><span class="font-semibold">Servico:</span> {{ selectedService?.name }}</p>
          <p><span class="font-semibold">Profissional:</span> {{ professionalName(bookingResult.professional_id) }}</p>
          <p><span class="font-semibold">Data:</span> {{ formatDateLong(bookingResult.starts_at) }}</p>
        </div>
      </AppSurface>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'booking' })

const route = useRoute()
const config = useRuntimeConfig()
const slug = route.params.slug as string

interface EstablishmentPublic {
  id: string
  name: string
  slug: string
}

interface WhitelabelPublic {
  logo_url: string | null
  primary_color: string
  secondary_color: string | null
}

interface PublicResponse {
  establishment: EstablishmentPublic
  whitelabel: WhitelabelPublic | null
}

interface ServiceItem {
  id: string
  name: string
  description: string | null
  duration_minutes: number
  price_cents: number | null
}

interface ProfessionalItem {
  id: string
  name: string
}

interface SlotItem {
  starts_at: string
  ends_at: string
  professional_id: string
}

interface BookingResult {
  id: string
  service_id: string
  professional_id: string
  client_name: string
  client_phone: string
  starts_at: string
  ends_at: string
  status: string
}

interface LookupResult extends BookingResult {
  can_cancel: boolean
  min_advance_cancel_hours: number
}

const { data: response, pending, error } = useFetch<{ data: PublicResponse }>(
  `${config.public.apiBaseUrl}/pub/${slug}`,
  {
    default: () => null,
  },
)

const { data: servicesResponse, pending: servicesLoading } = useFetch<{ data: ServiceItem[] }>(
  `${config.public.apiBaseUrl}/pub/${slug}/services`,
  {
    default: () => null,
  },
)

const establishment = computed(() => response.value?.data?.establishment ?? null)
const whitelabel = computed(() => response.value?.data?.whitelabel ?? null)
const services = computed(() => servicesResponse.value?.data ?? [])

const professionals = ref<ProfessionalItem[]>([])
const professionalsLoading = ref(false)

const availabilityByDate = ref<Record<string, SlotItem[]>>({})
const availabilityLoading = ref(false)
const availabilityError = ref('')

const selectedServiceId = ref('')
const selectedProfessionalId = ref('')
const selectedDate = ref('')
const selectedSlotKey = ref('')
const selectedSlot = ref<SlotItem | null>(null)
const bookingResult = ref<BookingResult | null>(null)
const bookingLoading = ref(false)
const bookingError = ref('')
const bookingForm = reactive({
  client_name: '',
  client_phone: '',
})
const lookupLoading = ref(false)
const lookupError = ref('')
const cancelLoading = ref(false)
const lookupResult = ref<LookupResult | null>(null)
const lookupForm = reactive({
  id: '',
  phone: '',
})
const rescheduleOpen = ref(false)
const rescheduleLoading = ref(false)
const rescheduleError = ref('')
const rescheduleDate = ref('')
const rescheduleAvailabilityByDate = ref<Record<string, SlotItem[]>>({})
const rescheduleSelectedSlot = ref<SlotItem | null>(null)
const rescheduleSlotKey = ref('')

const selectedService = computed(() => services.value.find(service => service.id === selectedServiceId.value) ?? null)
const activeStep = computed(() => {
  if (bookingResult.value || selectedSlot.value) return 4
  if (selectedService.value) return 3
  return 1
})

const dateOptions = computed(() => {
  const formatterWeekday = new Intl.DateTimeFormat('pt-BR', { weekday: 'short' })
  const formatterLabel = new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: '2-digit' })
  const today = new Date()

  return Array.from({ length: 7 }, (_, index) => {
    const date = new Date(today)
    date.setDate(today.getDate() + index)
    return {
      value: formatLocalDateInputValue(date),
      weekday: formatterWeekday.format(date).replace('.', ''),
      label: formatterLabel.format(date),
    }
  })
})

const slotsForSelectedDate = computed(() => availabilityByDate.value[selectedDate.value] ?? [])
const rescheduleSlotsForDate = computed(() => rescheduleAvailabilityByDate.value[rescheduleDate.value] ?? [])

watchEffect(() => {
  applyBrandTheme({
    primaryColor: whitelabel.value?.primary_color,
    secondaryColor: whitelabel.value?.secondary_color,
  })
})

onBeforeUnmount(() => {
  applyBrandTheme()
})

watchEffect(() => {
  if (!selectedDate.value && dateOptions.value.length > 0) {
    selectedDate.value = dateOptions.value[0].value
  }
})

watchEffect(() => {
  if (!rescheduleDate.value && dateOptions.value.length > 0) {
    rescheduleDate.value = dateOptions.value[0].value
  }
})

watch(selectedServiceId, async (newServiceId) => {
  selectedProfessionalId.value = ''
  selectedSlotKey.value = ''
  selectedSlot.value = null
  bookingResult.value = null
  bookingError.value = ''
  availabilityByDate.value = {}
  availabilityError.value = ''

  if (!newServiceId) {
    professionals.value = []
    return
  }

  await fetchProfessionals(newServiceId)
  await fetchAvailability()
})

watch(selectedProfessionalId, async () => {
  selectedSlotKey.value = ''
  selectedSlot.value = null
  bookingResult.value = null
  bookingError.value = ''

  if (!selectedServiceId.value) return
  await fetchAvailability()
})

watch(selectedDate, () => {
  selectedSlotKey.value = ''
  selectedSlot.value = null
  bookingResult.value = null
  bookingError.value = ''
})

function metaStyle(selected: boolean) {
  return {
    color: selected ? 'rgba(255, 255, 255, 0.72)' : 'var(--color-text-muted)',
  }
}

function selectService(service: ServiceItem) {
  selectedServiceId.value = service.id
}

function selectSlot(slot: SlotItem) {
  selectedSlot.value = slot
  selectedSlotKey.value = `${slot.professional_id}-${slot.starts_at}`
  bookingResult.value = null
  bookingError.value = ''
}

async function fetchProfessionals(serviceId: string) {
  professionalsLoading.value = true
  try {
    const data = await $fetch<{ data: ProfessionalItem[] }>(
      `${config.public.apiBaseUrl}/pub/${slug}/professionals`,
      {
        query: { service_id: serviceId },
      },
    )
    professionals.value = data.data
  } finally {
    professionalsLoading.value = false
  }
}

async function fetchAvailability() {
  if (!selectedServiceId.value) return

  availabilityLoading.value = true
  availabilityError.value = ''

  try {
    const query: Record<string, string> = {
      service_id: selectedServiceId.value,
      date_from: dateOptions.value[0]?.value ?? selectedDate.value,
      date_to: dateOptions.value[dateOptions.value.length - 1]?.value ?? selectedDate.value,
    }

    if (selectedProfessionalId.value) {
      query.professional_id = selectedProfessionalId.value
    }

    const data = await $fetch<{ data: Record<string, SlotItem[]> }>(
      `${config.public.apiBaseUrl}/pub/${slug}/availability`,
      { query },
    )

    availabilityByDate.value = data.data
  } catch (err: any) {
    availabilityByDate.value = {}
    availabilityError.value = err?.data?.error?.message ?? 'Nao foi possivel carregar os horarios.'
  } finally {
    availabilityLoading.value = false
  }
}

async function submitAppointment() {
  if (!selectedService.value || !selectedSlot.value) return

  bookingLoading.value = true
  bookingError.value = ''

  try {
    const idempotencyKey = globalThis.crypto?.randomUUID?.() ?? `${Date.now()}-${Math.random()}`
    const data = await $fetch<{ data: BookingResult }>(
      `${config.public.apiBaseUrl}/pub/${slug}/appointments`,
      {
        method: 'POST',
        body: {
          service_id: selectedService.value.id,
          professional_id: selectedSlot.value.professional_id,
          starts_at: selectedSlot.value.starts_at,
          client_name: bookingForm.client_name,
          client_phone: bookingForm.client_phone,
          idempotency_key: idempotencyKey,
        },
      },
    )

    bookingResult.value = data.data
    lookupResult.value = {
      ...data.data,
      can_cancel: true,
      min_advance_cancel_hours: 0,
    }
    lookupForm.id = data.data.id
    lookupForm.phone = data.data.client_phone
    await fetchAvailability()
  } catch (err: any) {
    bookingError.value = err?.data?.error?.message ?? 'Nao foi possivel concluir o agendamento.'
  } finally {
    bookingLoading.value = false
  }
}

async function fetchAppointmentLookup() {
  lookupLoading.value = true
  lookupError.value = ''

  try {
    const data = await $fetch<{ data: LookupResult }>(
      `${config.public.apiBaseUrl}/pub/${slug}/appointments/${lookupForm.id}`,
      {
        query: { phone: lookupForm.phone },
      },
    )
    lookupResult.value = data.data
    rescheduleOpen.value = false
    rescheduleError.value = ''
    rescheduleAvailabilityByDate.value = {}
    rescheduleSelectedSlot.value = null
    rescheduleSlotKey.value = ''
  } catch (err: any) {
    lookupResult.value = null
    lookupError.value = err?.data?.error?.message ?? 'Nao foi possivel localizar o agendamento.'
  } finally {
    lookupLoading.value = false
  }
}

async function cancelLookupAppointment() {
  if (!lookupResult.value) return

  cancelLoading.value = true
  lookupError.value = ''

  try {
    const data = await $fetch<{ data: LookupResult }>(
      `${config.public.apiBaseUrl}/pub/${slug}/appointments/${lookupResult.value.id}/cancel`,
      {
        method: 'PATCH',
        body: { phone: lookupForm.phone },
      },
    )
    lookupResult.value = data.data
    rescheduleOpen.value = false
    await fetchAvailability()
  } catch (err: any) {
    lookupError.value = err?.data?.error?.message ?? 'Nao foi possivel cancelar o agendamento.'
  } finally {
    cancelLoading.value = false
  }
}

async function toggleReschedule() {
  rescheduleOpen.value = !rescheduleOpen.value
  rescheduleError.value = ''

  if (!rescheduleOpen.value || !lookupResult.value) return

  rescheduleSelectedSlot.value = null
  rescheduleSlotKey.value = ''
  await fetchRescheduleAvailability()
}

async function fetchRescheduleAvailability() {
  if (!lookupResult.value) return

  rescheduleError.value = ''

  try {
    const query: Record<string, string> = {
      service_id: lookupResult.value.service_id,
      professional_id: lookupResult.value.professional_id,
      date_from: dateOptions.value[0]?.value ?? '',
      date_to: dateOptions.value[dateOptions.value.length - 1]?.value ?? '',
    }

    const data = await $fetch<{ data: Record<string, SlotItem[]> }>(
      `${config.public.apiBaseUrl}/pub/${slug}/availability`,
      { query },
    )

    rescheduleAvailabilityByDate.value = data.data
  } catch (err: any) {
    rescheduleAvailabilityByDate.value = {}
    rescheduleError.value = err?.data?.error?.message ?? 'Nao foi possivel carregar horarios para reagendamento.'
  }
}

function selectRescheduleSlot(slot: SlotItem) {
  rescheduleSelectedSlot.value = slot
  rescheduleSlotKey.value = `${slot.professional_id}-${slot.starts_at}`
}

async function submitReschedule() {
  if (!lookupResult.value || !rescheduleSelectedSlot.value) return

  rescheduleLoading.value = true
  rescheduleError.value = ''

  try {
    const data = await $fetch<{ data: BookingResult }>(
      `${config.public.apiBaseUrl}/pub/${slug}/appointments/${lookupResult.value.id}/reschedule`,
      {
        method: 'PATCH',
        body: {
          phone: lookupForm.phone,
          starts_at: rescheduleSelectedSlot.value.starts_at,
        },
      },
    )

    lookupForm.id = data.data.id
    const refreshed = await $fetch<{ data: LookupResult }>(
      `${config.public.apiBaseUrl}/pub/${slug}/appointments/${data.data.id}`,
      {
        query: { phone: lookupForm.phone },
      },
    )
    lookupResult.value = refreshed.data
    rescheduleOpen.value = false
    rescheduleSelectedSlot.value = null
    rescheduleSlotKey.value = ''
    await fetchAvailability()
  } catch (err: any) {
    rescheduleError.value = err?.data?.error?.message ?? 'Nao foi possivel reagendar.'
  } finally {
    rescheduleLoading.value = false
  }
}

function formatPrice(priceCents: number | null) {
  if (priceCents == null) return 'Preco sob consulta'
  return (priceCents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

function formatDuration(minutes: number) {
  if (minutes < 60) return `${minutes} min`
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return remainingMinutes === 0 ? `${hours}h` : `${hours}h ${remainingMinutes}min`
}

function formatTime(isoString: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(isoString))
}

function formatDateLong(isoString: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    weekday: 'long',
    day: '2-digit',
    month: 'long',
  }).format(new Date(isoString))
}

function professionalName(professionalId: string) {
  if (!professionalId) return 'Qualquer profissional'
  return professionals.value.find(professional => professional.id === professionalId)?.name ?? 'Profissional'
}

function humanStatus(status: string) {
  if (status === 'confirmed') return 'Confirmado'
  if (status === 'cancelled') return 'Cancelado'
  if (status === 'completed') return 'Concluido'
  if (status === 'no_show') return 'Nao compareceu'
  return status
}

function statusTone(status: string) {
  return {
    confirmed: 'info',
    cancelled: 'danger',
    completed: 'success',
    no_show: 'warning',
  }[status] as 'info' | 'danger' | 'success' | 'warning'
}

useHead(() => ({
  title: establishment.value?.name ?? 'Agendamento',
}))
</script>

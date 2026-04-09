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

    <div v-else class="space-y-6 pb-24">
      <AppSurface tone="brand" padding="lg">
        <div class="flex items-center justify-between gap-4" :style="brandThemeStyle({
          primaryColor: whitelabel?.primary_color,
          secondaryColor: whitelabel?.secondary_color,
        })">
          <div class="flex min-w-0 items-center gap-4">
            <div
              class="flex h-14 w-14 flex-shrink-0 items-center justify-center overflow-hidden rounded-[1.4rem] border bg-white/85 shadow-sm"
              style="border-color: rgba(var(--color-brand-primary-rgb), 0.14);"
            >
              <img
                v-if="whitelabel?.logo_url"
                :src="whitelabel.logo_url"
                :alt="establishment.name"
                class="h-full w-full object-cover"
              />
              <span
                v-else
                class="text-xs font-semibold uppercase tracking-[0.28em]"
                style="color: var(--color-brand-primary);"
              >
                {{ establishment.name.slice(0, 2).toUpperCase() }}
              </span>
            </div>

            <div class="min-w-0">
              <p class="ds-kicker">Agendamento online</p>
              <h1 class="truncate text-2xl font-semibold" style="color: var(--color-text);">
                {{ establishment.name }}
              </h1>
            </div>
          </div>

          <AppStatusPill tone="brand">Passo {{ activeStep }} de 5</AppStatusPill>
        </div>
      </AppSurface>

      <AppSurface v-if="summaryItems.length" tone="default" padding="lg">
        <div class="flex items-start justify-between gap-4">
          <div>
            <p class="ds-kicker">Resumo</p>
            <h2 class="mt-2 text-xl font-semibold" style="color: var(--color-text);">Seu agendamento em progresso</h2>
          </div>
          <AppStatusPill tone="brand">Resumo</AppStatusPill>
        </div>

        <div class="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <div
            v-for="item in summaryItems"
            :key="item.label"
            class="rounded-[1.2rem] border p-4"
            style="border-color: var(--color-border); background: var(--color-surface-muted);"
          >
            <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">{{ item.label }}</p>
            <p class="mt-2 text-sm font-semibold" style="color: var(--color-text);">{{ item.value }}</p>
          </div>
        </div>
      </AppSurface>

      <AppSurface v-if="bookingResult" tone="brand" padding="lg">
        <p class="ds-kicker">Agendamento confirmado</p>
        <h2 class="mt-2 font-display text-4xl font-semibold tracking-tight" style="color: var(--color-brand-text);">Tudo certo por aqui</h2>
        <p class="mt-3 text-sm" style="color: var(--color-text-muted);">
          Seu horario foi reservado e o codigo abaixo pode ser usado em contatos futuros com o estabelecimento.
        </p>

        <div class="mt-5 grid gap-3 text-sm md:grid-cols-2">
          <p><span class="font-semibold">Codigo:</span> {{ bookingResult.id }}</p>
          <p><span class="font-semibold">Servico:</span> {{ selectedService?.name }}</p>
          <p><span class="font-semibold">Profissional:</span> {{ professionalName(bookingResult.professional_id) }}</p>
          <p><span class="font-semibold">Data:</span> {{ formatDateLong(bookingResult.starts_at) }}</p>
          <p><span class="font-semibold">Horario:</span> {{ formatTime(bookingResult.starts_at) }}</p>
          <p><span class="font-semibold">Cliente:</span> {{ bookingForm.client_name }}</p>
        </div>
      </AppSurface>

      <AppSurface v-else tone="default" padding="lg">
        <div>
          <p class="ds-kicker">Passo {{ currentStep }}</p>
          <h2 class="mt-2 text-2xl font-semibold" style="color: var(--color-text);">{{ currentStepContent.title }}</h2>
          <p class="mt-2 text-sm" style="color: var(--color-text-muted);">{{ currentStepContent.description }}</p>
        </div>

        <div v-if="currentStep === 1" class="mt-6 space-y-4">
          <div v-if="servicesLoading" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
            Carregando servicos.
          </div>

          <div v-else-if="services.length === 0" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
            Nenhum servico disponivel no momento.
          </div>

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

        <div v-else-if="currentStep === 2" class="mt-6 space-y-3">
          <div v-if="professionalsLoading" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
            Carregando profissionais.
          </div>

          <div v-else-if="professionalsError" class="rounded-[1.2rem] border px-4 py-6 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
            {{ professionalsError }}
          </div>

          <div v-else-if="professionals.length === 0" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
            Nenhum profissional foi vinculado a este servico.
          </div>

          <AppSelectableCard
            v-for="professional in professionals"
            :key="professional.id"
            :title="professional.name"
            :description="selectedProfessionalId === professional.id ? 'Selecionado para este agendamento.' : 'Disponivel para este servico.'"
            :selected="selectedProfessionalId === professional.id"
            @select="selectProfessional(professional.id)"
          />
        </div>

        <div v-else-if="currentStep === 3" class="mt-6 space-y-6">
          <section class="space-y-3 rounded-[1.4rem] border p-4 sm:p-5" style="border-color: var(--color-border); background: var(--color-surface-muted);">
            <div class="flex items-center justify-between gap-3">
              <div>
                <p class="text-sm font-semibold" style="color: var(--color-text);">Escolha o dia</p>
                <p class="mt-1 text-xs" style="color: var(--color-text-muted);">
                  Primeiro selecione a data para ver os horarios disponiveis.
                </p>
              </div>
              <span class="rounded-full px-3 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.2em]" style="background: rgba(var(--color-brand-primary-rgb), 0.08); color: var(--color-brand-primary);">
                Datas
              </span>
            </div>

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
          </section>

          <section class="space-y-3 rounded-[1.4rem] border p-4 sm:p-5" style="border-color: var(--color-border); background: var(--color-surface);">
            <div class="flex items-center justify-between gap-3">
              <div>
                <p class="text-sm font-semibold" style="color: var(--color-text);">Horarios livres</p>
                <p class="mt-1 text-xs" style="color: var(--color-text-muted);">
                  {{ selectedDateLabel }}
                </p>
              </div>
              <span class="rounded-full px-3 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.2em]" style="background: rgba(15, 23, 42, 0.05); color: var(--color-text-soft);">
                Horarios
              </span>
            </div>

            <div v-if="availabilityLoading" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
              Buscando horarios disponiveis.
            </div>

            <div v-else-if="availabilityError" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(184, 107, 22, 0.24); background: var(--color-warning-soft); color: var(--color-warning);">
              {{ availabilityError }}
            </div>

            <div v-else-if="slotsForSelectedDate.length === 0" class="rounded-[1.2rem] border border-dashed px-4 py-6 text-sm" style="border-color: var(--color-border); color: var(--color-text-muted);">
              Nenhum horario disponivel para a data selecionada.
            </div>

            <div v-else class="grid grid-cols-2 gap-2 sm:grid-cols-4">
              <AppSelectableCard
                v-for="slot in slotsForSelectedDate"
                :key="`${slot.professional_id}-${slot.starts_at}`"
                :title="formatTime(slot.starts_at)"
                :selected="selectedSlotKey === `${slot.professional_id}-${slot.starts_at}`"
                @select="selectSlot(slot)"
              />
            </div>
          </section>
        </div>

        <div v-else-if="currentStep === 4" class="mt-6 space-y-4">
          <div
            v-for="field in customerFields"
            :key="field.key"
          >
            <label class="ds-label" :for="field.key">{{ field.label }}</label>
            <input
              :id="field.key"
              :value="bookingForm[field.key]"
              :type="field.type"
              :lang="field.type === 'date' ? 'pt-BR' : undefined"
              :autocomplete="field.autocomplete"
              :placeholder="field.placeholder"
              class="ds-input"
              @input="handleCustomerInput(field.key, $event)"
            />
          </div>

          <div v-if="stepError" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
            {{ stepError }}
          </div>
        </div>

        <div v-else class="mt-6 space-y-6">
          <BookingSummary
            v-if="selectedService && selectedSlot"
            title="Revise antes de confirmar"
            badge="Ultimo passo"
            :service-name="selectedService.name"
            :professional-name="professionalName(selectedSlot.professional_id)"
            :date-label="formatDateLong(selectedSlot.starts_at)"
            :time-label="formatTime(selectedSlot.starts_at)"
          />

          <div class="grid gap-4 md:grid-cols-2">
            <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
              <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Dados pessoais</p>
              <p class="mt-2 text-sm font-semibold" style="color: var(--color-text);">{{ bookingForm.client_name }}</p>
              <p class="mt-1 text-sm" style="color: var(--color-text-muted);">{{ bookingForm.client_email }}</p>
              <p class="mt-1 text-sm" style="color: var(--color-text-muted);">{{ formatBrazilianPhoneDisplay(bookingForm.client_phone) }}</p>
              <p class="mt-1 text-sm" style="color: var(--color-text-muted);">Nascimento: {{ formatBirthDate(bookingForm.client_birth_date) }}</p>
            </div>

            <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
              <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Como funciona</p>
              <p class="mt-2 text-sm" style="color: var(--color-text-muted);">
                Ao confirmar, vamos reservar o horario imediatamente e mostrar a tela final com o codigo do agendamento.
              </p>
            </div>
          </div>

          <div v-if="bookingError" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
            {{ bookingError }}
          </div>
        </div>
      </AppSurface>

      <div v-if="!bookingResult" class="sticky bottom-4 z-20">
        <AppSurface tone="default" padding="lg">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div class="text-sm" style="color: var(--color-text-muted);">
              {{ currentStepContent.footer }}
            </div>

            <div class="flex gap-2">
              <AppButton
                v-if="currentStep > 1"
                variant="secondary"
                @click="goBack"
              >
                Voltar
              </AppButton>

              <AppButton
                v-if="currentStep < 5"
                variant="primary"
                :disabled="!canContinue"
                @click="goNext"
              >
                Continuar
              </AppButton>

              <AppButton
                v-else
                variant="primary"
                :disabled="bookingLoading || !selectedSlot"
                @click="submitAppointment"
              >
                {{ bookingLoading ? 'Confirmando...' : 'Confirmar agendamento' }}
              </AppButton>
            </div>
          </div>
        </AppSurface>
      </div>
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

interface ProfessionalItemResponse {
  id?: string
  name?: string
  ID?: string
  Name?: string
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
  client_email?: string
  client_phone: string
  client_birth_date?: string
  starts_at: string
  ends_at: string
  status: string
}

type BookingFormKey = 'client_name' | 'client_email' | 'client_phone' | 'client_birth_date'

const customerFields: Array<{
  key: BookingFormKey
  label: string
  type: string
  placeholder: string
  autocomplete: string
}> = [
  {
    key: 'client_name',
    label: 'Nome',
    type: 'text',
    placeholder: 'Seu nome completo',
    autocomplete: 'name',
  },
  {
    key: 'client_email',
    label: 'Email',
    type: 'email',
    placeholder: 'voce@exemplo.com',
    autocomplete: 'email',
  },
  {
    key: 'client_phone',
    label: 'Telefone',
    type: 'tel',
    placeholder: '(11) 99999-9999',
    autocomplete: 'tel',
  },
  {
    key: 'client_birth_date',
    label: 'Data de nascimento',
    type: 'date',
    placeholder: '',
    autocomplete: 'bday',
  },
]

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

const currentStep = ref(1)
const professionals = ref<ProfessionalItem[]>([])
const professionalsLoading = ref(false)
const professionalsError = ref('')
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
const stepError = ref('')
const bookingForm = reactive<Record<BookingFormKey, string>>({
  client_name: '',
  client_email: '',
  client_phone: '',
  client_birth_date: '',
})

const selectedService = computed(() => services.value.find(service => service.id === selectedServiceId.value) ?? null)
const activeStep = computed(() => bookingResult.value ? 5 : currentStep.value)

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

const selectedDateLabel = computed(() => {
  const selected = dateOptions.value.find(date => date.value === selectedDate.value)
  if (!selected) {
    return 'Selecione uma data para continuar.'
  }

  return `${selected.label} - ${selected.weekday}`
})

const currentStepContent = computed(() => {
  return {
    1: {
      title: 'Escolha o servico',
      description: 'Selecione o atendimento que voce quer agendar.',
      footer: 'Escolha um servico para seguir.',
    },
    2: {
      title: 'Escolha o profissional',
      description: 'Nesta jornada o profissional e obrigatorio antes da busca de horarios.',
      footer: 'Selecione quem vai realizar o atendimento.',
    },
    3: {
      title: 'Escolha o horario',
      description: 'Mostramos apenas horarios disponiveis para o servico e profissional escolhidos.',
      footer: 'Escolha uma data e um horario disponivel.',
    },
    4: {
      title: 'Seus dados',
      description: 'Agora so faltam seus dados para finalizar a reserva com seguranca.',
      footer: 'Preencha nome, email, telefone e nascimento.',
    },
    5: {
      title: 'Confirme seu agendamento',
      description: 'Revise tudo antes de reservar o horario.',
      footer: 'Se estiver tudo certo, confirme agora.',
    },
  }[currentStep.value]!
})

const summaryItems = computed(() => {
  const items: Array<{ label: string; value: string }> = []

  if (bookingForm.client_name) items.push({ label: 'Cliente', value: bookingForm.client_name })
  if (selectedService.value) items.push({ label: 'Servico', value: selectedService.value.name })
  if (selectedProfessionalId.value) items.push({ label: 'Profissional', value: professionalName(selectedProfessionalId.value) })
  if (selectedSlot.value) {
    items.push({ label: 'Data', value: formatDateLong(selectedSlot.value.starts_at) })
    items.push({ label: 'Horario', value: formatTime(selectedSlot.value.starts_at) })
  }

  return items
})

const isDataStepValid = computed(() => {
  return (
    bookingForm.client_name.trim().length > 0 &&
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(bookingForm.client_email.trim()) &&
    bookingForm.client_phone.trim().length > 0 &&
    bookingForm.client_birth_date.trim().length > 0
  )
})

const canContinue = computed(() => {
  if (currentStep.value === 1) return Boolean(selectedService.value)
  if (currentStep.value === 2) return Boolean(selectedProfessionalId.value)
  if (currentStep.value === 3) return Boolean(selectedSlot.value)
  if (currentStep.value === 4) return isDataStepValid.value
  return Boolean(selectedSlot.value)
})

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

watch(selectedServiceId, async (newServiceId) => {
  selectedProfessionalId.value = ''
  selectedSlot.value = null
  selectedSlotKey.value = ''
  availabilityByDate.value = {}
  availabilityError.value = ''
  professionalsError.value = ''

  if (!newServiceId) {
    professionals.value = []
    return
  }

  await fetchProfessionals(newServiceId)
})

watch(selectedProfessionalId, async (newProfessionalId) => {
  selectedSlot.value = null
  selectedSlotKey.value = ''
  availabilityByDate.value = {}
  availabilityError.value = ''

  if (!selectedServiceId.value || !newProfessionalId) return
  await fetchAvailability()
})

watch(selectedDate, () => {
  selectedSlot.value = null
  selectedSlotKey.value = ''
})

function metaStyle(selected: boolean) {
  return {
    color: selected ? 'rgba(255, 255, 255, 0.72)' : 'var(--color-text-muted)',
  }
}

function selectService(service: ServiceItem) {
  selectedServiceId.value = service.id
}

function selectProfessional(professionalId: string) {
  selectedProfessionalId.value = professionalId
}

function selectSlot(slot: SlotItem) {
  selectedSlot.value = slot
  selectedSlotKey.value = `${slot.professional_id}-${slot.starts_at}`
}

function goBack() {
  stepError.value = ''
  bookingError.value = ''
  if (currentStep.value > 1) currentStep.value -= 1
}

function goNext() {
  stepError.value = ''

  if (currentStep.value === 4 && !isDataStepValid.value) {
    stepError.value = 'Preencha nome, email, telefone e data de nascimento para continuar.'
    return
  }

  if (!canContinue.value) return

  if (currentStep.value < 5) currentStep.value += 1
}

async function fetchProfessionals(serviceId: string) {
  professionalsLoading.value = true
  professionalsError.value = ''
  professionals.value = []

  try {
    const data = await $fetch<{ data: ProfessionalItemResponse[] }>(
      `${config.public.apiBaseUrl}/pub/${slug}/professionals`,
      {
        query: { service_id: serviceId },
      },
    )
    professionals.value = Array.isArray(data.data)
      ? data.data
          .map((professional) => ({
            id: professional.id ?? professional.ID ?? '',
            name: professional.name ?? professional.Name ?? '',
          }))
          .filter((professional) => professional.id && professional.name)
      : []
  } catch (err: any) {
    professionalsError.value = err?.data?.error?.message ?? 'Nao foi possivel carregar os profissionais para este servico.'
  } finally {
    professionalsLoading.value = false
  }
}

async function fetchAvailability() {
  if (!selectedServiceId.value || !selectedProfessionalId.value) return

  availabilityLoading.value = true
  availabilityError.value = ''

  try {
    const data = await $fetch<{ data: Record<string, SlotItem[]> }>(
      `${config.public.apiBaseUrl}/pub/${slug}/availability`,
      {
        query: {
          service_id: selectedServiceId.value,
          professional_id: selectedProfessionalId.value,
          date_from: dateOptions.value[0]?.value ?? selectedDate.value,
          date_to: dateOptions.value[dateOptions.value.length - 1]?.value ?? selectedDate.value,
        },
      },
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
          client_email: bookingForm.client_email,
          client_phone: normalizeBrazilianPhone(bookingForm.client_phone),
          client_birth_date: bookingForm.client_birth_date,
          idempotency_key: idempotencyKey,
        },
      },
    )

    bookingResult.value = data.data
    await fetchAvailability()
  } catch (err: any) {
    bookingError.value = err?.data?.error?.message ?? 'Nao foi possivel concluir o agendamento.'
  } finally {
    bookingLoading.value = false
  }
}

function professionalName(professionalId: string) {
  return professionals.value.find(professional => professional.id === professionalId)?.name ?? 'Profissional selecionado'
}

function formatPrice(priceCents: number | null) {
  if (priceCents == null) return 'Preco sob consulta'
  return (priceCents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

function formatDuration(durationMinutes: number) {
  if (durationMinutes >= 60 && durationMinutes % 60 === 0) {
    const hours = durationMinutes / 60
    return hours === 1 ? '1 hora' : `${hours} horas`
  }

  if (durationMinutes > 60) {
    const hours = Math.floor(durationMinutes / 60)
    const minutes = durationMinutes % 60
    return `${hours}h ${minutes}min`
  }

  return `${durationMinutes} min`
}

function formatDateLong(iso: string) {
  return new Date(iso).toLocaleDateString('pt-BR', {
    weekday: 'long',
    day: '2-digit',
    month: 'long',
  })
}

function formatBirthDate(value: string) {
  return new Date(`${value}T00:00:00`).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

function formatTime(iso: string) {
  return formatShortTimeBr(iso)
}

function handleCustomerInput(key: BookingFormKey, event: Event) {
  const input = event.target as HTMLInputElement
  if (key === 'client_phone') {
    bookingForm.client_phone = formatBrazilianPhoneInput(input.value)
    return
  }

  bookingForm[key] = input.value
}
</script>

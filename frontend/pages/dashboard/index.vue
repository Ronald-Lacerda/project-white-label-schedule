<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
      <div class="space-y-2">
        <h1 class="ds-title">Início</h1>
        <p class="max-w-2xl text-sm leading-6" style="color: var(--color-text-muted);">
          Acompanhe o dia atual, valide a prontidão do estabelecimento e acesse rapidamente o link público.
        </p>
      </div>

      <div class="flex flex-wrap gap-3 lg:justify-end">
        <AppButton to="/dashboard/professionals" variant="primary">Gerenciar profissionais</AppButton>
        <AppButton to="/dashboard/services" variant="secondary">Gerenciar serviços</AppButton>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-[1.5fr_1fr]">
      <AppSurface tone="default" padding="none">
        <div class="flex items-center justify-between border-b px-6 py-5" style="border-color: var(--color-border);">
          <div>
            <p class="ds-kicker">Hoje</p>
            <h2 class="mt-1 text-lg font-semibold" style="color: var(--color-text);">Agendamentos de hoje</h2>
            <p class="mt-0.5 text-xs" style="color: var(--color-text-soft);">{{ todayAppointments.length }} agendamento(s)</p>
          </div>
          <AppButton to="/dashboard/appointments" variant="ghost" size="sm">
            Ver todos
          </AppButton>
        </div>

        <div v-if="appointmentsLoading" class="p-6 text-center text-sm" style="color: var(--color-text-muted);">
          Carregando...
        </div>
        <div v-else-if="todayAppointments.length === 0" class="p-6 text-center text-sm" style="color: var(--color-text-muted);">
          Nenhum agendamento para hoje.
        </div>
        <ul v-else class="ds-grid-table divide-y" style="border-color: var(--color-border);">
          <li
            v-for="appointment in todayAppointments"
            :key="appointment.id"
            class="ds-grid-row flex items-center gap-4 px-6 py-3"
          >
            <div class="w-14 flex-shrink-0 text-center">
              <p class="text-sm font-semibold" style="color: var(--color-text);">{{ formatShortTimeBr(appointment.starts_at) }}</p>
              <p class="text-xs" style="color: var(--color-text-soft);">{{ formatShortTimeBr(appointment.ends_at) }}</p>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold" style="color: var(--color-text);">{{ appointment.client_name }}</p>
              <p class="truncate text-xs" style="color: var(--color-text-soft);">{{ appointment.service_name }} · {{ appointment.professional_name }}</p>
            </div>
            <AppStatusPill :tone="statusTone(appointment.status)" class="flex-shrink-0">
              {{ statusLabel(appointment.status) }}
            </AppStatusPill>
          </li>
        </ul>
      </AppSurface>

      <AppSurface tone="default" padding="lg">
        <p class="ds-kicker">Acesso público</p>
        <h2 class="mt-1 text-lg font-semibold" style="color: var(--color-text);">Link público</h2>
        <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
          Compartilhe este link para clientes agendarem sem acessar o painel.
        </p>

        <div
          class="mt-5 rounded-[1.2rem] border p-4"
          style="background: var(--color-surface-muted); border-color: var(--color-border);"
        >
          <p class="text-[0.68rem] font-semibold uppercase tracking-[0.28em]" style="color: var(--color-text-soft);">URL</p>
          <p class="mt-1.5 break-all text-sm font-semibold" style="color: var(--color-brand-primary);">{{ publicLink }}</p>
        </div>

        <div class="mt-4 flex flex-wrap gap-3">
          <AppButton :to="`/page/${slug}`" variant="primary" size="lg">Abrir página pública</AppButton>
        </div>
      </AppSurface>
    </div>

    <AppSurface tone="default" padding="lg">
      <div class="flex items-center justify-between gap-4">
        <div>
          <p class="ds-kicker">Checklist</p>
          <h2 class="mt-1 text-lg font-semibold" style="color: var(--color-text);">Status do estabelecimento</h2>
          <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
            O gestor só está pronto para operar quando todas as configurações obrigatórias abaixo estiverem concluídas.
          </p>
        </div>
        <AppButton to="/dashboard/settings" variant="ghost" size="sm">
          Configurações
        </AppButton>
      </div>

      <div
        class="mt-5 flex flex-col gap-3 rounded-[1.2rem] border px-4 py-4 md:flex-row md:items-center md:justify-between"
        :style="readinessSummaryStyle"
      >
        <div>
          <p class="text-sm font-semibold" style="color: var(--color-text);">
            {{ isReadyToOperate ? 'Estabelecimento pronto para uso' : 'Existem pendências bloqueando a operação' }}
          </p>
          <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
            {{ readinessSummaryText }}
          </p>
        </div>
        <AppStatusPill :tone="isReadyToOperate ? 'success' : 'warning'">
          {{ completedChecklistCount }}/{{ onboardingChecklist.length }} concluídos
        </AppStatusPill>
      </div>

      <ul class="mt-5 space-y-4">
        <li
          v-for="item in onboardingChecklist"
          :key="item.title"
          class="rounded-[1.2rem] border px-4 py-4"
          style="border-color: var(--color-border); background: var(--color-surface-muted);"
        >
          <div class="flex items-start justify-between gap-4 md:items-center">
            <div class="flex items-start gap-3">
              <span
                class="mt-1.5 h-2.5 w-2.5 flex-shrink-0 rounded-full"
                :style="item.done ? 'background: var(--color-success)' : 'background: var(--color-warning)'"
              />
              <div>
                <div class="flex flex-wrap items-center gap-2">
                  <p class="text-sm font-semibold" style="color: var(--color-text);">{{ item.title }}</p>
                  <AppStatusPill :tone="item.done ? 'success' : 'warning'">
                    {{ item.done ? 'Concluído' : 'Pendente' }}
                  </AppStatusPill>
                </div>
              </div>
            </div>

            <div class="flex shrink-0 items-center gap-2 self-start md:self-center">
              <AppButton :to="item.to" variant="ghost" size="sm" class="flex-shrink-0">
                {{ item.actionLabel }}
              </AppButton>
              <AppButton
                v-if="item.done"
                variant="secondary"
                size="sm"
                class="flex h-9 w-9 flex-shrink-0 items-center justify-center p-0"
                :aria-label="isChecklistItemExpanded(item.title) ? 'Recolher detalhes do checklist' : 'Expandir detalhes do checklist'"
                @click="toggleChecklistItem(item.title)"
              >
                <svg
                  class="h-4 w-4"
                  viewBox="0 0 20 20"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                  aria-hidden="true"
                >
                  <path
                    v-if="isChecklistItemExpanded(item.title)"
                    d="M5 12.5L10 7.5L15 12.5"
                    stroke="currentColor"
                    stroke-width="1.8"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                  <path
                    v-else
                    d="M5 7.5L10 12.5L15 7.5"
                    stroke="currentColor"
                    stroke-width="1.8"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </AppButton>
            </div>
          </div>

          <div v-if="!item.done || isChecklistItemExpanded(item.title)" class="ml-5 mt-3 space-y-1">
            <p class="text-sm" style="color: var(--color-text-muted);">{{ item.description }}</p>
            <p class="text-sm" style="color: var(--color-text-soft);">{{ item.impact }}</p>
          </div>
        </li>
      </ul>
    </AppSurface>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'manager', middleware: 'auth' })

const { establishment, fetchEstablishment } = useEstablishment()
const { professionals, fetchProfessionals } = useProfessionals()
const { services, fetchServices } = useServices()
const { hours, fetchHours } = useBusinessHours()
const { appointments: todayAppointments, loading: appointmentsLoading, fetchAppointments } = useAppointments()
const expandedChecklistItems = ref<string[]>([])

const today = new Date().toLocaleDateString('pt-BR', {
  weekday: 'long',
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

const establishmentName = computed(() => establishment.value?.name || 'Seu estabelecimento')
const slug = computed(() => establishment.value?.slug || 'seu-slug')
const publicLink = computed(() => `/page/${slug.value}`)
const openDays = computed(() => hours.value.filter(hour => !hour.is_closed).length)
const googleConnected = computed(() => Boolean(establishment.value?.google_calendar_connected))

const onboardingChecklist = computed(() => [
  {
    title: 'Google Agenda conectado',
    description: 'Conecte a conta Google para habilitar a sincronização das agendas dos profissionais.',
    impact: 'Sem essa conexão, o estabelecimento ainda não deve operar na plataforma.',
    done: googleConnected.value,
    to: '/dashboard/settings/google',
    actionLabel: 'Gerenciar',
  },
  {
    title: 'Profissionais cadastrados',
    description: 'Cadastre pelo menos um profissional para formar a agenda de atendimento.',
    impact: 'Sem profissionais, não há agenda, disponibilidade nem vínculo de calendário.',
    done: professionals.value.length > 0,
    to: '/dashboard/professionals',
    actionLabel: 'Gerenciar',
  },
  {
    title: 'Serviços publicados',
    description: 'Defina os serviços que o cliente poderá reservar pelo link público.',
    impact: 'Sem serviços, o fluxo público não consegue iniciar um novo agendamento.',
    done: services.value.length > 0,
    to: '/dashboard/services',
    actionLabel: 'Gerenciar',
  },
  {
    title: 'Horários de atendimento',
    description: 'Configure os dias e janelas de atendimento para liberar horários reais.',
    impact: 'Sem horários, o motor de disponibilidade não oferece slots para reserva.',
    done: openDays.value > 0,
    to: '/dashboard/hours',
    actionLabel: 'Gerenciar',
  },
])

const completedChecklistCount = computed(() => onboardingChecklist.value.filter(item => item.done).length)
const isReadyToOperate = computed(() => onboardingChecklist.value.every(item => item.done))
const pendingChecklistTitles = computed(() =>
  onboardingChecklist.value.filter(item => !item.done).map(item => item.title.toLowerCase()),
)
const readinessSummaryText = computed(() => {
  if (isReadyToOperate.value) {
    return 'O fluxo operacional e o agendamento público já podem ser usados com segurança.'
  }
  return `Conclua ${pendingChecklistTitles.value.join(', ')} antes de considerar o estabelecimento pronto para uso.`
})
const readinessSummaryStyle = computed(() =>
  isReadyToOperate.value
    ? 'border-color: rgba(20, 125, 100, 0.2); background: var(--color-success-soft);'
    : 'border-color: rgba(184, 107, 22, 0.2); background: var(--color-warning-soft);',
)

function isChecklistItemExpanded(title: string) {
  return expandedChecklistItems.value.includes(title)
}

function toggleChecklistItem(title: string) {
  if (isChecklistItemExpanded(title)) {
    expandedChecklistItems.value = expandedChecklistItems.value.filter(item => item !== title)
    return
  }

  expandedChecklistItems.value = [...expandedChecklistItems.value, title]
}

function todayStr() {
  return formatLocalDateInputValue()
}

function statusLabel(status: string) {
  const map: Record<string, string> = {
    confirmed: 'Confirmado',
    completed: 'Concluído',
    no_show: 'No-show',
    cancelled: 'Cancelado',
  }
  return map[status] ?? status
}

function statusTone(status: string): 'info' | 'success' | 'warning' | 'danger' | 'neutral' {
  return ({
    confirmed: 'info',
    completed: 'success',
    no_show: 'warning',
    cancelled: 'danger',
  } as const)[status] ?? 'neutral'
}

onMounted(async () => {
  await Promise.allSettled([
    fetchEstablishment(),
    fetchProfessionals(),
    fetchServices(),
    fetchHours(),
    fetchAppointments({ date_from: todayStr(), date_to: todayStr(), per_page: 50 }),
  ])
})
</script>

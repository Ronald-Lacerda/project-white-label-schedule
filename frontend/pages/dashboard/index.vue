<template>
  <div class="space-y-6">
    <AppSurface tone="default" padding="lg">
      <div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
        <div>
          <p class="ds-kicker">Painel do estabelecimento</p>
          <h1 class="ds-title mt-1">{{ establishmentName }}</h1>
          <p class="mt-2 text-sm" style="color: var(--color-text-soft);">{{ today }}</p>
        </div>
        <div class="flex flex-wrap gap-3">
          <AppButton to="/dashboard/professionals" variant="primary">Gerenciar profissionais</AppButton>
          <AppButton to="/dashboard/services" variant="secondary">Gerenciar servicos</AppButton>
        </div>
      </div>
    </AppSurface>

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
              <p class="text-sm font-semibold" style="color: var(--color-text);">{{ formatTime(appointment.starts_at) }}</p>
              <p class="text-xs" style="color: var(--color-text-soft);">{{ formatTime(appointment.ends_at) }}</p>
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
        <p class="ds-kicker">Acesso publico</p>
        <h2 class="mt-1 text-lg font-semibold" style="color: var(--color-text);">Link publico</h2>
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
          <AppButton :to="`/p/${slug}`" variant="primary" size="lg">Abrir pagina publica</AppButton>
          <AppButton to="/dashboard/hours" variant="secondary" size="lg">Ajustar horarios</AppButton>
        </div>
      </AppSurface>
    </div>

    <AppSurface tone="default" padding="lg">
      <div class="flex items-center justify-between gap-4">
        <div>
          <p class="ds-kicker">Checklist</p>
          <h2 class="mt-1 text-lg font-semibold" style="color: var(--color-text);">Status do estabelecimento</h2>
          <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
            O gestor so esta pronto para operar quando todas as configuracoes obrigatorias abaixo estiverem concluidas.
          </p>
        </div>
        <AppButton to="/dashboard/settings" variant="ghost" size="sm">
          Configuracoes
        </AppButton>
      </div>

      <div
        class="mt-5 flex flex-col gap-3 rounded-[1.2rem] border px-4 py-4 md:flex-row md:items-center md:justify-between"
        :style="readinessSummaryStyle"
      >
        <div>
          <p class="text-sm font-semibold" style="color: var(--color-text);">
            {{ isReadyToOperate ? 'Estabelecimento pronto para uso' : 'Existem pendencias bloqueando a operacao' }}
          </p>
          <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
            {{ readinessSummaryText }}
          </p>
        </div>
        <AppStatusPill :tone="isReadyToOperate ? 'success' : 'warning'">
          {{ completedChecklistCount }}/{{ onboardingChecklist.length }} concluidos
        </AppStatusPill>
      </div>

      <ul class="mt-5 space-y-4">
        <li
          v-for="item in onboardingChecklist"
          :key="item.title"
          class="flex items-start justify-between gap-4 rounded-[1.2rem] border px-4 py-4"
          style="border-color: var(--color-border); background: var(--color-surface-muted);"
        >
          <div class="flex items-start gap-3">
            <span
              class="mt-1.5 h-2.5 w-2.5 flex-shrink-0 rounded-full"
              :style="item.done ? 'background: var(--color-success)' : 'background: var(--color-warning)'"
            />
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <p class="text-sm font-semibold" style="color: var(--color-text);">{{ item.title }}</p>
                <AppStatusPill :tone="item.done ? 'success' : 'warning'">
                  {{ item.done ? 'Concluido' : 'Pendente' }}
                </AppStatusPill>
              </div>
              <p class="mt-1 text-sm" style="color: var(--color-text-muted);">{{ item.description }}</p>
              <p class="mt-1 text-sm" style="color: var(--color-text-soft);">{{ item.impact }}</p>
            </div>
          </div>

          <AppButton :to="item.to" variant="ghost" size="sm" class="flex-shrink-0">
            {{ item.actionLabel }}
          </AppButton>
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

const today = new Date().toLocaleDateString('pt-BR', {
  weekday: 'long',
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

const establishmentName = computed(() => establishment.value?.name || 'Seu estabelecimento')
const slug = computed(() => establishment.value?.slug || 'seu-slug')
const publicLink = computed(() => `/p/${slug.value}`)
const openDays = computed(() => hours.value.filter(hour => !hour.is_closed).length)
const googleConnected = computed(() => Boolean(establishment.value?.google_calendar_connected))

const onboardingChecklist = computed(() => [
  {
    title: 'Profissionais cadastrados',
    description: 'Cadastre pelo menos um profissional para formar a agenda de atendimento.',
    impact: 'Sem profissionais, nao ha agenda, disponibilidade nem vinculo de calendario.',
    done: professionals.value.length > 0,
    to: '/dashboard/professionals',
    actionLabel: professionals.value.length > 0 ? 'Revisar' : 'Configurar',
  },
  {
    title: 'Servicos publicados',
    description: 'Defina os servicos que o cliente podera reservar pelo link publico.',
    impact: 'Sem servicos, o fluxo publico nao consegue iniciar um novo agendamento.',
    done: services.value.length > 0,
    to: '/dashboard/services',
    actionLabel: services.value.length > 0 ? 'Revisar' : 'Configurar',
  },
  {
    title: 'Horarios de atendimento',
    description: 'Configure os dias e janelas de atendimento para liberar horarios reais.',
    impact: 'Sem horarios, o motor de disponibilidade nao oferece slots para reserva.',
    done: openDays.value > 0,
    to: '/dashboard/hours',
    actionLabel: openDays.value > 0 ? 'Ajustar' : 'Configurar',
  },
  {
    title: 'Google Agenda conectado',
    description: 'Conecte a conta Google para habilitar a sincronizacao das agendas dos profissionais.',
    impact: 'Sem essa conexao, o estabelecimento ainda nao deve operar na plataforma.',
    done: googleConnected.value,
    to: '/dashboard/settings/google',
    actionLabel: googleConnected.value ? 'Gerenciar' : 'Conectar',
  },
])

const completedChecklistCount = computed(() => onboardingChecklist.value.filter(item => item.done).length)
const isReadyToOperate = computed(() => onboardingChecklist.value.every(item => item.done))
const pendingChecklistTitles = computed(() =>
  onboardingChecklist.value.filter(item => !item.done).map(item => item.title.toLowerCase()),
)
const readinessSummaryText = computed(() => {
  if (isReadyToOperate.value) {
    return 'O fluxo operacional e o agendamento publico ja podem ser usados com seguranca.'
  }
  return `Conclua ${pendingChecklistTitles.value.join(', ')} antes de considerar o estabelecimento pronto para uso.`
})
const readinessSummaryStyle = computed(() =>
  isReadyToOperate.value
    ? 'border-color: rgba(20, 125, 100, 0.2); background: var(--color-success-soft);'
    : 'border-color: rgba(184, 107, 22, 0.2); background: var(--color-warning-soft);',
)

function todayStr() {
  return formatLocalDateInputValue()
}

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
}

function statusLabel(status: string) {
  const map: Record<string, string> = {
    confirmed: 'Confirmado',
    completed: 'Concluido',
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
    fetchAppointments({ date: todayStr(), per_page: 50 }),
  ])
})
</script>

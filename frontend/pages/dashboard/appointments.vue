<template>
  <div class="ds-page">
    <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
      <div class="space-y-2">
        <h1 class="ds-title">Agendamentos</h1>
        <p class="max-w-2xl text-sm leading-6" style="color: var(--color-text-muted);">
          Consulte a agenda por período, acompanhe os status e bloqueie intervalos da equipe quando necessário.
        </p>
      </div>

      <div class="flex flex-wrap gap-3 lg:justify-end">
        <AppButton variant="secondary" @click="clearFilters">
          Limpar filtros
        </AppButton>
        <AppButton variant="primary" @click="openBlockModal">
          Bloquear período
        </AppButton>
      </div>
    </div>

    <AppSurface tone="brand">
      <div class="flex flex-col gap-5">
        <div class="flex flex-col gap-2 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p class="ds-kicker">Filtro persistente</p>
            <h2 class="mt-2 text-xl font-semibold" style="color: var(--color-text);">Refine a agenda sem perder o panorama</h2>
          </div>
          <div class="flex flex-wrap gap-2">
            <AppStatusPill tone="info">Confirmado</AppStatusPill>
            <AppStatusPill tone="danger">Cancelado</AppStatusPill>
          </div>
        </div>

        <div class="grid gap-4 lg:grid-cols-[1fr_1fr_1fr_1fr_auto]">
          <div>
            <label class="ds-label">De</label>
            <input v-model="filters.date_from" type="date" lang="pt-BR" class="ds-input" />
          </div>

          <div>
            <label class="ds-label">Até</label>
            <input v-model="filters.date_to" type="date" lang="pt-BR" class="ds-input" />
          </div>

          <div>
            <label class="ds-label">Profissional</label>
            <select v-model="filters.professional_id" class="ds-select">
              <option value="">Todos</option>
              <option v-for="professional in professionals" :key="professional.id" :value="professional.id">
                {{ professional.name }}
              </option>
            </select>
          </div>

          <div>
            <label class="ds-label">Status</label>
            <select v-model="filters.status" class="ds-select">
              <option value="">Todos</option>
              <option value="confirmed">Confirmado</option>
              <option value="cancelled">Cancelado</option>
            </select>
          </div>

          <div class="flex items-end">
            <AppButton variant="secondary" block @click="applyFilters">
              Aplicar
            </AppButton>
          </div>
        </div>

        <div v-if="blockedPeriods.length > 0" class="rounded-[1.4rem] border p-4" style="border-color: rgba(var(--color-brand-primary-rgb), 0.12); background: rgba(255, 255, 255, 0.72);">
          <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p class="text-sm font-semibold" style="color: var(--color-text);">Bloqueios ativos no filtro atual</p>
              <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
                Use esta área para checar indisponibilidades antes de editar ou encaixar horários.
              </p>
            </div>
            <AppStatusPill tone="warning">{{ blockedPeriods.length }} bloqueio(s)</AppStatusPill>
          </div>

          <div class="mt-4 flex flex-wrap gap-3">
            <div
              v-for="blocked in blockedPeriods"
              :key="blocked.id"
              class="flex min-w-[260px] flex-1 items-start justify-between gap-3 rounded-[1.2rem] border p-4"
              style="border-color: var(--color-border); background: var(--color-surface);"
            >
              <div>
                <p class="text-sm font-semibold" style="color: var(--color-text);">{{ blocked.professional_name }}</p>
                <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
                  {{ formatDate(blocked.starts_at) }} · {{ formatTime(blocked.starts_at) }} até {{ formatTime(blocked.ends_at) }}
                </p>
                <p v-if="blocked.reason" class="mt-2 text-xs font-medium uppercase tracking-[0.18em]" style="color: var(--color-text-soft);">
                  {{ blocked.reason }}
                </p>
              </div>
              <AppButton variant="ghost" size="sm" @click="promptRemoveBlockedPeriod(blocked.id, blocked.professional_name)">
                Remover
              </AppButton>
            </div>
          </div>
        </div>
      </div>
    </AppSurface>

    <AppSurface tone="default" padding="none">
      <div v-if="loading" class="ds-empty-state">
        <AppStatusPill tone="info">Carregando</AppStatusPill>
        <p class="text-sm">Buscando agendamentos e bloqueios do filtro atual.</p>
      </div>

      <div v-else-if="error" class="ds-empty-state">
        <AppStatusPill tone="danger">Erro</AppStatusPill>
        <p class="text-sm">{{ error }}</p>
        <AppButton variant="secondary" @click="applyFilters">Tentar novamente</AppButton>
      </div>

      <div v-else-if="appointments.length === 0" class="ds-empty-state">
        <AppStatusPill tone="neutral">Sem resultados</AppStatusPill>
        <p class="text-sm">Nenhum agendamento encontrado para os filtros selecionados.</p>
      </div>

      <div v-else>
        <div class="hidden overflow-hidden xl:block">
          <table class="min-w-full">
            <thead>
              <tr class="border-b" style="border-color: var(--color-border); background: var(--color-surface-muted);">
                <th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-[0.22em]" style="color: var(--color-text-soft);">Horário</th>
                <th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-[0.22em]" style="color: var(--color-text-soft);">Cliente</th>
                <th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-[0.22em]" style="color: var(--color-text-soft);">Serviço</th>
                <th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-[0.22em]" style="color: var(--color-text-soft);">Profissional</th>
                <th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-[0.22em]" style="color: var(--color-text-soft);">Status</th>
                <th class="px-6 py-4 text-right text-xs font-semibold uppercase tracking-[0.22em]" style="color: var(--color-text-soft);">Ações</th>
              </tr>
            </thead>
            <tbody class="ds-grid-table">
              <tr
                v-for="appointment in appointments"
                :key="appointment.id"
                class="ds-grid-row cursor-pointer"
                @click="openDetail(appointment)"
              >
                <td class="px-6 py-4">
                  <p class="text-sm font-semibold" style="color: var(--color-text);">{{ formatTime(appointment.starts_at) }}</p>
                  <p class="mt-1 text-xs" style="color: var(--color-text-soft);">{{ formatDate(appointment.starts_at) }}</p>
                </td>
                <td class="px-6 py-4">
                  <p class="text-sm font-semibold" style="color: var(--color-text);">{{ appointment.client_name }}</p>
                  <p class="mt-1 text-xs" style="color: var(--color-text-soft);">{{ formatBrazilianPhoneDisplay(appointment.client_phone) }}</p>
                </td>
                <td class="px-6 py-4">
                  <p class="text-sm font-medium" style="color: var(--color-text);">{{ appointment.service_name }}</p>
                  <p class="mt-1 text-xs" style="color: var(--color-text-soft);">{{ appointment.duration_minutes }} min</p>
                </td>
                <td class="px-6 py-4">
                  <p class="text-sm font-medium" style="color: var(--color-text);">{{ appointment.professional_name }}</p>
                </td>
                <td class="px-6 py-4">
                  <AppStatusPill :tone="statusTone(appointment.status)">{{ statusLabel(appointment.status) }}</AppStatusPill>
                </td>
                <td class="px-6 py-4 align-middle" @click.stop>
                  <div v-if="appointment.status === 'confirmed'" class="flex items-center justify-end gap-2">
                    <AppButton size="sm" variant="danger" @click="changeStatus(appointment, 'cancelled')">Cancelar</AppButton>
                  </div>
                  <div v-else class="flex items-center justify-end">
                    <span class="text-xs font-medium uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Finalizado</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="grid gap-4 p-4 xl:hidden md:grid-cols-2">
          <div
            v-for="appointment in appointments"
            :key="`card-${appointment.id}`"
            class="rounded-[1.5rem] border p-4"
            style="border-color: var(--color-border); background: var(--color-surface);"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="text-xs font-semibold uppercase tracking-[0.24em]" style="color: var(--color-text-soft);">
                  {{ formatDate(appointment.starts_at) }}
                </p>
                <p class="mt-2 text-lg font-semibold" style="color: var(--color-text);">{{ formatTime(appointment.starts_at) }}</p>
              </div>
              <AppStatusPill :tone="statusTone(appointment.status)">{{ statusLabel(appointment.status) }}</AppStatusPill>
            </div>

            <div class="mt-4 space-y-3">
              <div>
                <p class="text-sm font-semibold" style="color: var(--color-text);">{{ appointment.client_name }}</p>
                <p class="text-sm" style="color: var(--color-text-muted);">{{ formatBrazilianPhoneDisplay(appointment.client_phone) }}</p>
              </div>
              <div class="grid gap-3 sm:grid-cols-2">
                <div>
                  <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Serviço</p>
                  <p class="mt-1 text-sm font-medium" style="color: var(--color-text);">{{ appointment.service_name }}</p>
                </div>
                <div>
                  <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Profissional</p>
                  <p class="mt-1 text-sm font-medium" style="color: var(--color-text);">{{ appointment.professional_name }}</p>
                </div>
              </div>
            </div>

            <div class="mt-4 flex flex-wrap gap-2">
              <AppButton size="sm" variant="secondary" @click="openDetail(appointment)">Detalhes</AppButton>
              <AppButton v-if="appointment.status === 'confirmed'" size="sm" variant="danger" @click="changeStatus(appointment, 'cancelled')">Cancelar</AppButton>
            </div>
          </div>
        </div>

        <div v-if="meta.total > meta.per_page" class="flex flex-col gap-4 border-t px-5 py-4 md:flex-row md:items-center md:justify-between" style="border-color: var(--color-border);">
          <p class="text-sm" style="color: var(--color-text-muted);">
            {{ (meta.page - 1) * meta.per_page + 1 }}-{{ Math.min(meta.page * meta.per_page, meta.total) }} de {{ meta.total }} agendamentos
          </p>
          <div class="flex gap-2">
            <AppButton size="sm" variant="secondary" :disabled="meta.page <= 1" @click="goPage(meta.page - 1)">
              Anterior
            </AppButton>
            <AppButton size="sm" variant="secondary" :disabled="meta.page * meta.per_page >= meta.total" @click="goPage(meta.page + 1)">
              Próxima
            </AppButton>
          </div>
        </div>
      </div>
    </AppSurface>

    <AppModal
      :open="Boolean(detailAppt)"
      title="Detalhe do agendamento"
      description="Consulte dados do cliente, status atual e realize ações rápidas sem sair da agenda."
      eyebrow="Consulta"
      width="md"
      @close="detailAppt = null"
    >
      <div v-if="detailAppt" class="space-y-6">
        <div class="grid gap-4 md:grid-cols-2">
          <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
            <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Cliente</p>
            <p class="mt-2 text-sm font-semibold" style="color: var(--color-text);">{{ detailAppt.client_name }}</p>
            <p class="mt-1 text-sm" style="color: var(--color-text-muted);">{{ formatBrazilianPhoneDisplay(detailAppt.client_phone) }}</p>
            <p v-if="detailAppt.client_email" class="mt-1 text-sm" style="color: var(--color-text-muted);">{{ detailAppt.client_email }}</p>
            <p v-if="detailAppt.client_birth_date" class="mt-1 text-sm" style="color: var(--color-text-muted);">Nascimento: {{ formatBirthDate(detailAppt.client_birth_date) }}</p>
          </div>
          <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
            <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Status</p>
            <div class="mt-2">
              <AppStatusPill :tone="statusTone(detailAppt.status)">{{ statusLabel(detailAppt.status) }}</AppStatusPill>
            </div>
          </div>
          <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
            <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Serviço</p>
            <p class="mt-2 text-sm font-semibold" style="color: var(--color-text);">
              {{ detailAppt.service_name }} ({{ detailAppt.duration_minutes }} min)
            </p>
          </div>
          <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
            <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Profissional</p>
            <p class="mt-2 text-sm font-semibold" style="color: var(--color-text);">{{ detailAppt.professional_name }}</p>
          </div>
        </div>

        <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface);">
          <p class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Data e hora</p>
          <p class="mt-2 text-sm font-semibold" style="color: var(--color-text);">
            {{ formatDate(detailAppt.starts_at) }} · {{ formatTime(detailAppt.starts_at) }} até {{ formatTime(detailAppt.ends_at) }}
          </p>
          <p v-if="detailAppt.notes" class="mt-4 text-sm" style="color: var(--color-text-muted);">
            {{ detailAppt.notes }}
          </p>
        </div>
      </div>

      <template #footer>
        <div v-if="detailAppt?.status === 'confirmed'" class="flex flex-wrap justify-end gap-2">
          <AppButton variant="danger" @click="changeStatus(detailAppt, 'cancelled')">Cancelar agendamento</AppButton>
        </div>
      </template>
    </AppModal>

    <AppModal
      :open="showBlockModal"
      title="Bloquear período"
      description="Reserve um intervalo para pausa, reunião, férias ou qualquer indisponibilidade da equipe."
      eyebrow="Disponibilidade"
      width="sm"
      @close="showBlockModal = false"
    >
      <div class="space-y-4">
        <div>
          <label class="ds-label">Profissional</label>
          <select v-model="blockForm.professional_id" class="ds-select">
            <option value="">Selecione...</option>
            <option v-for="professional in professionals" :key="professional.id" :value="professional.id">{{ professional.name }}</option>
          </select>
        </div>
        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="ds-label">Início</label>
            <input v-model="blockForm.starts_at" type="datetime-local" lang="pt-BR" class="ds-input" />
          </div>
          <div>
            <label class="ds-label">Fim</label>
            <input v-model="blockForm.ends_at" type="datetime-local" lang="pt-BR" class="ds-input" />
          </div>
        </div>
        <div>
          <label class="ds-label">Motivo</label>
          <input v-model="blockForm.reason" type="text" class="ds-input" placeholder="Ex.: reunião interna, férias, manutenção..." />
        </div>
        <p v-if="blockError" class="text-sm font-medium" style="color: var(--color-danger);">{{ blockError }}</p>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <AppButton variant="ghost" @click="showBlockModal = false">Cancelar</AppButton>
          <AppButton variant="primary" :disabled="blockSaving" @click="submitBlock">
            {{ blockSaving ? 'Salvando...' : 'Salvar bloqueio' }}
          </AppButton>
        </div>
      </template>
    </AppModal>

    <AppConfirmModal
      :open="removeBlockModal.open"
      title="Remover bloqueio"
      description="Essa ação libera novamente o período para o motor de disponibilidade."
      eyebrow="Disponibilidade"
      message="Deseja remover este bloqueio?"
      :details="removeBlockModal.professional ? `Profissional: ${removeBlockModal.professional}` : ''"
      confirm-label="Remover bloqueio"
      loading-label="Removendo..."
      :loading="removeBlockModal.loading"
      @cancel="closeRemoveBlockedPeriodModal"
      @confirm="removeBlockedPeriod"
    />
  </div>
</template>

<script setup lang="ts">
import type { ManagerAppointment } from '~/composables/useAppointments'

definePageMeta({ layout: 'manager', middleware: 'auth' })

const {
  appointments,
  meta,
  blockedPeriods,
  loading,
  error,
  fetchAppointments,
  updateStatus,
  fetchBlockedPeriods,
  createBlockedPeriod,
  deleteBlockedPeriod,
} = useAppointments()
const { professionals, fetchProfessionals } = useProfessionals()

const filters = reactive({
  date_from: formatLocalDateInputValue(),
  date_to: formatLocalDateInputValue(),
  professional_id: '',
  status: '',
})

const detailAppt = ref<ManagerAppointment | null>(null)
const showBlockModal = ref(false)
const blockForm = reactive({ professional_id: '', starts_at: '', ends_at: '', reason: '' })
const blockError = ref('')
const blockSaving = ref(false)
const removeBlockModal = reactive({
  open: false,
  id: '',
  professional: '',
  loading: false,
})

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

function formatTime(iso: string) {
  return formatShortTimeBr(iso)
}

function formatBirthDate(value: string) {
  return new Date(`${value}T00:00:00`).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
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

function statusTone(status: string) {
  return {
    confirmed: 'info',
    completed: 'success',
    no_show: 'warning',
    cancelled: 'danger',
  }[status] as 'info' | 'success' | 'warning' | 'danger'
}

async function applyFilters() {
  await Promise.allSettled([
    fetchAppointments({ ...filters, page: 1, per_page: 20 }),
    fetchBlockedPeriods(filters.professional_id || undefined, filters.date_from || undefined, filters.date_to || undefined),
  ])
}

async function clearFilters() {
  filters.date_from = formatLocalDateInputValue()
  filters.date_to = formatLocalDateInputValue()
  filters.professional_id = ''
  filters.status = ''
  await applyFilters()
}

async function goPage(page: number) {
  await fetchAppointments({ ...filters, page, per_page: meta.value.per_page })
}

function openDetail(appointment: ManagerAppointment) {
  detailAppt.value = appointment
}

async function changeStatus(appointment: ManagerAppointment, status: string) {
  await updateStatus(appointment.id, status)
  if (detailAppt.value?.id === appointment.id) {
    detailAppt.value = { ...detailAppt.value, status: status as ManagerAppointment['status'] }
  }
}

function openBlockModal() {
  const now = new Date()
  const end = new Date(now.getTime() + 60 * 60 * 1000)
  blockForm.professional_id = filters.professional_id || ''
  blockForm.starts_at = formatLocalDateTimeInputValue(now)
  blockForm.ends_at = formatLocalDateTimeInputValue(end)
  blockForm.reason = ''
  blockError.value = ''
  showBlockModal.value = true
}

async function submitBlock() {
  if (!blockForm.professional_id || !blockForm.starts_at || !blockForm.ends_at) {
    blockError.value = 'Preencha profissional, início e fim.'
    return
  }

  blockError.value = ''
  blockSaving.value = true
  try {
    await createBlockedPeriod({
      professional_id: blockForm.professional_id,
      starts_at: localDateTimeToIso(blockForm.starts_at),
      ends_at: localDateTimeToIso(blockForm.ends_at),
      reason: blockForm.reason || undefined,
    })
    showBlockModal.value = false
    await fetchBlockedPeriods(filters.professional_id || undefined, filters.date_from || undefined, filters.date_to || undefined)
  } catch (e: any) {
    blockError.value = e?.data?.error?.message ?? 'Erro ao criar bloqueio.'
  } finally {
    blockSaving.value = false
  }
}

function promptRemoveBlockedPeriod(id: string, professional: string) {
  removeBlockModal.open = true
  removeBlockModal.id = id
  removeBlockModal.professional = professional
}

function closeRemoveBlockedPeriodModal() {
  removeBlockModal.open = false
  removeBlockModal.id = ''
  removeBlockModal.professional = ''
  removeBlockModal.loading = false
}

async function removeBlockedPeriod() {
  if (!removeBlockModal.id) return
  removeBlockModal.loading = true
  try {
    await deleteBlockedPeriod(removeBlockModal.id)
    closeRemoveBlockedPeriodModal()
  } finally {
    removeBlockModal.loading = false
  }
}

onMounted(async () => {
  await Promise.allSettled([
    fetchProfessionals(),
    fetchAppointments({ date_from: filters.date_from, date_to: filters.date_to, page: 1, per_page: 20 }),
    fetchBlockedPeriods(undefined, filters.date_from, filters.date_to),
  ])
})
</script>

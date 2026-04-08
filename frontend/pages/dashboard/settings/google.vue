<template>
  <div class="max-w-2xl space-y-6">
    <div class="flex items-center gap-2">
      <AppButton to="/dashboard/settings" variant="ghost" size="sm">
        Voltar para configuracoes
      </AppButton>
    </div>

    <div
      v-if="callbackStatus === 'success'"
      class="rounded-[1.2rem] border px-5 py-4 text-sm font-medium"
      style="background: var(--color-success-soft); border-color: var(--color-success); color: var(--color-success);"
    >
      Google Agenda conectado com sucesso!
    </div>
    <div
      v-if="callbackStatus === 'error'"
      class="rounded-[1.2rem] border px-5 py-4 text-sm font-medium"
      style="background: var(--color-danger-soft); border-color: var(--color-danger); color: var(--color-danger);"
    >
      Nao foi possivel conectar o Google Agenda. Tente novamente.
    </div>

    <div v-if="loading" class="text-sm" style="color: var(--color-text-muted);">Carregando...</div>

    <div v-else class="space-y-6">
      <AppSurface tone="default" padding="lg">
        <div class="flex items-start justify-between gap-4">
          <div class="flex items-start gap-4">
            <div
              class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-[1rem] border"
              style="background: #f0f4ff; border-color: rgba(66, 133, 244, 0.2);"
            >
              <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="3" y="3" width="18" height="18" rx="2" stroke="#4285F4" stroke-width="1.5" />
                <path d="M3 9h18" stroke="#4285F4" stroke-width="1.5" />
                <path d="M9 3v6M15 3v6" stroke="#4285F4" stroke-width="1.5" stroke-linecap="round" />
                <rect x="7" y="13" width="4" height="4" rx="0.5" fill="#EA4335" />
              </svg>
            </div>

            <div>
              <h2 class="text-base font-semibold" style="color: var(--color-text);">Google Agenda</h2>
              <p class="mt-1 text-sm leading-6" style="color: var(--color-text-muted);">
                Sincronize os agendamentos com as agendas individuais de cada profissional.
              </p>
            </div>
          </div>

          <AppStatusPill :tone="status?.connected ? 'success' : 'neutral'" class="flex-shrink-0">
            {{ status?.connected ? 'Conectado' : 'Desconectado' }}
          </AppStatusPill>
        </div>

        <div class="mt-6 flex flex-wrap gap-3 border-t pt-5" style="border-color: var(--color-border);">
          <button
            v-if="!status?.connected"
            type="button"
            class="ds-button ds-button-primary"
            :disabled="connecting"
            @click="handleConnect"
          >
            <AppGoogleIcon />
            {{ connecting ? 'Redirecionando...' : 'Conectar Google Agenda' }}
          </button>
          <button
            v-else
            type="button"
            class="ds-button ds-button-danger"
            :disabled="disconnecting"
            @click="disconnectModal = true"
          >
            {{ disconnecting ? 'Desconectando...' : 'Desconectar conta' }}
          </button>
        </div>

        <p v-if="actionError" class="mt-3 text-sm" style="color: var(--color-danger);">{{ actionError }}</p>
      </AppSurface>

      <AppSurface v-if="status?.connected && status.professionals.length > 0" tone="default" padding="lg">
        <h2 class="mb-4 text-base font-semibold" style="color: var(--color-text);">Agendas dos profissionais</h2>

        <ul class="ds-grid-table divide-y" style="border-color: var(--color-border);">
          <li
            v-for="prof in status.professionals"
            :key="prof.id"
            class="ds-grid-row flex items-center justify-between px-1 py-3.5"
          >
            <div>
              <p class="text-sm font-semibold" style="color: var(--color-text);">{{ prof.name }}</p>
              <p class="mt-0.5 text-xs" :style="prof.google_calendar_id ? 'color: var(--color-text-soft)' : 'color: var(--color-warning)'">
                {{ prof.google_calendar_id ? 'Agenda sincronizada' : 'Agenda ainda nao criada' }}
              </p>
            </div>
            <AppStatusPill :tone="prof.google_calendar_id ? 'success' : 'warning'">
              {{ prof.google_calendar_id ? 'Ativo' : 'Pendente' }}
            </AppStatusPill>
          </li>
        </ul>
      </AppSurface>

      <AppSurface v-if="!status?.connected" tone="muted" padding="lg">
        <h3 class="mb-3 text-sm font-semibold" style="color: var(--color-text);">Como funciona</h3>
        <ul class="space-y-2 text-sm" style="color: var(--color-text-muted);">
          <li class="flex items-start gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 flex-shrink-0 rounded-full" style="background: var(--color-brand-primary);" />
            Cada profissional recebe uma agenda dedicada no Google Agenda.
          </li>
          <li class="flex items-start gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 flex-shrink-0 rounded-full" style="background: var(--color-brand-primary);" />
            Agendamentos confirmados aparecem automaticamente nas agendas.
          </li>
          <li class="flex items-start gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 flex-shrink-0 rounded-full" style="background: var(--color-brand-primary);" />
            Cancelamentos removem os eventos sincronizados.
          </li>
          <li class="flex items-start gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 flex-shrink-0 rounded-full" style="background: var(--color-brand-primary);" />
            Eventos externos bloqueiam slots no motor de disponibilidade.
          </li>
        </ul>
      </AppSurface>
    </div>

    <AppConfirmModal
      :open="disconnectModal"
      title="Desconectar Google Agenda"
      description="A conta sera desconectada do painel, sem alterar os agendamentos ja existentes."
      eyebrow="Integracao"
      message="Deseja desconectar a conta do Google Agenda?"
      details="A sincronizacao deixa de ocorrer ate uma nova conexao ser realizada."
      confirm-label="Desconectar conta"
      loading-label="Desconectando..."
      :loading="disconnecting"
      @cancel="disconnectModal = false"
      @confirm="handleDisconnect"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'manager', middleware: 'auth' })

const route = useRoute()
const { status, loading, fetchStatus, getAuthUrl, disconnect } = useGoogleCalendar()

const callbackStatus = computed(() => {
  const statusQuery = route.query.status
  if (statusQuery === 'success' || statusQuery === 'error') return statusQuery
  return null
})

const connecting = ref(false)
const disconnecting = ref(false)
const disconnectModal = ref(false)
const actionError = ref<string | null>(null)

await fetchStatus()

async function handleConnect() {
  connecting.value = true
  actionError.value = null
  try {
    const url = await getAuthUrl()
    window.location.href = url
  } catch (e: any) {
    actionError.value = e?.data?.error?.message ?? 'Erro ao iniciar conexao.'
    connecting.value = false
  }
}

async function handleDisconnect() {
  disconnecting.value = true
  actionError.value = null
  try {
    await disconnect()
    await fetchStatus()
    disconnectModal.value = false
  } catch (e: any) {
    actionError.value = e?.data?.error?.message ?? 'Erro ao desconectar.'
  } finally {
    disconnecting.value = false
  }
}
</script>

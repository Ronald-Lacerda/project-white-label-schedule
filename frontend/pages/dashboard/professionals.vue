<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="ds-kicker">Equipe</p>
        <h1 class="ds-title mt-1">Profissionais</h1>
        <p class="mt-2 text-sm leading-6" style="color: var(--color-text-muted);">
          Cadastre quem atende no estabelecimento e prepare a base do agendamento público.
        </p>
      </div>
      <AppButton variant="primary" @click="openModal()">+ Novo profissional</AppButton>
    </div>

    <div v-if="loading" class="text-sm" style="color: var(--color-text-muted);">Carregando profissionais...</div>

    <div
      v-if="error"
      class="rounded-[1.2rem] border px-5 py-4 text-sm"
      style="background: var(--color-danger-soft); border-color: var(--color-danger); color: var(--color-danger);"
    >
      <p>{{ error }}</p>
      <AppButton size="sm" variant="secondary" class="mt-3" @click="reload">Tentar novamente</AppButton>
    </div>

    <AppSurface v-else-if="!loading && professionals.length === 0" tone="default" padding="lg">
      <div class="ds-empty-state">
        <p>Nenhum profissional cadastrado ainda.</p>
        <AppButton variant="secondary" @click="openModal()">Cadastrar primeiro profissional</AppButton>
      </div>
    </AppSurface>

    <AppSurface v-else-if="professionals.length > 0" tone="default" padding="none">
      <ul class="ds-grid-table divide-y" style="border-color: var(--color-border);">
        <li
          v-for="professional in professionals"
          :key="professional.id"
          class="ds-grid-row flex items-center justify-between gap-4 px-6 py-4"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-full text-sm font-semibold"
              style="background: rgba(var(--color-brand-primary-rgb), 0.08); color: var(--color-brand-primary);"
            >
              {{ initialFor(professional.name) }}
            </div>
            <div>
              <p class="text-sm font-semibold" style="color: var(--color-text);">{{ professional.name }}</p>
              <p class="text-xs" style="color: var(--color-text-soft);">{{ professional.phone || 'Sem telefone cadastrado' }}</p>
              <p class="mt-1 text-xs" :style="{ color: professional.service_ids.length ? 'var(--color-text-soft)' : 'var(--color-danger)' }">
                {{ serviceSummary(professional) }}
              </p>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <AppStatusPill :tone="professional.active ? 'success' : 'neutral'">
              {{ professional.active ? 'Ativo' : 'Inativo' }}
            </AppStatusPill>
            <AppButton size="sm" variant="ghost" @click="openModal(professional)">Editar</AppButton>
            <AppButton size="sm" variant="danger" @click="handleDelete(professional)">Remover</AppButton>
          </div>
        </li>
      </ul>
    </AppSurface>

    <AppModal
      :open="modal.open"
      :title="modal.id ? 'Editar profissional' : 'Novo profissional'"
      eyebrow="Equipe"
      width="sm"
      @close="closeModal"
    >
      <form class="space-y-4" @submit.prevent="handleSave">
        <div>
          <label class="ds-label">Nome</label>
          <input v-model="modal.name" type="text" required class="ds-input" />
        </div>
        <div>
          <label class="ds-label">Telefone</label>
          <input v-model="modal.phone" type="tel" class="ds-input" />
        </div>
        <div class="space-y-3">
          <div class="flex items-center justify-between gap-3">
            <div>
              <label class="ds-label">Servicos atendidos</label>
              <p class="text-xs leading-5" style="color: var(--color-text-muted);">
                Selecione os servicos que este profissional atende.
              </p>
            </div>
            <span v-if="servicesLoading" class="text-xs" style="color: var(--color-text-muted);">Carregando...</span>
          </div>

          <div
            v-if="servicesError"
            class="rounded-[1rem] border px-4 py-3 text-sm"
            style="background: var(--color-danger-soft); border-color: var(--color-danger); color: var(--color-danger);"
          >
            {{ servicesError }}
          </div>

          <div
            v-else-if="activeServices.length === 0"
            class="rounded-[1rem] border px-4 py-3 text-sm"
            style="border-color: var(--color-border); color: var(--color-text-muted); background: var(--color-surface-muted);"
          >
            Nenhum servico ativo foi encontrado. Cadastre um servico para liberar este profissional no agendamento publico.
          </div>

          <div v-else class="grid grid-cols-1 gap-2 sm:grid-cols-2">
            <button
              v-for="service in activeServices"
              :key="service.id"
              type="button"
              class="rounded-[1rem] border px-3 py-3 text-left transition"
              :style="serviceOptionStyle(modal.selectedServiceIds.includes(service.id))"
              @click="toggleService(service.id)"
            >
              <div class="flex items-start gap-3">
                <div
                  class="mt-0.5 flex h-4 w-4 flex-shrink-0 items-center justify-center rounded border text-[10px] font-bold"
                  :style="serviceCheckStyle(modal.selectedServiceIds.includes(service.id))"
                >
                  {{ modal.selectedServiceIds.includes(service.id) ? 'x' : '' }}
                </div>
                <div class="min-w-0">
                  <p class="truncate text-sm font-semibold">{{ service.name }}</p>
                  <p class="mt-1 text-xs leading-5" :style="serviceMetaStyle(modal.selectedServiceIds.includes(service.id))">
                    {{ serviceCardDescription(service) }}
                  </p>
                </div>
              </div>
            </button>
          </div>

          <p
            v-if="!servicesLoading && modal.selectedServiceIds.length === 0"
            class="rounded-[1rem] border px-4 py-3 text-sm"
            style="border-color: rgba(180, 83, 9, 0.2); background: rgba(245, 158, 11, 0.12); color: rgb(146, 64, 14);"
          >
            Este profissional nao aparecera no agendamento publico ate ter ao menos um servico vinculado.
          </p>
        </div>
        <p v-if="modal.error" class="text-sm" style="color: var(--color-danger);">{{ modal.error }}</p>
        <div class="flex justify-end gap-3 pt-2">
          <AppButton type="button" variant="secondary" @click="closeModal">Cancelar</AppButton>
          <AppButton type="submit" variant="primary" :disabled="modal.saving">
            {{ modal.saving ? 'Salvando...' : 'Salvar' }}
          </AppButton>
        </div>
      </form>
    </AppModal>

    <AppConfirmModal
      :open="deleteModal.open"
      title="Remover profissional"
      description="Essa ação remove o cadastro da equipe e pode impactar a operação do estabelecimento."
      eyebrow="Equipe"
      message="Deseja remover este profissional?"
      :details="deleteModal.name ? `Profissional: ${deleteModal.name}` : ''"
      confirm-label="Remover profissional"
      loading-label="Removendo..."
      :loading="deleteModal.loading"
      @cancel="closeDeleteModal"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import type { Professional } from '~/composables/useProfessionals'

definePageMeta({ layout: 'manager', middleware: 'auth' })

const {
  professionals,
  loading,
  error,
  fetchProfessionals,
  getProfessional,
  createProfessional,
  updateProfessional,
  deleteProfessional,
  updateServices,
} = useProfessionals()

const {
  activeServices,
  loading: servicesLoading,
  error: servicesError,
  fetchServices,
  formatDuration,
  formatPrice,
} = useServices()

const modal = reactive({
  open: false,
  id: '',
  name: '',
  phone: '',
  selectedServiceIds: [] as string[],
  saving: false,
  error: '',
})

const deleteModal = reactive({
  open: false,
  id: '',
  name: '',
  loading: false,
})

onMounted(() => {
  reload()
})

function initialFor(name: string): string {
  return name.trim().charAt(0).toUpperCase() || '?'
}

async function reload() {
  await Promise.all([fetchProfessionals(), fetchServices()])
}

async function openModal(professional?: Professional) {
  modal.open = true
  modal.error = ''
  modal.id = professional?.id ?? ''
  modal.name = professional?.name ?? ''
  modal.phone = professional?.phone ?? ''
  modal.selectedServiceIds = professional?.service_ids ? [...professional.service_ids] : []

  if (professional?.id) {
    try {
      const detailed = await getProfessional(professional.id)
      modal.name = detailed.name
      modal.phone = detailed.phone ?? ''
      modal.selectedServiceIds = [...detailed.service_ids]
    } catch (e: any) {
      modal.error = e?.data?.error?.message ?? 'Erro ao carregar os servicos do profissional.'
    }
  }
}

function closeModal() {
  modal.open = false
  modal.id = ''
  modal.name = ''
  modal.phone = ''
  modal.selectedServiceIds = []
  modal.error = ''
}

function toggleService(serviceId: string) {
  if (modal.selectedServiceIds.includes(serviceId)) {
    modal.selectedServiceIds = modal.selectedServiceIds.filter(id => id !== serviceId)
    return
  }

  modal.selectedServiceIds = [...modal.selectedServiceIds, serviceId]
}

function serviceCardDescription(service: { duration_minutes: number; price_cents: number | null }) {
  return `${formatDuration(service.duration_minutes)} - ${formatPrice(service.price_cents)}`
}

function serviceOptionStyle(selected: boolean) {
  if (selected) {
    return {
      background: 'rgba(var(--color-brand-primary-rgb), 0.08)',
      borderColor: 'rgba(var(--color-brand-primary-rgb), 0.28)',
      color: 'var(--color-text)',
    }
  }

  return {
    background: 'var(--color-surface-muted)',
    borderColor: 'var(--color-border)',
    color: 'var(--color-text)',
  }
}

function serviceCheckStyle(selected: boolean) {
  if (selected) {
    return {
      background: 'var(--color-brand-primary)',
      borderColor: 'var(--color-brand-primary)',
      color: 'var(--color-brand-on-primary)',
    }
  }

  return {
    background: 'transparent',
    borderColor: 'var(--color-border-strong)',
    color: 'transparent',
  }
}

function serviceMetaStyle(selected: boolean) {
  return {
    color: selected ? 'var(--color-text-soft)' : 'var(--color-text-muted)',
  }
}

function serviceSummary(professional: Professional) {
  if (professional.service_ids.length === 0) {
    return 'Sem servicos vinculados ao agendamento publico'
  }

  const names = professional.service_ids
    .map(serviceId => activeServices.value.find(service => service.id === serviceId)?.name)
    .filter((name): name is string => Boolean(name))

  if (names.length === 0) {
    return `${professional.service_ids.length} servico(s) vinculado(s)`
  }

  const preview = names.slice(0, 2).join(', ')
  const remaining = names.length - 2
  return remaining > 0 ? `${preview} e mais ${remaining}` : preview
}

async function handleSave() {
  modal.saving = true
  modal.error = ''
  try {
    const data = {
      name: modal.name,
      phone: modal.phone || undefined,
    }
    if (modal.id) {
      await updateProfessional(modal.id, data)
      await updateServices(modal.id, modal.selectedServiceIds)
    } else {
      const created = await createProfessional(data)
      modal.id = created.id
      await updateServices(created.id, modal.selectedServiceIds)
    }
    closeModal()
  } catch (e: any) {
    modal.error = e?.data?.error?.message ?? 'Erro ao salvar.'
  } finally {
    modal.saving = false
  }
}

function handleDelete(professional: Professional) {
  deleteModal.open = true
  deleteModal.id = professional.id
  deleteModal.name = professional.name
}

function closeDeleteModal() {
  deleteModal.open = false
  deleteModal.id = ''
  deleteModal.name = ''
  deleteModal.loading = false
}

async function confirmDelete() {
  if (!deleteModal.id) return
  deleteModal.loading = true
  try {
    await deleteProfessional(deleteModal.id)
    closeDeleteModal()
  } finally {
    deleteModal.loading = false
  }
}
</script>

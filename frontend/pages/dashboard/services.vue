<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="ds-kicker">Catalogo</p>
        <h1 class="ds-title mt-1">Servicos</h1>
        <p class="mt-2 text-sm leading-6" style="color: var(--color-text-muted);">
          Monte o catalogo que aparece para o cliente escolher no fluxo publico.
        </p>
      </div>
      <AppButton variant="primary" @click="openModal()">+ Novo servico</AppButton>
    </div>

    <div v-if="loading" class="text-sm" style="color: var(--color-text-muted);">Carregando servicos...</div>

    <div
      v-if="error"
      class="rounded-[1.2rem] border px-5 py-4 text-sm"
      style="background: var(--color-danger-soft); border-color: var(--color-danger); color: var(--color-danger);"
    >
      <p>{{ error }}</p>
      <AppButton size="sm" variant="secondary" class="mt-3" @click="reload">Tentar novamente</AppButton>
    </div>

    <AppSurface v-else-if="!loading && services.length === 0" tone="default" padding="lg">
      <div class="ds-empty-state">
        <p>Nenhum servico cadastrado ainda.</p>
        <AppButton variant="secondary" @click="openModal()">Cadastrar primeiro servico</AppButton>
      </div>
    </AppSurface>

    <AppSurface v-else-if="services.length > 0" tone="default" padding="none">
      <ul class="ds-grid-table divide-y" style="border-color: var(--color-border);">
        <li
          v-for="service in services"
          :key="service.id"
          class="ds-grid-row flex items-center justify-between gap-4 px-6 py-4"
        >
          <div>
            <p class="text-sm font-semibold" style="color: var(--color-text);">{{ service.name }}</p>
            <p class="mt-0.5 text-xs" style="color: var(--color-text-soft);">
              {{ formatDuration(service.duration_minutes) }} · {{ formatPrice(service.price_cents) }}
            </p>
          </div>

          <div class="flex items-center gap-3">
            <AppStatusPill :tone="service.active ? 'success' : 'neutral'">
              {{ service.active ? 'Ativo' : 'Inativo' }}
            </AppStatusPill>
            <AppButton size="sm" variant="ghost" @click="openModal(service)">Editar</AppButton>
            <AppButton size="sm" variant="danger" @click="handleDelete(service)">Remover</AppButton>
          </div>
        </li>
      </ul>
    </AppSurface>

    <AppModal
      :open="modal.open"
      :title="modal.id ? 'Editar servico' : 'Novo servico'"
      eyebrow="Catalogo"
      width="sm"
      @close="closeModal"
    >
      <form class="space-y-4" @submit.prevent="handleSave">
        <div>
          <label class="ds-label">Nome</label>
          <input v-model="modal.name" type="text" required class="ds-input" />
        </div>
        <div>
          <label class="ds-label">Descricao (opcional)</label>
          <input v-model="modal.description" type="text" class="ds-input" />
        </div>
        <div>
          <label class="ds-label">Duracao</label>
          <select v-model.number="modal.duration_minutes" class="ds-select">
            <option :value="15">15 min</option>
            <option :value="30">30 min</option>
            <option :value="45">45 min</option>
            <option :value="60">1h</option>
            <option :value="90">1h 30min</option>
            <option :value="120">2h</option>
            <option :value="0">Personalizado</option>
          </select>
          <input
            v-if="modal.duration_minutes === 0"
            v-model.number="modal.custom_duration"
            type="number"
            min="5"
            class="ds-input mt-2"
            placeholder="Duracao em minutos"
          />
        </div>
        <div>
          <label class="ds-label">Preco (R$)</label>
          <input v-model="modal.price_reais" type="number" min="0" step="0.01" class="ds-input" placeholder="0,00" />
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
      title="Remover servico"
      description="Essa acao remove o item do catalogo e pode afetar o fluxo publico de reservas."
      eyebrow="Catalogo"
      message="Deseja remover este servico?"
      :details="deleteModal.name ? `Servico: ${deleteModal.name}` : ''"
      confirm-label="Remover servico"
      loading-label="Removendo..."
      :loading="deleteModal.loading"
      @cancel="closeDeleteModal"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import type { Service } from '~/composables/useServices'

definePageMeta({ layout: 'manager', middleware: 'auth' })

const {
  services,
  loading,
  error,
  fetchServices,
  createService,
  updateService,
  deleteService,
  formatPrice,
  formatDuration,
} = useServices()

const modal = reactive({
  open: false,
  id: '',
  name: '',
  description: '',
  duration_minutes: 30,
  custom_duration: 60,
  price_reais: '',
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

async function reload() {
  await fetchServices()
}

function openModal(service?: Service) {
  modal.open = true
  modal.id = service?.id ?? ''
  modal.name = service?.name ?? ''
  modal.description = service?.description ?? ''
  modal.duration_minutes = [15, 30, 45, 60, 90, 120].includes(service?.duration_minutes ?? 30)
    ? (service?.duration_minutes ?? 30)
    : 0
  modal.custom_duration = service?.duration_minutes ?? 60
  modal.price_reais = service?.price_cents != null ? String(service.price_cents / 100) : ''
  modal.error = ''
}

function closeModal() {
  modal.open = false
  modal.id = ''
  modal.name = ''
  modal.description = ''
  modal.duration_minutes = 30
  modal.custom_duration = 60
  modal.price_reais = ''
  modal.error = ''
}

function effectiveDuration(): number {
  return modal.duration_minutes === 0 ? modal.custom_duration : modal.duration_minutes
}

async function handleSave() {
  modal.saving = true
  modal.error = ''
  try {
    const data = {
      name: modal.name,
      description: modal.description || undefined,
      duration_minutes: effectiveDuration(),
      price_cents: modal.price_reais !== '' ? Math.round(Number(modal.price_reais) * 100) : undefined,
    }
    if (modal.id) {
      await updateService(modal.id, data)
    } else {
      await createService(data)
    }
    closeModal()
  } catch (e: any) {
    modal.error = e?.data?.error?.message ?? 'Erro ao salvar.'
  } finally {
    modal.saving = false
  }
}

function handleDelete(service: Service) {
  deleteModal.open = true
  deleteModal.id = service.id
  deleteModal.name = service.name
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
    await deleteService(deleteModal.id)
    closeDeleteModal()
  } finally {
    deleteModal.loading = false
  }
}
</script>

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
  createProfessional,
  updateProfessional,
  deleteProfessional,
} = useProfessionals()

const modal = reactive({
  open: false,
  id: '',
  name: '',
  phone: '',
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
  await fetchProfessionals()
}

function openModal(professional?: Professional) {
  modal.open = true
  modal.id = professional?.id ?? ''
  modal.name = professional?.name ?? ''
  modal.phone = professional?.phone ?? ''
  modal.error = ''
}

function closeModal() {
  modal.open = false
  modal.id = ''
  modal.name = ''
  modal.phone = ''
  modal.error = ''
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
    } else {
      await createProfessional(data)
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

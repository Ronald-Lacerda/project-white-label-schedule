<template>
  <AppModal
    :open="open"
    :title="title"
    :description="description"
    :eyebrow="eyebrow"
    width="sm"
    @close="$emit('cancel')"
  >
    <div class="flex items-start gap-4">
      <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-[1rem]" :class="iconClass">
        <span class="text-lg font-semibold">{{ icon }}</span>
      </div>
      <div class="space-y-2">
        <p class="text-sm font-semibold" style="color: var(--color-text);">{{ message }}</p>
        <p v-if="details" class="text-sm leading-6" style="color: var(--color-text-muted);">{{ details }}</p>
      </div>
    </div>

    <template #footer>
      <div class="flex flex-wrap justify-end gap-3">
        <AppButton variant="secondary" :disabled="loading" @click="$emit('cancel')">
          {{ cancelLabel }}
        </AppButton>
        <AppButton :variant="confirmVariant" :disabled="loading" @click="$emit('confirm')">
          {{ loading ? loadingLabel : confirmLabel }}
        </AppButton>
      </div>
    </template>
  </AppModal>
</template>

<script setup lang="ts">
import { copy } from '~/constants/copy'

const props = withDefaults(defineProps<{
  open: boolean
  title: string
  description?: string
  eyebrow?: string
  message: string
  details?: string
  confirmLabel?: string
  cancelLabel?: string
  loadingLabel?: string
  confirmVariant?: 'primary' | 'danger'
  tone?: 'danger' | 'warning' | 'info'
  loading?: boolean
}>(), {
  description: '',
  eyebrow: copy.confirmModal.eyebrow,
  details: '',
  confirmLabel: copy.common.confirm,
  cancelLabel: copy.common.cancel,
  loadingLabel: copy.common.processing,
  confirmVariant: 'danger',
  tone: 'danger',
  loading: false,
})

defineEmits<{
  cancel: []
  confirm: []
}>()

const icon = computed(() => ({
  danger: '!',
  warning: '?',
  info: 'i',
}[props.tone]))

const iconClass = computed(() => ({
  'bg-[var(--color-danger-soft)] text-[var(--color-danger)]': props.tone === 'danger',
  'bg-[var(--color-warning-soft)] text-[var(--color-warning)]': props.tone === 'warning',
  'bg-[var(--color-info-soft)] text-[var(--color-info)]': props.tone === 'info',
}))
</script>

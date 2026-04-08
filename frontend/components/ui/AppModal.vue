<template>
  <Teleport to="body">
    <div v-if="open" class="ds-modal-backdrop" @click.self="$emit('close')">
      <div class="ds-modal-panel" :class="widthClass">
        <div class="border-b px-6 py-5 md:px-7" style="border-color: var(--color-border);">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p v-if="eyebrow" class="ds-kicker">{{ eyebrow }}</p>
              <h2 class="mt-1 text-xl font-semibold" style="color: var(--color-text);">{{ title }}</h2>
              <p v-if="description" class="mt-2 text-sm" style="color: var(--color-text-muted);">{{ description }}</p>
            </div>
            <button
              type="button"
              class="ds-button ds-button-ghost ds-button-sm"
              aria-label="Fechar modal"
              @click="$emit('close')"
            >
              Fechar
            </button>
          </div>
        </div>

        <div class="max-h-[70vh] overflow-y-auto px-6 py-6 md:px-7">
          <slot />
        </div>

        <div v-if="$slots.footer" class="border-t px-6 py-4 md:px-7" style="border-color: var(--color-border);">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineEmits<{ close: [] }>()

const props = withDefaults(defineProps<{
  open: boolean
  title: string
  description?: string
  eyebrow?: string
  width?: 'sm' | 'md' | 'lg'
}>(), {
  description: '',
  eyebrow: '',
  width: 'md',
})

const widthClass = computed(() => ({
  sm: 'max-w-lg',
  md: 'max-w-2xl',
  lg: 'max-w-4xl',
}[props.width]))
</script>

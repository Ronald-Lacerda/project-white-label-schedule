<template>
  <button
    type="button"
    class="w-full rounded-[1.4rem] border px-4 py-4 text-left transition"
    :style="cardStyle"
    @click="$emit('select')"
  >
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0">
        <p class="text-base font-semibold">{{ title }}</p>
        <p v-if="description" class="mt-1 text-sm leading-6" :style="subtleStyle">
          {{ description }}
        </p>
        <slot />
      </div>

      <div v-if="$slots.meta || meta" class="shrink-0 text-right">
        <slot name="meta">
          <p class="text-sm font-semibold">{{ meta }}</p>
        </slot>
      </div>
    </div>
  </button>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  title: string
  description?: string
  meta?: string
  selected?: boolean
}>(), {
  description: '',
  meta: '',
  selected: false,
})

defineEmits<{ select: [] }>()

const cardStyle = computed(() => {
  if (props.selected) {
    return {
      background: 'var(--color-brand-primary)',
      color: 'var(--color-brand-on-primary)',
      borderColor: 'transparent',
      boxShadow: 'var(--shadow-soft)',
    }
  }

  return {
    background: 'var(--color-surface-muted)',
    color: 'var(--color-text)',
    borderColor: 'var(--color-border)',
  }
})

const subtleStyle = computed(() => ({
  color: props.selected ? 'rgba(255, 255, 255, 0.72)' : 'var(--color-text-muted)',
}))
</script>

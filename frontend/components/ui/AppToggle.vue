<template>
  <button
    type="button"
    class="inline-flex items-center gap-3 rounded-full border px-3 py-2 text-sm font-semibold transition"
    :style="wrapperStyle"
    :aria-pressed="checked ? 'true' : 'false'"
    @click="$emit('update:checked', !checked)"
  >
    <span class="relative inline-flex h-6 w-11 flex-shrink-0 rounded-full transition" :style="trackStyle">
      <span class="absolute top-0.5 h-5 w-5 rounded-full bg-white shadow-sm transition" :style="thumbStyle" />
    </span>
    <span>{{ checked ? checkedLabel : uncheckedLabel }}</span>
  </button>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  checked: boolean
  checkedLabel?: string
  uncheckedLabel?: string
}>(), {
  checkedLabel: 'Ativo',
  uncheckedLabel: 'Inativo',
})

defineEmits<{ 'update:checked': [value: boolean] }>()

const wrapperStyle = computed(() => ({
  borderColor: props.checked ? 'rgba(var(--color-brand-primary-rgb), 0.18)' : 'var(--color-border)',
  background: props.checked ? 'rgba(var(--color-brand-primary-rgb), 0.06)' : 'var(--color-surface)',
  color: props.checked ? 'var(--color-brand-primary)' : 'var(--color-text-muted)',
}))

const trackStyle = computed(() => ({
  background: props.checked ? 'var(--color-brand-primary)' : 'rgba(93, 103, 123, 0.26)',
}))

const thumbStyle = computed(() => ({
  left: props.checked ? 'calc(100% - 1.35rem)' : '0.125rem',
}))
</script>

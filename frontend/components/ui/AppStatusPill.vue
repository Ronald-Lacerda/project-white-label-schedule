<template>
  <span :class="pillClass">
    <slot />
  </span>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  tone?: 'neutral' | 'success' | 'warning' | 'danger' | 'info' | 'brand'
}>(), {
  tone: 'neutral',
})

const pillClass = computed(() => [
  'ds-status',
  props.tone === 'neutral' && 'ds-status-neutral',
  props.tone === 'success' && 'ds-status-success',
  props.tone === 'warning' && 'ds-status-warning',
  props.tone === 'danger' && 'ds-status-danger',
  props.tone === 'info' && 'ds-status-info',
  props.tone === 'brand' && 'ds-status-brand',
])
</script>

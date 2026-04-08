<template>
  <section :class="surfaceClass">
    <slot />
  </section>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  tone?: 'default' | 'muted' | 'brand'
  padding?: 'md' | 'lg' | 'none'
}>(), {
  tone: 'default',
  padding: 'md',
})

const surfaceClass = computed(() => [
  props.tone === 'default' && 'ds-surface',
  props.tone === 'muted' && 'ds-surface-muted',
  props.tone === 'brand' && 'ds-surface-brand',
  props.padding === 'md' && 'ds-panel',
  props.padding === 'lg' && 'ds-panel-lg',
])
</script>

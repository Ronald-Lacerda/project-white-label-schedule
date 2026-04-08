<template>
  <NuxtLink
    v-if="to"
    v-bind="$attrs"
    :to="to"
    :class="buttonClass"
    :aria-disabled="disabled ? 'true' : undefined"
    :tabindex="disabled ? -1 : undefined"
  >
    <slot />
  </NuxtLink>

  <button
    v-else
    v-bind="$attrs"
    :type="type"
    :disabled="disabled"
    :class="buttonClass"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  to?: string
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  block?: boolean
  disabled?: boolean
}>(), {
  to: undefined,
  type: 'button',
  variant: 'secondary',
  size: 'md',
  block: false,
  disabled: false,
})

const buttonClass = computed(() => [
  'ds-button',
  props.variant === 'primary' && 'ds-button-primary',
  props.variant === 'secondary' && 'ds-button-secondary',
  props.variant === 'ghost' && 'ds-button-ghost',
  props.variant === 'danger' && 'ds-button-danger',
  props.size === 'sm' && 'ds-button-sm',
  props.size === 'lg' && 'ds-button-lg',
  props.block && 'w-full',
  props.disabled && 'pointer-events-none',
])

const { to, type, disabled } = toRefs(props)
</script>

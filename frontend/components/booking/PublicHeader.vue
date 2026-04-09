<template>
  <AppSurface tone="brand" padding="lg">
    <div class="relative overflow-hidden" :style="themeStyle">
      <div class="absolute inset-y-0 right-0 hidden w-40 rounded-full blur-3xl md:block" :style="{ background: `rgba(${theme.secondaryRgb}, 0.16)` }" />

      <div class="relative flex flex-col gap-6">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="max-w-2xl">
            <p class="ds-kicker">{{ eyebrow }}</p>
            <h1 class="ds-display-title mt-3">{{ title }}</h1>
            <p class="ds-subtitle mt-3">{{ subtitle }}</p>
          </div>

          <div
            class="flex h-20 w-20 items-center justify-center overflow-hidden rounded-[1.8rem] border bg-white/80 shadow-sm"
            style="border-color: rgba(var(--color-brand-primary-rgb), 0.12);"
          >
            <img v-if="logoUrl" :src="logoUrl" :alt="title" class="h-full w-full object-cover" />
            <span v-else class="text-xs font-semibold uppercase tracking-[0.32em]" style="color: var(--color-brand-primary);">
              {{ initials }}
            </span>
          </div>
        </div>

        <div class="ds-stepper">
          <div
            v-for="item in steps"
            :key="item.step"
            :class="stepClass(item.step)"
          >
            {{ item.label }}
          </div>
        </div>
      </div>
    </div>
  </AppSurface>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  title: string
  subtitle: string
  logoUrl?: string | null
  primaryColor?: string | null
  secondaryColor?: string | null
  activeStep?: number
  eyebrow?: string
}>(), {
  logoUrl: null,
  primaryColor: null,
  secondaryColor: null,
  activeStep: 1,
  eyebrow: 'Agendamento online',
})

const steps = [
  { step: 1, label: '1. Servico' },
  { step: 2, label: '2. Profissional' },
  { step: 3, label: '3. Horario' },
  { step: 4, label: '4. Dados' },
  { step: 5, label: '5. Confirmacao' },
]

const theme = computed(() => buildBrandTheme({
  primaryColor: props.primaryColor,
  secondaryColor: props.secondaryColor,
}))

const themeStyle = computed(() => brandThemeStyle({
  primaryColor: props.primaryColor,
  secondaryColor: props.secondaryColor,
}))

const initials = computed(() => props.title.trim().slice(0, 2).toUpperCase())

function stepClass(step: number) {
  if (props.activeStep > step) return 'ds-step ds-step-complete'
  if (props.activeStep === step) return 'ds-step ds-step-active'
  return 'ds-step'
}
</script>

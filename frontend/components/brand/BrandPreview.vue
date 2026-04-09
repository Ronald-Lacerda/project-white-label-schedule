<template>
  <div class="ds-surface-brand ds-panel-lg" :style="themeStyle">
    <div class="space-y-6">
      <div class="flex items-start justify-between gap-4">
        <div class="space-y-3">
          <p class="ds-kicker">Preview da marca</p>
          <div>
            <h3 class="font-display text-2xl font-semibold tracking-tight" style="color: var(--color-brand-text);">
              {{ name || 'Seu estabelecimento' }}
            </h3>
            <p class="mt-1.5 max-w-sm text-sm leading-6" style="color: var(--color-text-muted);">
              Referência visual do link público com a identidade da marca aplicada.
            </p>
          </div>
        </div>

        <div
          class="flex h-14 w-14 flex-shrink-0 items-center justify-center overflow-hidden rounded-[1.2rem] border bg-white/80"
          style="border-color: rgba(var(--color-brand-primary-rgb), 0.14);"
        >
          <img v-if="logoUrl" :src="logoUrl" :alt="name || 'Logo'" class="h-full w-full object-cover" />
          <span v-else class="text-[10px] font-bold uppercase tracking-[0.3em]" style="color: var(--color-brand-primary);">
            {{ initials }}
          </span>
        </div>
      </div>

      <div class="grid gap-3 md:grid-cols-[1.4fr_1fr]">
        <!-- Cartão: link público + stepper -->
        <div class="rounded-[1.4rem] border bg-white/85 p-5" style="border-color: rgba(var(--color-brand-primary-rgb), 0.1);">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-[0.68rem] font-semibold uppercase tracking-[0.28em]" style="color: var(--color-text-soft);">Link público</p>
              <p class="mt-1.5 text-sm font-semibold break-all" style="color: var(--color-brand-primary);">
                /page/{{ slug || 'seu-link' }}
              </p>
            </div>
            <span class="ds-chip flex-shrink-0">24/7</span>
          </div>

          <div class="mt-4 ds-stepper" style="grid-template-columns: repeat(3, 1fr);">
            <div class="ds-step ds-step-complete">1. Serviço</div>
            <div class="ds-step ds-step-active">2. Profissional</div>
            <div class="ds-step">3. Horário</div>
          </div>
        </div>

        <!-- Cartão: swatches -->
        <div class="rounded-[1.4rem] border bg-white/88 p-5" style="border-color: rgba(var(--color-brand-primary-rgb), 0.1);">
          <p class="text-[0.68rem] font-semibold uppercase tracking-[0.28em]" style="color: var(--color-text-soft);">Paleta</p>
          <div class="mt-4 grid grid-cols-2 gap-3">
            <div class="space-y-1.5">
              <div class="h-10 rounded-xl" :style="{ background: theme.primary }" />
              <p class="text-[0.7rem] font-medium" style="color: var(--color-text-muted);">Primária</p>
            </div>
            <div class="space-y-1.5">
              <div class="h-10 rounded-xl" :style="{ background: theme.secondary }" />
              <p class="text-[0.7rem] font-medium" style="color: var(--color-text-muted);">Secundária</p>
            </div>
            <div class="space-y-1.5">
              <div class="h-10 rounded-xl border" :style="{ background: theme.surface, borderColor: 'rgba(var(--color-brand-primary-rgb), 0.1)' }" />
              <p class="text-[0.7rem] font-medium" style="color: var(--color-text-muted);">Superfície</p>
            </div>
            <div class="space-y-1.5">
              <div class="h-10 rounded-xl border" :style="{ background: theme.accent, borderColor: 'rgba(var(--color-brand-primary-rgb), 0.1)' }" />
              <p class="text-[0.7rem] font-medium" style="color: var(--color-text-muted);">Destaque</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  name?: string
  logoUrl?: string | null
  primaryColor?: string | null
  secondaryColor?: string | null
  slug?: string
}>(), {
  name: '',
  logoUrl: null,
  primaryColor: null,
  secondaryColor: null,
  slug: '',
})

const theme = computed(() => buildBrandTheme({
  primaryColor: props.primaryColor,
  secondaryColor: props.secondaryColor,
}))

const themeStyle = computed(() => brandThemeStyle({
  primaryColor: props.primaryColor,
  secondaryColor: props.secondaryColor,
}))

const initials = computed(() => (props.name || 'WL').trim().slice(0, 2).toUpperCase())
</script>

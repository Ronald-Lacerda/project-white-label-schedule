<template>
  <AppSurface tone="default" padding="lg">
    <div class="space-y-5">
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="ds-kicker">Resumo</p>
          <h3 class="mt-2 text-xl font-semibold" style="color: var(--color-text);">
            {{ title }}
          </h3>
        </div>
        <AppStatusPill tone="brand">{{ badge }}</AppStatusPill>
      </div>

      <dl class="grid gap-4 text-sm md:grid-cols-2">
        <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
          <dt class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Serviço</dt>
          <dd class="mt-2 font-semibold" style="color: var(--color-text);">{{ serviceName }}</dd>
        </div>
        <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
          <dt class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Profissional</dt>
          <dd class="mt-2 font-semibold" style="color: var(--color-text);">{{ professionalName }}</dd>
        </div>
        <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
          <dt class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Data</dt>
          <dd class="mt-2 font-semibold" style="color: var(--color-text);">{{ dateLabel }}</dd>
        </div>
        <div class="rounded-[1.2rem] border p-4" style="border-color: var(--color-border); background: var(--color-surface-muted);">
          <dt class="text-xs font-semibold uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Horário</dt>
          <dd class="mt-2 font-semibold" style="color: var(--color-text);">{{ timeLabel }}</dd>
        </div>
      </dl>
    </div>
  </AppSurface>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title?: string
  badge?: string
  serviceName: string
  professionalName: string
  dateLabel: string
  timeLabel: string
}>(), {
  title: 'Confira os detalhes antes de confirmar',
  badge: 'Pronto para reservar',
})
</script>

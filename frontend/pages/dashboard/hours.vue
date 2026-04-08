<template>
  <div class="ds-page max-w-4xl">
    <div>
      <p class="ds-kicker">Disponibilidade base</p>
      <h1 class="ds-title mt-1">Horarios de funcionamento</h1>
      <p class="mt-2 text-sm leading-6" style="color: var(--color-text-muted);">
        Defina os horarios base do estabelecimento para alimentar o motor de disponibilidade.
      </p>
    </div>

    <div v-if="loading" class="text-sm" style="color: var(--color-text-muted);">Carregando...</div>

    <div
      v-else-if="error"
      class="rounded-[1.2rem] border px-5 py-4 text-sm"
      style="background: var(--color-danger-soft); border-color: var(--color-danger); color: var(--color-danger);"
    >
      {{ error }}
    </div>

    <AppSurface v-else tone="default" padding="none">
      <ul class="ds-grid-table divide-y" style="border-color: var(--color-border);">
        <li
          v-for="hour in hours"
          :key="hour.day_of_week"
          class="flex flex-col gap-4 px-6 py-5 md:flex-row md:items-center md:justify-between"
        >
          <div class="space-y-1">
            <p class="text-sm font-semibold" style="color: var(--color-text);">
              {{ dayName(hour.day_of_week) }}
            </p>
            <p class="text-xs uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">
              {{ hour.is_closed ? 'Fechado' : `${formattedHour(hour.open_time)} as ${formattedHour(hour.close_time)}` }}
            </p>
          </div>

          <div class="flex flex-col gap-3 md:items-end">
            <AppToggle
              :checked="!hour.is_closed"
              checked-label="Aberto"
              unchecked-label="Fechado"
              @update:checked="hour.is_closed = !$event"
            />

            <div v-if="!hour.is_closed" class="flex flex-wrap items-center gap-3">
              <input v-model="hour.open_time" type="time" class="ds-input w-auto min-w-[8.5rem] px-3 py-2" />
              <span class="text-sm" style="color: var(--color-text-soft);">ate</span>
              <input v-model="hour.close_time" type="time" class="ds-input w-auto min-w-[8.5rem] px-3 py-2" />
            </div>
          </div>
        </li>
      </ul>
    </AppSurface>

    <AppSurface tone="muted" padding="lg">
      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm font-semibold" style="color: var(--color-text);">Revisao final</p>
          <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
            Esses horarios alimentam o motor de disponibilidade e o fluxo publico de reservas.
          </p>
        </div>

        <div class="flex items-center gap-4">
          <p v-if="success" class="text-sm" style="color: var(--color-success);">Salvo com sucesso.</p>
          <AppButton variant="primary" :disabled="saving || loading" @click="save">
            {{ saving ? 'Salvando...' : 'Salvar horarios' }}
          </AppButton>
        </div>
      </div>
    </AppSurface>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'manager', middleware: 'auth' })

const { hours, loading, error, fetchHours, saveHours, dayName } = useBusinessHours()

const saving = ref(false)
const success = ref(false)

onMounted(() => {
  fetchHours()
})

function formattedHour(value: string) {
  return value.slice(0, 5)
}

async function save() {
  saving.value = true
  success.value = false
  try {
    await saveHours()
    success.value = true
    setTimeout(() => { success.value = false }, 3000)
  } finally {
    saving.value = false
  }
}
</script>

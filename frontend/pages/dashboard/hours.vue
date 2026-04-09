<template>
  <div class="ds-page max-w-4xl">
    <div class="space-y-2">
      <h1 class="ds-title">Horários</h1>
      <p class="text-sm leading-6" style="color: var(--color-text-muted);">
        Defina a disponibilidade base do estabelecimento para alimentar o motor de horários e a agenda pública.
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
              {{ hour.is_closed ? 'Fechado' : `${formatShortTimeBr(hour.open_time)} às ${formatShortTimeBr(hour.close_time)}` }}
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
              <div class="flex items-center gap-2">
                <select
                  :value="getTimePart(hour.open_time, 'hour')"
                  class="ds-select w-auto min-w-[5.5rem] px-3 py-2"
                  @change="onTimeSelect(hour, 'open_time', 'hour', $event)"
                >
                  <option v-for="option in hourOptions" :key="`open-hour-${hour.day_of_week}-${option}`" :value="option">{{ option }}</option>
                </select>
                <span class="text-sm font-semibold" style="color: var(--color-text-soft);">:</span>
                <select
                  :value="getTimePart(hour.open_time, 'minute')"
                  class="ds-select w-auto min-w-[5.5rem] px-3 py-2"
                  @change="onTimeSelect(hour, 'open_time', 'minute', $event)"
                >
                  <option v-for="option in minuteOptions" :key="`open-minute-${hour.day_of_week}-${option}`" :value="option">{{ option }}</option>
                </select>
              </div>
              <span class="text-sm" style="color: var(--color-text-soft);">até</span>
              <div class="flex items-center gap-2">
                <select
                  :value="getTimePart(hour.close_time, 'hour')"
                  class="ds-select w-auto min-w-[5.5rem] px-3 py-2"
                  @change="onTimeSelect(hour, 'close_time', 'hour', $event)"
                >
                  <option v-for="option in hourOptions" :key="`close-hour-${hour.day_of_week}-${option}`" :value="option">{{ option }}</option>
                </select>
                <span class="text-sm font-semibold" style="color: var(--color-text-soft);">:</span>
                <select
                  :value="getTimePart(hour.close_time, 'minute')"
                  class="ds-select w-auto min-w-[5.5rem] px-3 py-2"
                  @change="onTimeSelect(hour, 'close_time', 'minute', $event)"
                >
                  <option v-for="option in minuteOptions" :key="`close-minute-${hour.day_of_week}-${option}`" :value="option">{{ option }}</option>
                </select>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </AppSurface>

    <AppSurface tone="muted" padding="lg">
      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm font-semibold" style="color: var(--color-text);">Revisão final</p>
          <p class="mt-1 text-sm" style="color: var(--color-text-muted);">
            Esses horários alimentam o motor de disponibilidade e o fluxo público de reservas.
          </p>
        </div>

        <div class="flex items-center gap-4">
          <p v-if="success" class="text-sm" style="color: var(--color-success);">Salvo com sucesso.</p>
          <AppButton variant="primary" :disabled="saving || loading" @click="save">
            {{ saving ? 'Salvando...' : 'Salvar horários' }}
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
const hourOptions = Array.from({ length: 24 }, (_, index) => String(index).padStart(2, '0'))
const minuteOptions = Array.from({ length: 60 }, (_, index) => String(index).padStart(2, '0'))

onMounted(() => {
  fetchHours()
})

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

function normalizeTime(value: string) {
  if (!value) return '08:00:00'
  const [hours = '08', minutes = '00', seconds = '00'] = value.split(':')
  return `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}:${seconds.padStart(2, '0')}`
}

function getTimePart(value: string, part: 'hour' | 'minute') {
  const [hours, minutes] = normalizeTime(value).split(':')
  return part === 'hour' ? hours : minutes
}

function updateTimePart(hour: { open_time: string; close_time: string }, field: 'open_time' | 'close_time', part: 'hour' | 'minute', nextValue: string) {
  const [hours, minutes, seconds] = normalizeTime(hour[field]).split(':')
  hour[field] = part === 'hour'
    ? `${nextValue}:${minutes}:${seconds}`
    : `${hours}:${nextValue}:${seconds}`
}

function onTimeSelect(hour: { open_time: string; close_time: string }, field: 'open_time' | 'close_time', part: 'hour' | 'minute', event: Event) {
  const target = event.target as HTMLSelectElement | null
  if (!target) return
  updateTimePart(hour, field, part, target.value)
}
</script>

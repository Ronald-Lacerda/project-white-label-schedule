<template>
  <div class="ds-page">
    <div v-if="loading" class="text-sm" style="color: var(--color-text-muted);">
      Carregando...
    </div>

    <div v-else class="space-y-8">
      <AppSurface tone="default" padding="lg">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="ds-kicker">Integração</p>
            <h2 class="mt-2 text-xl font-semibold" style="color: var(--color-text);">Google Agenda</h2>
            <p class="mt-2 text-sm leading-6" style="color: var(--color-text-muted);">
              Conecte sua conta para sincronizar agendas por profissional e preparar os próximos passos do MVP.
            </p>
          </div>

          <AppStatusPill :tone="establishment?.google_calendar_connected ? 'success' : 'neutral'">
            {{ establishment?.google_calendar_connected ? 'Conectado' : 'Não conectado' }}
          </AppStatusPill>
        </div>

        <div class="mt-4 flex flex-wrap gap-3">
          <AppButton
            v-if="establishment?.google_calendar_connected"
            to="/dashboard/settings/google"
            variant="primary"
          >
            Gerenciar conexão
          </AppButton>
          <button
            v-else
            type="button"
            class="ds-button ds-button-primary"
            :disabled="connecting"
            @click="handleConnect"
          >
            <AppGoogleIcon />
            {{ connecting ? 'Redirecionando...' : 'Conectar Google Agenda' }}
          </button>
        </div>
        <p v-if="connectError" class="mt-3 text-sm" style="color: var(--color-danger);">{{ connectError }}</p>
      </AppSurface>

      <div>
        <AppSurface tone="default" padding="lg">
          <h2 class="mb-4 text-xl font-semibold" style="color: var(--color-text);">Dados do estabelecimento</h2>

          <form class="space-y-4" @submit.prevent="saveEstablishment">
            <div class="grid grid-cols-1 gap-4 xl:grid-cols-12 xl:items-start">
              <div class="xl:col-span-6">
                <label class="ds-label">Nome</label>
                <input v-model="estForm.name" type="text" required class="ds-input" />
              </div>

              <div class="xl:col-span-3">
                <label class="ds-label">Fuso horário</label>
                <select v-model="estForm.timezone" class="ds-select">
                  <option value="America/Sao_Paulo">America/Sao Paulo</option>
                  <option value="America/Manaus">America/Manaus</option>
                  <option value="America/Belem">America/Belem</option>
                  <option value="America/Fortaleza">America/Fortaleza</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-1 gap-4 xl:grid-cols-12 xl:items-start">
              <div class="xl:col-span-8">
                <label class="ds-label">Slug (identificador da URL)</label>
                <div class="flex items-center gap-2">
                  <span class="text-sm" style="color: var(--color-text-soft);">/page/</span>
                  <input v-model="estForm.slug" type="text" required class="ds-input flex-1" />
                </div>
                <p class="mt-2 text-xs uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">
                  Link público: /page/{{ estForm.slug }}
                </p>
              </div>

              <div class="xl:col-span-4">
                <label class="ds-label">Antecedência mínima para cancelamento (horas)</label>
                <input v-model.number="estForm.min_advance_cancel_hours" type="number" min="0" class="ds-input" />
              </div>
            </div>

            <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
              <div>
                <label class="ds-label">E-mail de contato</label>
                <input v-model="estForm.contact_email" type="email" class="ds-input" />
              </div>
              <div class="lg:max-w-sm">
                <label class="ds-label">Telefone de contato</label>
                <input
                  :value="estForm.contact_phone"
                  type="tel"
                  inputmode="tel"
                  autocomplete="tel"
                  class="ds-input"
                  placeholder="(11) 99999-9999"
                  @input="onContactPhoneInput"
                />
              </div>
            </div>

            <p v-if="estError" class="text-sm" style="color: var(--color-danger);">{{ estError }}</p>
            <p v-if="estSuccess" class="text-sm" style="color: var(--color-success);">Salvo com sucesso!</p>

            <AppButton type="submit" variant="primary" :disabled="estSaving">
              {{ estSaving ? 'Salvando...' : 'Salvar' }}
            </AppButton>
          </form>
        </AppSurface>
      </div>

      <AppSurface tone="default" padding="lg">
        <h2 class="mb-4 text-xl font-semibold" style="color: var(--color-text);">Aparência (Whitelabel)</h2>

        <div class="space-y-6">
          <div>
            <label class="ds-label">Logo do estabelecimento</label>
            <div class="flex items-center gap-4">
              <div
                class="flex h-20 w-20 flex-shrink-0 items-center justify-center overflow-hidden rounded-[1.4rem] border bg-white/80"
                style="border-color: rgba(var(--color-brand-primary-rgb), 0.1);"
              >
                <img v-if="currentLogoUrl" :src="currentLogoUrl" alt="Logo" class="h-full w-full object-cover" />
                <span v-else class="text-xs" style="color: var(--color-text-soft);">Sem logo</span>
              </div>

              <div class="flex flex-col gap-2">
                <label class="cursor-pointer">
                  <span class="ds-button ds-button-secondary ds-button-sm">
                    {{ logoUploading ? 'Enviando...' : 'Escolher imagem' }}
                  </span>
                  <input
                    type="file"
                    accept=".jpg,.jpeg,.png,.webp,.svg"
                    class="hidden"
                    :disabled="logoUploading"
                    @change="handleLogoUpload"
                  />
                </label>

                <p class="text-xs uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">
                  JPG, PNG, WebP ou SVG. Max. 5 MB.
                </p>
                <p v-if="logoError" class="text-xs" style="color: var(--color-danger);">{{ logoError }}</p>
                <p v-if="logoSuccess" class="text-xs" style="color: var(--color-success);">Logo atualizada com sucesso!</p>
              </div>
            </div>
          </div>

          <form class="space-y-4" @submit.prevent="saveWhitelabel">
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label class="ds-label">Cor primária</label>
                <div class="flex items-center gap-2">
                  <input v-model="wlForm.primary_color" type="color" class="h-10 w-10 cursor-pointer rounded border" style="border-color: var(--color-border);" />
                  <input v-model="wlForm.primary_color" type="text" class="ds-input flex-1" placeholder="#000000" />
                </div>
              </div>

              <div>
                <label class="ds-label">Cor secundária</label>
                <div class="flex items-center gap-2">
                  <input v-model="wlForm.secondary_color" type="color" class="h-10 w-10 cursor-pointer rounded border" style="border-color: var(--color-border);" />
                  <input v-model="wlForm.secondary_color" type="text" class="ds-input flex-1" placeholder="#ffffff" />
                </div>
              </div>
            </div>

            <p class="text-sm leading-6" style="color: var(--color-text-muted);">
              O sistema deriva superficies, acentos e contraste a partir dessas duas cores para manter legibilidade.
            </p>
            <p v-if="wlSuccess" class="text-sm" style="color: var(--color-success);">Salvo com sucesso!</p>

            <AppButton type="submit" variant="primary" :disabled="wlSaving">
              {{ wlSaving ? 'Salvando...' : 'Salvar aparência' }}
            </AppButton>
          </form>
        </div>
      </AppSurface>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'manager', middleware: 'auth' })

const {
  establishment,
  whitelabel,
  loading,
  fetchEstablishment,
  updateEstablishment,
  fetchWhitelabel,
  updateWhitelabel,
  uploadLogo,
} = useEstablishment()

const { getAuthUrl } = useGoogleCalendar()

await Promise.all([fetchEstablishment(), fetchWhitelabel()])

const connecting = ref(false)
const connectError = ref<string | null>(null)

async function handleConnect() {
  connecting.value = true
  connectError.value = null
  try {
    const url = await getAuthUrl()
    window.location.href = url
  } catch (e: any) {
    connectError.value = e?.data?.error?.message ?? 'Erro ao iniciar conexão.'
    connecting.value = false
  }
}

const estForm = reactive({
  name: establishment.value?.name ?? '',
  slug: establishment.value?.slug ?? '',
  timezone: establishment.value?.timezone ?? 'America/Sao_Paulo',
  contact_email: establishment.value?.contact_email ?? '',
  contact_phone: formatBrazilianPhoneInput(establishment.value?.contact_phone ?? ''),
  min_advance_cancel_hours: establishment.value?.min_advance_cancel_hours ?? 0,
})

const wlForm = reactive({
  primary_color: whitelabel.value?.primary_color ?? '#000000',
  secondary_color: whitelabel.value?.secondary_color ?? '#ffffff',
})

const estSaving = ref(false)
const estError = ref('')
const estSuccess = ref(false)

const wlSaving = ref(false)
const wlSuccess = ref(false)

const currentLogoUrl = ref<string | null>(whitelabel.value?.logo_url ?? null)
const logoUploading = ref(false)
const logoError = ref('')
const logoSuccess = ref(false)

async function handleLogoUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  logoUploading.value = true
  logoError.value = ''
  logoSuccess.value = false

  try {
    currentLogoUrl.value = await uploadLogo(file)
    logoSuccess.value = true
    setTimeout(() => { logoSuccess.value = false }, 3000)
  } catch (e: any) {
    logoError.value = e?.data?.error?.message ?? 'Erro ao enviar a logo.'
  } finally {
    logoUploading.value = false
    input.value = ''
  }
}

async function saveEstablishment() {
  estSaving.value = true
  estError.value = ''
  estSuccess.value = false
  try {
    await updateEstablishment({
      ...estForm,
      contact_email: estForm.contact_email || null,
      contact_phone: normalizeBrazilianPhone(estForm.contact_phone) || null,
    })
    estSuccess.value = true
    setTimeout(() => {
      estSuccess.value = false
    }, 3000)
  } catch (e: any) {
    estError.value = e?.data?.error?.message ?? 'Erro ao salvar.'
  } finally {
    estSaving.value = false
  }
}

function onContactPhoneInput(event: Event) {
  const input = event.target as HTMLInputElement
  estForm.contact_phone = formatBrazilianPhoneInput(input.value)
}

async function saveWhitelabel() {
  wlSaving.value = true
  wlSuccess.value = false
  try {
    await updateWhitelabel({
      primary_color: wlForm.primary_color,
      secondary_color: wlForm.secondary_color || null,
    })
    wlSuccess.value = true
    setTimeout(() => {
      wlSuccess.value = false
    }, 3000)
  } finally {
    wlSaving.value = false
  }
}
</script>

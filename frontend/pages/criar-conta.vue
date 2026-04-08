<template>
  <div class="relative min-h-screen overflow-hidden">
    <div class="pointer-events-none absolute inset-0">
      <div class="absolute left-[-8%] top-0 h-72 w-72 rounded-full blur-3xl" style="background: rgba(var(--color-brand-secondary-rgb), 0.18);" />
      <div class="absolute right-[-10%] top-24 h-80 w-80 rounded-full blur-3xl" style="background: rgba(var(--color-brand-primary-rgb), 0.12);" />
      <div class="absolute inset-x-0 bottom-0 h-64 ds-dot-grid opacity-60" />
    </div>

    <div class="relative mx-auto flex min-h-screen max-w-6xl items-center px-4 py-10 sm:px-6 lg:px-8">
      <div class="grid w-full gap-6 lg:grid-cols-[1fr_560px] lg:items-start">
        <div class="hidden lg:block">
          <p class="ds-kicker">Onboarding</p>
          <h1 class="ds-display-title mt-4 max-w-xl">Crie sua operacao com identidade propria desde o primeiro acesso.</h1>
          <p class="ds-subtitle mt-4 max-w-lg">
            Cadastre o negocio, defina o link publico e entre direto no painel para concluir marca, equipe e horarios.
          </p>
        </div>

        <AppSurface tone="default" padding="lg">
          <p class="ds-kicker">Nova conta</p>
          <h2 class="ds-title mt-2">Criar conta</h2>
          <p class="mt-3 text-sm leading-6" style="color: var(--color-text-muted);">
            Cadastre seu estabelecimento e entre no painel automaticamente.
          </p>

          <form @submit.prevent="handleRegister" class="mt-6 space-y-4">
            <div>
              <label for="owner_name" class="ds-label">Nome do responsavel</label>
              <input
                id="owner_name"
                v-model="form.owner_name"
                type="text"
                required
                autocomplete="name"
                class="ds-input"
                placeholder="Seu nome"
              />
            </div>

            <div>
              <label for="establishment_name" class="ds-label">Nome do estabelecimento</label>
              <input
                id="establishment_name"
                v-model="form.establishment_name"
                type="text"
                required
                autocomplete="organization"
                class="ds-input"
                placeholder="Nome do seu negocio"
              />
            </div>

            <div>
              <label for="slug" class="ds-label">Slug da URL</label>
              <div class="flex items-center gap-2">
                <span class="text-sm" style="color: var(--color-text-soft);">/p/</span>
                <input
                  id="slug"
                  v-model="form.slug"
                  type="text"
                  required
                  class="ds-input"
                  placeholder="meu-estabelecimento"
                />
              </div>
              <p class="mt-2 text-xs uppercase tracking-[0.2em]" style="color: var(--color-text-soft);">Link publico: /p/{{ slugPreview }}</p>
            </div>

            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label for="email" class="ds-label">E-mail</label>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  required
                  autocomplete="email"
                  class="ds-input"
                  placeholder="voce@empresa.com"
                />
              </div>

              <div>
                <label for="contact_phone" class="ds-label">Telefone</label>
                <input
                  id="contact_phone"
                  v-model="form.contact_phone"
                  type="tel"
                  autocomplete="tel"
                  class="ds-input"
                  placeholder="(11) 99999-9999"
                />
              </div>
            </div>

            <div>
              <label for="password" class="ds-label">Senha</label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                required
                minlength="8"
                autocomplete="new-password"
                class="ds-input"
                placeholder="No minimo 8 caracteres"
              />
            </div>

            <div v-if="error" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
              {{ error }}
            </div>

            <AppButton type="submit" variant="primary" block :disabled="loading">
              {{ loading ? 'Criando conta...' : 'Criar conta' }}
            </AppButton>
          </form>

          <p class="mt-6 text-sm text-center" style="color: var(--color-text-muted);">
            Ja tem conta?
            <NuxtLink to="/login" class="font-semibold" style="color: var(--color-brand-primary);">
              Entrar
            </NuxtLink>
          </p>
        </AppSurface>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

const auth = useAuth()
const router = useRouter()

if (auth.isAuthenticated.value) {
  await navigateTo('/dashboard')
}

const form = reactive({
  owner_name: '',
  establishment_name: '',
  email: '',
  password: '',
  slug: '',
  contact_phone: '',
})

const loading = ref(false)
const error = ref('')
const slugTouched = ref(false)
const syncingSlug = ref(false)

const slugPreview = computed(() => form.slug || 'seu-link')

watch(
  () => form.establishment_name,
  (value) => {
    if (slugTouched.value) return
    syncingSlug.value = true
    form.slug = slugify(value)
  },
)

watch(
  () => form.slug,
  (value, previous) => {
    if (syncingSlug.value) {
      syncingSlug.value = false
      return
    }
    if (value !== previous) {
      slugTouched.value = true
      form.slug = slugify(value)
    }
  },
)

async function handleRegister() {
  loading.value = true
  error.value = ''

  try {
    await auth.register({
      owner_name: form.owner_name,
      establishment_name: form.establishment_name,
      email: form.email,
      password: form.password,
      slug: form.slug,
      contact_phone: form.contact_phone || null,
    })
    await router.push('/dashboard/settings')
  } catch (e: any) {
    const code = e?.data?.error?.code
    if (!e?.data?.error) {
      error.value = 'Nao foi possivel conectar ao servidor. Verifique se a API esta rodando em http://localhost:8080.'
    } else if (code === 'EMAIL_CONFLICT') {
      error.value = 'Este e-mail ja esta cadastrado.'
    } else if (code === 'SLUG_CONFLICT') {
      error.value = 'Este identificador de URL ja esta em uso.'
    } else if (code === 'INVALID_INPUT') {
      error.value = 'Revise os dados informados. A senha deve ter no minimo 8 caracteres.'
    } else {
      error.value = e?.data?.error?.message ?? 'Nao foi possivel criar a conta.'
    }
  } finally {
    loading.value = false
  }
}

function slugify(value: string) {
  return value
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9\s-_]/g, '')
    .replace(/[\s_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}
</script>

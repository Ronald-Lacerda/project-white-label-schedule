<template>
  <div class="relative min-h-screen overflow-hidden">
    <div class="pointer-events-none absolute inset-0">
      <div class="absolute left-[-10%] top-0 h-72 w-72 rounded-full blur-3xl" style="background: rgba(var(--color-brand-secondary-rgb), 0.18);" />
      <div class="absolute right-[-8%] top-20 h-80 w-80 rounded-full blur-3xl" style="background: rgba(var(--color-brand-primary-rgb), 0.12);" />
      <div class="absolute inset-x-0 bottom-0 h-56 ds-dot-grid opacity-60" />
    </div>

    <div class="relative mx-auto flex min-h-screen max-w-5xl items-center px-4 py-10 sm:px-6 lg:px-8">
      <div class="grid w-full gap-6 lg:grid-cols-[1.1fr_460px] lg:items-center">
        <div class="hidden lg:block">
          <p class="ds-kicker">White Label Schedule</p>
          <h1 class="ds-display-title mt-4 max-w-xl">Operacao clara para agenda, equipe e marca.</h1>
          <p class="ds-subtitle mt-4 max-w-lg">
            Entre no painel para organizar profissionais, servicos, horarios e acompanhar os agendamentos sem perder contexto.
          </p>
        </div>

        <AppSurface tone="default" padding="lg">
          <p class="ds-kicker">Acesso do gestor</p>
          <h2 class="ds-title mt-2">Entrar no painel</h2>
          <p class="mt-3 text-sm leading-6" style="color: var(--color-text-muted);">
            Use sua conta para acessar o ambiente administrativo do estabelecimento.
          </p>

          <form @submit.prevent="handleLogin" class="mt-6 space-y-4">
            <div>
              <label for="email" class="ds-label">E-mail</label>
              <input
                id="email"
                v-model="form.email"
                type="email"
                required
                autocomplete="email"
                class="ds-input"
                placeholder="seu@email.com"
              />
            </div>

            <div>
              <label for="password" class="ds-label">Senha</label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                required
                autocomplete="current-password"
                class="ds-input"
              />
            </div>

            <div v-if="error" class="rounded-[1.2rem] border px-4 py-3 text-sm" style="border-color: rgba(191, 58, 54, 0.24); background: var(--color-danger-soft); color: var(--color-danger);">
              {{ error }}
            </div>

            <AppButton type="submit" variant="primary" block :disabled="loading">
              {{ loading ? 'Entrando...' : 'Entrar' }}
            </AppButton>
          </form>

          <p class="mt-6 text-sm text-center" style="color: var(--color-text-muted);">
            Ainda nao tem conta?
            <NuxtLink to="/criar-conta" class="font-semibold" style="color: var(--color-brand-primary);">
              Criar conta
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
const route = useRoute()

if (auth.isAuthenticated.value) {
  await navigateTo('/dashboard')
}

const form = reactive({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await auth.login(form.email, form.password)
    const redirect = route.query.redirect as string | undefined
    await router.push(redirect || '/dashboard')
  } catch {
    error.value = 'E-mail ou senha incorretos.'
  } finally {
    loading.value = false
  }
}
</script>

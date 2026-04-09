<template>
  <div class="min-h-screen">
    <div class="mx-auto flex min-h-screen max-w-[1600px] gap-4 p-3 sm:p-4 lg:gap-6 lg:p-6">
      <aside class="hidden w-72 shrink-0 lg:flex">
        <div class="ds-surface ds-panel flex w-full flex-col">
          <div class="space-y-1 border-b pb-5" style="border-color: var(--color-border);">
            <div class="flex items-center gap-3">
              <div class="flex h-11 w-11 items-center justify-center rounded-[1rem] border text-[11px] font-bold tracking-[0.22em]" style="border-color: rgba(var(--color-brand-primary-rgb), 0.12); color: var(--color-brand-primary); background: rgba(var(--color-brand-primary-rgb), 0.05);">
                WLS
              </div>
              <div>
                <p class="text-xs font-semibold uppercase tracking-[0.24em]" style="color: var(--color-text-soft);">Manager</p>
                <h1 class="text-base font-semibold" style="color: var(--color-text);">White Label Schedule</h1>
              </div>
            </div>
          </div>

          <nav class="mt-5 flex-1 space-y-1.5">
            <NuxtLink
              v-for="item in navItems"
              :key="item.to"
              :to="item.to"
              :class="navLinkClass(item)"
            >
              <span class="mt-1 h-2.5 w-2.5 rounded-full transition" :class="navDotClass(item)" />
              <span class="min-w-0 flex-1">
                <span class="block text-sm font-semibold leading-5">{{ item.label }}</span>
                <span class="mt-0.5 block text-xs leading-5" style="color: inherit; opacity: 0.72;">{{ item.hint }}</span>
              </span>
              <span v-if="isActive(item)" class="h-8 w-1 rounded-full" style="background: rgba(var(--color-brand-primary-rgb), 0.9);" />
              <span v-else class="h-8 w-1 rounded-full bg-transparent" />
            </NuxtLink>
          </nav>

          <div class="mt-6 border-t pt-4" style="border-color: var(--color-border);">
            <div class="flex items-center gap-3 rounded-[1rem] px-2 py-2">
              <div class="flex h-10 w-10 items-center justify-center rounded-full text-xs font-semibold" style="background: rgba(var(--color-brand-primary-rgb), 0.08); color: var(--color-brand-primary);">
                {{ userInitials }}
              </div>
              <div class="min-w-0">
                <p class="truncate text-sm font-semibold" style="color: var(--color-text);">{{ auth.user.value?.name || 'Gestor' }}</p>
                <p class="text-xs" style="color: var(--color-text-soft);">Conta administrativa</p>
              </div>
            </div>

            <AppButton variant="ghost" size="sm" block @click="auth.logout">
              Sair do painel
            </AppButton>
          </div>
        </div>
      </aside>

      <div class="min-w-0 flex-1">
        <nav class="mt-4 flex gap-2 overflow-x-auto pb-1 lg:hidden">
          <NuxtLink
            v-for="item in navItems"
            :key="`mobile-${item.to}`"
            :to="item.to"
            :class="mobileNavLinkClass(item)"
            :style="isActive(item) ? activeMobileStyle : undefined"
          >
            {{ item.label }}
          </NuxtLink>
        </nav>

        <main class="mt-4 md:mt-6">
          <slot />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const auth = useAuth()
const route = useRoute()

if (import.meta.client) {
  applyBrandTheme()
}

const navItems = [
  { to: '/dashboard', label: 'Inicio', heading: 'Visao operacional', hint: 'Resumo do dia' },
  { to: '/dashboard/appointments', label: 'Agendamentos', heading: 'Agenda e status', hint: 'Fluxo e bloqueios' },
  { to: '/dashboard/professionals', label: 'Profissionais', heading: 'Equipe', hint: 'Escala e cadastro' },
  { to: '/dashboard/services', label: 'Servicos', heading: 'Catalogo', hint: 'Ofertas e duracao' },
  { to: '/dashboard/hours', label: 'Horarios', heading: 'Disponibilidade base', hint: 'Jornada da casa' },
  { to: '/dashboard/settings', label: 'Configuracoes', heading: 'Marca e integracoes', hint: 'Whitelabel e Google' },
]

const userInitials = computed(() => (auth.user.value?.name || 'G').trim().slice(0, 2).toUpperCase())
const activeMobileStyle = {
  background: 'var(--color-brand-primary)',
  color: 'var(--color-brand-on-primary)',
}

function isActive(item: { to: string }) {
  return item.to === '/dashboard' ? route.path === item.to : route.path.startsWith(item.to)
}

function navLinkClass(item: { to: string }) {
  return [
    'flex items-start gap-3 rounded-[1rem] px-3 py-3 transition',
    isActive(item)
      ? 'bg-white/70 text-[var(--color-text)] shadow-soft'
      : 'text-[var(--color-text-muted)] hover:bg-white/70 hover:text-[var(--color-text)]',
  ]
}

function navDotClass(item: { to: string }) {
  return isActive(item)
    ? 'bg-[rgb(var(--color-brand-primary-rgb))]'
    : 'bg-[rgba(93,103,123,0.35)]'
}

function mobileNavLinkClass(item: { to: string }) {
  return [
    'whitespace-nowrap rounded-full border px-4 py-2 text-sm font-semibold transition',
    isActive(item)
      ? 'border-transparent'
      : 'bg-white text-[var(--color-text-muted)]',
  ]
}
</script>

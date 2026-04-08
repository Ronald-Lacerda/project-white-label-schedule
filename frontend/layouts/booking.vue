<template>
  <div class="relative min-h-screen overflow-hidden">
    <div class="pointer-events-none absolute inset-0">
      <div class="absolute left-[-10%] top-0 h-72 w-72 rounded-full blur-3xl" style="background: rgba(var(--color-brand-secondary-rgb), 0.14);" />
      <div class="absolute right-[-8%] top-24 h-80 w-80 rounded-full blur-3xl" style="background: rgba(var(--color-brand-primary-rgb), 0.1);" />
      <div class="absolute inset-x-0 bottom-0 h-56 ds-dot-grid opacity-60" />
    </div>

    <main class="relative mx-auto max-w-6xl px-4 py-6 sm:px-6 sm:py-8 lg:px-8 lg:py-10">
      <slot />
    </main>
  </div>
</template>

<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <p class="text-sm" style="color: var(--color-text-muted);">Redirecionando...</p>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

await navigateTo('/login')
</script>

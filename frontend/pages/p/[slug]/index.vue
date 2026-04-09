<template>
  <div />
</template>

<script setup lang="ts">
definePageMeta({
  middleware: [
    () => navigateTo(`/page/${useRoute().params.slug as string}`, { redirectCode: 301 }),
  ],
})
</script>

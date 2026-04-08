export default defineNuxtConfig({
  devtools: { enabled: true },
  ssr: false,
  compatibilityDate: '2026-04-08',

  css: ['~/assets/css/main.css'],

  components: [
    {
      path: '~/components',
      pathPrefix: false,
    },
  ],

  experimental: {
    appManifest: false,
  },

  modules: ['@nuxtjs/tailwindcss'],

  typescript: {
    strict: true,
  },

  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
    },
  },

  app: {
    head: {
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      bodyAttrs: {
        class: 'antialiased',
      },
    },
  },

  vite: {
    optimizeDeps: {
      noDiscovery: true,
      include: [],
    },
  },

  nitro: {
    compressPublicAssets: true,
  },
})

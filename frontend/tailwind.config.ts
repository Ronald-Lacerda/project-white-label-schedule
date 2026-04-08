import type { Config } from 'tailwindcss'

export default {
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './app.vue',
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: 'var(--color-brand-primary)',
          secondary: 'var(--color-brand-secondary)',
          accent: 'var(--color-brand-accent)',
          surface: 'var(--color-brand-surface)',
          muted: 'var(--color-brand-surface-muted)',
        },
        app: {
          canvas: 'var(--color-canvas)',
          surface: 'var(--color-surface)',
          muted: 'var(--color-surface-muted)',
          soft: 'var(--color-surface-soft)',
          border: 'var(--color-border)',
          text: 'var(--color-text)',
          subtle: 'var(--color-text-muted)',
        },
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
        display: ['var(--font-display)'],
      },
      boxShadow: {
        soft: 'var(--shadow-soft)',
        card: 'var(--shadow-card)',
        float: 'var(--shadow-float)',
      },
      borderRadius: {
        shell: 'var(--radius-lg)',
        panel: 'var(--radius-md)',
      },
    },
  },
  plugins: [],
} satisfies Config

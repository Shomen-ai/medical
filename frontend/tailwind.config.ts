// Файл: tailwind.config.ts
// Назначение: конфигурация Tailwind CSS — фирменная палитра (primary, slate, muted), шрифт Inter и собственные тени card/card-lg.
import type { Config } from 'tailwindcss'

// Экспорт конфигурации Tailwind, типизированной через `satisfies Config`.
export default {
  content: ['./components/**/*.vue', './pages/**/*.vue', './app.vue'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#007C81', light: '#E6F5F5' },
        slate:   { DEFAULT: '#3C4F5B' },
        muted:   '#6B8290',
        border:  '#E6E6E6',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        card:    '0 4px 24px rgba(0, 110, 115, 0.12)',
        'card-lg': '0 6px 32px rgba(0, 0, 0, 0.10)',
      },
    },
  },
} satisfies Config

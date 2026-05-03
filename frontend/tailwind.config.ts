import type { Config } from 'tailwindcss'

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
    },
  },
} satisfies Config

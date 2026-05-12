// Файл: nuxt.config.ts
// Назначение: главная конфигурация Nuxt — модули (Tailwind, Pinia), runtime config с API base и публичными данными клиники, head-теги и SSR-настройки.
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  ssr: true,
  // pathPrefix: false — component names are just filenames regardless of directory depth.
  // e.g. components/booking/steps/StepSpecialty.vue → <StepSpecialty />, not <BookingStepsStepSpecialty />
  components: {
    dirs: [{ path: '~/components', pathPrefix: false }],
  },
  runtimeConfig: {
    apiBase: 'http://localhost:8080',   // overridden by NUXT_API_BASE in docker
    public: {
      clinicName: 'BeautyMed',
      clinicPhone: '+993 63 05-06-04',
      clinicAddress: 'Туркменабад, ул. Парахат 25/31',
      clinicHours: 'Пн–Сб: 9:00–20:00, Вс: выходной',
    },
  },
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap' },
      ],
      meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    },
  },
})

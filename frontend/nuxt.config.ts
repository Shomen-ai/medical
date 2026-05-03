export default defineNuxtConfig({
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
      clinicPhone: '+7 (4212) XX-XX-XX',
      clinicAddress: 'г. Хабаровск, ул. Примерная, 1',
      clinicHours: 'Пн–Сб: 9:00–20:00, Вс: выходной',
    },
  },
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    },
  },
})

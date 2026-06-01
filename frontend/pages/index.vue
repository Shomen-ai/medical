<!--
  Файл: pages/index.vue
  Назначение: главная страница сайта BeautyMed; собирает секции «герой», услуги, врачи, преимущества, отзывы и контакты, а также загружает справочники специальностей/врачей/услуг.
-->
<script setup lang="ts">
import type { Specialty, Doctor, Service } from '~/types'
import { SpecialtiesKey, DoctorsKey, ServicesKey } from '~/composables/injectionKeys'

// useApi handles base URL: server → config.apiBase (http://api:8080), client → '' (through nginx)
const { get } = useApi()

const config = useRuntimeConfig()
const { locale } = useI18n()

// SSR: parallel data fetch — fetcher runs server-side, result hydrated to client
const [
  { data: specialties },
  { data: doctors },
  { data: services },
] = await Promise.all([
  useAsyncData('specialties', () => get<Specialty[]>('/api/specialties')),
  useAsyncData('doctors',     () => get<Doctor[]>('/api/doctors')),
  useAsyncData('services',    () => get<Service[]>('/api/services')),
])

// Fallback to empty arrays on SSR error (site still renders with static content)
const safeSpecialties = computed(() => specialties.value ?? [])
const safeDoctors     = computed(() => doctors.value ?? [])
const safeServices    = computed(() => services.value ?? [])

// Provide to booking step components
provide(SpecialtiesKey, safeSpecialties)
provide(DoctorsKey, safeDoctors)
provide(ServicesKey, safeServices)

// SEO — user-visible strings switch with the current locale (ru/tk), SSR-safe via computed
const seo = computed(() =>
  locale.value === 'tk'
    ? {
        title: 'BeautyMed — Türkmenabatda gözellik we saglyk kliniki',
        description:
          'Kosmetologiýa, dermatologiýa, trihologiýa we estetiki lukmançylyk. Lukmana onlaýn ýazylyş 2 minutda.',
        ogTitle: 'BeautyMed — gözellik we saglyk kliniki',
        ogDescription: 'Türkmenabatda professional kosmetologiýa we estetiki lukmançylyk.',
      }
    : {
        title: 'BeautyMed — Клиника красоты и здоровья в Туркменабаде',
        description:
          'Косметология, дерматология, трихология и эстетическая медицина. Онлайн-запись к врачу за 2 минуты.',
        ogTitle: 'BeautyMed — Клиника красоты и здоровья',
        ogDescription: 'Профессиональная косметология и эстетическая медицина в Туркменабаде.',
      },
)

useHead({
  title: () => seo.value.title,
  meta: [
    {
      name: 'description',
      content: () => seo.value.description,
    },
    { property: 'og:title', content: () => seo.value.ogTitle },
    { property: 'og:description', content: () => seo.value.ogDescription },
    { property: 'og:image', content: '/clinic_3.png' },
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'MedicalClinic',
        name: 'BeautyMed',
        address: {
          '@type': 'PostalAddress',
          addressLocality: 'Туркменабад',
          streetAddress: 'ул. Парахат 25/31',
          addressCountry: 'TM',
        },
        telephone: config.public.clinicPhone,
        openingHours: 'Mo-Sa 09:00-20:00',
        url: 'https://beautymed.ru',
      }),
    },
  ],
})
</script>

<template>
  <main>
    <HeroSection />
    <AdvantagesSection />
    <ServicesSection :specialties="safeSpecialties" :services="safeServices" />
    <DoctorsSection :doctors="safeDoctors" />
    <ReviewsSection />
    <ContactsSection />
    <BookingModal />
  </main>
</template>

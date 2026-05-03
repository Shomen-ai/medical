<script setup lang="ts">
import type { Specialty, Doctor, Service } from '~/types'

// useApi handles base URL: server → config.apiBase (http://api:8080), client → '' (through nginx)
const { get } = useApi()

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
provide('specialties', safeSpecialties)
provide('doctors', safeDoctors)
provide('services', safeServices)

const config = useRuntimeConfig()

// SEO
useHead({
  title: 'BeautyMed — Клиника красоты и здоровья в Хабаровске',
  meta: [
    {
      name: 'description',
      content: 'Косметология, дерматология, трихология и эстетическая медицина. Онлайн-запись к врачу за 2 минуты.',
    },
    { property: 'og:title', content: 'BeautyMed — Клиника красоты и здоровья' },
    { property: 'og:description', content: 'Профессиональная косметология и эстетическая медицина в Хабаровске.' },
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
          addressLocality: 'Хабаровск',
          streetAddress: config.public.clinicAddress,
          addressCountry: 'RU',
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

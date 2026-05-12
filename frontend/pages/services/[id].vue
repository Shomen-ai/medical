<!--
  Файл: pages/services/[id].vue
  Назначение: детальная страница услуги по её ID; загружает данные услуги и связанную специальность, отображает описание, цену и кнопку записи.
-->
<script setup lang="ts">
import type { Service, Specialty } from '~/types'

const route = useRoute()
const { get } = useApi()
const booking = useBookingStore()
const { tMed } = useI18n()

const serviceId = computed(() => route.params.id as string)

const { data: service, error } = await useAsyncData(
  `service-${serviceId.value}`,
  () => get<Service>(`/api/services/${serviceId.value}`)
)

if (error.value || !service.value) {
  throw createError({ statusCode: 404, statusMessage: 'Услуга не найдена', fatal: true })
}

const { data: specialty } = await useAsyncData(
  `specialty-for-${serviceId.value}`,
  async () => {
    if (!service.value) return null
    const all = await get<Specialty[]>('/api/specialties')
    return all.find(s => s.id === service.value!.specialty_id) ?? null
  }
)

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(price)

const startBooking = () => {
  if (service.value) {
    booking.specialtyId = service.value.specialty_id
    booking.serviceId = service.value.id
    booking.servicePrice = service.value.price
    booking.finalPrice = service.value.price
  }
  booking.openModal(service.value?.specialty_id)
}

const svcName = computed(() => service.value?.name ?? '')
const svcDesc = computed(() => service.value?.description ?? '')

useHead({
  title: computed(() => svcName.value ? `${svcName.value} — BeautyMed` : 'Услуга — BeautyMed'),
  meta: [
    { name: 'description', content: svcDesc },
    { property: 'og:title', content: computed(() => `${svcName.value} — BeautyMed`) },
    { property: 'og:description', content: svcDesc },
  ],
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-12">
    <NuxtLink to="/#services" class="inline-flex items-center text-sm text-muted hover:text-primary mb-6">
      ← Все услуги
    </NuxtLink>

    <div v-if="service" class="space-y-6">
      <div>
        <div v-if="specialty" class="text-xs uppercase tracking-wide text-primary font-bold mb-2">
          {{ tMed(specialty.name) }}
        </div>
        <h1 class="text-3xl font-extrabold text-slate mb-3">{{ tMed(service.name) }}</h1>
        <div class="flex items-center gap-4 flex-wrap">
          <span class="text-xl font-bold text-primary">{{ formatPrice(service.price) }}</span>
          <span class="text-sm text-muted">Длительность: {{ service.duration_min }} мин</span>
        </div>
      </div>

      <article class="prose prose-slate max-w-none">
        <p class="text-slate leading-relaxed">{{ service.description }}</p>
      </article>

      <div class="pt-6 border-t border-border">
        <button
          type="button"
          class="bg-primary text-white px-6 py-3 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors"
          @click="startBooking"
        >
          Записаться на «{{ tMed(service.name) }}»
        </button>
      </div>
    </div>
  </div>
</template>

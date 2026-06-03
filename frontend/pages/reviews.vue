<!--
  Файл: pages/reviews.vue
  Назначение: страница всех отзывов пациентов с фильтрами по услуге и врачу + форма написания отзыва.
-->
<script setup lang="ts">
import type { ReviewItem, Doctor, Service } from '~/types'

const { t } = useI18n()
const { get } = useApi()

const doctorId = ref('')
const serviceId = ref('')
const limit = 20
const offset = ref(0)
const extra = ref<ReviewItem[]>([]) // дозагруженные «показать ещё»
const loadingMore = ref(false)
const reachedEnd = ref(false)

const { data: doctors } = await useAsyncData('reviews-doctors', () => get<Doctor[]>('/api/doctors'))
const { data: services } = await useAsyncData('reviews-services', () => get<Service[]>('/api/services'))

const buildQuery = (off: number) => {
  const p = new URLSearchParams()
  if (doctorId.value) p.set('doctor_id', doctorId.value)
  if (serviceId.value) p.set('service_id', serviceId.value)
  p.set('limit', String(limit)); p.set('offset', String(off))
  return p.toString()
}

// Первая страница: SSR (гидрируется без клиентского рефетча) + перезапрос при смене фильтра.
const { data: firstPage, refresh } = await useAsyncData(
  'reviews-list',
  () => get<ReviewItem[]>(`/api/reviews?${buildQuery(0)}`),
  { watch: [doctorId, serviceId], default: () => [] as ReviewItem[] },
)

// При смене первой страницы (фильтр) сбрасываем дозагруженное и флаг конца.
watch(firstPage, (p) => {
  extra.value = []
  offset.value = 0
  reachedEnd.value = (p?.length ?? 0) < limit
}, { immediate: true })

const reviews = computed(() => [...(firstPage.value ?? []), ...extra.value])

const loadMore = async () => {
  loadingMore.value = true
  try {
    offset.value += limit
    const batch = (await get<ReviewItem[]>(`/api/reviews?${buildQuery(offset.value)}`)) ?? []
    extra.value = [...extra.value, ...batch]
    if (batch.length < limit) reachedEnd.value = true
  } catch {
    reachedEnd.value = true
  } finally {
    loadingMore.value = false
  }
}

const formatDate = (iso: string) => {
  const d = new Date(iso)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}.${p(d.getMonth() + 1)}.${d.getFullYear()}`
}

useHead({ title: () => `${t('reviewsTitle')} — BeautyMed` })
</script>

<template>
  <main class="min-h-screen bg-[#F0FAFB] py-10">
    <div class="max-w-4xl mx-auto px-4 sm:px-8">
      <NuxtLink to="/" class="text-sm text-primary">{{ t('back') }}</NuxtLink>
      <h1 class="text-2xl sm:text-3xl font-extrabold text-slate mt-2 mb-6">{{ t('reviewsTitle') }}</h1>

      <ReviewForm class="mb-8" @created="() => refresh()" />

      <div class="flex flex-col sm:flex-row gap-3 mb-6">
        <select v-model="serviceId" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary bg-white">
          <option value="">{{ t('reviewsFilterAllServices') }}</option>
          <option v-for="s in services ?? []" :key="s.id" :value="s.id">{{ s.name }}</option>
        </select>
        <select v-model="doctorId" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary bg-white">
          <option value="">{{ t('reviewsFilterAllDoctors') }}</option>
          <option v-for="d in doctors ?? []" :key="d.id" :value="d.id">{{ d.full_name }}</option>
        </select>
      </div>

      <div v-if="reviews.length" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div v-for="r in reviews" :key="r.id" class="bg-white rounded-2xl shadow-card p-6 flex flex-col">
          <StarRating :model-value="r.rating" readonly size="text-sm" />
          <p class="text-sm text-slate leading-relaxed my-4 flex-1">{{ r.text }}</p>
          <div class="pt-3 border-t border-border">
            <div class="text-xs font-bold text-slate">{{ t('reviewsAnonymous') }}</div>
            <div class="flex items-center justify-between mt-0.5">
              <span class="text-[11px] text-muted truncate">{{ r.service_name }} · {{ r.doctor_name }}</span>
              <span class="text-[11px] text-muted shrink-0 ml-2">{{ formatDate(r.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
      <p v-else class="text-sm text-muted text-center py-10">{{ t('reviewsEmpty') }}</p>

      <div v-if="reviews.length && !reachedEnd" class="text-center mt-6">
        <button type="button" :disabled="loadingMore" class="text-sm font-semibold text-primary disabled:opacity-50" @click="loadMore">
          {{ loadingMore ? t('loading') : t('reviewsLoadMore') }}
        </button>
      </div>
    </div>
  </main>
</template>

<!--
  Файл: components/sections/ReviewsSection.vue
  Назначение: секция главной — карусель последних отзывов пациентов (до 10) из API + кнопка «Все отзывы».
-->
<script setup lang="ts">
import type { ReviewItem } from '~/types'

const { t } = useI18n()
const { get } = useApi()

const { data: reviews } = await useAsyncData('home-reviews', () =>
  get<ReviewItem[]>('/api/reviews?limit=10'))

const items = computed(() => reviews.value ?? [])

const formatDate = (iso: string) => {
  const d = new Date(iso)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}.${p(d.getMonth() + 1)}.${d.getFullYear()}`
}

const track = ref<HTMLElement | null>(null)
const scrollByCard = (dir: number) => track.value?.scrollBy({ left: dir * 300, behavior: 'smooth' })
</script>

<template>
  <section class="py-10 sm:py-14 bg-[#F0FAFB]">
    <div class="max-w-5xl mx-auto px-4 sm:px-8">
      <h2 class="text-2xl sm:text-3xl font-extrabold text-slate mb-6 sm:mb-10 text-center">{{ t('reviewsTitle') }}</h2>

      <div v-if="items.length" class="relative">
        <button
          type="button" :aria-label="t('back')"
          class="hidden sm:flex absolute -left-3 top-1/2 -translate-y-1/2 z-10 w-9 h-9 items-center justify-center rounded-full bg-white shadow-card text-slate hover:text-primary"
          @click="scrollByCard(-1)"
        >‹</button>
        <button
          type="button" :aria-label="t('next')"
          class="hidden sm:flex absolute -right-3 top-1/2 -translate-y-1/2 z-10 w-9 h-9 items-center justify-center rounded-full bg-white shadow-card text-slate hover:text-primary"
          @click="scrollByCard(1)"
        >›</button>

        <div ref="track" class="flex gap-4 overflow-x-auto snap-x snap-mandatory scroll-smooth pb-2 px-1" style="scrollbar-width:none">
          <div
            v-for="r in items"
            :key="r.id"
            class="snap-start shrink-0 w-[260px] sm:w-[280px] bg-white rounded-2xl shadow-card p-6 hover:shadow-card-lg transition-shadow flex flex-col"
          >
            <div class="text-5xl text-primary/20 font-serif leading-none mb-2 select-none">"</div>
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
      </div>

      <p v-else class="text-sm text-muted text-center">{{ t('reviewsEmpty') }}</p>

      <div class="text-center mt-6">
        <NuxtLink
          to="/reviews"
          class="inline-block text-white px-6 py-2.5 rounded-lg text-sm font-semibold hover:opacity-90 transition-opacity"
          style="background: linear-gradient(135deg, #005A5F, #00959D)"
        >
          {{ t('reviewsAll') }}
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

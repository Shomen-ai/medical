<script setup lang="ts">
import type { Specialty, Service } from '~/types'
import { getSpecialtyMeta } from '~/data/specialtyMeta'

const props = defineProps<{
  specialties: Specialty[]
  services: Service[]
}>()

const booking = useBookingStore()
const activeSpecialty = ref<string | null>(props.specialties[0]?.id ?? null)

const activeServices = computed(() =>
  props.services.filter(s => s.specialty_id === activeSpecialty.value)
)

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(price)
</script>

<template>
  <section id="services" class="py-14 bg-[#F0FAFB]">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-3xl font-extrabold text-slate mb-10 text-center">Наши услуги</h2>

      <!-- Specialty tabs -->
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="sp in specialties"
          :key="sp.id"
          type="button"
          class="px-4 py-2 rounded-full text-sm font-semibold transition-all"
          :class="activeSpecialty === sp.id
            ? 'text-white shadow-md'
            : 'bg-white border border-border text-muted hover:border-primary hover:text-primary'"
          :style="activeSpecialty === sp.id
            ? 'background: linear-gradient(135deg, #005A5F, #00959D)'
            : ''"
          @click="activeSpecialty = sp.id"
        >
          {{ getSpecialtyMeta(sp.name).icon }} {{ sp.name }}
        </button>
      </div>

      <!-- Service list -->
      <div class="bg-white rounded-2xl shadow-card overflow-hidden">
        <div
          v-for="(svc, i) in activeServices"
          :key="svc.id"
          class="flex items-center justify-between px-5 py-4 hover:bg-[#F0FAFB] transition-colors"
          :class="i < activeServices.length - 1 ? 'border-b border-border' : ''"
        >
          <div class="flex-1 pr-4">
            <NuxtLink :to="`/services/${svc.id}`" class="text-sm font-medium text-slate hover:text-primary transition-colors">
              {{ svc.name }}
            </NuxtLink>
            <div class="text-xs text-muted mt-0.5">{{ svc.duration_min }} мин</div>
          </div>
          <div class="flex items-center gap-4">
            <span class="text-sm font-semibold text-primary">{{ formatPrice(svc.price) }}</span>
            <NuxtLink
              :to="`/services/${svc.id}`"
              class="text-xs font-semibold text-muted hover:text-primary transition-colors hidden sm:block"
            >
              Подробнее
            </NuxtLink>
            <button
              type="button"
              class="text-xs font-semibold text-primary border border-primary px-3 py-1.5 rounded-lg hover:bg-primary hover:text-white transition-colors whitespace-nowrap"
              @click="booking.openModal(activeSpecialty ?? undefined)"
            >
              Записаться
            </button>
          </div>
        </div>
        <div v-if="activeServices.length === 0" class="px-5 py-8 text-center text-muted text-sm">
          Услуги не найдены
        </div>
      </div>
    </div>
  </section>
</template>

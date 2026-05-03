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
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(price)
</script>

<template>
  <section id="services" class="py-12 bg-gray-50">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Наши услуги</h2>

      <!-- Specialty tabs -->
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="sp in specialties"
          :key="sp.id"
          class="px-4 py-2 rounded-full text-sm font-medium transition-colors"
          :class="activeSpecialty === sp.id
            ? 'bg-primary text-white'
            : 'bg-white border border-border text-muted hover:border-primary hover:text-primary'"
          @click="activeSpecialty = sp.id"
        >
          {{ getSpecialtyMeta(sp.name).icon }} {{ sp.name }}
        </button>
      </div>

      <!-- Service list -->
      <div class="bg-white rounded-xl border border-border overflow-hidden">
        <div
          v-for="(svc, i) in activeServices"
          :key="svc.id"
          class="flex items-center justify-between px-5 py-4"
          :class="i < activeServices.length - 1 ? 'border-b border-border' : ''"
        >
          <div class="flex-1 pr-4">
            <div class="text-sm font-medium text-slate">{{ svc.name }}</div>
            <div class="text-xs text-muted mt-0.5">{{ svc.duration_min }} мин</div>
          </div>
          <div class="flex items-center gap-4">
            <span class="text-sm font-semibold text-primary">{{ formatPrice(svc.price) }}</span>
            <button
              class="text-xs font-semibold text-primary border border-primary px-3 py-1.5 rounded-lg hover:bg-primary/5 transition-colors whitespace-nowrap"
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

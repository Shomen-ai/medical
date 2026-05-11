<script setup lang="ts">
import { getSpecialtyMeta } from '~/data/specialtyMeta'
import { SpecialtiesKey, ServicesKey } from '~/composables/injectionKeys'

const booking = useBookingStore()
const specialties = inject(SpecialtiesKey, computed(() => []))
const services = inject(ServicesKey, computed(() => []))

const servicesForSelected = computed(() =>
  services.value.filter(s => s.specialty_id === booking.specialtyId)
)

const selectSpecialty = (id: string) => {
  booking.specialtyId = id
  booking.serviceId = null
  booking.servicePrice = 0
  booking.finalPrice = 0
  booking.promoValid = null
  booking.promoDiscountPct = 0
}

const selectService = (svc: { id: string; price: number }) => {
  booking.serviceId = svc.id
  booking.servicePrice = svc.price
  booking.finalPrice = svc.price
  booking.promoValid = null
  booking.promoDiscountPct = 0
}

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(price)
</script>

<template>
  <div class="flex flex-col sm:flex-row gap-4 sm:min-h-[280px]">
    <!-- Left: specialty list -->
    <div class="w-full sm:w-1/2 flex flex-col gap-2">
      <div class="text-xs font-semibold text-slate mb-1">Специальность</div>
      <button
        v-for="sp in specialties"
        :key="sp.id"
        type="button"
        class="flex items-center justify-between px-3 sm:px-4 py-2.5 sm:py-3 rounded-xl border text-sm transition-colors text-left"
        :class="booking.specialtyId === sp.id
          ? 'border-primary bg-primary/5 text-primary font-semibold'
          : 'border-border text-slate hover:border-primary'"
        @click="selectSpecialty(sp.id)"
      >
        <span>{{ getSpecialtyMeta(sp.name).icon }} {{ sp.name }}</span>
        <span class="ml-2 flex-shrink-0">{{ booking.specialtyId === sp.id ? '✓' : '›' }}</span>
      </button>
    </div>

    <!-- Right: service list (appears when specialty selected) -->
    <div class="w-full sm:w-1/2 flex flex-col">
      <Transition name="fade-right">
        <div v-if="booking.specialtyId && servicesForSelected.length" class="flex flex-col gap-1.5">
          <div class="text-xs font-semibold text-slate mb-1">Услуга</div>
          <button
            v-for="svc in servicesForSelected"
            :key="svc.id"
            type="button"
            class="flex items-center justify-between px-3 sm:px-4 py-2.5 rounded-lg border text-sm transition-colors text-left gap-2"
            :class="booking.serviceId === svc.id
              ? 'border-primary bg-primary/5 text-primary font-semibold'
              : 'border-border text-slate hover:border-primary'"
            @click="selectService(svc)"
          >
            <span class="text-left leading-snug">{{ svc.name }}</span>
            <span class="whitespace-nowrap text-xs flex-shrink-0">{{ formatPrice(svc.price) }}</span>
          </button>
        </div>
        <div v-else-if="!booking.specialtyId" class="hidden sm:flex items-center justify-center h-full text-sm text-muted text-center pt-8">
          Выберите<br>специальность
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.fade-right-enter-active { transition: all 0.2s ease; }
.fade-right-enter-from { opacity: 0; transform: translateX(-8px); }
</style>

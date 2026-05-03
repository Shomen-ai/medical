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
}

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(price)
</script>

<template>
  <div>
    <!-- Specialty options -->
    <div class="text-xs font-semibold text-slate mb-2">Выберите специальность</div>
    <div class="flex flex-col gap-2 mb-4">
      <button
        v-for="sp in specialties"
        :key="sp.id"
        class="flex items-center justify-between px-4 py-3 rounded-xl border text-sm transition-colors"
        :class="booking.specialtyId === sp.id
          ? 'border-primary bg-primary/5 text-primary font-semibold'
          : 'border-border text-slate hover:border-primary'"
        @click="selectSpecialty(sp.id)"
      >
        <span>{{ getSpecialtyMeta(sp.name).icon }} {{ sp.name }}</span>
        <span>{{ booking.specialtyId === sp.id ? '✓' : '›' }}</span>
      </button>
    </div>

    <!-- Service options (appear after specialty selected) -->
    <Transition name="slide-down">
      <div v-if="booking.specialtyId && servicesForSelected.length">
        <div class="text-xs font-semibold text-slate mb-2">Выберите услугу</div>
        <div class="flex flex-col gap-1.5">
          <button
            v-for="svc in servicesForSelected"
            :key="svc.id"
            class="flex items-center justify-between px-4 py-2.5 rounded-lg border text-sm transition-colors"
            :class="booking.serviceId === svc.id
              ? 'border-primary bg-primary/5 text-primary font-semibold'
              : 'border-border text-slate hover:border-primary'"
            @click="booking.serviceId = svc.id"
          >
            <span class="text-left">{{ svc.name }}</span>
            <span class="ml-3 whitespace-nowrap text-xs">{{ formatPrice(svc.price) }}</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.slide-down-enter-active { transition: all 0.2s ease; }
.slide-down-enter-from { opacity: 0; transform: translateY(-6px); }
</style>

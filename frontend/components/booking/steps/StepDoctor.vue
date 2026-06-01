<!--
  Файл: components/booking/steps/StepDoctor.vue
  Назначение: второй шаг мастера записи — выбор врача из списка специалистов, доступных по выбранной ранее специальности.
-->
<script setup lang="ts">
import { DoctorsKey } from '~/composables/injectionKeys'

const booking = useBookingStore()
const { t } = useI18n()
const doctors = inject(DoctorsKey, computed(() => []))

const filteredDoctors = computed(() =>
  doctors.value.filter(d => d.specialty_id === booking.specialtyId)
)
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-2">{{ t('bkSelectDoctor') }}</div>
    <div class="flex flex-col gap-2">
      <button
        v-for="doc in filteredDoctors"
        :key="doc.id"
        class="flex items-center gap-3 px-4 py-3 rounded-xl border text-left transition-colors"
        :class="booking.doctorId === doc.id
          ? 'border-primary bg-primary/5'
          : 'border-border hover:border-primary'"
        @click="booking.doctorId = doc.id"
      >
        <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0 overflow-hidden">
          <img v-if="doc.photo_url" :src="doc.photo_url" :alt="doc.full_name" class="w-full h-full object-cover" />
          <span v-else class="text-lg">👩‍⚕️</span>
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium text-slate truncate">{{ doc.full_name }}</div>
          <div class="text-xs text-muted">{{ t('doctorsExperience', { n: doc.experience_years }) }}</div>
        </div>
        <span v-if="booking.doctorId === doc.id" class="text-primary text-sm">✓</span>
      </button>
    </div>
    <div v-if="filteredDoctors.length === 0" class="text-sm text-muted text-center py-6">
      {{ t('bkNoDoctors') }}
    </div>
  </div>
</template>

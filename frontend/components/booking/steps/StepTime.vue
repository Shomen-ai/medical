<!--
  Файл: components/booking/steps/StepTime.vue
  Назначение: четвёртый шаг мастера записи — выбор свободного временного слота у врача на выбранную дату.
-->
<script setup lang="ts">
import type { TimeSlot } from '~/types'

const booking = useBookingStore()
const { get } = useApi()
const { t } = useI18n()

const slots = ref<TimeSlot[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const fetchSlots = async () => {
  if (!booking.doctorId || !booking.serviceId || !booking.date) return
  loading.value = true
  error.value = null
  try {
    slots.value = await get<TimeSlot[]>(
      `/api/doctors/${booking.doctorId}/slots?service_id=${booking.serviceId}&date=${booking.date}`
    )
  } catch {
    error.value = t('bkSlotsError')
    slots.value = []
  } finally {
    loading.value = false
  }
}

// Следим и за услугой: при её смене длительность/слоты меняются — иначе время «не подгружается».
watch([() => booking.date, () => booking.doctorId, () => booking.serviceId], fetchSlots, { immediate: true })

// booking.date is "YYYY-MM-DD" — format to Russian locale without timezone issues
const formattedDate = computed(() => {
  if (!booking.date) return ''
  const [y, m, d] = booking.date.split('-').map(Number)
  const months = t('monthsList').split(',')
  return `${d} ${months[m - 1]} ${y}`
})
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-3">
      {{ t('bkSelectTime', { date: formattedDate }) }}
    </div>

    <div v-if="loading" class="text-sm text-muted text-center py-6">{{ t('bkLoadingSchedule') }}</div>

    <div v-else-if="error" class="text-center py-6">
      <div class="text-sm text-red-500 mb-3">{{ error }}</div>
      <button type="button" class="text-xs text-primary underline" @click="fetchSlots">{{ t('bkRetry') }}</button>
    </div>

    <div v-else-if="slots.length === 0" class="text-sm text-muted text-center py-6">
      {{ t('bkNoSlots') }}
    </div>

    <div v-else class="grid grid-cols-3 gap-2">
      <button
        v-for="slot in slots"
        :key="slot.starts_at"
        type="button"
        class="py-2.5 rounded-lg border text-sm font-medium transition-colors"
        :class="booking.timeSlot === slot.starts_at
          ? 'bg-primary text-white border-primary'
          : 'border-border text-slate hover:border-primary'"
        @click="booking.timeSlot = slot.starts_at"
      >
        {{ slot.starts_at }}
      </button>
    </div>
  </div>
</template>

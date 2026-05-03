<script setup lang="ts">
import type { TimeSlot } from '~/types'

const booking = useBookingStore()
const { get } = useApi()

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
    error.value = 'Не удалось загрузить слоты'
    slots.value = []
  } finally {
    loading.value = false
  }
}

watch([() => booking.date, () => booking.doctorId], fetchSlots, { immediate: true })

const formatTime = (iso: string) =>
  new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' })
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-3">
      Выберите время — {{ booking.date }}
    </div>

    <div v-if="loading" class="text-sm text-muted text-center py-6">Загружаем расписание...</div>

    <div v-else-if="error" class="text-center py-6">
      <div class="text-sm text-red-500 mb-3">{{ error }}</div>
      <button class="text-xs text-primary underline" @click="fetchSlots">Попробовать снова</button>
    </div>

    <div v-else-if="slots.length === 0" class="text-sm text-muted text-center py-6">
      Нет доступных слотов на эту дату
    </div>

    <div v-else class="grid grid-cols-3 gap-2">
      <button
        v-for="slot in slots"
        :key="slot.starts_at"
        class="py-2 rounded-lg border text-sm font-medium transition-colors"
        :class="booking.timeSlot === slot.starts_at
          ? 'bg-primary text-white border-primary'
          : 'border-border text-slate hover:border-primary'"
        @click="booking.timeSlot = slot.starts_at"
      >
        {{ formatTime(slot.starts_at) }}
      </button>
    </div>
  </div>
</template>

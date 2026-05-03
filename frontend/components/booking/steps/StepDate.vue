<script setup lang="ts">
const booking = useBookingStore()

const today = new Date()
today.setHours(0, 0, 0, 0)

const viewDate = ref(new Date(today.getFullYear(), today.getMonth(), 1))

const monthLabel = computed(() =>
  viewDate.value.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })
)

const prevMonth = () => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() - 1)
  if (d >= new Date(today.getFullYear(), today.getMonth(), 1)) viewDate.value = d
}

const nextMonth = () => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() + 1)
  viewDate.value = d
}

const calendarDays = computed(() => {
  const year = viewDate.value.getFullYear()
  const month = viewDate.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)

  // Start grid on Monday (ISO week)
  const startOffset = (firstDay.getDay() + 6) % 7
  const days: { date: Date | null; disabled: boolean }[] = []

  for (let i = 0; i < startOffset; i++) days.push({ date: null, disabled: true })

  for (let d = 1; d <= lastDay.getDate(); d++) {
    const date = new Date(year, month, d)
    const isSunday = date.getDay() === 0
    const isPast = date < today
    days.push({ date, disabled: isSunday || isPast })
  }

  return days
})

const toISODate = (d: Date) => d.toISOString().slice(0, 10)

const selectDate = (d: Date) => {
  booking.date = toISODate(d)
  booking.timeSlot = null  // reset slot when date changes
}

const isSelected = (d: Date) => booking.date === toISODate(d)
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-3">Выберите дату</div>

    <!-- Month navigation -->
    <div class="flex items-center justify-between mb-3">
      <button class="text-muted hover:text-primary p-1 transition-colors" @click="prevMonth">‹</button>
      <span class="text-sm font-semibold text-slate capitalize">{{ monthLabel }}</span>
      <button class="text-muted hover:text-primary p-1 transition-colors" @click="nextMonth">›</button>
    </div>

    <!-- Day-of-week headers -->
    <div class="grid grid-cols-7 text-center mb-1">
      <div v-for="d in ['Пн','Вт','Ср','Чт','Пт','Сб','Вс']" :key="d" class="text-[10px] text-muted py-1">
        {{ d }}
      </div>
    </div>

    <!-- Calendar grid -->
    <div class="grid grid-cols-7 gap-0.5">
      <div v-for="(cell, i) in calendarDays" :key="i" class="aspect-square flex items-center justify-center">
        <button
          v-if="cell.date"
          class="w-8 h-8 rounded-full text-xs font-medium transition-colors"
          :class="{
            'bg-primary text-white': isSelected(cell.date),
            'text-muted cursor-not-allowed': cell.disabled,
            'text-slate hover:bg-primary/10': !cell.disabled && !isSelected(cell.date),
          }"
          :disabled="cell.disabled"
          @click="!cell.disabled && selectDate(cell.date!)"
        >
          {{ cell.date.getDate() }}
        </button>
      </div>
    </div>
  </div>
</template>

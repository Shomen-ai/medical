<!--
  Файл: components/DateField.vue
  Назначение: локализованный выбор даты (замена нативного <input type="date">, который
  браузер не переводит на туркменский). Поле + всплывающий календарь с месяцами/днями
  из i18n (monthsList, docWeekdays). Поддерживает v-model, :min и :max (YYYY-MM-DD).
-->
<script setup lang="ts">
const props = withDefaults(defineProps<{
  modelValue: string
  min?: string
  max?: string
  placeholder?: string
}>(), { modelValue: '', min: '', max: '', placeholder: '' })

const emit = defineEmits<{ 'update:modelValue': [v: string] }>()
const { t } = useI18n()

const open = ref(false)

const pad = (n: number) => String(n).padStart(2, '0')
const toISO = (d: Date) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
const parseBound = (s?: string) => (s ? new Date(s + 'T00:00:00') : null)

// Отображение в поле: DD.MM.YYYY (или плейсхолдер).
const display = computed(() => {
  if (!props.modelValue) return ''
  const [y, m, d] = props.modelValue.split('-')
  return `${d}.${m}.${y}`
})

const initMonth = () => {
  if (props.modelValue) {
    const [y, m] = props.modelValue.split('-').map(Number)
    return new Date(y, m - 1, 1)
  }
  // Пусто: открываем на месяце/годе верхней границы (для даты рождения это 2008 = сегодня−18),
  // чтобы не листать от текущего года; иначе — текущий месяц.
  if (props.max) {
    const [y, m] = props.max.split('-').map(Number)
    return new Date(y, m - 1, 1)
  }
  const now = new Date()
  return new Date(now.getFullYear(), now.getMonth(), 1)
}
const viewDate = ref(initMonth())
watch(() => props.modelValue, () => { viewDate.value = initMonth() })
watch(open, (v) => { if (v) viewDate.value = initMonth() })

const monthOptions = computed(() => t('monthsList').split(','))

// Диапазон годов для быстрого выбора: по min/max, иначе ±разумно вокруг текущего.
const yearOptions = computed(() => {
  const cur = new Date().getFullYear()
  const minY = props.min ? Number(props.min.slice(0, 4)) : cur - 100
  const maxY = props.max ? Number(props.max.slice(0, 4)) : cur + 10
  const years: number[] = []
  for (let y = maxY; y >= minY; y--) years.push(y)
  return years
})

const setMonth = (m: number) => { const d = new Date(viewDate.value); d.setDate(1); d.setMonth(m); viewDate.value = d }
const setYear = (y: number) => { const d = new Date(viewDate.value); d.setDate(1); d.setFullYear(y); viewDate.value = d }
const prevMonth = () => { const d = new Date(viewDate.value); d.setMonth(d.getMonth() - 1); viewDate.value = d }
const nextMonth = () => { const d = new Date(viewDate.value); d.setMonth(d.getMonth() + 1); viewDate.value = d }

const calendarDays = computed(() => {
  const year = viewDate.value.getFullYear()
  const month = viewDate.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startOffset = (firstDay.getDay() + 6) % 7 // понедельник первым
  const min = parseBound(props.min)
  const max = parseBound(props.max)
  const days: { date: Date | null; disabled: boolean }[] = []
  for (let i = 0; i < startOffset; i++) days.push({ date: null, disabled: true })
  for (let d = 1; d <= lastDay.getDate(); d++) {
    const date = new Date(year, month, d)
    const disabled = (min !== null && date < min) || (max !== null && date > max)
    days.push({ date, disabled })
  }
  return days
})

const isSelected = (d: Date) => props.modelValue === toISO(d)

const selectDay = (d: Date) => {
  emit('update:modelValue', toISO(d))
  open.value = false
}
const clear = () => { emit('update:modelValue', ''); open.value = false }
</script>

<template>
  <div class="relative">
    <!-- Поле -->
    <button
      type="button"
      class="w-full flex items-center justify-between border border-border rounded-lg px-3 py-2 text-sm text-left outline-none focus:border-primary"
      :class="display ? 'text-slate' : 'text-gray-400'"
      @click="open = !open"
    >
      <span>{{ display || placeholder || t('dateFieldPlaceholder') }}</span>
      <span class="text-muted">📅</span>
    </button>

    <!-- Календарь -->
    <div v-if="open">
      <!-- фон для закрытия по клику вне -->
      <div class="fixed inset-0 z-40" @click="open = false" />
      <div class="absolute z-50 mt-1 w-64 bg-white border border-border rounded-xl shadow-xl p-3">
        <div class="flex items-center gap-1 mb-2">
          <button type="button" class="text-muted hover:text-primary px-1 text-lg leading-none" @click="prevMonth">‹</button>
          <select
            :value="viewDate.getMonth()"
            class="flex-1 text-sm font-semibold text-slate bg-transparent outline-none cursor-pointer rounded px-1 py-0.5 hover:bg-gray-50"
            @change="setMonth(Number(($event.target as HTMLSelectElement).value))"
          >
            <option v-for="(m, i) in monthOptions" :key="i" :value="i">{{ m }}</option>
          </select>
          <select
            :value="viewDate.getFullYear()"
            class="text-sm font-semibold text-slate bg-transparent outline-none cursor-pointer rounded px-1 py-0.5 hover:bg-gray-50"
            @change="setYear(Number(($event.target as HTMLSelectElement).value))"
          >
            <option v-for="y in yearOptions" :key="y" :value="y">{{ y }}</option>
          </select>
          <button type="button" class="text-muted hover:text-primary px-1 text-lg leading-none" @click="nextMonth">›</button>
        </div>
        <div class="grid grid-cols-7 text-center mb-1">
          <div v-for="d in t('docWeekdays').split(',')" :key="d" class="text-[10px] text-muted py-1">{{ d }}</div>
        </div>
        <div class="grid grid-cols-7 gap-0.5">
          <div v-for="(cell, i) in calendarDays" :key="i" class="aspect-square flex items-center justify-center">
            <button
              v-if="cell.date"
              type="button"
              class="w-8 h-8 rounded-full text-sm font-medium transition-colors"
              :class="{
                'bg-primary text-white': isSelected(cell.date),
                'text-gray-300 cursor-not-allowed': cell.disabled,
                'text-slate hover:bg-primary/10': !cell.disabled && !isSelected(cell.date),
              }"
              :disabled="cell.disabled"
              @click="!cell.disabled && selectDay(cell.date!)"
            >
              {{ cell.date.getDate() }}
            </button>
          </div>
        </div>
        <div v-if="modelValue" class="text-center mt-2">
          <button type="button" class="text-xs text-muted hover:text-primary" @click="clear">{{ t('dateFieldClear') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

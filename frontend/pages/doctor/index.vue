<script setup lang="ts">
definePageMeta({ layout: 'staff' })

const auth = useAuthStore()
const router = useRouter()
const { get } = useApi()

onMounted(() => {
  auth.init()
  if (!auth.isDoctor) router.replace('/')
})

// ── Types ────────────────────────────────────────────────────────────
interface ScheduleCell {
  work_date: string; start_time: string; end_time: string
  is_day_off: boolean; has_appointments: boolean
}
interface Appointment {
  id: string; patient_name: string; patient_phone: string
  service_name: string; starts_at: string; status: string
}
interface DoctorStats {
  appointments_this_month: number
  unique_patients: number
  filled_records_pct: number
}

// ── Month navigation ─────────────────────────────────────────────────
const today = new Date()
const viewYear = ref(today.getFullYear())
const viewMonth = ref(today.getMonth() + 1)

const monthLabel = computed(() => {
  const d = new Date(viewYear.value, viewMonth.value - 1, 1)
  return d.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })
})

const prevMonth = () => {
  if (viewMonth.value === 1) { viewMonth.value = 12; viewYear.value-- }
  else viewMonth.value--
}
const nextMonth = () => {
  if (viewMonth.value === 12) { viewMonth.value = 1; viewYear.value++ }
  else viewMonth.value++
}

// ── Stats ────────────────────────────────────────────────────────────
const token = computed(() => auth.token ?? undefined)

const { data: stats } = await useAsyncData('doctor-stats', () =>
  get<DoctorStats>('/api/doctor/stats', token.value), { server: false })

// ── Schedule for month ───────────────────────────────────────────────
const { data: schedule, refresh: refreshSchedule } = await useAsyncData(
  'doctor-schedule',
  () => get<ScheduleCell[]>(`/api/doctor/schedule?year=${viewYear.value}&month=${viewMonth.value}`, token.value),
  { server: false }
)

watch([viewYear, viewMonth], () => refreshSchedule())

// All scheduled days sorted ascending
const workDays = computed(() =>
  (schedule.value ?? []).sort((a, b) => a.work_date.localeCompare(b.work_date))
)

// ── Appointments per selected day ─────────────────────────────────────
const selectedDate = ref('')
const dayAppointments = ref<Appointment[]>([])
const loadingDay = ref(false)

// Default to today if it's a work day
onMounted(() => {
  const todayStr = (() => {
    const y = today.getFullYear()
    const m = String(today.getMonth() + 1).padStart(2, '0')
    const d = String(today.getDate()).padStart(2, '0')
    return `${y}-${m}-${d}`
  })()
  selectedDate.value = todayStr
  loadDay(todayStr)
})

const loadDay = async (date: string) => {
  selectedDate.value = date
  loadingDay.value = true
  try {
    dayAppointments.value = (await get<Appointment[]>(
      `/api/doctor/appointments?date=${date}`, token.value
    )) ?? []
  } catch {
    dayAppointments.value = []
  } finally {
    loadingDay.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────────────────
const formatDate = (dateStr: string) => {
  const [y, m, d] = dateStr.split('-').map(Number)
  return new Date(y, m - 1, d).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', weekday: 'short' })
}

const formatTime = (iso: string) => {
  const d = new Date(iso)
  return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}

const isToday = (dateStr: string) => {
  const t = today
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const d = String(t.getDate()).padStart(2, '0')
  return dateStr === `${y}-${m}-${d}`
}

const statusLabel: Record<string, string> = {
  scheduled: 'Запланировано', completed: 'Завершено',
  cancelled: 'Отменено', rescheduled: 'Перенесено',
}

useHead({ title: 'Кабинет врача — BeautyMed' })
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-5xl mx-auto px-4 py-8 space-y-6">

      <!-- Stats -->
      <div class="grid grid-cols-3 gap-4">
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
          <div class="text-2xl font-extrabold text-primary">{{ stats?.appointments_this_month ?? '—' }}</div>
          <div class="text-xs text-gray-500 mt-1">Приёмов за месяц</div>
        </div>
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
          <div class="text-2xl font-extrabold text-primary">{{ stats?.unique_patients ?? '—' }}</div>
          <div class="text-xs text-gray-500 mt-1">Уникальных пациентов</div>
        </div>
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
          <div class="text-2xl font-extrabold text-primary">{{ stats?.filled_records_pct != null ? Math.round(stats.filled_records_pct) + '%' : '—' }}</div>
          <div class="text-xs text-gray-500 mt-1">Заполнено записей</div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">

        <!-- Left: month work days list -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <!-- Month nav -->
          <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100">
            <button type="button" class="text-gray-400 hover:text-primary p-1" @click="prevMonth">‹</button>
            <span class="text-sm font-semibold text-slate capitalize">{{ monthLabel }}</span>
            <button type="button" class="text-gray-400 hover:text-primary p-1" @click="nextMonth">›</button>
          </div>

          <div v-if="!workDays.length" class="text-center py-8 text-gray-400 text-sm">Нет записей в расписании</div>

          <div v-else class="divide-y divide-gray-50">
            <button
              v-for="cell in workDays"
              :key="cell.work_date"
              type="button"
              class="w-full text-left px-4 py-3 flex items-center justify-between transition-colors"
              :class="[
                selectedDate === cell.work_date
                  ? 'bg-primary/5 border-l-2 border-primary'
                  : 'hover:bg-gray-50 border-l-2 border-transparent',
                cell.is_day_off ? 'opacity-50' : ''
              ]"
              @click="loadDay(cell.work_date)"
            >
              <div>
                <div class="text-sm font-medium" :class="cell.is_day_off ? 'text-gray-400' : 'text-slate'">
                  {{ formatDate(cell.work_date) }}
                  <span v-if="isToday(cell.work_date)" class="ml-1 text-xs text-primary font-semibold">(сегодня)</span>
                  <span v-if="cell.is_day_off" class="ml-1 text-xs text-gray-400">(выходной)</span>
                </div>
                <div v-if="!cell.is_day_off" class="text-xs text-gray-400">{{ cell.start_time.slice(0, 5) }} – {{ cell.end_time.slice(0, 5) }}</div>
              </div>
            </button>
          </div>
        </div>

        <!-- Right: appointments for selected day -->
        <div class="md:col-span-2 bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div class="px-5 py-4 border-b border-gray-100">
            <div class="font-semibold text-slate">
              {{ selectedDate ? formatDate(selectedDate) : 'Выберите дату' }}
            </div>
          </div>

          <div v-if="loadingDay" class="text-center py-12 text-gray-400">Загружаем...</div>

          <div v-else-if="!dayAppointments.length" class="text-center py-12 text-gray-400">
            Нет приёмов на этот день
          </div>

          <div v-else class="divide-y divide-gray-50">
            <div
              v-for="apt in dayAppointments"
              :key="apt.id"
              class="px-5 py-4"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="flex items-center gap-3">
                  <div class="text-base font-bold text-primary w-14 flex-shrink-0">
                    {{ formatTime(apt.starts_at) }}
                  </div>
                  <div>
                    <div class="font-semibold text-slate">{{ apt.patient_name || 'Без имени' }}</div>
                    <div class="text-sm text-gray-400">{{ apt.patient_phone }}</div>
                    <div class="text-sm text-gray-500 mt-0.5">{{ apt.service_name }}</div>
                  </div>
                </div>
                <div class="flex flex-col items-end gap-2">
                  <span class="text-xs font-semibold px-2 py-0.5 rounded-full bg-blue-100 text-blue-700">
                    {{ statusLabel[apt.status] ?? apt.status }}
                  </span>
                  <NuxtLink
                    :to="`/doctor/appointment/${apt.id}`"
                    class="text-xs text-primary font-semibold underline"
                  >
                    Открыть приём
                  </NuxtLink>
                </div>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

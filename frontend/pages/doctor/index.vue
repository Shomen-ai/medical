<!--
  Файл: pages/doctor/index.vue
  Назначение: рабочее место врача — список приёмов на день, GitHub-подобный календарь активности и краткая статистика; использует layout 'staff'.
-->
<script setup lang="ts">
definePageMeta({ layout: 'staff' })

const auth = useAuthStore()
const router = useRouter()
const { get } = useApi()
const { t } = useI18n()

onMounted(() => {
  auth.init()
  if (!auth.isDoctor) router.replace('/')
})

// ── Types ────────────────────────────────────────────────────────────
interface ScheduleCell {
  work_date: string; start_time: string; end_time: string
  is_day_off: boolean
  appointments_count: number
  pending_records_count: number
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
interface DoctorPatient {
  patient_name: string; patient_phone: string
  visits: number; first_visit: string; last_visit: string
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

// ── Calendar grid: week-aligned month view (Mon-first), padded with nulls ──
const calendarCells = computed<(ScheduleCell | null)[]>(() => {
  const map = new Map<string, ScheduleCell>()
  for (const c of workDays.value) map.set(c.work_date, c)

  const first = new Date(viewYear.value, viewMonth.value - 1, 1)
  const last = new Date(viewYear.value, viewMonth.value, 0)
  // ISO week-day: Mon=0 ... Sun=6
  const leadingBlanks = (first.getDay() + 6) % 7
  const cells: (ScheduleCell | null)[] = []
  for (let i = 0; i < leadingBlanks; i++) cells.push(null)
  for (let day = 1; day <= last.getDate(); day++) {
    const y = viewYear.value
    const m = String(viewMonth.value).padStart(2, '0')
    const d = String(day).padStart(2, '0')
    const dateStr = `${y}-${m}-${d}`
    cells.push(map.get(dateStr) ?? {
      work_date: dateStr,
      start_time: '', end_time: '', is_day_off: true,
      appointments_count: 0, pending_records_count: 0,
    })
  }
  // Trailing pad so the grid ends on a Sunday
  while (cells.length % 7 !== 0) cells.push(null)
  return cells
})

// Heatmap intensity for the cell background: 0..4 buckets
const cellTone = (c: ScheduleCell) => {
  if (c.is_day_off && c.appointments_count === 0) return 'bg-gray-50 text-gray-300'
  const n = c.appointments_count
  if (n === 0) return 'bg-white text-slate'
  if (n === 1) return 'bg-primary/15 text-slate'
  if (n === 2) return 'bg-primary/30 text-slate'
  if (n <= 4) return 'bg-primary/50 text-white'
  return 'bg-primary/75 text-white'
}

// Day-of-week headers (Mon-first, in Russian short form)
const weekDayLabels = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

// Month-wide tallies for the legend strip
const monthTotals = computed(() => {
  let appts = 0, pending = 0
  for (const c of workDays.value) {
    appts += c.appointments_count
    pending += c.pending_records_count
  }
  return { appts, pending }
})

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

// ── Excel-отчёт врача (пациенты за сегодня + уникальные + статистика) ──
const generatingReport = ref(false)

// Локализованная подпись статуса записи для отчёта.
const tStatus = (s: string): string => (({
  scheduled: t('statusScheduled'), completed: t('statusCompleted'),
  cancelled: t('statusCancelled'), rescheduled: t('statusRescheduled'),
} as Record<string, string>)[s] ?? s)

// ДД.ММ.ГГГГ без зависимости от Intl-локали.
const reportDate = (iso: string) => {
  const d = new Date(iso)
  return `${String(d.getDate()).padStart(2, '0')}.${String(d.getMonth() + 1).padStart(2, '0')}.${d.getFullYear()}`
}

const generateReport = async () => {
  generatingReport.value = true
  try {
    const y = today.getFullYear()
    const m = String(today.getMonth() + 1).padStart(2, '0')
    const d = String(today.getDate()).padStart(2, '0')
    const ds = `${y}-${m}-${d}`

    const [todayAppts, patients] = await Promise.all([
      get<Appointment[]>(`/api/doctor/appointments?date=${ds}`, token.value),
      get<DoctorPatient[]>('/api/doctor/patients', token.value),
    ])

    await downloadXlsx(`BeautyMed_${ds}.xlsx`, [
      {
        name: t('reportSheetToday'),
        rows: [
          [t('reportColTime'), t('reportColPatient'), t('reportColPhone'), t('reportColService'), t('reportColStatus')],
          ...(todayAppts ?? []).map(a => [
            formatTime(a.starts_at), a.patient_name || '—', a.patient_phone, a.service_name, tStatus(a.status),
          ]),
        ],
      },
      {
        name: t('reportSheetPatients'),
        rows: [
          [t('reportColPatient'), t('reportColPhone'), t('reportColVisits'), t('reportColFirstVisit'), t('reportColLastVisit')],
          ...(patients ?? []).map(p => [
            p.patient_name || '—', p.patient_phone, p.visits, reportDate(p.first_visit), reportDate(p.last_visit),
          ]),
        ],
      },
      {
        name: t('reportSheetStats'),
        rows: [
          [t('reportColMetric'), t('reportColValue')],
          [t('reportStatAppointmentsMonth'), stats.value?.appointments_this_month ?? 0],
          [t('reportStatUniquePatients'), stats.value?.unique_patients ?? 0],
          [t('reportStatFilledPct'), stats.value?.filled_records_pct != null ? Math.round(stats.value.filled_records_pct) : 0],
        ],
      },
    ])
  } finally {
    generatingReport.value = false
  }
}

useHead({ title: 'Кабинет врача — BeautyMed' })
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-5xl mx-auto px-4 py-8 space-y-6">

      <!-- Report button -->
      <div class="flex justify-end">
        <button
          type="button"
          class="text-sm font-semibold text-white px-4 py-2 rounded-lg transition-opacity"
          :class="generatingReport ? 'opacity-50 cursor-not-allowed' : 'hover:opacity-90'"
          style="background: linear-gradient(135deg, #005A5F, #00959D)"
          :disabled="generatingReport"
          @click="generateReport"
        >
          📊 {{ generatingReport ? t('reportGenerating') : t('reportButton') }}
        </button>
      </div>

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

      <div class="grid grid-cols-1 md:grid-cols-5 gap-6">

        <!-- Left: month calendar grid -->
        <div class="md:col-span-2 bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <!-- Month nav -->
          <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100">
            <button type="button" class="text-gray-400 hover:text-primary p-1" @click="prevMonth">‹</button>
            <span class="text-sm font-semibold text-slate capitalize">{{ monthLabel }}</span>
            <button type="button" class="text-gray-400 hover:text-primary p-1" @click="nextMonth">›</button>
          </div>

          <!-- Weekday header -->
          <div class="grid grid-cols-7 gap-1 px-3 pt-3">
            <div
              v-for="(w, i) in weekDayLabels"
              :key="w"
              class="text-[10px] font-semibold text-center pb-1"
              :class="i === 6 ? 'text-red-400' : 'text-gray-400'"
            >{{ w }}</div>
          </div>

          <!-- Calendar cells -->
          <div class="grid grid-cols-7 gap-1 px-3 pb-3">
            <template v-for="(cell, idx) in calendarCells" :key="idx">
              <div v-if="!cell" class="aspect-square"></div>
              <button
                v-else
                type="button"
                class="aspect-square relative rounded-md text-xs font-semibold flex items-start justify-end p-1 transition-all"
                :class="[
                  cellTone(cell),
                  selectedDate === cell.work_date ? 'ring-2 ring-primary' : 'hover:ring-1 hover:ring-primary/40',
                  isToday(cell.work_date) ? 'outline outline-2 outline-offset-1 outline-amber-400' : ''
                ]"
                :title="cell.is_day_off
                  ? 'Выходной'
                  : `${cell.appointments_count} приём(ов), ${cell.pending_records_count} к заполнению`"
                @click="loadDay(cell.work_date)"
              >
                <span>{{ Number(cell.work_date.slice(-2)) }}</span>
                <span
                  v-if="cell.pending_records_count > 0"
                  class="absolute bottom-1 left-1 w-1.5 h-1.5 rounded-full bg-amber-500"
                ></span>
              </button>
            </template>
          </div>

          <!-- Legend -->
          <div class="px-4 py-3 border-t border-gray-100 text-xs text-gray-500 space-y-2">
            <div class="flex items-center gap-2">
              <span>Меньше</span>
              <span class="inline-block w-3 h-3 rounded-sm bg-white border border-gray-200"></span>
              <span class="inline-block w-3 h-3 rounded-sm bg-primary/15"></span>
              <span class="inline-block w-3 h-3 rounded-sm bg-primary/30"></span>
              <span class="inline-block w-3 h-3 rounded-sm bg-primary/50"></span>
              <span class="inline-block w-3 h-3 rounded-sm bg-primary/75"></span>
              <span>Больше приёмов</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="inline-block w-1.5 h-1.5 rounded-full bg-amber-500"></span>
              <span>Есть незаполненные записи</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="inline-block w-3 h-3 rounded-sm border-2 border-amber-400"></span>
              <span>Сегодня</span>
            </div>
            <div class="flex items-center justify-between text-slate font-medium pt-1">
              <span>За месяц: {{ monthTotals.appts }} приём(ов)</span>
              <span v-if="monthTotals.pending > 0" class="text-amber-600">К заполнению: {{ monthTotals.pending }}</span>
            </div>
          </div>
        </div>

        <!-- Right: appointments for selected day -->
        <div class="md:col-span-3 bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
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

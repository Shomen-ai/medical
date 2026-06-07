<!--
  Файл: pages/admin.vue
  Назначение: админ-панель клиники со сводкой записей, графиками (Bar/Doughnut) и управлением справочниками; использует layout 'staff'.
-->
<script setup lang="ts">
import { Bar, Doughnut } from 'vue-chartjs'

definePageMeta({ layout: 'staff' })

const auth = useAuthStore()
const router = useRouter()
const { get, post, patch } = useApi()
const { t, tMed } = useI18n()

onMounted(() => {
  auth.init()
  if (!auth.isAdmin) router.replace('/')
})

// ── Types ────────────────────────────────────────────────────────────────
interface KPI { appointments_today: number; revenue_today: number; free_slots: number }
interface Stats {
  total_patients: number; active_doctors: number
  appointments_month: number; revenue_quarter: number; top_service: string
}
interface Appointment {
  id: string; patient_name: string; patient_phone: string
  doctor_name: string; service_name: string
  starts_at: string; status: string; final_price: number
}
interface MonthlyPoint { month: string; appointments: number; revenue: number }
interface PeriodStats { appointments: number; revenue: number; unique_patients: number; cancelled: number }
interface Specialty { id: string; name: string }
interface DoctorReport {
  doctor_id: string; doctor_name: string; specialty_name: string
  appointments: number; unique_patients: number
}
interface AdminReview {
  id: string; rating: number; text: string; created_at: string
  doctor_name: string; service_name: string; is_hidden: boolean
}

// ── Data fetching ────────────────────────────────────────────────────────
const token = computed(() => auth.token ?? undefined)

const { data: dashboard, refresh: refreshDashboard } = await useAsyncData('admin-dashboard', () =>
  get<{ kpi: KPI; appointments: Appointment[] }>('/api/admin/dashboard', token.value), { server: false })

const { data: stats, refresh: refreshStats } = await useAsyncData('admin-stats', () =>
  get<Stats>('/api/admin/stats', token.value), { server: false })

// ── Report date range («с/по») ────────────────────────────────────────
const fmtDate = (dt: Date) => {
  const p = (n: number) => String(n).padStart(2, '0')
  return `${dt.getFullYear()}-${p(dt.getMonth() + 1)}-${p(dt.getDate())}`
}
const todayStr = fmtDate(new Date())
const dateFrom = ref(fmtDate(new Date(new Date().getFullYear(), new Date().getMonth(), 1)))
const dateTo = ref(todayStr)
const rangeQuery = computed(() => `from=${dateFrom.value}&to=${dateTo.value}`)

const { data: monthly, refresh: refreshMonthly } = await useAsyncData('admin-monthly', () =>
  get<MonthlyPoint[]>(`/api/admin/stats/monthly?${rangeQuery.value}`, token.value), { server: false })

const { data: specialties } = await useAsyncData('admin-specialties', () =>
  get<Specialty[]>('/api/specialties'), { server: false })

// ── Пресеты периода (заполняют поля «с/по»: от начала периода до сегодня) ──
const periods = computed(() => [
  { key: 'day',     label: t('adminPeriodDay') },
  { key: 'week',    label: t('adminPeriodWeek') },
  { key: 'month',   label: t('adminPeriodMonth') },
  { key: 'quarter', label: t('adminPeriodQuarter') },
  { key: 'year',    label: t('adminPeriodYear') },
])
const applyPreset = (key: string) => {
  const n = new Date()
  let from: Date
  switch (key) {
    case 'day':     from = new Date(n.getFullYear(), n.getMonth(), n.getDate()); break
    case 'week':    { const off = (n.getDay() + 6) % 7; from = new Date(n.getFullYear(), n.getMonth(), n.getDate() - off); break }
    case 'quarter': from = new Date(n.getFullYear(), Math.floor(n.getMonth() / 3) * 3, 1); break
    case 'year':    from = new Date(n.getFullYear(), 0, 1); break
    default:        from = new Date(n.getFullYear(), n.getMonth(), 1) // month
  }
  dateFrom.value = fmtDate(from)
  dateTo.value = todayStr
}
// «с» не должно быть больше «по».
watch(dateFrom, (v) => { if (v > dateTo.value) dateTo.value = v })

const { data: periodStats, refresh: refreshPeriod } = await useAsyncData(
  'admin-period-stats',
  () => get<PeriodStats>(`/api/admin/stats/period?${rangeQuery.value}`, token.value),
  { server: false }
)

// ── Reviews moderation ─────────────────────────────────────────────────
const { data: reviews, refresh: refreshReviews } = await useAsyncData('admin-reviews',
  () => get<AdminReview[]>('/api/admin/reviews', token.value), { server: false })

const toggleReview = async (r: AdminReview) => {
  await patch(`/api/admin/reviews/${r.id}`, { hidden: !r.is_hidden }, auth.token!)
  await refreshReviews()
}

// ── Per-doctor report (тот же диапазон «с/по») ────────────────────────
const { data: byDoctor, refresh: refreshByDoctor } = await useAsyncData(
  'admin-by-doctor',
  () => get<{ doctors: DoctorReport[] }>(`/api/admin/stats/by-doctor?${rangeQuery.value}`, token.value),
  { server: false }
)
// Диапазон «с/по» обновляет данные, зависящие от него.
watch([dateFrom, dateTo], () => { refreshPeriod(); refreshByDoctor(); refreshMonthly() })

// ── Расширенный отчёт: один Excel с несколькими листами ────────────────
interface FullReport {
  summary: { total: number; completed: number; scheduled: number; cancelled: number; rescheduled: number; revenue: number; unique_patients: number }
  by_service: { service_name: string; specialty_name: string; appointments: number; revenue: number }[]
  by_specialty: { specialty_name: string; appointments: number; unique_patients: number; revenue: number }[]
  by_doctor: { doctor_name: string; specialty_name: string; appointments: number; unique_patients: number; revenue: number; rating: number }[]
  daily: { day: string; appointments: number; revenue: number }[]
  by_weekday: { n: number; count: number }[]
  by_hour: { n: number; count: number }[]
  ratings_by_doctor: { name: string; avg: number; count: number }[]
  ratings_by_service: { name: string; avg: number; count: number }[]
  retention: { month: string; new_patients: number; returning_patients: number }[]
  gender: { label: string; count: number }[]
  age: { label: string; count: number }[]
}
const exportingFull = ref(false)
const exportFullReport = async () => {
  exportingFull.value = true
  try {
    const r = await get<FullReport>(`/api/admin/report/full?${rangeQuery.value}`, token.value)
    const weekdays = t('rfWeekdays').split(',')
    const monthNamesArr = t('adminMonthShort').split(',')
    const dmy = (iso: string) => { const [y, m, d] = iso.split('-'); return `${d}.${m}.${y}` }
    const monthLabel = (ym: string) => { const [y, m] = ym.split('-'); return `${monthNamesArr[+m - 1] ?? m} ${y}` }
    const periodStr = `${dmy(dateFrom.value)} — ${dmy(dateTo.value)}`
    const pct = (a: number, b: number) => (b > 0 ? Math.round((a / b) * 100) : 0)
    const s = r.summary
    type Row = (string | number)[]

    const sheets = [
      // 1. Сводка
      {
        name: t('rfSheetSummary'),
        mergeTitleRow: true,
        cols: [34, 18],
        decimal: [], money: [1],
        rows: [
          [`${t('rfSheetSummary').toUpperCase()} · ${periodStr}`] as Row,
          [] as Row,
          [t('rfMetric'), t('rfValue')] as Row,
          [t('rfTotalAppts'), s.total],
          [t('statusCompleted'), s.completed],
          [t('statusScheduled'), s.scheduled],
          [t('statusCancelled'), s.cancelled],
          [t('adminStatPatients'), s.unique_patients],
          [t('rfRevenueTmt'), s.revenue],
          [t('rfAvgCheck'), s.completed ? Math.round(s.revenue / s.completed) : 0],
          [t('rfCompletionRate'), `${pct(s.completed, s.total)}%`],
          [t('rfCancelRate'), `${pct(s.cancelled, s.total)}%`],
        ],
      },
      // 2. По услугам
      {
        name: t('rfSheetServices'), mergeTitleRow: true, cols: [22, 34, 12, 16], money: [3],
        rows: [
          [t('rfSheetServices').toUpperCase()] as Row, [] as Row,
          [t('reportColSpecialty'), t('rfService'), t('reportColAppointments'), t('rfRevenue')],
          ...r.by_service.map(x => [tMed(x.specialty_name), tMed(x.service_name), x.appointments, x.revenue] as Row),
        ],
      },
      // 3. По специальностям
      {
        name: t('rfSheetSpecialty'), mergeTitleRow: true, cols: [26, 12, 14, 16], money: [3],
        rows: [
          [t('rfSheetSpecialty').toUpperCase()] as Row, [] as Row,
          [t('reportColSpecialty'), t('reportColAppointments'), t('adminStatPatients'), t('rfRevenue')],
          ...r.by_specialty.map(x => [tMed(x.specialty_name), x.appointments, x.unique_patients, x.revenue] as Row),
        ],
      },
      // 4. По врачам
      {
        name: t('rfSheetDoctors'), mergeTitleRow: true, cols: [24, 24, 12, 14, 16, 12], money: [4], decimal: [5],
        rows: [
          [t('rfSheetDoctors').toUpperCase()] as Row, [] as Row,
          [t('reportColDoctor'), t('reportColSpecialty'), t('reportColAppointments'), t('adminStatPatients'), t('rfRevenue'), t('rfRating')],
          ...r.by_doctor.map(x => [x.doctor_name, tMed(x.specialty_name), x.appointments, x.unique_patients, x.revenue, x.rating] as Row),
        ],
      },
      // 5. Динамика по дням
      {
        name: t('rfSheetDaily'), mergeTitleRow: true, cols: [16, 12, 16], money: [2],
        rows: [
          [t('rfSheetDaily').toUpperCase()] as Row, [] as Row,
          [t('rfDate'), t('reportColAppointments'), t('rfRevenue')],
          ...r.daily.map(x => [dmy(x.day), x.appointments, x.revenue] as Row),
        ],
      },
      // 6. Загрузка по времени
      {
        name: t('rfSheetTime'), mergeTitleRow: true, cols: [18, 12],
        rows: [
          [t('rfSheetTime').toUpperCase()] as Row, [] as Row,
          [t('rfByWeekday')] as Row, [t('rfWeekday'), t('rfCount')],
          ...r.by_weekday.map(w => [weekdays[w.n - 1] ?? String(w.n), w.count] as Row),
          [] as Row,
          [t('rfByHour')] as Row, [t('rfHour'), t('rfCount')],
          ...r.by_hour.map(h => [`${String(h.n).padStart(2, '0')}:00`, h.count] as Row),
        ],
      },
      // 7. Рейтинги
      {
        name: t('rfSheetRatings'), mergeTitleRow: true, cols: [28, 14, 12], decimal: [1],
        rows: [
          [t('rfSheetRatings').toUpperCase()] as Row, [] as Row,
          [t('rfByDoctor')] as Row, [t('reportColDoctor'), t('rfRating'), t('rfReviews')],
          ...r.ratings_by_doctor.map(x => [x.name, x.avg, x.count] as Row),
          [] as Row,
          [t('rfByService')] as Row, [t('rfService'), t('rfRating'), t('rfReviews')],
          ...r.ratings_by_service.map(x => [tMed(x.name), x.avg, x.count] as Row),
        ],
      },
      // 8. Новые и повторные
      {
        name: t('rfSheetRetention'), mergeTitleRow: true, cols: [18, 16, 18],
        rows: [
          [t('rfSheetRetention').toUpperCase()] as Row, [] as Row,
          [t('rfMonth'), t('rfNew'), t('rfReturning')],
          ...r.retention.map(x => [monthLabel(x.month), x.new_patients, x.returning_patients] as Row),
        ],
      },
      // 9. Демография
      {
        name: t('rfSheetDemographics'), mergeTitleRow: true, cols: [18, 12],
        rows: [
          [t('rfSheetDemographics').toUpperCase()] as Row, [] as Row,
          [t('rfByGender')] as Row, [t('rfGender'), t('rfCount')],
          ...r.gender.map(x => [x.label, x.count] as Row),
          [] as Row,
          [t('rfByAge')] as Row, [t('rfAge'), t('rfCount')],
          ...r.age.map(x => [x.label, x.count] as Row),
        ],
      },
    ]
    await downloadXlsx(`BeautyMed_report_${dateFrom.value}_${dateTo.value}.xlsx`, sheets)
  } finally {
    exportingFull.value = false
  }
}

// ── Appointments with date filter ─────────────────────────────────────
const filterDate = ref('')
const filterStatus = ref('scheduled')
const { data: allAppointments, pending: loadingAppts, refresh: refreshAppts } = await useAsyncData(
  'admin-appointments',
  async () => {
    const params = new URLSearchParams()
    if (filterDate.value) params.set('date', filterDate.value)
    if (filterStatus.value) params.set('status', filterStatus.value)
    return (await get<Appointment[]>(`/api/admin/appointments?${params}`, token.value)) ?? []
  },
  { server: false }
)

watch([filterDate, filterStatus], () => refreshAppts())

// Появление токена (auth.init выполняется в onMounted — ПОСЛЕ первых фетчей с server:false)
// перезагружает ВСЕ панели; иначе при жёстком обновлении/прямом заходе они остаются пустыми
// («СЕГОДНЯ»-карточки, отчёт, статусы записей, таблица записей внизу).
watch(token, (t) => {
  if (!t) return
  refreshDashboard(); refreshStats(); refreshPeriod(); refreshByDoctor()
  refreshMonthly(); refreshReviews(); refreshAppts()
})

// ── Schedule generation ────────────────────────────────────────────────
const schedYear = ref(new Date().getFullYear())
const schedMonth = ref(new Date().getMonth() + 1)
const schedSpecialty = ref('')
const schedLoading = ref(false)
const schedMsg = ref('')
const schedError = ref(false)

const generateSchedule = async () => {
  if (!schedSpecialty.value) return
  schedLoading.value = true
  schedMsg.value = ''
  schedError.value = false
  try {
    const res = await post<{ generated: number }>(
      '/api/admin/schedule/generate',
      { year: schedYear.value, month: schedMonth.value, specialty_id: schedSpecialty.value },
      token.value
    )
    schedMsg.value = t('adminGenSuccess', { n: res.generated })
  } catch {
    schedMsg.value = t('adminGenError')
    schedError.value = true
  } finally {
    schedLoading.value = false
  }
}

// ── Chart data ────────────────────────────────────────────────────────
const monthNames = computed(() => {
  const list = t('adminMonthShort').split(',')
  return Object.fromEntries(list.map((name, i) => [String(i + 1).padStart(2, '0'), name]))
})

const barChartData = computed(() => {
  const pts = monthly.value ?? []
  const names = monthNames.value
  const multiYear = new Set(pts.map(p => p.month.slice(0, 4))).size > 1
  return {
    labels: pts.map(p => {
      const [y, m] = p.month.split('-')
      const name = names[m] ?? p.month
      return multiYear ? `${name} '${y.slice(2)}` : name
    }),
    datasets: [{
      label: t('adminChartAppts'),
      data: pts.map(p => p.appointments),
      backgroundColor: '#00959D',
      borderRadius: 6,
    }],
  }
})

const barChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    y: { beginAtZero: true, ticks: { stepSize: 1 }, grid: { color: '#f1f5f9' } },
    x: { grid: { display: false } },
  },
}

const statusChartData = computed(() => {
  const apts = allAppointments.value ?? []
  const counts = { scheduled: 0, completed: 0, cancelled: 0, rescheduled: 0 }
  apts.forEach(a => {
    if (a.status in counts) counts[a.status as keyof typeof counts]++
  })
  return {
    labels: [t('statusScheduled'), t('statusCompleted'), t('statusCancelled'), t('statusRescheduled')],
    datasets: [{
      data: [counts.scheduled, counts.completed, counts.cancelled, counts.rescheduled],
      backgroundColor: ['#3B82F6', '#10B981', '#EF4444', '#F59E0B'],
      borderWidth: 0,
      hoverOffset: 4,
    }],
  }
})

const donutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom' as const, labels: { padding: 12, font: { size: 11 } } },
  },
  cutout: '65%',
}

// ── Helpers ────────────────────────────────────────────────────────────
const formatMoney = (n: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(n)

const pad2 = (n: number) => String(n).padStart(2, '0')
const formatDateTime = (iso: string) => {
  const d = new Date(iso)
  return `${pad2(d.getDate())}.${pad2(d.getMonth() + 1)}.${d.getFullYear()} `
    + `${pad2(d.getHours())}:${pad2(d.getMinutes())}`
}

const statusLabel = computed<Record<string, string>>(() => ({
  scheduled: t('statusScheduled'), completed: t('statusCompleted'),
  cancelled: t('statusCancelled'), rescheduled: t('statusRescheduled'),
}))
const statusClass: Record<string, string> = {
  scheduled: 'bg-blue-100 text-blue-700', completed: 'bg-green-100 text-green-700',
  cancelled: 'bg-red-100 text-red-600', rescheduled: 'bg-yellow-100 text-yellow-700',
}

const months = computed(() => t('monthsList').split(','))

useHead({ title: t('adminPageTitle') })
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-8 space-y-8">

      <!-- KPI row -->
      <div>
        <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">{{ t('adminToday') }}</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-3xl font-extrabold text-primary">{{ dashboard?.kpi?.appointments_today ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminApptsToday') }}</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-3xl font-extrabold text-primary">{{ dashboard?.kpi?.revenue_today != null ? formatMoney(dashboard.kpi.revenue_today) : '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminRevenueCompleted') }}</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-3xl font-extrabold text-primary">{{ dashboard?.kpi?.free_slots ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminFreeSlots') }}</div>
          </div>
        </div>
      </div>

      <!-- Period stats -->
      <div>
        <div class="flex flex-col lg:flex-row lg:items-end lg:justify-between gap-3 mb-3">
          <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide">{{ t('adminStats') }}</h2>
          <div class="flex flex-wrap items-end gap-3">
            <!-- Период «с / по» (дата «по» не может быть больше сегодня) -->
            <div class="w-40">
              <label class="text-[11px] font-semibold text-gray-500 block mb-1">{{ t('adminDateFrom') }}</label>
              <DateField v-model="dateFrom" :max="todayStr" />
            </div>
            <div class="w-40">
              <label class="text-[11px] font-semibold text-gray-500 block mb-1">{{ t('adminDateTo') }}</label>
              <DateField v-model="dateTo" :min="dateFrom" :max="todayStr" />
            </div>
            <!-- Пресеты-«быстрый выбор» (заполняют поля выше) -->
            <div class="flex gap-1 bg-gray-100 rounded-xl p-1">
              <button
                v-for="p in periods"
                :key="p.key"
                type="button"
                class="px-3 py-1 text-xs font-semibold rounded-lg text-gray-500 hover:text-primary hover:bg-white transition-colors"
                @click="applyPreset(p.key)"
              >
                {{ p.label }}
              </button>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-slate">{{ periodStats?.appointments ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminStatAppts') }}</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-primary">
              {{ periodStats?.revenue != null ? formatMoney(periodStats.revenue) : '—' }}
            </div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminStatRevenue') }}</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-slate">{{ periodStats?.unique_patients ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminStatPatients') }}</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-red-500">{{ periodStats?.cancelled ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ t('adminStatCancelled') }}</div>
          </div>
        </div>

        <div v-if="stats?.top_service" class="mt-3 bg-white rounded-2xl border border-gray-100 shadow-sm px-5 py-3 text-sm text-gray-600">
          {{ t('adminTopService') }} <span class="font-semibold text-slate">{{ tMed(stats.top_service) }}</span>
        </div>
      </div>

      <!-- Charts row -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">

        <!-- Monthly appointments bar chart -->
        <div class="md:col-span-2 bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
          <div class="font-semibold text-slate mb-4">{{ t('adminApptsByMonth') }}</div>
          <div class="h-48">
            <Bar v-if="(monthly?.length ?? 0) > 0" :data="barChartData" :options="barChartOptions" />
            <div v-else class="h-full flex items-center justify-center text-gray-400 text-sm">{{ t('adminNoData') }}</div>
          </div>
        </div>

        <!-- Status donut chart -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
          <div class="font-semibold text-slate mb-4">{{ t('adminApptStatuses') }}</div>
          <div class="h-48">
            <Doughnut v-if="(allAppointments?.length ?? 0) > 0" :data="statusChartData" :options="donutOptions" />
            <div v-else class="h-full flex items-center justify-center text-gray-400 text-sm">{{ t('adminNoData') }}</div>
          </div>
        </div>

      </div>

      <!-- Doctors report (uses the period selector above) -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex flex-wrap gap-3 items-center justify-between">
          <h2 class="font-semibold text-slate">{{ t('reportAdminTitle') }}</h2>
          <button
            type="button"
            class="text-sm font-semibold text-white px-4 py-2 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50"
            style="background: linear-gradient(135deg, #005A5F, #00959D)"
            :disabled="exportingFull"
            @click="exportFullReport"
          >
            📑 {{ exportingFull ? t('loading') : t('rfExport') }}
          </button>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-xs text-gray-400 border-b border-gray-100 bg-gray-50">
                <th class="px-5 py-3 font-semibold">{{ t('reportColDoctor') }}</th>
                <th class="px-5 py-3 font-semibold">{{ t('reportColSpecialty') }}</th>
                <th class="px-5 py-3 font-semibold text-center">{{ t('reportColAppointments') }}</th>
                <th class="px-5 py-3 font-semibold text-center">{{ t('reportColUniquePatients') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="d in byDoctor?.doctors ?? []"
                :key="d.doctor_id"
                class="border-b border-gray-50 hover:bg-gray-50 transition-colors"
              >
                <td class="px-5 py-3 font-medium text-slate">{{ d.doctor_name }}</td>
                <td class="px-5 py-3 text-gray-500">{{ tMed(d.specialty_name) }}</td>
                <td class="px-5 py-3 text-center text-slate font-semibold">{{ d.appointments }}</td>
                <td class="px-5 py-3 text-center text-slate font-semibold">{{ d.unique_patients }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Отзывы (модерация) -->
      <div>
        <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">{{ t('adminReviews') }}</h2>
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm divide-y divide-gray-100">
          <div
            v-for="r in reviews ?? []"
            :key="r.id"
            class="p-4 flex items-start gap-3"
            :class="r.is_hidden ? 'opacity-50' : ''"
          >
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <StarRating :model-value="r.rating" readonly size="text-sm" />
                <span class="text-[11px] text-gray-400 truncate">{{ tMed(r.service_name) }} · {{ r.doctor_name }}</span>
              </div>
              <p class="text-sm text-slate">{{ r.text }}</p>
            </div>
            <button
              type="button"
              class="text-xs font-semibold px-3 py-1.5 rounded-lg border shrink-0"
              :class="r.is_hidden ? 'border-primary text-primary' : 'border-red-200 text-red-500'"
              @click="toggleReview(r)"
            >
              {{ r.is_hidden ? t('adminReviewShow') : t('adminReviewHide') }}
            </button>
          </div>
          <div v-if="(reviews ?? []).length === 0" class="p-4 text-sm text-gray-400">{{ t('reviewsEmpty') }}</div>
        </div>
      </div>

      <!-- Schedule generation -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
        <h2 class="font-semibold text-slate mb-4">{{ t('adminScheduleGen') }}</h2>
        <div class="flex flex-wrap gap-3 items-end">
          <div>
            <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('adminSpecialty') }}</label>
            <select v-model="schedSpecialty" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
              <option value="">{{ t('adminSelectPlaceholder') }}</option>
              <option v-for="sp in specialties" :key="sp.id" :value="sp.id">{{ sp.name }}</option>
            </select>
          </div>
          <div>
            <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('adminMonth') }}</label>
            <select v-model="schedMonth" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
              <option v-for="(name, i) in months" :key="i" :value="i + 1">{{ name }}</option>
            </select>
          </div>
          <div>
            <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('adminYear') }}</label>
            <select v-model="schedYear" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
              <option :value="2026">2026</option>
              <option :value="2027">2027</option>
            </select>
          </div>
          <button
            type="button"
            class="text-sm font-semibold text-white px-5 py-2 rounded-lg transition-opacity"
            :class="schedLoading || !schedSpecialty ? 'opacity-50 cursor-not-allowed' : 'hover:opacity-90'"
            :disabled="schedLoading || !schedSpecialty"
            style="background: linear-gradient(135deg, #005A5F, #00959D)"
            @click="generateSchedule"
          >
            {{ schedLoading ? t('adminGenerating') : t('adminGenerate') }}
          </button>
          <span v-if="schedMsg" class="text-sm" :class="schedError ? 'text-red-500' : 'text-green-600'">
            {{ schedMsg }}
          </span>
        </div>
      </div>

      <!-- Appointments table -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex flex-wrap gap-3 items-center justify-between">
          <h2 class="font-semibold text-slate">{{ t('adminAppts') }}</h2>
          <div class="flex gap-2 items-start">
            <div class="w-40"><DateField v-model="filterDate" /></div>
            <select v-model="filterStatus" class="border border-gray-200 rounded-lg px-3 py-1.5 text-sm outline-none focus:border-primary">
              <option value="">{{ t('adminAllStatuses') }}</option>
              <option value="scheduled">{{ t('statusScheduled') }}</option>
              <option value="completed">{{ t('statusCompleted') }}</option>
              <option value="cancelled">{{ t('statusCancelled') }}</option>
            </select>
          </div>
        </div>

        <div v-if="loadingAppts" class="text-center py-10 text-gray-400">{{ t('loading') }}</div>

        <div v-else-if="!allAppointments?.length" class="text-center py-10 text-gray-400">{{ t('adminNoAppts') }}</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-xs text-gray-400 border-b border-gray-100 bg-gray-50">
                <th class="px-5 py-3 font-semibold">{{ t('adminColPatient') }}</th>
                <th class="px-5 py-3 font-semibold">{{ t('adminColDoctorService') }}</th>
                <th class="px-5 py-3 font-semibold">{{ t('adminColDateTime') }}</th>
                <th class="px-5 py-3 font-semibold">{{ t('adminColPrice') }}</th>
                <th class="px-5 py-3 font-semibold">{{ t('adminColStatus') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="apt in allAppointments"
                :key="apt.id"
                class="border-b border-gray-50 hover:bg-gray-50 transition-colors"
              >
                <td class="px-5 py-3">
                  <div class="font-medium text-slate">{{ apt.patient_name || '—' }}</div>
                  <div class="text-xs text-gray-400">{{ apt.patient_phone }}</div>
                </td>
                <td class="px-5 py-3">
                  <div class="text-slate">{{ apt.doctor_name }}</div>
                  <div class="text-xs text-gray-400">{{ tMed(apt.service_name) }}</div>
                </td>
                <td class="px-5 py-3 text-slate whitespace-nowrap">{{ formatDateTime(apt.starts_at) }}</td>
                <td class="px-5 py-3 text-slate font-semibold">{{ formatMoney(apt.final_price) }}</td>
                <td class="px-5 py-3">
                  <span class="text-xs font-semibold px-2 py-0.5 rounded-full"
                    :class="statusClass[apt.status] ?? 'bg-gray-100 text-gray-600'">
                    {{ statusLabel[apt.status] ?? apt.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  </div>
</template>

<!--
  Файл: pages/admin.vue
  Назначение: админ-панель клиники со сводкой записей, графиками (Bar/Doughnut) и управлением справочниками; использует layout 'staff'.
-->
<script setup lang="ts">
import { Bar, Doughnut } from 'vue-chartjs'

definePageMeta({ layout: 'staff' })

const auth = useAuthStore()
const router = useRouter()
const { get, post } = useApi()

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

// ── Data fetching ────────────────────────────────────────────────────────
const token = computed(() => auth.token ?? undefined)

const { data: dashboard } = await useAsyncData('admin-dashboard', () =>
  get<{ kpi: KPI; appointments: Appointment[] }>('/api/admin/dashboard', token.value), { server: false })

const { data: stats } = await useAsyncData('admin-stats', () =>
  get<Stats>('/api/admin/stats', token.value), { server: false })

const { data: monthly } = await useAsyncData('admin-monthly', () =>
  get<MonthlyPoint[]>('/api/admin/stats/monthly', token.value), { server: false })

const { data: specialties } = await useAsyncData('admin-specialties', () =>
  get<Specialty[]>('/api/specialties'), { server: false })

// ── Period stats ──────────────────────────────────────────────────────
const periods = [
  { key: 'day',     label: 'День' },
  { key: 'week',    label: 'Неделя' },
  { key: 'month',   label: 'Месяц' },
  { key: 'quarter', label: 'Квартал' },
  { key: 'year',    label: 'Год' },
]
const selectedPeriod = ref('month')
const { data: periodStats, refresh: refreshPeriod } = await useAsyncData(
  'admin-period-stats',
  () => get<PeriodStats>(`/api/admin/stats/period?period=${selectedPeriod.value}`, token.value),
  { server: false }
)
watch(selectedPeriod, () => refreshPeriod())

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

// ── Schedule generation ────────────────────────────────────────────────
const schedYear = ref(new Date().getFullYear())
const schedMonth = ref(new Date().getMonth() + 1)
const schedSpecialty = ref('')
const schedLoading = ref(false)
const schedMsg = ref('')

const generateSchedule = async () => {
  if (!schedSpecialty.value) return
  schedLoading.value = true
  schedMsg.value = ''
  try {
    const res = await post<{ generated: number }>(
      '/api/admin/schedule/generate',
      { year: schedYear.value, month: schedMonth.value, specialty_id: schedSpecialty.value },
      token.value
    )
    schedMsg.value = `Сгенерировано ${res.generated} записей`
  } catch {
    schedMsg.value = 'Ошибка генерации'
  } finally {
    schedLoading.value = false
  }
}

// ── Chart data ────────────────────────────────────────────────────────
const monthNames: Record<string, string> = {
  '01': 'Янв', '02': 'Фев', '03': 'Мар', '04': 'Апр',
  '05': 'Май', '06': 'Июн', '07': 'Июл', '08': 'Авг',
  '09': 'Сен', '10': 'Окт', '11': 'Ноя', '12': 'Дек',
}

const barChartData = computed(() => {
  const pts = monthly.value ?? []
  return {
    labels: pts.map(p => {
      const [, m] = p.month.split('-')
      return monthNames[m] ?? p.month
    }),
    datasets: [{
      label: 'Записей',
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
    labels: ['Запланировано', 'Завершено', 'Отменено', 'Перенесено'],
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

const formatDateTime = (iso: string) => {
  const d = new Date(iso)
  return d.toLocaleDateString('ru-RU') + ' ' + d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}

const statusLabel: Record<string, string> = {
  scheduled: 'Запланировано', completed: 'Завершено',
  cancelled: 'Отменено', rescheduled: 'Перенесено',
}
const statusClass: Record<string, string> = {
  scheduled: 'bg-blue-100 text-blue-700', completed: 'bg-green-100 text-green-700',
  cancelled: 'bg-red-100 text-red-600', rescheduled: 'bg-yellow-100 text-yellow-700',
}

const months = ['Январь','Февраль','Март','Апрель','Май','Июнь',
  'Июль','Август','Сентябрь','Октябрь','Ноябрь','Декабрь']

useHead({ title: 'Панель администратора — BeautyMed' })
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-8 space-y-8">

      <!-- KPI row -->
      <div>
        <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Сегодня</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-3xl font-extrabold text-primary">{{ dashboard?.kpi?.appointments_today ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">Записей на сегодня</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-3xl font-extrabold text-primary">{{ dashboard?.kpi?.revenue_today != null ? formatMoney(dashboard.kpi.revenue_today) : '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">Выручка (завершённые)</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-3xl font-extrabold text-primary">{{ dashboard?.kpi?.free_slots ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">Свободных слотов</div>
          </div>
        </div>
      </div>

      <!-- Period stats -->
      <div>
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide">Статистика</h2>
          <!-- Period tabs -->
          <div class="flex gap-1 bg-gray-100 rounded-xl p-1">
            <button
              v-for="p in periods"
              :key="p.key"
              type="button"
              class="px-3 py-1 text-xs font-semibold rounded-lg transition-colors"
              :class="selectedPeriod === p.key
                ? 'bg-white text-primary shadow-sm'
                : 'text-gray-500 hover:text-slate'"
              @click="selectedPeriod = p.key"
            >
              {{ p.label }}
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-slate">{{ periodStats?.appointments ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">Записей</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-primary">
              {{ periodStats?.revenue != null ? formatMoney(periodStats.revenue) : '—' }}
            </div>
            <div class="text-xs text-gray-500 mt-1">Выручка</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-slate">{{ periodStats?.unique_patients ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">Пациентов</div>
          </div>
          <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 text-center">
            <div class="text-2xl font-extrabold text-red-500">{{ periodStats?.cancelled ?? '—' }}</div>
            <div class="text-xs text-gray-500 mt-1">Отменено</div>
          </div>
        </div>

        <div v-if="stats?.top_service" class="mt-3 bg-white rounded-2xl border border-gray-100 shadow-sm px-5 py-3 text-sm text-gray-600">
          🏆 Популярная услуга за месяц: <span class="font-semibold text-slate">{{ stats.top_service }}</span>
        </div>
      </div>

      <!-- Charts row -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">

        <!-- Monthly appointments bar chart -->
        <div class="md:col-span-2 bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
          <div class="font-semibold text-slate mb-4">Записи по месяцам</div>
          <div class="h-48">
            <Bar v-if="(monthly?.length ?? 0) > 0" :data="barChartData" :options="barChartOptions" />
            <div v-else class="h-full flex items-center justify-center text-gray-400 text-sm">Нет данных</div>
          </div>
        </div>

        <!-- Status donut chart -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
          <div class="font-semibold text-slate mb-4">Статусы записей</div>
          <div class="h-48">
            <Doughnut v-if="(allAppointments?.length ?? 0) > 0" :data="statusChartData" :options="donutOptions" />
            <div v-else class="h-full flex items-center justify-center text-gray-400 text-sm">Нет данных</div>
          </div>
        </div>

      </div>

      <!-- Schedule generation -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
        <h2 class="font-semibold text-slate mb-4">Генерация расписания</h2>
        <div class="flex flex-wrap gap-3 items-end">
          <div>
            <label class="text-xs font-semibold text-gray-500 block mb-1">Специальность</label>
            <select v-model="schedSpecialty" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
              <option value="">— выберите —</option>
              <option v-for="sp in specialties" :key="sp.id" :value="sp.id">{{ sp.name }}</option>
            </select>
          </div>
          <div>
            <label class="text-xs font-semibold text-gray-500 block mb-1">Месяц</label>
            <select v-model="schedMonth" class="border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
              <option v-for="(name, i) in months" :key="i" :value="i + 1">{{ name }}</option>
            </select>
          </div>
          <div>
            <label class="text-xs font-semibold text-gray-500 block mb-1">Год</label>
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
            {{ schedLoading ? 'Генерируем...' : 'Сгенерировать' }}
          </button>
          <span v-if="schedMsg" class="text-sm" :class="schedMsg.includes('Ошибка') ? 'text-red-500' : 'text-green-600'">
            {{ schedMsg }}
          </span>
        </div>
      </div>

      <!-- Appointments table -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex flex-wrap gap-3 items-center justify-between">
          <h2 class="font-semibold text-slate">Записи</h2>
          <div class="flex gap-2">
            <input
              v-model="filterDate"
              type="date"
              class="border border-gray-200 rounded-lg px-3 py-1.5 text-sm outline-none focus:border-primary"
            >
            <select v-model="filterStatus" class="border border-gray-200 rounded-lg px-3 py-1.5 text-sm outline-none focus:border-primary">
              <option value="">Все статусы</option>
              <option value="scheduled">Запланировано</option>
              <option value="completed">Завершено</option>
              <option value="cancelled">Отменено</option>
            </select>
          </div>
        </div>

        <div v-if="loadingAppts" class="text-center py-10 text-gray-400">Загружаем...</div>

        <div v-else-if="!allAppointments?.length" class="text-center py-10 text-gray-400">Нет записей</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-xs text-gray-400 border-b border-gray-100 bg-gray-50">
                <th class="px-5 py-3 font-semibold">Пациент</th>
                <th class="px-5 py-3 font-semibold">Врач / Услуга</th>
                <th class="px-5 py-3 font-semibold">Дата и время</th>
                <th class="px-5 py-3 font-semibold">Цена</th>
                <th class="px-5 py-3 font-semibold">Статус</th>
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
                  <div class="text-xs text-gray-400">{{ apt.service_name }}</div>
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

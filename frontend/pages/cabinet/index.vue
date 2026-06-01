<!--
  Файл: pages/cabinet/index.vue
  Назначение: личный кабинет пациента со списком его записей, возможностью отмены и переноса визита, а также навигацией к профилю и квитанциям.
-->
<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const { get, patch } = useApi()
const { t } = useI18n()

onMounted(() => {
  auth.init()
  if (!auth.isLoggedIn || auth.isAdmin || auth.isDoctor) {
    router.replace('/')
  }
})

interface Appointment {
  id: string
  doctor_id: string
  service_id: string
  doctor_name: string
  service_name: string
  starts_at: string
  status: string
  final_price: number
}

interface TimeSlot {
  starts_at: string  // "HH:MM"
  ends_at: string
}

const { data: appointments, pending, error, refresh } = await useAsyncData(
  'cabinet-appointments',
  () => get<Appointment[]>('/api/cabinet/appointments', auth.token ?? undefined),
  { server: false }
)

const pad2 = (n: number) => String(n).padStart(2, '0')
const formatDateTime = (iso: string) => {
  const d = new Date(iso)
  return `${pad2(d.getDate())}.${pad2(d.getMonth() + 1)}.${d.getFullYear()}, `
    + `${pad2(d.getHours())}:${pad2(d.getMinutes())}`
}

const formatMoney = (n: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(n)

const statusLabel = computed<Record<string, string>>(() => ({
  scheduled: t('statusScheduled'),
  completed: t('statusCompleted'),
  cancelled: t('statusCancelled'),
  rescheduled: t('statusRescheduled'),
}))
const statusClass: Record<string, string> = {
  scheduled: 'bg-blue-100 text-blue-700',
  completed: 'bg-green-100 text-green-700',
  cancelled: 'bg-red-100 text-red-600',
  rescheduled: 'bg-yellow-100 text-yellow-700',
}

const canReschedule = (apt: Appointment): boolean => {
  if (apt.status !== 'scheduled') return false
  return new Date(apt.starts_at).getTime() - Date.now() >= 2 * 60 * 60 * 1000
}

const cancelling = ref<string | null>(null)
const cancelAppointment = async (id: string) => {
  if (!auth.token) return
  cancelling.value = id
  try {
    await patch(`/api/cabinet/appointments/${id}/cancel`, {}, auth.token)
    if (appointments.value) {
      const apt = appointments.value.find(a => a.id === id)
      if (apt) apt.status = 'cancelled'
    }
  } catch { /* ignore */ } finally {
    cancelling.value = null
  }
}

// Reschedule state
const reschedulingId = ref<string | null>(null)
const availableDates = ref<string[]>([])
const slots = ref<TimeSlot[]>([])
const selectedDate = ref<string>('')
const selectedSlot = ref<string>('')
const loadingDates = ref(false)
const loadingSlots = ref(false)
const submittingReschedule = ref(false)
const rescheduleError = ref('')

const openReschedule = async (apt: Appointment) => {
  reschedulingId.value = apt.id
  availableDates.value = []
  slots.value = []
  selectedDate.value = ''
  selectedSlot.value = ''
  rescheduleError.value = ''
  loadingDates.value = true
  try {
    // Fetch dates for current month + next month, merged
    const d = new Date()
    const months: string[] = []
    for (let i = 0; i < 2; i++) {
      const m = new Date(d.getFullYear(), d.getMonth() + i, 1)
      months.push(`${m.getFullYear()}-${String(m.getMonth() + 1).padStart(2, '0')}`)
    }
    const results = await Promise.all(months.map(m =>
      get<string[]>(`/api/doctors/${apt.doctor_id}/available-dates?month=${m}`)
    ))
    const todayISO = new Date().toISOString().slice(0, 10)
    availableDates.value = results.flat().filter(date => date >= todayISO)
  } catch {
    rescheduleError.value = t('cabRescheduleDatesError')
  } finally {
    loadingDates.value = false
  }
}

const closeReschedule = () => {
  reschedulingId.value = null
  rescheduleError.value = ''
}

const selectDate = async (date: string) => {
  const apt = appointments.value?.find(a => a.id === reschedulingId.value)
  if (!apt) return
  selectedDate.value = date
  selectedSlot.value = ''
  slots.value = []
  loadingSlots.value = true
  rescheduleError.value = ''
  try {
    slots.value = await get<TimeSlot[]>(
      `/api/doctors/${apt.doctor_id}/slots?service_id=${apt.service_id}&date=${date}`
    )
  } catch {
    rescheduleError.value = t('cabRescheduleSlotsError')
  } finally {
    loadingSlots.value = false
  }
}

const submitReschedule = async () => {
  if (!auth.token || !reschedulingId.value || !selectedDate.value || !selectedSlot.value) return
  submittingReschedule.value = true
  rescheduleError.value = ''
  try {
    const startsAt = `${selectedDate.value}T${selectedSlot.value}:00+04:00`
    await patch(
      `/api/cabinet/appointments/${reschedulingId.value}/reschedule`,
      { starts_at: startsAt },
      auth.token
    )
    closeReschedule()
    await refresh()
  } catch (err: unknown) {
    const status = (err as { status?: number })?.status
    const msg = (err as { data?: { error?: string } })?.data?.error ?? ''
    if (status === 409 || msg.includes('taken')) {
      rescheduleError.value = t('cabSlotTaken')
    } else if (msg.includes('2 hours')) {
      rescheduleError.value = t('cabRescheduleTooLate')
    } else {
      rescheduleError.value = t('cabRescheduleError')
    }
  } finally {
    submittingReschedule.value = false
  }
}

const formatDateChip = (iso: string) => {
  const [, m, d] = iso.split('-').map(Number)
  return `${d} ${t('cabMonthsShort').split(',')[m - 1]}`
}

useHead({ title: t('cabPageTitle') })
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-10">
    <div class="flex items-center justify-between mb-8 flex-wrap gap-3">
      <h1 class="text-2xl font-bold text-slate">{{ t('cabMyAppts') }}</h1>
      <div class="flex items-center gap-4">
        <NuxtLink to="/cabinet/profile" class="text-sm text-muted hover:text-primary hover:underline">
          {{ t('cabProfile') }}
        </NuxtLink>
        <button
          type="button"
          class="text-sm text-muted underline"
          @click="auth.logout(); router.replace('/')"
        >
          {{ t('logout') }}
        </button>
      </div>
    </div>

    <div v-if="pending" class="text-center py-16 text-muted">{{ t('cabLoading') }}</div>
    <div v-else-if="error" class="text-center py-16 text-red-500">{{ t('cabLoadError') }}</div>
    <div v-else-if="!appointments?.length" class="text-center py-16 text-muted">
      {{ t('cabNoAppts') }}
    </div>

    <div v-else class="flex flex-col gap-4">
      <div
        v-for="apt in appointments"
        :key="apt.id"
        class="bg-white rounded-2xl border border-border p-5 shadow-sm"
      >
        <div class="flex items-start justify-between gap-4">
          <div>
            <div class="font-semibold text-slate">{{ apt.doctor_name }}</div>
            <div class="text-sm text-muted mt-0.5">{{ apt.service_name }}</div>
          </div>
          <span
            class="text-xs font-semibold px-2.5 py-1 rounded-full flex-shrink-0"
            :class="statusClass[apt.status] ?? 'bg-gray-100 text-gray-600'"
          >
            {{ statusLabel[apt.status] ?? apt.status }}
          </span>
        </div>
        <div class="mt-3 flex items-center justify-between flex-wrap gap-2">
          <span class="text-sm text-muted">{{ formatDateTime(apt.starts_at) }}</span>
          <div class="flex items-center gap-3">
            <span class="text-sm font-semibold text-primary">{{ formatMoney(apt.final_price) }}</span>
            <button
              v-if="canReschedule(apt) && reschedulingId !== apt.id"
              type="button"
              class="text-xs text-primary underline"
              @click="openReschedule(apt)"
            >
              {{ t('cabReschedule') }}
            </button>
            <button
              v-if="apt.status === 'scheduled'"
              type="button"
              class="text-xs text-red-500 underline"
              :disabled="cancelling === apt.id"
              @click="cancelAppointment(apt.id)"
            >
              {{ cancelling === apt.id ? '...' : t('cabCancel') }}
            </button>
          </div>
        </div>

        <!-- Reschedule panel -->
        <div
          v-if="reschedulingId === apt.id"
          class="mt-4 pt-4 border-t border-border space-y-3"
        >
          <div class="flex items-center justify-between">
            <div class="text-sm font-semibold text-slate">{{ t('cabSelectNewTime') }}</div>
            <button type="button" class="text-xs text-muted hover:text-slate" @click="closeReschedule">
              {{ t('close') }}
            </button>
          </div>

          <!-- Dates -->
          <div>
            <div class="text-xs font-semibold text-muted mb-1.5">{{ t('cabDate') }}</div>
            <div v-if="loadingDates" class="text-xs text-muted">{{ t('loading') }}</div>
            <div v-else-if="!availableDates.length" class="text-xs text-muted">
              {{ t('cabNoDates') }}
            </div>
            <div v-else class="flex flex-wrap gap-2">
              <button
                v-for="d in availableDates"
                :key="d"
                type="button"
                class="px-3 py-1.5 rounded-lg border text-xs transition-colors"
                :class="selectedDate === d
                  ? 'border-primary bg-primary/5 text-primary font-semibold'
                  : 'border-border text-slate hover:border-primary'"
                @click="selectDate(d)"
              >
                {{ formatDateChip(d) }}
              </button>
            </div>
          </div>

          <!-- Slots -->
          <div v-if="selectedDate">
            <div class="text-xs font-semibold text-muted mb-1.5">{{ t('cabTime') }}</div>
            <div v-if="loadingSlots" class="text-xs text-muted">{{ t('cabLoadingSlots') }}</div>
            <div v-else-if="!slots.length" class="text-xs text-muted">{{ t('cabNoSlots') }}</div>
            <div v-else class="flex flex-wrap gap-2">
              <button
                v-for="s in slots"
                :key="s.starts_at"
                type="button"
                class="px-3 py-1.5 rounded-lg border text-xs transition-colors"
                :class="selectedSlot === s.starts_at
                  ? 'border-primary bg-primary/5 text-primary font-semibold'
                  : 'border-border text-slate hover:border-primary'"
                @click="selectedSlot = s.starts_at"
              >
                {{ s.starts_at }}
              </button>
            </div>
          </div>

          <div v-if="rescheduleError" class="text-xs text-red-500">{{ rescheduleError }}</div>

          <button
            v-if="selectedDate && selectedSlot"
            type="button"
            class="w-full sm:w-auto bg-primary text-white px-4 py-2 rounded-lg text-sm font-semibold"
            :class="submittingReschedule ? 'opacity-50 cursor-not-allowed' : ''"
            :disabled="submittingReschedule"
            @click="submitReschedule"
          >
            {{ submittingReschedule ? t('cabRescheduling') : t('cabConfirmReschedule') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

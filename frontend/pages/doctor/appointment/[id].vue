<!--
  Файл: pages/doctor/appointment/[id].vue
  Назначение: карточка приёма для врача: ввод жалоб, диагноза и рекомендаций, сохранение черновика и завершение приёма с формированием медицинской записи.
-->
<script setup lang="ts">
definePageMeta({ layout: 'staff' })

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const { get, patch } = useApi()
const { t } = useI18n()

const appointmentId = computed(() => route.params.id as string)

onMounted(() => {
  auth.init()
  if (!auth.isDoctor) router.replace('/')
})

interface Appointment {
  id: string
  patient_name: string
  patient_phone: string
  service_name: string
  starts_at: string
  ends_at: string
  status: string
  final_price: number
}

interface AppointmentRecord {
  id: string
  appointment_id: string
  complaints: string
  diagnosis: string
  prescription: string | null
  recommendations: string | null
  is_draft: boolean
}

const token = computed(() => auth.token ?? undefined)

const appointment = ref<Appointment | null>(null)
const loading = ref(true)
const loadError = ref('')

const complaints = ref('')
const diagnosis = ref('')
const prescription = ref('')
const recommendations = ref('')

const saving = ref(false)
const saveError = ref('')
const justSaved = ref(false)

const loadAppointment = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const data = await get<{ appointment: Appointment; record: AppointmentRecord | null }>(
      `/api/doctor/appointments/${appointmentId.value}`,
      token.value
    )
    appointment.value = data.appointment
    if (data.record) {
      complaints.value = data.record.complaints
      diagnosis.value = data.record.diagnosis
      prescription.value = data.record.prescription ?? ''
      recommendations.value = data.record.recommendations ?? ''
    }
  } catch (err: unknown) {
    const status = (err as { status?: number })?.status
    loadError.value = status === 404
      ? t('recNotFound')
      : status === 403
        ? t('recNoAccess')
        : t('recLoadError')
  } finally {
    loading.value = false
  }
}

watch(appointmentId, loadAppointment, { immediate: true })

const canSubmit = computed(() => !!complaints.value.trim() && !!diagnosis.value.trim())

const saveRecord = async (isDraft: boolean) => {
  if (!canSubmit.value) {
    saveError.value = t('recRequiredError')
    return
  }
  if (!auth.token) return
  saving.value = true
  saveError.value = ''
  justSaved.value = false
  try {
    await patch(
      `/api/doctor/appointments/${appointmentId.value}/record`,
      {
        complaints: complaints.value.trim(),
        diagnosis: diagnosis.value.trim(),
        prescription: prescription.value.trim() || null,
        recommendations: recommendations.value.trim() || null,
        is_draft: isDraft,
      },
      auth.token
    )
    justSaved.value = true
    if (!isDraft) {
      // Completed → return to schedule
      setTimeout(() => router.push('/doctor'), 800)
    } else if (appointment.value) {
      // Refresh local status if we want to reflect any backend changes
      appointment.value.status = appointment.value.status
    }
  } catch (err: unknown) {
    const msg = (err as { data?: { error?: string } })?.data?.error ?? ''
    saveError.value = msg || t('recSaveError')
  } finally {
    saving.value = false
  }
}

const pad2 = (n: number) => String(n).padStart(2, '0')
const formattedDateTime = computed(() => {
  if (!appointment.value) return ''
  const d = new Date(appointment.value.starts_at)
  return `${pad2(d.getDate())}.${pad2(d.getMonth() + 1)}.${d.getFullYear()}, `
    + `${pad2(d.getHours())}:${pad2(d.getMinutes())}`
})

const statusLabel = computed<Record<string, string>>(() => ({
  scheduled: t('statusScheduled'),
  completed: t('statusCompleted'),
  cancelled: t('statusCancelled'),
  rescheduled: t('statusRescheduled'),
}))

useHead({ title: t('recPageTitle') })
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-3xl mx-auto px-4 py-8">
      <NuxtLink to="/doctor" class="inline-flex items-center text-sm text-muted hover:text-primary mb-4">
        {{ t('recBackToSchedule') }}
      </NuxtLink>

      <div v-if="loading" class="text-center py-16 text-muted">{{ t('recLoading') }}</div>

      <div v-else-if="loadError" class="text-center py-16">
        <div class="text-red-500 mb-3">{{ loadError }}</div>
        <NuxtLink to="/doctor" class="text-sm text-primary underline">{{ t('recReturnToSchedule') }}</NuxtLink>
      </div>

      <div v-else-if="appointment" class="space-y-6">
        <!-- Patient header -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
          <div class="flex items-start justify-between gap-4">
            <div>
              <div class="text-xl font-bold text-slate">{{ appointment.patient_name || t('recNoName') }}</div>
              <div class="text-sm text-muted mt-0.5">{{ appointment.patient_phone }}</div>
              <div class="text-sm text-slate mt-2">{{ appointment.service_name }}</div>
              <div class="text-sm text-muted">{{ formattedDateTime }}</div>
            </div>
            <span
              class="text-xs font-semibold px-2.5 py-1 rounded-full flex-shrink-0"
              :class="appointment.status === 'completed'
                ? 'bg-green-100 text-green-700'
                : 'bg-blue-100 text-blue-700'"
            >
              {{ statusLabel[appointment.status] ?? appointment.status }}
            </span>
          </div>
        </div>

        <!-- Record form -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 space-y-5">
          <div>
            <label class="text-sm font-semibold text-slate block mb-1.5">
              {{ t('recComplaintsLabel') }} <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="complaints"
              rows="3"
              :placeholder="t('recComplaintsPlaceholder')"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary resize-y"
            />
          </div>

          <div>
            <label class="text-sm font-semibold text-slate block mb-1.5">
              {{ t('recDiagnosisLabel') }} <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="diagnosis"
              rows="3"
              :placeholder="t('recDiagnosisPlaceholder')"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary resize-y"
            />
          </div>

          <div>
            <label class="text-sm font-semibold text-slate block mb-1.5">
              {{ t('recPrescriptionLabel') }} <span class="text-xs text-muted font-normal">{{ t('recOptional') }}</span>
            </label>
            <textarea
              v-model="prescription"
              rows="2"
              :placeholder="t('recPrescriptionPlaceholder')"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary resize-y"
            />
          </div>

          <div>
            <label class="text-sm font-semibold text-slate block mb-1.5">
              {{ t('recRecommendationsLabel') }} <span class="text-xs text-muted font-normal">{{ t('recOptional') }}</span>
            </label>
            <textarea
              v-model="recommendations"
              rows="2"
              :placeholder="t('recRecommendationsPlaceholder')"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary resize-y"
            />
          </div>

          <div v-if="saveError" class="text-sm text-red-500">{{ saveError }}</div>
          <div v-else-if="justSaved" class="text-sm text-emerald-600">{{ t('recSaved') }}</div>

          <div class="flex flex-wrap gap-3 pt-2">
            <button
              v-if="appointment.status !== 'completed'"
              type="button"
              class="bg-primary text-white px-5 py-2.5 rounded-lg text-sm font-semibold transition-opacity"
              :class="(!canSubmit || saving) ? 'opacity-50 cursor-not-allowed' : ''"
              :disabled="!canSubmit || saving"
              @click="saveRecord(false)"
            >
              {{ saving ? t('recSaving') : t('recComplete') }}
            </button>
            <button
              type="button"
              class="border border-border text-slate px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-gray-50 transition-colors"
              :class="(!canSubmit || saving) ? 'opacity-50 cursor-not-allowed' : ''"
              :disabled="!canSubmit || saving"
              @click="saveRecord(true)"
            >
              {{ appointment.status === 'completed' ? t('recSaveChanges') : t('recSaveDraft') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

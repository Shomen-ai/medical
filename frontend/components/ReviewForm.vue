<!--
  Файл: components/ReviewForm.vue
  Назначение: форма написания отзыва — выбор завершённого визита, оценка звёздами и текст.
  Доступна только авторизованному пациенту; визиты тянутся из /api/cabinet/reviewable.
-->
<script setup lang="ts">
import type { ReviewableAppt, ReviewItem } from '~/types'

const emit = defineEmits<{ created: [review: ReviewItem] }>()

const { t } = useI18n()
const { get, post } = useApi()
const auth = useAuthStore()

const appts = ref<ReviewableAppt[]>([])
const appointmentId = ref('')
const rating = ref(0)
const text = ref('')
const submitting = ref(false)
const error = ref('')
const done = ref(false)

const formatDate = (iso: string) => {
  const d = new Date(iso)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}.${p(d.getMonth() + 1)}.${d.getFullYear()}`
}

const loadAppts = async () => {
  if (!auth.isPatient || !auth.token) return
  try {
    appts.value = (await get<ReviewableAppt[]>('/api/cabinet/reviewable', auth.token)) ?? []
  } catch { /* оставляем пустым */ }
}

onMounted(() => { auth.init(); loadAppts() })

const canSubmit = computed(() =>
  !!appointmentId.value && rating.value >= 1 && text.value.trim().length >= 3 && !submitting.value)

const submit = async () => {
  error.value = ''
  if (!canSubmit.value) return
  submitting.value = true
  try {
    const review = await post<ReviewItem>('/api/reviews', {
      appointment_id: appointmentId.value,
      rating: rating.value,
      text: text.value.trim(),
    }, auth.token!)
    done.value = true
    appointmentId.value = ''; rating.value = 0; text.value = ''
    emit('created', review)
  } catch (e: any) {
    const status = e?.response?.status ?? e?.statusCode
    error.value = status === 403 ? t('reviewsNeedVisit') : t('reviewsError')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="bg-white rounded-2xl shadow-card p-6">
    <h3 class="text-lg font-bold text-slate mb-4">{{ t('reviewsWrite') }}</h3>

    <p v-if="!auth.isPatient" class="text-sm text-muted">{{ t('reviewsLoginHint') }}</p>
    <p v-else-if="appts.length === 0 && !done" class="text-sm text-muted">{{ t('reviewsNeedVisit') }}</p>
    <p v-else-if="done" class="text-sm text-primary font-semibold">{{ t('reviewsThanks') }}</p>

    <form v-else class="space-y-4" @submit.prevent="submit">
      <div>
        <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('reviewsFormVisit') }}</label>
        <select v-model="appointmentId" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
          <option value="">{{ t('reviewsFormVisitPick') }}</option>
          <option v-for="a in appts" :key="a.appointment_id" :value="a.appointment_id">
            {{ a.service_name }} · {{ a.doctor_name }} · {{ formatDate(a.starts_at) }}
          </option>
        </select>
      </div>

      <div>
        <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('reviewsFormRating') }}</label>
        <StarRating v-model="rating" />
      </div>

      <div>
        <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('reviewsFormText') }}</label>
        <textarea
          v-model="text" rows="4" maxlength="1000"
          class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary resize-none"
        />
      </div>

      <p v-if="error" class="text-sm text-red-500">{{ error }}</p>

      <button
        type="submit" :disabled="!canSubmit"
        class="w-full text-white px-5 py-2.5 rounded-lg text-sm font-semibold transition-opacity disabled:opacity-50"
        style="background: linear-gradient(135deg, #005A5F, #00959D)"
      >
        {{ submitting ? t('loading') : t('reviewsSubmit') }}
      </button>
    </form>
  </div>
</template>

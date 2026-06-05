<!--
  Файл: components/ReviewForm.vue
  Назначение: форма написания отзыва. Доступна любому авторизованному пациенту;
  врач и услуга указываются по желанию (для фильтрации на странице отзывов).
-->
<script setup lang="ts">
import type { Doctor, Service, ReviewItem } from '~/types'

const props = defineProps<{ doctors: Doctor[]; services: Service[] }>()
const emit = defineEmits<{ created: [review: ReviewItem] }>()

const { t } = useI18n()
const { post } = useApi()
const auth = useAuthStore()

const doctorId = ref('')
const serviceId = ref('')
const rating = ref(0)
const text = ref('')
const submitting = ref(false)
const error = ref('')
const done = ref(false)

onMounted(() => auth.init())

const canSubmit = computed(() =>
  rating.value >= 1 && text.value.trim().length >= 3 && !submitting.value)

const submit = async () => {
  error.value = ''
  if (!canSubmit.value) return
  submitting.value = true
  try {
    const review = await post<ReviewItem>('/api/reviews', {
      rating: rating.value,
      text: text.value.trim(),
      doctor_id: doctorId.value,
      service_id: serviceId.value,
    }, auth.token!)
    done.value = true
    doctorId.value = ''; serviceId.value = ''; rating.value = 0; text.value = ''
    emit('created', review)
  } catch {
    error.value = t('reviewsError')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="bg-white rounded-2xl shadow-card p-6">
    <h3 class="text-lg font-bold text-slate mb-4">{{ t('reviewsWrite') }}</h3>

    <p v-if="!auth.isPatient" class="text-sm text-muted">{{ t('reviewsLoginHint') }}</p>
    <p v-else-if="done" class="text-sm text-primary font-semibold">{{ t('reviewsThanks') }}</p>

    <form v-else class="space-y-4" @submit.prevent="submit">
      <div class="flex flex-col sm:flex-row gap-3">
        <div class="flex-1">
          <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('reviewsFormService') }}</label>
          <select v-model="serviceId" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
            <option value="">{{ t('reviewsFormOptional') }}</option>
            <option v-for="s in services" :key="s.id" :value="s.id">{{ s.name }}</option>
          </select>
        </div>
        <div class="flex-1">
          <label class="text-xs font-semibold text-gray-500 block mb-1">{{ t('reviewsFormDoctor') }}</label>
          <select v-model="doctorId" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:border-primary">
            <option value="">{{ t('reviewsFormOptional') }}</option>
            <option v-for="d in doctors" :key="d.id" :value="d.id">{{ d.full_name }}</option>
          </select>
        </div>
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

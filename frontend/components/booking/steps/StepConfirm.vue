<!--
  Файл: components/booking/steps/StepConfirm.vue
  Назначение: финальный шаг мастера записи — сводка по записи, ввод телефона и подтверждение SMS-кодом (OTP), создание записи в БД.
-->
<script setup lang="ts">
const booking = useBookingStore()
const auth = useAuthStore()
const { post } = useApi()
const { t } = useI18n()
const config = useRuntimeConfig()

const submittingOtp = ref(false)
const submittingBooking = ref(false)
const submittingPromo = ref(false)
const otpCode = ref('')
const otpError = ref('')
const bookingError = ref('')
const consent = ref(false)

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(price)

const applyPromo = async () => {
  if (!booking.serviceId) return
  submittingPromo.value = true
  try {
    const res = await post<{ valid: boolean; discount_pct: number; final_price: number }>(
      '/api/promo/check',
      { code: booking.promoCode.trim(), service_id: booking.serviceId }
    )
    booking.promoValid = booking.promoCode.trim() === '' ? null : res.valid
    booking.promoDiscountPct = res.discount_pct
    booking.finalPrice = res.final_price
  } catch {
    booking.promoValid = false
    booking.promoDiscountPct = 0
    booking.finalPrice = booking.servicePrice
  } finally {
    submittingPromo.value = false
  }
}

// Reactive countdown — updated every second via interval
const cooldownSecs = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const startTimer = () => {
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    const until = booking.otpCooldownUntil
    cooldownSecs.value = until ? Math.max(0, Math.ceil((until - Date.now()) / 1000)) : 0
    if (cooldownSecs.value === 0 && timer) { clearInterval(timer); timer = null }
  }, 500)
}

onUnmounted(() => { if (timer) clearInterval(timer) })

const phoneCheck = computed(() => validatePhone(booking.phone))

// Форматирует ввод по маске +993 65 12-34-56 на каждый input.
const onPhoneInput = (e: Event) => {
  booking.phone = formatTmPhone((e.target as HTMLInputElement).value)
  otpError.value = ''
}

const sendOtp = async () => {
  if (!booking.phone || !consent.value) return
  if (!phoneCheck.value.valid) {
    otpError.value = t('confirmPhoneErr')
    return
  }
  // Canonical form sent to API and stored for OTP verification.
  booking.phone = phoneCheck.value.e164
  submittingOtp.value = true
  otpError.value = ''
  try {
    await post('/api/auth/otp', { phone: booking.phone })
    booking.otpSent = true
    booking.startOtpCooldown()
    startTimer()
  } catch {
    otpError.value = t('confirmSendErr')
  } finally {
    submittingOtp.value = false
  }
}

const confirmBooking = async () => {
  if (!otpCode.value) return
  submittingBooking.value = true
  bookingError.value = ''
  try {
    // 1. Verify OTP
    const { access_token } = await post<{ access_token: string }>(
      '/api/auth/verify',
      { phone: booking.phone, code: otpCode.value }
    )
    booking.token = access_token
    // Save to auth store so cabinet page is accessible
    auth.token = access_token
    auth.role = 'patient'
    if (import.meta.client) {
      localStorage.setItem('auth_token', access_token)
      localStorage.setItem('auth_role', 'patient')
    }

    // 2. Build RFC3339 datetime: date "YYYY-MM-DD" + time "HH:MM" → "YYYY-MM-DDTHH:MM:00+04:00"
    const startsAt = `${booking.date}T${booking.timeSlot}:00+04:00`

    // 3. Create appointment
    await post(
      '/api/appointments',
      {
        doctor_id: booking.doctorId,
        service_id: booking.serviceId,
        starts_at: startsAt,
        promo_code: booking.promoCode.trim(),
      },
      access_token
    )

    booking.success = true
  } catch (err: unknown) {
    const status = (err as { status?: number })?.status
    const msg = (err as { data?: { error?: string } })?.data?.error ?? ''
    if (status === 401 || msg.toLowerCase().includes('otp') || msg.toLowerCase().includes('invalid or expired')) {
      bookingError.value = t('confirmCodeErr')
    } else if (msg.includes('taken') || msg.includes('conflict') || status === 409) {
      bookingError.value = t('confirmSlotTaken')
    } else {
      bookingError.value = t('confirmGenericErr') + config.public.clinicPhone
    }
  } finally {
    submittingBooking.value = false
  }
}

// Timezone-safe numeric date (locale-independent).
const formattedDate = computed(() => {
  if (!booking.date) return ''
  const [y, m, d] = booking.date.split('-')
  return `${d}.${m}.${y}`
})
</script>

<template>
  <div>
    <!-- Success state -->
    <div v-if="booking.success" class="text-center py-4">
      <div class="text-4xl mb-3">✅</div>
      <div class="text-base font-bold text-slate mb-1">{{ t('confirmSuccess') }}</div>
      <div class="text-sm text-muted mb-1">{{ formattedDate }}, {{ booking.timeSlot }}</div>
      <div v-if="booking.finalPrice > 0" class="text-sm text-slate font-semibold mb-1">
        {{ t('confirmToPay', { price: formatPrice(booking.finalPrice) }) }}
      </div>
      <div class="text-sm text-muted mb-6">{{ t('confirmSuccessFooter') }}</div>
      <button
        class="bg-primary text-white px-6 py-2.5 rounded-lg text-sm font-semibold"
        @click="booking.closeModal()"
      >
        {{ t('close') }}
      </button>
    </div>

    <!-- Form state -->
    <div v-else class="space-y-4">
      <!-- Summary -->
      <div class="bg-gray-50 rounded-xl p-4 text-sm space-y-2 text-slate">
        <div>📅 {{ formattedDate }}, {{ booking.timeSlot }}</div>
        <div v-if="booking.servicePrice > 0" class="flex items-center justify-between border-t border-border/60 pt-2">
          <span class="text-muted">{{ t('confirmCost') }}</span>
          <span class="font-semibold">
            <span
              v-if="booking.promoDiscountPct > 0"
              class="text-muted line-through mr-2 text-xs font-normal"
            >{{ formatPrice(booking.servicePrice) }}</span>
            {{ formatPrice(booking.finalPrice || booking.servicePrice) }}
          </span>
        </div>
      </div>

      <!-- Promo code -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">{{ t('confirmPromoLabel') }}</label>
        <div class="flex gap-2">
          <input
            v-model="booking.promoCode"
            type="text"
            placeholder="SUMMER10"
            class="flex-1 border border-border rounded-lg px-4 py-2.5 text-sm text-slate uppercase outline-none focus:border-primary"
            @keyup.enter="applyPromo"
          >
          <button
            type="button"
            class="px-4 py-2.5 rounded-lg text-sm font-semibold whitespace-nowrap transition-colors"
            :class="submittingPromo || !booking.serviceId
              ? 'bg-gray-100 text-muted cursor-not-allowed'
              : 'bg-primary text-white hover:bg-primary/90'"
            :disabled="submittingPromo || !booking.serviceId"
            @click="applyPromo"
          >
            {{ submittingPromo ? '...' : t('confirmPromoApply') }}
          </button>
        </div>
        <div v-if="booking.promoValid === true" class="text-xs text-emerald-600 mt-1">
          {{ t('confirmPromoValid', { pct: booking.promoDiscountPct }) }}
        </div>
        <div v-else-if="booking.promoValid === false" class="text-xs text-red-500 mt-1">
          {{ t('confirmPromoInvalid') }}
        </div>
      </div>

      <!-- Name -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">{{ t('confirmYourName') }}</label>
        <input
          v-model="booking.name"
          type="text"
          :placeholder="t('confirmNamePlaceholder')"
          class="w-full border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <!-- Phone + OTP send -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">{{ t('confirmPhone') }}</label>
        <div class="flex gap-2">
          <input
            :value="booking.phone"
            type="tel"
            inputmode="tel"
            autocomplete="tel"
            :placeholder="t('confirmPhonePlaceholder')"
            :disabled="booking.otpSent"
            class="flex-1 border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary disabled:bg-gray-50"
            :class="otpError ? 'border-red-400' : 'border-border'"
            @input="onPhoneInput"
          >
          <button
            class="px-4 py-2.5 rounded-lg text-sm font-semibold whitespace-nowrap transition-colors"
            :class="(cooldownSecs > 0 || submittingOtp || !consent)
              ? 'bg-gray-100 text-muted cursor-not-allowed'
              : 'bg-primary text-white hover:bg-primary/90'"
            :disabled="cooldownSecs > 0 || submittingOtp || !booking.phone || !consent"
            @click="sendOtp"
          >
            {{ cooldownSecs > 0 ? `${cooldownSecs}${t('secShort')}` : booking.otpSent ? t('confirmResend') : t('confirmGetCode') }}
          </button>
        </div>
        <div v-if="otpError" class="text-xs text-red-500 mt-1">{{ otpError }}</div>
      </div>

      <!-- FZ-152 consent -->
      <label v-if="!booking.otpSent" class="flex items-start gap-2 text-xs text-muted cursor-pointer">
        <input v-model="consent" type="checkbox" class="mt-0.5 flex-shrink-0 accent-primary">
        <span>
          {{ t('confirmConsent') }}
          <NuxtLink to="/privacy" target="_blank" class="text-primary underline" @click.stop>
            {{ t('confirmConsentLink') }}
          </NuxtLink>
          {{ t('confirmConsentLaw') }}
        </span>
      </label>

      <!-- OTP code input -->
      <Transition name="slide-down">
        <div v-if="booking.otpSent">
          <label class="text-xs font-semibold text-slate block mb-1">{{ t('confirmOtpLabel') }}</label>
          <input
            v-model="otpCode"
            type="text"
            inputmode="numeric"
            maxlength="6"
            placeholder="123456"
            class="w-full border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary tracking-widest text-center"
          >
          <div v-if="bookingError" class="text-xs text-red-500 mt-1">{{ bookingError }}</div>
        </div>
      </Transition>

      <!-- Submit -->
      <button
        v-if="booking.otpSent"
        class="w-full bg-primary text-white py-3 rounded-xl text-sm font-semibold transition-opacity"
        :class="(!otpCode || submittingBooking) ? 'opacity-50 cursor-not-allowed' : ''"
        :disabled="!otpCode || submittingBooking"
        @click="confirmBooking"
      >
        {{ submittingBooking ? t('confirmSubmitting') : t('confirmSubmit') }}
      </button>

      <div class="flex justify-start pt-1">
        <button class="text-sm text-muted hover:text-slate transition-colors" @click="booking.prevStep()">
          {{ t('back') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-down-enter-active { transition: all 0.2s ease; }
.slide-down-enter-from { opacity: 0; transform: translateY(-6px); }
</style>

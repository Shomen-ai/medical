<script setup lang="ts">
const booking = useBookingStore()
const { post } = useApi()
const config = useRuntimeConfig()

const submittingOtp = ref(false)
const submittingBooking = ref(false)
const otpCode = ref('')
const otpError = ref('')
const bookingError = ref('')
const cooldownDisplay = computed(() => booking.otpCooldownSecs)

const sendOtp = async () => {
  if (!booking.phone) return
  submittingOtp.value = true
  otpError.value = ''
  try {
    await post('/api/auth/otp', { phone: booking.phone })
    booking.otpSent = true
    booking.startOtpCooldown()
  } catch {
    otpError.value = 'Не удалось отправить код. Проверьте номер.'
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

    // 2. Create appointment
    await post(
      '/api/appointments',
      {
        doctor_id: booking.doctorId,
        service_id: booking.serviceId,
        starts_at: booking.timeSlot,
        promo_code: '',
      },
      access_token
    )

    booking.success = true
  } catch (err: unknown) {
    const msg = (err as { data?: { error?: string } })?.data?.error
    if (msg?.includes('invalid') || msg?.includes('OTP')) {
      bookingError.value = 'Неверный или истёкший код. Попробуйте ещё раз.'
    } else if (msg?.includes('taken') || msg?.includes('conflict')) {
      bookingError.value = 'Это время уже занято. Выберите другой слот.'
    } else {
      bookingError.value = 'Ошибка записи. Позвоните нам: ' + config.public.clinicPhone
    }
  } finally {
    submittingBooking.value = false
  }
}

const formatTime = (iso: string | null) =>
  iso ? new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' }) : ''
</script>

<template>
  <div>
    <!-- Success state -->
    <div v-if="booking.success" class="text-center py-4">
      <div class="text-4xl mb-3">✅</div>
      <div class="text-base font-bold text-slate mb-1">Вы записаны!</div>
      <div class="text-sm text-muted mb-1">{{ booking.date }}, {{ formatTime(booking.timeSlot) }}</div>
      <div class="text-sm text-muted mb-6">Ждём вас в клинике BeautyMed</div>
      <button
        class="bg-primary text-white px-6 py-2.5 rounded-lg text-sm font-semibold"
        @click="booking.closeModal()"
      >
        Закрыть
      </button>
    </div>

    <!-- Form state -->
    <div v-else class="space-y-4">
      <!-- Summary -->
      <div class="bg-gray-50 rounded-xl p-4 text-sm space-y-1 text-slate">
        <div>📅 {{ booking.date }}, {{ formatTime(booking.timeSlot) }}</div>
      </div>

      <!-- Name -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">Ваше имя</label>
        <input
          v-model="booking.name"
          type="text"
          placeholder="Иван Иванов"
          class="w-full border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <!-- Phone + OTP send -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">Номер телефона</label>
        <div class="flex gap-2">
          <input
            v-model="booking.phone"
            type="tel"
            placeholder="+79001234567"
            :disabled="booking.otpSent"
            class="flex-1 border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary disabled:bg-gray-50"
          >
          <button
            class="px-4 py-2.5 rounded-lg text-sm font-semibold whitespace-nowrap transition-colors"
            :class="cooldownDisplay > 0 || submittingOtp
              ? 'bg-gray-100 text-muted cursor-not-allowed'
              : 'bg-primary text-white hover:bg-primary/90'"
            :disabled="cooldownDisplay > 0 || submittingOtp || !booking.phone"
            @click="sendOtp"
          >
            {{ cooldownDisplay > 0 ? `${cooldownDisplay}с` : booking.otpSent ? 'Повтор' : 'Код' }}
          </button>
        </div>
        <div v-if="otpError" class="text-xs text-red-500 mt-1">{{ otpError }}</div>
      </div>

      <!-- OTP code input (appears after sending) -->
      <Transition name="slide-down">
        <div v-if="booking.otpSent">
          <label class="text-xs font-semibold text-slate block mb-1">Код из SMS</label>
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
        {{ submittingBooking ? 'Записываем...' : 'Подтвердить и записаться' }}
      </button>

      <!-- Back button -->
      <div class="flex justify-start pt-1">
        <button class="text-sm text-muted hover:text-slate transition-colors" @click="booking.prevStep()">
          ← Назад
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-down-enter-active { transition: all 0.2s ease; }
.slide-down-enter-from { opacity: 0; transform: translateY(-6px); }
</style>

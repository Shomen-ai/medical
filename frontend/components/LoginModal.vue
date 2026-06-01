<!--
  Файл: components/LoginModal.vue
  Назначение: модальное окно входа пациента по номеру телефона с подтверждением через SMS-код (OTP).
-->
<script setup lang="ts">
const emit = defineEmits<{ close: [] }>()
const auth = useAuthStore()
const router = useRouter()
const { t } = useI18n()

const phone = ref('')
const code = ref('')
const consent = ref(false)
const phoneError = ref('')

const phoneCheck = computed(() => validatePhone(phone.value))

// Форматирует ввод по маске +993 65 12-34-56 на каждый input.
const onPhoneInput = (e: Event) => {
  phone.value = formatTmPhone((e.target as HTMLInputElement).value)
  phoneError.value = ''
}

const handleSendOTP = () => {
  if (!consent.value) return
  if (!phoneCheck.value.valid) {
    phoneError.value = t('confirmPhoneErr')
    return
  }
  phoneError.value = ''
  auth.sendOTP(phoneCheck.value.e164, false)
}

const handleVerify = async () => {
  const role = await auth.verifyOTP(code.value, false)
  if (role) {
    emit('close')
    router.push('/cabinet')
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="emit('close')" />
      <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6">
        <!-- Header -->
        <div class="flex items-center justify-between mb-5">
          <h3 class="text-base font-bold text-slate">{{ t('loginTitle') }}</h3>
          <button type="button" class="text-muted hover:text-slate text-xl leading-none" @click="emit('close')">✕</button>
        </div>

        <!-- Phone step -->
        <div v-if="!auth.otpSent">
          <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('loginPhoneLabel') }}</label>
          <input
            :value="phone"
            type="tel"
            inputmode="tel"
            autocomplete="tel"
            :placeholder="t('confirmPhonePlaceholder')"
            class="w-full border rounded-xl px-4 py-2.5 text-sm mb-1 focus:outline-none focus:border-primary"
            :class="phoneError ? 'border-red-400' : 'border-border'"
            @keydown.enter="handleSendOTP"
            @input="onPhoneInput"
          >
          <div v-if="phoneError" class="text-xs text-red-500 mb-3">{{ phoneError }}</div>
          <div v-else class="mb-3" />
          <label class="flex items-start gap-2 text-xs text-muted mb-3 cursor-pointer">
            <input v-model="consent" type="checkbox" class="mt-0.5 flex-shrink-0 accent-primary">
            <span>
              {{ t('confirmConsent') }}
              <NuxtLink to="/privacy" target="_blank" class="text-primary underline" @click.stop>
                {{ t('confirmConsentLink') }}
              </NuxtLink>
              {{ t('confirmConsentLaw') }}
            </span>
          </label>
          <div v-if="auth.error" class="text-xs text-red-500 mb-3">{{ auth.error }}</div>
          <button
            type="button"
            class="w-full text-white py-2.5 rounded-xl text-sm font-bold transition-opacity"
            :class="(auth.loading || !consent) ? 'opacity-60 cursor-not-allowed' : 'hover:opacity-90'"
            style="background: linear-gradient(135deg, #005A5F, #00959D)"
            :disabled="auth.loading || !consent"
            @click="handleSendOTP"
          >
            {{ auth.loading ? t('loginSending') : t('loginGetCode') }}
          </button>
          <div class="mt-4 text-center">
            <NuxtLink to="/staff-login" class="text-xs text-muted underline" @click="emit('close')">
              {{ t('loginStaffLink') }}
            </NuxtLink>
          </div>
        </div>

        <!-- OTP step -->
        <div v-else>
          <p class="text-xs text-muted mb-4">{{ t('loginCodeSentTo', { phone: auth.phone }) }}</p>
          <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('loginCodeLabel') }}</label>
          <input
            v-model="code"
            type="text"
            inputmode="numeric"
            maxlength="6"
            class="w-full border border-border rounded-xl px-4 py-2.5 text-sm mb-4 text-center tracking-[0.5em] focus:outline-none focus:border-primary"
            @keydown.enter="handleVerify"
          >
          <div v-if="auth.error" class="text-xs text-red-500 mb-3">{{ auth.error }}</div>
          <button
            type="button"
            class="w-full text-white py-2.5 rounded-xl text-sm font-bold transition-opacity mb-2"
            :class="auth.loading ? 'opacity-60 cursor-not-allowed' : 'hover:opacity-90'"
            style="background: linear-gradient(135deg, #005A5F, #00959D)"
            :disabled="auth.loading"
            @click="handleVerify"
          >
            {{ auth.loading ? t('loginChecking') : t('loginSubmit') }}
          </button>
          <button type="button" class="w-full text-xs text-muted underline" @click="auth.otpSent = false">
            {{ t('loginChangeNumber') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const emit = defineEmits<{ close: [] }>()
const auth = useAuthStore()
const router = useRouter()

const phone = ref('')
const code = ref('')
const consent = ref(false)
const phoneError = ref('')

const phoneCheck = computed(() => validatePhone(phone.value))

const handleSendOTP = () => {
  if (!consent.value) return
  if (!phoneCheck.value.valid) {
    phoneError.value = 'Введите телефон в формате +7XXXXXXXXXX (Россия) или +993XXXXXXXX (Туркменистан)'
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
          <h3 class="text-base font-bold text-slate">Вход в личный кабинет</h3>
          <button type="button" class="text-muted hover:text-slate text-xl leading-none" @click="emit('close')">✕</button>
        </div>

        <!-- Phone step -->
        <div v-if="!auth.otpSent">
          <label class="text-xs font-semibold text-slate block mb-1.5">Номер телефона</label>
          <input
            v-model="phone"
            type="tel"
            inputmode="tel"
            autocomplete="tel"
            class="w-full border rounded-xl px-4 py-2.5 text-sm mb-1 focus:outline-none focus:border-primary"
            :class="phoneError ? 'border-red-400' : 'border-border'"
            @keydown.enter="handleSendOTP"
            @input="phoneError = ''"
          >
          <div v-if="phoneError" class="text-xs text-red-500 mb-3">{{ phoneError }}</div>
          <div v-else class="mb-3" />
          <label class="flex items-start gap-2 text-xs text-muted mb-3 cursor-pointer">
            <input v-model="consent" type="checkbox" class="mt-0.5 flex-shrink-0 accent-primary">
            <span>
              Я даю согласие на обработку информации о личной жизни в соответствии с
              <NuxtLink to="/privacy" target="_blank" class="text-primary underline" @click.stop>
                Политикой конфиденциальности
              </NuxtLink>
              (Закон Туркменистана от 20.03.2017).
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
            {{ auth.loading ? 'Отправляем...' : 'Получить код' }}
          </button>
          <div class="mt-4 text-center">
            <NuxtLink to="/staff-login" class="text-xs text-muted underline" @click="emit('close')">
              Вход для персонала клиники
            </NuxtLink>
          </div>
        </div>

        <!-- OTP step -->
        <div v-else>
          <p class="text-xs text-muted mb-4">Код отправлен на {{ auth.phone }}</p>
          <label class="text-xs font-semibold text-slate block mb-1.5">Код из SMS</label>
          <input
            v-model="code"
            type="text"
            inputmode="numeric"
            maxlength="6"
            placeholder="_ _ _ _ _ _"
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
            {{ auth.loading ? 'Проверяем...' : 'Войти' }}
          </button>
          <button type="button" class="w-full text-xs text-muted underline" @click="auth.otpSent = false">
            Изменить номер
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

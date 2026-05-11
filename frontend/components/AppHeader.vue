<script setup lang="ts">
const config = useRuntimeConfig()
const auth = useAuthStore()
const booking = useBookingStore()
const router = useRouter()
const { locale, setLocale, t } = useI18n()

const showLogin = ref(false)

onMounted(() => auth.init())

const handleCabinetClick = () => {
  if (auth.isLoggedIn) {
    if (auth.isAdmin) router.push('/admin')
    else if (auth.isDoctor) router.push('/doctor')
    else router.push('/cabinet')
  } else {
    showLogin.value = true
  }
}
</script>

<template>
  <header class="sticky top-0 z-50 bg-white border-b border-gray-100 shadow-sm">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 h-14 sm:h-16 flex items-center justify-between gap-3">
      <!-- Left: logo -->
      <NuxtLink to="/" class="flex items-center gap-2 flex-shrink-0">
        <span
          class="text-lg sm:text-xl font-extrabold"
          style="background: linear-gradient(135deg,#005A5F,#00959D); -webkit-background-clip:text; -webkit-text-fill-color:transparent;"
        >
          {{ config.public.clinicName }}
        </span>
        <span class="hidden sm:block text-sm text-gray-400 font-medium">· {{ t('cityName') }}</span>
      </NuxtLink>

      <!-- Center: hours (only desktop) -->
      <div class="hidden md:flex items-center justify-center gap-2 text-sm text-gray-500 flex-1">
        <svg class="w-4 h-4 text-primary flex-shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>
        </svg>
        <span>{{ t('workingHoursValue') }}</span>
      </div>

      <!-- Right: locale + buttons -->
      <div class="flex items-center gap-2 sm:gap-3 flex-shrink-0">
        <!-- Locale switcher -->
        <div class="flex items-center text-xs font-semibold text-muted border border-border rounded-full overflow-hidden">
          <button
            type="button"
            class="px-2.5 py-1 transition-colors"
            :class="locale === 'ru' ? 'bg-primary text-white' : 'hover:text-primary'"
            @click="setLocale('ru')"
          >
            RU
          </button>
          <button
            type="button"
            class="px-2.5 py-1 transition-colors"
            :class="locale === 'tk' ? 'bg-primary text-white' : 'hover:text-primary'"
            @click="setLocale('tk')"
          >
            TK
          </button>
        </div>
        <button
          type="button"
          class="hidden sm:block text-sm font-semibold text-white px-5 py-2.5 rounded-xl transition-opacity hover:opacity-90"
          style="background: linear-gradient(135deg, #005A5F, #00959D)"
          @click="booking.openModal()"
        >
          {{ t('bookShort') }}
        </button>
        <button
          type="button"
          class="text-xs sm:text-sm font-semibold text-primary border-2 border-primary px-3 sm:px-5 py-2 sm:py-2.5 rounded-xl hover:bg-primary/5 transition-colors whitespace-nowrap"
          @click="handleCabinetClick"
        >
          {{ auth.isLoggedIn ? t('cabinet') : t('signIn') }}
        </button>
      </div>
    </div>
  </header>

  <LoginModal v-if="showLogin" @close="showLogin = false" />
</template>

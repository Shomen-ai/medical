<script setup lang="ts">
const config = useRuntimeConfig()
const auth = useAuthStore()
const booking = useBookingStore()
const router = useRouter()

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
    <div class="max-w-6xl mx-auto px-6 h-16 grid grid-cols-3 items-center">
      <!-- Left: logo -->
      <NuxtLink to="/" class="flex items-center gap-2">
        <span
          class="text-xl font-extrabold"
          style="background: linear-gradient(135deg,#005A5F,#00959D); -webkit-background-clip:text; -webkit-text-fill-color:transparent;"
        >
          {{ config.public.clinicName }}
        </span>
        <span class="hidden sm:block text-sm text-gray-400 font-medium">· Туркменабад</span>
      </NuxtLink>

      <!-- Center: hours -->
      <div class="flex items-center justify-center gap-2 text-sm text-gray-500">
        <svg class="w-4 h-4 text-primary flex-shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>
        </svg>
        <span class="hidden md:block">{{ config.public.clinicHours }}</span>
      </div>

      <!-- Right: buttons -->
      <div class="flex items-center justify-end gap-3">
        <button
          type="button"
          class="hidden sm:block text-sm font-semibold text-white px-5 py-2.5 rounded-xl transition-opacity hover:opacity-90"
          style="background: linear-gradient(135deg, #005A5F, #00959D)"
          @click="booking.openModal()"
        >
          Записаться
        </button>
        <button
          type="button"
          class="text-sm font-semibold text-primary border-2 border-primary px-5 py-2.5 rounded-xl hover:bg-primary/5 transition-colors"
          @click="handleCabinetClick"
        >
          {{ auth.isLoggedIn ? 'Кабинет' : 'Войти' }}
        </button>
      </div>
    </div>
  </header>

  <LoginModal v-if="showLogin" @close="showLogin = false" />
</template>

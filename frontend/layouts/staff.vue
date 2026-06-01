<!--
  Файл: layouts/staff.vue
  Назначение: макет для внутреннего портала персонала (админ/доктор); скрывает публичную шапку и показывает служебную панель с переключателем языка, ролью и выходом.
-->
<script setup lang="ts">
const auth = useAuthStore()
const { locale, setLocale, t } = useI18n()

const roleLabel = computed(() => {
  if (auth.isAdmin) return t('staffRoleAdmin')
  if (auth.isDoctor) return t('staffRoleDoctor')
  return t('staffRoleDefault')
})
</script>

<template>
  <div>
    <!-- Staff header: gradient, no booking button -->
    <header style="background: linear-gradient(135deg, #005A5F, #00959D)">
      <div class="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between gap-3">
        <NuxtLink to="/" class="text-lg font-extrabold text-white tracking-tight">BeautyMed</NuxtLink>
        <span class="hidden sm:block text-sm font-semibold text-white/90">{{ roleLabel }}</span>
        <div class="flex items-center gap-3">
          <!-- Locale switcher -->
          <div class="flex items-center text-xs font-semibold text-white/80 border border-white/40 rounded-full overflow-hidden">
            <button
              type="button"
              class="px-2.5 py-1 transition-colors"
              :class="locale === 'ru' ? 'bg-white text-primary' : 'hover:text-white'"
              @click="setLocale('ru')"
            >
              RU
            </button>
            <button
              type="button"
              class="px-2.5 py-1 transition-colors"
              :class="locale === 'tk' ? 'bg-white text-primary' : 'hover:text-white'"
              @click="setLocale('tk')"
            >
              TK
            </button>
          </div>
          <button
            type="button"
            class="text-sm font-semibold text-white/80 hover:text-white transition-colors"
            @click="auth.logout(); navigateTo('/')"
          >
            {{ t('logout') }}
          </button>
        </div>
      </div>
    </header>
    <slot />
  </div>
</template>

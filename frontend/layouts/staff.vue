<!--
  Файл: layouts/staff.vue
  Назначение: макет для внутреннего портала персонала (админ/доктор); скрывает публичную шапку и показывает служебную панель с информацией о роли и выходом.
-->
<script setup lang="ts">
const auth = useAuthStore()
const route = useRoute()

const roleLabel = computed(() => {
  if (auth.isAdmin) return 'Кабинет администратора'
  if (auth.isDoctor) return 'Кабинет врача'
  return 'Кабинет'
})
</script>

<template>
  <div>
    <!-- Staff header: gradient, no booking button -->
    <header style="background: linear-gradient(135deg, #005A5F, #00959D)">
      <div class="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
        <NuxtLink to="/" class="text-lg font-extrabold text-white tracking-tight">BeautyMed</NuxtLink>
        <span class="text-sm font-semibold text-white/90">{{ roleLabel }}</span>
        <button
          type="button"
          class="text-sm font-semibold text-white/80 hover:text-white transition-colors"
          @click="auth.logout(); navigateTo('/')"
        >
          Выйти
        </button>
      </div>
    </header>
    <slot />
  </div>
</template>

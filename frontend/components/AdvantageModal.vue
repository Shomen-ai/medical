<!--
  Файл: components/AdvantageModal.vue
  Назначение: модальное окно с подробным описанием выбранного преимущества клиники, открываемое из секции преимуществ на главной.
-->
<script setup lang="ts">
import type { Advantage } from '~/data/advantages'

const props = defineProps<{ advantage: Advantage | null }>()
const emit = defineEmits<{ close: [] }>()

const booking = useBookingStore()
const { t, locale } = useI18n()

const handleBooking = () => {
  emit('close')
  booking.openModal()
}

// Close on Escape
onMounted(() => {
  const onKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape') emit('close')
  }
  window.addEventListener('keydown', onKey)
  onUnmounted(() => window.removeEventListener('keydown', onKey))
})

// Lock body scroll while modal is open
watch(() => props.advantage, (adv) => {
  if (!import.meta.client) return
  document.body.style.overflow = adv ? 'hidden' : ''
}, { immediate: true })

onUnmounted(() => {
  if (import.meta.client) document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="advantage"
        class="fixed inset-0 z-[60] flex items-center justify-center p-4"
      >
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />
        <div
          class="relative bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto"
          @click.stop
        >
          <!-- Close button -->
          <button
            type="button"
            aria-label="Закрыть"
            class="absolute top-3 right-3 w-9 h-9 flex items-center justify-center text-muted hover:text-slate hover:bg-gray-100 rounded-full transition-colors z-10"
            @click="emit('close')"
          >
            ✕
          </button>

          <!-- Header -->
          <div class="px-6 pt-7 pb-4 flex items-center gap-4">
            <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-3xl flex-shrink-0">
              {{ advantage.icon }}
            </div>
            <div>
              <h2 class="text-xl font-extrabold text-slate leading-tight">{{ advantage.title[locale] }}</h2>
              <p class="text-xs text-muted mt-0.5">{{ advantage.text[locale] }}</p>
            </div>
          </div>

          <!-- Body -->
          <div class="px-6 pb-4 space-y-3">
            <p
              v-for="(para, i) in advantage.body[locale]"
              :key="i"
              class="text-sm text-slate leading-relaxed"
            >
              {{ para }}
            </p>
          </div>

          <!-- Footer -->
          <div class="px-6 pb-6 pt-3 border-t border-border">
            <button
              type="button"
              class="w-full bg-primary text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors"
              @click="handleBooking"
            >
              {{ t('bookOnline') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active { transition: opacity 0.18s ease; }
.fade-enter-from,
.fade-leave-to { opacity: 0; }
</style>

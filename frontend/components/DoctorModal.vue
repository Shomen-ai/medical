<!--
  Файл: components/DoctorModal.vue
  Назначение: модальное окно с информацией о враче (фото, специальность, стаж, образование),
  открываемое по клику на фотографию врача в секции «Наши врачи».
-->
<script setup lang="ts">
import type { Doctor } from '~/types'

const props = defineProps<{ doctor: Doctor | null; photo: string }>()
const emit = defineEmits<{ close: []; book: [doctor: Doctor] }>()

const { t, tMed } = useI18n()

const handleBooking = () => {
  if (props.doctor) emit('book', props.doctor)
}

// Закрытие по Escape
onMounted(() => {
  const onKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape') emit('close')
  }
  window.addEventListener('keydown', onKey)
  onUnmounted(() => window.removeEventListener('keydown', onKey))
})

// Блокируем прокрутку body, пока модалка открыта
watch(() => props.doctor, (doc) => {
  if (!import.meta.client) return
  document.body.style.overflow = doc ? 'hidden' : ''
}, { immediate: true })

onUnmounted(() => {
  if (import.meta.client) document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="doctor"
        class="fixed inset-0 z-[60] flex items-center justify-center p-4"
      >
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />
        <div
          class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md max-h-[90vh] overflow-y-auto"
          @click.stop
        >
          <!-- Close button -->
          <button
            type="button"
            :aria-label="t('close')"
            class="absolute top-3 right-3 w-9 h-9 flex items-center justify-center text-white/90 hover:text-white hover:bg-black/20 rounded-full transition-colors z-10"
            @click="emit('close')"
          >
            ✕
          </button>

          <!-- Photo -->
          <div class="relative h-64 overflow-hidden rounded-t-2xl">
            <img :src="photo" :alt="doctor.full_name" class="w-full h-full object-cover object-top">
            <div
              class="absolute bottom-0 inset-x-0 px-5 pt-10 pb-4"
              style="background: linear-gradient(to top, rgba(0,0,0,0.7), transparent)"
            >
              <h2 class="text-xl font-extrabold text-white leading-tight">{{ tMed(doctor.full_name) }}</h2>
            </div>
          </div>

          <!-- Info -->
          <div class="px-6 py-5 space-y-4">
            <!-- Специальность -->
            <div class="flex items-start gap-3">
              <span class="text-lg flex-shrink-0">🩺</span>
              <div>
                <div class="text-[11px] uppercase tracking-wide text-muted font-semibold">{{ t('doctorsSpecialty') }}</div>
                <div class="text-sm font-semibold text-slate">{{ tMed(doctor.specialty_name) }}</div>
              </div>
            </div>

            <!-- Стаж -->
            <div class="flex items-start gap-3">
              <span class="text-lg flex-shrink-0">⏳</span>
              <div>
                <div class="text-[11px] uppercase tracking-wide text-muted font-semibold">{{ t('doctorsExperienceLabel') }}</div>
                <div class="text-sm font-semibold text-slate">{{ t('doctorsExperienceYears', { n: doctor.experience_years }) }}</div>
              </div>
            </div>

            <!-- Образование -->
            <div v-if="doctor.education" class="flex items-start gap-3">
              <span class="text-lg flex-shrink-0">🎓</span>
              <div>
                <div class="text-[11px] uppercase tracking-wide text-muted font-semibold">{{ t('doctorsEducation') }}</div>
                <div class="text-sm text-slate leading-snug">{{ tMed(doctor.education) }}</div>
              </div>
            </div>

            <!-- Био (если есть, переводим через tMed) -->
            <p v-if="doctor.bio" class="text-sm text-muted leading-relaxed pt-1 border-t border-border">
              {{ tMed(doctor.bio) }}
            </p>
          </div>

          <!-- Footer -->
          <div class="px-6 pb-6 pt-1">
            <button
              type="button"
              class="w-full text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:opacity-90 transition-opacity"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
              @click="handleBooking"
            >
              {{ t('bookShort') }}
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

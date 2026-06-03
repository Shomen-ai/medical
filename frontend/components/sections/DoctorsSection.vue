<!--
  Файл: components/sections/DoctorsSection.vue
  Назначение: секция главной страницы с карточками врачей клиники (фото, специализация, опыт) и быстрой записью к конкретному специалисту.
-->
<script setup lang="ts">
import type { Doctor, Service } from '~/types'
import { ServicesKey } from '~/composables/injectionKeys'

const props = defineProps<{ doctors: Doctor[] }>()
const booking = useBookingStore()
const { t, tMed } = useI18n()

// Services are provided by pages/index.vue; default to empty list if not present.
const services = inject(ServicesKey, computed(() => [] as Service[]))

// Pre-select specialty + doctor + first service of that specialty in code,
// then open the booking modal on the Date step.
const startBooking = (doc: Doctor) => {
  const svc = services.value.find(s => s.specialty_id === doc.specialty_id)
  if (!svc) {
    // No service available for this specialty — fall back to legacy flow
    // (opens at Doctor step with specialty pre-set).
    booking.openModal(doc.specialty_id)
    return
  }
  booking.openModalForDoctor(doc.specialty_id, doc.id, svc.id, svc.price)
}

const STOCK_PHOTOS = [
  '/doctors/doctor-1.jpg',
  '/doctors/doctor-2.jpg',
  '/doctors/doctor-3.jpg',
  '/doctors/doctor-4.jpg',
  '/doctors/doctor-5.jpg',
  '/doctors/doctor-6.jpg',
  '/doctors/doctor-7.jpg',
  '/doctors/doctor-8.jpg',
]

const photoSrc = (doc: Doctor, index: number) =>
  (doc.photo_url && doc.photo_url.trim()) || STOCK_PHOTOS[index % STOCK_PHOTOS.length]

// Модалка с информацией о враче (открывается по клику на фото).
const selectedDoctor = ref<Doctor | null>(null)
const selectedPhoto = ref('')

const openDoctor = (doc: Doctor, index: number) => {
  selectedDoctor.value = doc
  selectedPhoto.value = photoSrc(doc, index)
}

const bookFromModal = (doc: Doctor) => {
  selectedDoctor.value = null
  startBooking(doc)
}
</script>

<template>
  <section class="py-10 sm:py-14 bg-white">
    <div class="max-w-5xl mx-auto px-4 sm:px-8">
      <h2 class="text-2xl sm:text-3xl font-extrabold text-slate mb-6 sm:mb-10 text-center">{{ t('doctorsTitle') }}</h2>
      <div v-if="doctors.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 sm:gap-5">
        <div
          v-for="(doc, i) in doctors"
          :key="doc.id"
          class="bg-white rounded-2xl overflow-hidden shadow-card hover:shadow-card-lg transition-shadow"
        >
          <!-- Photo with specialty badge bar (clickable → doctor modal) -->
          <button
            type="button"
            class="relative h-36 sm:h-48 overflow-hidden w-full block group/photo cursor-pointer"
            :aria-label="doc.full_name"
            @click="openDoctor(doc, i)"
          >
            <img
              :src="photoSrc(doc, i)"
              :alt="doc.full_name"
              class="w-full h-full object-cover object-top transition-transform duration-300 group-hover/photo:scale-105"
            >
            <!-- Hover hint -->
            <div class="absolute inset-0 bg-black/0 group-hover/photo:bg-black/20 transition-colors flex items-center justify-center">
              <span class="opacity-0 group-hover/photo:opacity-100 transition-opacity text-white text-2xl">🔍</span>
            </div>
            <div
              class="absolute bottom-0 inset-x-0 py-1.5 text-center text-white text-[9px] font-bold uppercase tracking-wide"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
            >
              {{ tMed(doc.specialty_name) }}
            </div>
          </button>
          <!-- Info -->
          <div class="p-3 sm:p-4">
            <div class="text-sm font-bold text-slate leading-snug mb-1">{{ doc.full_name }}</div>
            <div class="text-xs text-muted mb-1.5">{{ t('doctorsExperience', { n: doc.experience_years }) }}</div>
            <div
              v-if="doc.education"
              class="text-[11px] leading-snug text-muted/90 mb-3 flex items-start gap-1"
              :title="`${t('doctorsEducation')}: ${tMed(doc.education)}`"
            >
              <span class="flex-shrink-0">🎓</span>
              <span class="line-clamp-2">{{ tMed(doc.education) }}</span>
            </div>
            <button
              type="button"
              class="w-full text-white text-[11px] font-bold py-2 rounded-lg transition-opacity hover:opacity-90"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
              @click="startBooking(doc)"
            >
              {{ t('bookShort') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Fallback when no doctors from backend -->
      <div v-if="doctors.length === 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 sm:gap-5">
        <div
          v-for="photo in STOCK_PHOTOS"
          :key="photo"
          class="bg-white rounded-2xl overflow-hidden shadow-card hover:shadow-card-lg transition-shadow"
        >
          <div class="relative h-36 sm:h-48 overflow-hidden">
            <img :src="photo" :alt="t('secDoctorAlt')" class="w-full h-full object-cover object-top">
            <div
              class="absolute bottom-0 inset-x-0 py-1.5 text-center text-white text-[9px] font-bold uppercase tracking-wide"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
            >
              {{ t('secClinicSpecialist') }}
            </div>
          </div>
          <div class="p-3 sm:p-4">
            <div class="text-sm font-bold text-slate mb-1">{{ t('secDoctorBeautyMed') }}</div>
            <div class="text-xs text-muted mb-3">{{ t('secExperiencedSpecialist') }}</div>
            <button
              type="button"
              class="w-full text-white text-[11px] font-bold py-2 rounded-lg hover:opacity-90 transition-opacity"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
              @click="booking.openModal()"
            >
              {{ t('bookShort') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Модалка с подробной информацией о враче (специальность, стаж, образование) -->
    <DoctorModal
      :doctor="selectedDoctor"
      :photo="selectedPhoto"
      @close="selectedDoctor = null"
      @book="bookFromModal"
    />
  </section>
</template>

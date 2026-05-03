<script setup lang="ts">
import type { Doctor } from '~/types'

const props = defineProps<{ doctors: Doctor[] }>()
const booking = useBookingStore()

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
</script>

<template>
  <section class="py-14 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-3xl font-extrabold text-slate mb-10 text-center">Наши врачи</h2>
      <div v-if="doctors.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-5">
        <div
          v-for="(doc, i) in doctors"
          :key="doc.id"
          class="bg-white rounded-2xl overflow-hidden shadow-card hover:shadow-card-lg transition-shadow"
        >
          <!-- Photo with specialty badge bar -->
          <div class="relative h-48 overflow-hidden">
            <img
              :src="photoSrc(doc, i)"
              :alt="doc.full_name"
              class="w-full h-full object-cover object-top"
            >
            <div
              class="absolute bottom-0 inset-x-0 py-1.5 text-center text-white text-[9px] font-bold uppercase tracking-wide"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
            >
              {{ doc.specialty_name }}
            </div>
          </div>
          <!-- Info -->
          <div class="p-4">
            <div class="text-sm font-bold text-slate leading-snug mb-1">{{ doc.full_name }}</div>
            <div class="text-xs text-muted mb-3">Стаж {{ doc.experience_years }} лет</div>
            <button
              type="button"
              class="w-full text-white text-[11px] font-bold py-2 rounded-lg transition-opacity hover:opacity-90"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
              @click="booking.openModal(doc.specialty_id)"
            >
              Записаться
            </button>
          </div>
        </div>
      </div>

      <!-- Fallback when no doctors from backend -->
      <div v-if="doctors.length === 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-5">
        <div
          v-for="photo in STOCK_PHOTOS"
          :key="photo"
          class="bg-white rounded-2xl overflow-hidden shadow-card hover:shadow-card-lg transition-shadow"
        >
          <div class="relative h-48 overflow-hidden">
            <img :src="photo" alt="Врач клиники" class="w-full h-full object-cover object-top">
            <div
              class="absolute bottom-0 inset-x-0 py-1.5 text-center text-white text-[9px] font-bold uppercase tracking-wide"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
            >
              Специалист клиники
            </div>
          </div>
          <div class="p-4">
            <div class="text-sm font-bold text-slate mb-1">Врач BeautyMed</div>
            <div class="text-xs text-muted mb-3">Опытный специалист</div>
            <button
              type="button"
              class="w-full text-white text-[11px] font-bold py-2 rounded-lg hover:opacity-90 transition-opacity"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
              @click="booking.openModal()"
            >
              Записаться
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

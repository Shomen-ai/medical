<script setup lang="ts">
import type { Doctor } from '~/types'

defineProps<{ doctors: Doctor[] }>()
const booking = useBookingStore()
</script>

<template>
  <section class="py-12 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Наши врачи</h2>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
        <div
          v-for="doc in doctors"
          :key="doc.id"
          class="bg-white border border-border rounded-xl overflow-hidden"
        >
          <!-- Photo with specialty badge -->
          <div class="relative h-40 bg-primary/10 flex items-center justify-center overflow-hidden">
            <img
              v-if="doc.photo_url"
              :src="doc.photo_url"
              :alt="doc.full_name"
              class="w-full h-full object-cover object-top"
            >
            <span v-else class="text-5xl">👩‍⚕️</span>
            <span class="absolute bottom-2 left-2 bg-primary text-white text-[9px] font-bold uppercase tracking-wide px-2 py-1 rounded-full">
              {{ doc.specialty_name }}
            </span>
          </div>
          <!-- Info -->
          <div class="p-3">
            <div class="text-xs font-semibold text-slate leading-snug mb-1">{{ doc.full_name }}</div>
            <div class="text-[11px] text-muted mb-3">Стаж {{ doc.experience_years }} лет</div>
            <button
              class="w-full bg-primary/10 text-primary text-[11px] font-semibold py-1.5 rounded-lg hover:bg-primary/20 transition-colors"
              @click="booking.openModal(doc.specialty_id)"
            >
              Записаться
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

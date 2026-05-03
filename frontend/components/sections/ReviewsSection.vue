<script setup lang="ts">
import { REVIEWS } from '~/data/reviews'
</script>

<template>
  <section class="py-12 bg-gray-50">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Отзывы пациентов</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <div
          v-for="review in REVIEWS"
          :key="review.author"
          class="bg-white rounded-xl border border-border p-5"
        >
          <div class="flex gap-0.5 mb-3">
            <span v-for="i in 5" :key="i" class="text-yellow-400 text-sm">
              {{ i <= review.rating ? '★' : '☆' }}
            </span>
          </div>
          <p class="text-sm text-slate leading-relaxed mb-4">{{ review.text }}</p>
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-muted">{{ review.author }}</span>
            <span class="text-xs text-muted">{{ review.date }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

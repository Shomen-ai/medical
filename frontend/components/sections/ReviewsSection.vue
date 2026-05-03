<script setup lang="ts">
import { REVIEWS } from '~/data/reviews'
</script>

<template>
  <section class="py-14 bg-[#F0FAFB]">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-3xl font-extrabold text-slate mb-10 text-center">Отзывы пациентов</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <div
          v-for="review in REVIEWS"
          :key="review.author"
          class="bg-white rounded-2xl shadow-card p-6 hover:shadow-card-lg transition-shadow"
        >
          <div class="text-5xl text-primary/20 font-serif leading-none mb-2 select-none">"</div>
          <div class="flex gap-0.5 mb-3">
            <span v-for="i in 5" :key="i" class="text-yellow-400 text-sm">
              {{ i <= review.rating ? '★' : '☆' }}
            </span>
          </div>
          <p class="text-sm text-slate leading-relaxed mb-5">{{ review.text }}</p>
          <div class="flex items-center justify-between pt-3 border-t border-border">
            <span class="text-xs font-bold text-slate">{{ review.author }}</span>
            <span class="text-xs text-muted">{{ review.date }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

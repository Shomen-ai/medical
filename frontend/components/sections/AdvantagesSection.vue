<script setup lang="ts">
import { ADVANTAGES } from '~/data/advantages'
</script>

<template>
  <section class="py-14 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-3xl font-extrabold text-slate mb-10 text-center">Почему выбирают нас</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-5">
        <div
          v-for="adv in ADVANTAGES"
          :key="adv.title"
          class="bg-white rounded-2xl shadow-card p-6 text-center hover:shadow-card-lg transition-shadow"
        >
          <div class="w-14 h-14 rounded-full bg-primary/10 flex items-center justify-center text-2xl mx-auto mb-4">
            {{ adv.icon }}
          </div>
          <div class="text-sm font-bold text-slate mb-1">{{ adv.title }}</div>
          <div class="text-xs text-muted leading-relaxed">{{ adv.text }}</div>
        </div>
      </div>
    </div>
  </section>
</template>

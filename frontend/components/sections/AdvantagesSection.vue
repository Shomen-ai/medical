<script setup lang="ts">
import { ADVANTAGES } from '~/data/advantages'
</script>

<template>
  <section class="py-12 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Почему выбирают нас</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div
          v-for="adv in ADVANTAGES"
          :key="adv.title"
          class="text-center"
        >
          <div class="text-3xl mb-3">{{ adv.icon }}</div>
          <div class="text-sm font-semibold text-slate mb-1">{{ adv.title }}</div>
          <div class="text-xs text-muted leading-relaxed">{{ adv.text }}</div>
        </div>
      </div>
    </div>
  </section>
</template>

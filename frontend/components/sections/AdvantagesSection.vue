<script setup lang="ts">
import { ADVANTAGES } from '~/data/advantages'
</script>

<template>
  <section class="py-14 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-3xl font-extrabold text-slate mb-10 text-center">Почему выбирают нас</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <NuxtLink
          v-for="adv in ADVANTAGES"
          :key="adv.slug"
          :to="`/advantages/${adv.slug}`"
          class="bg-white rounded-2xl shadow-card p-8 text-center hover:shadow-card-lg transition-shadow block"
        >
          <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-3xl mx-auto mb-5">
            {{ adv.icon }}
          </div>
          <div class="text-base font-bold text-slate mb-2">{{ adv.title }}</div>
          <div class="text-sm text-muted leading-relaxed">{{ adv.text }}</div>
          <div class="text-xs text-primary font-semibold mt-3">Подробнее →</div>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

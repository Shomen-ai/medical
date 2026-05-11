<script setup lang="ts">
import { ADVANTAGES, type Advantage } from '~/data/advantages'

const { t, locale } = useI18n()
const activeAdvantage = ref<Advantage | null>(null)
</script>

<template>
  <section class="py-10 sm:py-14 bg-white">
    <div class="max-w-5xl mx-auto px-4 sm:px-8">
      <h2 class="text-2xl sm:text-3xl font-extrabold text-slate mb-6 sm:mb-10 text-center">
        {{ t('advantagesTitle') }}
      </h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-6">
        <button
          v-for="adv in ADVANTAGES"
          :key="adv.slug"
          type="button"
          class="bg-white rounded-2xl shadow-card p-5 sm:p-8 text-center hover:shadow-card-lg transition-shadow block text-left w-full"
          @click="activeAdvantage = adv"
        >
          <div class="w-12 h-12 sm:w-16 sm:h-16 rounded-full bg-primary/10 flex items-center justify-center text-2xl sm:text-3xl mx-auto mb-3 sm:mb-5">
            {{ adv.icon }}
          </div>
          <div class="text-sm sm:text-base font-bold text-slate mb-1 sm:mb-2 text-center">
            {{ adv.title[locale] }}
          </div>
          <div class="text-xs sm:text-sm text-muted leading-relaxed text-center">
            {{ adv.text[locale] }}
          </div>
          <div class="text-xs text-primary font-semibold mt-2 sm:mt-3 text-center">{{ t('more') }}</div>
        </button>
      </div>
    </div>

    <AdvantageModal :advantage="activeAdvantage" @close="activeAdvantage = null" />
  </section>
</template>

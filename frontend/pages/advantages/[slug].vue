<script setup lang="ts">
import { ADVANTAGES } from '~/data/advantages'

const route = useRoute()
const booking = useBookingStore()
const { t, locale } = useI18n()

const advantage = computed(() => ADVANTAGES.find(a => a.slug === route.params.slug))

if (!advantage.value) {
  throw createError({ statusCode: 404, statusMessage: 'Not found', fatal: true })
}

useHead({
  title: computed(() => `${advantage.value!.title[locale.value]} — BeautyMed`),
  meta: [
    { name: 'description', content: computed(() => advantage.value!.body[locale.value][0].slice(0, 160)) },
    { property: 'og:title', content: computed(() => `${advantage.value!.title[locale.value]} — BeautyMed`) },
    { property: 'og:description', content: computed(() => advantage.value!.body[locale.value][0].slice(0, 200)) },
  ],
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-12">
    <NuxtLink to="/" class="inline-flex items-center text-sm text-muted hover:text-primary mb-6">
      {{ t('back') }}
    </NuxtLink>

    <div v-if="advantage" class="space-y-6">
      <div class="flex items-center gap-4">
        <div class="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center text-4xl flex-shrink-0">
          {{ advantage.icon }}
        </div>
        <div>
          <h1 class="text-3xl font-extrabold text-slate">{{ advantage.title[locale] }}</h1>
          <p class="text-sm text-muted mt-1">{{ advantage.text[locale] }}</p>
        </div>
      </div>

      <article class="prose prose-slate max-w-none space-y-4">
        <p v-for="(para, i) in advantage.body[locale]" :key="i" class="text-slate leading-relaxed">
          {{ para }}
        </p>
      </article>

      <div class="pt-6 border-t border-border">
        <button
          type="button"
          class="bg-primary text-white px-6 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors"
          @click="booking.openModal()"
        >
          {{ t('bookOnline') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
</script>

<template>
  <section class="py-14 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-3xl font-extrabold text-slate mb-10 text-center">Контакты</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- Info -->
        <div class="space-y-5">
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary flex-shrink-0 text-lg">
              📍
            </div>
            <div>
              <div class="text-sm font-bold text-slate mb-0.5">Адрес</div>
              <div class="text-sm text-muted">{{ config.public.clinicAddress }}</div>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary flex-shrink-0 text-lg">
              📞
            </div>
            <div>
              <div class="text-sm font-bold text-slate mb-0.5">Телефон</div>
              <a :href="`tel:${config.public.clinicPhone}`" class="text-sm text-primary font-semibold hover:underline">
                {{ config.public.clinicPhone }}
              </a>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary flex-shrink-0 text-lg">
              🕐
            </div>
            <div>
              <div class="text-sm font-bold text-slate mb-0.5">Часы работы</div>
              <div class="text-sm text-muted">{{ config.public.clinicHours }}</div>
            </div>
          </div>
        </div>
        <!-- Map placeholder -->
        <div class="rounded-2xl overflow-hidden shadow-card bg-gray-100 min-h-[200px] flex items-center justify-center">
          <span class="text-muted text-sm">Карта (добавить Yandex Maps iframe)</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
</script>

<template>
  <section class="py-12 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Контакты</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- Info -->
        <div class="space-y-4">
          <div class="flex items-start gap-3">
            <span class="text-primary mt-0.5">📍</span>
            <div>
              <div class="text-sm font-semibold text-slate mb-0.5">Адрес</div>
              <div class="text-sm text-muted">{{ config.public.clinicAddress }}</div>
            </div>
          </div>
          <div class="flex items-start gap-3">
            <span class="text-primary mt-0.5">📞</span>
            <div>
              <div class="text-sm font-semibold text-slate mb-0.5">Телефон</div>
              <a :href="`tel:${config.public.clinicPhone}`" class="text-sm text-primary font-medium">
                {{ config.public.clinicPhone }}
              </a>
            </div>
          </div>
          <div class="flex items-start gap-3">
            <span class="text-primary mt-0.5">🕐</span>
            <div>
              <div class="text-sm font-semibold text-slate mb-0.5">Часы работы</div>
              <div class="text-sm text-muted">{{ config.public.clinicHours }}</div>
            </div>
          </div>
        </div>
        <!-- Map placeholder -->
        <div class="rounded-xl overflow-hidden border border-border bg-gray-100 min-h-[200px] flex items-center justify-center">
          <span class="text-muted text-sm">Карта (добавить Yandex Maps iframe)</span>
        </div>
      </div>
    </div>
  </section>
</template>

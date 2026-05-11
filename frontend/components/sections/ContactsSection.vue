<script setup lang="ts">
const config = useRuntimeConfig()
</script>

<template>
  <section class="py-10 sm:py-14 bg-white">
    <div class="max-w-5xl mx-auto px-4 sm:px-8">
      <h2 class="text-2xl sm:text-3xl font-extrabold text-slate mb-6 sm:mb-10 text-center">Контакты</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8">
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
        <!-- Map -->
        <div class="rounded-2xl overflow-hidden shadow-card bg-gray-100 min-h-[200px]">
          <iframe
            src="https://maps.google.com/maps?q=Parahat%2025%2F31%20Turkmenabat&t=&z=14&ie=UTF8&iwloc=&output=embed"
            class="w-full h-full min-h-[200px] border-0"
            loading="lazy"
            referrerpolicy="no-referrer-when-downgrade"
            title="Карта клиники BeautyMed"
          />
        </div>
      </div>
      <div class="mt-8 pt-6 border-t border-border text-center text-xs text-muted">
        © {{ new Date().getFullYear() }} BeautyMed ·
        <NuxtLink to="/privacy" class="underline hover:text-primary">Политика конфиденциальности</NuxtLink>
      </div>
    </div>
  </section>
</template>

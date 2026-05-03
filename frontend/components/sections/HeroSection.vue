<script setup lang="ts">
const booking = useBookingStore()
</script>

<template>
  <section class="bg-white border-b border-border">
    <div class="max-w-5xl mx-auto">
      <!-- Main hero row -->
      <div class="flex items-stretch min-h-[280px]">
        <!-- Text block -->
        <div class="flex-1 px-8 py-10 flex flex-col justify-center">
          <p class="text-xs font-semibold tracking-[2.5px] text-primary uppercase mb-3">
            Клиника красоты и здоровья · Хабаровск
          </p>
          <h1 class="text-3xl font-extrabold text-slate leading-tight mb-3">
            Красота и здоровье<br>в надёжных руках
          </h1>
          <p class="text-muted text-sm leading-relaxed mb-6">
            Косметология, дерматология,<br>трихология и эстетическая медицина
          </p>
          <div class="flex gap-3">
            <button
              class="bg-primary text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors"
              @click="booking.openModal()"
            >
              Записаться онлайн
            </button>
            <a
              href="#services"
              class="border border-primary text-primary px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/5 transition-colors"
            >
              Наши услуги
            </a>
          </div>
        </div>
        <!-- Photo block -->
        <div class="w-[300px] flex-shrink-0 overflow-hidden">
          <img
            src="/clinic_3.png"
            alt="Клиника BeautyMed"
            class="w-full h-full object-cover object-top"
          >
        </div>
      </div>
      <!-- Stats row -->
      <div class="flex border-t border-border bg-gray-50">
        <div class="flex-1 py-3 text-center border-r border-border">
          <span class="block text-xl font-extrabold text-primary">10+</span>
          <span class="text-xs text-slate">лет работы</span>
        </div>
        <div class="flex-1 py-3 text-center border-r border-border">
          <span class="block text-xl font-extrabold text-primary">8</span>
          <span class="text-xs text-slate">врачей</span>
        </div>
        <div class="flex-1 py-3 text-center">
          <span class="block text-xl font-extrabold text-primary">3000+</span>
          <span class="text-xs text-slate">пациентов</span>
        </div>
      </div>
    </div>
  </section>
</template>

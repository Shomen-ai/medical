<script setup lang="ts">
const booking = useBookingStore()
</script>

<template>
  <section style="background: linear-gradient(135deg, #005A5F 0%, #00959D 100%)">
    <div class="max-w-5xl mx-auto">
      <!-- Main hero row -->
      <div class="flex items-stretch min-h-[280px]">
        <!-- Text block -->
        <div class="flex-1 px-8 py-10 flex flex-col justify-center">
          <div class="inline-flex items-center bg-white/15 text-white text-xs font-semibold tracking-[2px] uppercase px-3 py-1.5 rounded-full mb-4 w-fit">
            Клиника красоты · Ульяновск
          </div>
          <h1 class="text-3xl font-black text-white leading-tight mb-3">
            Красота и здоровье<br>в надёжных руках
          </h1>
          <p class="text-white/75 text-sm leading-relaxed mb-6">
            Косметология, дерматология,<br>трихология и эстетическая медицина
          </p>
          <div class="flex gap-3">
            <button
              type="button"
              class="bg-white text-[#005A5F] px-5 py-2.5 rounded-lg text-sm font-bold shadow-lg hover:shadow-xl transition-shadow"
              @click="booking.openModal()"
            >
              Записаться онлайн
            </button>
            <a
              href="#services"
              class="border border-white/50 text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-white/10 transition-colors"
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
      <div class="flex bg-black/20 border-t border-white/10">
        <div class="flex-1 py-3 text-center border-r border-white/10">
          <span class="block text-xl font-extrabold text-white">10+</span>
          <span class="text-[11px] text-white/55 uppercase tracking-wide">лет работы</span>
        </div>
        <div class="flex-1 py-3 text-center border-r border-white/10">
          <span class="block text-xl font-extrabold text-white">8</span>
          <span class="text-[11px] text-white/55 uppercase tracking-wide">врачей</span>
        </div>
        <div class="flex-1 py-3 text-center">
          <span class="block text-xl font-extrabold text-white">3000+</span>
          <span class="text-[11px] text-white/55 uppercase tracking-wide">пациентов</span>
        </div>
      </div>
    </div>
  </section>
</template>

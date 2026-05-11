<script setup lang="ts">
const booking = useBookingStore()

const scrollToServices = () => {
  const el = document.getElementById('services')
  if (!el) return
  const headerH = document.querySelector('header')?.offsetHeight ?? 64
  const top = el.getBoundingClientRect().top + window.scrollY - headerH - 8
  window.scrollTo({ top, behavior: 'smooth' })
}
</script>

<template>
  <section style="background: linear-gradient(135deg, #005A5F 0%, #00959D 100%)">
    <div class="max-w-5xl mx-auto">
      <!-- Main hero row -->
      <div class="flex items-stretch min-h-[280px]">
        <!-- Text block -->
        <div class="flex-1 px-8 py-10 flex flex-col justify-center">
          <div class="inline-flex items-center bg-white/15 text-white text-xs font-semibold tracking-[2px] uppercase px-3 py-1.5 rounded-full mb-4 w-fit">
            Клиника красоты · Туркменабад
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
            <button
              type="button"
              class="border border-white/50 text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-white/10 transition-colors"
              @click="scrollToServices"
            >
              Наши услуги
            </button>
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
    </div>
  </section>
</template>

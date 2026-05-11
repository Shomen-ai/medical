<script setup lang="ts">
const booking = useBookingStore()

const stepComponents = {
  1: resolveComponent('StepSpecialty'),
  2: resolveComponent('StepDoctor'),
  3: resolveComponent('StepDate'),
  4: resolveComponent('StepTime'),
  5: resolveComponent('StepConfirm'),
}

const stepTitles: Record<number, string> = {
  1: 'Специальность',
  2: 'Врач',
  3: 'Дата',
  4: 'Время',
  5: 'Подтверждение',
}

const handleNext = () => {
  if (booking.step < 5 && booking.canProceed) booking.nextStep()
}

const handleBack = () => booking.prevStep()

// Close on Escape
const escHandler = (e: KeyboardEvent) => { if (e.key === 'Escape') booking.closeModal() }
onMounted(() => window.addEventListener('keydown', escHandler))
onUnmounted(() => window.removeEventListener('keydown', escHandler))
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="booking.open"
        class="fixed inset-0 z-50 flex items-end sm:items-center justify-center sm:p-4"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/40 backdrop-blur-sm"
          @click="booking.closeModal()"
        />

        <!-- Modal box -->
        <div class="relative bg-white rounded-t-2xl sm:rounded-2xl shadow-2xl w-full max-w-2xl max-h-[92vh] sm:max-h-[90vh] flex flex-col overflow-hidden">
          <!-- Header -->
          <div class="flex items-center justify-between px-4 sm:px-6 pt-4 sm:pt-5 pb-0">
            <span class="text-sm font-bold text-slate">Онлайн-запись</span>
            <button class="text-muted hover:text-slate transition-colors text-lg p-1 -m-1" @click="booking.closeModal()">✕</button>
          </div>

          <!-- Body -->
          <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
            <!-- Dot indicator -->
            <div class="flex items-center gap-1.5 mb-5">
              <div
                v-for="n in 5"
                :key="n"
                class="h-2 rounded-full transition-all duration-300"
                :class="{
                  'bg-primary w-6': n === booking.step,
                  'bg-primary w-2': n < booking.step,
                  'bg-border w-2': n > booking.step,
                }"
              />
            </div>

            <!-- Step title -->
            <div class="text-xs font-semibold text-muted uppercase tracking-wide mb-3">
              Шаг {{ booking.step }} — {{ stepTitles[booking.step] }}
            </div>

            <!-- Dynamic step component -->
            <component :is="stepComponents[booking.step as keyof typeof stepComponents]" />
          </div>

          <!-- Footer (hidden on step 5, which has its own submit) -->
          <div v-if="booking.step < 5" class="px-4 sm:px-6 py-3 sm:py-4 border-t border-border flex justify-between items-center">
            <button
              v-if="booking.step > 1"
              class="text-sm text-muted hover:text-slate transition-colors"
              @click="handleBack"
            >
              ← Назад
            </button>
            <span v-else />
            <button
              class="bg-primary text-white px-5 py-2 rounded-lg text-sm font-semibold transition-opacity"
              :class="booking.canProceed ? 'opacity-100' : 'opacity-40 cursor-not-allowed'"
              :disabled="!booking.canProceed"
              @click="handleNext"
            >
              Далее →
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>

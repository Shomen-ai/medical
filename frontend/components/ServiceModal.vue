<!--
  Файл: components/ServiceModal.vue
  Назначение: модальное окно с подробным описанием услуги (специальность, цена, длительность,
  развёрнутый текст), открываемое по клику на «Подробнее» в секции «Наши услуги».
-->
<script setup lang="ts">
import type { Service } from '~/types'
import { SERVICE_DESCRIPTIONS } from '~/data/serviceDescriptions'

const props = defineProps<{ service: Service | null; specialtyName: string }>()
const emit = defineEmits<{ close: []; book: [service: Service] }>()

const { t, tMed, locale } = useI18n()

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'TMT', maximumFractionDigits: 0 }).format(price)

// Подробное описание из словаря (ru/tk); фолбэк на короткое из БД.
const desc = computed(() => {
  if (!props.service) return ''
  const rich = SERVICE_DESCRIPTIONS[props.service.name]?.[locale.value]
  return rich || (props.service.description ? tMed(props.service.description) : '')
})

const handleBook = () => {
  if (props.service) emit('book', props.service)
}

// Закрытие по Escape
onMounted(() => {
  const onKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape') emit('close')
  }
  window.addEventListener('keydown', onKey)
  onUnmounted(() => window.removeEventListener('keydown', onKey))
})

// Блокируем прокрутку body, пока модалка открыта
watch(() => props.service, (svc) => {
  if (!import.meta.client) return
  document.body.style.overflow = svc ? 'hidden' : ''
}, { immediate: true })

onUnmounted(() => {
  if (import.meta.client) document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="service"
        class="fixed inset-0 z-[60] flex items-center justify-center p-4"
      >
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />
        <div
          class="relative bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto"
          @click.stop
        >
          <!-- Close button -->
          <button
            type="button"
            :aria-label="t('close')"
            class="absolute top-3 right-3 w-9 h-9 flex items-center justify-center text-muted hover:text-slate hover:bg-black/5 rounded-full transition-colors z-10"
            @click="emit('close')"
          >
            ✕
          </button>

          <!-- Header + meta -->
          <div class="px-6 pt-6 pb-2">
            <div v-if="specialtyName" class="text-xs uppercase tracking-wide text-primary font-bold mb-2">
              {{ tMed(specialtyName) }}
            </div>
            <h2 class="text-2xl font-extrabold text-slate leading-tight mb-3 pr-8">{{ tMed(service.name) }}</h2>
            <div class="flex items-center gap-4 flex-wrap">
              <span class="text-xl font-bold text-primary">{{ formatPrice(service.price) }}</span>
              <span class="text-sm text-muted">{{ t('svcDuration', { n: service.duration_min }) }}</span>
            </div>
          </div>

          <!-- Description -->
          <div class="px-6 py-4">
            <p class="text-sm text-slate leading-relaxed whitespace-pre-line">{{ desc }}</p>
          </div>

          <!-- Footer -->
          <div class="px-6 pb-6 pt-1">
            <button
              type="button"
              class="w-full text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:opacity-90 transition-opacity"
              style="background: linear-gradient(135deg, #005A5F, #00959D)"
              @click="handleBook"
            >
              {{ t('bookShort') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active { transition: opacity 0.18s ease; }
.fade-enter-from,
.fade-leave-to { opacity: 0; }
</style>

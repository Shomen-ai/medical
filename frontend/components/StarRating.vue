<!--
  Файл: components/StarRating.vue
  Назначение: компонент звёздного рейтинга 1–5. Режим только-чтение (отображение) и
  интерактивный (выбор через v-model).
-->
<script setup lang="ts">
const props = withDefaults(defineProps<{ modelValue?: number; readonly?: boolean; size?: string }>(), {
  modelValue: 0,
  readonly: false,
  size: 'text-xl',
})
const emit = defineEmits<{ 'update:modelValue': [value: number] }>()

const set = (n: number) => { if (!props.readonly) emit('update:modelValue', n) }
</script>

<template>
  <div class="flex gap-0.5">
    <button
      v-for="i in 5"
      :key="i"
      type="button"
      :disabled="readonly"
      :aria-label="`${i}`"
      :class="[
        size,
        i <= modelValue ? 'text-yellow-400' : 'text-gray-300',
        readonly ? 'cursor-default' : 'cursor-pointer hover:scale-110 transition-transform',
      ]"
      @click="set(i)"
    >
      {{ i <= modelValue ? '★' : '☆' }}
    </button>
  </div>
</template>

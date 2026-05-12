// Файл: plugins/chartjs.client.ts
// Назначение: клиентский Nuxt-плагин — регистрирует все компоненты Chart.js для графиков статистики в админ-панели.
import { Chart, registerables } from 'chart.js'

// Экспорт Nuxt-плагина, выполняющего регистрацию модулей Chart.js на старте клиентского рантайма.
export default defineNuxtPlugin(() => {
  Chart.register(...registerables)
})

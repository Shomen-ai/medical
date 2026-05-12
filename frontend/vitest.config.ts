// Файл: vitest.config.ts
// Назначение: конфигурация Vitest для юнит-тестов — исключает e2e-сценарии Playwright и node_modules.
import { defineConfig } from 'vitest/config'

// Экспорт конфигурации Vitest.
export default defineConfig({
  test: {
    exclude: ['tests/e2e/**', 'node_modules/**'],
  },
})

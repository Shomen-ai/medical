// Файл: playwright.config.ts
// Назначение: конфигурация Playwright для e2e-тестов — каталог tests/e2e, базовый URL dev-сервера, проект Chromium и автоматический запуск `npm run dev`.
import { defineConfig, devices } from '@playwright/test'

// Экспорт конфигурации Playwright.
export default defineConfig({
  testDir: './tests/e2e',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: true,
  },
})

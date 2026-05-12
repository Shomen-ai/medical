// Файл: tests/e2e/booking.spec.ts
// Назначение: e2e-тесты Playwright для главной страницы и модального окна онлайн-записи — открытие модалки, закрытие по бэкдропу/Escape, дизейбл кнопки «Далее» на шаге 1.
import { test, expect, type Page } from '@playwright/test'

// The modal backdrop container is a fixed overlay with z-50.
// AdvantagesSection also contains "Онлайн-запись" text which is always in the DOM.
// We scope modal assertions to the modal container to avoid strict-mode violations.
const modal = (page: Page) => page.locator('.fixed.inset-0.z-50')

test('homepage loads with hero section', async ({ page }) => {
  await page.goto('/')
  await expect(page.locator('h1')).toContainText('Красота и здоровье')
  await expect(page.locator('text=Записаться онлайн').first()).toBeVisible()
})

test('booking modal opens when CTA clicked', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  await expect(modal(page)).toBeVisible()
  await expect(modal(page).locator('text=Онлайн-запись')).toBeVisible()
  await expect(modal(page).locator('text=Шаг 1')).toBeVisible()
})

test('booking modal closes on backdrop click', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  await expect(modal(page)).toBeVisible()
  // Click backdrop (outside modal box)
  await page.mouse.click(10, 10)
  await expect(modal(page)).not.toBeVisible()
})

test('booking modal closes on Escape', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  await expect(modal(page)).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(modal(page)).not.toBeVisible()
})

test('step 1 Next button disabled until specialty and service selected', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  const nextBtn = page.locator('button', { hasText: 'Далее →' })
  await expect(nextBtn).toBeDisabled()
})

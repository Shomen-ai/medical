// Аудит локализации: ставит locale=tk и снимает ключевые экраны, чтобы найти
// непереведённые (русские) фрагменты в туркменском режиме.
import { chromium } from 'playwright'
import { execSync } from 'child_process'
import { resolve, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const DIR = resolve(__dirname, '../../screenshots/tk-audit')
execSync(`mkdir -p '${DIR}'`)
const BASE = process.env.BASE_URL || 'http://85.198.80.245'
const SSH = `sshpass -p 'vZdfA-cj3jHZX8' ssh -o StrictHostKeyChecking=no -o ConnectTimeout=12 root@85.198.80.245`

const browser = await chromium.launch({ headless: true })
const ctx = await browser.newContext({ viewport: { width: 1280, height: 900 }, locale: 'ru-RU' })
// заранее ставим cookie локали = tk
const host = new URL(BASE).hostname
await ctx.addCookies([{ name: 'beautymed_locale', value: 'tk', domain: host, path: '/' }])
const page = await ctx.newPage()
const shot = (n, full = false) => page.screenshot({ path: resolve(DIR, `${n}.png`), fullPage: full })

async function main() {
  // 1. Главная (TK)
  await page.goto(`${BASE}/`, { waitUntil: 'networkidle' }); await page.waitForTimeout(900)
  await shot('01-home', true)

  // 2. Страница услуги (берём первый id услуги из API)
  const svc = await page.evaluate(async (b) => (await fetch(b + '/api/services').then(r => r.json()))[0], BASE)
  if (svc?.id) {
    await page.goto(`${BASE}/services/${svc.id}`, { waitUntil: 'networkidle' }); await page.waitForTimeout(700)
    await shot('02-service-detail', true)
  }

  // 3. Мастер записи — шаги
  await page.goto(`${BASE}/`, { waitUntil: 'networkidle' }); await page.waitForTimeout(500)
  await page.locator('button', { hasText: 'Bellige durmak' }).first().click().catch(()=>{})
  await page.waitForTimeout(800)
  await shot('03-booking-step1')

  // 4. Страница преимущества (advantages/[slug])
  await page.goto(`${BASE}/`, { waitUntil: 'networkidle' }); await page.waitForTimeout(400)
  const adv = page.locator('a[href*="/advantages/"]').first()
  if (await adv.count()) {
    const href = await adv.getAttribute('href')
    await page.goto(`${BASE}${href}`, { waitUntil: 'networkidle' }); await page.waitForTimeout(500)
    await shot('04-advantage', true)
  }

  // 5. Кабинет врача (TK)
  await page.goto(`${BASE}/staff-login`, { waitUntil: 'networkidle' }); await page.waitForTimeout(400)
  await page.locator('input[type=text]').fill('doctor1')
  await page.locator('input[type=password]').fill('doctor123')
  await page.locator('button[type=submit]').click()
  await page.waitForURL('**/doctor', { timeout: 15000 }); await page.waitForTimeout(2000)
  await shot('05-doctor', true)
  await ctx.clearCookies(); await ctx.addCookies([{ name: 'beautymed_locale', value: 'tk', domain: host, path: '/' }])

  // 6. Админ (TK)
  await page.goto(`${BASE}/staff-login`, { waitUntil: 'networkidle' }); await page.waitForTimeout(400)
  await page.locator('input[type=text]').fill('admin')
  await page.locator('input[type=password]').fill('admin123')
  await page.locator('button[type=submit]').click()
  await page.waitForURL('**/admin', { timeout: 15000 }); await page.waitForTimeout(2500)
  await shot('06-admin', true)

  // 7. Кабинет пациента (OTP)
  const ctx2 = await browser.newContext({ viewport: { width: 1280, height: 900 }, locale: 'ru-RU' })
  await ctx2.addCookies([{ name: 'beautymed_locale', value: 'tk', domain: host, path: '/' }])
  const pp = await ctx2.newPage()
  try {
    await pp.goto(`${BASE}/`, { waitUntil: 'networkidle' }); await pp.waitForTimeout(700)
    await pp.locator('button', { hasText: 'Girmek' }).first().click()
    await pp.waitForTimeout(600)
    await pp.locator('input[type=tel]').fill('65200001')
    await pp.locator('input[type=checkbox]').first().check()
    await pp.locator('button', { hasText: 'Kod almak' }).click()
    await pp.waitForTimeout(1500)
    const keys = execSync(`${SSH} "docker exec beautymed-redis-1 redis-cli --scan --pattern 'otp:*'"`).toString().trim().split('\n').map(s=>s.trim()).filter(Boolean)
    const key = keys.find(k => k.includes('65200001')) || keys[0]
    const code = execSync(`${SSH} "docker exec beautymed-redis-1 redis-cli GET '${key}'"`).toString().trim()
    await pp.locator('input[inputmode=numeric], input[type=text]').last().fill(code)
    await pp.locator('button', { hasText: 'Girmek' }).last().click()
    await pp.waitForURL('**/cabinet', { timeout: 15000 }); await pp.waitForTimeout(1800)
    await pp.screenshot({ path: resolve(DIR, '07-patient.png'), fullPage: true })
    // профиль
    await pp.goto(`${BASE}/cabinet/profile`, { waitUntil: 'networkidle' }); await pp.waitForTimeout(1200)
    await pp.screenshot({ path: resolve(DIR, '08-patient-profile.png'), fullPage: true })
  } catch(e) { console.log('patient err', e.message) }
  await ctx2.close()

  console.log('done →', DIR)
}
main().catch(e => { console.error(e); process.exit(1) }).finally(() => browser.close())

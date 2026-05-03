# BeautyMed Frontend Public Site — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the BeautyMed public-facing website as a Nuxt 3 SSR app with a 5-step booking modal.

**Architecture:** Nuxt 3 SSR fetches specialties, doctors, and services server-side from the existing Go backend (all endpoints are already public, no auth required). Static content (advantages, reviews, contacts) lives in frontend data files. The booking modal runs fully client-side using the existing OTP → JWT → appointment flow.

**Tech Stack:** Nuxt 3, TypeScript, Pinia, Tailwind CSS v3, Vitest, Playwright

**Important — no backend changes for content:** `GET /api/specialties`, `GET /api/doctors`, `GET /api/services`, and `GET /api/doctors/:id/slots` already exist and require no auth. The only required sequence involving auth is the booking itself: `POST /api/auth/otp` → `POST /api/auth/verify` (returns `access_token`) → `POST /api/appointments` (Bearer token).

---

## File Map

**Create:**
```
frontend/
├── package.json
├── nuxt.config.ts
├── tailwind.config.ts
├── app.vue
├── types/index.ts
├── stores/booking.ts
├── composables/useApi.ts
├── data/specialtyMeta.ts
├── data/reviews.ts
├── data/advantages.ts
├── components/sections/HeroSection.vue
├── components/sections/AdvantagesSection.vue
├── components/sections/ServicesSection.vue
├── components/sections/DoctorsSection.vue
├── components/sections/ReviewsSection.vue
├── components/sections/ContactsSection.vue
├── components/booking/BookingModal.vue
├── components/booking/steps/StepSpecialty.vue
├── components/booking/steps/StepDoctor.vue
├── components/booking/steps/StepDate.vue
├── components/booking/steps/StepTime.vue
├── components/booking/steps/StepConfirm.vue
├── pages/index.vue
├── public/clinic_3.png  (copied from project root)
├── public/clinic_2.png  (copied from project root)
├── tests/booking.store.test.ts
├── tests/e2e/booking.spec.ts
├── playwright.config.ts
└── Dockerfile
```

**Modify:**
- `docker-compose.yml` — add `frontend` service
- `nginx/nginx.conf` — add `location /` → frontend, `location /health` → api

---

## Task 1: Nuxt 3 Project Scaffold

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/nuxt.config.ts`
- Create: `frontend/tailwind.config.ts`
- Create: `frontend/app.vue`

- [ ] **Step 1: Create `frontend/package.json`**

```json
{
  "name": "beautymed-frontend",
  "private": true,
  "type": "module",
  "scripts": {
    "build": "nuxt build",
    "dev": "nuxt dev",
    "typecheck": "nuxi typecheck",
    "test": "vitest run",
    "test:e2e": "playwright test"
  },
  "dependencies": {
    "@pinia/nuxt": "^0.9.0",
    "nuxt": "^3.15.0",
    "pinia": "^2.3.0",
    "vue": "^3.5.0"
  },
  "devDependencies": {
    "@nuxt/test-utils": "^3.15.0",
    "@nuxtjs/tailwindcss": "^6.13.0",
    "@playwright/test": "^1.50.0",
    "@vue/test-utils": "^2.4.6",
    "typescript": "^5.7.0",
    "vitest": "^2.1.0"
  }
}
```

- [ ] **Step 2: Create `frontend/nuxt.config.ts`**

```ts
export default defineNuxtConfig({
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  ssr: true,
  // pathPrefix: false — component names are just filenames regardless of directory depth.
  // e.g. components/booking/steps/StepSpecialty.vue → <StepSpecialty />, not <BookingStepsStepSpecialty />
  components: {
    dirs: [{ path: '~/components', pathPrefix: false }],
  },
  runtimeConfig: {
    apiBase: 'http://localhost:8080',   // overridden by NUXT_API_BASE in docker
    public: {
      clinicName: 'BeautyMed',
      clinicPhone: '+7 (4212) XX-XX-XX',
      clinicAddress: 'г. Хабаровск, ул. Примерная, 1',
      clinicHours: 'Пн–Сб: 9:00–20:00, Вс: выходной',
    },
  },
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    },
  },
})
```

- [ ] **Step 3: Create `frontend/tailwind.config.ts`**

```ts
import type { Config } from 'tailwindcss'

export default {
  content: ['./components/**/*.vue', './pages/**/*.vue', './app.vue'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#007C81', light: '#E6F5F5' },
        slate:   { DEFAULT: '#3C4F5B' },
        muted:   '#6B8290',
        border:  '#E6E6E6',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
} satisfies Config
```

- [ ] **Step 4: Create `frontend/app.vue`**

```vue
<template>
  <div>
    <NuxtPage />
  </div>
</template>
```

- [ ] **Step 5: Install dependencies**

Run from `frontend/` directory:
```bash
npm install
```

Expected: `node_modules/` created, no errors.

- [ ] **Step 6: Copy clinic photos to `frontend/public/`**

```bash
cp clinic_3.png frontend/public/clinic_3.png
cp clinic_2.png frontend/public/clinic_2.png
```

- [ ] **Step 7: Verify dev server starts**

Run from `frontend/`:
```bash
npm run dev
```

Expected: `Nuxt 3 ... ready on http://localhost:3000` — no TS errors.

- [ ] **Step 8: Commit**

```bash
git add frontend/package.json frontend/nuxt.config.ts frontend/tailwind.config.ts frontend/app.vue frontend/public/
git commit -m "feat(frontend): scaffold Nuxt 3 project"
```

---

## Task 2: TypeScript Types

**Files:**
- Create: `frontend/types/index.ts`

- [ ] **Step 1: Create `frontend/types/index.ts`**

```ts
export interface Specialty {
  id: string
  name: string
  slot_duration_min: number
}

export interface Doctor {
  id: string
  full_name: string
  specialty_id: string
  specialty_name: string
  bio: string
  photo_url: string
  experience_years: number
  is_active: boolean
}

export interface Service {
  id: string
  name: string
  description: string
  price: number
  duration_min: number
  specialty_id: string
  is_active: boolean
}

export interface TimeSlot {
  starts_at: string   // ISO 8601 datetime
  ends_at: string
}

export interface Review {
  author: string
  text: string
  rating: number
  date: string
}

export interface SpecialtyMeta {
  icon: string
  description: string
  color: string
}
```

- [ ] **Step 2: Run typecheck to verify no errors**

```bash
npm run typecheck
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/types/
git commit -m "feat(frontend): add TypeScript types"
```

---

## Task 3: Booking Pinia Store (TDD)

**Files:**
- Create: `frontend/stores/booking.ts`
- Create: `frontend/tests/booking.store.test.ts`

- [ ] **Step 1: Create `frontend/tests/booking.store.test.ts`**

```ts
import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach } from 'vitest'
import { useBookingStore } from '../stores/booking'

describe('booking store', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('opens at step 1 by default', () => {
    const s = useBookingStore()
    s.openModal()
    expect(s.open).toBe(true)
    expect(s.step).toBe(1)
  })

  it('opens at step 2 when specialtyId is pre-selected', () => {
    const s = useBookingStore()
    s.openModal('spec-123')
    expect(s.step).toBe(2)
    expect(s.specialtyId).toBe('spec-123')
  })

  it('step 1 canProceed requires both specialtyId and serviceId', () => {
    const s = useBookingStore()
    s.openModal()
    s.specialtyId = 'spec-1'
    expect(s.canProceed).toBe(false)
    s.serviceId = 'svc-1'
    expect(s.canProceed).toBe(true)
  })

  it('step 2 canProceed requires doctorId', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 2
    expect(s.canProceed).toBe(false)
    s.doctorId = 'doc-1'
    expect(s.canProceed).toBe(true)
  })

  it('step 3 canProceed requires date', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 3
    expect(s.canProceed).toBe(false)
    s.date = '2026-06-10'
    expect(s.canProceed).toBe(true)
  })

  it('step 4 canProceed requires timeSlot', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 4
    expect(s.canProceed).toBe(false)
    s.timeSlot = '2026-06-10T10:00:00Z'
    expect(s.canProceed).toBe(true)
  })

  it('nextStep increments step', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 2
    s.nextStep()
    expect(s.step).toBe(3)
  })

  it('nextStep stops at 5', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 5
    s.nextStep()
    expect(s.step).toBe(5)
  })

  it('prevStep decrements step', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 3
    s.prevStep()
    expect(s.step).toBe(2)
  })

  it('prevStep stops at 1', () => {
    const s = useBookingStore()
    s.openModal()
    s.step = 1
    s.prevStep()
    expect(s.step).toBe(1)
  })

  it('closeModal resets all state', () => {
    const s = useBookingStore()
    s.openModal('spec-1')
    s.doctorId = 'doc-1'
    s.timeSlot = '2026-06-10T10:00:00Z'
    s.closeModal()
    expect(s.open).toBe(false)
    expect(s.specialtyId).toBeNull()
    expect(s.doctorId).toBeNull()
    expect(s.timeSlot).toBeNull()
    expect(s.step).toBe(1)
  })

  it('otpCooldownSecs returns 0 with no cooldown', () => {
    const s = useBookingStore()
    expect(s.otpCooldownSecs).toBe(0)
  })

  it('startOtpCooldown sets ~60s cooldown', () => {
    const s = useBookingStore()
    s.startOtpCooldown()
    expect(s.otpCooldownSecs).toBeGreaterThan(58)
    expect(s.otpCooldownSecs).toBeLessThanOrEqual(60)
  })
})
```

- [ ] **Step 2: Run tests — verify they all fail**

```bash
npm run test
```

Expected: all tests FAIL with "Cannot find module '../stores/booking'".

- [ ] **Step 3: Create `frontend/stores/booking.ts`**

```ts
import { defineStore } from 'pinia'

export type BookingStep = 1 | 2 | 3 | 4 | 5

interface BookingState {
  open: boolean
  step: BookingStep
  specialtyId: string | null
  serviceId: string | null
  doctorId: string | null
  date: string | null
  timeSlot: string | null   // RFC3339 starts_at from slot
  name: string
  phone: string
  otpSent: boolean
  otpCooldownUntil: number | null
  token: string | null
  success: boolean
}

export const useBookingStore = defineStore('booking', {
  state: (): BookingState => ({
    open: false,
    step: 1,
    specialtyId: null,
    serviceId: null,
    doctorId: null,
    date: null,
    timeSlot: null,
    name: '',
    phone: '',
    otpSent: false,
    otpCooldownUntil: null,
    token: null,
    success: false,
  }),
  getters: {
    canProceed: (state): boolean => {
      switch (state.step) {
        case 1: return !!state.specialtyId && !!state.serviceId
        case 2: return !!state.doctorId
        case 3: return !!state.date
        case 4: return !!state.timeSlot
        default: return false
      }
    },
    otpCooldownSecs: (state): number => {
      if (!state.otpCooldownUntil) return 0
      return Math.max(0, Math.ceil((state.otpCooldownUntil - Date.now()) / 1000))
    },
  },
  actions: {
    openModal(specialtyId?: string) {
      this.open = true
      if (specialtyId) {
        this.specialtyId = specialtyId
        this.step = 2
      } else {
        this.step = 1
      }
    },
    closeModal() {
      this.$reset()
    },
    nextStep() {
      if (this.step < 5) this.step = (this.step + 1) as BookingStep
    },
    prevStep() {
      if (this.step > 1) this.step = (this.step - 1) as BookingStep
    },
    startOtpCooldown() {
      this.otpCooldownUntil = Date.now() + 60_000
    },
  },
})
```

- [ ] **Step 4: Run tests — all should pass**

```bash
npm run test
```

Expected: all 13 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add frontend/stores/booking.ts frontend/tests/booking.store.test.ts
git commit -m "feat(frontend): add booking Pinia store with unit tests"
```

---

## Task 4: useApi Composable + Static Data Files

**Files:**
- Create: `frontend/composables/useApi.ts`
- Create: `frontend/data/specialtyMeta.ts`
- Create: `frontend/data/reviews.ts`
- Create: `frontend/data/advantages.ts`

- [ ] **Step 1: Create `frontend/composables/useApi.ts`**

```ts
// SSR: uses runtimeConfig.apiBase (internal Docker URL)
// Client: uses '' (relative URL, nginx proxies /api/ to backend)
export const useApi = () => {
  const config = useRuntimeConfig()
  const base = import.meta.server ? config.apiBase : ''

  const get = <T>(path: string) => $fetch<T>(`${base}${path}`)

  const post = <T>(path: string, body: unknown, token?: string) =>
    $fetch<T>(`${base}${path}`, {
      method: 'POST',
      body,
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    })

  const patch = <T>(path: string, body: unknown, token: string) =>
    $fetch<T>(`${base}${path}`, {
      method: 'PATCH',
      body,
      headers: { Authorization: `Bearer ${token}` },
    })

  return { get, post, patch }
}
```

- [ ] **Step 2: Create `frontend/data/specialtyMeta.ts`**

```ts
import type { SpecialtyMeta } from '~/types'

export const SPECIALTY_META: Record<string, SpecialtyMeta> = {
  'Косметология': {
    icon: '✨',
    description: 'Аппаратная косметология, инъекции, пилинги и комплексный уход за кожей.',
    color: 'bg-pink-50',
  },
  'Дерматология': {
    icon: '🔬',
    description: 'Диагностика и лечение акне, дерматитов, розацеа. Удаление новообразований.',
    color: 'bg-blue-50',
  },
  'Трихология': {
    icon: '💆',
    description: 'Лечение выпадения волос, заболеваний кожи головы, PRP-терапия.',
    color: 'bg-green-50',
  },
  'Эстетическая медицина': {
    icon: '💎',
    description: 'Нитевой лифтинг, объёмное моделирование, коррекция возрастных изменений.',
    color: 'bg-purple-50',
  },
}

export const getSpecialtyMeta = (name: string): SpecialtyMeta =>
  SPECIALTY_META[name] ?? { icon: '🏥', description: '', color: 'bg-gray-50' }
```

- [ ] **Step 3: Create `frontend/data/reviews.ts`**

```ts
import type { Review } from '~/types'

export const REVIEWS: Review[] = [
  {
    author: 'Анна К.',
    text: 'Отличная клиника! Посещаю косметолога уже год. Кожа стала намного лучше, результат виден после первой процедуры.',
    rating: 5,
    date: '2026-03-15',
  },
  {
    author: 'Михаил Т.',
    text: 'Обратился к дерматологу с давней проблемой. Врач подобрал лечение — за месяц всё прошло. Спасибо!',
    rating: 5,
    date: '2026-02-28',
  },
  {
    author: 'Елена С.',
    text: 'Прохожу курс трихологии. Видны результаты уже после 3 сеансов. Вежливый персонал, уютная обстановка.',
    rating: 5,
    date: '2026-04-10',
  },
]
```

- [ ] **Step 4: Create `frontend/data/advantages.ts`**

```ts
export const ADVANTAGES = [
  { icon: '👩‍⚕️', title: 'Опытные врачи', text: '8 специалистов с опытом от 5 до 15 лет' },
  { icon: '🏥', title: 'Современное оборудование', text: 'Аппаратура последнего поколения' },
  { icon: '📱', title: 'Онлайн-запись', text: 'Запишитесь к врачу за 2 минуты, без звонков' },
  { icon: '💬', title: 'Индивидуальный подход', text: 'Программа лечения под каждого пациента' },
]
```

- [ ] **Step 5: Commit**

```bash
git add frontend/composables/ frontend/data/
git commit -m "feat(frontend): add useApi composable and static data files"
```

---

## Task 5: HeroSection

**Files:**
- Create: `frontend/components/sections/HeroSection.vue`

- [ ] **Step 1: Create `frontend/components/sections/HeroSection.vue`**

```vue
<script setup lang="ts">
const booking = useBookingStore()
</script>

<template>
  <section class="bg-white border-b border-border">
    <div class="max-w-5xl mx-auto">
      <!-- Main hero row -->
      <div class="flex items-stretch min-h-[280px]">
        <!-- Text block -->
        <div class="flex-1 px-8 py-10 flex flex-col justify-center">
          <p class="text-xs font-semibold tracking-[2.5px] text-primary uppercase mb-3">
            Клиника красоты и здоровья · Хабаровск
          </p>
          <h1 class="text-3xl font-extrabold text-slate leading-tight mb-3">
            Красота и здоровье<br>в надёжных руках
          </h1>
          <p class="text-muted text-sm leading-relaxed mb-6">
            Косметология, дерматология,<br>трихология и эстетическая медицина
          </p>
          <div class="flex gap-3">
            <button
              class="bg-primary text-white px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors"
              @click="booking.openModal()"
            >
              Записаться онлайн
            </button>
            <a
              href="#services"
              class="border border-primary text-primary px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/5 transition-colors"
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
      <div class="flex border-t border-border bg-gray-50">
        <div class="flex-1 py-3 text-center border-r border-border">
          <span class="block text-xl font-extrabold text-primary">10+</span>
          <span class="text-xs text-slate">лет работы</span>
        </div>
        <div class="flex-1 py-3 text-center border-r border-border">
          <span class="block text-xl font-extrabold text-primary">8</span>
          <span class="text-xs text-slate">врачей</span>
        </div>
        <div class="flex-1 py-3 text-center">
          <span class="block text-xl font-extrabold text-primary">3000+</span>
          <span class="text-xs text-slate">пациентов</span>
        </div>
      </div>
    </div>
  </section>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/sections/HeroSection.vue
git commit -m "feat(frontend): add HeroSection"
```

---

## Task 6: AdvantagesSection

**Files:**
- Create: `frontend/components/sections/AdvantagesSection.vue`

- [ ] **Step 1: Create `frontend/components/sections/AdvantagesSection.vue`**

```vue
<script setup lang="ts">
import { ADVANTAGES } from '~/data/advantages'
</script>

<template>
  <section class="py-12 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Почему выбирают нас</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div
          v-for="adv in ADVANTAGES"
          :key="adv.title"
          class="text-center"
        >
          <div class="text-3xl mb-3">{{ adv.icon }}</div>
          <div class="text-sm font-semibold text-slate mb-1">{{ adv.title }}</div>
          <div class="text-xs text-muted leading-relaxed">{{ adv.text }}</div>
        </div>
      </div>
    </div>
  </section>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/sections/AdvantagesSection.vue
git commit -m "feat(frontend): add AdvantagesSection"
```

---

## Task 7: ServicesSection

**Files:**
- Create: `frontend/components/sections/ServicesSection.vue`

- [ ] **Step 1: Create `frontend/components/sections/ServicesSection.vue`**

```vue
<script setup lang="ts">
import type { Specialty, Service } from '~/types'
import { getSpecialtyMeta } from '~/data/specialtyMeta'

const props = defineProps<{
  specialties: Specialty[]
  services: Service[]
}>()

const booking = useBookingStore()
const activeSpecialty = ref<string | null>(props.specialties[0]?.id ?? null)

const activeServices = computed(() =>
  props.services.filter(s => s.specialty_id === activeSpecialty.value)
)

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(price)
</script>

<template>
  <section id="services" class="py-12 bg-gray-50">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Наши услуги</h2>

      <!-- Specialty tabs -->
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="sp in specialties"
          :key="sp.id"
          class="px-4 py-2 rounded-full text-sm font-medium transition-colors"
          :class="activeSpecialty === sp.id
            ? 'bg-primary text-white'
            : 'bg-white border border-border text-muted hover:border-primary hover:text-primary'"
          @click="activeSpecialty = sp.id"
        >
          {{ getSpecialtyMeta(sp.name).icon }} {{ sp.name }}
        </button>
      </div>

      <!-- Service list -->
      <div class="bg-white rounded-xl border border-border overflow-hidden">
        <div
          v-for="(svc, i) in activeServices"
          :key="svc.id"
          class="flex items-center justify-between px-5 py-4"
          :class="i < activeServices.length - 1 ? 'border-b border-border' : ''"
        >
          <div class="flex-1 pr-4">
            <div class="text-sm font-medium text-slate">{{ svc.name }}</div>
            <div class="text-xs text-muted mt-0.5">{{ svc.duration_min }} мин</div>
          </div>
          <div class="flex items-center gap-4">
            <span class="text-sm font-semibold text-primary">{{ formatPrice(svc.price) }}</span>
            <button
              class="text-xs font-semibold text-primary border border-primary px-3 py-1.5 rounded-lg hover:bg-primary/5 transition-colors whitespace-nowrap"
              @click="booking.openModal(activeSpecialty ?? undefined)"
            >
              Записаться
            </button>
          </div>
        </div>
        <div v-if="activeServices.length === 0" class="px-5 py-8 text-center text-muted text-sm">
          Услуги не найдены
        </div>
      </div>
    </div>
  </section>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/sections/ServicesSection.vue
git commit -m "feat(frontend): add ServicesSection with specialty tabs"
```

---

## Task 8: DoctorsSection

**Files:**
- Create: `frontend/components/sections/DoctorsSection.vue`

- [ ] **Step 1: Create `frontend/components/sections/DoctorsSection.vue`**

```vue
<script setup lang="ts">
import type { Doctor } from '~/types'

defineProps<{ doctors: Doctor[] }>()
const booking = useBookingStore()
</script>

<template>
  <section class="py-12 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Наши врачи</h2>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
        <div
          v-for="doc in doctors"
          :key="doc.id"
          class="bg-white border border-border rounded-xl overflow-hidden"
        >
          <!-- Photo with specialty badge -->
          <div class="relative h-40 bg-primary/10 flex items-center justify-center overflow-hidden">
            <img
              v-if="doc.photo_url"
              :src="doc.photo_url"
              :alt="doc.full_name"
              class="w-full h-full object-cover object-top"
            >
            <span v-else class="text-5xl">👩‍⚕️</span>
            <span class="absolute bottom-2 left-2 bg-primary text-white text-[9px] font-bold uppercase tracking-wide px-2 py-1 rounded-full">
              {{ doc.specialty_name }}
            </span>
          </div>
          <!-- Info -->
          <div class="p-3">
            <div class="text-xs font-semibold text-slate leading-snug mb-1">{{ doc.full_name }}</div>
            <div class="text-[11px] text-muted mb-3">Стаж {{ doc.experience_years }} лет</div>
            <button
              class="w-full bg-primary/10 text-primary text-[11px] font-semibold py-1.5 rounded-lg hover:bg-primary/20 transition-colors"
              @click="booking.openModal(doc.specialty_id)"
            >
              Записаться
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/sections/DoctorsSection.vue
git commit -m "feat(frontend): add DoctorsSection grid"
```

---

## Task 9: ReviewsSection

**Files:**
- Create: `frontend/components/sections/ReviewsSection.vue`

- [ ] **Step 1: Create `frontend/components/sections/ReviewsSection.vue`**

```vue
<script setup lang="ts">
import { REVIEWS } from '~/data/reviews'
</script>

<template>
  <section class="py-12 bg-gray-50">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Отзывы пациентов</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <div
          v-for="review in REVIEWS"
          :key="review.author"
          class="bg-white rounded-xl border border-border p-5"
        >
          <div class="flex gap-0.5 mb-3">
            <span v-for="i in 5" :key="i" class="text-yellow-400 text-sm">
              {{ i <= review.rating ? '★' : '☆' }}
            </span>
          </div>
          <p class="text-sm text-slate leading-relaxed mb-4">{{ review.text }}</p>
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-muted">{{ review.author }}</span>
            <span class="text-xs text-muted">{{ review.date }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/sections/ReviewsSection.vue
git commit -m "feat(frontend): add ReviewsSection"
```

---

## Task 10: ContactsSection

**Files:**
- Create: `frontend/components/sections/ContactsSection.vue`

- [ ] **Step 1: Create `frontend/components/sections/ContactsSection.vue`**

```vue
<script setup lang="ts">
const config = useRuntimeConfig()
</script>

<template>
  <section class="py-12 bg-white">
    <div class="max-w-5xl mx-auto px-8">
      <h2 class="text-xl font-bold text-slate mb-8 text-center">Контакты</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- Info -->
        <div class="space-y-4">
          <div class="flex items-start gap-3">
            <span class="text-primary mt-0.5">📍</span>
            <div>
              <div class="text-sm font-semibold text-slate mb-0.5">Адрес</div>
              <div class="text-sm text-muted">{{ config.public.clinicAddress }}</div>
            </div>
          </div>
          <div class="flex items-start gap-3">
            <span class="text-primary mt-0.5">📞</span>
            <div>
              <div class="text-sm font-semibold text-slate mb-0.5">Телефон</div>
              <a :href="`tel:${config.public.clinicPhone}`" class="text-sm text-primary font-medium">
                {{ config.public.clinicPhone }}
              </a>
            </div>
          </div>
          <div class="flex items-start gap-3">
            <span class="text-primary mt-0.5">🕐</span>
            <div>
              <div class="text-sm font-semibold text-slate mb-0.5">Часы работы</div>
              <div class="text-sm text-muted">{{ config.public.clinicHours }}</div>
            </div>
          </div>
        </div>
        <!-- Map placeholder (replace iframe src with real Yandex Maps embed URL) -->
        <div class="rounded-xl overflow-hidden border border-border bg-gray-100 min-h-[200px] flex items-center justify-center">
          <span class="text-muted text-sm">Карта (добавить Yandex Maps iframe)</span>
        </div>
      </div>
    </div>
  </section>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/sections/ContactsSection.vue
git commit -m "feat(frontend): add ContactsSection"
```

---

## Task 11: BookingModal Shell

**Files:**
- Create: `frontend/components/booking/BookingModal.vue`

The modal wraps all 5 steps. It renders the correct step component, shows dot indicators, and handles Back/Next footer.

- [ ] **Step 1: Create `frontend/components/booking/BookingModal.vue`**

```vue
<script setup lang="ts">
const booking = useBookingStore()

// resolveComponent uses the filename (pathPrefix: false in nuxt.config.ts)
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
onMounted(() => {
  const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') booking.closeModal() }
  window.addEventListener('keydown', handler)
  onUnmounted(() => window.removeEventListener('keydown', handler))
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="booking.open"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/40 backdrop-blur-sm"
          @click="booking.closeModal()"
        />

        <!-- Modal box -->
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md max-h-[90vh] flex flex-col overflow-hidden">
          <!-- Header -->
          <div class="flex items-center justify-between px-6 pt-5 pb-0">
            <span class="text-sm font-bold text-slate">Онлайн-запись</span>
            <button class="text-muted hover:text-slate transition-colors text-lg" @click="booking.closeModal()">✕</button>
          </div>

          <!-- Body -->
          <div class="flex-1 overflow-y-auto px-6 py-4">
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
          <div v-if="booking.step < 5" class="px-6 py-4 border-t border-border flex justify-between items-center">
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
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/booking/BookingModal.vue
git commit -m "feat(frontend): add BookingModal shell with dot indicator"
```

---

## Task 12: StepSpecialty (Specialty + Service Selection)

**Files:**
- Create: `frontend/components/booking/steps/StepSpecialty.vue`

Step 1 shows specialty options; after selecting a specialty, a service list appears below. Both `specialtyId` and `serviceId` must be set for `canProceed` to be true.

Props passed via the page: `specialties` and `services` are injected via `provide`/`inject` from `pages/index.vue`.

- [ ] **Step 1: Create `frontend/components/booking/steps/StepSpecialty.vue`**

```vue
<script setup lang="ts">
import type { Specialty, Service } from '~/types'
import { getSpecialtyMeta } from '~/data/specialtyMeta'

const booking = useBookingStore()
const specialties = inject<Ref<Specialty[]>>('specialties', ref([]))
const services = inject<Ref<Service[]>>('services', ref([]))

const servicesForSelected = computed(() =>
  services.value.filter(s => s.specialty_id === booking.specialtyId)
)

const selectSpecialty = (id: string) => {
  booking.specialtyId = id
  booking.serviceId = null
}

const formatPrice = (price: number) =>
  new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', maximumFractionDigits: 0 }).format(price)
</script>

<template>
  <div>
    <!-- Specialty options -->
    <div class="text-xs font-semibold text-slate mb-2">Выберите специальность</div>
    <div class="flex flex-col gap-2 mb-4">
      <button
        v-for="sp in specialties"
        :key="sp.id"
        class="flex items-center justify-between px-4 py-3 rounded-xl border text-sm transition-colors"
        :class="booking.specialtyId === sp.id
          ? 'border-primary bg-primary/5 text-primary font-semibold'
          : 'border-border text-slate hover:border-primary'"
        @click="selectSpecialty(sp.id)"
      >
        <span>{{ getSpecialtyMeta(sp.name).icon }} {{ sp.name }}</span>
        <span>{{ booking.specialtyId === sp.id ? '✓' : '›' }}</span>
      </button>
    </div>

    <!-- Service options (appear after specialty selected) -->
    <Transition name="slide-down">
      <div v-if="booking.specialtyId && servicesForSelected.length">
        <div class="text-xs font-semibold text-slate mb-2">Выберите услугу</div>
        <div class="flex flex-col gap-1.5">
          <button
            v-for="svc in servicesForSelected"
            :key="svc.id"
            class="flex items-center justify-between px-4 py-2.5 rounded-lg border text-sm transition-colors"
            :class="booking.serviceId === svc.id
              ? 'border-primary bg-primary/5 text-primary font-semibold'
              : 'border-border text-slate hover:border-primary'"
            @click="booking.serviceId = svc.id"
          >
            <span class="text-left">{{ svc.name }}</span>
            <span class="ml-3 whitespace-nowrap text-xs">{{ formatPrice(svc.price) }}</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.slide-down-enter-active { transition: all 0.2s ease; }
.slide-down-enter-from { opacity: 0; transform: translateY(-6px); }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/booking/steps/StepSpecialty.vue
git commit -m "feat(frontend): add StepSpecialty with inline service selection"
```

---

## Task 13: StepDoctor

**Files:**
- Create: `frontend/components/booking/steps/StepDoctor.vue`

Filters doctors by selected specialty. User picks a doctor to proceed.

- [ ] **Step 1: Create `frontend/components/booking/steps/StepDoctor.vue`**

```vue
<script setup lang="ts">
import type { Doctor } from '~/types'

const booking = useBookingStore()
const doctors = inject<Ref<Doctor[]>>('doctors', ref([]))

const filteredDoctors = computed(() =>
  doctors.value.filter(d => d.specialty_id === booking.specialtyId)
)
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-2">Выберите врача</div>
    <div class="flex flex-col gap-2">
      <button
        v-for="doc in filteredDoctors"
        :key="doc.id"
        class="flex items-center gap-3 px-4 py-3 rounded-xl border text-left transition-colors"
        :class="booking.doctorId === doc.id
          ? 'border-primary bg-primary/5'
          : 'border-border hover:border-primary'"
        @click="booking.doctorId = doc.id"
      >
        <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0 overflow-hidden">
          <img v-if="doc.photo_url" :src="doc.photo_url" :alt="doc.full_name" class="w-full h-full object-cover" />
          <span v-else class="text-lg">👩‍⚕️</span>
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium text-slate truncate">{{ doc.full_name }}</div>
          <div class="text-xs text-muted">Стаж {{ doc.experience_years }} лет</div>
        </div>
        <span v-if="booking.doctorId === doc.id" class="text-primary text-sm">✓</span>
      </button>
    </div>
    <div v-if="filteredDoctors.length === 0" class="text-sm text-muted text-center py-6">
      Нет доступных врачей
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/booking/steps/StepDoctor.vue
git commit -m "feat(frontend): add StepDoctor"
```

---

## Task 14: StepDate

**Files:**
- Create: `frontend/components/booking/steps/StepDate.vue`

A simple month calendar. User clicks a date; past dates and Sundays are disabled.

- [ ] **Step 1: Create `frontend/components/booking/steps/StepDate.vue`**

```vue
<script setup lang="ts">
const booking = useBookingStore()

const today = new Date()
today.setHours(0, 0, 0, 0)

const viewDate = ref(new Date(today.getFullYear(), today.getMonth(), 1))

const monthLabel = computed(() =>
  viewDate.value.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })
)

const prevMonth = () => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() - 1)
  if (d >= new Date(today.getFullYear(), today.getMonth(), 1)) viewDate.value = d
}

const nextMonth = () => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() + 1)
  viewDate.value = d
}

const calendarDays = computed(() => {
  const year = viewDate.value.getFullYear()
  const month = viewDate.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)

  // Start grid on Monday (ISO week)
  const startOffset = (firstDay.getDay() + 6) % 7
  const days: { date: Date | null; disabled: boolean }[] = []

  for (let i = 0; i < startOffset; i++) days.push({ date: null, disabled: true })

  for (let d = 1; d <= lastDay.getDate(); d++) {
    const date = new Date(year, month, d)
    const isSunday = date.getDay() === 0
    const isPast = date < today
    days.push({ date, disabled: isSunday || isPast })
  }

  return days
})

const toISODate = (d: Date) => d.toISOString().slice(0, 10)

const selectDate = (d: Date) => {
  booking.date = toISODate(d)
  booking.timeSlot = null  // reset slot when date changes
}

const isSelected = (d: Date) => booking.date === toISODate(d)
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-3">Выберите дату</div>

    <!-- Month navigation -->
    <div class="flex items-center justify-between mb-3">
      <button class="text-muted hover:text-primary p-1 transition-colors" @click="prevMonth">‹</button>
      <span class="text-sm font-semibold text-slate capitalize">{{ monthLabel }}</span>
      <button class="text-muted hover:text-primary p-1 transition-colors" @click="nextMonth">›</button>
    </div>

    <!-- Day-of-week headers -->
    <div class="grid grid-cols-7 text-center mb-1">
      <div v-for="d in ['Пн','Вт','Ср','Чт','Пт','Сб','Вс']" :key="d" class="text-[10px] text-muted py-1">
        {{ d }}
      </div>
    </div>

    <!-- Calendar grid -->
    <div class="grid grid-cols-7 gap-0.5">
      <div v-for="(cell, i) in calendarDays" :key="i" class="aspect-square flex items-center justify-center">
        <button
          v-if="cell.date"
          class="w-8 h-8 rounded-full text-xs font-medium transition-colors"
          :class="{
            'bg-primary text-white': isSelected(cell.date),
            'text-muted cursor-not-allowed': cell.disabled,
            'text-slate hover:bg-primary/10': !cell.disabled && !isSelected(cell.date),
          }"
          :disabled="cell.disabled"
          @click="!cell.disabled && selectDate(cell.date!)"
        >
          {{ cell.date.getDate() }}
        </button>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/booking/steps/StepDate.vue
git commit -m "feat(frontend): add StepDate calendar"
```

---

## Task 15: StepTime

**Files:**
- Create: `frontend/components/booking/steps/StepTime.vue`

Fetches available slots from `GET /api/doctors/:id/slots?service_id=&date=` when date changes. Slots are `{starts_at, ends_at}` ISO strings.

- [ ] **Step 1: Create `frontend/components/booking/steps/StepTime.vue`**

```vue
<script setup lang="ts">
import type { TimeSlot } from '~/types'

const booking = useBookingStore()
const { get } = useApi()

const slots = ref<TimeSlot[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const fetchSlots = async () => {
  if (!booking.doctorId || !booking.serviceId || !booking.date) return
  loading.value = true
  error.value = null
  try {
    slots.value = await get<TimeSlot[]>(
      `/api/doctors/${booking.doctorId}/slots?service_id=${booking.serviceId}&date=${booking.date}`
    )
  } catch {
    error.value = 'Не удалось загрузить слоты'
    slots.value = []
  } finally {
    loading.value = false
  }
}

watch(() => booking.date, fetchSlots, { immediate: true })

const formatTime = (iso: string) =>
  new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' })
</script>

<template>
  <div>
    <div class="text-xs font-semibold text-slate mb-3">
      Выберите время — {{ booking.date }}
    </div>

    <div v-if="loading" class="text-sm text-muted text-center py-6">Загружаем расписание...</div>

    <div v-else-if="error" class="text-center py-6">
      <div class="text-sm text-red-500 mb-3">{{ error }}</div>
      <button class="text-xs text-primary underline" @click="fetchSlots">Попробовать снова</button>
    </div>

    <div v-else-if="slots.length === 0" class="text-sm text-muted text-center py-6">
      Нет доступных слотов на эту дату
    </div>

    <div v-else class="grid grid-cols-3 gap-2">
      <button
        v-for="slot in slots"
        :key="slot.starts_at"
        class="py-2 rounded-lg border text-sm font-medium transition-colors"
        :class="booking.timeSlot === slot.starts_at
          ? 'bg-primary text-white border-primary'
          : 'border-border text-slate hover:border-primary'"
        @click="booking.timeSlot = slot.starts_at"
      >
        {{ formatTime(slot.starts_at) }}
      </button>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/booking/steps/StepTime.vue
git commit -m "feat(frontend): add StepTime with slot fetching"
```

---

## Task 16: StepConfirm (OTP Flow)

**Files:**
- Create: `frontend/components/booking/steps/StepConfirm.vue`

This step has two sub-phases:
1. Name + phone → "Получить код" → `POST /api/auth/otp`
2. OTP input → "Записаться" → `POST /api/auth/verify` → store token → `POST /api/appointments`

The component also shows a success state after booking.

- [ ] **Step 1: Create `frontend/components/booking/steps/StepConfirm.vue`**

```vue
<script setup lang="ts">
const booking = useBookingStore()
const { post } = useApi()

const submittingOtp = ref(false)
const submittingBooking = ref(false)
const otpCode = ref('')
const otpError = ref('')
const bookingError = ref('')
const cooldownInterval = ref<ReturnType<typeof setInterval> | null>(null)
const cooldownDisplay = ref(0)

// Tick cooldown display every second
onMounted(() => {
  cooldownInterval.value = setInterval(() => {
    cooldownDisplay.value = booking.otpCooldownSecs
  }, 1000)
})
onUnmounted(() => {
  if (cooldownInterval.value) clearInterval(cooldownInterval.value)
})

const sendOtp = async () => {
  if (!booking.phone) return
  submittingOtp.value = true
  otpError.value = ''
  try {
    await post('/api/auth/otp', { phone: booking.phone })
    booking.otpSent = true
    booking.startOtpCooldown()
    cooldownDisplay.value = 60
  } catch {
    otpError.value = 'Не удалось отправить код. Проверьте номер.'
  } finally {
    submittingOtp.value = false
  }
}

const confirmBooking = async () => {
  if (!otpCode.value) return
  submittingBooking.value = true
  bookingError.value = ''
  try {
    // 1. Verify OTP
    const { access_token } = await post<{ access_token: string }>(
      '/api/auth/verify',
      { phone: booking.phone, code: otpCode.value }
    )
    booking.token = access_token

    // 2. Create appointment
    await post(
      '/api/appointments',
      {
        doctor_id: booking.doctorId,
        service_id: booking.serviceId,
        starts_at: booking.timeSlot,
        promo_code: '',
      },
      access_token
    )

    booking.success = true
  } catch (err: unknown) {
    const msg = (err as { data?: { error?: string } })?.data?.error
    if (msg?.includes('invalid') || msg?.includes('OTP')) {
      bookingError.value = 'Неверный или истёкший код. Попробуйте ещё раз.'
    } else if (msg?.includes('taken') || msg?.includes('conflict')) {
      bookingError.value = 'Это время уже занято. Выберите другой слот.'
    } else {
      bookingError.value = 'Ошибка записи. Позвоните нам: ' + useRuntimeConfig().public.clinicPhone
    }
  } finally {
    submittingBooking.value = false
  }
}

const formatTime = (iso: string | null) =>
  iso ? new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' }) : ''
</script>

<template>
  <div>
    <!-- Success state -->
    <div v-if="booking.success" class="text-center py-4">
      <div class="text-4xl mb-3">✅</div>
      <div class="text-base font-bold text-slate mb-1">Вы записаны!</div>
      <div class="text-sm text-muted mb-1">{{ booking.date }}, {{ formatTime(booking.timeSlot) }}</div>
      <div class="text-sm text-muted mb-6">Ждём вас в клинике BeautyMed</div>
      <button
        class="bg-primary text-white px-6 py-2.5 rounded-lg text-sm font-semibold"
        @click="booking.closeModal()"
      >
        Закрыть
      </button>
    </div>

    <!-- Form state -->
    <div v-else class="space-y-4">
      <!-- Summary -->
      <div class="bg-gray-50 rounded-xl p-4 text-sm space-y-1 text-slate">
        <div>📅 {{ booking.date }}, {{ formatTime(booking.timeSlot) }}</div>
      </div>

      <!-- Name -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">Ваше имя</label>
        <input
          v-model="booking.name"
          type="text"
          placeholder="Иван Иванов"
          class="w-full border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <!-- Phone + OTP send -->
      <div>
        <label class="text-xs font-semibold text-slate block mb-1">Номер телефона</label>
        <div class="flex gap-2">
          <input
            v-model="booking.phone"
            type="tel"
            placeholder="+79001234567"
            :disabled="booking.otpSent"
            class="flex-1 border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary disabled:bg-gray-50"
          >
          <button
            class="px-4 py-2.5 rounded-lg text-sm font-semibold whitespace-nowrap transition-colors"
            :class="cooldownDisplay > 0 || submittingOtp
              ? 'bg-gray-100 text-muted cursor-not-allowed'
              : 'bg-primary text-white hover:bg-primary/90'"
            :disabled="cooldownDisplay > 0 || submittingOtp || !booking.phone"
            @click="sendOtp"
          >
            {{ cooldownDisplay > 0 ? `${cooldownDisplay}с` : booking.otpSent ? 'Повтор' : 'Код' }}
          </button>
        </div>
        <div v-if="otpError" class="text-xs text-red-500 mt-1">{{ otpError }}</div>
      </div>

      <!-- OTP code input (appears after sending) -->
      <Transition name="slide-down">
        <div v-if="booking.otpSent">
          <label class="text-xs font-semibold text-slate block mb-1">Код из SMS</label>
          <input
            v-model="otpCode"
            type="text"
            inputmode="numeric"
            maxlength="6"
            placeholder="123456"
            class="w-full border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary tracking-widest text-center"
          >
          <div v-if="bookingError" class="text-xs text-red-500 mt-1">{{ bookingError }}</div>
        </div>
      </Transition>

      <!-- Submit -->
      <button
        v-if="booking.otpSent"
        class="w-full bg-primary text-white py-3 rounded-xl text-sm font-semibold transition-opacity"
        :class="(!otpCode || submittingBooking) ? 'opacity-50 cursor-not-allowed' : ''"
        :disabled="!otpCode || submittingBooking"
        @click="confirmBooking"
      >
        {{ submittingBooking ? 'Записываем...' : 'Подтвердить и записаться' }}
      </button>

      <!-- Back button (step 5 has no footer) -->
      <div class="flex justify-start pt-1">
        <button class="text-sm text-muted hover:text-slate transition-colors" @click="booking.prevStep()">
          ← Назад
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-down-enter-active { transition: all 0.2s ease; }
.slide-down-enter-from { opacity: 0; transform: translateY(-6px); }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/booking/steps/StepConfirm.vue
git commit -m "feat(frontend): add StepConfirm with OTP flow and booking submission"
```

---

## Task 17: pages/index.vue — Assembly + SEO

**Files:**
- Create: `frontend/pages/index.vue`

This is the main SSR page. It fetches all data in parallel and provides it to child components via `provide`. It also assembles all sections and mounts `BookingModal`.

- [ ] **Step 1: Create `frontend/pages/index.vue`**

```vue
<script setup lang="ts">
import type { Specialty, Doctor, Service } from '~/types'

// useApi handles base URL: server → config.apiBase (http://api:8080), client → '' (through nginx)
const { get } = useApi()

// SSR: parallel data fetch — fetcher runs server-side, result hydrated to client
const [
  { data: specialties },
  { data: doctors },
  { data: services },
] = await Promise.all([
  useAsyncData('specialties', () => get<Specialty[]>('/api/specialties')),
  useAsyncData('doctors',     () => get<Doctor[]>('/api/doctors')),
  useAsyncData('services',    () => get<Service[]>('/api/services')),
])

// Fallback to empty arrays on SSR error (site still renders with static content)
const safeSpecialties = computed(() => specialties.value ?? [])
const safeDoctors     = computed(() => doctors.value ?? [])
const safeServices    = computed(() => services.value ?? [])

// Provide to booking step components
provide('specialties', safeSpecialties)
provide('doctors', safeDoctors)
provide('services', safeServices)

const config = useRuntimeConfig()

// SEO
useHead({
  title: 'BeautyMed — Клиника красоты и здоровья в Хабаровске',
  meta: [
    {
      name: 'description',
      content: 'Косметология, дерматология, трихология и эстетическая медицина. Онлайн-запись к врачу за 2 минуты.',
    },
    { property: 'og:title', content: 'BeautyMed — Клиника красоты и здоровья' },
    { property: 'og:description', content: 'Профессиональная косметология и эстетическая медицина в Хабаровске.' },
    { property: 'og:image', content: '/clinic_3.png' },
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'MedicalClinic',
        name: 'BeautyMed',
        address: {
          '@type': 'PostalAddress',
          addressLocality: 'Хабаровск',
          streetAddress: config.public.clinicAddress,
          addressCountry: 'RU',
        },
        telephone: config.public.clinicPhone,
        openingHours: 'Mo-Sa 09:00-20:00',
        url: 'https://beautymed.ru',
      }),
    },
  ],
})
</script>

<template>
  <main>
    <HeroSection />
    <AdvantagesSection />
    <ServicesSection :specialties="safeSpecialties" :services="safeServices" />
    <DoctorsSection :doctors="safeDoctors" />
    <ReviewsSection />
    <ContactsSection />
    <BookingModal />
  </main>
</template>
```

- [ ] **Step 2: Start dev server and verify the page renders**

Run from `frontend/`:
```bash
npm run dev
```

Open http://localhost:3000. Expected:
- Hero section with photo renders
- Advantages section shows 4 cards
- If backend is running locally: services and doctors sections populate
- If backend is not running: sections render empty (graceful fallback)
- No console errors
- BookingModal is invisible until triggered

- [ ] **Step 3: Verify booking modal opens**

Click any "Записаться" button. Expected: modal opens, dot indicator at step 1, specialties list visible.

- [ ] **Step 4: Commit**

```bash
git add frontend/pages/index.vue
git commit -m "feat(frontend): assemble index page with SSR data fetch and SEO"
```

---

## Task 18: Playwright E2E Test

**Files:**
- Create: `frontend/playwright.config.ts`
- Create: `frontend/tests/e2e/booking.spec.ts`

- [ ] **Step 1: Create `frontend/playwright.config.ts`**

```ts
import { defineConfig, devices } from '@playwright/test'

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
```

- [ ] **Step 2: Create `frontend/tests/e2e/booking.spec.ts`**

```ts
import { test, expect } from '@playwright/test'

test('homepage loads with hero section', async ({ page }) => {
  await page.goto('/')
  await expect(page.locator('h1')).toContainText('Красота и здоровье')
  await expect(page.locator('text=Записаться онлайн').first()).toBeVisible()
})

test('booking modal opens when CTA clicked', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  await expect(page.locator('text=Онлайн-запись')).toBeVisible()
  await expect(page.locator('text=Шаг 1')).toBeVisible()
})

test('booking modal closes on backdrop click', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  await expect(page.locator('text=Онлайн-запись')).toBeVisible()
  // Click backdrop (outside modal box)
  await page.mouse.click(10, 10)
  await expect(page.locator('text=Онлайн-запись')).not.toBeVisible()
})

test('booking modal closes on Escape', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  await page.keyboard.press('Escape')
  await expect(page.locator('text=Онлайн-запись')).not.toBeVisible()
})

test('step 1 Next button disabled until specialty and service selected', async ({ page }) => {
  await page.goto('/')
  await page.locator('button', { hasText: 'Записаться онлайн' }).first().click()
  const nextBtn = page.locator('button', { hasText: 'Далее →' })
  await expect(nextBtn).toBeDisabled()
})
```

- [ ] **Step 3: Install Playwright browsers**

```bash
npx playwright install chromium
```

- [ ] **Step 4: Run E2E tests**

Make sure dev server is running on port 3000, then:
```bash
npm run test:e2e
```

Expected: all 5 tests pass.

- [ ] **Step 5: Commit**

```bash
git add frontend/playwright.config.ts frontend/tests/e2e/
git commit -m "test(frontend): add Playwright E2E tests for homepage and booking modal"
```

---

## Task 19: Dockerfile + docker-compose + nginx

**Files:**
- Create: `frontend/Dockerfile`
- Modify: `docker-compose.yml`
- Modify: `nginx/nginx.conf`

- [ ] **Step 1: Create `frontend/Dockerfile`**

```dockerfile
FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:22-alpine
WORKDIR /app
COPY --from=builder /app/.output ./
EXPOSE 3000
ENV NITRO_PORT=3000
CMD ["node", "server/index.mjs"]
```

- [ ] **Step 2: Update `docker-compose.yml` — add frontend service**

Replace the `docker-compose.yml` with:

```yaml
version: "3.9"

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: beautymed
      POSTGRES_PASSWORD: beautymed
      POSTGRES_DB: beautymed
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U beautymed"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  api:
    build: ./backend
    env_file: .env
    environment:
      DATABASE_URL: postgres://beautymed:beautymed@postgres:5432/beautymed?sslmode=disable
      REDIS_URL: redis://redis:6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  frontend:
    build: ./frontend
    environment:
      NUXT_API_BASE: http://api:8080
    depends_on:
      - api

  nginx:
    image: nginx:alpine
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "80:80"
    depends_on:
      - api
      - frontend

volumes:
  pgdata:
```

Note: removed `ports: - "8080:8080"` from `api` (no longer exposed directly; nginx proxies) and removed `ports: - "6379:6379"` from `redis` (internal only in production).

- [ ] **Step 3: Update `nginx/nginx.conf`**

```nginx
events { worker_connections 1024; }

http {
  upstream api {
    server api:8080;
  }

  upstream frontend {
    server frontend:3000;
  }

  server {
    listen 80;

    # Backend health check
    location /health {
      proxy_pass http://api;
      proxy_set_header Host $host;
    }

    # Backend API
    location /api/ {
      proxy_pass http://api;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
    }

    # Frontend (Nuxt SSR)
    location / {
      proxy_pass http://frontend;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection 'upgrade';
      proxy_http_version 1.1;
    }
  }
}
```

- [ ] **Step 4: Verify docker-compose builds locally**

```bash
docker compose build frontend
```

Expected: build completes without errors. The Nuxt standalone build outputs to `.output/`.

- [ ] **Step 5: Commit**

```bash
git add frontend/Dockerfile docker-compose.yml nginx/nginx.conf
git commit -m "feat(infra): add frontend Docker service and update nginx routing"
```

---

## Final Verification

After all tasks are complete, run the full stack locally:

```bash
docker compose up --build
```

Open http://localhost:80. Verify:
- [ ] Homepage renders with hero photo
- [ ] Advantages, services, doctors sections load from backend
- [ ] Clicking "Записаться онлайн" opens booking modal
- [ ] Step indicator dots animate correctly
- [ ] Steps 1–4 navigate with Back/Next
- [ ] Step 5 (confirm) shows phone + OTP flow
- [ ] All Vitest unit tests pass: `npm run test`
- [ ] All Playwright E2E tests pass: `npm run test:e2e`
- [ ] `npm run typecheck` passes with no errors

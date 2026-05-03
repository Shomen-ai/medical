# BeautyMed Plan 3 — Frontend Public Site Design

## Overview

Public-facing website for BeautyMed clinic (Khabarovsk) built with Nuxt 3 SSR. The site serves as the patient acquisition funnel: it presents the clinic, doctors, and services, and lets patients book appointments online via a 5-step modal.

All content is fetched server-side from the existing Go backend. The frontend adds no separate database.

---

## Stack

| Concern | Choice |
|---|---|
| Framework | Nuxt 3 (SSR mode) |
| Language | TypeScript |
| Styling | Tailwind CSS v3 |
| State | Pinia |
| Testing (unit) | Vitest + @nuxt/test-utils |
| Testing (e2e) | Playwright |
| Type checking | nuxi typecheck |

---

## Architecture

**SSR + Backend as Single Source of Truth.**

`pages/index.vue` fetches all content server-side using `useAsyncData` with parallel requests:

```ts
const [{ data: doctors }, { data: specialties }, { data: services }, { data: reviews }] =
  await Promise.all([
    useAsyncData('doctors',     () => $fetch('/api/public/doctors')),
    useAsyncData('specialties', () => $fetch('/api/public/specialties')),
    useAsyncData('services',    () => $fetch('/api/public/services')),
    useAsyncData('reviews',     () => $fetch('/api/public/reviews')),
  ])
```

Data is passed as props to section components. No additional client-side fetches for content.

---

## New Public Endpoints (Go Backend)

These endpoints require no authentication and are added as part of Plan 3:

| Method | Path | Description |
|---|---|---|
| GET | `/api/public/specialties` | 4 specialties with name, description, icon |
| GET | `/api/public/doctors` | All doctors: name, photo URL, specialty, years of experience |
| GET | `/api/public/services` | Services grouped by specialty with prices |
| GET | `/api/public/reviews` | Patient reviews (text, rating, author) |

Existing endpoints used by the booking modal (no changes needed):

- `GET /api/slots?doctor_id=&date=` — available time slots
- `POST /api/auth/request-otp` — send OTP to phone
- `POST /api/auth/verify-otp` — verify OTP
- `POST /api/appointments` — create appointment

---

## File Structure

```
frontend/
├── pages/
│   └── index.vue                         # Main page, SSR data fetch
├── components/
│   ├── sections/
│   │   ├── HeroSection.vue
│   │   ├── AdvantagesSection.vue
│   │   ├── ServicesSection.vue
│   │   ├── DoctorsSection.vue
│   │   ├── ReviewsSection.vue
│   │   └── ContactsSection.vue
│   └── booking/
│       ├── BookingModal.vue              # Modal shell, step routing, dot indicator
│       └── steps/
│           ├── StepSpecialty.vue
│           ├── StepDoctor.vue
│           ├── StepDate.vue
│           ├── StepTime.vue
│           └── StepConfirm.vue
├── stores/
│   └── booking.ts                        # Pinia store for modal state
├── composables/
│   └── useApi.ts                         # Base $fetch wrapper with base URL
├── public/
│   ├── clinic_3.png                      # Hero photo (reception with logo)
│   └── clinic_2.png                      # Team photo (available as fallback)
├── nuxt.config.ts
├── tailwind.config.ts
├── Dockerfile
└── package.json
```

---

## Homepage Sections

| # | Section | Content source |
|---|---|---|
| 1 | **Hero** | Static (config): headline, subline, stats (10+ лет · 8 врачей · 3000+ пациентов). Photo: `clinic_3.png` right, text left. Two CTAs: «Записаться онлайн» (opens modal) + «Наши услуги» (scrolls). |
| 2 | **Advantages** | Static: 4–5 advantage cards with icon and short text (e.g., «Опытные врачи», «Современное оборудование») |
| 3 | **Services** | From `/api/public/services` — 4 specialty tabs with service list and prices |
| 4 | **Doctors** | From `/api/public/doctors` — grid of vertical cards (photo top, specialty badge overlay, name, experience, «Записаться» button) |
| 5 | **Reviews** | From `/api/public/reviews` — 3-column grid of review cards |
| 6 | **Contacts** | Static: address, phone, hours. Embedded Yandex Maps iframe. |

---

## Booking Modal

### Visual design
- Opens as centered modal overlay
- Step indicator: dot indicators (active dot is wider pill, completed dots filled teal)
- Footer: «← Назад» (ghost) + «Далее →» (teal button)

### 5 steps

| Step | Component | API |
|---|---|---|
| 1 | StepSpecialty | Uses `specialties` already fetched at SSR |
| 2 | StepDoctor | Filters `doctors` by selected specialty (client-side, no request) |
| 3 | StepDate | Calendar UI, client-side date picker |
| 4 | StepTime | `GET /api/slots?doctor_id=&date=` (client fetch on date select) |
| 5 | StepConfirm | Name + phone input → `POST /api/auth/request-otp` → OTP input → `POST /api/auth/verify-otp` → `POST /api/appointments` |

### Booking store (Pinia)

```ts
state: {
  open: boolean
  step: 1 | 2 | 3 | 4 | 5
  specialtyId: string | null
  doctorId: string | null
  date: string | null          // ISO date
  timeSlot: string | null
  name: string
  phone: string
  otpSent: boolean
}
actions: open(specialtyId?), close(), nextStep(), prevStep(), confirmBooking()
```

`open(specialtyId?)` — if specialtyId is provided, pre-selects specialty and jumps to step 2.

### Error handling
- Slot fetch failure: inline "Попробовать снова" button
- OTP wrong: inline error, retry after 60s cooldown
- Appointment POST failure: inline error with support phone number

---

## SEO

- `useHead()` per page: title, meta description, og:image
- JSON-LD structured data on index page: `MedicalClinic` schema with address, phone, opening hours

---

## Error Handling (SSR)

If the backend is unavailable during SSR, the page renders a fallback that shows the clinic phone number and address (hardcoded in `nuxt.config.ts` as `runtimeConfig.public`). All sections that failed to load are hidden gracefully.

---

## Testing

**Unit (Vitest):**
- `stores/booking.ts`: step transitions, pre-selected specialty skipping step 1, OTP cooldown state

**E2E (Playwright):**
- Open homepage → click «Записаться» → complete all 5 steps with test data → assert `POST /api/appointments` was called with correct payload

**CI:**
- `nuxi typecheck` — TypeScript check
- `vitest run` — unit tests
- Playwright tests against a running dev server

---

## Deployment

New `frontend` service added to existing `docker-compose.yml`:

```yaml
frontend:
  build: ./frontend
  environment:
    - NUXT_PUBLIC_API_BASE=http://api:8080
  depends_on:
    - api
```

Nginx updated: `location /` → `http://frontend:3000`, existing `location /api/` → `http://api:8080` unchanged.

Nuxt output mode: `standalone` (self-contained Node.js server, no separate static file serving needed).

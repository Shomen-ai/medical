# BeautyMed — Redesign Frontend (Plan 4)

## Overview

Visual redesign of the existing Nuxt 3 public site. No architectural changes — only Tailwind classes, component markup, and static assets are updated. The booking modal is untouched.

**Goal:** Replace the flat, cheap-looking default style with a Bold Modern aesthetic — gradient teal hero, card shadows, real stock doctor photos, improved typography.

---

## Design Decisions

| Question | Decision |
|---|---|
| Direction | Bold Modern (teal gradient, white CTAs, confident typography) |
| Hero scale | Compact Bold (A) |
| Doctor cards | Style A — gradient badge overlay, box-shadow, gradient button |
| Doctor photos | Unsplash stock (free commercial license), downloaded to `frontend/public/doctors/` |
| City | **Ульяновск** (replace «Хабаровск» everywhere: nuxt.config.ts, ContactsSection, HeroSection) |

---

## Color & Typography Changes

**Tailwind config** — no new colors needed; existing `primary` (#007C81), `slate`, `muted`, `border` stay. New utility: add `shadow-card` to tailwind config:

```ts
boxShadow: {
  card: '0 4px 24px rgba(0, 110, 115, 0.12)',
  'card-lg': '0 6px 32px rgba(0, 0, 0, 0.10)',
}
```

**Typography:** Section `<h2>` headings change from `text-xl` → `text-3xl font-extrabold`.

**Section backgrounds:** Alternating pattern:
- Hero: gradient
- Advantages: `bg-white`
- Services: `bg-[#F0FAFB]` (very light teal)
- Doctors: `bg-white`
- Reviews: `bg-[#F0FAFB]`
- Contacts: `bg-white`

---

## Component Changes

### `nuxt.config.ts`
- `clinicAddress`: replace «Хабаровск» → «Ульяновск»
- `clinicPhone`: no change (placeholder stays until real number provided)

### `HeroSection.vue`
**Before:** White background, plain border-b, small h1 (text-3xl), gray stat row.

**After:**
- Section bg: `background: linear-gradient(135deg, #005A5F 0%, #00959D 100%)`
- Eyebrow text: badge pill `bg-white/15 text-white` with «Клиника красоты · Ульяновск»
- h1: `text-3xl font-black text-white` (keep size, change weight + color)
- Subtext: `text-white/75`
- Primary CTA: `bg-white text-primary font-bold shadow-lg hover:shadow-xl`
- Secondary CTA: `border border-white/50 text-white hover:bg-white/10`
- Stats row: `bg-black/20` with white numbers and muted white labels
- Photo: add `rounded-xl shadow-2xl` to img

### `AdvantagesSection.vue`
**Before:** Plain emoji + text, no visual container.

**After:**
- Section heading: `text-3xl font-extrabold`
- Each advantage: white card with `rounded-2xl shadow-card p-6`
- Icon: replace emoji string with a `<div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center text-primary text-2xl mb-4">` wrapper
- Keep existing icons from `data/advantages.ts`

### `ServicesSection.vue`
**Before:** Plain pill tabs with border.

**After:**
- Section heading: `text-3xl font-extrabold`
- Active tab: `bg-gradient-to-r from-[#005A5F] to-[#00959D] text-white shadow-md` (instead of flat `bg-primary`)
- Service rows: add `hover:bg-[#F0FAFB]` transition on hover
- «Записаться» button per row: `bg-primary/10 text-primary hover:bg-primary hover:text-white transition-colors`

### `DoctorsSection.vue`
**Before:** Small cards, emoji placeholder, plain `bg-primary/10` button, no shadows.

**After:**
- Section heading: `text-3xl font-extrabold`
- Card: `rounded-2xl overflow-hidden shadow-card hover:shadow-card-lg transition-shadow`
- Photo area height: `h-48` (was `h-40`)
- `<img>`: real Unsplash URL per doctor (4 photos downloaded to `frontend/public/doctors/`)
- Specialty badge: `absolute bottom-0 inset-x-0 bg-gradient-to-r from-[#005A5F] to-[#00959D] text-white text-[9px] font-bold uppercase tracking-wide py-1.5 text-center` (full-width bar at bottom of photo)
- «Записаться» button: `bg-gradient-to-r from-[#005A5F] to-[#00959D] text-white font-bold rounded-lg py-2 hover:opacity-90 transition-opacity`

### `ReviewsSection.vue`
**Before:** Border-only cards.

**After:**
- Section heading: `text-3xl font-extrabold`
- Card: remove `border border-border`, add `shadow-card rounded-2xl`
- Add large opening quote `"` in `text-5xl text-primary/20 font-serif leading-none mb-2` at top of each card

### `ContactsSection.vue`
**Before:** Plain emoji inline with text.

**After:**
- Section heading: `text-3xl font-extrabold`
- Each contact item icon: wrap in `<div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary flex-shrink-0">` (teal circle)
- Address text: «Ульяновск» (corrected city)

---

## Doctor Stock Photos

Download 4 photos from Unsplash to `frontend/public/doctors/`:

| File | Unsplash URL | Role |
|---|---|---|
| `doctor-1.jpg` | `https://images.unsplash.com/photo-1559839734-2b71ea197ec2?w=400&h=480&fit=crop&q=85` | Female, cosmetology |
| `doctor-2.jpg` | `https://images.unsplash.com/photo-1612349317150-e413f6a5b16d?w=400&h=480&fit=crop&q=85` | Male, dermatology |
| `doctor-3.jpg` | `https://images.unsplash.com/photo-1594824476967-48c8b964273f?w=400&h=480&fit=crop&q=85` | Female, trichology |
| `doctor-4.jpg` | `https://images.unsplash.com/photo-1622253692010-333f2da6031d?w=400&h=480&fit=crop&q=85` | Male, aesthetic medicine |

In `DoctorsSection.vue`, map `doctor.photo_url` fallback: if `photo_url` is empty, use `/doctors/doctor-{n}.jpg` cycling by index.

---

## Files Changed

| File | Change type |
|---|---|
| `frontend/tailwind.config.ts` | Add `boxShadow.card` and `boxShadow.card-lg` |
| `frontend/nuxt.config.ts` | City: Ульяновск |
| `frontend/components/sections/HeroSection.vue` | Full restyle |
| `frontend/components/sections/AdvantagesSection.vue` | Card wrap + icon circle |
| `frontend/components/sections/ServicesSection.vue` | Gradient tab + hover row |
| `frontend/components/sections/DoctorsSection.vue` | New card style + stock photos |
| `frontend/components/sections/ReviewsSection.vue` | Shadow card + quote mark |
| `frontend/components/sections/ContactsSection.vue` | Icon circles + city fix |
| `frontend/public/doctors/doctor-{1-4}.jpg` | New assets (downloaded) |

**Not changed:** `BookingModal.vue`, all step components, `stores/booking.ts`, `pages/index.vue`, `composables/`, `types/`, `Dockerfile`, `docker-compose.yml`, `nginx.conf`.

---

## Testing

- `npm run typecheck` — no new types introduced, should pass unchanged
- `npm run test` — no store logic changed, 13 tests pass unchanged  
- Visual check: `npm run dev` → verify each section on desktop and mobile
- Deploy: `rsync` + `docker compose up --build -d frontend` on server `85.198.80.245`

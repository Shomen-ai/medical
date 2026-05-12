// Файл: tests/booking.store.test.ts
// Назначение: юнит-тесты Vitest для Pinia-стора booking — проверяют initial state, переходы между шагами, валидацию canProceed и сброс модалки.
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

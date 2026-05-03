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

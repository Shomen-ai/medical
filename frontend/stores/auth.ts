// Файл: stores/auth.ts
// Назначение: Pinia-стор авторизации — хранит JWT-токен и роль пользователя, выполняет отправку/проверку OTP-кода и выход из системы.
import { defineStore } from 'pinia'

type Role = 'patient' | 'admin' | 'doctor' | null

interface AuthState {
  token: string | null
  role: Role
  phone: string
  otpSent: boolean
  loading: boolean
  error: string | null
}

// Стор авторизации: токен, роль, состояние формы OTP и действия sendOTP/verifyOTP/logout.
export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: null,
    role: null,
    phone: '',
    otpSent: false,
    loading: false,
    error: null,
  }),

  getters: {
    isLoggedIn: (s) => !!s.token,
    isAdmin: (s) => s.role === 'admin',
    isDoctor: (s) => s.role === 'doctor',
    isPatient: (s) => s.role === 'patient',
  },

  actions: {
    init() {
      if (!import.meta.client) return
      const token = localStorage.getItem('auth_token')
      const role = localStorage.getItem('auth_role') as Role
      if (token && role) { this.token = token; this.role = role }
    },

    async sendOTP(phone: string, isStaff: boolean) {
      this.loading = true
      this.error = null
      try {
        const path = isStaff ? '/api/staff/auth/otp' : '/api/auth/otp'
        await $fetch(path, { method: 'POST', body: { phone } })
        this.phone = phone
        this.otpSent = true
      } catch {
        this.error = 'Ошибка при отправке кода'
      } finally {
        this.loading = false
      }
    },

    async verifyOTP(code: string, isStaff: boolean): Promise<Role> {
      this.loading = true
      this.error = null
      try {
        const path = isStaff ? '/api/staff/auth/verify' : '/api/auth/verify'
        const resp = await $fetch<{ access_token: string; role?: string }>(path, {
          method: 'POST',
          body: { phone: this.phone, code },
        })
        this.token = resp.access_token
        this.role = (resp.role as Role) ?? 'patient'
        if (import.meta.client) {
          localStorage.setItem('auth_token', this.token)
          localStorage.setItem('auth_role', this.role!)
        }
        this.otpSent = false
        return this.role
      } catch {
        this.error = 'Неверный код, попробуйте ещё раз'
        return null
      } finally {
        this.loading = false
      }
    },

    logout() {
      this.token = null
      this.role = null
      this.phone = ''
      this.otpSent = false
      if (import.meta.client) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_role')
      }
    },
  },
})

// Phone validation + normalization for RU (+7) and TM (+993).
//
// Accepted forms (digits only after country code):
//   RU:  +7 followed by 10 digits → 11 digits total starting with 7
//   TM:  +993 followed by 8 digits → 11 digits total starting with 993
// Russian "8XXXXXXXXXX" is also normalized to "+7XXXXXXXXXX".

export type PhoneCountry = 'ru' | 'tm'

export interface PhoneCheck {
  valid: boolean
  e164: string        // canonical form starting with +
  country: PhoneCountry | null
}

export function normalizePhoneInput(raw: string): string {
  const digits = raw.replace(/\D+/g, '')
  if (digits.startsWith('8') && digits.length === 11) return '+7' + digits.slice(1)
  if (digits.startsWith('7') && digits.length === 11) return '+' + digits
  if (digits.startsWith('993') && digits.length === 11) return '+' + digits
  if (raw.trim().startsWith('+')) return '+' + digits
  return digits ? '+' + digits : ''
}

export function validatePhone(raw: string): PhoneCheck {
  const e164 = normalizePhoneInput(raw)
  const digits = e164.replace(/\D+/g, '')

  if (digits.startsWith('7') && digits.length === 11) {
    return { valid: true, e164: '+' + digits, country: 'ru' }
  }
  if (digits.startsWith('993') && digits.length === 11) {
    return { valid: true, e164: '+' + digits, country: 'tm' }
  }
  return { valid: false, e164, country: null }
}

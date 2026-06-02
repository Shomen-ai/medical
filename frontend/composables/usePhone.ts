// Файл: composables/usePhone.ts
// Назначение: валидация и нормализация туркменских номеров (+993) —
// приведение к каноническому E.164. Префикс +993 в форме фиксирован,
// пользователь вводит только 8 национальных цифр (без пробелов и тире).

// Только Туркменистан: код страны +993 + 8 национальных цифр.
//   E.164 (каноника):  +99365123456     — 11 цифр, начинается с 993
//   Ввод в поле:        65123456         — 8 цифр, только цифры, без разделителей

// Результат валидации номера: флаг корректности и канонический E.164.
export interface PhoneCheck {
  valid: boolean
  e164: string
}

// Извлекает национальную часть (до 8 цифр без кода страны 993) из произвольного ввода.
// Код 993 отбрасывается только если цифр больше восьми (т.е. он реально присутствует).
function nationalDigits(raw: string): string {
  let d = raw.replace(/\D+/g, '')
  if (d.length > 8 && d.startsWith('993')) d = d.slice(3)
  return d.slice(0, 8)
}

// Приводит ввод к каноническому E.164 «+993XXXXXXXX» (пустая строка, если цифр нет).
export function normalizePhoneInput(raw: string): string {
  const d = nationalDigits(raw)
  return d ? '+993' + d : ''
}

// Оставляет в поле только цифры (до 8 национальных). Без пробелов и тире.
// Используется на каждое событие input, чтобы поле содержало строго цифры.
export function formatTmPhone(raw: string): string {
  return nationalDigits(raw)
}

// Валидирует номер: корректен только туркменский (8 национальных цифр,
// при наличии кода страны — ровно 993). Российские +7 и прочие отвергаются.
export function validatePhone(raw: string): PhoneCheck {
  const digits = raw.replace(/\D+/g, '')
  let national = ''
  if (digits.length === 8) national = digits
  else if (digits.length === 11 && digits.startsWith('993')) national = digits.slice(3)

  if (national.length === 8) return { valid: true, e164: '+993' + national }
  return { valid: false, e164: digits ? '+' + digits : '' }
}

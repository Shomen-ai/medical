// Файл: tests/phone.test.ts
// Назначение: юнит-тесты валидации/маски туркменских номеров (composables/usePhone).
import { describe, it, expect } from 'vitest'
import { validatePhone, normalizePhoneInput, formatTmPhone } from '../composables/usePhone'

describe('usePhone (Туркменистан, TM-only)', () => {
  it('форматирует 8 цифр в маску +993 65 12-34-56', () => {
    expect(formatTmPhone('65123456')).toBe('+993 65 12-34-56')
  })

  it('форматирует частичный ввод по мере набора', () => {
    expect(formatTmPhone('6')).toBe('+993 6')
    expect(formatTmPhone('65')).toBe('+993 65')
    expect(formatTmPhone('651')).toBe('+993 65 1')
    expect(formatTmPhone('6512')).toBe('+993 65 12')
    expect(formatTmPhone('651234')).toBe('+993 65 12-34')
  })

  it('отбрасывает уже введённый код страны 993', () => {
    expect(formatTmPhone('99365123456')).toBe('+993 65 12-34-56')
    expect(formatTmPhone('+993 65 12-34-56')).toBe('+993 65 12-34-56')
  })

  it('ограничивает 8 национальными цифрами', () => {
    expect(formatTmPhone('651234567890')).toBe('+993 65 12-34-56')
  })

  it('пустой ввод → пустая строка', () => {
    expect(formatTmPhone('')).toBe('')
    expect(normalizePhoneInput('')).toBe('')
  })

  it('нормализует к E.164', () => {
    expect(normalizePhoneInput('+993 65 12-34-56')).toBe('+99365123456')
    expect(normalizePhoneInput('65123456')).toBe('+99365123456')
  })

  it('валидирует полный туркменский номер', () => {
    const r = validatePhone('+993 65 12-34-56')
    expect(r.valid).toBe(true)
    expect(r.e164).toBe('+99365123456')
  })

  it('отклоняет неполный номер', () => {
    expect(validatePhone('+993 65 12').valid).toBe(false)
    expect(validatePhone('').valid).toBe(false)
  })

  it('отклоняет российский +7 (только TM)', () => {
    expect(validatePhone('+79161234567').valid).toBe(false)
  })
})

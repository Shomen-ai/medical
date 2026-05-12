// Файл: composables/useI18n.ts
// Назначение: i18n-композабл — хранит выбранный язык в cookie, возвращает функции t() (UI-переводы) и tMed() (медицинский глоссарий), переключает локаль.
import ru from '~/locales/ru'
import tk from '~/locales/tk'
import { medicalDict } from '~/locales/medical'

// Поддерживаемые локали интерфейса.
export type Locale = 'ru' | 'tk'

const dictionaries: Record<Locale, typeof ru> = { ru, tk }

const LOCALE_COOKIE = 'beautymed_locale'

// Substitutes {name} placeholders with values from params.
function interpolate(template: string, params?: Record<string, string | number>): string {
  if (!params) return template
  return template.replace(/\{(\w+)\}/g, (_, key) => {
    const v = params[key]
    return v === undefined ? `{${key}}` : String(v)
  })
}

// Композабл локализации: реактивная локаль из cookie + переводчики t/tMed/setLocale.
export function useI18n() {
  // useCookie is SSR-safe and synced to browser; default 'ru'.
  const locale = useCookie<Locale>(LOCALE_COOKIE, {
    default: () => 'ru',
    maxAge: 60 * 60 * 24 * 365,
  })

  const t = (key: keyof typeof ru, params?: Record<string, string | number>): string => {
    const dict = dictionaries[locale.value as Locale] ?? ru
    const tpl = dict[key] ?? ru[key] ?? key
    return interpolate(tpl, params)
  }

  // Look up a medical term (specialty / service name) in the current locale.
  // Falls back to the original Russian name if no translation exists.
  const tMed = (name: string): string => {
    const dict = medicalDict[locale.value as Locale] ?? {}
    return dict[name] ?? name
  }

  const setLocale = (next: Locale) => {
    locale.value = next
  }

  return { locale, t, tMed, setLocale }
}

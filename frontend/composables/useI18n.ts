import ru from '~/locales/ru'
import tk from '~/locales/tk'

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

  const setLocale = (next: Locale) => {
    locale.value = next
  }

  return { locale, t, setLocale }
}

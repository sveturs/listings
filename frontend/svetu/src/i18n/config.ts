export const locales = ['ru', 'en', 'sr'] as const;
export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = 'sr'; // Сербский по умолчанию для .rs домена

export const i18n = {
  locales,
  defaultLocale,
  localeDetection: {
    enabled: true,
    cookieName: 'locale-preference',
    cookieMaxAge: 365 * 24 * 60 * 60, // 1 год
  },
} as const;

export function getLocaleMessages(locale: Locale) {
  switch (locale) {
    case 'ru':
      return import('../messages/ru.json').then((module) => module.default);
    case 'en':
      return import('../messages/en.json').then((module) => module.default);
    case 'sr':
      return import('../messages/sr.json').then((module) => module.default);
    default:
      return import('../messages/sr.json').then((module) => module.default);
  }
}

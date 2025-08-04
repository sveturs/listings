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

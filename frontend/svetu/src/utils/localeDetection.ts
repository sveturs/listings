import { NextRequest } from 'next/server';
import { Locale, i18n } from '@/i18n/config';

interface LocaleDetectionConfig {
  cookieName: string;
  cookieMaxAge: number;
  defaultLocale: Locale;
}

export const localeConfig: LocaleDetectionConfig = {
  cookieName: 'locale-preference',
  cookieMaxAge: 365 * 24 * 60 * 60, // 1 год
  defaultLocale: 'sr', // Сербский по умолчанию для .rs домена
};

/**
 * Парсит заголовок Accept-Language и возвращает список языков с весами
 * Пример: "sr-RS,sr;q=0.9,ru;q=0.8,en-US;q=0.7,en;q=0.6"
 */
export function parseAcceptLanguage(
  acceptLanguage: string | null
): Array<{ locale: string; quality: number }> {
  if (!acceptLanguage) return [];

  const languages = acceptLanguage
    .split(',')
    .map((lang) => {
      const parts = lang.trim().split(';');
      const locale = parts[0].toLowerCase().split('-')[0]; // Берем только код языка без региона
      const quality = parts[1] ? parseFloat(parts[1].split('=')[1]) : 1.0;
      return { locale, quality };
    })
    .sort((a, b) => b.quality - a.quality);

  return languages;
}

/**
 * Определяет наиболее подходящий язык на основе:
 * 1. Cookie с сохраненным выбором пользователя
 * 2. Accept-Language заголовка браузера
 * 3. Дефолтного языка (сербский)
 */
export function detectLocale(request: NextRequest): Locale {
  // 1. Проверяем cookie с сохраненным выбором
  const cookieLocale = request.cookies.get(localeConfig.cookieName)?.value as
    | Locale
    | undefined;
  if (cookieLocale && i18n.locales.includes(cookieLocale)) {
    return cookieLocale;
  }

  // 2. Анализируем Accept-Language заголовок
  const acceptLanguage = request.headers.get('accept-language');
  const parsedLanguages = parseAcceptLanguage(acceptLanguage);

  // Ищем первый поддерживаемый язык из списка предпочтений
  for (const { locale } of parsedLanguages) {
    if (i18n.locales.includes(locale as Locale)) {
      return locale as Locale;
    }
  }

  // 3. Возвращаем дефолтный язык (сербский)
  return localeConfig.defaultLocale;
}

/**
 * Создает cookie для сохранения выбранного языка
 */
export function createLocaleCookie(locale: Locale): string {
  return `${localeConfig.cookieName}=${locale}; Max-Age=${localeConfig.cookieMaxAge}; Path=/; SameSite=Lax`;
}

/**
 * Проверяет, является ли строка валидной локалью
 */
export function isValidLocale(locale: string): locale is Locale {
  return i18n.locales.includes(locale as Locale);
}

/**
 * Извлекает локаль из пути URL
 * Пример: /ru/products -> ru
 */
export function getLocaleFromPathname(pathname: string): Locale | null {
  const segments = pathname.split('/');
  const potentialLocale = segments[1];

  if (potentialLocale && isValidLocale(potentialLocale)) {
    return potentialLocale;
  }

  return null;
}

/**
 * Убирает локаль из пути URL
 * Пример: /ru/products -> /products
 */
export function removeLocaleFromPathname(pathname: string): string {
  const locale = getLocaleFromPathname(pathname);
  if (locale) {
    const segments = pathname.split('/');
    segments.splice(1, 1);
    return segments.join('/') || '/';
  }
  return pathname;
}

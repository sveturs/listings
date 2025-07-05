import { NextRequest, NextResponse } from 'next/server';
import createIntlMiddleware from 'next-intl/middleware';
import { routing } from './i18n/routing';
import {
  detectLocale,
  getLocaleFromPathname,
  createLocaleCookie,
} from './utils/localeDetection';
import { i18n } from './i18n/config';

const intlMiddleware = createIntlMiddleware(routing);

export default function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  // Проверяем, есть ли локаль в URL
  const pathnameLocale = getLocaleFromPathname(pathname);

  // Если локаль уже есть в URL, просто передаем управление next-intl middleware
  if (pathnameLocale) {
    return intlMiddleware(request);
  }

  // Определяем подходящую локаль
  const detectedLocale = detectLocale(request);

  // Создаем новый URL с локалью
  const newUrl = new URL(`/${detectedLocale}${pathname}`, request.url);

  // Создаем response с редиректом
  const response = NextResponse.redirect(newUrl);

  // Если определили локаль по Accept-Language (а не из cookie),
  // сохраняем её в cookie для будущих визитов
  const cookieLocale = request.cookies.get(
    i18n.localeDetection.cookieName
  )?.value;
  if (!cookieLocale) {
    response.headers.set('Set-Cookie', createLocaleCookie(detectedLocale));
  }

  return response;
}

export const config = {
  // Исключаем из обработки:
  // - API routes
  // - Статические файлы (_next/static, _next/image, favicon.ico)
  // - Изображения и другие ресурсы
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp|ico|txt|pdf)).*)',
  ],
};

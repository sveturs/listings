import { NextRequest, NextResponse } from 'next/server';

/**
 * Получает правильный base URL для редиректов
 * Приоритет:
 * 1. Заголовок x-forwarded-host (для прокси/nginx)
 * 2. Заголовок host
 * 3. request.url origin
 * 4. Fallback на localhost (только для development)
 */
function getBaseUrl(request: NextRequest): string {
  // Проверяем заголовки от прокси/nginx
  const forwardedHost = request.headers.get('x-forwarded-host');
  const forwardedProto = request.headers.get('x-forwarded-proto') || 'https';

  if (forwardedHost) {
    return `${forwardedProto}://${forwardedHost}`;
  }

  // Используем заголовок Host
  const host = request.headers.get('host');
  if (host) {
    const protocol =
      request.headers.get('x-forwarded-proto') ||
      (host.includes('localhost') ? 'http' : 'https');
    return `${protocol}://${host}`;
  }

  // Fallback на origin из request.url
  const url = new URL(request.url);
  return url.origin;
}

/**
 * Google OAuth callback handler
 *
 * Этот эндпоинт получает токены от backend после успешной OAuth авторизации,
 * устанавливает их в httpOnly cookies и редиректит пользователя на frontend callback страницу.
 *
 * Flow:
 * 1. User нажимает "Login with Google"
 * 2. Backend редиректит на Google OAuth
 * 3. Google редиректит обратно на backend /api/v1/auth/google/callback
 * 4. Backend получает токены от Auth Service
 * 5. Backend редиректит сюда с токенами в query параметрах
 * 6. Мы устанавливаем токены в httpOnly cookies
 * 7. Редиректим на frontend /[locale]/auth/callback?success=true
 */
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);

    // Логируем все важные детали для отладки
    console.log('[GoogleCallback] Request details:', {
      url: request.url,
      host: request.headers.get('host'),
      'x-forwarded-host': request.headers.get('x-forwarded-host'),
      'x-forwarded-proto': request.headers.get('x-forwarded-proto'),
      origin: request.headers.get('origin'),
      referer: request.headers.get('referer'),
    });

    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');
    const locale = searchParams.get('locale') || 'en';
    const returnUrl = searchParams.get('return_url') || '/';

    // Валидация: токены обязательны
    if (!accessToken || !refreshToken) {
      console.error('[GoogleCallback] Missing tokens in callback');

      // Используем явный base URL для редиректа
      const baseUrl = getBaseUrl(request);
      return NextResponse.redirect(
        new URL(`/${locale}/login?error=auth_failed`, baseUrl)
      );
    }

    console.log('[GoogleCallback] Tokens received, setting cookies');

    // Создаем response с редиректом на frontend callback
    // ВАЖНО: используем явный base URL вместо request.url
    const baseUrl = getBaseUrl(request);
    const redirectUrl = new URL(`/${locale}/auth/callback`, baseUrl);
    redirectUrl.searchParams.set('success', 'true');
    redirectUrl.searchParams.set('return_url', returnUrl);

    console.log('[GoogleCallback] Redirecting to:', redirectUrl.toString());

    const response = NextResponse.redirect(redirectUrl);

    // Устанавливаем access token в httpOnly cookie
    response.cookies.set({
      name: 'access_token',
      value: accessToken,
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 60 * 60 * 48, // 48 hours (соответствует времени жизни токена)
    });

    // Устанавливаем refresh token в httpOnly cookie
    response.cookies.set({
      name: 'refresh_token',
      value: refreshToken,
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 60 * 60 * 24 * 30, // 30 days
    });

    console.log(
      '[GoogleCallback] Cookies set, redirecting to frontend callback'
    );

    return response;
  } catch (error) {
    console.error('[GoogleCallback] Error processing callback:', error);

    // Используем явный base URL для редиректа при ошибке
    const baseUrl = getBaseUrl(request);
    return NextResponse.redirect(new URL('/login?error=auth_failed', baseUrl));
  }
}

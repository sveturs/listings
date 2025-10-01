import { NextRequest, NextResponse } from 'next/server';

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

    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');
    const locale = searchParams.get('locale') || 'en';
    const returnUrl = searchParams.get('return_url') || '/';

    // Валидация: токены обязательны
    if (!accessToken || !refreshToken) {
      console.error('[GoogleCallback] Missing tokens in callback');
      return NextResponse.redirect(
        new URL(`/${locale}/login?error=auth_failed`, request.url)
      );
    }

    console.log('[GoogleCallback] Tokens received, setting cookies');

    // Создаем response с редиректом на frontend callback
    const redirectUrl = new URL(`/${locale}/auth/callback`, request.url);
    redirectUrl.searchParams.set('success', 'true');
    redirectUrl.searchParams.set('return_url', returnUrl);

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

    console.log('[GoogleCallback] Cookies set, redirecting to frontend callback');

    return response;
  } catch (error) {
    console.error('[GoogleCallback] Error processing callback:', error);
    return NextResponse.redirect(
      new URL('/login?error=auth_failed', request.url)
    );
  }
}

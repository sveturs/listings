import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

export async function GET() {
  const cookieStore = await cookies();
  const allCookies = cookieStore.getAll();

  return NextResponse.json({
    cookies: allCookies.map(c => ({
      name: c.name,
      value: c.value.substring(0, 20) + '...', // Первые 20 символов для безопасности
      hasValue: !!c.value,
    })),
    access_token_exists: !!cookieStore.get('access_token')?.value,
    refresh_token_exists: !!cookieStore.get('refresh_token')?.value,
  });
}

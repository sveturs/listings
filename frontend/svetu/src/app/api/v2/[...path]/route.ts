import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const BACKEND_URL =
  process.env.BACKEND_INTERNAL_URL || 'http://localhost:33423';

/**
 * Универсальный BFF прокси для всех backend API запросов
 *
 * Маппинг:
 * /api/v2/* → backend /api/v1/*
 *
 * Преимущества:
 * - Автоматически добавляет JWT токены из httpOnly cookies
 * - Нет CORS проблем (все на одном домене)
 * - Централизованная обработка ошибок
 * - Легко добавить rate limiting, logging, caching
 */

async function proxyRequest(
  request: NextRequest,
  method: string,
  path: string[]
) {
  try {
    // Получаем cookies
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;
    const allCookies = cookieStore.getAll();

    // Строим URL для backend
    const backendPath = path.join('/');
    const searchParams = request.nextUrl.searchParams.toString();
    const backendUrl = `${BACKEND_URL}/api/v1/${backendPath}${searchParams ? `?${searchParams}` : ''}`;

    console.log(`[BFF Proxy] ${method} /api/v2/${backendPath} → ${backendUrl}`);
    console.log(
      `[BFF Proxy] Access token: ${accessToken ? 'present (length: ' + accessToken.length + ')' : 'MISSING'}`
    );
    console.log(
      `[BFF Proxy] All cookies:`,
      allCookies.map((c) => c.name).join(', ')
    );

    // Подготавливаем заголовки
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    // Добавляем Authorization header если есть токен
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    // Копируем важные заголовки из оригинального запроса
    const acceptLanguage = request.headers.get('Accept-Language');
    if (acceptLanguage) {
      headers['Accept-Language'] = acceptLanguage;
    }

    // Подготавливаем body для методов с телом запроса
    let body: string | undefined;
    if (['POST', 'PUT', 'PATCH'].includes(method)) {
      try {
        const requestBody = await request.json();
        body = JSON.stringify(requestBody);
      } catch {
        // Если не удалось распарсить JSON, пропускаем body
        body = undefined;
      }
    }

    // Выполняем запрос к backend
    const response = await fetch(backendUrl, {
      method,
      headers,
      body,
      // Не используем credentials здесь, так как это server-to-server запрос
    });

    // Получаем данные ответа
    const contentType = response.headers.get('content-type');
    let data;

    if (contentType?.includes('application/json')) {
      const text = await response.text();
      if (text) {
        try {
          data = JSON.parse(text);
        } catch (e) {
          console.error('[BFF Proxy] Failed to parse JSON response:', e);
          data = { error: 'Invalid JSON response from backend' };
        }
      }
    } else {
      // Для не-JSON ответов просто возвращаем текст
      data = await response.text();
    }

    // Логируем ошибки
    if (!response.ok) {
      console.error(`[BFF Proxy] Backend error: ${response.status}`, data);
    }

    // Возвращаем ответ с тем же статусом что и backend
    return NextResponse.json(data, {
      status: response.status,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  } catch (error) {
    console.error('[BFF Proxy] Request failed:', error);

    return NextResponse.json(
      {
        error: 'Proxy request failed',
        message: error instanceof Error ? error.message : 'Unknown error',
      },
      { status: 500 }
    );
  }
}

// GET handler
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, 'GET', path);
}

// POST handler
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, 'POST', path);
}

// PUT handler
export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, 'PUT', path);
}

// DELETE handler
export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, 'DELETE', path);
}

// PATCH handler
export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, 'PATCH', path);
}

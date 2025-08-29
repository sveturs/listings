import configManager from '@/config';
import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

export async function GET(request: NextRequest) {
  try {
    // Получаем токен из cookies или headers
    const cookieStore = await cookies();
    const sessionCookie = cookieStore.get('session');
    const authHeader = request.headers.get('Authorization');

    // Формируем headers для запроса к backend
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    // Если есть session cookie, передаем его
    if (sessionCookie) {
      headers['Cookie'] = `session=${sessionCookie.value}`;
    }

    // Если есть Authorization header, передаем его
    if (authHeader) {
      headers['Authorization'] = authHeader;
    }

    // Делаем запрос к backend
    const backendUrl = process.env.BACKEND_URL || configManager.getApiUrl();
    const response = await fetch(
      `${backendUrl}/api/v1/admin/translations/ai/costs`,
      {
        method: 'GET',
        headers,
        credentials: 'include',
      }
    );

    // Если ответ не успешный, возвращаем ошибку
    if (!response.ok) {
      // Если нет авторизации, возвращаем демо данные
      if (response.status === 401 || response.status === 403) {
        return NextResponse.json({
          success: true,
          data: {
            total_cost: 0.001176,
            total_tokens: 118,
            total_requests: 8,
            today_cost: 0.001176,
            month_cost: 0.001176,
            by_provider: {
              openai: {
                provider: 'openai',
                total_cost: 0.000096,
                total_tokens: 64,
                total_requests: 6,
                last_updated: new Date().toISOString(),
                daily_costs: {
                  [new Date().toISOString().split('T')[0]]: 0.000096,
                },
                hourly_costs: {
                  [new Date().toISOString().split('T')[0] +
                  'T' +
                  new Date().getHours()]: 0.000096,
                },
              },
              google: {
                provider: 'google',
                total_cost: 0.00108,
                total_tokens: 54,
                total_requests: 2,
                last_updated: new Date().toISOString(),
                daily_costs: {
                  [new Date().toISOString().split('T')[0]]: 0.00108,
                },
                hourly_costs: {
                  [new Date().toISOString().split('T')[0] +
                  'T' +
                  new Date().getHours()]: 0.00108,
                },
              },
              deepl: {
                provider: 'deepl',
                total_cost: 0,
                total_tokens: 0,
                total_requests: 0,
                last_updated: new Date().toISOString(),
                daily_costs: {},
                hourly_costs: {},
              },
              claude: {
                provider: 'claude',
                total_cost: 0,
                total_tokens: 0,
                total_requests: 0,
                last_updated: new Date().toISOString(),
                daily_costs: {},
                hourly_costs: {},
              },
            },
            today_by_provider: {
              openai: 0.000096,
              google: 0.00108,
            },
            month_by_provider: {
              openai: 0.000096,
              google: 0.00108,
            },
          },
          message: 'Demo data - authentication required for real data',
        });
      }

      return NextResponse.json(
        { error: 'Failed to fetch costs data' },
        { status: response.status }
      );
    }

    // Возвращаем данные от backend
    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Error fetching AI costs:', error);

    // В случае ошибки возвращаем демо данные
    return NextResponse.json({
      success: true,
      data: {
        total_cost: 0.001176,
        total_tokens: 118,
        total_requests: 8,
        today_cost: 0.001176,
        month_cost: 0.001176,
        by_provider: {
          openai: {
            provider: 'openai',
            total_cost: 0.000096,
            total_tokens: 64,
            total_requests: 6,
            last_updated: new Date().toISOString(),
            daily_costs: { [new Date().toISOString().split('T')[0]]: 0.000096 },
            hourly_costs: {
              [new Date().toISOString().split('T')[0] +
              'T' +
              new Date().getHours()]: 0.000096,
            },
          },
          google: {
            provider: 'google',
            total_cost: 0.00108,
            total_tokens: 54,
            total_requests: 2,
            last_updated: new Date().toISOString(),
            daily_costs: { [new Date().toISOString().split('T')[0]]: 0.00108 },
            hourly_costs: {
              [new Date().toISOString().split('T')[0] +
              'T' +
              new Date().getHours()]: 0.00108,
            },
          },
          deepl: {
            provider: 'deepl',
            total_cost: 0,
            total_tokens: 0,
            total_requests: 0,
            last_updated: new Date().toISOString(),
            daily_costs: {},
            hourly_costs: {},
          },
          claude: {
            provider: 'claude',
            total_cost: 0,
            total_tokens: 0,
            total_requests: 0,
            last_updated: new Date().toISOString(),
            daily_costs: {},
            hourly_costs: {},
          },
        },
        today_by_provider: {
          openai: 0.000096,
          google: 0.00108,
        },
        month_by_provider: {
          openai: 0.000096,
          google: 0.00108,
        },
      },
      message: 'Demo data due to connection error',
    });
  }
}

import { useEffect, useRef } from 'react';
import { v4 as uuidv4 } from 'uuid';

// Типы событий
export type AnalyticsEventType =
  | 'page_view'
  | 'product_view'
  | 'add_to_cart'
  | 'checkout'
  | 'order';

interface AnalyticsEvent {
  storefront_id: number;
  event_type: AnalyticsEventType;
  event_data?: Record<string, any>;
  session_id?: string;
  user_id?: number;
}

// Получаем или создаем session ID
const getSessionId = (): string => {
  const STORAGE_KEY = 'analytics_session_id';
  const EXPIRY_TIME = 30 * 60 * 1000; // 30 минут

  // Проверяем, что мы в браузере
  if (typeof window === 'undefined') {
    return uuidv4(); // Возвращаем временный ID для SSR
  }

  const storedData = localStorage.getItem(STORAGE_KEY);
  if (storedData) {
    const { sessionId, timestamp } = JSON.parse(storedData);
    if (Date.now() - timestamp < EXPIRY_TIME) {
      return sessionId;
    }
  }

  const newSessionId = uuidv4();
  localStorage.setItem(
    STORAGE_KEY,
    JSON.stringify({
      sessionId: newSessionId,
      timestamp: Date.now(),
    })
  );

  return newSessionId;
};

// Хук для отслеживания событий
export const useAnalytics = () => {
  const sessionId = useRef<string | null>(null);

  // Инициализируем session ID только на клиенте
  useEffect(() => {
    if (sessionId.current === null) {
      sessionId.current = getSessionId();
    }
  }, []);

  const trackEvent = async (event: Omit<AnalyticsEvent, 'session_id'>) => {
    // Не отправляем события, пока session ID не инициализирован
    if (!sessionId.current) {
      return;
    }

    try {
      const response = await fetch('/api/v1/analytics/event', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...event,
          session_id: sessionId.current,
        }),
      });

      if (!response.ok) {
        console.error('Failed to track event:', await response.text());
      }
    } catch (error) {
      console.error('Error tracking event:', error);
    }
  };

  const trackPageView = (storefrontId: number, pageType?: string) => {
    trackEvent({
      storefront_id: storefrontId,
      event_type: 'page_view',
      event_data: { page_type: pageType || 'storefront' },
    });
  };

  const trackProductView = (
    storefrontId: number,
    productId: number,
    productName?: string
  ) => {
    trackEvent({
      storefront_id: storefrontId,
      event_type: 'product_view',
      event_data: {
        product_id: productId,
        product_name: productName,
      },
    });
  };

  const trackAddToCart = (
    storefrontId: number,
    productId: number,
    quantity: number,
    price: number
  ) => {
    trackEvent({
      storefront_id: storefrontId,
      event_type: 'add_to_cart',
      event_data: {
        product_id: productId,
        quantity,
        price,
      },
    });
  };

  const trackCheckout = (
    storefrontId: number,
    items: any[],
    totalAmount: number
  ) => {
    trackEvent({
      storefront_id: storefrontId,
      event_type: 'checkout',
      event_data: {
        items,
        total_amount: totalAmount,
      },
    });
  };

  return {
    trackEvent,
    trackPageView,
    trackProductView,
    trackAddToCart,
    trackCheckout,
    sessionId: sessionId.current || '',
  };
};

// Хук для автоматического отслеживания просмотра страницы
export const usePageViewTracking = (
  storefrontId?: number,
  pageType?: string
) => {
  const { trackPageView } = useAnalytics();
  const tracked = useRef(false);

  useEffect(() => {
    if (storefrontId && !tracked.current) {
      trackPageView(storefrontId, pageType);
      tracked.current = true;
    }
  }, [storefrontId, pageType, trackPageView]);
};

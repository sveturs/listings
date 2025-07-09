'use client';

/**
 * Hook для поведенческого трекинга пользователей
 *
 * Особенности:
 * - Автоматическая генерация session_id
 * - Батчевая отправка событий (каждые 5 секунд или при 10 событиях)
 * - Обработка ошибок и retry логика
 * - Поддержка анонимных пользователей
 * - Дебаунс для предотвращения спама событий
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { useLocalStorage } from './useLocalStorage';
import { apiClient } from '../services/api-client';
import { getUserAgentInfo } from '../utils/deviceDetection';
import type {
  BehaviorEvent,
  BehaviorEventType,
  BehaviorEventBatch,
  BehaviorTrackingConfig,
  BehaviorTrackingResponse,
  BehaviorTrackingState,
  TrackingContext,
  TrackingOptions,
  UseBehaviorTrackingOptions,
  UseBehaviorTrackingReturn,
  BaseBehaviorEvent,
  SearchPerformedEvent,
  ResultClickedEvent,
  ItemViewedEvent,
  ItemPurchasedEvent,
  SearchFilterAppliedEvent,
  SearchSortChangedEvent,
} from '../types/behavior';

/**
 * Конфигурация по умолчанию
 */
const DEFAULT_CONFIG: BehaviorTrackingConfig = {
  enabled: true,
  batchSize: 10,
  batchTimeout: 5000, // 5 секунд
  maxRetries: 3,
  endpoint: '/api/v1/analytics/track',
};

/**
 * Генерирует уникальный session_id
 */
function generateSessionId(): string {
  return `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Генерирует уникальный batch_id
 */
function generateBatchId(): string {
  return `batch_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Дебаунс для предотвращения спама событий
 */
function debounce<T extends (...args: any[]) => void>(
  func: T,
  wait: number
): T {
  let timeout: NodeJS.Timeout;
  return ((...args: any[]) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  }) as T;
}

/**
 * Hook для поведенческого трекинга
 */
export function useBehaviorTracking(
  options: UseBehaviorTrackingOptions = {}
): UseBehaviorTrackingReturn {
  const {
    enabled = true,
    batchSize = DEFAULT_CONFIG.batchSize,
    batchTimeout = DEFAULT_CONFIG.batchTimeout,
    maxRetries = DEFAULT_CONFIG.maxRetries,
    endpoint = DEFAULT_CONFIG.endpoint,
    userId,
    autoInit = true,
  } = options;

  // Конфигурация
  const config: BehaviorTrackingConfig = {
    enabled,
    batchSize,
    batchTimeout,
    maxRetries,
    endpoint,
  };

  // Состояние
  const [state, setState] = useState<BehaviorTrackingState>({
    isEnabled: enabled,
    isInitialized: false,
    pendingEvents: [],
    context: {
      sessionId: '',
    },
    retryCount: 0,
  });

  // Хранение session_id в localStorage
  const [sessionId, setSessionId] = useLocalStorage<string>(
    'behavior_tracking_session_id',
    '',
    { userId, namespace: 'svetu_tracking' }
  );

  // Refs для таймеров и состояния
  const batchTimeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const retryTimeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const userAgentInfoRef = useRef(getUserAgentInfo());

  /**
   * Инициализация hook'а
   */
  const initialize = useCallback(() => {
    if (state.isInitialized) return;

    let currentSessionId = sessionId;

    // Генерируем session_id если его нет
    if (!currentSessionId) {
      currentSessionId = generateSessionId();
      setSessionId(currentSessionId);
    }

    setState((prev) => ({
      ...prev,
      isInitialized: true,
      context: {
        ...prev.context,
        sessionId: currentSessionId,
      },
    }));
  }, [sessionId, setSessionId, state.isInitialized]);

  /**
   * Создает базовое событие с общими данными
   */
  const createBaseEvent = useCallback(
    <T extends BehaviorEventType>(
      eventType: T
    ): Omit<BaseBehaviorEvent, 'event_type'> & { event_type: T } => {
      const userAgentInfo = userAgentInfoRef.current;

      return {
        session_id: state.context.sessionId,
        event_type: eventType,
        timestamp: new Date().toISOString(),
        user_id: userId,
        device_type: userAgentInfo.device.type,
        browser: userAgentInfo.browser.full,
        referer: typeof window !== 'undefined' ? document.referrer : undefined,
        user_agent: userAgentInfo.userAgent,
      };
    },
    [state.context.sessionId, userId]
  );

  /**
   * Отправляет события на сервер
   */
  const sendEvents = useCallback(
    async (events: BehaviorEvent[]): Promise<boolean> => {
      if (!config.enabled || events.length === 0) return true;

      const batch: BehaviorEventBatch = {
        events,
        batch_id: generateBatchId(),
        created_at: new Date().toISOString(),
      };

      try {
        const response = await apiClient.post<BehaviorTrackingResponse>(
          config.endpoint,
          batch,
          {
            timeout: 10000, // 10 секунд
            retries: 0, // Retry логика обрабатывается отдельно
          }
        );

        if (response.error) {
          console.error('[BehaviorTracking] API Error:', response.error);
          return false;
        }

        console.log('[BehaviorTracking] Events sent successfully:', {
          count: events.length,
          batchId: batch.batch_id,
          processed: response.data?.processed_count,
        });

        return true;
      } catch (error) {
        console.error('[BehaviorTracking] Network Error:', error);
        return false;
      }
    },
    [config.enabled, config.endpoint]
  );

  /**
   * Отправляет события с retry логикой
   */
  const sendEventsWithRetry = useCallback(
    async (events: BehaviorEvent[]): Promise<void> => {
      if (!config.enabled || events.length === 0) return;

      const success = await sendEvents(events);

      if (!success && state.retryCount < config.maxRetries) {
        setState((prev) => ({
          ...prev,
          retryCount: prev.retryCount + 1,
          lastError: `Failed to send events, retry ${prev.retryCount + 1}/${config.maxRetries}`,
        }));

        // Экспоненциальная задержка для повторных попыток
        const delay = Math.min(1000 * Math.pow(2, state.retryCount), 10000);

        retryTimeoutRef.current = setTimeout(() => {
          sendEventsWithRetry(events);
        }, delay);
      } else if (success) {
        setState((prev) => ({
          ...prev,
          retryCount: 0,
          lastError: undefined,
        }));
      } else {
        setState((prev) => ({
          ...prev,
          lastError: `Failed to send events after ${config.maxRetries} retries`,
        }));
      }
    },
    [config.enabled, config.maxRetries, sendEvents, state.retryCount]
  );

  /**
   * Добавляет событие в очередь
   */
  const addEvent = useCallback(
    (event: BehaviorEvent, options: TrackingOptions = {}) => {
      if (!config.enabled || !state.isInitialized) return;

      setState((prev) => ({
        ...prev,
        pendingEvents: [...prev.pendingEvents, event],
      }));

      // Отправка немедленно, если запрошено
      if (options.immediate) {
        sendEventsWithRetry([event]);
        return;
      }

      // Отправка при достижении лимита размера батча
      if (state.pendingEvents.length + 1 >= config.batchSize) {
        flushEvents();
        return;
      }

      // Установка таймера для отправки по времени
      if (batchTimeoutRef.current) {
        clearTimeout(batchTimeoutRef.current);
      }

      batchTimeoutRef.current = setTimeout(() => {
        flushEvents();
      }, config.batchTimeout);
    },
    [
      config.enabled,
      config.batchSize,
      config.batchTimeout,
      state.isInitialized,
      state.pendingEvents.length,
    ]
  );

  /**
   * Отправляет все накопленные события
   */
  const flushEvents = useCallback(async () => {
    if (state.pendingEvents.length === 0) return;

    const eventsToSend = [...state.pendingEvents];

    setState((prev) => ({
      ...prev,
      pendingEvents: [],
    }));

    if (batchTimeoutRef.current) {
      clearTimeout(batchTimeoutRef.current);
    }

    await sendEventsWithRetry(eventsToSend);
  }, [state.pendingEvents, sendEventsWithRetry]);

  /**
   * Очищает контекст трекинга
   */
  const clearContext = useCallback(() => {
    setState((prev) => ({
      ...prev,
      context: {
        sessionId: prev.context.sessionId,
      },
    }));
  }, []);

  /**
   * Обновляет контекст трекинга
   */
  const updateContext = useCallback((updates: Partial<TrackingContext>) => {
    setState((prev) => ({
      ...prev,
      context: {
        ...prev.context,
        ...updates,
      },
    }));
  }, []);

  // Дебаунсированные методы трекинга
  const debouncedAddEvent = useCallback(debounce(addEvent, 100), [addEvent]);

  /**
   * Методы трекинга конкретных событий
   */
  const trackSearchPerformed = useCallback(
    async (data: Omit<SearchPerformedEvent, keyof BaseBehaviorEvent>) => {
      const event: SearchPerformedEvent = {
        ...createBaseEvent('search_performed'),
        ...data,
      };

      // Обновляем контекст
      updateContext({
        currentSearchQuery: data.search_query,
        currentSearchStartTime: Date.now(),
        lastSearchResults: data.results_count,
        searchFilters: data.search_filters,
        searchSort: data.search_sort,
      });

      // Используем прямой addEvent вместо debounced для важных событий поиска
      addEvent(event);
    },
    [createBaseEvent, updateContext, addEvent]
  );

  const trackResultClicked = useCallback(
    async (data: Omit<ResultClickedEvent, keyof BaseBehaviorEvent>) => {
      const event: ResultClickedEvent = {
        ...createBaseEvent('result_clicked'),
        ...data,
      };

      addEvent(event);
    },
    [createBaseEvent, addEvent]
  );

  const trackItemViewed = useCallback(
    async (data: Omit<ItemViewedEvent, keyof BaseBehaviorEvent>) => {
      const event: ItemViewedEvent = {
        ...createBaseEvent('item_viewed'),
        ...data,
      };

      // Обновляем контекст
      updateContext({
        currentItemId: data.item_id,
        currentItemViewStartTime: Date.now(),
      });

      debouncedAddEvent(event);
    },
    [createBaseEvent, updateContext, debouncedAddEvent]
  );

  const trackItemPurchased = useCallback(
    async (data: Omit<ItemPurchasedEvent, keyof BaseBehaviorEvent>) => {
      const event: ItemPurchasedEvent = {
        ...createBaseEvent('item_purchased'),
        ...data,
      };

      // Отправляем покупки немедленно
      addEvent(event, { immediate: true });
    },
    [createBaseEvent, addEvent]
  );

  const trackSearchFilterApplied = useCallback(
    async (data: Omit<SearchFilterAppliedEvent, keyof BaseBehaviorEvent>) => {
      const event: SearchFilterAppliedEvent = {
        ...createBaseEvent('search_filter_applied'),
        ...data,
      };

      // Обновляем контекст
      updateContext({
        searchFilters: {
          ...state.context.searchFilters,
          [data.filter_type]: data.filter_value,
        },
      });

      debouncedAddEvent(event);
    },
    [
      createBaseEvent,
      updateContext,
      debouncedAddEvent,
      state.context.searchFilters,
    ]
  );

  const trackSearchSortChanged = useCallback(
    async (data: Omit<SearchSortChangedEvent, keyof BaseBehaviorEvent>) => {
      const event: SearchSortChangedEvent = {
        ...createBaseEvent('search_sort_changed'),
        ...data,
      };

      // Обновляем контекст
      updateContext({
        searchSort: data.sort_type,
      });

      debouncedAddEvent(event);
    },
    [createBaseEvent, updateContext, debouncedAddEvent]
  );

  /**
   * Вспомогательные методы для специфических use cases
   */
  const startSearch = useCallback(
    (query: string, filters?: Record<string, any>, sort?: string) => {
      updateContext({
        currentSearchQuery: query,
        currentSearchStartTime: Date.now(),
        searchFilters: filters,
        searchSort: sort,
      });
    },
    [updateContext]
  );

  const startItemView = useCallback(
    (itemId: string) => {
      updateContext({
        currentItemId: itemId,
        currentItemViewStartTime: Date.now(),
      });
    },
    [updateContext]
  );

  const endItemView = useCallback(() => {
    const { currentItemId, currentItemViewStartTime } = state.context;

    if (currentItemId && currentItemViewStartTime) {
      const viewDuration = Date.now() - currentItemViewStartTime;

      // Можно отправить событие с продолжительностью просмотра
      console.log('[BehaviorTracking] Item view ended:', {
        itemId: currentItemId,
        duration: viewDuration,
      });
    }

    updateContext({
      currentItemId: undefined,
      currentItemViewStartTime: undefined,
    });
  }, [state.context, updateContext]);

  // Инициализация при монтировании
  useEffect(() => {
    if (autoInit) {
      initialize();
    }
  }, [autoInit, initialize]);

  // Отправка событий при unmount
  useEffect(() => {
    return () => {
      if (batchTimeoutRef.current) {
        clearTimeout(batchTimeoutRef.current);
      }
      if (retryTimeoutRef.current) {
        clearTimeout(retryTimeoutRef.current);
      }

      // Отправляем оставшиеся события
      if (state.pendingEvents.length > 0) {
        sendEvents(state.pendingEvents);
      }
    };
  }, [state.pendingEvents, sendEvents]);

  // Отправка событий при закрытии страницы
  useEffect(() => {
    const handleBeforeUnload = () => {
      if (state.pendingEvents.length > 0) {
        // Используем navigator.sendBeacon для надежной отправки
        if (navigator.sendBeacon) {
          const batch: BehaviorEventBatch = {
            events: state.pendingEvents,
            batch_id: generateBatchId(),
            created_at: new Date().toISOString(),
          };

          navigator.sendBeacon(config.endpoint, JSON.stringify(batch));
        }
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [state.pendingEvents, config.endpoint]);

  return {
    state,
    config,
    trackSearchPerformed,
    trackResultClicked,
    trackItemViewed,
    trackItemPurchased,
    trackSearchFilterApplied,
    trackSearchSortChanged,
    flushEvents,
    clearContext,
    updateContext,
    startSearch,
    startItemView,
    endItemView,
  };
}

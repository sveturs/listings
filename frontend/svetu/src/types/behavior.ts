/**
 * Типы для поведенческого трекинга пользователей
 */

/**
 * Типы событий для трекинга
 */
export type BehaviorEventType =
  | 'search_performed'
  | 'result_clicked'
  | 'item_viewed'
  | 'item_purchased'
  | 'search_filter_applied'
  | 'search_sort_changed';

/**
 * Типы устройств
 */
export type DeviceType = 'desktop' | 'mobile' | 'tablet';

/**
 * Базовые данные события
 */
export interface BaseBehaviorEvent {
  session_id: string;
  event_type: BehaviorEventType;
  timestamp: string;
  user_id?: string;
  device_type: DeviceType;
  browser: string;
  referer?: string;
  user_agent: string;
  ip_address?: string;
}

/**
 * Данные события поиска
 */
export interface SearchPerformedEvent extends BaseBehaviorEvent {
  event_type: 'search_performed';
  search_query: string;
  search_filters?: Record<string, any>;
  search_sort?: string;
  results_count: number;
  search_duration_ms?: number;
  item_type?: 'marketplace' | 'storefront';
}

/**
 * Данные события клика по результату поиска
 */
export interface ResultClickedEvent extends BaseBehaviorEvent {
  event_type: 'result_clicked';
  search_query: string;
  clicked_item_id: string;
  click_position: number;
  total_results: number;
  click_time_from_search_ms: number;
  item_type: 'marketplace' | 'storefront';
}

/**
 * Данные события просмотра товара
 */
export interface ItemViewedEvent extends BaseBehaviorEvent {
  event_type: 'item_viewed';
  item_id: string;
  item_title: string;
  item_category?: string;
  item_price?: number;
  view_duration_ms?: number;
  previous_page?: string;
}

/**
 * Данные события покупки
 */
export interface ItemPurchasedEvent extends BaseBehaviorEvent {
  event_type: 'item_purchased';
  item_id: string;
  item_title: string;
  item_category?: string;
  purchase_amount: number;
  currency: string;
  payment_method?: string;
  time_from_first_view_ms?: number;
}

/**
 * Данные события применения фильтра
 */
export interface SearchFilterAppliedEvent extends BaseBehaviorEvent {
  event_type: 'search_filter_applied';
  search_query: string;
  filter_type: string;
  filter_value: string;
  results_count_before: number;
  results_count_after: number;
}

/**
 * Данные события изменения сортировки
 */
export interface SearchSortChangedEvent extends BaseBehaviorEvent {
  event_type: 'search_sort_changed';
  search_query: string;
  sort_type: string;
  previous_sort?: string;
  results_count: number;
}

/**
 * Объединенный тип для всех событий
 */
export type BehaviorEvent =
  | SearchPerformedEvent
  | ResultClickedEvent
  | ItemViewedEvent
  | ItemPurchasedEvent
  | SearchFilterAppliedEvent
  | SearchSortChangedEvent;

/**
 * Данные для отправки событий пакетом
 */
export interface BehaviorEventBatch {
  events: BehaviorEvent[];
  batch_id: string;
  created_at: string;
}

/**
 * Конфигурация для трекинга
 */
export interface BehaviorTrackingConfig {
  enabled: boolean;
  batchSize: number;
  batchTimeout: number;
  maxRetries: number;
  endpoint: string;
}

/**
 * Ответ от API при отправке событий
 */
export interface BehaviorTrackingResponse {
  success: boolean;
  processed_count: number;
  failed_count: number;
  errors?: string[];
}

/**
 * Контекст трекинга для хранения временных данных
 */
export interface TrackingContext {
  sessionId: string;
  currentSearchQuery?: string;
  currentSearchStartTime?: number;
  currentItemViewStartTime?: number;
  currentItemId?: string;
  lastSearchResults?: number;
  searchFilters?: Record<string, any>;
  searchSort?: string;
}

/**
 * Опции для методов трекинга
 */
export interface TrackingOptions {
  immediate?: boolean; // Отправить сразу, не дожидаясь батча
  retry?: boolean; // Повторять при ошибках
  context?: Record<string, any>; // Дополнительные данные
}

/**
 * Состояние hook'а useBehaviorTracking
 */
export interface BehaviorTrackingState {
  isEnabled: boolean;
  isInitialized: boolean;
  pendingEvents: BehaviorEvent[];
  context: TrackingContext;
  lastError?: string;
  retryCount: number;
}

/**
 * Методы hook'а useBehaviorTracking
 */
export interface BehaviorTrackingMethods {
  // Методы трекинга событий
  trackSearchPerformed: (
    data: Omit<SearchPerformedEvent, keyof BaseBehaviorEvent>
  ) => Promise<void>;
  trackResultClicked: (
    data: Omit<ResultClickedEvent, keyof BaseBehaviorEvent>
  ) => Promise<void>;
  trackItemViewed: (
    data: Omit<ItemViewedEvent, keyof BaseBehaviorEvent>
  ) => Promise<void>;
  trackItemPurchased: (
    data: Omit<ItemPurchasedEvent, keyof BaseBehaviorEvent>
  ) => Promise<void>;
  trackSearchFilterApplied: (
    data: Omit<SearchFilterAppliedEvent, keyof BaseBehaviorEvent>
  ) => Promise<void>;
  trackSearchSortChanged: (
    data: Omit<SearchSortChangedEvent, keyof BaseBehaviorEvent>
  ) => Promise<void>;

  // Методы управления
  flushEvents: () => Promise<void>;
  clearContext: () => void;
  updateContext: (updates: Partial<TrackingContext>) => void;

  // Методы для специфических use cases
  startSearch: (
    query: string,
    filters?: Record<string, any>,
    sort?: string
  ) => void;
  startItemView: (itemId: string) => void;
  endItemView: () => void;
}

/**
 * Возвращаемый тип hook'а useBehaviorTracking
 */
export interface UseBehaviorTrackingReturn extends BehaviorTrackingMethods {
  state: BehaviorTrackingState;
  config: BehaviorTrackingConfig;
}

/**
 * Опции для инициализации hook'а
 */
export interface UseBehaviorTrackingOptions {
  enabled?: boolean;
  batchSize?: number;
  batchTimeout?: number;
  maxRetries?: number;
  endpoint?: string;
  userId?: string;
  autoInit?: boolean;
}

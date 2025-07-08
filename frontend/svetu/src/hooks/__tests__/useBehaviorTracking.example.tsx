'use client';

/**
 * Пример использования hook'а useBehaviorTracking
 * Этот файл показывает, как правильно использовать hook
 */

import React from 'react';
import { useBehaviorTracking } from '../useBehaviorTracking';

interface ExampleComponentProps {
  userId?: string;
}

export const ExampleBehaviorTrackingComponent: React.FC<
  ExampleComponentProps
> = ({ userId }) => {
  const {
    state,
    trackSearchPerformed,
    trackResultClicked,
    trackItemViewed,
    trackItemPurchased,
    trackSearchFilterApplied,
    trackSearchSortChanged,
    flushEvents,
    startSearch,
    startItemView,
    endItemView,
  } = useBehaviorTracking({
    enabled: true,
    userId,
    batchSize: 5,
    batchTimeout: 3000,
  });

  const handleSearch = async () => {
    // Начинаем отслеживание поиска
    startSearch('iPhone 15', { category: 'electronics' }, 'price-asc');

    // Отправляем событие поиска
    await trackSearchPerformed({
      search_query: 'iPhone 15',
      search_filters: { category: 'electronics' },
      search_sort: 'price-asc',
      results_count: 25,
      search_duration_ms: 120,
    });
  };

  const handleResultClick = async () => {
    await trackResultClicked({
      search_query: 'iPhone 15',
      clicked_item_id: 'item-123',
      click_position: 3,
      total_results: 25,
      click_time_from_search_ms: 5000,
    });
  };

  const handleItemView = async () => {
    // Начинаем отслеживание просмотра товара
    startItemView('item-123');

    await trackItemViewed({
      item_id: 'item-123',
      item_title: 'iPhone 15 Pro Max',
      item_category: 'electronics',
      item_price: 1199.99,
      previous_page: 'search-results',
    });
  };

  const handleItemPurchase = async () => {
    await trackItemPurchased({
      item_id: 'item-123',
      item_title: 'iPhone 15 Pro Max',
      item_category: 'electronics',
      purchase_amount: 1199.99,
      currency: 'USD',
      payment_method: 'credit_card',
      time_from_first_view_ms: 300000, // 5 минут
    });
  };

  const handleFilterApply = async () => {
    await trackSearchFilterApplied({
      search_query: 'iPhone 15',
      filter_type: 'price_range',
      filter_value: '1000-1500',
      results_count_before: 25,
      results_count_after: 8,
    });
  };

  const handleSortChange = async () => {
    await trackSearchSortChanged({
      search_query: 'iPhone 15',
      sort_type: 'price-desc',
      previous_sort: 'price-asc',
      results_count: 8,
    });
  };

  const handleEndItemView = () => {
    endItemView();
  };

  const handleFlushEvents = async () => {
    await flushEvents();
  };

  return (
    <div className="p-6 space-y-4">
      <h2 className="text-2xl font-bold">Behavior Tracking Example</h2>

      <div className="bg-gray-100 p-4 rounded">
        <h3 className="font-semibold mb-2">Tracking State:</h3>
        <ul className="text-sm space-y-1">
          <li>Enabled: {state.isEnabled ? 'Yes' : 'No'}</li>
          <li>Initialized: {state.isInitialized ? 'Yes' : 'No'}</li>
          <li>Pending Events: {state.pendingEvents.length}</li>
          <li>Session ID: {state.context.sessionId}</li>
          <li>Retry Count: {state.retryCount}</li>
          {state.lastError && (
            <li className="text-red-600">Error: {state.lastError}</li>
          )}
        </ul>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <button onClick={handleSearch} className="btn btn-primary">
          Track Search
        </button>

        <button onClick={handleResultClick} className="btn btn-secondary">
          Track Result Click
        </button>

        <button onClick={handleItemView} className="btn btn-accent">
          Track Item View
        </button>

        <button onClick={handleItemPurchase} className="btn btn-success">
          Track Purchase
        </button>

        <button onClick={handleFilterApply} className="btn btn-warning">
          Track Filter Apply
        </button>

        <button onClick={handleSortChange} className="btn btn-info">
          Track Sort Change
        </button>

        <button onClick={handleEndItemView} className="btn btn-neutral">
          End Item View
        </button>

        <button onClick={handleFlushEvents} className="btn btn-ghost">
          Flush Events
        </button>
      </div>

      <div className="bg-blue-50 p-4 rounded">
        <h3 className="font-semibold mb-2">Usage Notes:</h3>
        <ul className="text-sm space-y-1 list-disc list-inside">
          <li>
            Events are automatically batched and sent every 5 seconds or when 10
            events accumulate
          </li>
          <li>Purchase events are sent immediately for better reliability</li>
          <li>
            Session ID is automatically generated and stored in localStorage
          </li>
          <li>Device and browser information is automatically detected</li>
          <li>Retry logic handles temporary network failures</li>
          <li>Events are sent using navigator.sendBeacon on page unload</li>
        </ul>
      </div>
    </div>
  );
};

export default ExampleBehaviorTrackingComponent;

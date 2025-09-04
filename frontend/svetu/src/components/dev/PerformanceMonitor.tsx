'use client';

import { useEffect, useState, useCallback } from 'react';

interface PerformanceMetrics {
  renderTime: number;
  componentCount: number;
  memoryUsage: number;
  localStorageSize: number;
  cacheHitRate: number;
}

interface PerformanceEntry extends Performance {
  memory?: {
    usedJSHeapSize: number;
    totalJSHeapSize: number;
    jsHeapSizeLimit: number;
  };
}

export function PerformanceMonitor({ enabled = false }: { enabled?: boolean }) {
  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    renderTime: 0,
    componentCount: 0,
    memoryUsage: 0,
    localStorageSize: 0,
    cacheHitRate: 0,
  });

  const [isVisible, setIsVisible] = useState(false);

  // Подсчет размера localStorage
  const calculateLocalStorageSize = useCallback(() => {
    try {
      let total = 0;
      for (const key in localStorage) {
        if (localStorage.hasOwnProperty(key)) {
          total += localStorage[key].length + key.length;
        }
      }
      return total;
    } catch {
      return 0;
    }
  }, []);

  // Подсчет cache hit rate для атрибутов
  const calculateCacheHitRate = useCallback(() => {
    try {
      const cacheStats = localStorage.getItem('attribute_cache_stats');
      if (cacheStats) {
        const stats = JSON.parse(cacheStats);
        return (stats.hits / (stats.hits + stats.misses)) * 100;
      }
      return 0;
    } catch {
      return 0;
    }
  }, []);

  // Обновление метрик
  const updateMetrics = useCallback(() => {
    if (!enabled) return;

    const perf = performance as PerformanceEntry;

    // Подсчет времени рендеринга (примерный)
    const navigationEntries = performance.getEntriesByType('navigation');
    const renderTime =
      navigationEntries.length > 0
        ? (navigationEntries[0] as PerformanceNavigationTiming).loadEventEnd -
          (navigationEntries[0] as PerformanceNavigationTiming).responseEnd
        : 0;

    // Память
    const memoryUsage = perf.memory
      ? Math.round(perf.memory.usedJSHeapSize / 1048576)
      : 0; // В МБ

    // LocalStorage размер
    const localStorageSize = Math.round(calculateLocalStorageSize() / 1024); // В КБ

    // Cache hit rate
    const cacheHitRate = calculateCacheHitRate();

    // Подсчет компонентов в DOM (приблизительно)
    const componentCount = document.querySelectorAll(
      '[data-testid], [class*="Component"]'
    ).length;

    setMetrics({
      renderTime: Math.round(renderTime),
      componentCount,
      memoryUsage,
      localStorageSize,
      cacheHitRate: Math.round(cacheHitRate),
    });
  }, [enabled, calculateLocalStorageSize, calculateCacheHitRate]);

  // Автообновление метрик
  useEffect(() => {
    if (!enabled) return;

    const interval = setInterval(updateMetrics, 2000); // Каждые 2 секунды
    updateMetrics(); // Первоначальное обновление

    return () => clearInterval(interval);
  }, [enabled, updateMetrics]);

  // Горячие клавиши для показа/скрытия
  useEffect(() => {
    if (!enabled) return;

    const handleKeyPress = (event: KeyboardEvent) => {
      // Ctrl+Shift+P для toggle
      if (event.ctrlKey && event.shiftKey && event.key === 'P') {
        event.preventDefault();
        setIsVisible((prev) => !prev);
      }
    };

    document.addEventListener('keydown', handleKeyPress);
    return () => document.removeEventListener('keydown', handleKeyPress);
  }, [enabled]);

  // Очистка кэша атрибутов
  const clearAttributeCache = useCallback(() => {
    try {
      const keysToRemove: string[] = [];
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (
          key &&
          (key.startsWith('recent_') ||
            key.startsWith('popular_') ||
            key.startsWith('count_'))
        ) {
          keysToRemove.push(key);
        }
      }
      keysToRemove.forEach((key) => localStorage.removeItem(key));
      updateMetrics();
      alert(`Очищено ${keysToRemove.length} ключей кэша атрибутов`);
    } catch (error) {
      console.error('Error clearing attribute cache:', error);
    }
  }, [updateMetrics]);

  if (!enabled || !isVisible) {
    return enabled ? (
      <div className="fixed bottom-4 right-4 bg-neutral text-neutral-content px-2 py-1 rounded text-xs opacity-50 hover:opacity-100 transition-opacity">
        Ctrl+Shift+P для метрик
      </div>
    ) : null;
  }

  return (
    <div className="fixed top-4 right-4 bg-base-100 border border-base-300 rounded-lg p-4 w-72 shadow-lg z-50">
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-bold text-sm">⚡ Performance Monitor</h3>
        <button
          onClick={() => setIsVisible(false)}
          className="btn btn-xs btn-ghost"
        >
          ✕
        </button>
      </div>

      <div className="space-y-2 text-xs">
        {/* Время рендеринга */}
        <div className="flex justify-between">
          <span>Render Time:</span>
          <span
            className={metrics.renderTime > 500 ? 'text-error' : 'text-success'}
          >
            {metrics.renderTime}ms
          </span>
        </div>

        {/* Память */}
        <div className="flex justify-between">
          <span>Memory Usage:</span>
          <span
            className={
              metrics.memoryUsage > 100 ? 'text-warning' : 'text-success'
            }
          >
            {metrics.memoryUsage}MB
          </span>
        </div>

        {/* localStorage размер */}
        <div className="flex justify-between">
          <span>LocalStorage:</span>
          <span
            className={
              metrics.localStorageSize > 1024 ? 'text-warning' : 'text-success'
            }
          >
            {metrics.localStorageSize}KB
          </span>
        </div>

        {/* Cache hit rate */}
        <div className="flex justify-between">
          <span>Cache Hit Rate:</span>
          <span
            className={
              metrics.cacheHitRate < 50 ? 'text-error' : 'text-success'
            }
          >
            {metrics.cacheHitRate}%
          </span>
        </div>

        {/* Количество компонентов */}
        <div className="flex justify-between">
          <span>Components:</span>
          <span
            className={
              metrics.componentCount > 200 ? 'text-warning' : 'text-success'
            }
          >
            {metrics.componentCount}
          </span>
        </div>
      </div>

      {/* Действия */}
      <div className="mt-4 space-y-2">
        <button
          onClick={updateMetrics}
          className="btn btn-xs btn-outline w-full"
        >
          🔄 Обновить метрики
        </button>

        <button
          onClick={clearAttributeCache}
          className="btn btn-xs btn-outline btn-warning w-full"
        >
          🗑️ Очистить кэш атрибутов
        </button>
      </div>

      {/* Советы */}
      <div className="mt-3 p-2 bg-base-200 rounded text-xs">
        <div className="font-bold mb-1">💡 Советы:</div>
        <ul className="space-y-1 text-xs opacity-75">
          <li>• Render Time &lt; 200ms - хорошо</li>
          <li>• Memory &lt; 50MB - отлично</li>
          <li>• Cache Hit Rate &gt; 80% - цель</li>
          <li>• LocalStorage &lt; 500KB - норма</li>
        </ul>
      </div>

      <div className="mt-2 text-xs opacity-50 text-center">
        Обновляется каждые 2 сек
      </div>
    </div>
  );
}

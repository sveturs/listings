import { useState, useCallback, useRef } from 'react';
import { apiClient } from '@/services/api-client';

export type LoadingStage = 'initial' | 'basic' | 'detailed' | 'complete';

interface ProgressiveLoadingOptions {
  onStageComplete?: (stage: LoadingStage, data: any) => void;
  onError?: (error: Error, stage: LoadingStage) => void;
}

/**
 * Хук для прогрессивной загрузки данных карты
 * Загружает данные в несколько этапов для улучшения UX
 */
export function useProgressiveLoading(options: ProgressiveLoadingOptions = {}) {
  const [loadingStage, setLoadingStage] = useState<LoadingStage>('initial');
  const [isLoading, setIsLoading] = useState(false);
  const abortControllerRef = useRef<AbortController | null>(null);

  /**
   * Отменить текущую загрузку
   */
  const cancelLoading = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
    setIsLoading(false);
    setLoadingStage('initial');
  }, []);

  /**
   * Прогрессивная загрузка данных
   * @param bounds - Границы видимой области карты
   * @param filters - Фильтры поиска
   * @param zoom - Уровень зума карты
   */
  const loadProgressively = useCallback(
    async (
      bounds: { north: number; south: number; east: number; west: number },
      filters: any = {},
      zoom: number = 10
    ) => {
      // Отменяем предыдущую загрузку если она есть
      cancelLoading();

      // Создаем новый контроллер для отмены
      const controller = new AbortController();
      abortControllerRef.current = controller;

      setIsLoading(true);

      try {
        // Этап 1: Загрузка кластеров для текущего зума
        setLoadingStage('basic');

        if (zoom < 12) {
          // На малых зумах загружаем только кластеры
          const params = new URLSearchParams({
            bounds: `${bounds.south},${bounds.west},${bounds.north},${bounds.east}`,
            zoom: zoom.toString(),
          });
          const clustersResponse = await apiClient.get(
            `/api/v1/gis/clusters?${params}`
          );

          if (!controller.signal.aborted) {
            options.onStageComplete?.('basic', {
              clusters: clustersResponse.data.data,
              type: 'clusters',
            });
          }
        }

        // Этап 2: Загрузка основных данных в видимой области
        setLoadingStage('detailed');

        const detailParams = new URLSearchParams({
          latitude: ((bounds.north + bounds.south) / 2).toString(),
          longitude: ((bounds.east + bounds.west) / 2).toString(),
          radius: calculateRadiusFromBounds(bounds).toString(),
          limit: (zoom > 15 ? 500 : 200).toString(),
          ...Object.fromEntries(
            Object.entries(filters).map(([k, v]) => [k, String(v)])
          ),
        });
        const mainDataResponse = await apiClient.get(
          `/api/v1/gis/search/radius?${detailParams}`
        );

        if (!controller.signal.aborted) {
          options.onStageComplete?.('detailed', {
            listings: mainDataResponse.data.data?.listings || [],
            totalCount: mainDataResponse.data.data?.totalCount || 0,
            type: 'listings',
          });
        }

        // Этап 3: Догрузка оставшихся данных если нужно
        if (mainDataResponse.data.data?.hasMore && zoom > 14) {
          setLoadingStage('complete');

          const completeParams = new URLSearchParams({
            latitude: ((bounds.north + bounds.south) / 2).toString(),
            longitude: ((bounds.east + bounds.west) / 2).toString(),
            radius: calculateRadiusFromBounds(bounds).toString(),
            offset: '200',
            limit: '800',
            ...Object.fromEntries(
              Object.entries(filters).map(([k, v]) => [k, String(v)])
            ),
          });
          const remainingDataResponse = await apiClient.get(
            `/api/v1/gis/search/radius?${completeParams}`
          );

          if (!controller.signal.aborted) {
            options.onStageComplete?.('complete', {
              listings: remainingDataResponse.data.data?.listings || [],
              totalCount: remainingDataResponse.data.data?.totalCount || 0,
              type: 'listings_additional',
            });
          }
        }

        setLoadingStage('complete');
      } catch (error: any) {
        if (error.name !== 'AbortError') {
          console.error('Progressive loading error:', error);
          options.onError?.(error, loadingStage);
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    },
    [cancelLoading, loadingStage, options]
  );

  return {
    loadProgressively,
    loadingStage,
    isLoading,
    cancelLoading,
  };
}

/**
 * Вычислить радиус поиска из границ карты
 * @param bounds - Границы видимой области
 * @returns Радиус в метрах
 */
function calculateRadiusFromBounds(bounds: {
  north: number;
  south: number;
  east: number;
  west: number;
}): number {
  // Вычисляем диагональ видимой области
  const latDiff = bounds.north - bounds.south;
  const lngDiff = bounds.east - bounds.west;

  // Приблизительный расчет в метрах (1 градус ≈ 111км)
  const latDistance = latDiff * 111000;
  const lngDistance =
    lngDiff *
    111000 *
    Math.cos((((bounds.north + bounds.south) / 2) * Math.PI) / 180);

  // Радиус = половина диагонали
  const diagonal = Math.sqrt(latDistance ** 2 + lngDistance ** 2);
  return Math.round(diagonal / 2);
}

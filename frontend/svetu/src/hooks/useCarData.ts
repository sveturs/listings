import React from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { CarsService } from '@/services/CarsService';
import type { CarMake, CarModel, CarGeneration } from '@/types/cars';

// Время кэширования - 1 час
const CACHE_TIME = 60 * 60 * 1000;
const STALE_TIME = 30 * 60 * 1000; // 30 минут

/**
 * Hook для получения списка марок автомобилей с кэшированием
 */
export const useCarMakes = () => {
  return useQuery<CarMake[], Error>({
    queryKey: ['car-makes'],
    queryFn: async () => {
      const response = await CarsService.getMakes();
      if (response.success && response.data) {
        return response.data;
      }
      throw new Error('Failed to load car makes');
    },
    staleTime: STALE_TIME,
    gcTime: CACHE_TIME,
    refetchOnWindowFocus: false,
  });
};

/**
 * Hook для получения моделей конкретной марки с кэшированием
 */
export const useCarModels = (makeSlug: string, enabled = true) => {
  return useQuery<CarModel[], Error>({
    queryKey: ['car-models', makeSlug],
    queryFn: async () => {
      if (!makeSlug) return [];
      const response = await CarsService.getModels(makeSlug);
      if (response.success && response.data) {
        return response.data;
      }
      throw new Error('Failed to load car models');
    },
    enabled: enabled && !!makeSlug,
    staleTime: STALE_TIME,
    gcTime: CACHE_TIME,
    refetchOnWindowFocus: false,
  });
};

/**
 * Hook для получения поколений модели с кэшированием
 */
export const useCarGenerations = (modelId: number | null, enabled = true) => {
  return useQuery<CarGeneration[], Error>({
    queryKey: ['car-generations', modelId],
    queryFn: async () => {
      if (!modelId) return [];
      const response = await CarsService.getGenerations(modelId);
      if (response.success && response.data) {
        return response.data;
      }
      throw new Error('Failed to load car generations');
    },
    enabled: enabled && !!modelId,
    staleTime: STALE_TIME,
    gcTime: CACHE_TIME,
    refetchOnWindowFocus: false,
  });
};

/**
 * Hook для предзагрузки популярных марок
 */
export const usePrefetchPopularMakes = () => {
  const queryClient = useQueryClient();

  const popularMakeSlugs = [
    'volkswagen',
    'mercedes-benz',
    'bmw',
    'audi',
    'toyota',
    'ford',
    'opel',
    'renault',
    'peugeot',
    'fiat',
  ];

  // Предзагрузка моделей для популярных марок
  React.useEffect(() => {
    popularMakeSlugs.forEach((makeSlug) => {
      void queryClient.prefetchQuery({
        queryKey: ['car-models', makeSlug],
        queryFn: async () => {
          const response = await CarsService.getModels(makeSlug);
          return response.data || [];
        },
        staleTime: STALE_TIME,
        gcTime: CACHE_TIME,
      });
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryClient]);
};

/**
 * Hook для инвалидации кэша автомобильных данных
 */
export const useInvalidateCarCache = () => {
  const queryClient = useQueryClient();

  return {
    invalidateMakes: () =>
      queryClient.invalidateQueries({ queryKey: ['car-makes'] }),
    invalidateModels: (makeSlug?: string) => {
      if (makeSlug) {
        queryClient.invalidateQueries({ queryKey: ['car-models', makeSlug] });
      } else {
        queryClient.invalidateQueries({ queryKey: ['car-models'] });
      }
    },
    invalidateGenerations: (modelSlug?: string) => {
      if (modelSlug) {
        queryClient.invalidateQueries({
          queryKey: ['car-generations', modelSlug],
        });
      } else {
        queryClient.invalidateQueries({ queryKey: ['car-generations'] });
      }
    },
    invalidateAll: () => {
      queryClient.invalidateQueries({ queryKey: ['car-makes'] });
      queryClient.invalidateQueries({ queryKey: ['car-models'] });
      queryClient.invalidateQueries({ queryKey: ['car-generations'] });
    },
  };
};

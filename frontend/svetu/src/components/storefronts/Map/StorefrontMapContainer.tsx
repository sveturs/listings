'use client';

import { useState, useCallback, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import dynamic from 'next/dynamic';
import { useStorefronts } from '@/hooks/useStorefronts';
import type {
  StorefrontMapData,
  StorefrontFilters,
  GeocodeResult,
  StorefrontMapConfig,
} from '@/types/storefront';
import { storefrontApi } from '@/services/storefrontApi';
import MapFilters from './MapFilters';
import AddressSearch from './AddressSearch';

// Динамически загружаем карту для избежания SSR проблем
const StorefrontMap = dynamic(() => import('./StorefrontMap'), {
  ssr: false,
  loading: () => (
    <div className="bg-base-200 animate-pulse rounded-lg h-96">
      <div className="flex items-center justify-center h-full">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    </div>
  ),
});

interface StorefrontMapContainerProps {
  initialCenter?: { lat: number; lng: number };
  initialZoom?: number;
  height?: string;
  showFilters?: boolean;
  showSearch?: boolean;
  clustering?: boolean;
  mapConfig?: Partial<StorefrontMapConfig>;
  className?: string;
}

const StorefrontMapContainer: React.FC<StorefrontMapContainerProps> = ({
  initialCenter = { lat: 44.7866, lng: 20.4489 }, // Белград
  initialZoom = 12,
  height = '500px',
  showFilters = true,
  showSearch = true,
  clustering = false,
  // mapConfig,
  className = '',
}) => {
  const router = useRouter();
  const locale = useLocale();
  const { filters, updateFilters, resetFilters, isLoading } = useStorefronts();

  const [mapStorefronts, setMapStorefronts] = useState<StorefrontMapData[]>([]);
  const [mapCenter, setMapCenter] = useState(initialCenter);
  const [mapZoom, setMapZoom] = useState(initialZoom);
  const [isLoadingMap, setIsLoadingMap] = useState(false);
  const [filtersVisible, setFiltersVisible] = useState(false);

  // Загрузка витрин для карты
  const loadMapStorefronts = useCallback(
    async (bounds?: L.LatLngBounds) => {
      setIsLoadingMap(true);

      try {
        if (bounds) {
          // Загрузка витрин в пределах карты
          const mapData = await storefrontApi.getStorefrontsForMap({
            north: bounds.getNorth(),
            south: bounds.getSouth(),
            east: bounds.getEast(),
            west: bounds.getWest(),
          });
          setMapStorefronts(mapData);
        } else {
          // Загрузка всех витрин или по фильтрам
          const response = await storefrontApi.getStorefronts({
            ...filters,
            limit: 1000, // Лимит для карты
          });

          // Преобразуем в формат для карты
          const mapData: StorefrontMapData[] = (response.storefronts || []).map(
            (storefront) => ({
              id: storefront.id,
              name: storefront.name,
              latitude: storefront.latitude,
              longitude: storefront.longitude,
              address: storefront.address,
              phone: storefront.phone,
              logo_url: storefront.logo_url,
              rating: storefront.rating,
              accepts_cards: true, // Пока заглушка
              working_now: true, // Пока заглушка
              has_delivery: true, // Пока заглушка
              has_self_pickup: true, // Пока заглушка
              supports_cod: true, // Пока заглушка
              products_count: storefront.products_count,
              slug: storefront.slug,
            })
          );

          setMapStorefronts(mapData);
        }
      } catch (error) {
        console.error('Ошибка загрузки витрин для карты:', error);
        setMapStorefronts([]);
      } finally {
        setIsLoadingMap(false);
      }
    },
    [filters]
  );

  // Загрузка витрин при изменении фильтров
  useEffect(() => {
    loadMapStorefronts();
  }, [loadMapStorefronts]);

  const handleFiltersChange = useCallback(
    (newFilters: Partial<StorefrontFilters>) => {
      updateFilters(newFilters);
    },
    [updateFilters]
  );

  const handleResetFilters = useCallback(() => {
    resetFilters();
  }, [resetFilters]);

  const handleLocationSelect = useCallback(
    (location: GeocodeResult) => {
      setMapCenter({ lat: location.latitude, lng: location.longitude });
      setMapZoom(14);

      // Обновляем фильтры с новым местоположением
      handleFiltersChange({
        latitude: location.latitude,
        longitude: location.longitude,
        radiusKm: 10,
        city: location.city,
      });
    },
    [handleFiltersChange]
  );

  const handleStorefrontClick = useCallback(
    (storefront: StorefrontMapData) => {
      // Переход на страницу витрины
      router.push(`/${locale}/storefronts/${storefront.slug}`);
    },
    [router, locale]
  );

  const handleBoundsChange = useCallback((_bounds: L.LatLngBounds) => {
    // Опционально можно обновлять витрины при изменении границ карты
    // loadMapStorefronts(bounds);
  }, []);

  return (
    <div className={`relative ${className}`}>
      {/* Заголовок и поиск */}
      <div className="mb-4">
        <div className="flex flex-col lg:flex-row gap-4 items-start lg:items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold mb-2">Карта витрин</h2>
            <p className="text-gray-600">
              Найдено витрин: {mapStorefronts.length}
            </p>
          </div>

          {showSearch && (
            <div className="w-full lg:w-96">
              <AddressSearch
                onLocationSelect={handleLocationSelect}
                placeholder="Поиск по адресу..."
                className="w-full"
              />
            </div>
          )}
        </div>

        {/* Кнопки управления */}
        <div className="flex items-center gap-2 mt-4">
          {showFilters && (
            <button
              className={`btn btn-outline btn-sm ${filtersVisible ? 'btn-active' : ''}`}
              onClick={() => setFiltersVisible(!filtersVisible)}
            >
              🔍 Фильтры
              {Object.values(filters).filter(
                (v) => v !== undefined && v !== null && v !== ''
              ).length > 0 && (
                <span className="badge badge-primary ml-1">
                  {
                    Object.values(filters).filter(
                      (v) => v !== undefined && v !== null && v !== ''
                    ).length
                  }
                </span>
              )}
            </button>
          )}

          <button
            className="btn btn-outline btn-sm"
            onClick={() => loadMapStorefronts()}
            disabled={isLoadingMap}
          >
            {isLoadingMap ? (
              <span className="loading loading-spinner loading-sm"></span>
            ) : (
              '🔄'
            )}
            Обновить
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
        {/* Боковая панель с фильтрами */}
        {showFilters && filtersVisible && (
          <div className="lg:col-span-3">
            <MapFilters
              filters={filters}
              onFiltersChange={handleFiltersChange}
              onResetFilters={handleResetFilters}
              isLoading={isLoading}
            />
          </div>
        )}

        {/* Карта */}
        <div
          className={`${showFilters && filtersVisible ? 'lg:col-span-9' : 'lg:col-span-12'}`}
        >
          <div className="relative">
            {isLoadingMap && (
              <div className="absolute top-4 left-1/2 transform -translate-x-1/2 z-[1000]">
                <div className="bg-white rounded-lg shadow-lg px-4 py-2 flex items-center gap-2">
                  <span className="loading loading-spinner loading-sm"></span>
                  <span>Загрузка витрин...</span>
                </div>
              </div>
            )}

            <StorefrontMap
              storefronts={mapStorefronts}
              center={mapCenter}
              zoom={mapZoom}
              height={height}
              onStorefrontClick={handleStorefrontClick}
              onBoundsChange={handleBoundsChange}
              showSearch={false} // Поиск выносим в заголовок
              clustering={clustering}
              className="rounded-lg shadow-lg"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default StorefrontMapContainer;

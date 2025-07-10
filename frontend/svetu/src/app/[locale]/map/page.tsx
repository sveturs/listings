'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { InteractiveMap } from '@/components/GIS';
import { useGeoSearch } from '@/components/GIS/hooks/useGeoSearch';
import { useApi } from '@/hooks/useApi';
import { MapViewState, MapMarkerData } from '@/components/GIS/types/gis';
import { useDebounce } from '@/hooks/useDebounce';
import { SearchBar } from '@/components/SearchBar';
import { useRouter } from '@/i18n/routing';
import { toast } from 'react-hot-toast';

interface ListingData {
  id: number;
  name: string;
  price: number;
  location: {
    lat: number;
    lng: number;
    city?: string;
    country?: string;
  };
  category: {
    id: number;
    name: string;
    slug: string;
  };
  images: string[];
  created_at: string;
}

interface MapFilters {
  category: string;
  priceFrom: number;
  priceTo: number;
  radius: number;
}

const MapPage: React.FC = () => {
  const t = useTranslations('map');
  const router = useRouter();
  const { get } = useApi();
  const { search: geoSearch } = useGeoSearch();

  // Состояние карты
  const [viewState, setViewState] = useState<MapViewState>({
    longitude: 20.4649, // Сербия - Белград
    latitude: 44.8176,
    zoom: 10,
    pitch: 0,
    bearing: 0,
  });

  // Данные и фильтры
  const [listings, setListings] = useState<ListingData[]>([]);
  const [markers, setMarkers] = useState<MapMarkerData[]>([]);
  const [filters, setFilters] = useState<MapFilters>({
    category: '',
    priceFrom: 0,
    priceTo: 0,
    radius: 10000, // 10 км в метрах
  });

  // Поиск
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedSearchQuery = useDebounce(searchQuery, 500);

  // Состояние загрузки
  const [isLoading, setIsLoading] = useState(false);
  const [isSearching, setIsSearching] = useState(false);

  // Загрузка объявлений для карты
  const loadListings = useCallback(async () => {
    setIsLoading(true);
    try {
      const params = new URLSearchParams({
        limit: '100',
        page: '1',
        sort_by: 'date',
        sort_order: 'desc',
        ...(filters.category && { category_id: filters.category }),
        ...(filters.priceFrom > 0 && {
          price_min: filters.priceFrom.toString(),
        }),
        ...(filters.priceTo > 0 && { price_max: filters.priceTo.toString() }),
      });

      const response = await get(`/api/v1/search?${params}`);

      if (response.items) {
        // Преобразуем данные API в формат, ожидаемый компонентом
        const transformedListings = response.items.map((item: any) => ({
          id: item.product_id,
          name: item.name,
          price: item.price,
          location: {
            lat: item.location.lat,
            lng: item.location.lng,
            city: item.location.city,
            country: item.location.country,
          },
          category: item.category,
          images: item.images || [],
          created_at: item.created_at,
        }));
        setListings(transformedListings);
      }
    } catch (error) {
      console.error('Error loading listings:', error);
      toast.error(t('errors.loadingFailed'));
    } finally {
      setIsLoading(false);
    }
  }, [filters, get, t]);

  // Преобразование объявлений в маркеры
  const createMarkers = useCallback(
    (listingsData: ListingData[]): MapMarkerData[] => {
      return listingsData
        .filter((listing) => listing.location?.lat && listing.location?.lng)
        .map((listing) => ({
          id: listing.id.toString(),
          position: [listing.location.lng, listing.location.lat] as [
            number,
            number,
          ],
          type: 'listing',
          data: {
            title: listing.name,
            price: listing.price,
            category: listing.category?.name || 'Unknown',
            image: listing.images?.[0],
            address:
              `${listing.location.city || ''}, ${listing.location.country || ''}`
                .trim()
                .replace(/^,\s*|,\s*$/, ''),
            id: listing.id,
          },
          color: getCategoryColor(listing.category?.slug || 'default'),
          size: 'medium',
        }));
    },
    []
  );

  // Получение цвета для категории
  const getCategoryColor = (categorySlug: string): string => {
    const colors: { [key: string]: string } = {
      'real-estate': '#3B82F6', // blue
      vehicles: '#EF4444', // red
      electronics: '#10B981', // green
      clothing: '#F59E0B', // amber
      services: '#8B5CF6', // violet
      jobs: '#F97316', // orange
      'children-goods-toys': '#EC4899', // pink
      'home-garden': '#16A34A', // green
      appliances: '#0EA5E9', // sky
      default: '#6B7280', // gray
    };
    return colors[categorySlug] || colors.default;
  };

  // Поиск по адресу
  const handleAddressSearch = useCallback(
    async (query: string) => {
      if (!query.trim()) return;

      setIsSearching(true);
      try {
        const results = await geoSearch({
          query,
          limit: 1,
          language: 'ru',
        });

        if (results.length > 0) {
          const result = results[0];
          const newViewState = {
            ...viewState,
            longitude: parseFloat(result.lon),
            latitude: parseFloat(result.lat),
            zoom: 14,
          };
          setViewState(newViewState);
          toast.success(t('search.found'));
        } else {
          toast.error(t('search.notFound'));
        }
      } catch (error) {
        console.error('Search error:', error);
        toast.error(t('search.error'));
      } finally {
        setIsSearching(false);
      }
    },
    [geoSearch, viewState, t]
  );

  // Обработка поиска
  useEffect(() => {
    if (debouncedSearchQuery) {
      handleAddressSearch(debouncedSearchQuery);
    }
  }, [debouncedSearchQuery, handleAddressSearch]);

  // Обработка изменений фильтров
  useEffect(() => {
    loadListings();
  }, [loadListings]);

  // Создание маркеров при изменении объявлений
  useEffect(() => {
    const newMarkers = createMarkers(listings);
    setMarkers(newMarkers);
  }, [listings, createMarkers]);

  // Обработка клика по маркеру
  const handleMarkerClick = useCallback(
    (marker: MapMarkerData) => {
      if (marker.data?.id) {
        router.push(`/marketplace/${marker.data.id}`);
      }
    },
    [router]
  );

  // Обработка изменения области просмотра
  const handleViewStateChange = useCallback((newViewState: MapViewState) => {
    setViewState(newViewState);
  }, []);

  // Обработка изменения фильтров
  const handleFiltersChange = useCallback((newFilters: Partial<MapFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
  }, []);

  return (
    <div className="min-h-screen bg-base-100">
      {/* Заголовок */}
      <div className="bg-white border-b border-base-300 px-4 py-4">
        <div className="container mx-auto">
          <h1 className="text-2xl font-bold text-base-content mb-2">
            {t('title')}
          </h1>
          <p className="text-base-content-secondary">{t('description')}</p>
        </div>
      </div>

      {/* Контейнер с картой и фильтрами */}
      <div className="relative h-screen">
        {/* Боковая панель с фильтрами */}
        <div className="absolute left-4 top-4 z-10 w-80 bg-white rounded-lg shadow-lg">
          {/* Поиск по адресу */}
          <div className="p-4 border-b border-base-300">
            <label className="block text-sm font-medium text-base-content mb-2">
              {t('search.address')}
            </label>
            <SearchBar
              value={searchQuery}
              onChange={(value) => setSearchQuery(value)}
              placeholder={t('search.addressPlaceholder')}
              isLoading={isSearching}
              className="w-full"
            />
          </div>

          {/* Фильтры */}
          <div className="p-4">
            <h3 className="text-lg font-medium text-base-content mb-3">
              {t('filters.title')}
            </h3>

            {/* Категория */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-base-content mb-1">
                {t('filters.category')}
              </label>
              <select
                className="select select-bordered w-full"
                value={filters.category}
                onChange={(e) =>
                  handleFiltersChange({ category: e.target.value })
                }
              >
                <option value="">{t('filters.allCategories')}</option>
                <option value="real-estate">
                  {t('categories.realEstate')}
                </option>
                <option value="vehicles">{t('categories.vehicles')}</option>
                <option value="electronics">
                  {t('categories.electronics')}
                </option>
                <option value="clothing">{t('categories.clothing')}</option>
                <option value="services">{t('categories.services')}</option>
                <option value="jobs">{t('categories.jobs')}</option>
              </select>
            </div>

            {/* Цена от */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-base-content mb-1">
                {t('filters.priceFrom')}
              </label>
              <input
                type="number"
                className="input input-bordered w-full"
                value={filters.priceFrom || ''}
                onChange={(e) =>
                  handleFiltersChange({
                    priceFrom: parseInt(e.target.value) || 0,
                  })
                }
                placeholder="0"
              />
            </div>

            {/* Цена до */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-base-content mb-1">
                {t('filters.priceTo')}
              </label>
              <input
                type="number"
                className="input input-bordered w-full"
                value={filters.priceTo || ''}
                onChange={(e) =>
                  handleFiltersChange({
                    priceTo: parseInt(e.target.value) || 0,
                  })
                }
                placeholder="∞"
              />
            </div>

            {/* Радиус поиска */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-base-content mb-1">
                {t('filters.radius')}: {Math.round(filters.radius / 1000)} км
              </label>
              <input
                type="range"
                className="range range-primary"
                min="1000"
                max="50000"
                step="1000"
                value={filters.radius}
                onChange={(e) =>
                  handleFiltersChange({ radius: parseInt(e.target.value) })
                }
              />
            </div>

            {/* Статистика */}
            <div className="text-sm text-base-content-secondary">
              {t('results.showing')}: {markers.length} {t('results.listings')}
            </div>
          </div>
        </div>

        {/* Карта */}
        <div className="absolute inset-0">
          <InteractiveMap
            initialViewState={viewState}
            markers={markers}
            onMarkerClick={handleMarkerClick}
            onViewStateChange={handleViewStateChange}
            className="w-full h-full"
            controlsConfig={{
              showNavigation: true,
              showFullscreen: true,
              showGeolocate: true,
              showScale: true,
              showSearch: false, // Отключаем встроенный поиск, используем свой
            }}
          />
        </div>

        {/* Индикатор загрузки */}
        {isLoading && (
          <div className="absolute top-4 right-4 z-10 bg-white rounded-lg shadow-lg p-3">
            <div className="flex items-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
              <span className="text-sm text-base-content">{t('loading')}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default MapPage;

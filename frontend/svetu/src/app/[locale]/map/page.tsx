'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { InteractiveMap } from '@/components/GIS';
import { useGeoSearch } from '@/components/GIS/hooks/useGeoSearch';
import { MapViewState, MapMarkerData } from '@/components/GIS/types/gis';
import { useDebounce } from '@/hooks/useDebounce';
import { SearchBar } from '@/components/SearchBar';
import { useRouter } from '@/i18n/routing';
import { useSearchParams } from 'next/navigation';
import { toast } from 'react-hot-toast';
import { apiClient } from '@/services/api-client';
import { MobileFiltersDrawer } from '@/components/GIS/Mobile';
import WalkingAccessibilityControl from '@/components/GIS/Map/WalkingAccessibilityControl';
import { isPointInIsochrone } from '@/components/GIS/utils/mapboxIsochrone';
import type { Feature, Polygon } from 'geojson';

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
  const searchParams = useSearchParams();
  const { search: geoSearch } = useGeoSearch();

  // Функция для получения начальных значений из URL
  const getInitialFiltersFromURL = (): MapFilters => {
    return {
      category: searchParams.get('category') || '',
      priceFrom: parseInt(searchParams.get('priceFrom') || '0') || 0,
      priceTo: parseInt(searchParams.get('priceTo') || '0') || 0,
      radius: parseInt(searchParams.get('radius') || '10000') || 10000,
    };
  };

  // Функция для получения начального состояния карты из URL
  const getInitialViewStateFromURL = (): MapViewState => {
    const lat = parseFloat(searchParams.get('lat') || '44.8176');
    const lng = parseFloat(searchParams.get('lng') || '20.4649');
    const zoom = parseFloat(searchParams.get('zoom') || '10');

    return {
      longitude: lng,
      latitude: lat,
      zoom: zoom,
      pitch: 0,
      bearing: 0,
    };
  };

  // Состояние карты
  const [viewState, setViewState] = useState<MapViewState>(
    getInitialViewStateFromURL()
  );
  const [isInitialized, setIsInitialized] = useState(false);

  // Состояние маркера покупателя
  const [buyerLocation, setBuyerLocation] = useState({
    longitude: viewState.longitude,
    latitude: viewState.latitude,
  });

  // Данные и фильтры
  const [listings, setListings] = useState<ListingData[]>([]);
  const [markers, setMarkers] = useState<MapMarkerData[]>([]);
  const [filters, setFilters] = useState<MapFilters>(
    getInitialFiltersFromURL()
  );

  // Поиск
  const [searchQuery, setSearchQuery] = useState(searchParams.get('q') || '');
  const debouncedSearchQuery = useDebounce(searchQuery, 500);

  // Создаем debounced версию фильтров для оптимизации запросов
  const debouncedFilters = useDebounce(filters, 800);

  // Создаем debounced версию viewState для оптимизации обновления URL
  const debouncedViewState = useDebounce(viewState, 500);

  // Состояние загрузки
  const [isLoading, setIsLoading] = useState(false);
  const [isSearching, setIsSearching] = useState(false);

  // Состояние для WalkingAccessibilityControl
  const [walkingMode, setWalkingMode] = useState<'radius' | 'walking'>(
    'radius'
  );
  const [walkingTime, setWalkingTime] = useState(15);

  // Состояние мобильных элементов
  const [isMobileFiltersOpen, setIsMobileFiltersOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  // Состояние для текущего изохрона
  const [currentIsochrone, setCurrentIsochrone] =
    useState<Feature<Polygon> | null>(null);

  // Функция для обновления URL без перезагрузки страницы
  const updateURL = useCallback(
    (newFilters: MapFilters, newViewState: MapViewState, query?: string) => {
      const params = new URLSearchParams();

      // Добавляем только непустые значения
      if (newFilters.category) params.set('category', newFilters.category);
      if (newFilters.priceFrom > 0)
        params.set('priceFrom', newFilters.priceFrom.toString());
      if (newFilters.priceTo > 0)
        params.set('priceTo', newFilters.priceTo.toString());
      if (newFilters.radius !== 10000)
        params.set('radius', newFilters.radius.toString());

      // Координаты карты
      params.set('lat', newViewState.latitude.toFixed(6));
      params.set('lng', newViewState.longitude.toFixed(6));
      params.set('zoom', newViewState.zoom.toFixed(2));

      // Поисковый запрос
      if (query) params.set('q', query);

      // Обновляем URL без перезагрузки
      const newURL = `${window.location.pathname}${params.toString() ? '?' + params.toString() : ''}`;
      window.history.replaceState({}, '', newURL);
    },
    []
  );

  // Определение мобильного устройства
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();
    window.addEventListener('resize', checkMobile);

    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // Отмечаем, что компонент инициализирован
  useEffect(() => {
    setIsInitialized(true);
  }, []);

  // Загрузка объявлений для карты
  const loadListings = useCallback(async () => {
    setIsLoading(true);
    try {
      const params = new URLSearchParams({
        limit: '100',
        page: '1',
        sort_by: 'date',
        sort_order: 'desc',
        ...(debouncedFilters.category && {
          categories: debouncedFilters.category,
        }),
        ...(debouncedFilters.priceFrom > 0 && {
          min_price: debouncedFilters.priceFrom.toString(),
        }),
        ...(debouncedFilters.priceTo > 0 && {
          max_price: debouncedFilters.priceTo.toString(),
        }),
      });

      // Используем позицию покупателя для географического поиска
      if (buyerLocation.latitude && buyerLocation.longitude) {
        params.append('latitude', buyerLocation.latitude.toString());
        params.append('longitude', buyerLocation.longitude.toString());

        // Преобразуем радиус из метров в формат для backend (например, "10km")
        if (debouncedFilters.radius) {
          const radiusKm = Math.round(debouncedFilters.radius / 1000);
          params.append('distance', `${radiusKm}km`);
        }
      }

      // Используем GIS API если есть координаты покупателя, иначе обычный search
      const endpoint =
        buyerLocation.latitude && buyerLocation.longitude
          ? '/api/v1/gis/search'
          : '/api/v1/search';

      // Логируем полный URL запроса и параметры
      const fullUrl = `${endpoint}?${params}`;
      console.log(
        '[Map] Using endpoint:',
        endpoint,
        'with params:',
        Object.fromEntries(params)
      );

      const response = await apiClient.get(fullUrl);
      console.log('[Map] API response:', response.data);
      console.log(
        '[Map] Listings count:',
        response.data?.data?.listings?.length ||
          response.data?.data?.length ||
          0
      );

      // Логируем цены объявлений для отладки
      if (response.data?.data?.listings) {
        const prices = response.data.data.listings.map((l: any) => ({
          id: l.id,
          price: l.price,
          title: l.title,
        }));
        console.log('[Map] Listings prices:', prices);
      }

      // Обрабатываем ответ в зависимости от используемого API
      if (endpoint === '/api/v1/gis/search' && response.data?.data?.listings) {
        // GIS API возвращает data.listings
        const transformedListings = response.data.data.listings
          .filter(
            (item: any) =>
              item.location && item.location.lat && item.location.lng
          )
          .map((item: any) => ({
            id: item.id,
            name: item.title,
            price: item.price,
            location: {
              lat: item.location.lat,
              lng: item.location.lng,
              city: item.address || '',
              country: 'Serbia',
            },
            category: {
              id: 0,
              name: item.category || 'Unknown',
              slug: '',
            },
            images: [],
            created_at: item.created_at,
          }));
        setListings(transformedListings);
      } else if (response.data?.items) {
        // Обычный search API возвращает items
        const transformedListings = response.data.items
          .filter(
            (item: any) =>
              item.location && item.location.lat && item.location.lng
          )
          .map((item: any) => ({
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
      } else {
        console.warn('[Map] Unknown API response format:', response.data);
        setListings([]);
      }
    } catch (error) {
      console.error('Error loading listings:', error);
      toast.error(t('errors.loadingFailed'));
    } finally {
      setIsLoading(false);
    }
  }, [debouncedFilters, buyerLocation, t]);

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
          longitude: listing.location.lng,
          latitude: listing.location.lat,
          title: listing.name,
          type: 'listing' as const,
          imageUrl: listing.images?.[0],
          metadata: {
            price: listing.price,
            currency: 'RSD',
            category: listing.category?.name || 'Unknown',
          },
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
        }));
    },
    []
  );

  // Получение цвета для категории (если понадобится в будущем)
  // const getCategoryColor = (categorySlug: string): string => {
  //   const colors: { [key: string]: string } = {
  //     'real-estate': '#3B82F6', // blue
  //     vehicles: '#EF4444', // red
  //     electronics: '#10B981', // green
  //     clothing: '#F59E0B', // amber
  //     services: '#8B5CF6', // violet
  //     jobs: '#F97316', // orange
  //     'children-goods-toys': '#EC4899', // pink
  //     'home-garden': '#16A34A', // green
  //     appliances: '#0EA5E9', // sky
  //     default: '#6B7280', // gray
  //   };
  //   return colors[categorySlug] || colors.default;
  // };

  // Поиск по адресу
  const handleAddressSearch = useCallback(
    async (query: string) => {
      if (!query.trim()) return;

      setIsSearching(true);
      setSearchQuery(query);

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

  // Обработка изменений фильтров и позиции покупателя
  useEffect(() => {
    loadListings();
  }, [
    loadListings,
    debouncedFilters.category,
    debouncedFilters.priceFrom,
    debouncedFilters.priceTo,
    debouncedFilters.radius,
    buyerLocation.latitude,
    buyerLocation.longitude,
  ]);

  // Создание маркеров при изменении объявлений с фильтрацией по изохрону
  useEffect(() => {
    let newMarkers = createMarkers(listings);

    // Фильтруем маркеры по изохрону если включен режим walking и есть изохрон
    if (walkingMode === 'walking' && currentIsochrone) {
      console.log('[Map] Filtering markers by isochrone');
      const filteredMarkers = newMarkers.filter((marker) => {
        const isInside = isPointInIsochrone(
          [marker.longitude, marker.latitude],
          currentIsochrone
        );
        return isInside;
      });
      console.log(
        `[Map] Filtered ${newMarkers.length} markers to ${filteredMarkers.length} within isochrone`
      );
      newMarkers = filteredMarkers;
    }

    console.log('[Map] Setting markers:', newMarkers);
    setMarkers(newMarkers);
  }, [listings, createMarkers, walkingMode, currentIsochrone]);

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

  // Обработчик изменения позиции покупателя
  const handleBuyerLocationChange = useCallback(
    (newLocation: { longitude: number; latitude: number }) => {
      setBuyerLocation(newLocation);
    },
    []
  );

  // Обработка изменения фильтров
  const handleFiltersChange = useCallback((newFilters: Partial<MapFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
  }, []);

  // Обновление URL при изменении фильтров, viewState или searchQuery
  useEffect(() => {
    if (isInitialized) {
      updateURL(filters, debouncedViewState, searchQuery);
    }
  }, [filters, debouncedViewState, searchQuery, updateURL, isInitialized]);

  return (
    <div className="min-h-screen bg-base-100">
      {/* Заголовок - скрываем на мобильном */}
      <div className="bg-white border-b border-base-300 px-4 py-4 hidden md:block">
        <div className="container mx-auto">
          <h1 className="text-2xl font-bold text-base-content mb-2">
            {t('title')}
          </h1>
          <p className="text-base-content-secondary">{t('description')}</p>
        </div>
      </div>

      {/* Контейнер с картой и фильтрами */}
      <div className="relative h-screen md:h-[calc(100vh-140px)]">
        {/* Десктопная боковая панель с фильтрами */}
        <div className="absolute left-4 top-4 z-10 w-80 bg-white rounded-lg shadow-lg hidden md:block">
          {/* Поиск по адресу */}
          <div className="p-4 border-b border-base-300">
            <label className="block text-sm font-medium text-base-content mb-2">
              {t('search.address')}
            </label>
            <SearchBar
              initialQuery={searchQuery}
              onSearch={handleAddressSearch}
              placeholder={t('search.addressPlaceholder')}
              className="w-full"
              geoLocation={
                viewState.latitude && viewState.longitude
                  ? {
                      lat: viewState.latitude,
                      lon: viewState.longitude,
                      radius: filters.radius,
                    }
                  : undefined
              }
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
                <option value="1100">Квартира</option>
                <option value="1200">Комната</option>
                <option value="1300">Дом, дача, коттедж</option>
                <option value="2000">Автомобили</option>
                <option value="3000">Электроника</option>
                <option value="9000">Работа</option>
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
              <WalkingAccessibilityControl
                mode={walkingMode}
                searchRadius={filters.radius}
                walkingTime={walkingTime}
                onModeChange={setWalkingMode}
                onRadiusChange={(radius) => handleFiltersChange({ radius })}
                onWalkingTimeChange={setWalkingTime}
              />
            </div>

            {/* Статистика */}
            <div className="text-sm text-base-content-secondary">
              {t('results.showing')}: {markers.length} {t('results.listings')}
            </div>
          </div>
        </div>

        {/* Мобильная кнопка фильтров */}
        <div className="absolute top-4 left-4 z-[1000] md:hidden">
          <button
            onClick={() => setIsMobileFiltersOpen(true)}
            className="bg-white rounded-lg shadow-lg p-3 flex items-center space-x-2 hover:bg-gray-50 transition-all duration-200 active:scale-95"
            aria-label="Открыть фильтры"
          >
            <svg
              className="w-5 h-5 text-gray-700"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.414A1 1 0 013 6.707V4z"
              />
            </svg>
            <span className="text-sm font-medium text-gray-700">
              {t('filters.title')}
            </span>
            {(filters.category ||
              filters.priceFrom > 0 ||
              filters.priceTo > 0) && (
              <span className="bg-primary text-white text-xs px-2 py-1 rounded-full min-w-[20px] h-5 flex items-center justify-center">
                {[
                  filters.category ? 1 : 0,
                  filters.priceFrom > 0 ? 1 : 0,
                  filters.priceTo > 0 ? 1 : 0,
                ].reduce((a, b) => a + b, 0)}
              </span>
            )}
          </button>
        </div>

        {/* Мобильный поиск */}
        <div className="absolute top-4 right-4 left-20 z-[1000] md:hidden">
          <SearchBar
            initialQuery={searchQuery}
            onSearch={handleAddressSearch}
            placeholder={t('search.addressPlaceholder')}
            className="w-full"
          />
        </div>

        {/* Карта */}
        <div className="absolute inset-0">
          <InteractiveMap
            initialViewState={viewState}
            markers={markers}
            onMarkerClick={handleMarkerClick}
            onViewStateChange={handleViewStateChange}
            className="w-full h-full"
            mapboxAccessToken={process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN}
            controlsConfig={{
              showNavigation: true,
              showFullscreen: true,
              showGeolocate: true,
              position: isMobile ? 'bottom-right' : 'top-right',
            }}
            isMobile={isMobile}
            showBuyerMarker={true}
            buyerLocation={buyerLocation}
            searchRadius={filters.radius}
            walkingMode={walkingMode}
            walkingTime={walkingTime}
            onBuyerLocationChange={handleBuyerLocationChange}
            onIsochroneChange={setCurrentIsochrone}
          />
        </div>

        {/* Мобильный drawer с фильтрами */}
        <MobileFiltersDrawer
          isOpen={isMobileFiltersOpen}
          onClose={() => setIsMobileFiltersOpen(false)}
          filters={filters}
          onFiltersChange={handleFiltersChange}
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          onSearch={handleAddressSearch}
          isSearching={isSearching}
          markersCount={markers.length}
          translations={{
            title: t('filters.title'),
            search: {
              address: t('search.address'),
              placeholder: t('search.addressPlaceholder'),
            },
            filters: {
              category: t('filters.category'),
              allCategories: t('filters.allCategories'),
              priceFrom: t('filters.priceFrom'),
              priceTo: t('filters.priceTo'),
              radius: t('filters.radius'),
            },
            categories: {
              realEstate: t('categories.realEstate'),
              vehicles: t('categories.vehicles'),
              electronics: t('categories.electronics'),
              clothing: t('categories.clothing'),
              services: t('categories.services'),
              jobs: t('categories.jobs'),
            },
            results: {
              showing: t('results.showing'),
              listings: t('results.listings'),
            },
            actions: {
              apply: t('actions.apply'),
              reset: t('actions.reset'),
            },
          }}
        />

        {/* Индикатор загрузки */}
        {isLoading && (
          <div className="absolute top-20 right-4 z-10 bg-white rounded-lg shadow-lg p-3 md:top-4">
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

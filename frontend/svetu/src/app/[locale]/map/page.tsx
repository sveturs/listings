'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { InteractiveMap } from '@/components/GIS';
import MarkerClickPopup from '@/components/GIS/Map/MarkerClickPopup';
import { useGeoSearch } from '@/components/GIS/hooks/useGeoSearch';
import {
  MapViewState,
  MapMarkerData,
  MapBounds,
} from '@/components/GIS/types/gis';
import { useDebounce } from '@/hooks/useDebounce';
import { SearchBar } from '@/components/SearchBar';
import { useRouter } from '@/i18n/routing';
import { useSearchParams } from 'next/navigation';
import { toast } from 'react-hot-toast';
import { apiClient } from '@/services/api-client';
import { MobileFiltersDrawer } from '@/components/GIS/Mobile';
// import WalkingAccessibilityControl from '@/components/GIS/Map/WalkingAccessibilityControl'; // Заменен на NativeSliderControl
import { isPointInIsochrone } from '@/components/GIS/utils/mapboxIsochrone';
import type { Feature, Polygon } from 'geojson';
// import { DistrictMapSelector } from '@/components/search';
import { SmartFilters } from '@/components/marketplace/SmartFilters';
import { QuickFilters } from '@/components/marketplace/QuickFilters';

// Функция для проверки, находится ли точка внутри полигона (Ray Casting Algorithm)
function isPointInPolygon(
  point: [number, number],
  polygon: [number, number][]
): boolean {
  const [x, y] = point;
  let inside = false;

  for (let i = 0, j = polygon.length - 1; i < polygon.length; j = i++) {
    const [xi, yi] = polygon[i];
    const [xj, yj] = polygon[j];

    if (yi > y !== yj > y && x < ((xj - xi) * (y - yi)) / (yj - yi) + xi) {
      inside = !inside;
    }
  }

  return inside;
}

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
  views_count?: number;
  rating?: number;
}

interface MapFilters {
  category: string;
  priceFrom: number;
  priceTo: number;
  radius: number;
  attributes?: Record<string, any>;
}

const MapPage: React.FC = () => {
  const t = useTranslations('map');
  const _router = useRouter();
  const searchParams = useSearchParams();
  const { search: geoSearch } = useGeoSearch();

  // Получаем язык из URL безопасно для SSR
  const [currentLang, setCurrentLang] = useState('sr');

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const lang = window.location.pathname.split('/')[1] || 'sr';
      setCurrentLang(lang);
    }
  }, []);

  // Функция для получения начальных значений из URL
  const getInitialFiltersFromURL = (): MapFilters => {
    const attributesStr = searchParams?.get('attributes');
    let attributes: Record<string, any> = {};

    if (attributesStr) {
      try {
        attributes = JSON.parse(decodeURIComponent(attributesStr));
      } catch (e) {
        console.error('Failed to parse attributes from URL', e);
      }
    }

    return {
      category: searchParams?.get('category') || '',
      priceFrom: parseInt(searchParams?.get('priceFrom') || '0') || 0,
      priceTo: parseInt(searchParams?.get('priceTo') || '0') || 0,
      radius: parseInt(searchParams?.get('radius') || '5000') || 5000,
      attributes,
    };
  };

  // Функция для получения начального состояния карты из URL
  const getInitialViewStateFromURL = (): MapViewState => {
    const lat = parseFloat(searchParams?.get('lat') || '44.8176');
    const lng = parseFloat(searchParams?.get('lng') || '20.4649');
    const zoom = parseFloat(searchParams?.get('zoom') || '10');

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

  // Состояние маркера покупателя - инициализируем с фиксированными координатами
  const [buyerLocation, setBuyerLocation] = useState({
    longitude: 20.457273, // Центр Белграда
    latitude: 44.787197,
  });

  // Дебаунсированная позиция покупателя
  const debouncedBuyerLocation = useDebounce(buyerLocation, 1000);

  // Данные и фильтры
  const [listings, setListings] = useState<ListingData[]>([]);
  const [markers, setMarkers] = useState<MapMarkerData[]>([]);
  const [filters, setFilters] = useState<MapFilters>(
    getInitialFiltersFromURL()
  );

  // Поиск
  const [searchQuery, setSearchQuery] = useState(searchParams?.get('q') || '');
  const [isSearchFromUser, setIsSearchFromUser] = useState(false);
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
  const [selectedMarker, setSelectedMarker] = useState<MapMarkerData | null>(
    null
  );
  const [isMobile, setIsMobile] = useState(false);

  // Состояние для текущего изохрона
  const [currentIsochrone, setCurrentIsochrone] =
    useState<Feature<Polygon> | null>(null);

  // Состояние для границ районов
  const [districtBoundary, setDistrictBoundary] =
    useState<Feature<Polygon> | null>(null);

  // Состояние для типа поиска (адрес или район)
  const [searchType, setSearchType] = useState<'address' | 'district'>(
    'address'
  );

  // Включить поиск по районам
  const _enableDistrictSearch = searchType === 'district';

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
      if (newFilters.radius !== 5000)
        params.set('radius', newFilters.radius.toString());
      if (
        newFilters.attributes &&
        Object.keys(newFilters.attributes).length > 0
      )
        params.set(
          'attributes',
          encodeURIComponent(JSON.stringify(newFilters.attributes))
        );

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
    // Проверяем, есть ли активный районный поиск без радиуса/изохрона
    const hasDistrictOnly =
      typeof window !== 'undefined' &&
      ((window as any).__DISTRICT_MARKERS_SET__ ||
        (window as any).__DISTRICT_PAGE_ACTIVE__) &&
      !debouncedBuyerLocation.latitude &&
      !debouncedBuyerLocation.longitude;

    if (hasDistrictOnly) {
      console.log('🚫 loadListings blocked: District-only search is active');
      return;
    }

    setIsLoading(true);
    try {
      // Определяем тип поиска
      const hasRadiusSearch =
        debouncedBuyerLocation.latitude && debouncedBuyerLocation.longitude;
      const hasDistrictBoundary = districtBoundary !== null;
      const isCombinedSearch = hasRadiusSearch && hasDistrictBoundary;

      console.log('🔍 Search type analysis:', {
        hasRadiusSearch,
        hasDistrictBoundary,
        isCombinedSearch,
        searchType,
      });

      // Используем специализированный радиусный поиск если есть координаты покупателя, иначе обычный search
      const useRadiusSearch = hasRadiusSearch;
      const endpoint = useRadiusSearch
        ? '/api/v1/gis/search/radius'
        : '/api/v1/search';

      let response;

      if (useRadiusSearch) {
        // Для радиусного поиска используем GET с query параметрами
        const params = new URLSearchParams({
          latitude: debouncedBuyerLocation.latitude.toString(),
          longitude: debouncedBuyerLocation.longitude.toString(),
          radius: debouncedFilters.radius.toString(), // в метрах
          limit: '100',
          ...(debouncedFilters.category && {
            category: debouncedFilters.category,
          }),
          ...(debouncedFilters.priceFrom > 0 && {
            min_price: debouncedFilters.priceFrom.toString(),
          }),
          ...(debouncedFilters.priceTo > 0 && {
            max_price: debouncedFilters.priceTo.toString(),
          }),
          ...(debouncedFilters.attributes &&
            Object.keys(debouncedFilters.attributes).length > 0 && {
              attributes: JSON.stringify(debouncedFilters.attributes),
            }),
        });

        const fullUrl = `${endpoint}?${params}`;

        // Добавляем заголовок для комбинированного поиска
        const headers: Record<string, string> = {};
        if (isCombinedSearch) {
          headers['X-Combined-Search'] = 'true';
          console.log(
            '🔍 Adding combined search header for district+radius search'
          );
        }

        response = await apiClient.get(fullUrl, { headers });
      } else {
        // Для обычного поиска используем GET с параметрами
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
          ...(debouncedFilters.attributes &&
            Object.keys(debouncedFilters.attributes).length > 0 && {
              attributes: JSON.stringify(debouncedFilters.attributes),
            }),
        });

        const fullUrl = `${endpoint}?${params}`;
        response = await apiClient.get(fullUrl);
      }

      // Обрабатываем ответ в зависимости от используемого API
      if (useRadiusSearch && response.data?.data?.listings) {
        // GIS API возвращает data.listings
        let filteredListings = response.data.data.listings.filter(
          (item: any) => item.location && item.location.lat && item.location.lng
        );

        // Если есть границы района, фильтруем по ним (комбинированный поиск)
        if (isCombinedSearch && districtBoundary) {
          console.log(
            '🔍 Applying district boundary filter to radius search results'
          );
          console.log(
            '📍 Before district filter:',
            filteredListings.length,
            'listings'
          );

          filteredListings = filteredListings.filter((item: any) => {
            const point: [number, number] = [
              item.location.lng,
              item.location.lat,
            ];
            const isInside = isPointInPolygon(
              point,
              districtBoundary.geometry.coordinates[0] as [number, number][]
            );
            return isInside;
          });

          console.log(
            '📍 After district filter:',
            filteredListings.length,
            'listings'
          );
        }

        const transformedListings = filteredListings.map((item: any) => ({
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
          images: item.images || [],
          created_at: item.created_at,
          // Добавляем новые поля
          views_count: item.views_count || 0,
          rating: item.rating || 0,
        }));

        console.log(
          '🗺️ Setting combined search results:',
          transformedListings.length
        );
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
            // Добавляем новые поля
            views_count: item.views_count || 0,
            rating: item.rating || 0,
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
  }, [
    debouncedFilters,
    debouncedBuyerLocation,
    districtBoundary,
    searchType,
    t,
  ]);

  // Функция для получения иконки по категории
  const getCategoryIcon = (categoryName: string | undefined): string => {
    if (!categoryName) return '🏠';

    const category = categoryName.toLowerCase();

    // Автомобили
    if (
      category.includes('автомобил') ||
      category.includes('car') ||
      category.includes('vozilo')
    )
      return '🚗';
    // Недвижимость
    if (
      category.includes('квартир') ||
      category.includes('apartment') ||
      category.includes('stan')
    )
      return '🏠';
    if (
      category.includes('дом') ||
      category.includes('house') ||
      category.includes('kuća')
    )
      return '🏘️';
    if (
      category.includes('комнат') ||
      category.includes('room') ||
      category.includes('soba')
    )
      return '🛏️';
    // Электроника
    if (
      category.includes('телефон') ||
      category.includes('phone') ||
      category.includes('telefon')
    )
      return '📱';
    if (
      category.includes('компьютер') ||
      category.includes('computer') ||
      category.includes('računar')
    )
      return '💻';
    if (
      category.includes('телевизор') ||
      category.includes('tv') ||
      category.includes('televizor')
    )
      return '📺';
    // Работа
    if (
      category.includes('работ') ||
      category.includes('job') ||
      category.includes('posao')
    )
      return '💼';
    // Услуги
    if (
      category.includes('услуг') ||
      category.includes('service') ||
      category.includes('usluga')
    )
      return '🔧';
    // Одежда
    if (
      category.includes('одежд') ||
      category.includes('cloth') ||
      category.includes('odeća')
    )
      return '👕';
    // Спорт
    if (category.includes('спорт') || category.includes('sport')) return '⚽';
    // По умолчанию
    return '📦';
  };

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
            icon: getCategoryIcon(listing.category?.name),
          },
          data: {
            title: listing.name,
            price: listing.price,
            category: listing.category?.name || 'Unknown',
            image: (listing as any).images?.[0] || listing.images?.[0],
            address:
              `${listing.location.city || ''}, ${listing.location.country || ''}`
                .trim()
                .replace(/^,\s*|,\s*$/, ''),
            id: listing.id,
            icon: getCategoryIcon(listing.category?.name),
            views_count: (listing as any).views_count || 0,
            rating: (listing as any).rating || 0,
            created_at: listing.created_at,
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
      setIsSearchFromUser(true);
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

          // Обновляем позицию покупателя на найденную локацию
          setBuyerLocation({
            longitude: parseFloat(result.lon),
            latitude: parseFloat(result.lat),
          });

          toast.success(t('search.found'));
        } else {
          toast.error(t('search.notFound'));
        }
      } catch (error) {
        console.error('Search error:', error);
        toast.error(t('search.error'));
      } finally {
        setIsSearching(false);
        setIsSearchFromUser(false);
      }
    },
    [geoSearch, viewState, t]
  );

  // Обработка поиска
  useEffect(() => {
    if (debouncedSearchQuery && isSearchFromUser) {
      handleAddressSearch(debouncedSearchQuery);
    }
  }, [debouncedSearchQuery, handleAddressSearch, isSearchFromUser]);

  // Обработка изменений фильтров и позиции покупателя
  useEffect(() => {
    loadListings();
  }, [
    loadListings,
    debouncedFilters.category,
    debouncedFilters.priceFrom,
    debouncedFilters.priceTo,
    debouncedFilters.radius,
    debouncedFilters.attributes,
    debouncedBuyerLocation.latitude,
    debouncedBuyerLocation.longitude,
  ]);

  // Создание маркеров при изменении объявлений с фильтрацией по изохрону
  useEffect(() => {
    let newMarkers = createMarkers(listings);
    console.log(
      '🗺️ Creating markers from listings:',
      listings.length,
      '→',
      newMarkers.length
    );

    // Фильтруем маркеры по изохрону если включен режим walking и есть изохрон
    if (walkingMode === 'walking' && currentIsochrone) {
      console.log('🚶 Applying isochrone filter, mode:', walkingMode);
      const filteredMarkers = newMarkers.filter((marker) => {
        const isInside = isPointInIsochrone(
          [marker.longitude, marker.latitude],
          currentIsochrone
        );
        return isInside;
      });
      console.log(
        '🚶 After isochrone filter:',
        newMarkers.length,
        '→',
        filteredMarkers.length
      );
      newMarkers = filteredMarkers;
    }

    console.log('🗺️ Final markers count:', newMarkers.length);
    setMarkers(newMarkers);
  }, [listings, createMarkers, walkingMode, currentIsochrone]);

  // Обработка клика по маркеру
  const handleMarkerClick = useCallback((marker: MapMarkerData) => {
    // Показываем расширенный popup вместо мгновенного перехода
    setSelectedMarker(marker);
  }, []);

  // Обработчик результатов поиска по районам
  const handleDistrictSearchResults = useCallback((results: any[]) => {
    console.log(
      '🔍 District search results received:',
      results.length,
      'items'
    );
    console.log('First result example:', results[0]);

    // Устанавливаем флаг, что активен районный поиск
    if (typeof window !== 'undefined') {
      (window as any).__DISTRICT_MARKERS_SET__ = true;
      (window as any).__DISTRICT_PAGE_ACTIVE__ = true;
      setTimeout(() => {
        delete (window as any).__DISTRICT_MARKERS_SET__;
      }, 3000); // Защита на 3 секунды
    }

    // Преобразуем результаты в формат ListingData
    const transformedListings = results
      .filter((item: any) => item.location?.lat && item.location?.lng)
      .map((item: any) => ({
        id: parseInt(item.id),
        name: item.title || 'Untitled',
        price: item.price || 0,
        location: {
          lat: item.location.lat,
          lng: item.location.lng,
          city: item.location.address || '',
          country: 'Serbia',
        },
        category: {
          id: 0,
          name: item.category || 'Unknown',
          slug: '',
        },
        images: item.images || [],
        created_at: item.created_at || new Date().toISOString(),
        // Добавляем новые поля
        views_count: item.views_count || 0,
        rating: item.rating || 0,
      }));

    console.log('📍 Transformed listings:', transformedListings.length);
    console.log(
      '🗺️ Setting district listings on main page:',
      transformedListings.length
    );
    setListings(transformedListings);
  }, []);

  // Обработчик изменения границ района
  const handleDistrictBoundsChange = useCallback(
    (bounds: [number, number, number, number] | null) => {
      if (bounds) {
        const [minLng, minLat, maxLng, maxLat] = bounds;
        const centerLng = (minLng + maxLng) / 2;
        const centerLat = (minLat + maxLat) / 2;

        // Рассчитываем уровень зума чтобы вместить весь район
        const lngDiff = maxLng - minLng;
        const latDiff = maxLat - minLat;
        const maxDiff = Math.max(lngDiff, latDiff);

        let zoom = 12;
        if (maxDiff < 0.05) zoom = 14;
        else if (maxDiff < 0.1) zoom = 13;
        else if (maxDiff < 0.2) zoom = 12;
        else if (maxDiff < 0.4) zoom = 11;
        else zoom = 10;

        setViewState({
          ...viewState,
          longitude: centerLng,
          latitude: centerLat,
          zoom: zoom,
        });

        // Также обновляем позицию покупателя на центр района
        setBuyerLocation({
          longitude: centerLng,
          latitude: centerLat,
        });
      }
    },
    [viewState]
  );

  // Текущий viewport для передачи в DistrictMapSelector
  const [currentMapViewport, setCurrentMapViewport] = useState<{
    bounds: MapBounds;
    center: { lat: number; lng: number };
  } | null>(null);

  // Обработка изменения области просмотра
  const handleViewStateChange = useCallback((newViewState: MapViewState) => {
    setViewState(newViewState);
  }, []);

  // Дебаунсированное обновление viewport для DistrictMapSelector
  useEffect(() => {
    const timer = setTimeout(() => {
      // Вычисляем bounds из viewport
      const zoomFactor = Math.pow(2, 14 - viewState.zoom) * 0.01;
      const bounds: MapBounds = {
        north: viewState.latitude + zoomFactor,
        south: viewState.latitude - zoomFactor,
        east: viewState.longitude + zoomFactor,
        west: viewState.longitude - zoomFactor,
      };

      setCurrentMapViewport({
        bounds,
        center: {
          lat: viewState.latitude,
          lng: viewState.longitude,
        },
      });
    }, 500); // Дебаунс в 500мс

    return () => clearTimeout(timer);
  }, [viewState.latitude, viewState.longitude, viewState.zoom]);

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

  // Обработчик для быстрых фильтров
  const handleQuickFilterSelect = useCallback(
    (quickFilters: Record<string, any>) => {
      setFilters((prev) => ({
        ...prev,
        attributes: {
          ...prev.attributes,
          ...quickFilters,
        },
      }));
    },
    []
  );

  // Обработчик изменения радиуса поиска
  const handleSearchRadiusChange = useCallback(
    (radius: number) => {
      handleFiltersChange({ radius });
    },
    [handleFiltersChange]
  );

  // Обновление URL при изменении фильтров, viewState или searchQuery
  useEffect(() => {
    if (isInitialized) {
      updateURL(filters, debouncedViewState, searchQuery);
    }
  }, [filters, debouncedViewState, searchQuery, updateURL, isInitialized]);

  // Memoized переводы для контролов
  const controlTranslations = useMemo(
    () => ({
      walkingAccessibility: t('controls.walkingAccessibility'),
      searchRadius: t('controls.searchRadius'),
      minutes: t('controls.minutes'),
      km: t('controls.km'),
      m: t('controls.m'),
      changeModeHint: t('controls.changeModeHint'),
      holdForSettings: t('controls.holdForSettings'),
      singleClickHint: t('controls.singleClickHint'),
      mobileHint: t('controls.mobileHint'),
      desktopHint: t('controls.desktopHint'),
      updatingIsochrone: t('controls.updatingIsochrone'),
    }),
    [t]
  );

  return (
    <div className="min-h-screen bg-base-100">
      {/* Контейнер с картой и фильтрами */}
      <div className="relative h-[100dvh] md:h-[calc(100vh-140px)]">
        {/* Десктопная боковая панель с фильтрами */}
        <div className="absolute left-4 top-4 z-10 w-80 bg-white rounded-lg shadow-lg hidden md:block">
          {/* Поиск по адресу */}
          <div className="p-4 border-b border-base-300">
            {/* Поиск по адресу - всегда показываем, так как районы отключены */}
            <label className="block text-sm font-medium text-base-content mb-2">
              {t('search.address')}
            </label>
            <SearchBar
              initialQuery={searchQuery}
              onSearch={(query) => {
                setIsSearchFromUser(true);
                handleAddressSearch(query);
              }}
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

            {/* Динамические фильтры по атрибутам категории */}
            {filters.category && (
              <div className="mb-4">
                <SmartFilters
                  categoryId={parseInt(filters.category) || null}
                  onChange={(attributeFilters) =>
                    handleFiltersChange({ attributes: attributeFilters })
                  }
                  lang={currentLang}
                  className="space-y-3"
                />
              </div>
            )}

            {/* Контроль радиуса поиска */}
            <div className="mb-4 space-y-3">
              <label className="block text-sm font-medium text-base-content">
                {t('controls.radiusControl')}
              </label>

              {/* Переключатель типа радиуса */}
              <div className="flex gap-1 p-1 bg-base-200 rounded-lg">
                <button
                  type="button"
                  className={`flex-1 px-3 py-1.5 text-xs font-medium rounded transition-colors ${
                    walkingMode === 'walking'
                      ? 'bg-primary text-primary-content'
                      : 'text-base-content hover:bg-base-300'
                  }`}
                  onClick={() => setWalkingMode('walking')}
                >
                  🚶 {t('controls.walkingMode')}
                </button>
                <button
                  type="button"
                  className={`flex-1 px-3 py-1.5 text-xs font-medium rounded transition-colors ${
                    walkingMode === 'radius'
                      ? 'bg-primary text-primary-content'
                      : 'text-base-content hover:bg-base-300'
                  }`}
                  onClick={() => setWalkingMode('radius')}
                >
                  📏 {t('controls.distanceMode')}
                </button>
              </div>

              {/* Слайдер радиуса */}
              <div className="space-y-2">
                <div className="flex justify-between text-sm text-base-content/70">
                  <span>
                    {walkingMode === 'walking'
                      ? `5 ${t('controls.minUnit')}`
                      : `0.1 ${t('controls.kmUnit')}`}
                  </span>
                  <span className="font-medium">
                    {walkingMode === 'walking'
                      ? `${walkingTime} ${t('controls.minUnit')}`
                      : `${filters.radius >= 1000 ? (filters.radius / 1000).toFixed(1) : (filters.radius / 1000).toFixed(1)} ${t('controls.kmUnit')}`}
                  </span>
                  <span>
                    {walkingMode === 'walking'
                      ? `60 ${t('controls.minUnit')}`
                      : `50 ${t('controls.kmUnit')}`}
                  </span>
                </div>
                <input
                  type="range"
                  className="range range-primary range-sm"
                  min={walkingMode === 'walking' ? 5 : 100}
                  max={walkingMode === 'walking' ? 60 : 50000}
                  step={walkingMode === 'walking' ? 5 : 100}
                  value={
                    walkingMode === 'walking' ? walkingTime : filters.radius
                  }
                  onChange={(e) => {
                    const value = Number(e.target.value);
                    if (walkingMode === 'walking') {
                      setWalkingTime(value);
                    } else {
                      handleFiltersChange({ radius: value });
                    }
                  }}
                />
              </div>
            </div>

            {/* Быстрые фильтры для категории */}
            {filters.category && (
              <div className="mb-4">
                <QuickFilters
                  categoryId={filters.category}
                  onSelectFilter={handleQuickFilterSelect}
                />
              </div>
            )}

            {/* Умные фильтры по атрибутам категории */}
            {filters.category && (
              <div className="mb-4 border-t pt-4">
                <SmartFilters
                  categoryId={parseInt(filters.category) || null}
                  onChange={(attributeFilters) =>
                    handleFiltersChange({ attributes: attributeFilters })
                  }
                  lang={currentLang}
                />
              </div>
            )}

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
            onSearch={(query) => {
              setIsSearchFromUser(true);
              handleAddressSearch(query);
            }}
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
            onWalkingModeChange={setWalkingMode}
            onWalkingTimeChange={setWalkingTime}
            onSearchRadiusChange={handleSearchRadiusChange}
            useNativeControl={true} // Используем нативный контрол по умолчанию
            controlTranslations={controlTranslations}
            districtBoundary={districtBoundary}
          />

          {/* Расширенный popup при клике */}
          {selectedMarker && (
            <MarkerClickPopup
              marker={selectedMarker}
              onClose={() => setSelectedMarker(null)}
            />
          )}
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
          enableDistrictSearch={searchType === 'district'}
          onDistrictSearchResults={handleDistrictSearchResults}
          onDistrictBoundsChange={handleDistrictBoundsChange}
          onDistrictBoundaryChange={setDistrictBoundary}
          currentViewport={currentMapViewport}
          searchType={searchType}
          onSearchTypeChange={setSearchType}
          translations={{
            title: t('filters.title'),
            search: {
              address: t('search.address'),
              placeholder: t('search.addressPlaceholder'),
              byAddress: t('search.byAddress'),
              byDistrict: t('search.byDistrict'),
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

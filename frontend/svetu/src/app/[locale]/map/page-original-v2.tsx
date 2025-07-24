'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { InteractiveMap } from '@/components/GIS';
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
import { isPointInIsochrone } from '@/components/GIS/utils/mapboxIsochrone';
import type { Feature, Polygon } from 'geojson';
import { SmartFilters } from '@/components/marketplace/SmartFilters';
import { QuickFilters } from '@/components/marketplace/QuickFilters';
import { CategoryTreeSelector } from '@/components/common/CategoryTreeSelector';

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
  name?: string;
  title?: string;
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
  individual_address?: string;
  location_privacy?: string;
}

interface MapFilters {
  categories: number[];
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

    // Получаем категории из URL (поддерживаем как одну, так и несколько)
    const categoriesParam =
      searchParams?.get('categories') || searchParams?.get('category') || '';
    let categories: number[] = [];
    if (categoriesParam) {
      // Если это строка с запятыми, разбиваем на массив
      if (categoriesParam.includes(',')) {
        categories = categoriesParam
          .split(',')
          .map((c) => parseInt(c))
          .filter((c) => !isNaN(c));
      } else {
        const parsed = parseInt(categoriesParam);
        if (!isNaN(parsed)) {
          categories = [parsed];
        }
      }
    }

    return {
      categories,
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
    const zoom = parseFloat(searchParams?.get('zoom') || '11');

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

  // Получаем начальные координаты из URL для buyerLocation
  const initialViewState = getInitialViewStateFromURL();

  // Состояние маркера покупателя - инициализируем с координатами из URL
  const [buyerLocation, setBuyerLocation] = useState({
    longitude: initialViewState.longitude,
    latitude: initialViewState.latitude,
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

  // Состояние для сворачивания левой панели
  const [isLeftPanelCollapsed, setIsLeftPanelCollapsed] = useState(false);

  // Состояние для раскрытия секции фильтров
  const [isFiltersExpanded, setIsFiltersExpanded] = useState(false);

  // Функция для обновления URL без перезагрузки страницы
  const updateURL = useCallback(
    (newFilters: MapFilters, newViewState: MapViewState, query?: string) => {
      const params = new URLSearchParams();

      // Добавляем только непустые значения
      if (newFilters.categories && newFilters.categories.length > 0) {
        params.set('categories', newFilters.categories.join(','));
      }
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

  // Отмечаем, что компонент инициализирован после небольшой задержки
  // чтобы избежать перезаписи URL параметров при первой загрузке
  useEffect(() => {
    const timer = setTimeout(() => {
      setIsInitialized(true);
    }, 1000);
    return () => clearTimeout(timer);
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

    console.log('🔍 Loading listings with filters:', {
      categories: debouncedFilters.categories,
      priceFrom: debouncedFilters.priceFrom,
      priceTo: debouncedFilters.priceTo,
      radius: debouncedFilters.radius,
      buyerLocation: debouncedBuyerLocation,
    });

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
        buyerLat: debouncedBuyerLocation.latitude,
        buyerLng: debouncedBuyerLocation.longitude,
        endpoint: hasRadiusSearch
          ? '/api/v1/gis/search/radius'
          : '/api/v1/search',
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
          ...(debouncedFilters.categories &&
            debouncedFilters.categories.length > 0 && {
              categories: debouncedFilters.categories.join(','),
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
        console.log('📡 GIS API Request:', fullUrl);

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
          ...(debouncedFilters.categories &&
            debouncedFilters.categories.length > 0 && {
              categories: debouncedFilters.categories.join(','),
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
        console.log('📡 Search API Request:', fullUrl);
        response = await apiClient.get(fullUrl);
      }

      // Обрабатываем ответ в зависимости от используемого API
      if (useRadiusSearch && response.data?.data) {
        // GIS API возвращает data.listings (может быть null)
        console.log('[Map] GIS API response:', {
          success: response.data.success,
          totalCount: response.data.data.total_count,
          hasListings: !!response.data.data.listings,
          listingsCount: response.data.data.listings?.length || 0,
        });
        const apiListings = response.data.data.listings || [];
        let filteredListings = apiListings.filter(
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
          individual_address: item.individual_address || item.address,
          location_privacy: item.privacy_level || item.location_privacy,
        }));

        console.log(
          '🗺️ GIS API results:',
          transformedListings.length,
          'listings',
          'First few listings:',
          transformedListings.slice(0, 3).map((l: any) => ({
            id: l.id,
            name: l.name,
            category: l.category,
            location: l.location,
          }))
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
            individual_address: item.individual_address || item.address,
            location_privacy: item.privacy_level || item.location_privacy,
          }));
        console.log(
          '🗺️ Search API results:',
          transformedListings.length,
          'listings',
          'Requested categories:',
          debouncedFilters.categories,
          'First few listings:',
          transformedListings.slice(0, 3).map((l: any) => ({
            id: l.id,
            name: l.name,
            category: l.category,
            location: l.location,
          }))
        );
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
          title: listing.title || listing.name || 'Untitled',
          type: 'listing' as const,
          imageUrl: listing.images?.[0],
          metadata: {
            price: listing.price,
            currency: 'RSD',
            category: listing.category?.name || 'Unknown',
            icon: getCategoryIcon(listing.category?.name),
          },
          data: {
            title: listing.title || listing.name || 'Untitled',
            price: listing.price,
            category: listing.category?.name || 'Unknown',
            image: (listing as any).images?.[0] || listing.images?.[0],
            address:
              listing.individual_address ||
              (listing as any).address ||
              `${listing.location.city || ''}, ${listing.location.country || ''}`
                .trim()
                .replace(/^,\s*|,\s*$/, ''),
            locationPrivacy: listing.location_privacy,
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
    // Используем JSON.stringify для массива категорий, чтобы правильно отслеживать изменения
    JSON.stringify(debouncedFilters.categories),
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

  // Популярные районы
  const popularDistricts = [
    { name: 'Нови Београд', lat: 44.8094, lng: 20.3864, zoom: 13 },
    { name: 'Земун', lat: 44.8433, lng: 20.4011, zoom: 13 },
    { name: 'Врачар', lat: 44.7988, lng: 20.4724, zoom: 14 },
    { name: 'Савски венац', lat: 44.7879, lng: 20.4573, zoom: 13 },
  ];

  // Быстрые категории
  const quickCategories = [
    { icon: '🏠', name: t('categories.realEstate'), id: 1 },
    { icon: '🚗', name: t('categories.vehicles'), id: 2 },
    { icon: '📱', name: t('categories.electronics'), id: 3 },
    { icon: '👕', name: t('categories.clothing'), id: 4 },
    { icon: '🔧', name: t('categories.services'), id: 5 },
    { icon: '💼', name: t('categories.jobs'), id: 6 },
  ];

  return (
    <div className="relative w-full h-screen overflow-hidden bg-base-100">
      {/* Карта на весь экран */}
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
            position: 'top-right',
          }}
          isMobile={isMobile}
          selectedMarker={selectedMarker}
          onMarkerClose={() => setSelectedMarker(null)}
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
          useNativeControl={true}
          controlTranslations={controlTranslations}
          districtBoundary={districtBoundary}
        />
      </div>

      {/* Левая панель - поиск и категории */}
      <div
        className={`absolute left-0 top-0 bottom-0 ${isLeftPanelCollapsed ? 'w-12' : 'w-80'} bg-base-100 shadow-2xl flex flex-col z-20 transition-all duration-300 ${isMobile ? '-translate-x-full' : ''}`}
      >
        {/* Кнопка сворачивания/разворачивания */}
        <button
          onClick={() => setIsLeftPanelCollapsed(!isLeftPanelCollapsed)}
          className={`absolute ${isLeftPanelCollapsed ? 'left-3' : '-right-3'} top-6 z-30 btn btn-circle btn-sm bg-base-100 hover:bg-base-200 shadow-md`}
          title={isLeftPanelCollapsed ? 'Развернуть панель' : 'Свернуть панель'}
        >
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d={isLeftPanelCollapsed ? 'M9 5l7 7-7 7' : 'M15 19l-7-7 7-7'}
            />
          </svg>
        </button>

        {/* Лого и поиск */}
        <div
          className={`p-4 border-b border-base-300 ${isLeftPanelCollapsed ? 'hidden' : ''}`}
        >
          <div className="flex items-center gap-2 mb-4">
            <h1 className="text-2xl font-bold">SveTu</h1>
            <div className="badge badge-primary">{markers.length}</div>
          </div>

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

        {/* Быстрые фильтры */}
        <div
          className={`p-4 border-b border-base-300 ${isLeftPanelCollapsed ? 'hidden' : ''}`}
        >
          <h3 className="text-sm font-semibold mb-3 text-base-content/70">
            {t('categories.title')}
          </h3>
          <div className="grid grid-cols-3 gap-2">
            {quickCategories.map((cat) => (
              <button
                key={cat.id}
                onClick={() => {
                  const isSelected = filters.categories.includes(cat.id);
                  handleFiltersChange({
                    categories: isSelected
                      ? filters.categories.filter((c) => c !== cat.id)
                      : [...filters.categories, cat.id],
                  });
                }}
                className={`btn btn-sm ${filters.categories.includes(cat.id) ? 'btn-primary' : 'btn-ghost'} flex flex-col h-auto py-2`}
              >
                <span className="text-xl">{cat.icon}</span>
                <span className="text-xs">{cat.name}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Популярные районы */}
        <div
          className={`p-4 border-b border-base-300 ${isLeftPanelCollapsed ? 'hidden' : ''}`}
        >
          <h3 className="text-sm font-semibold mb-3 text-base-content/70">
            {t('popularDistricts')}
          </h3>
          <div className="flex flex-wrap gap-2">
            {popularDistricts.map((district) => (
              <button
                key={district.name}
                onClick={() => {
                  setViewState({
                    ...viewState,
                    latitude: district.lat,
                    longitude: district.lng,
                    zoom: district.zoom,
                  });
                  setBuyerLocation({
                    latitude: district.lat,
                    longitude: district.lng,
                  });
                }}
                className="btn btn-xs btn-outline"
              >
                {district.name}
              </button>
            ))}
          </div>
        </div>

        {/* Дополнительные фильтры */}
        <div
          className={`flex-1 overflow-y-auto ${isLeftPanelCollapsed ? 'hidden' : ''}`}
        >
          {/* Кнопка-заголовок для фильтров */}
          <div className="px-4 pt-4">
            <button
              onClick={() => setIsFiltersExpanded(!isFiltersExpanded)}
              className="w-full flex items-center justify-between p-3 rounded-lg hover:bg-base-200 transition-colors"
            >
              <div className="flex items-center gap-2">
                <svg
                  className="w-4 h-4"
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
                <span className="font-medium">{t('filters.title')}</span>
                {(filters.priceFrom > 0 ||
                  filters.priceTo > 0 ||
                  filters.categories.length > 0) && (
                  <div className="badge badge-primary badge-sm">
                    {filters.categories.length +
                      (filters.priceFrom > 0 ? 1 : 0) +
                      (filters.priceTo > 0 ? 1 : 0)}
                  </div>
                )}
              </div>
              <svg
                className={`w-4 h-4 transition-transform ${isFiltersExpanded ? 'rotate-180' : ''}`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </button>
          </div>

          {/* Содержимое фильтров */}
          {isFiltersExpanded && (
            <div className="p-4 space-y-4">
              {/* Категория */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('filters.category')}</span>
                </label>
                <CategoryTreeSelector
                  value={filters.categories}
                  onChange={(value) => {
                    const categories = Array.isArray(value)
                      ? value
                      : value
                        ? [value]
                        : [];
                    handleFiltersChange({ categories });
                  }}
                  multiple={true}
                  placeholder={t('filters.allCategories')}
                  showPath={true}
                  className="w-full"
                />
              </div>

              {/* Цена */}
              <div>
                <label className="label">
                  <span className="label-text">{t('filters.price')}</span>
                </label>
                <div className="grid grid-cols-2 gap-2">
                  <div className="form-control">
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      value={filters.priceFrom || ''}
                      onChange={(e) =>
                        handleFiltersChange({
                          priceFrom: parseInt(e.target.value) || 0,
                        })
                      }
                      placeholder={t('filters.priceFrom')}
                    />
                  </div>
                  <div className="form-control">
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      value={filters.priceTo || ''}
                      onChange={(e) =>
                        handleFiltersChange({
                          priceTo: parseInt(e.target.value) || 0,
                        })
                      }
                      placeholder={t('filters.priceTo')}
                    />
                  </div>
                </div>
              </div>

              {/* Радиус поиска */}
              <div>
                <label className="label">
                  <span className="label-text">
                    {t('controls.radiusControl')}
                  </span>
                </label>

                {/* Переключатель режима */}
                <div className="tabs tabs-boxed tabs-sm mb-3">
                  <a
                    className={`tab ${walkingMode === 'walking' ? 'tab-active' : ''}`}
                    onClick={() => setWalkingMode('walking')}
                  >
                    🚶 {t('controls.walkingMode')}
                  </a>
                  <a
                    className={`tab ${walkingMode === 'radius' ? 'tab-active' : ''}`}
                    onClick={() => setWalkingMode('radius')}
                  >
                    📏 {t('controls.distanceMode')}
                  </a>
                </div>

                {/* Слайдер */}
                <div className="space-y-2">
                  <div className="flex justify-between text-xs">
                    <span>
                      {walkingMode === 'walking'
                        ? `5 ${t('controls.minUnit')}`
                        : `0.1 ${t('controls.kmUnit')}`}
                    </span>
                    <span className="font-medium badge badge-primary badge-sm">
                      {walkingMode === 'walking'
                        ? `${walkingTime} ${t('controls.minUnit')}`
                        : `${(filters.radius / 1000).toFixed(1)} ${t('controls.kmUnit')}`}
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

              {/* Динамические фильтры */}
              {filters.categories && filters.categories.length > 0 && (
                <div>
                  <SmartFilters
                    categoryId={filters.categories[0]}
                    onChange={(attributeFilters) =>
                      handleFiltersChange({ attributes: attributeFilters })
                    }
                    lang={currentLang}
                    className="space-y-3"
                  />
                </div>
              )}

              {/* Быстрые фильтры */}
              {filters.categories && filters.categories.length > 0 && (
                <div>
                  <QuickFilters
                    categoryId={filters.categories[0].toString()}
                    onSelectFilter={handleQuickFilterSelect}
                  />
                </div>
              )}

              {/* Кнопка сброса */}
              <button
                onClick={() => {
                  setFilters({
                    categories: [],
                    priceFrom: 0,
                    priceTo: 0,
                    radius: 5000,
                    attributes: {},
                  });
                }}
                className="btn btn-outline btn-sm btn-block"
              >
                {t('filters.resetFilters')}
              </button>
            </div>
          )}
        </div>

        {/* Свёрнутое состояние - вертикальные иконки */}
        {isLeftPanelCollapsed && (
          <div className="flex flex-col items-center py-4 gap-3">
            {/* Поиск */}
            <button
              onClick={() => setIsLeftPanelCollapsed(false)}
              className="btn btn-ghost btn-sm btn-square"
              title="Поиск"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </button>

            {/* Категории */}
            <button
              onClick={() => setIsLeftPanelCollapsed(false)}
              className="btn btn-ghost btn-sm btn-square"
              title="Категории"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
                />
              </svg>
            </button>

            {/* Фильтры */}
            <button
              onClick={() => {
                setIsLeftPanelCollapsed(false);
                setIsFiltersExpanded(true);
              }}
              className="btn btn-ghost btn-sm btn-square"
              title="Фильтры"
            >
              <svg
                className="w-5 h-5"
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
            </button>
          </div>
        )}
      </div>

      {/* Мобильная кнопка меню */}
      {isMobile && (
        <button
          onClick={() => setIsMobileFiltersOpen(true)}
          className="btn btn-circle btn-primary fixed top-4 left-4 shadow-xl z-30"
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </button>
      )}

      {/* Плавающие кнопки быстрых действий */}
      <div className="absolute bottom-6 right-6 flex flex-col gap-3 z-10">
        {/* Геолокация */}
        <button
          onClick={() => {
            if (navigator.geolocation) {
              navigator.geolocation.getCurrentPosition(
                (position) => {
                  const { latitude, longitude } = position.coords;
                  setViewState({
                    ...viewState,
                    latitude,
                    longitude,
                    zoom: 15,
                  });
                  setBuyerLocation({ latitude, longitude });
                },
                () => {
                  toast.error(t('geolocation.error'));
                }
              );
            }
          }}
          className="btn btn-circle btn-lg bg-base-100 shadow-xl hover:shadow-2xl"
          title={t('geolocation.findMe')}
        >
          📍
        </button>

        {/* Показать все маркеры */}
        {markers.length > 0 && (
          <button
            onClick={() => {
              // Вычисляем границы всех маркеров
              const lats = markers.map((m) => m.latitude);
              const lngs = markers.map((m) => m.longitude);
              const minLat = Math.min(...lats);
              const maxLat = Math.max(...lats);
              const minLng = Math.min(...lngs);
              const maxLng = Math.max(...lngs);

              const centerLat = (minLat + maxLat) / 2;
              const centerLng = (minLng + maxLng) / 2;

              // Рассчитываем zoom чтобы показать все маркеры
              const latDiff = maxLat - minLat;
              const lngDiff = maxLng - minLng;
              const maxDiff = Math.max(latDiff, lngDiff);

              let zoom = 10;
              if (maxDiff < 0.01) zoom = 15;
              else if (maxDiff < 0.05) zoom = 13;
              else if (maxDiff < 0.1) zoom = 12;
              else if (maxDiff < 0.5) zoom = 10;
              else zoom = 8;

              setViewState({
                ...viewState,
                latitude: centerLat,
                longitude: centerLng,
                zoom,
              });
            }}
            className="btn btn-circle btn-lg bg-base-100 shadow-xl hover:shadow-2xl"
            title={t('showAll')}
          >
            🔍
          </button>
        )}
      </div>

      {/* Индикатор загрузки */}
      {isLoading && (
        <div className="absolute top-20 left-1/2 transform -translate-x-1/2 z-30">
          <div className="alert alert-info shadow-lg">
            <span className="loading loading-spinner loading-sm"></span>
            <span>{t('loading')}</span>
          </div>
        </div>
      )}

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
    </div>
  );
};

export default MapPage;

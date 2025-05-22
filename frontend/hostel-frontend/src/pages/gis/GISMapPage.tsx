import React, { useState, useEffect, useRef, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useLocation } from '../../contexts/LocationContext';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet.markercluster/dist/MarkerCluster.css';
import 'leaflet.markercluster/dist/MarkerCluster.Default.css';
import 'leaflet.markercluster/dist/leaflet.markercluster';
import '../../styles/marker-cluster.css';
import {
  Box,
  Fab,
  CircularProgress,
  Typography,
  Button,
  useTheme,
  useMediaQuery
} from '@mui/material';

import {
  Navigation,
  Menu as MenuIcon
} from '@mui/icons-material';

import axios from '../../api/axios';
import {
  TILE_LAYER_URL,
  TILE_LAYER_ATTRIBUTION,
  CARTO_VOYAGER_URL,
  CARTO_VOYAGER_ATTRIBUTION
} from '../../components/maps/map-constants';
import '../../components/maps/leaflet-icons';
// Import GIS components
import {
  GISCategoryPanel,
  GISLayerControl,
  GISSearchPanel,
  GISListingCard,
  GISFilterPanel,
  GISResultsDrawer
} from '../../components/gis';

// Define drawer width to match the drawer component
const drawerWidth = 400;

// Гарантированно работающие источники тайлов
const FALLBACK_TILE_URL = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
const FALLBACK_ATTRIBUTION = '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors';

interface Listing {
  id: number;
  title: string;
  price: number;
  description: string;
  condition: string;
  category_id: number;
  status: string;
  latitude: number;
  longitude: number;
  city: string;
  country: string;
  storefront_id?: number;
  storefront_name?: string;
  isStorefront?: boolean;
  storefrontItemCount?: number;
  storefrontName?: string;
  [key: string]: any;
}

interface MapCenter {
  latitude: number;
  longitude: number;
}

interface PriceRange {
  min: number;
  max: number;
}

interface Filters {
  query: string;
  category_id: string;
  price: PriceRange;
  city: string;
  radius: number;
  condition: string;
  sort: string;
  sort_direction: string;
  attributeFilters: Record<string, any>;
  [key: string]: any;
}

interface StorefrontInfo {
  listings: Listing[];
  name: string;
  latitude: number;
  longitude: number;
}

type TileLayerType = 'standard' | 'satellite' | 'terrain' | 'traffic' | 'heatmap' | 'carto';

const GISMapPage: React.FC = () => {
  const { t } = useTranslation(['marketplace', 'common', 'gis']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const { userLocation, detectUserLocation } = useLocation();

  // Состояния для карты
  const mapRef = useRef<L.Map | null>(null);
  const mapContainerRef = useRef<HTMLDivElement | null>(null);
  const markersLayerRef = useRef<L.LayerGroup | null>(null);
  const markerClusterRef = useRef<L.MarkerClusterGroup | null>(null);
  const [mapReady, setMapReady] = useState<boolean>(false);
  const [mapCenter, setMapCenter] = useState<MapCenter | null>(null);
  const [mapZoom, setMapZoom] = useState<number>(13);
  const [mapInitialized, setMapInitialized] = useState<boolean>(false);
  const [clusteringEnabled, setClusteringEnabled] = useState<boolean>(true);

  // Состояния для интерфейса
  const [drawerOpen, setDrawerOpen] = useState<boolean>(!isMobile);
  const [searchDrawerOpen, setSearchDrawerOpen] = useState<boolean>(false);
  const [filterDrawerOpen, setFilterDrawerOpen] = useState<boolean>(false);
  const [layersOpen, setLayersOpen] = useState<boolean>(false);
  const [activeCategoryId, setActiveCategoryId] = useState<number | null>(null);
  const [selectedListingId, setSelectedListingId] = useState<number | null>(null);
  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);

  // Состояния для данных
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [totalListings, setTotalListings] = useState<number>(0);
  const [activeLayers, setActiveLayers] = useState<string[]>(['listings']);

  // Фильтры
  const [filters, setFilters] = useState<Filters>({
    query: searchParams.get('query') || '',
    category_id: searchParams.get('category_id') || '',
    price: {
      min: parseInt(searchParams.get('min_price') || '0', 10),
      max: parseInt(searchParams.get('max_price') || '100000', 10)
    },
    city: searchParams.get('city') || '',
    radius: parseInt(searchParams.get('radius') || '0', 10),
    condition: searchParams.get('condition') || 'any',
    sort: searchParams.get('sort_by') || 'id', // Используем id вместо created_at для совместимости с OpenSearch
    sort_direction: searchParams.get('sort_direction') || 'desc',
    attributeFilters: {}
  });

  // Также сохраняем id избранных объявлений
  const [favoriteListings, setFavoriteListings] = useState<number[]>([]);

  // Определяем начальный центр карты на основе параметров URL или местоположения пользователя
  useEffect(() => {
    const lat = searchParams.get('latitude');
    const lon = searchParams.get('longitude');

    if (lat && lon) {
      setMapCenter({
        latitude: parseFloat(lat),
        longitude: parseFloat(lon)
      });
    } else if (userLocation) {
      setMapCenter({
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
    } else {
      // Позиция по умолчанию - Нови-Сад, Сербия
      setMapCenter({
        latitude: 45.2671,
        longitude: 19.8335
      });
    }

    const zoom = searchParams.get('zoom');
    if (zoom) {
      setMapZoom(parseFloat(zoom));
    }
  }, [searchParams, userLocation]);

  // Супер простая инициализация карты - минимальный надежный вариант
  useEffect(() => {
    // Если нет контейнера или координат, ничего не делаем
    if (!mapContainerRef.current || !mapCenter) return;

    console.log("Инициализация карты...");

    try {
      // Если карта уже существует, удаляем ее
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }

      // Очищаем контейнер
      mapContainerRef.current.innerHTML = '';

      // Создаем карту
      const map = L.map(mapContainerRef.current, {
        center: [mapCenter.latitude, mapCenter.longitude],
        zoom: mapZoom,
        zoomControl: false,        // Отключаем стандартный контрол зума
        preferCanvas: true,        // Используем canvas для лучшей производительности
        zoomAnimation: true,       // Включаем плавную анимацию зума
        fadeAnimation: true,       // Включаем плавную анимацию появления
        markerZoomAnimation: true, // Включаем анимацию маркеров при зуме
        zoomSnap: 1,              // Возвращаем стандартный шаг зума
        zoomDelta: 1,             // Стандартная дельта зума
        wheelPxPerZoomLevel: 60,  // Стандартная чувствительность колесика
        renderer: L.canvas(),     // Принудительно используем canvas renderer
        inertia: true,            // Плавная инерция при перетаскивании
        inertiaDeceleration: 3000, // Быстрая остановка инерции
        inertiaMaxSpeed: 1500,    // Максимальная скорость инерции
        worldCopyJump: false      // Отключаем копирование мира для стабильности
      });

      // Добавляем тайловый слой - используем OpenStreetMap как самый надежный
      L.tileLayer(FALLBACK_TILE_URL, {
        attribution: FALLBACK_ATTRIBUTION,
        maxZoom: 19,
        updateWhenZooming: false,    // НЕ обновлять при зуме для плавности
        keepBuffer: 4,              // Стандартный буфер тайлов
        updateWhenIdle: true,       // Обновлять только когда карта неактивна
        updateInterval: 200         // Интервал обновления в мс
      }).addTo(map);

      // Добавляем слой маркеров
      const markersLayer = L.layerGroup().addTo(map);
      
      // Создаем группу кластеризации маркеров
      const markerCluster = L.markerClusterGroup({
        maxClusterRadius: 50,       // Максимальное расстояние для кластеризации
        disableClusteringAtZoom: 16, // Отключать кластеризацию на зуме >= 16
        spiderfyOnMaxZoom: true,    // Раскрывать кластер при максимальном зуме
        showCoverageOnHover: false, // Не показывать область покрытия
        zoomToBoundsOnClick: true,  // Зумить к границам кластера при клике
        removeOutsideVisibleBounds: false, // Не удалять маркеры за пределами видимости
        animate: true,              // Анимировать кластеризацию
        animateAddingMarkers: false, // Не анимировать добавление отдельных маркеров
        chunkedLoading: true,       // Загружать маркеры порциями для производительности
        iconCreateFunction: function(cluster) {
          const count = cluster.getChildCount();
          let size = 'small';
          
          if (count > 50) {
            size = 'large';
          } else if (count > 10) {
            size = 'medium';
          }
          
          return L.divIcon({
            html: `<div><span>${count}</span></div>`,
            className: `marker-cluster marker-cluster-${size}`,
            iconSize: [40, 40]
          });
        }
      }).addTo(map);

      // Добавляем контрол зума
      L.control.zoom({
        position: 'bottomright'
      }).addTo(map);

      // Сохраняем ссылки
      mapRef.current = map;
      markersLayerRef.current = markersLayer;
      markerClusterRef.current = markerCluster;

      // Устанавливаем обработчики событий только после полной загрузки карты
      map.on('load', function() {
        console.log("Карта полностью загружена (событие load)");
      });

      // Обработчик moveend - срабатывает при изменении положения карты
      map.on('moveend', function() {
        // Дополнительная проверка готовности карты
        if (!map._loaded || !mapRef.current || !mapContainerRef.current) {
          return;
        }

        try {
          const center = map.getCenter();
          const currentZoom = map.getZoom();
          
          if (center && isFinite(center.lat) && isFinite(center.lng)) {
            setMapCenter({
              latitude: center.lat,
              longitude: center.lng
            });
            setMapZoom(currentZoom);

            // Дебаунс для URL и загрузки данных
            if (window.mapUpdateTimeout) {
              clearTimeout(window.mapUpdateTimeout);
            }

            window.mapUpdateTimeout = setTimeout(() => {
              try {
                // Обновляем URL
                setSearchParams(prev => {
                  const newParams = new URLSearchParams(prev);
                  newParams.set('latitude', center.lat.toFixed(6));
                  newParams.set('longitude', center.lng.toFixed(6));
                  newParams.set('zoom', currentZoom.toString());
                  return newParams;
                }, { replace: true });

                // Загружаем данные только если карта не в процессе анимации
                if (!map._animatingZoom && !map._zooming && !map._panAnim && !map.dragging) {
                  fetchListingsInViewport();
                }
              } catch (urlError) {
                console.error("Ошибка при обновлении URL:", urlError);
              }
            }, 300);
          }
        } catch (e) {
          console.error("Ошибка в обработчике moveend:", e);
        }
      });

      // Добавляем отдельный обработчик zoomend для обновления маркеров
      map.on('zoomend', function() {
        if (!map._loaded || !mapRef.current || !mapContainerRef.current) {
          return;
        }

        try {
          // Принудительно обновляем маркеры после завершения зума
          setTimeout(() => {
            if (mapRef.current && mapContainerRef.current && listings.length > 0) {
              console.log("Принудительное обновление маркеров после зума");
              updateMapMarkers(listings);
            }
          }, 100);
        } catch (e) {
          console.error("Ошибка в обработчике zoomend:", e);
        }
      });

      // Добавляем обработчик dragend для обновления маркеров при перетаскивании
      map.on('dragend', function() {
        if (!map._loaded || !mapRef.current || !mapContainerRef.current) {
          return;
        }

        try {
          // Принудительно обновляем маркеры после перетаскивания
          setTimeout(() => {
            if (mapRef.current && mapContainerRef.current && listings.length > 0) {
              console.log("Принудительное обновление маркеров после перетаскивания");
              updateMapMarkers(listings);
            }
          }, 50);
        } catch (e) {
          console.error("Ошибка в обработчике dragend:", e);
        }
      });

      // Используем whenReady для более надежной инициализации
      map.whenReady(() => {
        if (!mapReady) {
          setMapReady(true);
          setMapInitialized(true);
          console.log("Карта готова");
          
          // Немедленно загружаем данные
          fetchListingsInViewport();
        }
      });

      // Обработка ошибок загрузки тайлов
      map.on('tileerror', function(error) {
        console.error("Ошибка загрузки тайла:", error);
      });

      // Резервный таймаут
      setTimeout(() => {
        if (!mapReady && mapRef.current) {
          console.log("Принудительная инициализация карты");
          setMapReady(true);
          setMapInitialized(true);
          fetchListingsInViewport();
        }
      }, 1000);

      console.log("Карта инициализирована успешно");
    } catch (error) {
      console.error("Ошибка при инициализации карты:", error);
      setError("Не удалось загрузить карту");
    }

    // Очистка при размонтировании
    return () => {
      console.log("Размонтирование компонента карты");

      // Очищаем все таймауты
      if (window.mapUpdateTimeout) {
        clearTimeout(window.mapUpdateTimeout);
      }
      if (window.fetchListingsTimeout) {
        clearTimeout(window.fetchListingsTimeout);
      }

      // Удаляем карту
      if (mapRef.current) {
        try {
          mapRef.current.remove();
          mapRef.current = null;
        } catch (e) {
          console.error("Ошибка при удалении карты:", e);
        }
      }

      // Сбрасываем состояние
      setMapReady(false);
      setMapInitialized(false);
    };
  }, [mapCenter, mapZoom]); // Перерисовываем карту при изменении центра или зума

  // Загрузка объявлений в текущей области видимости
  const fetchListingsInViewport = async (): Promise<void> => {
    // Проверка готовности карты
    if (!mapRef.current) {
      console.warn("Попытка загрузить объявления без инициализированной карты");
      return;
    }

    if (!mapContainerRef.current) {
      console.warn("Попытка загрузить объявления без контейнера карты");
      return;
    }

    // Дебаунс для избежания слишком частых вызовов API
    if (window.fetchListingsTimeout) {
      clearTimeout(window.fetchListingsTimeout);
    }

    window.fetchListingsTimeout = setTimeout(async () => {
      try {
        // Включаем индикатор загрузки
        setLoading(true);

        // Проверяем готовность карты
        if (!mapRef.current._loaded && mapReady === false) {
          console.warn("Карта еще не готова для загрузки данных");
          setLoading(false);
          return;
        }

      // Получаем границы видимой области с проверкой ошибок
      let bounds, northEast, southWest;

      try {
        bounds = mapRef.current.getBounds();
        if (!bounds) {
          console.warn("Не удалось получить границы карты");
          setLoading(false);
          return;
        }

        northEast = bounds.getNorthEast();
        southWest = bounds.getSouthWest();

        // Проверка валидности координат
        if (!northEast || !southWest ||
            !isFinite(northEast.lat) || !isFinite(northEast.lng) ||
            !isFinite(southWest.lat) || !isFinite(southWest.lng)) {
          console.warn("Получены невалидные координаты границ карты");
          setLoading(false);
          return;
        }
      } catch (boundsError) {
        console.error("Ошибка при получении границ карты:", boundsError);
        setLoading(false);
        return;
      }

      // Формируем параметры запроса
      const params = {
        view_mode: 'map',
        bbox: `${southWest.lat},${southWest.lng},${northEast.lat},${northEast.lng}`,
        size: 1000, // Большое количество объявлений для карты
        status: 'active',
        query: filters.query || '',
        category_id: activeCategoryId || filters.category_id || '',
        min_price: filters.price?.min || 0,
        max_price: filters.price?.max || 100000,
        condition: filters.condition !== 'any' ? filters.condition : '',
        sort_by: filters.sort || 'id',
        sort_direction: filters.sort_direction || 'desc',
        radius: filters.radius > 0 ? filters.radius : '',
        city: filters.city || ''
      };

      try {
        // Делаем запрос к API с таймаутом
        const response = await axios.get('/api/v1/marketplace/search', {
          params,
          timeout: 15000 // 15-секундный таймаут
        });

        // Проверяем, что компонент всё еще смонтирован
        if (!mapRef.current || !mapContainerRef.current) {
          console.warn("Компонент размонтирован во время загрузки данных");
          return;
        }

        // Обрабатываем ответ
        if (response.data) {
          let fetchedListings = [];

          // Извлекаем данные из ответа с учетом возможных форматов
          if (Array.isArray(response.data.data)) {
            fetchedListings = response.data.data;
          } else if (response.data.data && response.data.data.data) {
            fetchedListings = response.data.data.data;
          } else if (response.data.data) {
            fetchedListings = [response.data.data];
          }

          // Проверяем и фильтруем объявления
          const validListings = fetchedListings.filter(listing =>
              listing &&
              typeof listing.id === 'number' &&
              typeof listing.latitude === 'number' &&
              typeof listing.longitude === 'number' &&
              isFinite(listing.latitude) &&
              isFinite(listing.longitude)
          );

          console.log(`Получено ${validListings.length} валидных объявлений из ${fetchedListings.length} от API`);

          // Обновляем состояние и маркеры
          setListings(validListings);
          setTotalListings(response.data.meta?.total || validListings.length);

          // Обновляем маркеры только если компонент всё еще смонтирован
          if (mapRef.current && mapContainerRef.current && markersLayerRef.current) {
            updateMapMarkers(validListings);
          }
        }
      } catch (apiError) {
        // Обработка ошибки API
        console.error("Ошибка при запросе к API:", apiError);
        setError("Ошибка при загрузке объявлений");

        // Убираем рекурсивный вызов для избежания бесконечного цикла
        if (apiError.code === 'ECONNABORTED' || !apiError.response) {
          console.log("Ошибка сети при загрузке объявлений, повторная попытка отменена");
        }
      }
      } catch (error) {
        console.error("Общая ошибка при загрузке объявлений:", error);
        setError("Ошибка при загрузке объявлений");
      } finally {
        // Выключаем индикатор загрузки
        setLoading(false);
      }
    }, 200); // Дебаунс 200ms
  };

  // Обновление маркеров на карте
  const updateMapMarkers = (listings: Listing[]): void => {
    try {
      // Проверяем все необходимые зависимости
      if (!mapRef.current || !mapContainerRef.current) {
        console.warn("Попытка обновить маркеры без инициализированной карты");
        return;
      }

      // Валидация входных данных
      if (!Array.isArray(listings)) {
        console.warn("updateMapMarkers: получен неверный формат данных listings", listings);
        return;
      }

      // Очищаем существующие маркеры
      if (clusteringEnabled && markerClusterRef.current) {
        markerClusterRef.current.clearLayers();
      } else if (markersLayerRef.current) {
        markersLayerRef.current.clearLayers();
      }

      // Создаем или пересоздаем слои, если нужно
      if (!markersLayerRef.current) {
        try {
          markersLayerRef.current = L.layerGroup().addTo(mapRef.current);
        } catch (error) {
          console.error("Ошибка при создании слоя маркеров:", error);
          return;
        }
      }

      if (clusteringEnabled && !markerClusterRef.current) {
        try {
          markerClusterRef.current = L.markerClusterGroup({
            maxClusterRadius: 50,
            disableClusteringAtZoom: 16,
            spiderfyOnMaxZoom: true,
            showCoverageOnHover: false,
            zoomToBoundsOnClick: true,
            removeOutsideVisibleBounds: false,
            animate: true,
            animateAddingMarkers: false,
            chunkedLoading: true,
            iconCreateFunction: function(cluster) {
              const count = cluster.getChildCount();
              let size = 'small';
              
              if (count > 50) {
                size = 'large';
              } else if (count > 10) {
                size = 'medium';
              }
              
              return L.divIcon({
                html: `<div><span>${count}</span></div>`,
                className: `marker-cluster marker-cluster-${size}`,
                iconSize: [40, 40]
              });
            }
          }).addTo(mapRef.current);
        } catch (error) {
          console.error("Ошибка при создании группы кластеризации:", error);
          return;
        }
      }

      // Используем валидные объявления (только с корректными координатами)
      const validListings = listings.filter(listing =>
          listing &&
          typeof listing.latitude === 'number' &&
          typeof listing.longitude === 'number' &&
          isFinite(listing.latitude) &&
          isFinite(listing.longitude)
      );

      console.log(`Обрабатываем ${validListings.length} валидных объявлений из ${listings.length} для отображения на карте`);

      // Группируем объявления по витринам
      const storefronts = new Map<number, StorefrontInfo>();
      const individualListings: Listing[] = [];

      // Безопасно группируем объявления
      validListings.forEach(listing => {
        if (listing.storefront_id) {
          if (!storefronts.has(listing.storefront_id)) {
            storefronts.set(listing.storefront_id, {
              listings: [],
              name: listing.storefront_name || "Магазин",
              latitude: listing.latitude,
              longitude: listing.longitude
            });
          }
          storefronts.get(listing.storefront_id)!.listings.push(listing);
        } else {
          individualListings.push(listing);
        }
      });

      // Безопасно создаем все маркеры с обработкой ошибок
      const createMarkerSafely = (
          latitude: number,
          longitude: number,
          icon: L.Icon | L.DivIcon,
          clickHandler: () => void
      ): L.Marker | null => {
        try {
          const marker = L.marker([latitude, longitude], { icon });
          marker.on('click', clickHandler);
          return marker;
        } catch (e) {
          console.error("Ошибка при создании маркера:", e);
          return null;
        }
      };

      // Обработка маркеров витрин
      storefronts.forEach((storefront, id) => {
        try {
          if (storefront.listings.length > 0) {
            const clickHandler = () => {
              if (mapRef.current && mapContainerRef.current && storefront.listings.length > 0) {
                setSelectedListing({
                  ...storefront.listings[0],
                  isStorefront: true,
                  storefrontItemCount: storefront.listings.length,
                  storefrontName: storefront.name
                });
                setSelectedListingId(storefront.listings[0].id);
              }
            };

            const storeMarker = createMarkerSafely(
                storefront.latitude,
                storefront.longitude,
                createCustomIcon('store', storefront.listings.length),
                clickHandler
            );

            if (storeMarker) {
              if (clusteringEnabled && markerClusterRef.current) {
                markerClusterRef.current.addLayer(storeMarker);
              } else if (markersLayerRef.current) {
                markersLayerRef.current.addLayer(storeMarker);
              }
            }
          }
        } catch (e) {
          console.error("Ошибка при обработке маркера для витрины:", e);
        }
      });

      // Создаем маркеры для отдельных объявлений пакетами
      const chunkSize = 50; // Оптимальный размер пакета для производительности

      // Обрабатываем маркеры отдельных объявлений асинхронно пакетами
      const processListingsChunk = (startIdx: number) => {
        try {
          // Проверяем, что все системы все еще доступны
          if (!mapRef.current || !markersLayerRef.current || !mapContainerRef.current) {
            return;
          }

          const endIdx = Math.min(startIdx + chunkSize, individualListings.length);
          if (startIdx >= individualListings.length) return;

          const chunk = individualListings.slice(startIdx, endIdx);

          chunk.forEach(listing => {
            const clickHandler = () => {
              if (mapRef.current && mapContainerRef.current) {
                setSelectedListing(listing);
                setSelectedListingId(listing.id);
              }
            };

            const marker = createMarkerSafely(
                listing.latitude,
                listing.longitude,
                createMarkerIcon(listing),
                clickHandler
            );

            if (marker) {
              if (clusteringEnabled && markerClusterRef.current) {
                markerClusterRef.current.addLayer(marker);
              } else if (markersLayerRef.current) {
                markersLayerRef.current.addLayer(marker);
              }
            }
          });

          // Обрабатываем следующий пакет асинхронно, чтобы не блокировать UI
          if (endIdx < individualListings.length) {
            setTimeout(() => processListingsChunk(endIdx), 0);
          }
        } catch (e) {
          console.error("Ошибка при обработке пакета маркеров:", e);
        }
      };

      // Запускаем обработку первого пакета
      if (individualListings.length > 0) {
        processListingsChunk(0);
      }
    } catch (e) {
      console.error("Общая ошибка в updateMapMarkers:", e);
    }
  };

  // Создание иконки маркера в зависимости от типа объявления
  const createMarkerIcon = (listing: Listing): L.DivIcon => {
    let iconType = 'default';

    // Определяем тип иконки на основе категории
    if (listing.category_id) {
      const categoryId = listing.category_id;
      if (categoryId === 1) iconType = 'realestate';
      else if (categoryId === 2) iconType = 'car';
      else if (categoryId === 3) iconType = 'electronics';
      else iconType = 'shopping';
    }

    return createCustomIcon(iconType);
  };

  // Создание настраиваемой иконки маркера
  const createCustomIcon = (type: string, count: number = 0): L.DivIcon => {
    // Определяем цвет и иконку в зависимости от типа
    let color, iconHtml;

    switch (type) {
      case 'store':
        color = '#4CAF50';
        iconHtml = `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
          ${count > 0 ? `<div style="position:absolute;top:-10px;right:-10px;background:#FF5722;border-radius:50%;width:20px;height:20px;display:flex;align-items:center;justify-content:center;font-size:12px;">${count}</div>` : ''}
        </div>`;
        break;
      case 'realestate':
        color = '#2196F3';
        iconHtml = `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
        </div>`;
        break;
      case 'car':
        color = '#FF5722';
        iconHtml = `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="6" width="18" height="12" rx="2"></rect><path d="M6 12h12"></path></svg>
        </div>`;
        break;
      case 'electronics':
        color = '#9C27B0';
        iconHtml = `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="6" y="2" width="12" height="20" rx="2"></rect><line x1="12" y1="18" x2="12" y2="18"></line></svg>
        </div>`;
        break;
      default:
        color = '#FFC107';
        iconHtml = `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10h-8a1 1 0 01-1-1V3L21 10z"></path><path d="M21 18v-8h-8v8z"></path><path d="M3 3v8h8V3"></path><path d="M3 21h8v-8H3z"></path></svg>
        </div>`;
    }

    return L.divIcon({
      className: `custom-marker-${type}`,
      html: `<div style="background-color:${color};width:100%;height:100%;border-radius:50%;display:flex;align-items:center;justify-content:center;color:white;box-shadow:0 2px 5px rgba(0,0,0,0.3);">${iconHtml}</div>`,
      iconSize: [30, 30],
      iconAnchor: [15, 15]
    });
  };

  // Обработчик изменения категории
  const handleCategoryChange = (categoryId: number | null): void => {
    setActiveCategoryId(categoryId);
    setFilters(prev => ({
      ...prev,
      category_id: categoryId ? categoryId.toString() : ''
    }));

    // Обновляем URL
    setSearchParams(prev => {
      const newParams = new URLSearchParams(prev);
      if (categoryId) {
        newParams.set('category_id', categoryId.toString());
      } else {
        newParams.delete('category_id');
      }
      return newParams;
    });

    // Загружаем новые данные
    fetchListingsInViewport();
  };

  // Обработчик изменения фильтров
  const handleFilterChange = (newFilters: Partial<Filters>): void => {
    setFilters(prev => {
      // Преобразуем поле sort если необходимо для совместимости с OpenSearch
      let updated = { ...prev, ...newFilters };

      // Проверяем, если sort установлен в "newest" или "created_at", меняем на "id" для OpenSearch
      if (updated.sort === 'newest' || updated.sort === 'created_at') {
        updated.sort = 'id';
        updated.sort_direction = 'desc';
      }

      // Обновляем URL
      setSearchParams(prev => {
        const newParams = new URLSearchParams(prev);

        Object.entries(updated).forEach(([key, value]) => {
          if (value && key !== 'attributeFilters') {
            if (typeof value === 'object' && 'min' in value && 'max' in value) {
              // Обработка диапазона цен
              newParams.set('min_price', value.min.toString());
              newParams.set('max_price', value.max.toString());
            } else {
              newParams.set(key, value.toString());
            }
          } else if (!value && key !== 'attributeFilters') {
            newParams.delete(key);
          }
        });

        // Добавляем атрибутные фильтры
        if (updated.attributeFilters) {
          Object.entries(updated.attributeFilters).forEach(([attrKey, attrValue]) => {
            if (attrValue) {
              newParams.set(`attr_${attrKey}`, attrValue.toString());
            } else {
              newParams.delete(`attr_${attrKey}`);
            }
          });
        }

        return newParams;
      });

      return updated;
    });

    // Загружаем новые данные
    fetchListingsInViewport();
  };

  // Обработчик получения местоположения пользователя
  const handleDetectLocation = async (): Promise<void> => {
    try {
      await detectUserLocation();

      if (userLocation) {
        // Центрируем карту на местоположении пользователя
        if (mapRef.current) {
          mapRef.current.setView([userLocation.lat, userLocation.lon], 15);
        }

        // Обновляем фильтры с местоположением
        setFilters(prev => ({
          ...prev,
          latitude: userLocation.lat,
          longitude: userLocation.lon,
          city: userLocation.city || '',
          country: userLocation.country || ''
        }));

        // Загружаем новые данные
        fetchListingsInViewport();
      }
    } catch (error) {
      console.error("Error detecting location:", error);
      setError("Не удалось определить местоположение");
    }
  };

  // Обработчик поиска
  const handleSearch = (query: string): void => {
    setFilters(prev => ({
      ...prev,
      query
    }));

    // Обновляем URL
    setSearchParams(prev => {
      const newParams = new URLSearchParams(prev);
      if (query) {
        newParams.set('query', query);
      } else {
        newParams.delete('query');
      }
      return newParams;
    });

    // Загружаем новые данные
    fetchListingsInViewport();

    // Закрываем панель поиска на мобильных устройствах
    if (isMobile) {
      setSearchDrawerOpen(false);
    }
  };

  // Handler for favorite toggle
  const handleFavoriteToggle = (listingId: number, isFavorite: boolean): void => {
    if (isFavorite) {
      setFavoriteListings(prev => [...prev, listingId]);
    } else {
      setFavoriteListings(prev => prev.filter(id => id !== listingId));
    }
  };

  // Handler for show on map
  const handleShowOnMap = (listing: Listing): void => {
    if (listing && listing.latitude && listing.longitude && mapRef.current) {
      mapRef.current.setView([listing.latitude, listing.longitude], 16);
      setSelectedListing(listing);
      setSelectedListingId(listing.id);
    }
  };

  // Handler for contact click
  const handleContactClick = (listing: Listing): void => {
    // Navigate to chat or show contact info
    navigate(`/marketplace/listings/${listing.id}`);
  };

  // Handler for layer change
  const handleLayerChange = (newLayer: TileLayerType): void => {
    // Change map tile layer based on selection
    if (!mapRef.current) return;

    // Remove existing tile layers
    mapRef.current.eachLayer(layer => {
      if (layer instanceof L.TileLayer) {
        mapRef.current!.removeLayer(layer);
      }
    });

    // Общие оптимизированные параметры тайлов для всех типов карт
    const optimizedTileOptions = {
      maxZoom: 19,
      minZoom: 3,
      subdomains: 'abc',          // Поддомены для распределения запросов
      updateWhenZooming: false,    // НЕ обновлять при зуме для плавности
      keepBuffer: 2,              // Стандартный буфер тайлов
      updateWhenIdle: true,       // Обновлять только когда карта неактивна
      updateInterval: 200         // Интервал обновления в мс
    };

    // Add new tile layer based on selection
    switch(newLayer) {
      case 'satellite':
        L.tileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}', {
          ...optimizedTileOptions,
          attribution: 'Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community'
        }).addTo(mapRef.current);
        break;
      case 'terrain':
        L.tileLayer('https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png', {
          ...optimizedTileOptions,
          attribution: 'Map data: &copy; OpenStreetMap contributors, SRTM | Map style: &copy; OpenTopoMap (CC-BY-SA)'
        }).addTo(mapRef.current);
        break;
      case 'traffic':
        // Default map with traffic layer
        L.tileLayer(FALLBACK_TILE_URL, {
          ...optimizedTileOptions,
          attribution: FALLBACK_ATTRIBUTION
        }).addTo(mapRef.current);
        // Add traffic layer if available (example only)
        break;
      case 'heatmap':
        // Default map
        L.tileLayer(FALLBACK_TILE_URL, {
          ...optimizedTileOptions,
          attribution: FALLBACK_ATTRIBUTION
        }).addTo(mapRef.current);
        // Add heatmap layer if data available
        break;
      case 'carto':
        // Использование альтернативного провайдера
        L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
          ...optimizedTileOptions,
          attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        }).addTo(mapRef.current);
        break;
      default: // standard
        L.tileLayer(FALLBACK_TILE_URL, {
          ...optimizedTileOptions,
          attribution: FALLBACK_ATTRIBUTION
        }).addTo(mapRef.current);
    }
  };

  // Handler for clustering toggle
  const handleClusteringToggle = (enableClustering: boolean): void => {
    if (!mapRef.current) return;

    setClusteringEnabled(enableClustering);

    // Переключаем между кластеризованными и обычными маркерами
    if (enableClustering) {
      // Создаем группу кластеризации если её нет
      if (!markerClusterRef.current) {
        markerClusterRef.current = L.markerClusterGroup({
          maxClusterRadius: 50,
          disableClusteringAtZoom: 16,
          spiderfyOnMaxZoom: true,
          showCoverageOnHover: false,
          zoomToBoundsOnClick: true,
          removeOutsideVisibleBounds: false,
          animate: true,
          animateAddingMarkers: false,
          chunkedLoading: true,
          iconCreateFunction: function(cluster) {
            const count = cluster.getChildCount();
            let size = 'small';
            
            if (count > 50) {
              size = 'large';
            } else if (count > 10) {
              size = 'medium';
            }
            
            return L.divIcon({
              html: `<div><span>${count}</span></div>`,
              className: `marker-cluster marker-cluster-${size}`,
              iconSize: [40, 40]
            });
          }
        }).addTo(mapRef.current);
      } else {
        mapRef.current.addLayer(markerClusterRef.current);
      }
      
      // Убираем обычный слой маркеров
      if (markersLayerRef.current && mapRef.current.hasLayer(markersLayerRef.current)) {
        mapRef.current.removeLayer(markersLayerRef.current);
      }
    } else {
      // Убираем кластеризацию
      if (markerClusterRef.current && mapRef.current.hasLayer(markerClusterRef.current)) {
        mapRef.current.removeLayer(markerClusterRef.current);
      }
      
      // Добавляем обычный слой маркеров
      if (markersLayerRef.current) {
        mapRef.current.addLayer(markersLayerRef.current);
      }
    }

    // Обновляем маркеры на карте
    updateMapMarkers(listings);
  };

  // Render the selected listing info
  const renderListingInfo = (): React.ReactNode => {
    if (!selectedListing) return null;

    return (
        <Box
            sx={{
              position: 'absolute',
              bottom: 16,
              left: isMobile ? 16 : (drawerOpen ? 336 : 16),
              right: isMobile ? 16 : 'auto',
              maxWidth: isMobile ? 'auto' : 400,
              transition: 'left 0.3s',
              zIndex: 1000,
            }}
        >
          <GISListingCard
              listing={selectedListing as any}
              isFavorite={favoriteListings.includes(selectedListing.id)}
              onFavoriteToggle={handleFavoriteToggle}
              onShowOnMap={(listing: any) => handleShowOnMap(listing)}
              onContactClick={(listing: any) => handleContactClick(listing)}
              compact={false}
          />
        </Box>
    );
  };

  return (
      <div style={{
        position: 'fixed',
        top: '80px', // Увеличил отступ сверху на 5px
        left: 0,
        right: 0,
        bottom: 0,
        display: 'flex',
        padding: 0,
        margin: 0,
        overflow: 'hidden'
      }}>
        {/* Results Drawer - Putting this first in layout flow */}
        <GISResultsDrawer
            open={drawerOpen}
            onToggleDrawer={() => setDrawerOpen(!drawerOpen)}
            listings={listings as any}
            loading={loading}
            onShowOnMap={(listing: any) => handleShowOnMap(listing)}
            onFilterClick={() => setFilterDrawerOpen(true)}
            onSortClick={() => setFilterDrawerOpen(true)}
            onRefresh={fetchListingsInViewport}
            favoriteListings={favoriteListings}
            onFavoriteToggle={handleFavoriteToggle}
            onContactClick={(listing: any) => handleContactClick(listing)}
            totalCount={totalListings}
            expandToEdge={true} // Signal to expand drawer to full width
        />

        {/* Left Categories Panel */}
        <GISCategoryPanel
            open={searchDrawerOpen}
            onClose={() => setSearchDrawerOpen(false)}
            onCategorySelect={(category: any) => handleCategoryChange(category.id)}
        />

        {/* Right Filter Panel */}
        <GISFilterPanel
            open={filterDrawerOpen}
            onClose={() => setFilterDrawerOpen(false)}
            filters={filters as any}
            onFiltersChange={handleFilterChange}
            onApplyFilters={() => fetchListingsInViewport()}
            onResetFilters={() => {
              handleFilterChange({
                price: { min: 0, max: 100000 },
                condition: 'any',
                radius: 0,
                sort: 'id',  // Используем id вместо created_at для совместимости с OpenSearch
                sort_direction: 'desc'  // Новые объявления имеют бóльшие id
              });
            }}
        />

        {/* Main map container - takes all remaining space */}
        <div style={{
          flexGrow: 1,
          height: '100%',
          width: '100%',
          position: 'relative',
          margin: 0,
          padding: 0,
          overflow: 'hidden',
          zIndex: 5
        }}
        >
          {/* Map container with enhanced styling for reliable rendering */}
          <div
              ref={mapContainerRef}
              style={{
                height: '100%',
                width: '100%',
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                zIndex: 5,
                background: '#f5f5f5', // Фон, видимый во время загрузки карты
                minHeight: '300px', // Минимальная высота для гарантии отрисовки
                border: '1px solid #e0e0e0', // Визуальное обозначение границ карты
                visibility: 'visible', // Гарантия видимости
                display: 'block', // Гарантия отображения
                overflow: 'hidden' // Предотвращение прокрутки внутри контейнера
              }}
              data-testid="map-container" // Атрибут для тестирования
          />

          {/* Search Panel */}
          <GISSearchPanel
              onSearch={handleSearch}
              onLocationSelect={(location) => {
                if (mapRef.current && location) {
                  mapRef.current.setView([location.latitude, location.longitude], 15);
                }
              }}
              drawerOpen={drawerOpen}
              drawerWidth={drawerWidth}
          />

          {/* Layer Control */}
          <GISLayerControl
              layers="standard"
              onLayerChange={handleLayerChange}
              clusterMarkers={clusteringEnabled}
              onClusteringToggle={handleClusteringToggle}
              minimized={!layersOpen}
              onToggleMinimize={() => setLayersOpen(!layersOpen)}
              onClose={() => setLayersOpen(false)}
          />

          {/* Mobile menu button */}
          {isMobile && !drawerOpen && (
              <Fab
                  size="medium"
                  color="primary"
                  aria-label="menu"
                  onClick={() => setDrawerOpen(true)}
                  sx={{
                    position: 'absolute',
                    top: 16,
                    left: 16,
                    zIndex: 1000
                  }}
              >
                <MenuIcon />
              </Fab>
          )}

          {/* My location button */}
          <Fab
              size="medium"
              color="default"
              aria-label="my location"
              onClick={handleDetectLocation}
              sx={{
                position: 'absolute',
                bottom: 90,
                right: 16,
                zIndex: 1000
              }}
          >
            <Navigation />
          </Fab>

          {/* Loading indicator */}
          {loading && (
              <Box
                  sx={{
                    position: 'absolute',
                    top: 80,
                    left: '50%',
                    transform: 'translateX(-50%)',
                    zIndex: 1000,
                    bgcolor: 'background.paper',
                    borderRadius: 8,
                    p: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1,
                    boxShadow: 2,
                  }}
              >
                <CircularProgress size={24} />
                <Typography variant="body2">
                  {t('common:loading')}
                </Typography>
              </Box>
          )}

          {/* Selected listing information */}
          {renderListingInfo()}
        </div>
      </div>
  );
};

// Global Declarations to avoid TypeScript errors with window properties

export default GISMapPage;
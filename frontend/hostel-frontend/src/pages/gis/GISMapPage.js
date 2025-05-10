// frontend/hostel-frontend/src/pages/gis/GISMapPage.js
import React, { useState, useEffect, useRef, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useLocation } from '../../contexts/LocationContext';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
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

const GISMapPage = () => {
  const { t } = useTranslation(['marketplace', 'common', 'gis']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const { userLocation, detectUserLocation } = useLocation();
  
  // Состояния для карты
  const mapRef = useRef(null);
  const mapContainerRef = useRef(null);
  const markersLayerRef = useRef(null);
  const clustersRef = useRef(null);
  const wheelHandlerRef = useRef(null); // Создаем wheelHandlerRef на уровне компонента
  const [mapReady, setMapReady] = useState(false);
  const [mapCenter, setMapCenter] = useState(null);
  const [mapZoom, setMapZoom] = useState(13);
  
  // Состояния для интерфейса
  const [drawerOpen, setDrawerOpen] = useState(!isMobile);
  const [searchDrawerOpen, setSearchDrawerOpen] = useState(false);
  const [filterDrawerOpen, setFilterDrawerOpen] = useState(false);
  const [layersOpen, setLayersOpen] = useState(false);
  const [activeCategoryId, setActiveCategoryId] = useState(null);
  const [selectedListingId, setSelectedListingId] = useState(null);
  const [selectedListing, setSelectedListing] = useState(null);
  
  // Состояния для данных
  const [listings, setListings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [totalListings, setTotalListings] = useState(0);
  const [activeLayers, setActiveLayers] = useState(['listings']);
  
  // Фильтры
  const [filters, setFilters] = useState({
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
  const [favoriteListings, setFavoriteListings] = useState([]);

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
  }, [searchParams, userLocation]);

  // Глобальная настройка для предотвращения ошибки _leaflet_pos
  useEffect(() => {
    // Патч 1: Безопасная обработка _leaflet_pos
    const originalGetPosition = L.DomUtil.getPosition;
    L.DomUtil.getPosition = function(el) {
      if (!el || typeof el !== 'object' || el === null) {
        return new L.Point(0, 0);
      }

      try {
        if (!el._leaflet_pos) {
          return new L.Point(0, 0);
        }
        return originalGetPosition(el);
      } catch (e) {
        return new L.Point(0, 0);
      }
    };

    // Патч 2: Безопасная установка позиции
    const originalSetPosition = L.DomUtil._setPosition;
    L.DomUtil._setPosition = function(el, point) {
      if (!el) return;
      try {
        originalSetPosition(el, point);
      } catch (e) {}
    };

    // Патч 3: Безопасное преобразование координат
    const originalContainerPointToLayerPoint = L.Map.prototype.containerPointToLayerPoint;
    L.Map.prototype.containerPointToLayerPoint = function(point) {
      try {
        return originalContainerPointToLayerPoint.call(this, point);
      } catch (e) {
        return point;
      }
    };

    // Патч 4: Безопасный зум
    const originalPerformZoom = L.Map.prototype._performZoom;
    L.Map.prototype._performZoom = function() {
      try {
        return originalPerformZoom.apply(this, arguments);
      } catch (e) {
        // В случае ошибки отменяем анимацию
        if (this._animatingZoom) {
          this._animatingZoom = false;
        }
        return this;
      }
    };
  }, []); // Пустой массив зависимостей - эффект выполнится только один раз

  // Отключение нативного обработчика масштабирования колесиком не требуется
  // теперь Leaflet полностью патчирован для предотвращения ошибок

  // Инициализация карты
  useEffect(() => {
    // Если нет mapCenter или DOM контейнера, выходим
    if (!mapContainerRef.current || !mapCenter) return;

    // Если карта уже создана, просто обновляем вид и выходим
    if (mapRef.current) {
      try {
        // Обновляем вид безопасно, с отключением анимации для предотвращения ошибок
        mapRef.current.setView([mapCenter.latitude, mapCenter.longitude], mapZoom, {
          animate: false
        });
      } catch (e) {
        console.error("Ошибка при обновлении вида карты:", e);
      }
      return;
    }

    try {
      // Инициализируем карту с безопасными настройками
      mapRef.current = L.map(mapContainerRef.current, {
        // Основные настройки рендера
        preferCanvas: true,            // Canvas для производительности
        renderer: L.canvas(),          // Принудительно используем Canvas

        // Настройки отображения
        zoomControl: false,            // Отключаем встроенные контролы
        attributionControl: false,     // Отключаем копирайты

        // Настройки для стабильности
        zoomSnap: 0.5,                 // Более плавный зум
        wheelPxPerZoomLevel: 120,      // Более грубое изменение зума колесом
        wheelDebounceTime: 100,        // Задержка для обработки вращения колеса

        // Анимации и эффекты
        zoomAnimation: false,         // Отключаем анимацию зума для стабильности
        fadeAnimation: false,         // Отключаем анимации для стабильности
        markerZoomAnimation: false,   // Отключаем анимацию маркеров
        inertia: false,               // Отключаем инерцию для стабильности

        // Другие настройки
        worldCopyJump: true,          // Перемещение через границу
        tap: false,                   // Отключаем tap для мобильных
        trackResize: true,            // Отслеживание изменения размеров
        maxZoom: 19,                  // Максимальный зум
        minZoom: 3,                   // Минимальный зум
        doubleClickZoom: true,        // Включаем стандартный зум по двойному клику
        scrollWheelZoom: true,        // Включаем стандартный зум колесиком
        dragging: true,                // Оставляем возможность перетаскивания
        keyboard: false                // Отключаем зум с клавиатуры
      });

      // Устанавливаем начальный вид безопасно, без анимации
      mapRef.current.setView([mapCenter.latitude, mapCenter.longitude], mapZoom, {
        animate: false
      });

      // Добавляем базовый слой
      L.tileLayer(TILE_LAYER_URL, {
        attribution: TILE_LAYER_ATTRIBUTION,
        maxZoom: 19,
        subdomains: 'abc',
        updateWhenZooming: false,   // Отключаем обновление при зуме
        keepBuffer: 4               // Буфер тайлов
      }).addTo(mapRef.current);

      // Полагаемся на стандартное поведение scrollWheelZoom Leaflet

      // Используем стандартные обработчики зума с нашими патчами для Leaflet
      // которые должны предотвратить ошибки _leaflet_pos

      // Добавляем стандартный контрол зума без переопределений
      L.control.zoom({
        position: 'bottomright',
        zoomInTitle: 'Приблизить',
        zoomOutTitle: 'Отдалить'
      }).addTo(mapRef.current);

      /* Удаляем переопределение контрола зума, используем стандартное поведение
      const originalZoomIn = zoomControl.getContainer;
      zoomControl.getContainer = function() {
        const container = originalZoomIn.call(this);
        if (container) {
          // Находим кнопки зума
          const zoomInButton = container.querySelector('.leaflet-control-zoom-in');
          const zoomOutButton = container.querySelector('.leaflet-control-zoom-out');

          if (zoomInButton) {
            // Очищаем существующие обработчики
            zoomInButton.parentNode.replaceChild(zoomInButton.cloneNode(true), zoomInButton);
            const newZoomInButton = container.querySelector('.leaflet-control-zoom-in');

            // Добавляем новый безопасный обработчик
            newZoomInButton.addEventListener('click', (e) => {
              e.preventDefault();
              e.stopPropagation();

              try {
                // Безопасное увеличение зума
                if (mapRef.current) {
                  const currentZoom = mapRef.current.getZoom();
                  mapRef.current.setZoom(currentZoom + 1, { animate: false });
                }
              } catch (error) {
                console.warn("Ошибка при увеличении зума:", error);
              }
            });
          }

          if (zoomOutButton) {
            // Очищаем существующие обработчики
            zoomOutButton.parentNode.replaceChild(zoomOutButton.cloneNode(true), zoomOutButton);
            const newZoomOutButton = container.querySelector('.leaflet-control-zoom-out');

            // Добавляем новый безопасный обработчик
            newZoomOutButton.addEventListener('click', (e) => {
              e.preventDefault();
              e.stopPropagation();

              try {
                // Безопасное уменьшение зума
                if (mapRef.current) {
                  const currentZoom = mapRef.current.getZoom();
                  mapRef.current.setZoom(currentZoom - 1, { animate: false });
                }
              } catch (error) {
                console.warn("Ошибка при уменьшении зума:", error);
              }
            });
          }
        }
        return container;
      };

      // Добавляем контрол на карту
      */

      // Создаем слой маркеров
      markersLayerRef.current = L.layerGroup().addTo(mapRef.current);

      // Обновляем размеры для уверенности
      setTimeout(() => {
        if (mapRef.current) {
          mapRef.current.invalidateSize({ animate: false });
        }
      }, 100);

      // Подписываемся на события с проверкой состояния карты
      const safeEventHandler = (handler) => {
        return (e) => {
          if (mapRef.current && mapContainerRef.current) {
            handler(e);
          }
        };
      };

      mapRef.current.on('moveend', safeEventHandler(handleMapMoveEnd));
      mapRef.current.on('zoomend', safeEventHandler(handleMapZoomEnd));

      // Отмечаем, что карта готова
      setMapReady(true);

      // Загружаем данные
      setTimeout(() => {
        if (mapRef.current) {
          fetchListingsInViewport();
        }
      }, 500);
    } catch (error) {
      console.error("Ошибка при инициализации карты:", error);
    }

    // wheelHandlerRef уже объявлен на уровне компонента

    // Очистка при размонтировании
    return () => {
      // Используем стандартные обработчики Leaflet, так что очистка не требуется

      // Удаляем карту
      if (mapRef.current) {
        try {
          mapRef.current.off('moveend', handleMapMoveEnd);
          mapRef.current.off('zoomend', handleMapZoomEnd);
          mapRef.current.remove();
        } catch (error) {
          console.error("Ошибка при удалении карты:", error);
        }
        mapRef.current = null;
      }
      setMapReady(false);
    };
  }, [mapCenter]); // Зависимость от mapCenter с проверкой mapRef.current внутри

  // Обработчик перемещения карты
  const handleMapMoveEnd = () => {
    try {
      if (!mapRef.current || !mapContainerRef.current) return;

      const center = mapRef.current.getCenter();
      const zoom = mapRef.current.getZoom();

      // Обновляем mapCenter только при существенном изменении
      if (
        !mapCenter ||
        Math.abs(center.lat - mapCenter.latitude) > 0.001 ||
        Math.abs(center.lng - mapCenter.longitude) > 0.001
      ) {
        setMapCenter({
          latitude: center.lat,
          longitude: center.lng
        });
      }

      // Обновляем URL с задержкой, чтобы не мешать плавности
      if (window.updateUrlTimeout) {
        clearTimeout(window.updateUrlTimeout);
      }

      window.updateUrlTimeout = setTimeout(() => {
        try {
          if (mapRef.current) {
            setSearchParams(prev => {
              const newParams = new URLSearchParams(prev);
              const currentCenter = mapRef.current.getCenter();
              const currentZoom = mapRef.current.getZoom();
              newParams.set('latitude', currentCenter.lat.toFixed(6));
              newParams.set('longitude', currentCenter.lng.toFixed(6));
              newParams.set('zoom', currentZoom);
              return newParams;
            }, { replace: true });
          }
        } catch (e) {
          console.error("Ошибка при обновлении URL:", e);
        }
      }, 300);

      // Загружаем новые данные с небольшой задержкой
      if (window.fetchListingsTimeout) {
        clearTimeout(window.fetchListingsTimeout);
      }

      window.fetchListingsTimeout = setTimeout(() => {
        try {
          if (mapRef.current && mapReady) {
            fetchListingsInViewport();
          }
        } catch (e) {
          console.error("Ошибка при загрузке объявлений:", e);
        }
      }, 400);
    } catch (e) {
      console.error("Ошибка в обработчике moveend:", e);
    }
  };

  // Обработчик масштабирования карты
  const handleMapZoomEnd = () => {
    try {
      if (!mapRef.current || !mapContainerRef.current) return;

      // Получаем текущий зум
      const newZoom = mapRef.current.getZoom();

      // Обновляем состояние зума только при реальном изменении
      if (Math.abs(newZoom - mapZoom) > 0.1) {
        setMapZoom(newZoom);
      }

      // Обновляем URL с задержкой
      if (window.updateUrlTimeout) {
        clearTimeout(window.updateUrlTimeout);
      }

      window.updateUrlTimeout = setTimeout(() => {
        try {
          if (mapRef.current) {
            setSearchParams(prev => {
              const newParams = new URLSearchParams(prev);
              const currentCenter = mapRef.current.getCenter();
              const currentZoom = mapRef.current.getZoom();
              newParams.set('latitude', currentCenter.lat.toFixed(6));
              newParams.set('longitude', currentCenter.lng.toFixed(6));
              newParams.set('zoom', currentZoom);
              return newParams;
            }, { replace: true });
          }
        } catch (e) {
          console.error("Ошибка при обновлении URL после зума:", e);
        }
      }, 300);

      // Загружаем новые данные с небольшой задержкой
      if (window.fetchListingsTimeout) {
        clearTimeout(window.fetchListingsTimeout);
      }

      window.fetchListingsTimeout = setTimeout(() => {
        try {
          if (mapRef.current && mapReady) {
            fetchListingsInViewport();
          }
        } catch (e) {
          console.error("Ошибка при загрузке объявлений после зума:", e);
        }
      }, 400);
    } catch (e) {
      console.error("Ошибка в обработчике zoomend:", e);
    }
  };

  // Загрузка объявлений в текущей области видимости
  const fetchListingsInViewport = async () => {
    if (!mapRef.current || !mapReady || !mapContainerRef.current) return;

    try {
      setLoading(true);

      // Получаем границы видимой области
      try {
        const bounds = mapRef.current.getBounds();
        const northEast = bounds.getNorthEast();
        const southWest = bounds.getSouthWest();

        // Проверка валидности координат
        if (!isFinite(northEast.lat) || !isFinite(northEast.lng) ||
            !isFinite(southWest.lat) || !isFinite(southWest.lng)) {
          console.warn("Получены невалидные координаты границ карты.");
          setLoading(false);
          return;
        }

        // Формируем параметры запроса
        const params = {
          view_mode: 'map',
          bbox: `${southWest.lat},${southWest.lng},${northEast.lat},${northEast.lng}`,
          size: 1000, // Запрашиваем большое количество объявлений для карты
          status: 'active', // Только активные объявления
          query: filters.query,
          category_id: activeCategoryId || filters.category_id,
          min_price: filters.price.min,
          max_price: filters.price.max,
          condition: filters.condition !== 'any' ? filters.condition : '',
          sort_by: filters.sort, // Используем значение из фильтров (которое уже преобразовано для OpenSearch)
          sort_direction: filters.sort_direction || 'desc', // Используем направление из фильтров или desc по умолчанию
          radius: filters.radius > 0 ? filters.radius : '',
          city: filters.city
        };

        // Делаем запрос к API
        const response = await axios.get('/api/v1/marketplace/search', { params });

        if (response.data && response.data.data) {
          const fetchedListings = Array.isArray(response.data.data)
            ? response.data.data
            : response.data.data.data || [];

          // Проверяем, что карта и компонент все еще существуют
          if (mapRef.current && mapContainerRef.current) {
            // Логируем количество полученных объявлений для отладки
            console.log(`Получено ${fetchedListings.length} объявлений от API`);

            setListings(fetchedListings);
            setTotalListings(response.data.meta?.total || fetchedListings.length);

            // Обновляем маркеры на карте
            updateMapMarkers(fetchedListings);
          }
        }
      } catch (boundsError) {
        console.error("Ошибка при получении границ карты:", boundsError);
      }
    } catch (error) {
      console.error("Ошибка при загрузке объявлений:", error);
      setError("Ошибка при загрузке объявлений");
    } finally {
      setLoading(false);
    }
  };

  // Обновление маркеров на карте - быстрый способ без задержек и анимаций
  const updateMapMarkers = (listings) => {
    try {
      if (!mapRef.current || !markersLayerRef.current || !mapContainerRef.current) return;

      // Очищаем текущие маркеры безопасно
      try {
        markersLayerRef.current.clearLayers();
      } catch (e) {
        console.error("Ошибка при очистке слоя маркеров:", e);
        // В случае ошибки пробуем пересоздать слой маркеров
        if (mapRef.current) {
          markersLayerRef.current = L.layerGroup().addTo(mapRef.current);
        } else {
          return; // Если нет карты, просто выходим
        }
      }

      // Проверка входных данных
      if (!Array.isArray(listings)) {
        console.warn("updateMapMarkers: получен неверный формат данных listings", listings);
        return;
      }

      // Логируем количество объявлений для обработки
      console.log(`Обрабатываем ${listings.length} объявлений для отображения на карте`);

      // Группируем объявления по витринам
      const storefronts = new Map();
      const individualListings = [];

      listings.forEach(listing => {
        if (!listing || !listing.latitude || !listing.longitude ||
            !isFinite(listing.latitude) || !isFinite(listing.longitude)) return;

        if (listing.storefront_id) {
          if (!storefronts.has(listing.storefront_id)) {
            storefronts.set(listing.storefront_id, {
              listings: [],
              name: listing.storefront_name || "Магазин",
              latitude: listing.latitude,
              longitude: listing.longitude
            });
          }
          storefronts.get(listing.storefront_id).listings.push(listing);
        } else {
          individualListings.push(listing);
        }
      });

      // Добавляем все маркеры непосредственно, без отложенной анимации
      // Создаем маркеры для витрин
      storefronts.forEach((storefront, id) => {
        try {
          if (storefront.listings.length > 0) {
            const storeMarker = L.marker([storefront.latitude, storefront.longitude], {
              icon: createCustomIcon('store', storefront.listings.length)
            });

            storeMarker.on('click', () => {
              if (mapRef.current && mapContainerRef.current) {
                setSelectedListing({
                  ...storefront.listings[0],
                  isStorefront: true,
                  storefrontItemCount: storefront.listings.length,
                  storefrontName: storefront.name
                });
                setSelectedListingId(storefront.listings[0].id);
              }
            });

            if (markersLayerRef.current) {
              markersLayerRef.current.addLayer(storeMarker);
            }
          }
        } catch (e) {
          console.error("Ошибка при создании маркера для витрины:", e);
        }
      });

      // Создаем маркеры для отдельных объявлений - пакетами, чтобы не блокировать рендеринг
      const chunkSize = 50; // Обрабатываем по 50 маркеров за раз

      for (let i = 0; i < individualListings.length; i += chunkSize) {
        const chunk = individualListings.slice(i, i + chunkSize);

        chunk.forEach(listing => {
          try {
            const marker = L.marker([listing.latitude, listing.longitude], {
              icon: createMarkerIcon(listing)
            });

            marker.on('click', () => {
              if (mapRef.current && mapContainerRef.current) {
                setSelectedListing(listing);
                setSelectedListingId(listing.id);
              }
            });

            if (markersLayerRef.current) {
              markersLayerRef.current.addLayer(marker);
            }
          } catch (e) {
            console.error("Ошибка при создании маркера для объявления:", e);
          }
        });
      }
    } catch (e) {
      console.error("Общая ошибка в updateMapMarkers:", e);
    }
  };

  // Создание иконки маркера в зависимости от типа объявления
  const createMarkerIcon = (listing) => {
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
  const createCustomIcon = (type, count = 0) => {
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
  const handleCategoryChange = (categoryId) => {
    setActiveCategoryId(categoryId);
    setFilters(prev => ({
      ...prev,
      category_id: categoryId
    }));
    
    // Обновляем URL
    setSearchParams(prev => {
      const newParams = new URLSearchParams(prev);
      if (categoryId) {
        newParams.set('category_id', categoryId);
      } else {
        newParams.delete('category_id');
      }
      return newParams;
    });
    
    // Загружаем новые данные
    fetchListingsInViewport();
  };

  // Обработчик изменения фильтров
  const handleFilterChange = (newFilters) => {
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
            newParams.set(key, value);
          } else if (!value && key !== 'attributeFilters') {
            newParams.delete(key);
          }
        });

        // Добавляем атрибутные фильтры
        if (updated.attributeFilters) {
          Object.entries(updated.attributeFilters).forEach(([attrKey, attrValue]) => {
            if (attrValue) {
              newParams.set(`attr_${attrKey}`, attrValue);
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
  const handleDetectLocation = async () => {
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
  const handleSearch = (query) => {
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
  const handleFavoriteToggle = (listingId, isFavorite) => {
    if (isFavorite) {
      setFavoriteListings(prev => [...prev, listingId]);
    } else {
      setFavoriteListings(prev => prev.filter(id => id !== listingId));
    }
  };

  // Handler for show on map
  const handleShowOnMap = (listing) => {
    if (listing && listing.latitude && listing.longitude && mapRef.current) {
      mapRef.current.setView([listing.latitude, listing.longitude], 16);
      setSelectedListing(listing);
      setSelectedListingId(listing.id);
    }
  };

  // Handler for contact click
  const handleContactClick = (listing) => {
    // Navigate to chat or show contact info
    navigate(`/marketplace/listings/${listing.id}`);
  };

  // Handler for layer change
  const handleLayerChange = (newLayer) => {
    // Change map tile layer based on selection
    if (!mapRef.current) return;
    
    // Remove existing tile layers
    mapRef.current.eachLayer(layer => {
      if (layer instanceof L.TileLayer) {
        mapRef.current.removeLayer(layer);
      }
    });
    
    // Общие оптимизированные параметры тайлов для всех типов карт
    const optimizedTileOptions = {
      maxZoom: 19,
      minZoom: 3,
      subdomains: 'abc',          // Поддомены для распределения запросов
      updateWhenZooming: true,    // Обновление при зуме
      keepBuffer: 4               // Буфер тайлов
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
        L.tileLayer(TILE_LAYER_URL, {
          ...optimizedTileOptions,
          attribution: TILE_LAYER_ATTRIBUTION
        }).addTo(mapRef.current);
        // Add traffic layer if available (example only)
        break;
      case 'heatmap':
        // Default map
        L.tileLayer(TILE_LAYER_URL, {
          ...optimizedTileOptions,
          attribution: TILE_LAYER_ATTRIBUTION
        }).addTo(mapRef.current);
        // Add heatmap layer if data available
        break;
      case 'carto':
        // Использование альтернативного провайдера - CartoDB Voyager для более быстрой загрузки
        L.tileLayer(CARTO_VOYAGER_URL, {
          ...optimizedTileOptions,
          attribution: CARTO_VOYAGER_ATTRIBUTION
        }).addTo(mapRef.current);
        break;
      default: // standard
        L.tileLayer(TILE_LAYER_URL, {
          ...optimizedTileOptions,
          attribution: TILE_LAYER_ATTRIBUTION
        }).addTo(mapRef.current);
    }
  };

  // Handler for clustering toggle
  const handleClusteringToggle = (enableClustering) => {
    if (!mapRef.current || !markersLayerRef.current) return;
    
    // Implement marker clustering based on the enableClustering flag
    if (enableClustering) {
      // Convert existing markers to a cluster group
      if (!clustersRef.current) {
        // Create cluster group if it doesn't exist
        // This would use Leaflet.markercluster in a real implementation
        console.log("Clustering enabled - would implement MarkerClusterGroup");
      }
    } else {
      // Remove clustering and show individual markers
      console.log("Clustering disabled - would remove MarkerClusterGroup");
    }
    
    // Refresh markers on the map
    updateMapMarkers(listings);
  };

  // Render the selected listing info
  const renderListingInfo = () => {
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
          listing={selectedListing}
          isFavorite={favoriteListings.includes(selectedListing.id)}
          onFavoriteToggle={handleFavoriteToggle}
          onShowOnMap={handleShowOnMap}
          onContactClick={handleContactClick}
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
        listings={listings}
        loading={loading}
        onShowOnMap={handleShowOnMap}
        onFilterClick={() => setFilterDrawerOpen(true)}
        onSortClick={() => setFilterDrawerOpen(true)}
        onRefresh={fetchListingsInViewport}
        favoriteListings={favoriteListings}
        onFavoriteToggle={handleFavoriteToggle}
        onContactClick={handleContactClick}
        totalCount={totalListings}
        expandToEdge={true} // Signal to expand drawer to full width
      />
      
      {/* Left Categories Panel */}
      <GISCategoryPanel
        open={searchDrawerOpen}
        onClose={() => setSearchDrawerOpen(false)}
        onCategorySelect={handleCategoryChange}
      />
      
      {/* Right Filter Panel */}
      <GISFilterPanel
        open={filterDrawerOpen}
        onClose={() => setFilterDrawerOpen(false)}
        filters={filters}
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
        {/* Map */}
        <div
          ref={mapContainerRef}
          style={{
            height: '100%',
            width: '100%',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0
          }}
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
          clusterMarkers={false}
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

export default GISMapPage;
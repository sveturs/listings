// frontend/hostel-frontend/src/components/marketplace/MapView.js
import React, { useState, useEffect, useRef, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Navigation, X, List, Maximize2, Store } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import {
  Box,
  Paper,
  Typography,
  Chip,
  Button,
  Card,
  CardContent,
  CardMedia,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  useTheme,
  useMediaQuery,
  Collapse,
  IconButton,
  Modal
} from '@mui/material';
import { TILE_LAYER_URL, TILE_LAYER_ATTRIBUTION } from '../maps/map-constants';
import '../maps/leaflet-icons'; // Для исправления иконок Leaflet
import FullscreenMap from '../maps/FullscreenMap';
import { useLocation } from '../../contexts/LocationContext';

// Компонент для предпросмотра объявления при клике по маркеру
const ListingPreview = ({ listing, onClose, onNavigate }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const navigate = useNavigate();

  if (!listing) return null;

  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };

  // Нормализация пути к изображению
  const getImageUrl = (images) => {
    if (!images || !Array.isArray(images) || images.length === 0) {
      return null;
    }

    // Сначала пытаемся найти главное изображение
    let mainImage = images.find(img => img && img.is_main === true);

    // Если главное не найдено, берем первое изображение
    if (!mainImage) {
      mainImage = images[0];
    }

    // Если изображение - это строка с путем
    if (typeof mainImage === 'string') {
      const baseUrl = process.env.REACT_APP_BACKEND_URL || '';
      return `${baseUrl}/uploads/${mainImage}`;
    }

    // Если изображение - это объект с file_path
    if (mainImage && typeof mainImage === 'object' && mainImage.file_path) {
      const baseUrl = process.env.REACT_APP_BACKEND_URL || '';
      return `${baseUrl}/uploads/${mainImage.file_path}`;
    }

    return null;
  };

  const imageUrl = getImageUrl(listing.images);

  // Определяем, показываем ли мы карточку витрины или объявления
  const isStorefrontCard = Boolean(listing.isPartOfStorefront ||
    (listing.storefront_id && listing.isUniqueLocation));

  console.log('ListingPreview: определение типа карточки:', {
    isStorefrontCard: isStorefrontCard,
    isPartOfStorefront: listing.isPartOfStorefront,
    hasStorefrontId: Boolean(listing.storefront_id),
    isUnique: listing.isUniqueLocation
  })
  return (
    <Card
      sx={{
        position: 'absolute',
        bottom: isMobile ? 0 : 16,
        left: isMobile ? 0 : 16,
        maxWidth: isMobile ? '100%' : 400,
        width: isMobile ? '100%' : 'auto',
        zIndex: 1000,
        borderRadius: isMobile ? '8px 8px 0 0' : 1,
      }}
    >
      <Box sx={{ position: 'relative' }}>
        <IconButton
          onClick={onClose}
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
            bgcolor: 'background.paper',
            opacity: 0.8,
            '&:hover': { bgcolor: 'background.paper', opacity: 1 },
            zIndex: 10
          }}
        >
          <X size={16} />
        </IconButton>

        {isStorefrontCard && (
          <Box
            sx={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              backgroundColor: theme.palette.primary.main,
              color: 'white',
              padding: '4px 12px',
              display: 'flex',
              alignItems: 'center',
              zIndex: 5
            }}
          >
            <Store size={16} style={{ marginRight: 6 }} />
            <Typography variant="subtitle2" fontWeight="bold">
              {listing.storefrontName || t('common:map.storefront')}
              {listing.isUniqueLocation && (
                <Chip
                  size="small"
                  label={t('common:map.uniqueLocation', { defaultValue: 'Уникальный адрес' })}
                  sx={{ ml: 1, height: 18, fontSize: '0.65rem' }}
                  color="secondary"
                />
              )}
            </Typography>
          </Box>
        )}

        {imageUrl && (
          <CardMedia
            component="img"
            height={isStorefrontCard ? 160 : 140}
            image={imageUrl}
            alt={listing.title}
            sx={{
              pt: isStorefrontCard ? '24px' : 0
            }}
          />
        )}

        <CardContent>
          {isStorefrontCard && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="text.secondary">
                {t('common:map.items', { count: listing.storefrontItemCount })}
              </Typography>
              <Typography variant="h6" fontWeight="bold" sx={{ mt: 0.5 }}>
                {listing.title}
              </Typography>
            </Box>
          )}

          {!isStorefrontCard && (
            <Typography variant="subtitle1" noWrap gutterBottom>
              {listing.title}
            </Typography>
          )}

          <Typography variant="h6" color="primary" gutterBottom>
            {formatPrice(listing.price)}
          </Typography>

          <Box display="flex" justifyContent="space-between" alignItems="center" mt={1}>
            <Chip
              label={listing.condition === 'new' ? t('listings.conditions.new') : t('listings.conditions.used')}
              size="small"
              color={listing.condition === 'new' ? 'success' : 'default'}
            />

            {isStorefrontCard ? (
              <Button
                variant="contained"
                color="primary"
                size="small"
                fullWidth
                sx={{ ml: 1 }}
                onClick={() => {
                  console.log(`Переход на страницу витрины: /shop/${listing.storefront_id}?highlightedListingId=${listing.id}`, {
                    storefront_id: listing.storefront_id,
                    listing_id: listing.id
                  });
                  navigate(`/shop/${listing.storefront_id}?highlightedListingId=${listing.id}`);
                }}
              >
                {listing.isUniqueLocation
                  ? t('common:map.viewItemInStorefront', { defaultValue: 'Открыть в витрине' })
                  : t('common:map.viewStorefront')}
              </Button>
            ) : (
              <Button
                variant="contained"
                size="small"
                onClick={() => {
                  console.log(`Переход на страницу объявления: ${listing.id}`);
                  onNavigate(listing.id);
                }}
              >
                {t('listings.details.viewDetails')}
              </Button>
            )}
          </Box>
        </CardContent>
      </Box>
    </Card>
  );
};

const MapView = ({ listings, filters, onFilterChange, onMapClose }) => {
  const { t } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const navigate = useNavigate();
  const { userLocation, detectUserLocation } = useLocation();

  // Создаем правильный объект с координатами из userLocation
  const locationCoordinates = userLocation ? {
    latitude: userLocation.lat,
    longitude: userLocation.lon
  } : null;

  // Логируем для отладки
  console.log("MapView userLocation:", userLocation);
  console.log("MapView locationCoordinates:", locationCoordinates);

  const mapRef = useRef(null);
  const markersRef = useRef([]);
  const mapContainerRef = useRef(null);
  const [selectedListing, setSelectedListing] = useState(null);
  const [mapReady, setMapReady] = useState(false);
  useEffect(() => {
    if (selectedListing) {
      console.log("Состояние selectedListing изменилось:", {
        id: selectedListing.id,
        isPartOfStorefront: selectedListing.isPartOfStorefront,
        storefront_id: selectedListing.storefront_id,
        isUniqueLocation: selectedListing.isUniqueLocation
      });
    }
  }, [selectedListing]);
  // Состояние для модального окна с полноэкранной картой
  const [expandedMapOpen, setExpandedMapOpen] = useState(false);
  // Центр для полноэкранной карты
  const [expandedMapCenter, setExpandedMapCenter] = useState(null);
  // Маркеры для полноэкранной карты
  const [expandedMapMarkers, setExpandedMapMarkers] = useState([]);

  // Убедимся, что у нас всегда есть userLocationState
  const [userLocationState, setUserLocationState] = useState(locationCoordinates);

  // Обновляем userLocationState при изменении userLocation
  useEffect(() => {
    if (userLocation) {
      setUserLocationState({
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
      console.log("Обновлен userLocationState из userLocation:", {
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
    }
  }, [userLocation]);

  // Инициализация карты
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    // Добавим задержку, чтобы убедиться, что DOM готов
    const timer = setTimeout(() => {
      try {
        // Устанавливаем центр карты на основе данных из контекста
        const initialPosition = userLocation
          ? [userLocation.lat, userLocation.lon]  // Используем данные из контекста
          : [45.2671, 19.8335]; // Координаты Нови-Сада по умолчанию

        // Проверка, что контейнер существует и имеет размеры
        if (!mapContainerRef.current ||
          mapContainerRef.current.clientWidth === 0 ||
          mapContainerRef.current.clientHeight === 0) {
          console.log("Map container not ready yet, retrying...");
          return;
        }

        // Создаем карту с опцией preferCanvas для лучшей производительности
        mapRef.current = L.map(mapContainerRef.current, {
          preferCanvas: true,
          attributionControl: false,
          zoomControl: true,
          inertia: true,
          fadeAnimation: true,
          zoomAnimation: true,
          renderer: L.canvas()
        }).setView(initialPosition, 13);

        // Добавляем слой тайлов
        L.tileLayer(TILE_LAYER_URL, {
          attribution: TILE_LAYER_ATTRIBUTION,
          maxZoom: 19
        }).addTo(mapRef.current);

        // Если есть местоположение пользователя, добавляем маркер и круг
        if (userLocation) {
          try {
            L.circle(initialPosition, {
              color: theme.palette.primary.main,
              fillColor: theme.palette.primary.light,
              fillOpacity: 0.2,
              radius: getRadiusInMeters(filters.distance || '5km')
            }).addTo(mapRef.current);

            L.marker(initialPosition, {
              icon: L.divIcon({
                html: `<div style="
                  background-color: ${theme.palette.primary.main};
                  width: 16px;
                  height: 16px;
                  border-radius: 50%;
                  border: 2px solid white;
                  box-shadow: 0 0 4px rgba(0,0,0,0.3);
                "></div>`,
                className: 'my-location-marker',
                iconSize: [20, 20],
                iconAnchor: [10, 10]
              })
            }).addTo(mapRef.current)
              .bindTooltip(t('listings.map.yourLocation'), { permanent: false });
          } catch (innerError) {
            console.error("Error adding user marker:", innerError);
          }
        }

        // Обработчик для исправления проблемы _leaflet_pos
        mapRef.current.on('zoomanim', (e) => {
          // Ничего не делаем, но это помогает предотвратить ошибку
        });

        setMapReady(true);
      } catch (error) {
        console.error("Error initializing map:", error);
      }
    }, 300); // Увеличиваем задержку для гарантии готовности DOM

    return () => {
      clearTimeout(timer);
      if (mapRef.current) {
        try {
          // Явно удаляем все слои перед удалением карты
          mapRef.current.eachLayer((layer) => {
            mapRef.current.removeLayer(layer);
          });
          mapRef.current.remove();
        } catch (error) {
          console.error("Error removing map:", error);
        }
        mapRef.current = null;
      }
      setMapReady(false); // Сбрасываем состояние готовности карты
    };
  }, [userLocation, theme, t, filters.distance]);

  // Функция для преобразования расстояния (например, "5km") в метры
  const getRadiusInMeters = (distanceString) => {
    if (!distanceString) return 5000; // По умолчанию 5 км

    const match = distanceString.match(/^(\d+)km$/);
    if (match) {
      return parseInt(match[1]) * 1000;
    }
    return 5000;
  };

  // Обновляем круг радиуса при изменении фильтра расстояния
  useEffect(() => {
    if (!mapRef.current || !userLocation || !filters.distance || !mapReady) return;

    try {
      // Удаляем старые круги
      mapRef.current.eachLayer(layer => {
        if (layer instanceof L.Circle) {
          mapRef.current.removeLayer(layer);
        }
      });

      // Добавляем новый круг с актуальным радиусом
      const radiusInMeters = getRadiusInMeters(filters.distance);
      L.circle([userLocation.lat, userLocation.lon], {
        color: theme.palette.primary.main,
        fillColor: theme.palette.primary.light,
        fillOpacity: 0.2,
        radius: radiusInMeters
      }).addTo(mapRef.current);
    } catch (error) {
      console.error("Error updating radius circle:", error);
    }
  }, [filters.distance, userLocation, theme, mapReady]);

  // Обновляем маркеры объявлений
  useEffect(() => {
    if (!mapRef.current || !mapReady) return;

    try {
      // Удаляем старые маркеры
      markersRef.current.forEach(marker => {
        try {
          mapRef.current.removeLayer(marker);
        } catch (error) {
          console.error("Error removing marker:", error);
        }
      });
      markersRef.current = [];

      // Проверяем наличие объявлений с координатами
      const validListings = listings.filter(listing =>
        listing.latitude && listing.longitude &&
        listing.show_on_map !== false
      );

      if (validListings.length === 0) return;

      // Создаем группу маркеров для автомасштабирования
      const markerGroup = L.featureGroup();

      // Группировка объявлений по витринам и определение уникальных местоположений
      const storefrontListings = new Map(); // Map для хранения объявлений по витринам по основному адресу
      const storefrontUniqueLocations = new Map(); // Map для хранения объявлений витрин с уникальными адресами
      const individualListings = []; // Объявления без витрины

      // Функция для проверки, являются ли координаты уникальными
      const isUniqueLocation = (storefrontId, lat, lng) => {
        // Если у нас нет записи о главном местоположении витрины, то считаем уникальным
        if (!storefrontListings.has(storefrontId)) return true;

        const mainLocation = storefrontListings.get(storefrontId).location;

        // Проверяем, отличаются ли координаты от основного местоположения витрины
        // Используем небольшую погрешность, чтобы избежать проблем с точностью
        const threshold = 0.0001; // примерно 10 метров
        const isDifferent = Math.abs(mainLocation[0] - lat) > threshold ||
          Math.abs(mainLocation[1] - lng) > threshold;

        // Выводим отладочную информацию при обнаружении уникальных координат
        if (isDifferent) {
          console.log(`Обнаружены уникальные координаты для объявления витрины ${storefrontId}: [${lat}, ${lng}] отличаются от основных [${mainLocation[0]}, ${mainLocation[1]}]`);
        }

        return isDifferent;
      };

      // Сначала проходим по всем объявлениям для сбора основной информации о витринах
      validListings.forEach(listing => {
        if (listing.storefront_id) {
          if (!storefrontListings.has(listing.storefront_id)) {
            // Создаем новую запись о витрине
            storefrontListings.set(listing.storefront_id, {
              listings: [],
              name: listing.storefront_name || t('listings.map.storefront'),
              location: [listing.latitude, listing.longitude],
              address: listing.location || listing.address
            });
          }
          // Добавляем объявление в общий список витрины
          storefrontListings.get(listing.storefront_id).listings.push(listing);
        } else {
          // Объявления без витрины сразу добавляем в список индивидуальных
          individualListings.push(listing);
        }
      });

      // Ограничиваем количество индивидуальных объявлений
      const MAX_INDIVIDUAL_LISTINGS = 1000;
      const limitedIndividualListings = individualListings.slice(0, MAX_INDIVIDUAL_LISTINGS);

      // Ограничиваем количество витрин
      const MAX_STOREFRONTS = 200;
      // Берем только первые MAX_STOREFRONTS витрин если их больше
      const limitedStorefrontIds = Array.from(storefrontListings.keys()).slice(0, MAX_STOREFRONTS);

      // Ограничиваем количество объявлений в одной витрине
      const MAX_ITEMS_PER_STOREFRONT = 10000000;

      console.log(`Обрабатываем ${limitedStorefrontIds.length} витрин из ${storefrontListings.size}`);
      console.log(`Обрабатываем ${limitedIndividualListings.length} индивидуальных объявлений (лимит: ${MAX_INDIVIDUAL_LISTINGS})`);

      // Теперь проходим по всем объявлениям снова для определения уникальных местоположений
      console.log("Начинаем поиск объявлений с уникальными местоположениями...");
      validListings.forEach(listing => {
        if (listing.storefront_id && listing.latitude && listing.longitude) {
          // Логируем информацию о каждом объявлении витрины
          console.log(`Проверка объявления: ID=${listing.id}, storefront_id=${listing.storefront_id}, координаты=[${listing.latitude}, ${listing.longitude}]`);

          // Проверяем, уникально ли местоположение для данной витрины
          if (isUniqueLocation(listing.storefront_id, listing.latitude, listing.longitude)) {
            // Создаем уникальный ключ для этого местоположения
            const locationKey = `${listing.storefront_id}_${listing.latitude}_${listing.longitude}`;
            console.log(`НАЙДЕНО УНИКАЛЬНОЕ МЕСТОПОЛОЖЕНИЕ: key=${locationKey}`);

            if (!storefrontUniqueLocations.has(locationKey)) {
              // Если это первое объявление с такими координатами, создаем новую запись
              const storefrontName = listing.storefront_name ||
                (storefrontListings.has(listing.storefront_id) ?
                  storefrontListings.get(listing.storefront_id).name :
                  t('listings.map.storefront'));

              storefrontUniqueLocations.set(locationKey, {
                listings: [],
                storefront_id: listing.storefront_id,
                name: storefrontName,
                location: [listing.latitude, listing.longitude],
                address: listing.location || listing.address
              });

              console.log(`Создана новая запись уникального местоположения: ${storefrontName}, координаты=[${listing.latitude}, ${listing.longitude}]`);
            }

            // Добавляем объявление в список уникальных местоположений
            storefrontUniqueLocations.get(locationKey).listings.push(listing);
            console.log(`Объявление ID=${listing.id} добавлено в список уникальных местоположений, всего: ${storefrontUniqueLocations.get(locationKey).listings.length}`);

            // И удаляем его из основного списка витрины, чтобы избежать дублирования
            if (storefrontListings.has(listing.storefront_id)) {
              const mainStorefrontListings = storefrontListings.get(listing.storefront_id).listings;
              const index = mainStorefrontListings.findIndex(item => item.id === listing.id);
              if (index >= 0) {
                mainStorefrontListings.splice(index, 1);
                console.log(`Объявление ID=${listing.id} удалено из основного списка витрины, оставшихся: ${mainStorefrontListings.length}`);
              } else {
                console.log(`Объявление ID=${listing.id} не найдено в основном списке витрины`);
              }
            } else {
              console.log(`Витрина ID=${listing.storefront_id} не найдена в списке основных витрин`);
            }
          }
        }
      });

      // Отладочная информация
      console.log(`Найдено ${storefrontListings.size} витрин на карте`);
      storefrontListings.forEach((storefront, id) => {
        console.log(`Витрина ID ${id}: ${storefront.name}, объявлений: ${storefront.listings.length}`);
      });

      console.log(`Найдено ${storefrontUniqueLocations.size} уникальных местоположений товаров витрин`);
      storefrontUniqueLocations.forEach((location, key) => {
        console.log(`Уникальное местоположение: ${key}, витрина ID: ${location.storefront_id}, объявлений: ${location.listings.length}`);
      });

      // Добавляем маркеры для витрин
      limitedStorefrontIds.forEach(storefrontId => {
        try {
          const storefront = storefrontListings.get(storefrontId);
          // Ограничиваем количество товаров в одной витрине до разумного предела
          const limitedStorefrontListings = storefront.listings;

          // Создаем маркер для витрины в виде магазина
          const storeMarker = L.marker(storefront.location, {
            icon: L.divIcon({
              html: `
                <div style="
                  width: 42px;
                  height: 42px;
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                  justify-content: center;
                ">
                  <div style="
                    background-color: ${theme.palette.primary.main};
                    color: white;
                    width: 40px;
                    height: 32px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    border-radius: 4px;
                    position: relative;
                    border: 2px solid white;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.3);
                  ">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"></path>
                      <polyline points="9 22 9 12 15 12 15 22"></polyline>
                    </svg>
                    <div style="
                      position: absolute;
                      top: -10px;
                      right: -10px;
                      background-color: ${theme.palette.error.main};
                      color: white;
                      border-radius: 50%;
                      width: 22px;
                      height: 22px;
                      font-size: 12px;
                      font-weight: bold;
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      border: 1px solid white;
                    ">${storefront.listings.length > 999 ? '999+' : storefront.listings.length}</div>
                  </div>
                  <div style="
                    width: 0;
                    height: 0;
                    border-left: 8px solid transparent;
                    border-right: 8px solid transparent;
                    border-top: 10px solid ${theme.palette.primary.main};
                    margin-top: -2px;
                  "></div>
                </div>
              `,
              className: 'storefront-marker',
              iconSize: [42, 42],
              iconAnchor: [21, 42]
            })
          })
          .bindTooltip(`${storefront.name} (${storefront.listings.length} ${t('common:map.items')})`)
            .on('click', () => {
              // При клике показываем первое объявление из витрины с меткой витрины
              if (limitedStorefrontListings.length > 0) {
                const firstListing = limitedStorefrontListings[0];
                // Добавляем информацию о том, что это часть витрины
                firstListing.isPartOfStorefront = true;
                firstListing.storefrontName = storefront.name;
                firstListing.storefrontItemCount = storefront.listings.length;
                firstListing.storefront_id = storefrontId; // Добавляем ID витрины для перехода
                // Добавляем ID объявления для выделения при переходе на страницу витрины
                firstListing.id = firstListing.id || limitedStorefrontListings[0].id;
                setSelectedListing(firstListing);
              }
            });

          storeMarker.addTo(mapRef.current);
          markerGroup.addLayer(storeMarker);
          markersRef.current.push(storeMarker);
        } catch (error) {
          console.error("Error adding storefront marker:", error);
        }
      });

      // Добавляем маркеры для уникальных местоположений объявлений витрины
      storefrontUniqueLocations.forEach((location, locationKey) => {
        try {
          // Создаем маркер с индикацией, что это часть витрины, но с уникальным местоположением
          const uniqueStoreMarker = L.marker(location.location, {
            icon: L.divIcon({
              html: `
                <div style="
                  width: 36px;
                  height: 36px;
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                  justify-content: center;
                ">
                  <div style="
                    background-color: ${theme.palette.secondary.main};
                    color: white;
                    width: 34px;
                    height: 28px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    border-radius: 4px;
                    position: relative;
                    border: 2px solid white;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.3);
                  ">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                      <circle cx="12" cy="10" r="3"></circle>
                    </svg>
                    <div style="
                      position: absolute;
                      top: -8px;
                      right: -8px;
                      background-color: ${theme.palette.error.main};
                      color: white;
                      border-radius: 50%;
                      width: 20px;
                      height: 20px;
                      font-size: 10px;
                      font-weight: bold;
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      border: 1px solid white;
                    ">${location.listings.length > 99 ? '99+' : location.listings.length}</div>
                  </div>
                  <div style="
                    width: 0;
                    height: 0;
                    border-left: 7px solid transparent;
                    border-right: 7px solid transparent;
                    border-top: 8px solid ${theme.palette.secondary.main};
                    margin-top: -2px;
                  "></div>
                </div>
              `,
              className: 'storefront-unique-marker',
              iconSize: [36, 36],
              iconAnchor: [18, 36]
            })
          })
            .bindTooltip(`${location.name} (${t('common:map.uniqueLocation')}) - ${location.listings.length} ${t('common:map.items')}`)

            .on('click', () => {
              // При клике показываем первое объявление с этого уникального местоположения
              if (location.listings.length > 0) {
                // Создаем полностью новый объект (важно для избежания проблем с React)
                const uniqueListing = JSON.parse(JSON.stringify(location.listings[0]));

                // Устанавливаем необходимые атрибуты
                uniqueListing.isPartOfStorefront = true; // Ключевое свойство!
                uniqueListing.storefrontName = location.name;
                uniqueListing.storefrontItemCount = location.listings.length;
                uniqueListing.storefront_id = location.storefront_id;
                uniqueListing.isUniqueLocation = true;

                // Убеждаемся, что у листинга есть ID
                if (!uniqueListing.id && location.listings[0].id) {
                  uniqueListing.id = location.listings[0].id;
                }

                console.log("Подготовлен объект для уникального местоположения:", {
                  id: uniqueListing.id,
                  title: uniqueListing.title,
                  storefront_id: uniqueListing.storefront_id,
                  isPartOfStorefront: uniqueListing.isPartOfStorefront,
                  isUniqueLocation: uniqueListing.isUniqueLocation
                });

                // Проверяем наличие id и storefront_id для перехода
                if (uniqueListing.id && uniqueListing.storefront_id) {
                  console.log("Устанавливаем выделенный листинг для уникального местоположения");
                  setSelectedListing(uniqueListing);
                } else {
                  console.error("Не удалось получить ID объявления или витрины для уникального местоположения");
                }
              }
            })

          uniqueStoreMarker.addTo(mapRef.current);
          markerGroup.addLayer(uniqueStoreMarker);
          markersRef.current.push(uniqueStoreMarker);

          console.log(`Добавлен маркер для товаров с уникальным местоположением: ${location.name}, координаты=[${location.location[0]}, ${location.location[1]}], объявлений: ${location.listings.length}`);
        } catch (error) {
          console.error("Error adding unique location storefront marker:", error);
        }
      });

      // Добавляем маркеры для индивидуальных объявлений
      individualListings.forEach(listing => {
        try {
          const marker = L.marker([listing.latitude, listing.longitude])
            .bindTooltip(`${listing.price.toLocaleString()} RSD`)
            .on('click', () => {
              setSelectedListing(listing);
            });

          marker.addTo(mapRef.current);
          markerGroup.addLayer(marker);
          markersRef.current.push(marker);
        } catch (error) {
          console.error("Error adding marker:", error);
        }
      });

      // Устанавливаем границы карты, чтобы были видны все маркеры
      // если нет пользовательского местоположения
      if (!userLocation && markerGroup.getLayers().length > 0) {
        try {
          mapRef.current.fitBounds(markerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
          });
        } catch (error) {
          console.error("Error fitting bounds:", error);
        }
      }
    } catch (error) {
      console.error("Error updating markers:", error);
    }
  }, [listings, mapReady, userLocation, t, theme]);

  // Обработчик изменения радиуса поиска
  const handleRadiusChange = (event) => {
    onFilterChange({ ...filters, distance: event.target.value });
  };

  // Навигация к подробностям объявления
  const handleNavigateToListing = (listingId) => {
    navigate(`/marketplace/listings/${listingId}`);
  };

  // Обработчик для открытия полноэкранной карты
  const handleExpandMap = () => {
    // Получаем список всех объявлений с координатами
    const validListings = listings.filter(listing =>
      listing.latitude && listing.longitude && listing.show_on_map !== false
    );

    // Формируем маркеры для полноэкранной карты
    const markersForFullscreen = validListings.map(listing => ({
      latitude: listing.latitude,
      longitude: listing.longitude,
      title: listing.title,
      tooltip: `${listing.price.toLocaleString()} RSD`,
      id: listing.id,
      listing: listing // Передаем полные данные о листинге
    }));

    // Определяем центр для полноэкранной карты с гарантированными координатами
    let center = null;

    // Первый случай: выбранное объявление
    if (selectedListing && selectedListing.latitude && selectedListing.longitude) {
      center = {
        latitude: selectedListing.latitude,
        longitude: selectedListing.longitude,
        title: selectedListing.title
      };
      console.log("Используем координаты выбранного объявления:", center);
    }

    // Второй случай: местоположение пользователя
    else if (userLocation && userLocation.lat && userLocation.lon) {
      center = {
        latitude: userLocation.lat,
        longitude: userLocation.lon,
        title: t('listings.map.yourLocation')
      };
      console.log("Используем координаты пользователя:", center);
    }
    // Третий случай: текущий центр карты
    else if (mapRef.current) {
      try {
        const mapCenter = mapRef.current.getCenter();
        center = {
          latitude: mapCenter.lat,
          longitude: mapCenter.lng,
          title: t('listings.map.mapCenter')
        };
        console.log("Используем текущий центр карты:", center);
      } catch (error) {
        console.error("Ошибка при получении центра карты:", error);
      }
    }
    // Четвертый случай: первое объявление из списка
    else if (validListings.length > 0) {
      const firstListing = validListings[0];
      center = {
        latitude: firstListing.latitude,
        longitude: firstListing.longitude,
        title: firstListing.title
      };
      console.log("Используем координаты первого объявления:", center);
    }
    // Пятый случай: фиксированные координаты по умолчанию (Нови-Сад)
    else {
      center = {
        latitude: 45.2671,
        longitude: 19.8335,
        title: "Нови-Сад"
      };
      console.log("Используем координаты по умолчанию:", center);
    }

    // Дополнительная проверка перед установкой состояния
    if (!center || !center.latitude || !center.longitude) {
      console.error("Не удалось определить координаты для карты:", center);
      // Устанавливаем координаты по умолчанию
      center = {
        latitude: 45.2671,
        longitude: 19.8335,
        title: "Нови-Сад"
      };
    }

    // Проверяем, что у нас есть числовые значения для координат
    center.latitude = Number(center.latitude);
    center.longitude = Number(center.longitude);

    console.log("Итоговые координаты для полноэкранной карты:", center);

    // Устанавливаем состояние и открываем модальное окно
    setExpandedMapCenter(center);
    setExpandedMapMarkers(markersForFullscreen);
    setExpandedMapOpen(true);
  };

  // Функция для определения местоположения пользователя
  const handleDetectLocation = async () => {
    try {
      // Используем функцию из контекста местоположения
      const locationData = await detectUserLocation();

      // Если успешно получили местоположение, обновляем фильтры
      onFilterChange({
        ...filters,
        latitude: locationData.lat,
        longitude: locationData.lon,
        distance: filters.distance || '5km'
      });

      // Центрируем карту на новых координатах
      if (mapRef.current) {
        mapRef.current.setView([locationData.lat, locationData.lon], 13);
      }
    } catch (error) {
      console.error("Error getting location:", error);
      alert(t('listings.map.locationError'));
    }
  };

  // Определяем, доступна ли карта (не используется в запросе distance без координат)
  const isMapAvailable = useMemo(() => {
    return !filters.distance || (filters.latitude && filters.longitude);
  }, [filters.distance, filters.latitude, filters.longitude]);

  const isDistanceWithoutCoordinates = filters.distance && (!filters.latitude || !filters.longitude);

  return (
    <Box
      sx={{
        position: 'relative',
        height: isMobile ? 'calc(100vh - 120px)' : '90vh',
        width: '100%',
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      {/* Панель инструментов карты */}
      <Paper
        elevation={3}
        sx={{
          p: 2,
          mb: 2,
          zIndex: 1000,
          display: 'flex',
          flexWrap: 'wrap',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 2
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>


          {/* Показываем информацию о количестве объявлений и витрин на карте */}
          {(() => {
            // Фильтруем объявления с координатами и флагом show_on_map
            const validListings = listings.filter(l => l.latitude && l.longitude && l.show_on_map !== false);

            // Подсчитываем количество уникальных витрин
            const storefrontsSet = new Set();
            validListings.forEach(l => {
              if (l.storefront_id) {
                storefrontsSet.add(l.storefront_id);
              }
            });

            // Подсчитываем количество частных объявлений (без витрины)
            const individualListingsCount = validListings.filter(l => !l.storefront_id).length;

            return (
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Chip
                  label={`${validListings.length} ${t('listings.map.itemsOnMap')}`}
                  color="primary"
                  variant="outlined"
                />
                <Chip
                  label={`${storefrontsSet.size} ${t('listings.map.storefrontsOnMap', { defaultValue: 'витрин' })}`}
                  color="secondary"
                  variant="outlined"
                />
                <Chip
                  label={`${individualListingsCount} ${t('listings.map.individualItemsOnMap', { defaultValue: 'объявлений' })}`}
                  color="info"
                  variant="outlined"
                />
              </Box>
            );
          })()}
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>


          <Button
            variant="outlined"
            startIcon={<List />}
            onClick={onMapClose}
          >
            {isMobile ? t('listings.map.list') : t('listings.map.backToList')}
          </Button>
        </Box>
      </Paper>

      {/* Контейнер карты */}
      <Box
        sx={{
          flex: 1,
          borderRadius: 1,
          overflow: 'hidden',
          position: 'relative'
        }}
      >
        {/* Добавляем кнопку "Развернуть" в стиле MiniMap */}
        <IconButton
          onClick={handleExpandMap}
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
            bgcolor: 'background.paper',
            '&:hover': {
              bgcolor: 'background.paper',
            },
            zIndex: 1000,
            boxShadow: '0 2px 6px rgba(0,0,0,0.1)'
          }}
          size="small"
        >
          <Maximize2 size={20} />
        </IconButton>

        <div
          ref={mapContainerRef}
          style={{ width: '100%', height: '100%' }}
        />

        {!mapReady && (
          <Box
            sx={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              bgcolor: 'background.paper',
              zIndex: 999
            }}
          >
            <Typography variant="h6">
              {t('listings.map.loadingMap', { defaultValue: 'Загрузка карты...' })}
            </Typography>
          </Box>
        )}
      </Box>

      {/* Информация о выбранном объявлении */}
      {selectedListing && (
        <ListingPreview
          listing={selectedListing}
          onClose={() => setSelectedListing(null)}
          onNavigate={handleNavigateToListing}
        />
      )}

      {/* Кнопка определения местоположения */}
      {!userLocation && (
        <Button
          variant="contained"
          color="primary"
          startIcon={<Navigation />}
          sx={{
            position: 'absolute',
            bottom: 16,
            right: 16,
            zIndex: 1000
          }}
          onClick={handleDetectLocation}
        >
          {t('listings.map.useMyLocation')}
        </Button>
      )}

      {/* Модальное окно с полноэкранной картой */}
      <Modal
        open={expandedMapOpen}
        onClose={() => setExpandedMapOpen(false)}
        aria-labelledby="expanded-map-modal"
        aria-describedby="expanded-map-view"
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          p: 3
        }}
      >
        <Paper
          sx={{
            width: '90%',
            maxWidth: 1200,
            maxHeight: '90vh',
            bgcolor: 'background.paper',
            borderRadius: 2,
            boxShadow: 24,
            position: 'relative',
            overflow: 'hidden'
          }}
        >
          <Box sx={{ position: 'absolute', top: 8, right: 8, zIndex: 1050 }}>
            <IconButton
              onClick={() => setExpandedMapOpen(false)}
              sx={{
                bgcolor: 'background.paper',
                '&:hover': {
                  bgcolor: 'background.paper',
                },
                boxShadow: '0 2px 6px rgba(0,0,0,0.2)'
              }}
            >
              <X size={20} />
            </IconButton>
          </Box>

          {expandedMapCenter && (
            <FullscreenMap
              latitude={expandedMapCenter.latitude}
              longitude={expandedMapCenter.longitude}
              title={expandedMapCenter.title}
              markers={expandedMapMarkers}
            />
          )}
        </Paper>
      </Modal>
    </Box>
  );
};

export default MapView;
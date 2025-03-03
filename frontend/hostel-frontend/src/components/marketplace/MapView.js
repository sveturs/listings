// frontend/hostel-frontend/src/components/marketplace/MapView.js
import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Navigation, X, List, Maximize2 } from 'lucide-react';
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

// Компонент для предпросмотра объявления при клике по маркеру
const ListingPreview = ({ listing, onClose, onNavigate }) => {
  const { t } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  if (!listing) return null;
  
  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };
  
  const imageUrl = listing.images && listing.images.length > 0
    ? `${process.env.REACT_APP_BACKEND_URL}/uploads/${listing.images[0].file_path}`
    : null;
    
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
            '&:hover': { bgcolor: 'background.paper', opacity: 1 }
          }}
        >
          <X size={16} />
        </IconButton>
        
        {imageUrl && (
          <CardMedia
            component="img"
            height="140"
            image={imageUrl}
            alt={listing.title}
          />
        )}
        
        <CardContent>
          <Typography variant="subtitle1" noWrap gutterBottom>
            {listing.title}
          </Typography>
          
          <Typography variant="h6" color="primary" gutterBottom>
            {formatPrice(listing.price)}
          </Typography>
          
          <Box display="flex" justifyContent="space-between" alignItems="center" mt={1}>
            <Chip 
              label={listing.condition === 'new' ? t('listings.conditions.new') : t('listings.conditions.used')}
              size="small"
              color={listing.condition === 'new' ? 'success' : 'default'}
            />
            
            <Button 
              variant="contained" 
              size="small" 
              onClick={() => onNavigate(listing.id)}
            >
              {t('listings.details.viewDetails')}
            </Button>
          </Box>
        </CardContent>
      </Box>
    </Card>
  );
};

const MapView = ({ listings, userLocation, filters, onFilterChange, onMapClose }) => {
  const { t } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const navigate = useNavigate();
  const mapRef = useRef(null);
  const markersRef = useRef([]);
  const mapContainerRef = useRef(null);
  const [selectedListing, setSelectedListing] = useState(null);
  const [mapReady, setMapReady] = useState(false);
  
  // Состояние для модального окна с полноэкранной картой
  const [expandedMapOpen, setExpandedMapOpen] = useState(false);
  // Центр для полноэкранной карты
  const [expandedMapCenter, setExpandedMapCenter] = useState(null);
  // Маркеры для полноэкранной карты
  const [expandedMapMarkers, setExpandedMapMarkers] = useState([]);
  
  // Варианты радиуса поиска
  const radiusOptions = [
    { value: "1km", label: "1 км" },
    { value: "3km", label: "3 км" },
    { value: "5km", label: "5 км" },
    { value: "10km", label: "10 км" },
    { value: "15km", label: "15 км" },
    { value: "30km", label: "30 км" }
  ];
  
  // Инициализация карты
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;
    
    // Устанавливаем центр карты
    const initialPosition = userLocation 
      ? [userLocation.latitude, userLocation.longitude] 
      : [45.2671, 19.8335]; // Координаты Нови-Сада по умолчанию
      
    // Создаем карту
    mapRef.current = L.map(mapContainerRef.current).setView(initialPosition, 13);
    
    // Добавляем слой тайлов
    L.tileLayer(TILE_LAYER_URL, {
      attribution: TILE_LAYER_ATTRIBUTION,
      maxZoom: 19
    }).addTo(mapRef.current);
    
    // Если есть местоположение пользователя, добавляем маркер
    if (userLocation) {
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
    }
    
    setMapReady(true);
    
    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
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
    if (!mapRef.current || !userLocation || !filters.distance) return;
    
    // Удаляем старые круги
    mapRef.current.eachLayer(layer => {
      if (layer instanceof L.Circle) {
        mapRef.current.removeLayer(layer);
      }
    });
    
    // Добавляем новый круг с актуальным радиусом
    const radiusInMeters = getRadiusInMeters(filters.distance);
    L.circle([userLocation.latitude, userLocation.longitude], {
      color: theme.palette.primary.main,
      fillColor: theme.palette.primary.light,
      fillOpacity: 0.2,
      radius: radiusInMeters
    }).addTo(mapRef.current);
    
  }, [filters.distance, userLocation, theme]);
  
  // Обновляем маркеры объявлений
  useEffect(() => {
    if (!mapRef.current || !mapReady) return;
    
    // Удаляем старые маркеры
    markersRef.current.forEach(marker => {
      mapRef.current.removeLayer(marker);
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
    
    // Добавляем новые маркеры
    validListings.forEach(listing => {
      const marker = L.marker([listing.latitude, listing.longitude])
        .bindTooltip(`${listing.price.toLocaleString()} RSD`)
        .on('click', () => {
          setSelectedListing(listing);
        });
      
      marker.addTo(mapRef.current);
      markerGroup.addLayer(marker);
      markersRef.current.push(marker);
    });
    
    // Устанавливаем границы карты, чтобы были видны все маркеры
    // если нет пользовательского местоположения
    if (!userLocation) {
      mapRef.current.fitBounds(markerGroup.getBounds(), {
        padding: [50, 50],
        maxZoom: 15
      });
    }
    
  }, [listings, mapReady, userLocation]);
  
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
    const markersForFullscreen = listings
      .filter(listing => listing.latitude && listing.longitude && listing.show_on_map !== false)
      .map(listing => ({
        latitude: listing.latitude,
        longitude: listing.longitude,
        title: listing.title,
        tooltip: `${listing.price.toLocaleString()} RSD`
      }));
    
    // Определяем центр для полноэкранной карты
    let center = null;
    
    if (selectedListing && selectedListing.latitude && selectedListing.longitude) {
      // Если выбрано объявление, используем его координаты как центр
      center = {
        latitude: selectedListing.latitude,
        longitude: selectedListing.longitude,
        title: selectedListing.title
      };
    } else if (userLocation) {
      // Если есть местоположение пользователя, используем его как центр
      center = {
        latitude: userLocation.latitude,
        longitude: userLocation.longitude,
        title: t('listings.map.yourLocation')
      };
    } else if (mapRef.current) {
      // Иначе используем текущий центр карты
      const mapCenter = mapRef.current.getCenter();
      center = {
        latitude: mapCenter.lat,
        longitude: mapCenter.lng,
        title: t('listings.map.mapCenter')
      };
    }
    
    setExpandedMapCenter(center);
    setExpandedMapMarkers(markersForFullscreen);
    setExpandedMapOpen(true);
  };
  
  // Обработчик для закрытия полноэкранной карты
  const handleCloseExpandedMap = () => {
    setExpandedMapOpen(false);
  };
  
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
          <Typography variant="subtitle1">
            {t('listings.map.showOnMap')}
          </Typography>
          
          <Chip 
            label={`${listings.filter(l => l.latitude && l.longitude && l.show_on_map !== false).length} ${t('listings.map.itemsOnMap')}`}
            color="primary"
            variant="outlined"
          />
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          {userLocation && (
            <FormControl size="small" sx={{ minWidth: 120 }}>
              <InputLabel id="radius-select-label">{t('listings.map.radius')}</InputLabel>
              <Select
                labelId="radius-select-label"
                value={filters.distance || '5km'}
                label={t('listings.map.radius')}
                onChange={handleRadiusChange}
              >
                {radiusOptions.map(option => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}
          
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
            zIndex: 1000, // Убедимся, что кнопка отображается поверх карты
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
          onClick={() => {
            if (navigator.geolocation) {
              navigator.geolocation.getCurrentPosition(
                (position) => {
                  const { latitude, longitude } = position.coords;
                  
                  // Обновляем фильтры с новыми координатами
                  onFilterChange({
                    ...filters,
                    latitude,
                    longitude,
                    distance: filters.distance || '5km'
                  });
                  
                  // Центрируем карту на новых координатах
                  if (mapRef.current) {
                    mapRef.current.setView([latitude, longitude], 13);
                  }
                },
                (error) => {
                  console.error("Error getting location:", error);
                  alert(t('listings.map.locationError'));
                }
              );
            } else {
              alert(t('listings.map.geolocationNotSupported'));
            }
          }}
        >
          {t('listings.map.useMyLocation')}
        </Button>
      )}
      
      {/* Модальное окно с полноэкранной картой */}
      <Modal
        open={expandedMapOpen}
        onClose={handleCloseExpandedMap}
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
              onClick={handleCloseExpandedMap}
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
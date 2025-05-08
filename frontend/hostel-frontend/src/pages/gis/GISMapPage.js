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
import { TILE_LAYER_URL, TILE_LAYER_ATTRIBUTION } from '../../components/maps/map-constants';
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
    sort: searchParams.get('sort_by') || 'newest',
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

  // Инициализация карты
  useEffect(() => {
    // Проверяем все необходимые условия для инициализации карты
    if (!mapContainerRef.current || !mapCenter || mapRef.current) return;
    
    // Добавляем задержку для уверенности, что DOM-элемент полностью готов
    const initMapTimer = setTimeout(() => {
      try {
        console.log("Initializing GIS map with center:", mapCenter);
        
        // Создаем карту с настройками для высокой производительности
        mapRef.current = L.map(mapContainerRef.current, {
          preferCanvas: true,
          zoomControl: false,
          attributionControl: false
        }).setView([mapCenter.latitude, mapCenter.longitude], mapZoom);
        
        // Добавляем базовый слой карты
        L.tileLayer(TILE_LAYER_URL, {
          attribution: TILE_LAYER_ATTRIBUTION,
          maxZoom: 19
        }).addTo(mapRef.current);
        
        // Добавляем элементы управления
        L.control.zoom({
          position: 'bottomright'
        }).addTo(mapRef.current);
        
        // Создаем слой маркеров
        markersLayerRef.current = L.layerGroup().addTo(mapRef.current);
        
        // Перерисовываем карту для уверенности в корректных размерах
        mapRef.current.invalidateSize();
        
        // Настраиваем обработчики событий
        mapRef.current.on('moveend', handleMapMoveEnd);
        mapRef.current.on('zoomend', handleMapZoomEnd);
        
        setMapReady(true);
        console.log("GIS map successfully initialized");
        
      } catch (error) {
        console.error("Error initializing map:", error);
        setError("Не удалось инициализировать карту: " + error.message);
      }
    }, 100); // Небольшая задержка для уверенности, что DOM готов
    
    // Функция очистки
    return () => {
      clearTimeout(initMapTimer);
      
      if (mapRef.current) {
        try {
          console.log("Cleaning up GIS map");
          mapRef.current.off('moveend', handleMapMoveEnd);
          mapRef.current.off('zoomend', handleMapZoomEnd);
          mapRef.current.remove();
        } catch (err) {
          console.error("Error removing map:", err);
        }
        mapRef.current = null;
      }
      setMapReady(false);
    };
  }, [mapCenter]);

  // Обработчики событий карты
  const handleMapMoveEnd = () => {
    if (!mapRef.current) return;
    
    const center = mapRef.current.getCenter();
    const zoom = mapRef.current.getZoom();
    
    setMapCenter({
      latitude: center.lat,
      longitude: center.lng
    });
    
    // Обновляем параметры URL без перезагрузки страницы
    setSearchParams(prev => {
      const newParams = new URLSearchParams(prev);
      newParams.set('latitude', center.lat.toFixed(6));
      newParams.set('longitude', center.lng.toFixed(6));
      newParams.set('zoom', zoom);
      return newParams;
    }, { replace: true });
    
    // При перемещении карты загружаем новые данные
    fetchListingsInViewport();
  };
  
  const handleMapZoomEnd = () => {
    if (!mapRef.current) return;
    setMapZoom(mapRef.current.getZoom());
  };

  // Загрузка объявлений в текущей области видимости
  const fetchListingsInViewport = async () => {
    if (!mapRef.current || !mapReady) return;
    
    try {
      setLoading(true);
      
      // Получаем границы видимой области
      const bounds = mapRef.current.getBounds();
      const northEast = bounds.getNorthEast();
      const southWest = bounds.getSouthWest();
      
      // Формируем параметры запроса
      const params = {
        view_mode: 'map',
        bbox: `${southWest.lat},${southWest.lng},${northEast.lat},${northEast.lng}`,
        size: 1000, // Запрашиваем большое количество объявлений для карты
        query: filters.query,
        category_id: activeCategoryId || filters.category_id,
        min_price: filters.price.min,
        max_price: filters.price.max,
        condition: filters.condition !== 'any' ? filters.condition : '',
        sort_by: filters.sort,
        radius: filters.radius > 0 ? filters.radius : '',
        city: filters.city
      };
      
      // Делаем запрос к API
      const response = await axios.get('/api/v1/marketplace/search', { params });
      
      if (response.data && response.data.data) {
        const fetchedListings = Array.isArray(response.data.data) 
          ? response.data.data 
          : response.data.data.data || [];
        
        setListings(fetchedListings);
        setTotalListings(response.data.meta?.total || fetchedListings.length);
        
        // Обновляем маркеры на карте
        updateMapMarkers(fetchedListings);
      }
    } catch (error) {
      console.error("Error fetching listings:", error);
      setError("Ошибка при загрузке объявлений");
    } finally {
      setLoading(false);
    }
  };

  // Обновление маркеров на карте
  const updateMapMarkers = (listings) => {
    if (!mapRef.current || !markersLayerRef.current) return;
    
    // Очищаем текущие маркеры
    markersLayerRef.current.clearLayers();
    
    // Группируем объявления по витринам
    const storefronts = new Map();
    const individualListings = [];
    
    listings.forEach(listing => {
      if (!listing.latitude || !listing.longitude) return;
      
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
    
    // Создаем маркеры для витрин
    storefronts.forEach((storefront, id) => {
      if (storefront.listings.length > 0) {
        const storeMarker = L.marker([storefront.latitude, storefront.longitude], {
          icon: createCustomIcon('store', storefront.listings.length)
        });
        
        storeMarker.on('click', () => {
          setSelectedListing({
            ...storefront.listings[0],
            isStorefront: true,
            storefrontItemCount: storefront.listings.length,
            storefrontName: storefront.name
          });
          setSelectedListingId(storefront.listings[0].id);
        });
        
        markersLayerRef.current.addLayer(storeMarker);
      }
    });
    
    // Создаем маркеры для отдельных объявлений
    individualListings.forEach(listing => {
      const marker = L.marker([listing.latitude, listing.longitude], {
        icon: createMarkerIcon(listing)
      });
      
      marker.on('click', () => {
        setSelectedListing(listing);
        setSelectedListingId(listing.id);
      });
      
      markersLayerRef.current.addLayer(marker);
    });
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
      const updated = { ...prev, ...newFilters };
      
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
    
    // Add new tile layer based on selection
    switch(newLayer) {
      case 'satellite':
        L.tileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}', {
          attribution: 'Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community'
        }).addTo(mapRef.current);
        break;
      case 'terrain':
        L.tileLayer('https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png', {
          attribution: 'Map data: &copy; OpenStreetMap contributors, SRTM | Map style: &copy; OpenTopoMap (CC-BY-SA)'
        }).addTo(mapRef.current);
        break;
      case 'traffic':
        // Default map with traffic layer
        L.tileLayer(TILE_LAYER_URL, {
          attribution: TILE_LAYER_ATTRIBUTION
        }).addTo(mapRef.current);
        // Add traffic layer if available (example only)
        break;
      case 'heatmap':
        // Default map
        L.tileLayer(TILE_LAYER_URL, {
          attribution: TILE_LAYER_ATTRIBUTION
        }).addTo(mapRef.current);
        // Add heatmap layer if data available
        break;
      default: // standard
        L.tileLayer(TILE_LAYER_URL, {
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
            sort: 'newest'
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
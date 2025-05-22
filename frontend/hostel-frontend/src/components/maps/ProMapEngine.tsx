import React, { useEffect, useRef, useState, useCallback, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Box, IconButton, Paper, Fade, Chip, CircularProgress, Typography } from '@mui/material';
import { styled, useTheme, alpha } from '@mui/material/styles';
import { 
  MyLocation, 
  Fullscreen, 
  FullscreenExit, 
  ZoomIn, 
  ZoomOut,
  FilterList,
  ViewList,
  Refresh
} from '@mui/icons-material';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import Supercluster from 'supercluster';
import debounce from 'lodash.debounce';
import axios from '../../api/axios';

// Fix для иконок Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
  iconUrl: require('leaflet/dist/images/marker-icon.png'),
  shadowUrl: require('leaflet/dist/images/marker-shadow.png'),
});

// Типы данных
interface MapMarker {
  id: number;
  latitude: number;
  longitude: number;
  title: string;
  price?: number;
  condition?: 'new' | 'used' | string;
  category_id?: number;
  main_image?: string;
  city?: string;
  country?: string;
  storefront_id?: number;
  storefront_name?: string;
  type?: 'listing' | 'storefront';
  status?: string;
}

interface MapCluster {
  id: string;
  latitude: number;
  longitude: number;
  count: number;
  avg_price?: number;
  categories?: number[];
}

interface MapBounds {
  northEast: { lat: number; lng: number };
  southWest: { lat: number; lng: number };
}

interface MapFilters {
  categories?: string;
  condition?: string;
  min_price?: string;
  max_price?: string;
  query?: string;
  category_id?: string;
}

interface ProMapEngineProps {
  initialCenter?: [number, number];
  initialZoom?: number;
  filters?: MapFilters;
  onMarkerClick?: (marker: MapMarker) => void;
  onBoundsChange?: (bounds: MapBounds) => void;
  onZoomChange?: (zoom: number) => void;
  onFilterChange?: (filters: MapFilters) => void;
  className?: string;
  fullscreenDefault?: boolean;
  loading?: boolean;
  markers?: MapMarker[];
  enableClustering?: boolean;
}

// Styled компоненты
const MapContainer = styled(Box)(({ theme }) => ({
  position: 'relative',
  width: '100%',
  height: '100%',
  borderRadius: theme.spacing(1),
  overflow: 'hidden',
  backgroundColor: '#f8fafc',
  '& .leaflet-container': {
    width: '100%',
    height: '100%',
    borderRadius: theme.spacing(1),
    zIndex: 1,
  },
  '& .leaflet-control-zoom': {
    display: 'none', // Скрываем стандартные контролы
  },
  '& .leaflet-control-attribution': {
    display: 'none',
  }
}));

const ControlsContainer = styled(motion.div)(({ theme }) => ({
  position: 'absolute',
  top: theme.spacing(2),
  right: theme.spacing(2),
  zIndex: 1000,
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
}));

const MapStats = styled(motion.div)(({ theme }) => ({
  position: 'absolute',
  bottom: theme.spacing(2),
  left: theme.spacing(2),
  zIndex: 1000,
  display: 'flex',
  gap: theme.spacing(1),
  flexWrap: 'wrap',
}));

const LoadingOverlay = styled(Box)(({ theme }) => ({
  position: 'absolute',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  backgroundColor: 'rgba(255, 255, 255, 0.8)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  zIndex: 2000,
  borderRadius: theme.spacing(1),
}));

// Создаем кастомные иконки для кластеров
const createClusterIcon = (count: number, avg_price?: number) => {
  const size = Math.min(60, 30 + Math.log(count) * 5);
  const opacity = Math.min(0.9, 0.6 + count / 50);
  
  return L.divIcon({
    html: `
      <div style="
        width: ${size}px;
        height: ${size}px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: bold;
        font-size: ${Math.max(12, size / 4)}px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        border: 3px solid white;
        opacity: ${opacity};
        transition: all 0.3s ease;
        cursor: pointer;
      ">
        ${count}
        ${avg_price ? `<div style="font-size: 8px; position: absolute; bottom: -2px;">${Math.round(avg_price)}€</div>` : ''}
      </div>
    `,
    className: 'custom-cluster-icon',
    iconSize: L.point(size, size),
    iconAnchor: L.point(size / 2, size / 2)
  });
};

// Создаем кастомные иконки для маркеров
const createMarkerIcon = (marker: MapMarker) => {
  const priceColor = marker.condition === 'new' ? '#4caf50' : '#ff9800';
  
  return L.divIcon({
    html: `
      <div style="
        width: 40px;
        height: 40px;
        background: ${priceColor};
        border-radius: 50% 50% 50% 0;
        transform: rotate(-45deg);
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 3px 10px rgba(0,0,0,0.2);
        border: 2px solid white;
        cursor: pointer;
        transition: all 0.2s ease;
      ">
        <div style="
          color: white;
          font-weight: bold;
          font-size: 10px;
          transform: rotate(45deg);
          text-align: center;
          line-height: 1;
        ">
          ${marker.price && marker.price > 1000 ? `${Math.round(marker.price/1000)}K` : marker.price || ''}€
        </div>
      </div>
    `,
    className: 'custom-marker-icon',
    iconSize: L.point(40, 40),
    iconAnchor: L.point(20, 35)
  });
};

const ProMapEngine: React.FC<ProMapEngineProps> = ({
  initialCenter = [45.2671, 19.8335], // Нови-Сад по умолчанию
  initialZoom = 13,
  filters = {},
  onMarkerClick,
  onBoundsChange,
  onZoomChange,
  onFilterChange,
  className,
  fullscreenDefault = false,
  loading: externalLoading = false,
  markers: externalMarkers = [],
  enableClustering = true
}) => {
  const theme = useTheme();
  // State
  const [isLoading, setIsLoading] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(fullscreenDefault);
  const [currentZoom, setCurrentZoom] = useState(initialZoom);
  const [markersData, setMarkersData] = useState<MapMarker[]>(externalMarkers);
  const [clustersData, setClustersData] = useState<MapCluster[]>([]);
  const [dataType, setDataType] = useState<'markers' | 'clusters'>('markers');
  const [totalCount, setTotalCount] = useState(0);
  const [mapReady, setMapReady] = useState(false);

  // Refs
  const mapRef = useRef<any>(null);
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const markersLayerRef = useRef<L.LayerGroup>(null);
  const superclusterRef = useRef<Supercluster<any> | null>(null);

  // Инициализация Supercluster
  useEffect(() => {
    superclusterRef.current = new Supercluster({
      radius: 60,
      maxZoom: 16,
      minZoom: 2,
      nodeSize: 64,
    });
  }, []);

  // Функция загрузки данных с API
  const fetchMapData = useCallback(async (bounds: L.LatLngBounds, zoom: number) => {
    if (!bounds) return;

    setIsLoading(true);
    try {
      const ne = bounds.getNorthEast();
      const sw = bounds.getSouthWest();
      
      const params = {
        ne_lat: ne.lat.toString(),
        ne_lng: ne.lng.toString(),
        sw_lat: sw.lat.toString(),
        sw_lng: sw.lng.toString(),
        zoom: zoom.toString(),
        view_mode: 'map',
        bbox: `${sw.lat},${sw.lng},${ne.lat},${ne.lng}`,
        size: zoom >= 15 ? 500 : 100,
        status: 'active',
        ...filters
      };

      // Используем существующий поиск API для карты
      const response = await axios.get('/api/v1/marketplace/search', { params });
      
      if (response.data) {
        let fetchedListings = [];
        if (Array.isArray(response.data.data)) {
          fetchedListings = response.data.data;
        } else if (response.data.data && response.data.data.data) {
          fetchedListings = response.data.data.data;
        }

        const validListings = fetchedListings.filter(listing =>
          listing &&
          typeof listing.latitude === 'number' &&
          typeof listing.longitude === 'number' &&
          isFinite(listing.latitude) &&
          isFinite(listing.longitude)
        );

        setMarkersData(validListings);
        setDataType('markers');
        setTotalCount(response.data.meta?.total || validListings.length);
      }
    } catch (error) {
      console.error('Error fetching map data:', error);
    } finally {
      setIsLoading(false);
    }
  }, [filters]);

  // Debounced версия загрузки данных для плавности
  const debouncedFetchMapData = useMemo(
    () => debounce(fetchMapData, 300),
    [fetchMapData]
  );

  // Инициализация карты
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    // Создаем карту с улучшенными настройками
    const map = L.map(mapContainerRef.current, {
      center: initialCenter,
      zoom: initialZoom,
      zoomControl: false,
      attributionControl: false,
      preferCanvas: true,
      zoomAnimation: true,
      fadeAnimation: true,
      markerZoomAnimation: true,
      transform3DLimit: 2048,
      renderer: L.canvas({ tolerance: 5 })
    });

    // Добавляем современные тайлы
    L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
      attribution: '',
      subdomains: 'abcd',
      maxZoom: 19
    }).addTo(map);

    // Создаем слой для маркеров
    const markersLayer = L.layerGroup().addTo(map);
    markersLayerRef.current = markersLayer;

    mapRef.current = map;

    // События карты
    map.on('moveend zoomend', () => {
      const zoom = map.getZoom();
      setCurrentZoom(zoom);
      
      const bounds = map.getBounds();
      onBoundsChange?.({
        northEast: { lat: bounds.getNorthEast().lat, lng: bounds.getNorthEast().lng },
        southWest: { lat: bounds.getSouthWest().lat, lng: bounds.getSouthWest().lng }
      });
      debouncedFetchMapData(bounds, zoom);
    });

    // Первоначальная загрузка данных
    const bounds = map.getBounds();
    fetchMapData(bounds, initialZoom);

    return () => {
      map.remove();
      mapRef.current = null;
    };
  }, [initialCenter, initialZoom]); // eslint-disable-line

  // Обновление маркеров и кластеров на карте
  useEffect(() => {
    if (!mapRef.current || !markersLayerRef.current) return;

    // Очищаем предыдущие маркеры
    markersLayerRef.current.clearLayers();

    if (dataType === 'markers') {
      // Добавляем отдельные маркеры
      markersData.forEach(marker => {
        const leafletMarker = L.marker([marker.latitude, marker.longitude], {
          icon: createMarkerIcon(marker)
        });

        leafletMarker.on('click', () => {
          onMarkerClick?.(marker);
        });

        // Добавляем hover эффект
        leafletMarker.on('mouseover', function(e) {
          const icon = e.target.getElement();
          if (icon) {
            icon.style.transform = 'scale(1.1)';
            icon.style.zIndex = '1000';
          }
        });

        leafletMarker.on('mouseout', function(e) {
          const icon = e.target.getElement();
          if (icon) {
            icon.style.transform = 'scale(1)';
            icon.style.zIndex = 'auto';
          }
        });

        markersLayerRef.current?.addLayer(leafletMarker);
      });
    } else {
      // Добавляем кластеры
      clustersData.forEach(cluster => {
        const leafletMarker = L.marker([cluster.latitude, cluster.longitude], {
          icon: createClusterIcon(cluster.count, cluster.avg_price)
        });

        leafletMarker.on('click', () => {
          // При клике на кластер - увеличиваем масштаб
          mapRef.current?.setView([cluster.latitude, cluster.longitude], currentZoom + 2, {
            animate: true,
            duration: 0.5
          });
        });

        // Hover эффект для кластеров
        leafletMarker.on('mouseover', function(e) {
          const icon = e.target.getElement();
          if (icon) {
            icon.style.transform = 'scale(1.15)';
            icon.style.zIndex = '1000';
          }
        });

        leafletMarker.on('mouseout', function(e) {
          const icon = e.target.getElement();
          if (icon) {
            icon.style.transform = 'scale(1)';
            icon.style.zIndex = 'auto';
          }
        });

        markersLayerRef.current?.addLayer(leafletMarker);
      });
    }
  }, [markersData, clustersData, dataType, onMarkerClick, currentZoom]);

  // Обработчики контролов
  const handleZoomIn = () => {
    mapRef.current?.zoomIn();
  };

  const handleZoomOut = () => {
    mapRef.current?.zoomOut();
  };

  const handleMyLocation = () => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords;
          mapRef.current?.setView([latitude, longitude], 15, {
            animate: true,
            duration: 1
          });
        },
        (error) => {
          console.error('Error getting location:', error);
        }
      );
    }
  };

  const handleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  const handleRefresh = () => {
    if (mapRef.current) {
      const bounds = mapRef.current.getBounds();
      fetchMapData(bounds, currentZoom);
    }
  };

  return (
    <MapContainer 
      className={className}
      sx={{
        height: isFullscreen ? '100vh' : '70vh',
        position: isFullscreen ? 'fixed' : 'relative',
        top: isFullscreen ? 0 : 'auto',
        left: isFullscreen ? 0 : 'auto',
        right: isFullscreen ? 0 : 'auto',
        bottom: isFullscreen ? 0 : 'auto',
        zIndex: isFullscreen ? 9999 : 'auto',
      }}
    >
      <div ref={mapContainerRef} style={{ width: '100%', height: '100%' }} />
      
      {/* Загрузка */}
      <AnimatePresence>
        {isLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <LoadingOverlay>
              <CircularProgress size={48} />
            </LoadingOverlay>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Контролы карты */}
      <ControlsContainer
        initial={{ x: 20, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        <Paper elevation={3}>
          <IconButton onClick={handleZoomIn} size="small">
            <ZoomIn />
          </IconButton>
        </Paper>
        
        <Paper elevation={3}>
          <IconButton onClick={handleZoomOut} size="small">
            <ZoomOut />
          </IconButton>
        </Paper>
        
        <Paper elevation={3}>
          <IconButton onClick={handleMyLocation} size="small">
            <MyLocation />
          </IconButton>
        </Paper>
        
        <Paper elevation={3}>
          <IconButton onClick={handleRefresh} size="small">
            <Refresh />
          </IconButton>
        </Paper>
        
        <Paper elevation={3}>
          <IconButton onClick={handleFullscreen} size="small">
            {isFullscreen ? <FullscreenExit /> : <Fullscreen />}
          </IconButton>
        </Paper>
      </ControlsContainer>

      {/* Статистика */}
      <MapStats
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <Fade in={totalCount > 0}>
          <Chip 
            label={`${totalCount} объявлений`}
            size="small"
            color="primary"
            variant="filled"
          />
        </Fade>
        
        <Fade in={true}>
          <Chip 
            label={`Zoom: ${currentZoom}`}
            size="small"
            variant="outlined"
          />
        </Fade>
        
        <Fade in={dataType === 'clusters'}>
          <Chip 
            label="Кластеры"
            size="small"
            color="secondary"
            variant="outlined"
          />
        </Fade>
      </MapStats>
    </MapContainer>
  );
};

export default ProMapEngine;
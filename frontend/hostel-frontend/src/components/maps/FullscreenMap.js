// frontend/hostel-frontend/src/components/maps/FullscreenMap.js
import React, { useEffect, useRef, useState } from 'react';
import { Paper, Box, Typography, IconButton, Card, CardContent, CardMedia, Button, Chip } from '@mui/material';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import '../maps/leaflet-icons'; // Импортируем файл с фиксом иконок
import { X } from 'lucide-react';

// Компонент для предпросмотра объявления при клике по маркеру
const ListingPreview = ({ listing, onClose, onNavigate }) => {
  const { t } = useTranslation('marketplace');

  if (!listing) return null;

  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };

  const getImageUrl = () => {
    if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
      return null;
    }

    let mainImage = listing.images.find(img => img && img.is_main === true) || listing.images[0];
    const baseUrl = window.ENV?.REACT_APP_MINIO_URL || window.ENV?.REACT_APP_BACKEND_URL || '';

    if (mainImage && typeof mainImage === 'object') {
      // Если есть публичный URL, используем его напрямую
      if (mainImage.public_url && mainImage.public_url !== '') {
        // Проверяем, абсолютный или относительный URL
        if (mainImage.public_url.startsWith('http')) {
          return mainImage.public_url;
        } else {
          return `${baseUrl}${mainImage.public_url}`;
        }
      }

      // Для MinIO-объектов формируем URL на основе storage_type
      if (mainImage.storage_type === 'minio' ||
          (mainImage.file_path && mainImage.file_path.includes('listings/'))) {
        console.log('Using MinIO URL:', `${baseUrl}${mainImage.public_url}`);
        return `${baseUrl}${mainImage.public_url}`;
      }

      // Обычный файл
      if (mainImage.file_path) {
        return `${baseUrl}/uploads/${mainImage.file_path}`;
      }
    }

    // Для строк (обратная совместимость)
    if (mainImage && typeof mainImage === 'string') {
      return `${baseUrl}/uploads/${mainImage}`;
    }

    return null;
  };

  const imageUrl = getImageUrl();

  return (
    <Card
      sx={{
        position: 'absolute',
        bottom: 16,
        left: 16,
        maxWidth: 400,
        width: 'auto',
        zIndex: 1000,
        borderRadius: 1,
        boxShadow: 3
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

const FullscreenMap = ({ latitude, longitude, title, markers = [] }) => {
  const { t } = useTranslation('marketplace');
  const navigate = useNavigate();
  const mapContainerRef = useRef(null);
  const mapRef = useRef(null);
  const [selectedListing, setSelectedListing] = useState(null);
  const [error, setError] = useState(null);

  // Если координат нет, компонент всё равно отрендерится, но карта не будет инициализирована
  const hasCoordinates = latitude && longitude;

  useEffect(() => {
    // Проверяем условия внутри хука
    if (!hasCoordinates || !mapContainerRef.current) return;

    // Инициализируем карту только если её еще нет
    if (!mapRef.current) {
      mapRef.current = L.map(mapContainerRef.current).setView([latitude, longitude], 15);

      // Добавляем слой тайлов OpenStreetMap
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19
      }).addTo(mapRef.current);

      // Если есть список маркеров, добавляем их
      if (markers && markers.length > 0) {
        const markerGroup = L.featureGroup();

        markers.forEach(marker => {
          const leafletMarker = L.marker([marker.latitude, marker.longitude]);

          // Если у маркера есть всплывающая подсказка
          if (marker.tooltip) {
            leafletMarker.bindTooltip(marker.tooltip);
          }

          // Обработчик клика по маркеру для показа карточки
          leafletMarker.on('click', () => {
            // Если у маркера есть полные данные объявления
            if (marker.listing) {
              setSelectedListing(marker.listing);
            }
            // Если есть только ID, то можно запросить данные с сервера
            else if (marker.id) {
              // Перенаправляем на страницу объявления
              navigate(`/marketplace/listings/${marker.id}`);
            }
            // Если нет полных данных, но есть заголовок, показываем его в попапе
            else if (marker.title) {
              leafletMarker.bindPopup(marker.title).openPopup();
            }
          });

          leafletMarker.addTo(mapRef.current);
          markerGroup.addLayer(leafletMarker);
        });

        // Масштабируем карту, чтобы видеть все маркеры
        if (markers.length > 1) {
          mapRef.current.fitBounds(markerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
          });
        }
      } else if (title) {
        // Если маркеров нет, но есть центральная точка с заголовком
        const marker = L.marker([latitude, longitude]).addTo(mapRef.current);
        marker.bindPopup(title);
      }
    } else {
      // Если карта уже существует, обновляем вид и маркеры
      console.log("Обновление карты с координатами:", [latitude, longitude]);
      // Обновляем центр и масштаб карты
      mapRef.current.setView([latitude, longitude], 15);

      // Сначала очищаем все существующие маркеры
      mapRef.current.eachLayer(layer => {
        if (layer instanceof L.Marker) {
          mapRef.current.removeLayer(layer);
        }
      });

      // Затем добавляем новые маркеры
      if (markers && markers.length > 0) {
        const markerGroup = L.featureGroup();

        markers.forEach(marker => {
          const leafletMarker = L.marker([marker.latitude, marker.longitude]);

          if (marker.tooltip) {
            leafletMarker.bindTooltip(marker.tooltip);
          }

          // Обработчик клика по маркеру для показа карточки
          leafletMarker.on('click', () => {
            // Если у маркера есть полные данные объявления
            if (marker.listing) {
              setSelectedListing(marker.listing);
            }
            // Если есть только ID, то можно запросить данные с сервера
            else if (marker.id) {
              // Перенаправляем на страницу объявления
              navigate(`/marketplace/listings/${marker.id}`);
            }
            // Если нет полных данных, но есть заголовок, показываем его в попапе
            else if (marker.title) {
              leafletMarker.bindPopup(marker.title).openPopup();
            }
          });

          leafletMarker.addTo(mapRef.current);
          markerGroup.addLayer(leafletMarker);
        });

        // Масштабируем карту, чтобы видеть все маркеры
        if (markers.length > 1) {
          mapRef.current.fitBounds(markerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
          });
        }
      } else if (title) {
        // Если маркеров нет, но есть центральная точка с заголовком
        const marker = L.marker([latitude, longitude]).addTo(mapRef.current);
        marker.bindPopup(title);
      }
    }

    // Очистка при размонтировании компонента
    return () => {
      if (mapRef.current) {
        try {
          mapRef.current.remove();
        } catch (err) {
          console.error("Ошибка при удалении карты:", err);
        }
        mapRef.current = null;
      }
    };
  }, [latitude, longitude, title, markers, hasCoordinates, navigate]);

  // Навигация к подробностям объявления
  const handleNavigateToListing = (listingId) => {
    navigate(`/marketplace/listings/${listingId}`);
  };

  return (
    <Paper
      sx={{
        position: 'relative',
        width: '100%',
        maxWidth: 1200,
        maxHeight: '90vh',
        overflow: 'hidden'
      }}
    >
      {hasCoordinates ? (
        <div
          ref={mapContainerRef}
          style={{ width: '100%', height: '80vh' }}
        />
      ) : (
        <div style={{
          width: '100%',
          height: '80vh',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          backgroundColor: '#f5f5f5'
        }}>
          <Typography variant="h6" color="error" gutterBottom>
            Координаты отсутствуют
          </Typography>
          {error && (
            <Typography variant="body2" color="text.secondary">
              {error}
            </Typography>
          )}
          <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
            Координаты: {latitude}, {longitude}
          </Typography>
        </div>
      )}

      {/* Карточка объявления при клике на маркер */}
      {selectedListing && (
        <ListingPreview
          listing={selectedListing}
          onClose={() => setSelectedListing(null)}
          onNavigate={handleNavigateToListing}
        />
      )}
    </Paper>
  );
};

export default FullscreenMap;
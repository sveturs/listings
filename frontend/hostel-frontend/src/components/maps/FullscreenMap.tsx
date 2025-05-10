// frontend/hostel-frontend/src/components/maps/FullscreenMap.tsx
import React, { useEffect, useRef, useState } from 'react';
import { Paper, Box, Typography, IconButton, Card, CardContent, CardMedia, Button, Chip } from '@mui/material';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import '../maps/leaflet-icons'; // Импортируем файл с фиксом иконок
import { X } from 'lucide-react';

interface Image {
  is_main?: boolean;
  public_url?: string;
  file_path?: string;
  storage_type?: string;
}

interface Listing {
  id: string | number;
  title: string;
  price: number;
  condition: 'new' | 'used';
  images?: Array<Image | string>;
}

interface ListingPreviewProps {
  listing: Listing;
  onClose: () => void;
  onNavigate: (id: string | number) => void;
}

interface Marker {
  latitude: number;
  longitude: number;
  title?: string;
  tooltip?: string;
  id?: string | number;
  listing?: Listing;
}

interface FullscreenMapProps {
  latitude: number;
  longitude: number;
  title?: string;
  markers?: Marker[];
}

// Компонент для предпросмотра объявления при клике по маркеру
const ListingPreview: React.FC<ListingPreviewProps> = ({ listing, onClose, onNavigate }) => {
  const { t } = useTranslation('marketplace');

  if (!listing) return null;

  const formatPrice = (price: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };

  const getImageUrl = (): string => {
    if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
      return '/placeholder-listing.jpg';
    }

    // Находим главное изображение или используем первое в списке
    let mainImage = listing.images.find(img => 
      typeof img === 'object' && img !== null && img.is_main === true
    ) || listing.images[0];

    // Используем переменную окружения из window.ENV вместо process.env
    const baseUrl = (window as any).ENV?.REACT_APP_MINIO_URL || (window as any).ENV?.REACT_APP_BACKEND_URL || '';
    console.log('FullscreenMap: Using baseUrl from env:', baseUrl);

    // 1. Строковые пути (для обратной совместимости)
    if (typeof mainImage === 'string') {
      console.log('FullscreenMap: Processing string image path:', mainImage);

      // Относительный путь MinIO
      if (mainImage.startsWith('/listings/')) {
        const url = `${baseUrl}${mainImage}`;
        console.log('FullscreenMap: Using MinIO relative path:', url);
        return url;
      }

      // ID/filename.jpg (прямой путь MinIO)
      if (mainImage.match(/^\d+\/[^\/]+$/)) {
        const url = `${baseUrl}/listings/${mainImage}`;
        console.log('FullscreenMap: Using direct MinIO path pattern:', url);
        return url;
      }

      // Локальное хранилище (обратная совместимость)
      const url = `${baseUrl}/uploads/${mainImage}`;
      console.log('FullscreenMap: Using local storage path:', url);
      return url;
    }

    // 2. Объекты с информацией о файле
    if (typeof mainImage === 'object' && mainImage !== null) {
      console.log('FullscreenMap: Processing image object:', mainImage);

      // Приоритет 1: Используем PublicURL если он доступен
      if (mainImage.public_url && typeof mainImage.public_url === 'string' && mainImage.public_url.trim() !== '') {
        const publicUrl = mainImage.public_url;
        console.log('FullscreenMap: Found public_url string:', publicUrl);

        // Абсолютный URL
        if (publicUrl.startsWith('http')) {
          console.log('FullscreenMap: Using absolute URL:', publicUrl);
          return publicUrl;
        }
        // Относительный URL с /listings/
        else if (publicUrl.startsWith('/listings/')) {
          const url = `${baseUrl}${publicUrl}`;
          console.log('FullscreenMap: Using public_url with listings path:', url);
          return url;
        }
        // Другой относительный URL
        else {
          const url = `${baseUrl}${publicUrl}`;
          console.log('FullscreenMap: Using general relative public_url:', url);
          return url;
        }
      }

      // Приоритет 2: Формируем URL на основе типа хранилища и пути к файлу
      if (mainImage.file_path) {
        if (mainImage.storage_type === 'minio' || mainImage.file_path.includes('listings/')) {
          // Учитываем возможность наличия префикса listings/ в пути
          const filePath = mainImage.file_path.includes('listings/')
            ? mainImage.file_path.replace('listings/', '') 
            : mainImage.file_path;

          const url = `${baseUrl}/listings/${filePath}`;
          console.log('FullscreenMap: Constructed MinIO URL from path:', url);
          return url;
        }

        // Локальное хранилище
        const url = `${baseUrl}/uploads/${mainImage.file_path}`;
        console.log('FullscreenMap: Using local storage path from object:', url);
        return url;
      }
    }

    console.log('FullscreenMap: Could not determine image URL, using placeholder');
    return '/placeholder-listing.jpg';
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

const FullscreenMap: React.FC<FullscreenMapProps> = ({ latitude, longitude, title, markers = [] }) => {
  const { t } = useTranslation('marketplace');
  const navigate = useNavigate();
  const mapContainerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<L.Map | null>(null);
  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Если координат нет, компонент всё равно отрендерится, но карта не будет инициализирована
  const hasCoordinates = Boolean(latitude && longitude);

  useEffect(() => {
    // Проверяем условия внутри хука
    if (!hasCoordinates || !mapContainerRef.current) return;

    // Добавляем небольшую задержку перед инициализацией карты, чтобы DOM успел полностью смонтироваться
    const initMapTimer = setTimeout(() => {
      // Инициализируем карту только если её еще нет
      if (!mapRef.current) {
        try {
          // Создаем карту с улучшенными параметрами для более плавного зума
          mapRef.current = L.map(mapContainerRef.current!, {
            zoomAnimation: true,         // Включаем анимацию зума
            fadeAnimation: true,         // Включаем анимацию прозрачности
            markerZoomAnimation: true,   // Анимация маркеров при зуме
            zoom: 15,                    // Начальный уровень масштабирования
            center: [latitude, longitude], // Начальный центр
            zoomSnap: 0.5,               // Более плавный зум (шаг масштабирования)
            zoomDelta: 0.5,              // Меньший шаг колесика мыши
            wheelDebounceTime: 40,       // Уменьшаем задержку колесика мыши
            wheelPxPerZoomLevel: 60,     // Больше пикселей на уровень масштаба
            preferCanvas: true,          // Использовать Canvas для лучшей производительности
            renderer: L.canvas()          // Явно указываем Canvas renderer
          });

          // Добавляем слой тайлов OpenStreetMap с улучшенными параметрами для плавности
          L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
            maxZoom: 19,
            subdomains: 'abc',        // Используем все поддомены для распределения запросов
            updateWhenIdle: true,     // Обновляем тайлы только когда карта неподвижна
            updateWhenZooming: false, // Не обновляем тайлы во время масштабирования
            keepBuffer: 4             // Сохраняем больше тайлов в буфере
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

              leafletMarker.addTo(mapRef.current!);
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
        } catch (error) {
          console.error("Ошибка при инициализации карты:", error);
          setError("Ошибка при инициализации карты");
        }
      } else {
        try {
          // Если карта уже существует, обновляем вид и маркеры
          console.log("Обновление карты с координатами:", [latitude, longitude]);
          // Обновляем центр и масштаб карты
          mapRef.current.setView([latitude, longitude], 15);

          // Сначала очищаем все существующие маркеры
          mapRef.current.eachLayer(layer => {
            if (layer instanceof L.Marker) {
              mapRef.current!.removeLayer(layer);
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

              leafletMarker.addTo(mapRef.current!);
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
        } catch (error) {
          console.error("Ошибка при обновлении карты:", error);
        }
      }
    }, 100); // Задержка в 100 мс для полной инициализации DOM

    // Очистка при размонтировании компонента
    return () => {
      clearTimeout(initMapTimer);
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
  const handleNavigateToListing = (listingId: string | number): void => {
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
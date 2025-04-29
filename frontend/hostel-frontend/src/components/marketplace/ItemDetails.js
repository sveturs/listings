// frontend/hostel-frontend/src/components/marketplace/ItemDetails.js
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Grid,
  Paper,
  Divider,
  CircularProgress,
  Container,
  Chip,
} from '@mui/material';
import { Favorite, Share, ShoppingCart } from '@mui/icons-material';
import axios from '../../api/axios';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000';

const ItemDetails = () => {
  const { id } = useParams();
  const [item, setItem] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  useEffect(() => {
    const fetchItemDetails = async () => {
      try {
        console.log("Fetching details for item ID:", id);
        const response = await axios.get(`/api/v1/marketplace/listings/${id}`);
        console.log("API response for item details:", response.data);
        // Используем правильное имя переменной (было ошибочно setListing)
        setItem(response.data.data);
      } catch (err) {
        console.error("Error fetching item details:", err);
        setError('Не удалось загрузить данные товара.');
      } finally {
        setLoading(false);
      }
    };
    fetchItemDetails();
  }, [id]);

  // Функция для получения URL изображения
  const getImageUrl = (image) => {
    if (!image) {
        return '/placeholder.png';
    }

    // Обработка строкового формата (обратная совместимость)
    if (typeof image === 'string') {
        return `${BACKEND_URL}/uploads/${image}`;
    }

    // Обработка объектного формата
    if (typeof image === 'object') {
        // Если есть публичный URL, используем его напрямую
        if (image.public_url) {
            return image.public_url;
        }
        
        // Для MinIO используем специальный путь
        if (image.storage_type === 'minio') {
            return `${BACKEND_URL}/listings/${image.file_path.split('/').pop()}`;
        }
        
        // Для локального хранилища - старый формат
        if (image.file_path) {
            return `${BACKEND_URL}/uploads/${image.file_path}`;
        }
    }

    return '/placeholder.png';
}


  // Переключение на предыдущее изображение
  const handlePrevImage = () => {
    if (item && item.images && item.images.length > 0) {
      setCurrentImageIndex((prev) => 
        prev > 0 ? prev - 1 : item.images.length - 1
      );
    }
  };

  // Переключение на следующее изображение
  const handleNextImage = () => {
    if (item && item.images && item.images.length > 0) {
      setCurrentImageIndex((prev) => 
        prev < item.images.length - 1 ? prev + 1 : 0
      );
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" mt={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box display="flex" justifyContent="center" mt={4}>
        <Typography color="error">{error}</Typography>
      </Box>
    );
  }

  // Получаем текущее изображение для отображения
  const currentImage = item && item.images && item.images.length > 0 
    ? item.images[currentImageIndex] 
    : null;

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Grid container spacing={4}>
        {/* Изображения */}
        <Grid item xs={12} md={6}>
          <Box 
            component={Paper} 
            elevation={3} 
            sx={{ 
              position: 'relative',
              height: 500, // Фиксированная высота для контейнера изображения
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              overflow: 'hidden'
            }}
          >
            <img
              src={currentImage ? getImageUrl(currentImage) : (item.image || '/placeholder.png')}
              alt={item.title}
              style={{ 
                maxWidth: '100%',
                maxHeight: '100%',
                objectFit: 'contain', // Важное свойство - изображение поместится полностью
                width: 'auto',
                height: 'auto'
              }}
            />
            
            {/* Кнопки для навигации между изображениями */}
            {item && item.images && item.images.length > 1 && (
              <>
                <Button 
                  variant="contained" 
                  sx={{ 
                    position: 'absolute', 
                    left: 10, 
                    backgroundColor: 'rgba(0,0,0,0.5)',
                    minWidth: '40px',
                    '&:hover': { backgroundColor: 'rgba(0,0,0,0.7)' }
                  }}
                  onClick={handlePrevImage}
                >
                  &lt;
                </Button>
                <Button 
                  variant="contained" 
                  sx={{ 
                    position: 'absolute', 
                    right: 10, 
                    backgroundColor: 'rgba(0,0,0,0.5)',
                    minWidth: '40px',
                    '&:hover': { backgroundColor: 'rgba(0,0,0,0.7)' }
                  }}
                  onClick={handleNextImage}
                >
                  &gt;
                </Button>
              </>
            )}
          </Box>
          
          {/* Миниатюры изображений */}
          {item && item.images && item.images.length > 1 && (
            <Box 
              sx={{ 
                display: 'flex', 
                gap: 1, 
                mt: 2, 
                overflowX: 'auto',
                py: 1
              }}
            >
              {item.images.map((image, index) => (
                <Box 
                  key={index}
                  component="img"
                  src={getImageUrl(image)}
                  alt={`Thumbnail ${index}`}
                  sx={{
                    height: 80,
                    width: 80,
                    objectFit: 'cover',
                    cursor: 'pointer',
                    borderRadius: 1,
                    border: currentImageIndex === index ? '2px solid #1976d2' : '2px solid transparent',
                    opacity: currentImageIndex === index ? 1 : 0.7,
                    transition: 'all 0.2s',
                    '&:hover': { opacity: 1 }
                  }}
                  onClick={() => setCurrentImageIndex(index)}
                />
              ))}
            </Box>
          )}
        </Grid>
        
        {/* Основная информация */}
        <Grid item xs={12} md={6}>
          <Box>
            <Typography variant="h4" fontWeight="bold">
              {item.title}
            </Typography>
            <Typography variant="h6" color="text.secondary" gutterBottom>
              {item.price} RSD
            </Typography>
            <Typography variant="body1" sx={{ my: 2 }}>
              {item.description}
            </Typography>
            <Divider sx={{ my: 2 }} />
            <Box display="flex" gap={2} flexWrap="wrap">
              <Button
                id="buyButton"
                variant="contained"
                startIcon={<ShoppingCart />}
                size="large"
              >
                Купить
              </Button>
              <Button
                id="bookmarkButton"
                variant="outlined"
                startIcon={<Favorite />}
                size="large"
              >
                В избранное
              </Button>
              <Button
                id="shareButton"
                variant="outlined"
                startIcon={<Share />}
                size="large"
              >
                Поделиться
              </Button>
            </Box>
          </Box>
        </Grid>
      </Grid>
      
      {/* Отзывы */}
      <Box mt={4}>
        <Typography variant="h5" fontWeight="bold">
          Отзывы
        </Typography>
        {/* Здесь можно добавить список отзывов */}
        <Typography variant="body2" color="text.secondary">
          Пока нет отзывов
        </Typography>
      </Box>
    </Container>
  );
};

export default ItemDetails;
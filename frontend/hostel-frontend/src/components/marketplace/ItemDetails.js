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
  Chip,
} from '@mui/material';
import { Favorite, Share, ShoppingCart } from '@mui/icons-material';
import axios from '../api/axios';

const ItemDetails = () => {
  const { id } = useParams();
  const [item, setItem] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchItemDetails = async () => {
      try {
        const response = await axios.get(`/api/v1/marketplace/listings/${id}`);
setListing(response.data.data);
      } catch (err) {
        setError('Не удалось загрузить данные товара.');
      } finally {
        setLoading(false);
      }
    };

    fetchItemDetails();
  }, [id]);

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

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Grid container spacing={4}>
        {/* Изображения */}
        <Grid item xs={12} md={6}>
          <Box component={Paper} elevation={3} sx={{ overflow: 'hidden' }}>
            <img
              src={item.image || '/placeholder.png'}
              alt={item.title}
              style={{ width: '100%', height: 'auto' }}
            />
          </Box>
        </Grid>

        {/* Основная информация */}
        <Grid item xs={12} md={6}>
          <Box>
            <Typography variant="h4" fontWeight="bold">
              {item.title}
            </Typography>
            <Typography variant="h6" color="text.secondary" gutterBottom>
              {item.price} ₽
            </Typography>
            <Typography variant="body1" sx={{ my: 2 }}>
              {item.description}
            </Typography>
            <Divider sx={{ my: 2 }} />
            <Box display="flex" gap={2}>
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
          Пока нет отзывов.
        </Typography>
      </Box>
    </Container>
  );
};

export default ItemDetails;

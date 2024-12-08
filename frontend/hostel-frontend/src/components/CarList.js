import React from 'react';
import {
  Grid,
  Card,
  CardContent,
  CardMedia,
  Typography,
  Box,
  Chip,
  Skeleton
} from '@mui/material';
import {
  DirectionsCar,
  LocalGasStation,
  AirlineSeatReclineNormal,
  Settings
} from '@mui/icons-material';
import axios from '../api/axios';

const CarList = () => {
  const [cars, setCars] = React.useState([]);
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    const fetchCars = async () => {
      try {
        const response = await axios.get('/api/v1/cars');
        // Берем только последние 8 машин
        setCars(response.data.data.slice(0, 8));
      } catch (error) {
        console.error('Error fetching cars:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchCars();
  }, []);

  const getTransmissionLabel = (trans) => {
    return trans === 'automatic' ? 'Автомат' : 'Механика';
  };

  const getFuelTypeLabel = (fuel) => {
    const types = {
      petrol: 'Бензин',
      diesel: 'Дизель',
      electric: 'Электро',
      hybrid: 'Гибрид'
    };
    return types[fuel] || fuel;
  };

  if (loading) {
    return (
      <Grid container spacing={3}>
        {[...Array(8)].map((_, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Skeleton variant="rectangular" height={200} />
            <Skeleton height={30} sx={{ mt: 1 }} />
            <Skeleton height={20} width="60%" />
          </Grid>
        ))}
      </Grid>
    );
  }

  return (
    <Grid container spacing={3}>
      {cars.map((car) => (
        <Grid item xs={12} sm={6} md={3} key={car.id}>
          <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <CardMedia
              component="img"
              height="200"
              image={car.images?.[0]?.file_path || '/placeholder-car.jpg'}
              alt={`${car.make} ${car.model}`}
              sx={{ objectFit: 'cover' }}
            />
            <CardContent sx={{ flexGrow: 1 }}>
              <Typography variant="h6" component="div" gutterBottom noWrap>
                {car.make} {car.model}
              </Typography>
              
              <Box sx={{ display: 'flex', gap: 1, mb: 1, flexWrap: 'wrap' }}>
                <Chip
                  size="small"
                  icon={<Settings sx={{ fontSize: 16 }} />}
                  label={getTransmissionLabel(car.transmission)}
                />
                <Chip
                  size="small"
                  icon={<LocalGasStation sx={{ fontSize: 16 }} />}
                  label={getFuelTypeLabel(car.fuel_type)}
                />
                <Chip
                  size="small"
                  icon={<AirlineSeatReclineNormal sx={{ fontSize: 16 }} />}
                  label={`${car.seats} мест`}
                />
              </Box>

              <Typography variant="h6" color="primary" sx={{ mt: 2 }}>
                {car.price_per_day.toLocaleString()} ₽/день
              </Typography>

              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }} noWrap>
                {car.location}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};

export default CarList;
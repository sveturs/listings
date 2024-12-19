import React, { useState, useEffect, useCallback } from 'react';
import { LoadScript } from '@react-google-maps/api';
import MapView from '../../components/car/CarMapView';
import CarBookingDialog from '../../components/car/CarBookingDialog';
import CarDetailsDialog from '../../components/car/CarDetailsDialog';
import {
  Box,
  Button,
  Card,
  Rating,
  CardContent,
  CardMedia,
  Chip,
  Container,
  Divider,
  Grid,
  IconButton,
  MenuItem,
  Paper,
  Skeleton,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
  useTheme,
  useMediaQuery
} from '@mui/material';
import {
  Map as MapIcon,
  ViewList as ListIcon,
  Search as SearchIcon,
  LocalGasStation as FuelIcon,
  Speed as TransmissionIcon,
  AirlineSeatReclineNormal as SeatsIcon,
  CalendarMonth as CalendarIcon,
  PhotoLibrary as GalleryIcon,
  LocationOn as LocationIcon,
  Tune as FilterIcon,
  Clear as ClearIcon
} from '@mui/icons-material';
import { debounce } from 'lodash';

import axios from '../../api/axios';


const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

export default function CarListPage() {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const isTablet = useMediaQuery(theme.breakpoints.down('md'));

  const [viewMode, setViewMode] = useState('list');
  const [selectedCarDetails, setSelectedCarDetails] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showFilters, setShowFilters] = useState(!isMobile);
  const [cars, setCars] = useState([]);
  const [selectedCar, setSelectedCar] = useState(null);
  const [bookingDialogOpen, setBookingDialogOpen] = useState(false);
  const [filters, setFilters] = useState({
    start_date: '',
    end_date: '',
    make: '',
    transmission: '',
    fuel_type: '',
    min_price: '',
    max_price: ''
  });

  const fetchCars = useCallback(async () => {
    try {
      setIsLoading(true);
      const params = Object.fromEntries(
        Object.entries(filters).filter(([_, v]) => v !== '')
      );
      const { data } = await axios.get('/api/v1/cars/available', { params });
      setCars(data?.data || []);
    } catch (error) {
      console.error('Error fetching cars:', error);
    } finally {
      setIsLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchCars();
  }, []);

  const debouncedFetch = debounce(fetchCars, 500);

  const handleFilterChange = (name, value) => {
    setFilters(prev => ({ ...prev, [name]: value }));
    debouncedFetch();
  };

  const clearFilters = () => {
    setFilters({
      start_date: '',
      end_date: '',
      make: '',
      transmission: '',
      fuel_type: '',
      min_price: '',
      max_price: ''
    });
    fetchCars();
  };

  const renderFilters = () => (
    <Paper
      elevation={0}
      sx={{
        p: 2,
        border: 1,
        borderColor: 'divider',
        display: showFilters ? 'block' : 'none'
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
        <Typography variant="h6">Фильтры</Typography>
        <Button
          size="small"
          startIcon={<ClearIcon />}
          onClick={clearFilters}
        >
          Сбросить
        </Button>
      </Box>

      <Grid container spacing={2}>
        <Grid item xs={12} md={6}>
          <TextField
            label="Дата начала"
            type="date"
            fullWidth
            size="small"
            value={filters.start_date}
            onChange={(e) => handleFilterChange('start_date', e.target.value)}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>

        <Grid item xs={12} md={6}>
          <TextField
            label="Дата окончания"
            type="date"
            fullWidth
            size="small"
            value={filters.end_date}
            onChange={(e) => handleFilterChange('end_date', e.target.value)}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>

        <Grid item xs={12} md={6}>
          <TextField
            select
            label="Тип топлива"
            fullWidth
            size="small"
            value={filters.fuel_type}
            onChange={(e) => handleFilterChange('fuel_type', e.target.value)}
          >
            <MenuItem value="">Все</MenuItem>
            <MenuItem value="petrol">Бензин</MenuItem>
            <MenuItem value="diesel">Дизель</MenuItem>
            <MenuItem value="electric">Электро</MenuItem>
            <MenuItem value="hybrid">Гибрид</MenuItem>
          </TextField>
        </Grid>

        <Grid item xs={12} md={6}>
          <TextField
            select
            label="Коробка передач"
            fullWidth
            size="small"
            value={filters.transmission}
            onChange={(e) => handleFilterChange('transmission', e.target.value)}
          >
            <MenuItem value="">Все</MenuItem>
            <MenuItem value="automatic">Автомат</MenuItem>
            <MenuItem value="manual">Механика</MenuItem>
          </TextField>
        </Grid>

        <Grid item xs={12} md={6}>
          <TextField
            label="Минимальная цена"
            type="number"
            fullWidth
            size="small"
            value={filters.min_price}
            onChange={(e) => handleFilterChange('min_price', e.target.value)}
          />
        </Grid>

        <Grid item xs={12} md={6}>
          <TextField
            label="Максимальная цена"
            type="number"
            fullWidth
            size="small"
            value={filters.max_price}
            onChange={(e) => handleFilterChange('max_price', e.target.value)}
          />
        </Grid>
      </Grid>
    </Paper>
  );

  const renderCarCard = (car) => (
    <Card
      onClick={() => setSelectedCarDetails(car)}
      elevation={0}
      sx={{
        cursor: 'pointer',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        border: 1,
        borderColor: 'divider',
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: theme.shadows[4]
        }
      }}
    >
      <Box sx={{ position: 'relative', pt: '56.25%' }}>
        <CardMedia
          component="img"
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            objectFit: 'cover'
          }}
          image={car.images?.[0] ?
            `${BACKEND_URL}/uploads/${car.images[0].file_path}` :
            '/placeholder-car.jpg'
          }
          alt={`${car.make} ${car.model}`}
        />
        {car.images?.length > 1 && (
          <Button
            variant="contained"
            size="small"
            startIcon={<GalleryIcon />}
            sx={{
              position: 'absolute',
              bottom: 8,
              right: 8,
              bgcolor: 'rgba(0, 0, 0, 0.7)'
            }}
          >
            {car.images.length} фото
          </Button>
        )}
      </Box>

      <CardContent sx={{ flexGrow: 1, p: 2 }}>
        <Typography variant="h6" gutterBottom>
          {car.make} {car.model}
          <Typography
            component="span"
            color="text.secondary"
            sx={{ ml: 1 }}
          >
            {car.year}
          </Typography>
        </Typography>

        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 2 }}>
          <Chip
            icon={<FuelIcon />}
            label={
              car.fuel_type === 'petrol' ? 'Бензин' :
                car.fuel_type === 'diesel' ? 'Дизель' :
                  car.fuel_type === 'electric' ? 'Электро' : 'Гибрид'
            }
            size="small"
          />
          <Chip
            icon={<TransmissionIcon />}
            label={car.transmission === 'automatic' ? 'Автомат' : 'Механика'}
            size="small"
          />
          <Chip
            icon={<SeatsIcon />}
            label={`${car.seats} мест`}
            size="small"
          />
        </Box>

        <Typography variant="body2" color="text.secondary" gutterBottom>
          <LocationIcon sx={{ fontSize: 16, mr: 0.5, verticalAlign: 'text-bottom' }} />
          {car.location}
        </Typography>

        {car.features?.length > 0 && (
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden',
              mb: 2
            }}
          >
            {car.features.join(' • ')}
          </Typography>
        )}
        {car.rating > 0 && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
            <Rating value={car.rating} precision={0.1} readOnly size="small" />
            <Typography variant="body2" color="text.secondary">
              {car.reviews_count} {car.reviews_count === 1 ? 'отзыв' :
                car.reviews_count < 5 ? 'отзыва' : 'отзывов'}
            </Typography>
          </Box>
        )}
      </CardContent>

      <Divider />

      <Box sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <Box>
          <Typography variant="h6" color="primary">
            {car.price_per_day} ₽
          </Typography>
          <Typography variant="caption" color="text.secondary">
            в день
          </Typography>
        </Box>
        <Button
          variant="contained"
          onClick={(e) => {
            e.stopPropagation(); // Предотвращаем открытие CarDetailsDialog
            setSelectedCar(car);
            setBookingDialogOpen(true);
          }}
          disabled={!filters.start_date || !filters.end_date}
        >
          Забронировать
        </Button>
      </Box>
    </Card>
  );

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <Box sx={{ mb: 4 }}>
        <Box sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          mb: 3,
          flexWrap: 'wrap',
          gap: 2
        }}>
          <Typography variant="h4" component="h1">
            Аренда автомобилей
          </Typography>

          <Box sx={{ display: 'flex', gap: 2 }}>
            <ToggleButtonGroup
              value={viewMode}
              exclusive
              onChange={(_, mode) => mode && setViewMode(mode)}
              size="small"
            >
              <ToggleButton value="list">
                <ListIcon sx={{ mr: 1 }} /> Список
              </ToggleButton>
              <ToggleButton value="map">
                <MapIcon sx={{ mr: 1 }} /> Карта
              </ToggleButton>
            </ToggleButtonGroup>

            {isMobile && (
              <Button
                variant="outlined"
                onClick={() => setShowFilters(!showFilters)}
                startIcon={<FilterIcon />}
              >
                Фильтры
              </Button>
            )}
          </Box>
        </Box>

        {renderFilters()}
      </Box>

      {isLoading ? (
        <Grid container spacing={3}>
          {[1, 2, 3, 4, 5, 6].map(i => (
            <Grid item xs={12} sm={6} md={4} key={i}>
              <Skeleton variant="rectangular" height={200} />
              <Box sx={{ pt: 0.5 }}>
                <Skeleton />
                <Skeleton width="60%" />
              </Box>
            </Grid>
          ))}
        </Grid>
      ) : viewMode === 'list' ? (
        <Grid container spacing={3}>
          {cars.map(car => (
            <Grid item xs={12} sm={6} md={4} key={car.id}>
              {renderCarCard(car)}
            </Grid>
          ))}
        </Grid>
      ) : (
        <LoadScript googleMapsApiKey={process.env.REACT_APP_GOOGLE_MAPS_API_KEY}>
          <MapView
            cars={cars}
            onCarSelect={(car) => {
              setSelectedCar(car);
              setBookingDialogOpen(true);
            }}
            onViewDetails={(car) => setSelectedCarDetails(car)} // Добавьте этот обработчик
            onOpenGallery={(car) => setSelectedCarDetails(car)}
          />
        </LoadScript>
      )}

      {selectedCar && (
        <CarBookingDialog
          open={bookingDialogOpen}
          onClose={() => {
            setBookingDialogOpen(false);
            setSelectedCar(null);
          }}
          car={selectedCar}
          startDate={filters.start_date}
          endDate={filters.end_date}
        />
      )}

      {/* Добавляем диалог с подробной информацией */}
      <CarDetailsDialog
        open={Boolean(selectedCarDetails)}
        onClose={() => setSelectedCarDetails(null)}
        car={selectedCarDetails}
        onBook={(car) => {
          setSelectedCarDetails(null);
          setSelectedCar(car);
          setBookingDialogOpen(true);
        }}
      />
    </Container>
  );
}
import React, { useState, useEffect, useCallback } from 'react';
import { LoadScript } from '@react-google-maps/api';
import MapView from '../components/CarMapView';
import CarBookingDialog from '../components/CarBookingDialog';
import {
    Container,
    Grid,
    Card,
    CardContent,
    Typography,
    TextField,
    Button,
    Box,
    CardMedia,
    Chip,
    ToggleButtonGroup,
    ToggleButton,
    Divider,
    Paper,
    MenuItem,
    CircularProgress
} from '@mui/material';
import { debounce } from 'lodash';
import {
    ViewList as ViewListIcon,
    Map as MapIcon,
    Search as SearchIcon,
    LocalGasStation as FuelIcon,
    Speed as TransmissionIcon,
    LocationOn as LocationIcon,
    PhotoLibrary as GalleryIcon,
} from '@mui/icons-material';
import axios from '../api/axios';

const CarListPage = () => {
    const [cars, setCars] = useState([]);
    const [filteredCars, setFilteredCars] = useState([]);
    const [viewMode, setViewMode] = useState('list');
    const [filters, setFilters] = useState({
        category: '',
        make: '',
        model: '',
        transmission: '',
        fuel_type: '',
        min_price: '',
        max_price: '',
        location: '',
        start_date: '',
        end_date: '',
    });
    const [selectedCar, setSelectedCar] = useState(null);
    const [galleryOpen, setGalleryOpen] = useState(false);
    const [bookingDialogOpen, setBookingDialogOpen] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    // Функция загрузки автомобилей
    const fetchCars = useCallback(async (currentFilters) => {
        try {
            setIsLoading(true);
            const params = Object.entries(currentFilters).reduce((acc, [key, value]) => {
                if (value !== '' && value !== null && value !== undefined) {
                    acc[key] = value;
                }
                return acc;
            }, {});

            const response = await axios.get('/api/v1/cars/available', { params });
            const carsData = response.data?.data || [];
            setCars(carsData);
            setFilteredCars(carsData);
        } catch (error) {
            console.error('Ошибка при получении списка автомобилей:', error);
            setCars([]);
            setFilteredCars([]);
        } finally {
            setIsLoading(false);
        }
    }, []);

    // Загрузка автомобилей при монтировании компонента
    useEffect(() => {
        fetchCars(filters);
    }, []); // Пустой массив зависимостей для загрузки только при монтировании

    // Отложенный поиск при изменении фильтров
    const debouncedSearch = useCallback(
        debounce((filters) => {
            fetchCars(filters);
        }, 500),
        [fetchCars]
    );

    // Обработчик изменения фильтров
    const handleFilterChange = (newFilters) => {
        setFilters(prev => {
            const updatedFilters = { ...prev, ...newFilters };
            debouncedSearch(updatedFilters);
            return updatedFilters;
        });
    };

    // Обработчик кнопки поиска
    const handleSearch = () => {
        fetchCars(filters);
    };

    const handleBooking = (car) => {
        if (!filters.start_date || !filters.end_date) {
            alert('Пожалуйста, выберите даты аренды');
            return;
        }
        setSelectedCar(car);
        setBookingDialogOpen(true);
    };
    return (
        <Container maxWidth="xl">
            {isLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                    <CircularProgress />
                </Box>
            ) : (
                <>
                    <Paper sx={{ p: 2, mt: 4, mb: 3 }}>
                        <Grid container spacing={2}>
                            <Grid item xs={12} sm={6} md={2}>
                                <TextField
                                    label="Дата начала"
                                    type="date"
                                    size="small"
                                    fullWidth
                                    InputLabelProps={{ shrink: true }}
                                    value={filters.start_date}
                                    onChange={(e) => setFilters({ ...filters, start_date: e.target.value })}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={2}>
                                <TextField
                                    label="Дата окончания"
                                    type="date"
                                    size="small"
                                    fullWidth
                                    InputLabelProps={{ shrink: true }}
                                    value={filters.end_date}
                                    onChange={(e) => setFilters({ ...filters, end_date: e.target.value })}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={2}>
                                <TextField
                                    label="Марка"
                                    size="small"
                                    fullWidth
                                    value={filters.make}
                                    onChange={(e) => setFilters({ ...filters, make: e.target.value })}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={2}>
                                <TextField
                                    select
                                    label="Тип топлива"
                                    size="small"
                                    fullWidth
                                    value={filters.fuel_type}
                                    onChange={(e) => setFilters({ ...filters, fuel_type: e.target.value })}
                                >
                                    <MenuItem value="">Все</MenuItem>
                                    <MenuItem value="petrol">Бензин</MenuItem>
                                    <MenuItem value="diesel">Дизель</MenuItem>
                                    <MenuItem value="electric">Электро</MenuItem>
                                    <MenuItem value="hybrid">Гибрид</MenuItem>
                                </TextField>
                            </Grid>
                            <Grid item xs={12} sm={6} md={2}>
                                <TextField
                                    select
                                    label="Коробка передач"
                                    size="small"
                                    fullWidth
                                    value={filters.transmission}
                                    onChange={(e) => setFilters({ ...filters, transmission: e.target.value })}
                                >
                                    <MenuItem value="">Все</MenuItem>
                                    <MenuItem value="automatic">Автомат</MenuItem>
                                    <MenuItem value="manual">Механика</MenuItem>
                                </TextField>
                            </Grid>
                            <Grid item xs={12} sm={6} md={2}>
                                <Button
                                    variant="contained"
                                    fullWidth
                                    onClick={handleSearch}
                                    startIcon={<SearchIcon />}
                                    disabled={isLoading}
                                >
                                    {isLoading ? 'Поиск...' : 'Найти'}
                                </Button>
                            </Grid>
                        </Grid>
                    </Paper>

                    <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 2 }}>
                        <ToggleButtonGroup
                            value={viewMode}
                            exclusive
                            onChange={(e, newMode) => newMode && setViewMode(newMode)}
                            size="small"
                        >
                            <ToggleButton value="list">
                                <ViewListIcon sx={{ mr: 1 }} /> Список
                            </ToggleButton>
                            <ToggleButton value="map">
                                <MapIcon sx={{ mr: 1 }} /> Карта
                            </ToggleButton>
                        </ToggleButtonGroup>
                    </Box>

                    {viewMode === 'list' ? (
                        <Grid container spacing={3}>
                            {filteredCars.map((car) => (
                                <Grid item xs={12} sm={6} md={4} key={car.id}>
                                    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                                        <Box sx={{ position: 'relative', pt: '56.25%' }}>
                                            <CardMedia
                                                component="img"
                                                sx={{
                                                    position: 'absolute',
                                                    top: 0,
                                                    left: 0,
                                                    width: '100%',
                                                    height: '100%',
                                                    objectFit: 'cover',
                                                }}
                                                image={car.images?.length > 0 ?
                                                    `${process.env.REACT_APP_BACKEND_URL}/uploads/${car.images[0].file_path}` :
                                                    '/placeholder-car.jpg'}
                                                alt={`${car.make} ${car.model}`}
                                                onClick={() => {
                                                    if (car.images?.length) {
                                                        setSelectedCar(car);
                                                        setGalleryOpen(true);
                                                    }
                                                }}
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
                                                        bgcolor: 'rgba(0, 0, 0, 0.7)',
                                                    }}
                                                    onClick={() => {
                                                        setSelectedCar(car);
                                                        setGalleryOpen(true);
                                                    }}
                                                >
                                                    {car.images.length} фото
                                                </Button>
                                            )}
                                        </Box>

                                        <CardContent sx={{ flexGrow: 1 }}>
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

                                            <Box sx={{ mb: 2 }}>
                                                <Chip
                                                    icon={<FuelIcon />}
                                                    label={car.fuel_type === 'petrol' ? 'Бензин' :
                                                        car.fuel_type === 'diesel' ? 'Дизель' :
                                                            car.fuel_type === 'electric' ? 'Электро' : 'Гибрид'}
                                                    size="small"
                                                    sx={{ mr: 1, mb: 1 }}
                                                />
                                                <Chip
                                                    icon={<TransmissionIcon />}
                                                    label={car.transmission === 'automatic' ? 'Автомат' : 'Механика'}
                                                    size="small"
                                                    sx={{ mr: 1, mb: 1 }}
                                                />
                                                <Chip
                                                    label={`${car.seats} мест`}
                                                    size="small"
                                                    sx={{ mb: 1 }}
                                                />
                                            </Box>

                                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                                <LocationIcon sx={{ fontSize: 16, mr: 0.5, verticalAlign: 'text-bottom' }} />
                                                {car.location}
                                            </Typography>

                                            {car.features?.length > 0 && (
                                                <Typography variant="body2" color="text.secondary">
                                                    {car.features.join(' • ')}
                                                </Typography>
                                            )}
                                        </CardContent>

                                        <Divider />

                                        <Box sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                                            <Typography variant="h6" color="primary">
                                                {car.price_per_day} ₽
                                                <Typography component="span" variant="body2" color="text.secondary">
                                                    /день
                                                </Typography>
                                            </Typography>
                                            <Button
                                                variant="contained"
                                                onClick={() => handleBooking(car)}
                                                disabled={!filters.start_date || !filters.end_date}
                                            >
                                                Забронировать
                                            </Button>
                                        </Box>
                                    </Card>
                                </Grid>
                            ))}
                        </Grid>
                    ) : (
                        <LoadScript googleMapsApiKey={process.env.REACT_APP_GOOGLE_MAPS_API_KEY}>
                            <MapView
                                cars={filteredCars}
                                onCarSelect={(car) => {
                                    setSelectedCar(car);
                                    setBookingDialogOpen(true);
                                }}
                                onOpenGallery={(car) => {
                                    setSelectedCar(car);
                                    setGalleryOpen(true);
                                }}
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
                </>
            )}
        </Container>
    );
};

export default CarListPage;
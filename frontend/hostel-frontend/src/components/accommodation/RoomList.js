//hostel-booking-system/frontend/hostel-frontend/src/components/accommodation/RoomList.js
import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { debounce } from 'lodash';
import MapView from './MapView';
import BookingDialog from './BookingDialog';
import RoomDetailsDialog from './RoomDetailsDialog';
import AdvancedFilters from './AdvancedFilters';
import {
    Grid,
    Card,
    CardContent,
    Typography,
    Box,
    Button,
    ToggleButtonGroup,
    ToggleButton,
    Chip,
    CardMedia,
    Divider,
    CircularProgress,
    Alert,
} from '@mui/material';
import {
    ViewList as ViewListIcon,
    Map as MapIcon,
    LocationOn as LocationIcon,
    PhotoLibrary as GalleryIcon,
    SingleBed as SingleBedIcon,
    Hotel as HotelIcon,
    Apartment as ApartmentIcon,
} from '@mui/icons-material';
import axios from '../../api/axios';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const RoomList = () => {
    const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
    const [rooms, setRooms] = useState([]);
    const [viewMode, setViewMode] = useState('list');
    const [selectedRoom, setSelectedRoom] = useState(null);
    const [galleryOpen, setGalleryOpen] = useState(false);
    const [bookingDialogOpen, setBookingDialogOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const [totalCount, setTotalCount] = useState(0);
    const [filters, setFilters] = useState({
        start_date: '',
        end_date: '',
        city: '',
        country: '',
        min_price: '',
        max_price: '',
        capacity: '',
        accommodation_type: '',
        sort_by: 'created_at',
        sort_direction: 'desc'
    });
    const handleRoomClick = (room) => {
        setSelectedRoom(room);
        setDetailsDialogOpen(true);
    };

    const fetchRooms = useCallback(async (isLoadMore = false) => {
        try {
            setLoading(true);
            setError(null);

            const params = new URLSearchParams();

            // Добавляем только непустые фильтры
            Object.entries(filters).forEach(([key, value]) => {
                if (value) params.append(key, value);
            });

            // Добавляем параметры пагинации
            params.append('page', isLoadMore ? page + 1 : 1);
            params.append('limit', 12);

            const response = await axios.get(`${BACKEND_URL}/rooms`, { params });

            const { data, meta } = response.data.data;
            setRooms(prev => isLoadMore ? [...prev, ...data] : data);
            setTotalCount(meta.total);
            setHasMore(meta.has_more);

            if (isLoadMore) {
                setPage(p => p + 1);
            } else {
                setPage(1);
            }

        } catch (error) {
            console.error("Ошибка при получении списка комнат:", error);
            setError(error.response?.data?.error || "Ошибка при загрузке комнат");
            setRooms([]); 
        } finally {
            setLoading(false);
        }
    }, [filters, page]);

    // Используем debounce для fetchRooms
    const debouncedFetchRooms = useMemo(
        () => debounce(fetchRooms, 300),
        [fetchRooms]
    );

    useEffect(() => {
        // Вызываем функцию загрузки данных
        debouncedFetchRooms();
        
        // Возвращаем функцию очистки, которая будет вызвана при размонтировании компонента
        return () => {
            debouncedFetchRooms.cancel();
        };
    }, [debouncedFetchRooms]); // Зависимость указывает, что эффект должен перезапускаться при изменении debouncedFetchRooms

    const handleFilterChange = (newFilters) => {
        setFilters(newFilters);
    };

    const handleSortChange = (sortBy, direction) => {
        setFilters(prev => ({
            ...prev,
            sort_by: sortBy,
            sort_direction: direction
        }));
    };

    const handleLoadMore = () => {
        if (!loading && hasMore) {
            fetchRooms(true);
        }
    };

    const handleBooking = (room) => {
        if (!filters.start_date || !filters.end_date) {
            setError('Пожалуйста, выберите даты заезда и выезда');
            return;
        }
        setSelectedRoom(room);
        setBookingDialogOpen(true);
    };

    const renderRoomCard = (room) => (
        <Card
            onClick={() => handleRoomClick(room)}
            sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                '&:hover': {
                    boxShadow: 6,
                    transform: 'translateY(-2px)',
                    transition: 'all 0.2s ease-in-out'
                }
            }}>
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
                        cursor: 'pointer',
                    }}
                    image={room.images?.[0] ?
                        `${BACKEND_URL}/uploads/${room.images[0].file_path}` :
                        '/placeholder-room.jpg'}
                    alt={room.name}
                />
                {room.images?.length > 1 && (
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
                    >
                        {room.images.length} фото
                    </Button>
                )}
            </Box>
    
            <CardContent sx={{ flexGrow: 1 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    {room.accommodation_type === 'bed' ? (
                        <SingleBedIcon color="primary" sx={{ mr: 1 }} />
                    ) : room.accommodation_type === 'apartment' ? (
                        <ApartmentIcon color="primary" sx={{ mr: 1 }} />
                    ) : (
                        <HotelIcon color="primary" sx={{ mr: 1 }} />
                    )}
                    <Typography variant="h6" noWrap>
                        {room.name}
                    </Typography>
                </Box>
    
                <Typography variant="body2" color="text.secondary" sx={{
                    display: 'flex',
                    alignItems: 'center',
                    mb: 1
                }}>
                    <LocationIcon sx={{ fontSize: 18, mr: 0.5 }} />
                    {room.address_street}, {room.address_city}
                </Typography>
    
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 2 }}>
                    <Chip
                        label={`${room.price_per_night} ₽/ночь`}
                        color="primary"
                        size="small"
                    />
                    {room.accommodation_type === 'bed' && (
                        <Chip
                            label={`${room.available_beds}/${room.total_beds} мест`}
                            color="secondary"
                            size="small"
                        />
                    )}
                    {room.rating > 0 && (
                        <Chip
                            label={`Рейтинг: ${room.rating.toFixed(1)}`}
                            size="small"
                            color="default"
                        />
                    )}
                </Box>
            </CardContent>
    
            <Divider />
    
            <Box sx={{ p: 2 }}>
                <Button
                    fullWidth
                    variant="contained"
                    onClick={(e) => {
                        e.stopPropagation();
                        handleBooking(room);
                    }}
                    disabled={!filters.start_date || !filters.end_date}
                >
                    Забронировать
                </Button>
            </Box>
        </Card>
    );

    return (
        <Box>
            <AdvancedFilters
                initialFilters={filters}
                onFilterChange={handleFilterChange}
                onSortChange={handleSortChange}
                isLoading={loading}
            />

            {error && (
                <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="body2" color="text.secondary">
                    Найдено: {totalCount}
                </Typography>

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

            {loading && page === 1 ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
                    <CircularProgress />
                </Box>
            ) : viewMode === 'list' ? (
                <>
                    <Grid container spacing={3}>
                        {Array.isArray(rooms) && rooms.map((room) => (
                            <Grid item xs={12} sm={6} md={4} key={room.id}>
                                {renderRoomCard(room)}
                            </Grid>
                        ))}
                    </Grid>

                    {hasMore && (
                        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                            <Button
                                variant="outlined"
                                onClick={handleLoadMore}
                                disabled={loading}
                                startIcon={loading && <CircularProgress size={20} />}
                            >
                                {loading ? 'Загрузка...' : 'Загрузить еще'}
                            </Button>
                        </Box>
                    )}
                </>
            ) : (
                <MapView
                    rooms={rooms}
                    onRoomSelect={handleBooking}
                    onOpenGallery={(room) => {
                        setSelectedRoom(room);
                        setGalleryOpen(true);
                    }}
                />
            )}

            {selectedRoom && (
                <>
                    <RoomDetailsDialog
                        open={detailsDialogOpen}
                        onClose={() => setDetailsDialogOpen(false)}
                        room={selectedRoom}
                        onBook={handleBooking}
                    />
                    <BookingDialog
                        open={bookingDialogOpen}
                        onClose={() => {
                            setBookingDialogOpen(false);
                            setSelectedRoom(null);
                        }}
                        room={selectedRoom}
                        startDate={filters.start_date}
                        endDate={filters.end_date}
                    />
                </>
            )}
        </Box>
    );
};

export default RoomList;
import React, { useState, useEffect, useCallback } from 'react';
import axios from "../api/axios";
import { Chip } from '@mui/material';
import { LoadScript } from '@react-google-maps/api';
import {
    Grid, Card, CardContent, Typography, TextField,
    Button, Divider, Box, Dialog, DialogContent, IconButton,
    MobileStepper, CardMedia,
    ToggleButton,
    ToggleButtonGroup,
    Paper
} from "@mui/material";
import {
    KeyboardArrowLeft, KeyboardArrowRight,
    Close as CloseIcon,
    SingleBed as SingleBedIcon,
    Hotel as HotelIcon,
    Apartment as ApartmentIcon,
    Home as HomeIcon,
    ViewList as ViewListIcon,
    Map as MapIcon
} from '@mui/icons-material';
import BookingDialog from "./BookingDialog";
import MapView from './MapView';


const BACKEND_URL = 'http://localhost:3000';

const ImageGallery = ({ images, open, onClose }) => {
    const [activeStep, setActiveStep] = useState(0);
    const maxSteps = images.length;

    const handleNext = () => {
        setActiveStep((prevStep) => (prevStep + 1) % maxSteps);
    };

    const handleBack = () => {
        setActiveStep((prevStep) => (prevStep - 1 + maxSteps) % maxSteps);
    };

    if (!images.length) return null;

    return (
        <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
            <DialogContent sx={{ position: 'relative', p: 0 }}>
                <IconButton
                    onClick={onClose}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: 'white',
                        bgcolor: 'rgba(0, 0, 0, 0.5)',
                        '&:hover': {
                            bgcolor: 'rgba(0, 0, 0, 0.7)',
                        },
                    }}
                >
                    <CloseIcon />
                </IconButton>
                <Box sx={{ height: '80vh', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                    <img
                        src={`${BACKEND_URL}/uploads/${images[activeStep].file_path}`} // Исправлено здесь
                        alt={images[activeStep].file_name}
                        style={{
                            width: '100%',
                            height: '100%',
                            objectFit: 'contain',
                        }}
                    />
                    <MobileStepper
                        steps={maxSteps}
                        position="static"
                        activeStep={activeStep}
                        sx={{
                            bgcolor: 'background.default',
                            position: 'absolute',
                            bottom: 0,
                            width: '100%',
                        }}
                        nextButton={
                            <Button size="small" onClick={handleNext}>
                                Следующее
                                <KeyboardArrowRight />
                            </Button>
                        }
                        backButton={
                            <Button size="small" onClick={handleBack}>
                                <KeyboardArrowLeft />
                                Предыдущее
                            </Button>
                        }
                    />
                </Box>
            </DialogContent>
        </Dialog>
    );
};

const RoomList = () => {
    const [rooms, setRooms] = useState([]);
    const [filters, setFilters] = useState({
        capacity: "",
        min_price: "",
        max_price: "",
        city: "",
        country: "",
        start_date: "",
        end_date: ""
    });
    const [viewMode, setViewMode] = useState('list'); // 'list' или 'map'
    const [roomsWithCoordinates, setRoomsWithCoordinates] = useState([]);
    const geocodeRooms = async (rooms) => {
        if (!window.google) return rooms;

        const geocoder = new window.google.maps.Geocoder();
        const geocodeAddress = async (room) => {
            const address = `${room.address_street}, ${room.address_city}, ${room.address_country}`;
            try {
                const result = await new Promise((resolve, reject) => {
                    geocoder.geocode({ address }, (results, status) => {
                        if (status === 'OK') {
                            resolve(results[0].geometry.location);
                        } else {
                            reject(status);
                        }
                    });
                });

                return {
                    ...room,
                    latitude: result.lat(),
                    longitude: result.lng()
                };
            } catch (error) {
                console.error(`Error geocoding address: ${address}`, error);
                return room;
            }
        };

        const roomsWithCoords = await Promise.all(rooms.map(geocodeAddress));
        return roomsWithCoords;
    };
    const [selectedRoom, setSelectedRoom] = useState(null);
    const [galleryOpen, setGalleryOpen] = useState(false);
    const [bookingDialogOpen, setBookingDialogOpen] = useState(false);

    const fetchRooms = useCallback(async () => {
        try {
            const params = new URLSearchParams();
            if (filters.capacity) params.append('capacity', filters.capacity);
            if (filters.min_price) params.append('min_price', filters.min_price);
            if (filters.max_price) params.append('max_price', filters.max_price);
            if (filters.city) params.append('city', filters.city);
            if (filters.country) params.append('country', filters.country);
            if (filters.start_date && filters.end_date) {
                params.append('start_date', filters.start_date);
                params.append('end_date', filters.end_date);
            }
    
            const response = await axios.get(`/rooms?${params.toString()}`);
            const roomsData = response.data || [];
    
            // Получаем изображения и геокодируем адреса для каждой комнаты
            const roomsWithImagesAndCoords = await Promise.all(
                roomsData.map(async (room) => {
                    // Получаем изображения
                    const imagesResponse = await axios.get(`/rooms/${room.id}/images`);
                    const images = imagesResponse.data || [];
    
                    // Геокодируем адрес
                    const address = `${room.address_street}, ${room.address_city}, ${room.address_country}`;
                    try {
                        const geocoder = new window.google.maps.Geocoder();
                        const result = await new Promise((resolve, reject) => {
                            geocoder.geocode({ address }, (results, status) => {
                                if (status === 'OK') {
                                    resolve(results[0].geometry.location);
                                } else {
                                    reject(status);
                                }
                            });
                        });
    
                        return {
                            ...room,
                            images,
                            latitude: result.lat(),
                            longitude: result.lng()
                        };
                    } catch (error) {
                        console.error(`Ошибка геокодирования для ${address}:`, error);
                        return {
                            ...room,
                            images
                        };
                    }
                })
            );
    
            setRooms(roomsWithImagesAndCoords);
            setRoomsWithCoordinates(roomsWithImagesAndCoords);
        } catch (error) {
            console.error("Ошибка при получении списка комнат:", error);
        }
    }, [filters]);
    
    const handleDateChange = (field, value) => {
        setFilters(prev => {
            const newFilters = { ...prev, [field]: value };

            if (field === 'end_date' && newFilters.start_date && value < newFilters.start_date) {
                return prev;
            }

            if (field === 'start_date' && newFilters.end_date && value > newFilters.end_date) {
                return prev;
            }

            return newFilters;
        });
    };

    const handleBooking = (room) => {
        if (!filters.start_date || !filters.end_date) {
            alert('Пожалуйста, выберите даты заезда и выезда');
            return;
        }
        setSelectedRoom(room);
        setBookingDialogOpen(true);
    };


    useEffect(() => {
        fetchRooms();
    }, [fetchRooms]);
    const viewToggle = (
        <Paper sx={{ p: 1, mb: 2 }}>
            <ToggleButtonGroup
                value={viewMode}
                exclusive
                onChange={(e, newMode) => newMode && setViewMode(newMode)}
                aria-label="view mode"
            >
                <ToggleButton value="list" aria-label="list view">
                    <ViewListIcon /> Список
                </ToggleButton>
                <ToggleButton value="map" aria-label="map view">
                    <MapIcon /> Карта
                </ToggleButton>
            </ToggleButtonGroup>
        </Paper>
    );
    const AccommodationInfo = ({ room }) => {
        const getAccommodationInfo = () => {
            switch (room.accommodation_type) {
                case 'bed':
                    return {
                        title: 'Койко-место',
                        details: `Доступно ${room.available_beds} из ${room.total_beds} мест`,
                        icon: <SingleBedIcon />,
                        shared: true
                    };
                case 'room':
                    return {
                        title: room.is_shared ? 'Общая комната' : 'Отдельная комната',
                        details: `Вместимость: ${room.capacity} чел.`,
                        icon: <HotelIcon />,
                        shared: room.is_shared
                    };
                case 'apartment':
                    return {
                        title: 'Квартира',
                        details: `${room.capacity} комнат`,
                        icon: <ApartmentIcon />,
                        shared: false
                    };
                default:
                    return {
                        title: 'Помещение',
                        details: `Вместимость: ${room.capacity} чел.`,
                        icon: <HomeIcon />,
                        shared: false
                    };
            }
        };

        const info = getAccommodationInfo();

        return (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                {info.icon}
                <Box>
                    <Typography variant="subtitle2" sx={{ fontWeight: 'medium' }}>
                        {info.title}
                        {info.shared && (
                            <Chip
                                size="small"
                                label="Общее помещение"
                                color="secondary"
                                sx={{ ml: 1 }}
                            />
                        )}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {info.details}
                    </Typography>
                </Box>
            </Box>
        );
    };


    const today = new Date().toISOString().split('T')[0];

    return (
        <div>

            <Grid container spacing={2} sx={{ marginBottom: 4 }}>
                <Grid item xs={12} sm={6}>
                    <TextField
                        label="Дата заезда"
                        type="date"
                        fullWidth
                        InputLabelProps={{ shrink: true }}
                        value={filters.start_date}
                        onChange={(e) => handleDateChange('start_date', e.target.value)}
                        inputProps={{ min: today }}
                    />
                </Grid>
                <Grid item xs={12} sm={6}>
                    <TextField
                        label="Дата выезда"
                        type="date"
                        fullWidth
                        InputLabelProps={{ shrink: true }}
                        value={filters.end_date}
                        onChange={(e) => handleDateChange('end_date', e.target.value)}
                        inputProps={{ min: filters.start_date || today }}
                    />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <TextField
                        label="Вместимость"
                        type="number"
                        fullWidth
                        value={filters.capacity}
                        onChange={(e) => setFilters({ ...filters, capacity: e.target.value })}
                    />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <TextField
                        label="Мин. цена"
                        type="number"
                        fullWidth
                        value={filters.min_price}
                        onChange={(e) => setFilters({ ...filters, min_price: e.target.value })}
                    />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <TextField
                        label="Макс. цена"
                        type="number"
                        fullWidth
                        value={filters.max_price}
                        onChange={(e) => setFilters({ ...filters, max_price: e.target.value })}
                    />
                </Grid>
                <Grid item xs={12} sm={6}>
                    <TextField
                        label="Город"
                        fullWidth
                        value={filters.city}
                        onChange={(e) => setFilters({ ...filters, city: e.target.value })}
                    />
                </Grid>
                <Grid item xs={12} sm={6}>
                    <TextField
                        label="Страна"
                        fullWidth
                        value={filters.country}
                        onChange={(e) => setFilters({ ...filters, country: e.target.value })}
                    />
                </Grid>
                <Grid item xs={12}>
                    <Button variant="contained" color="primary" onClick={fetchRooms}>
                        Фильтровать
                    </Button>
                </Grid>
            </Grid>
            {viewToggle}

            {viewMode === 'map' ? (
                <LoadScript 
                    googleMapsApiKey={process.env.REACT_APP_GOOGLE_MAPS_API_KEY}
                    libraries={["places", "geometry"]}
                >
                    <MapView
                        rooms={roomsWithCoordinates}
                        onRoomSelect={(room) => {
                            setSelectedRoom(room);
                            setBookingDialogOpen(true);
                        }}
                        onOpenGallery={(room) => {
                            setSelectedRoom(room);
                            setGalleryOpen(true);
                        }}
                    />
                </LoadScript>
            ) : (
                <Grid container spacing={2}>
                    {rooms.map((room) => (
                        <Grid item xs={12} md={6} lg={4} key={room.id}>
                            <Card sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                height: '100%',
                                '& .MuiCardContent-root': {
                                    padding: '12px',
                                },
                                '& .MuiTypography-root': {
                                    lineHeight: '1.3',
                                }
                            }}>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', p: 1.5 }}>
                                    <Box sx={{ flex: 1, pr: 1.5 }}>
                                        <AccommodationInfo room={room} />
                                        <Typography variant="h6" sx={{
                                            mb: 0.5,
                                            fontSize: '1.1rem'
                                        }}>
                                            {room.name}
                                        </Typography>
                                        {room.accommodation_type === 'bed' ? (
                                            <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                                Цена за койко-место: {room.price_per_night} руб./ночь
                                            </Typography>
                                        ) : (
                                            <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                                                Цена за {room.accommodation_type === 'apartment' ? 'квартиру' : 'комнату'}: {room.price_per_night} руб./ночь
                                            </Typography>
                                        )}
                                    </Box>

                                    {/* Правый верхний угол: эскиз */}
                                    <Box sx={{
                                        width: '100px',
                                        height: '100px',
                                        flexShrink: 0,
                                        p: room.images?.length ? 0 : 1
                                    }}>
                                        {room.images && room.images.length > 0 ? (
                                            <CardMedia
                                                component="img"
                                                sx={{
                                                    width: '100%',
                                                    height: '100%',
                                                    objectFit: 'cover',
                                                    borderRadius: '4px',
                                                    '&:hover': {
                                                        opacity: 0.8,
                                                        transition: 'opacity 0.2s ease-in-out',
                                                    },
                                                }}
                                                image={`${BACKEND_URL}/uploads/${room.images[0].file_path}`}
                                                alt={room.name}
                                                onClick={() => {
                                                    setSelectedRoom(room);
                                                    setGalleryOpen(true);
                                                }}
                                            />
                                        ) : (
                                            <Box
                                                sx={{
                                                    width: '100%',
                                                    height: '100%',
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    bgcolor: 'grey.100',
                                                    borderRadius: '4px',
                                                    fontSize: '0.8rem'
                                                }}
                                            >
                                                <Typography variant="body2" color="text.secondary">
                                                    Нет фото
                                                </Typography>
                                            </Box>
                                        )}
                                    </Box>
                                </Box>

                                <Divider />

                                {/* Нижняя часть карточки с адресом и кнопками */}
                                <CardContent sx={{
                                    pt: 1,
                                    pb: '8px !important',
                                    display: 'flex',
                                    justifyContent: 'space-between',
                                    alignItems: 'center'
                                }}>
                                    <Typography variant="body2" color="text.secondary" sx={{
                                        fontSize: '0.85rem',
                                        flex: 1,
                                        mr: 1
                                    }}>
                                        {room.address_street}
                                        {room.address_city && `, ${room.address_city}`}
                                        {room.address_state && `, ${room.address_state}`}
                                        {room.address_country && `, ${room.address_country}`}
                                        {room.address_postal_code && ` (${room.address_postal_code})`}
                                    </Typography>
                                    <Box sx={{ display: 'flex', gap: 0.5 }}>
                                        {room.images && room.images.length > 1 && (
                                            <Button
                                                size="small"
                                                sx={{
                                                    minWidth: 'auto',
                                                    padding: '4px 8px',
                                                    fontSize: '0.8rem'
                                                }}
                                                onClick={() => {
                                                    setSelectedRoom(room);
                                                    setGalleryOpen(true);
                                                }}
                                            >
                                                Все ({room.images.length})
                                            </Button>
                                        )}
                                        <Button
                                            variant="contained"
                                            color="primary"
                                            size="small"
                                            sx={{
                                                minWidth: 'auto',
                                                padding: '4px 8px',
                                                fontSize: '0.8rem'
                                            }}
                                            onClick={() => handleBooking(room)}
                                            disabled={!filters.start_date || !filters.end_date}
                                        >
                                            Забронировать
                                        </Button>
                                    </Box>
                                </CardContent>
                            </Card>
                        </Grid>

                    ))}
                </Grid>
            )}
            {selectedRoom && (
                <>
                    <ImageGallery
                        images={selectedRoom.images || []}
                        open={galleryOpen}
                        onClose={() => {
                            setGalleryOpen(false);
                            setSelectedRoom(null);
                        }}
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
        </div>
    );
};

export default RoomList;
import React, { useState, useCallback } from 'react';
import { GoogleMap, Marker, InfoWindow } from '@react-google-maps/api';
import {
    Card,
    CardContent,
    Typography,
    Box,
    Button,
    Chip,
    CardMedia,
} from '@mui/material';
import {
    SingleBed as SingleBedIcon,
    Hotel as HotelIcon,
    Apartment as ApartmentIcon,
    PhotoLibrary as PhotoLibraryIcon,
    MyLocation as MyLocationIcon,
} from '@mui/icons-material';

const BACKEND_URL = 'http://localhost:3000';

const mapContainerStyle = {
    width: '100%',
    height: '700px'
};

const defaultCenter = {
    lat: 45.2671, // Нови-Сад
    lng: 19.8335
};

const MapView = ({ rooms, onRoomSelect, onOpenGallery }) => {
    const [map, setMap] = useState(null);
    const [selectedRoom, setSelectedRoom] = useState(null);

    const getMapOptions = () => ({
        scrollwheel: true,
        mapTypeControl: true,
        mapTypeControlOptions: {
            style: window.google?.maps.MapTypeControlStyle.DROPDOWN_MENU,
            mapTypeIds: ["roadmap", "satellite", "hybrid"]
        },
        styles: [
            {
                featureType: "poi",
                elementType: "labels",
                stylers: [{ visibility: "off" }]
            }
        ],
        fullscreenControl: true,
        streetViewControl: false,
        zoomControl: true,
    });

    const onMapLoad = useCallback((map) => {
        setMap(map);
        if (rooms.length > 0) {
            const bounds = new window.google.maps.LatLngBounds();
            rooms.forEach(room => {
                if (room.latitude && room.longitude) {
                    bounds.extend({
                        lat: parseFloat(room.latitude),
                        lng: parseFloat(room.longitude)
                    });
                }
            });
            map.fitBounds(bounds);
        }
    }, [rooms]);

    const handleMyLocation = () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (position) => {
                    const pos = {
                        lat: position.coords.latitude,
                        lng: position.coords.longitude
                    };
                    map?.panTo(pos);
                    map?.setZoom(15);
                },
                (error) => {
                    console.error("Error getting location:", error);
                    alert("Не удалось получить местоположение");
                }
            );
        } else {
            alert("Геолокация не поддерживается вашим браузером");
        }
    };

    const InfoWindowContent = ({ room }) => {
        const hasImages = room.images && room.images.length > 0;
        const mainImage = hasImages ? room.images.find(img => img.is_main) || room.images[0] : null;

        return (
            <Card sx={{ width: 300, border: 'none', boxShadow: 'none' }}>
                {mainImage && (
                    <Box
                        sx={{
                            position: 'relative',
                            cursor: 'pointer',
                            '&:hover': {
                                '& .overlay': {
                                    opacity: 1
                                }
                            }
                        }}
                        onClick={(e) => {
                            e.stopPropagation();
                            onOpenGallery(room);
                        }}
                    >
                        <CardMedia
                            component="img"
                            height="160"
                            image={`${BACKEND_URL}/uploads/${mainImage.file_path}`}
                            alt={room.name}
                            sx={{
                                borderRadius: '4px 4px 0 0',
                                objectFit: 'cover'
                            }}
                        />
                        {room.images.length > 1 && (
                            <Box
                                className="overlay"
                                sx={{
                                    position: 'absolute',
                                    bottom: 0,
                                    right: 0,
                                    bgcolor: 'rgba(0, 0, 0, 0.6)',
                                    color: 'white',
                                    px: 1,
                                    py: 0.5,
                                    borderRadius: '4px 0 0 0',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 0.5,
                                    opacity: 0,
                                    transition: 'opacity 0.2s',
                                }}
                            >
                                <PhotoLibraryIcon fontSize="small" />
                                <Typography variant="caption">
                                    +{room.images.length - 1}
                                </Typography>
                            </Box>
                        )}
                    </Box>
                )}
                <CardContent sx={{ p: 1.5 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                        {room.accommodation_type === 'bed' ? (
                            <SingleBedIcon sx={{ mr: 1 }} />
                        ) : room.accommodation_type === 'apartment' ? (
                            <ApartmentIcon sx={{ mr: 1 }} />
                        ) : (
                            <HotelIcon sx={{ mr: 1 }} />
                        )}
                        <Typography variant="subtitle1" component="div">
                            {room.name}
                        </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                        {room.address_street}, {room.address_city}
                    </Typography>
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 1.5 }}>
                        <Chip
                            size="small"
                            label={`${room.price_per_night} ₽/ночь`}
                            color="primary"
                        />
                        {room.accommodation_type === 'bed' && (
                            <Chip
                                size="small"
                                label={`${room.available_beds}/${room.total_beds} мест`}
                                color="secondary"
                            />
                        )}
                    </Box>
                    <Button
                        variant="contained"
                        size="small"
                        fullWidth
                        onClick={(e) => {
                            e.stopPropagation();
                            onRoomSelect(room);
                        }}
                    >
                        Забронировать
                    </Button>
                </CardContent>
            </Card>
        );
    };

    return (
        <Box sx={{ position: 'relative' }}>
            <GoogleMap
                mapContainerStyle={mapContainerStyle}
                center={defaultCenter}
                zoom={13}
                onLoad={onMapLoad}
                options={getMapOptions()}
            >
                {rooms.map((room) => (
                    room.latitude && room.longitude ? (
                        <Marker
                            key={room.id}
                            position={{
                                lat: parseFloat(room.latitude),
                                lng: parseFloat(room.longitude)
                            }}
                            onClick={() => setSelectedRoom(room)}
                            icon={{
                                path: window.google?.maps.SymbolPath.CIRCLE,
                                fillColor: room.accommodation_type === 'bed'
                                    ? '#1976d2'
                                    : room.accommodation_type === 'apartment'
                                        ? '#dc004e'
                                        : '#4caf50',
                                fillOpacity: 1,
                                strokeWeight: 1,
                                strokeColor: '#ffffff',
                                scale: 10,
                            }}
                        />
                    ) : null
                ))}

                {selectedRoom && (
                    <InfoWindow
                        position={{
                            lat: parseFloat(selectedRoom.latitude),
                            lng: parseFloat(selectedRoom.longitude)
                        }}
                        onCloseClick={() => setSelectedRoom(null)}
                    >
                        <InfoWindowContent room={selectedRoom} />
                    </InfoWindow>
                )}
            </GoogleMap>
            <Button
                variant="contained"
                startIcon={<MyLocationIcon />}
                onClick={handleMyLocation}
                sx={{
                    position: 'absolute',
                    top: '10px',
                    right: '60px',
                    backgroundColor: 'white',
                    color: 'black',
                    '&:hover': {
                        backgroundColor: '#f5f5f5',
                    },
                    boxShadow: '0 2px 6px rgba(0,0,0,.3)',
                }}
            >
                Моё местоположение
            </Button>
        </Box>
    );
};

export default MapView;
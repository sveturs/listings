import React, { useState, useEffect, useCallback } from 'react';
import { GoogleMap, LoadScript, Marker, InfoWindow } from '@react-google-maps/api';
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
} from '@mui/icons-material';

const GOOGLE_MAPS_API_KEY = process.env.REACT_APP_GOOGLE_MAPS_API_KEY;
const BACKEND_URL = 'http://localhost:3000';

const MapView = ({ rooms, onRoomSelect, onOpenGallery }) => {
    const [selectedRoom, setSelectedRoom] = useState(null);

    const mapContainerStyle = {
        width: '100%',
        height: '700px'
    };

    const defaultCenter = {
        lat: 55.7558,
        lng: 37.6173
    };

    const onMapLoad = useCallback((map) => {
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
        <GoogleMap
            mapContainerStyle={mapContainerStyle}
            center={defaultCenter}
            zoom={10}
            onLoad={onMapLoad}
            options={{
                styles: [
                    {
                        featureType: "poi",
                        elementType: "labels",
                        stylers: [{ visibility: "off" }]
                    }
                ],
                fullscreenControl: false,
                streetViewControl: false,
            }}
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
                            path: window.google.maps.SymbolPath.CIRCLE,
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
    );
};

export default MapView;
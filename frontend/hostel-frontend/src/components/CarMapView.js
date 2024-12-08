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
    LocalGasStation as FuelIcon,
    Speed as TransmissionIcon,
    PhotoLibrary as PhotoLibraryIcon,
} from '@mui/icons-material';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const mapContainerStyle = {
    width: '100%',
    height: '700px'
};

const defaultCenter = {
    lat: 45.2671, // Нови-Сад
    lng: 19.8335
};

const CarMapView = ({ cars, onCarSelect, onOpenGallery }) => {
    const [map, setMap] = useState(null);
    const [selectedCar, setSelectedCar] = useState(null);

    const getMapOptions = useCallback(() => ({
        scrollwheel: true,
        mapTypeControl: true,
        mapTypeControlOptions: {
            style: 'DEFAULT',
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
    }), []);

    const onMapLoad = useCallback((map) => {
        setMap(map);
        if (cars?.length > 0) {
            const bounds = new window.google.maps.LatLngBounds();
            let hasValidCoords = false;

            cars.forEach(car => {
                if (car.latitude && car.longitude) {
                    bounds.extend({ lat: car.latitude, lng: car.longitude });
                    hasValidCoords = true;
                }
            });

            if (hasValidCoords) {
                map.fitBounds(bounds);
            }
        }
    }, [cars]);

    const InfoWindowContent = ({ car }) => {
        const hasImages = car.images && car.images.length > 0;
        const mainImage = hasImages ? car.images[0] : null;

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
                            onOpenGallery(car);
                        }}
                    >
                        <CardMedia
                            component="img"
                            height="160"
                            image={`${BACKEND_URL}/uploads/${mainImage.file_path}`}
                            alt={`${car.make} ${car.model}`}
                            sx={{
                                borderRadius: '4px 4px 0 0',
                                objectFit: 'cover'
                            }}
                        />
                        {car.images.length > 1 && (
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
                                    gap: 0.5
                                }}
                            >
                                <PhotoLibraryIcon fontSize="small" />
                                <Typography variant="caption">
                                    +{car.images.length - 1}
                                </Typography>
                            </Box>
                        )}
                    </Box>
                )}
                <CardContent sx={{ p: 1.5 }}>
                    <Typography variant="h6" gutterBottom>
                        {car.make} {car.model}
                        <Typography variant="body2" color="text.secondary" component="span" sx={{ ml: 1 }}>
                            {car.year}
                        </Typography>
                    </Typography>

                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 1.5 }}>
                        <Chip
                            icon={<FuelIcon />}
                            label={car.fuel_type === 'petrol' ? 'Бензин' :
                                  car.fuel_type === 'diesel' ? 'Дизель' :
                                  car.fuel_type === 'electric' ? 'Электро' : 'Гибрид'}
                            size="small"
                        />
                        <Chip
                            icon={<TransmissionIcon />}
                            label={car.transmission === 'automatic' ? 'Автомат' : 'Механика'}
                            size="small"
                        />
                        <Chip
                            size="small"
                            label={`${car.price_per_day} ₽/день`}
                            color="primary"
                        />
                    </Box>

                    <Button
                        variant="contained"
                        size="small"
                        fullWidth
                        onClick={(e) => {
                            e.stopPropagation();
                            onCarSelect(car);
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
                {cars?.map((car) => {
                    if (car.latitude && car.longitude) {
                        return (
                            <Marker
                                key={car.id}
                                position={{ lat: car.latitude, lng: car.longitude }}
                                onClick={() => setSelectedCar(car)}
                                icon={{
                                    path: "M10 2C6.68 2 4 4.68 4 8c0 2.04 1.01 3.84 2.56 4.94L10 20l3.44-7.06C14.99 11.84 16 10.04 16 8c0-3.32-2.68-6-6-6zm0 8.5c-1.38 0-2.5-1.12-2.5-2.5S8.62 5.5 10 5.5s2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z",
                                    fillColor: '#1976d2',
                                    fillOpacity: 1,
                                    strokeWeight: 1,
                                    strokeColor: '#ffffff',
                                    scale: 2,
                                    anchor: new window.google.maps.Point(10, 20),
                                }}
                            />
                        );
                    }
                    return null;
                })}

                {selectedCar && (
                    <InfoWindow
                        position={{ lat: selectedCar.latitude, lng: selectedCar.longitude }}
                        onCloseClick={() => setSelectedCar(null)}
                    >
                        <InfoWindowContent car={selectedCar} />
                    </InfoWindow>
                )}
            </GoogleMap>
        </Box>
    );
};

export default CarMapView;
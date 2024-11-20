import React, { useState, useCallback } from 'react';
import { GoogleMap, LoadScript, Marker } from '@react-google-maps/api';
import {
    Box,
    TextField,
    Paper,
    Typography,
    InputAdornment,
    IconButton,
} from '@mui/material';
import { Search as SearchIcon, MyLocation as MyLocationIcon } from '@mui/icons-material';

const libraries = ["places", "geometry"];

const LocationPicker = ({ onLocationSelect }) => {
    const [map, setMap] = useState(null);
    const [marker, setMarker] = useState(null);
    const [address, setAddress] = useState('');

    const mapContainerStyle = {
        width: '100%',
        height: '400px'
    };

    const defaultCenter = {
        lat: 45.2671,  // Нови-Сад
        lng: 19.8335
    };

    const mapOptions = {
        scrollwheel: true,
        mapTypeControl: true,
        mapTypeControlOptions: {
            style: window.google.maps.MapTypeControlStyle.DROPDOWN_MENU,
            mapTypeIds: ["roadmap", "satellite", "hybrid"]
        },
        streetViewControl: false,
        fullscreenControl: true,
    };

    const onMapLoad = useCallback((map) => {
        setMap(map);
        // Инициализируем поисковую строку
        const searchBoxElement = new window.google.maps.places.SearchBox(
            document.getElementById('location-search')
        );

        // Слушаем изменения в поисковой строке
        searchBoxElement.addListener('places_changed', () => {
            const places = searchBoxElement.getPlaces();
            if (places.length === 0) return;

            const place = places[0];
            if (!place.geometry) return;

            // Центрируем карту на найденном месте
            map.setCenter(place.geometry.location);
            map.setZoom(17);

            // Устанавливаем маркер
            setMarker({
                lat: place.geometry.location.lat(),
                lng: place.geometry.location.lng()
            });

            // Обновляем адрес
            setAddress(place.formatted_address);

            // Вызываем callback с данными
            onLocationSelect({
                latitude: place.geometry.location.lat(),
                longitude: place.geometry.location.lng(),
                formatted_address: place.formatted_address,
                address_components: place.address_components
            });
        });
    }, [onLocationSelect]);

    const handleCurrentLocation = () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (position) => {
                    const lat = position.coords.latitude;
                    const lng = position.coords.longitude;

                    // Центрируем карту и устанавливаем маркер
                    if (map) {
                        map.setCenter({ lat, lng });
                        map.setZoom(17);
                        setMarker({ lat, lng });
                    }

                    // Получаем адрес по координатам
                    const geocoder = new window.google.maps.Geocoder();
                    geocoder.geocode(
                        { location: { lat, lng } },
                        (results, status) => {
                            if (status === 'OK' && results[0]) {
                                const place = results[0];
                                // Обновляем адрес в поле ввода
                                setAddress(place.formatted_address);
                                // Вызываем callback с данными
                                onLocationSelect({
                                    latitude: lat,
                                    longitude: lng,
                                    formatted_address: place.formatted_address,
                                    address_components: place.address_components
                                });
                            }
                        }
                    );
                },
                (error) => {
                    console.error("Error getting current location:", error);
                    alert("Не удалось получить текущее местоположение");
                },
                {
                    enableHighAccuracy: true,
                    timeout: 5000,
                    maximumAge: 0
                }
            );
        } else {
            alert("Геолокация не поддерживается вашим браузером");
        }
    };

    const handleMapClick = useCallback((e) => {
        const lat = e.latLng.lat();
        const lng = e.latLng.lng();
        
        setMarker({ lat, lng });

        // Получаем адрес по координатам
        const geocoder = new window.google.maps.Geocoder();
        geocoder.geocode(
            { location: { lat, lng } }, 
            (results, status) => {
                if (status === 'OK' && results[0]) {
                    const place = results[0];
                    setAddress(place.formatted_address);
                    onLocationSelect({
                        latitude: lat,
                        longitude: lng,
                        formatted_address: place.formatted_address,
                        address_components: place.address_components
                    });
                }
            }
        );
    }, [onLocationSelect]);

    return (
        <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
                Выберите местоположение объекта
            </Typography>
            <Box sx={{ mb: 2 }}>
                <TextField
                    id="location-search"
                    fullWidth
                    placeholder="Поиск по адресу..."
                    value={address}
                    onChange={(e) => setAddress(e.target.value)}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon />
                            </InputAdornment>
                        ),
                        endAdornment: (
                            <InputAdornment position="end">
                                <IconButton 
                                    onClick={handleCurrentLocation}
                                    title="Мое местоположение"
                                >
                                    <MyLocationIcon />
                                </IconButton>
                            </InputAdornment>
                        )
                    }}
                />
            </Box>
            <LoadScript
                googleMapsApiKey={process.env.REACT_APP_GOOGLE_MAPS_API_KEY}
                libraries={libraries}
            >
                <GoogleMap
                    mapContainerStyle={mapContainerStyle}
                    center={defaultCenter}
                    zoom={10}
                    onLoad={onMapLoad}
                    onClick={handleMapClick}
                    options={mapOptions}
                >
                    {marker && (
                        <Marker
                            position={marker}
                            draggable={true}
                            onDragEnd={(e) => handleMapClick(e)}
                        />
                    )}
                </GoogleMap>
            </LoadScript>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Кликните по карте или введите адрес для выбора местоположения
            </Typography>
        </Paper>
    );
};

export default LocationPicker;
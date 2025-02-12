import React, { useState, useCallback } from 'react';
import { GoogleMap, Marker } from '@react-google-maps/api';
import {
    Box,
    TextField,
    Paper,
    Typography,
    InputAdornment,
    IconButton,
    useTheme,
    useMediaQuery
} from '@mui/material';
import { Search as SearchIcon, MyLocation as MyLocationIcon } from '@mui/icons-material';

const LocationPicker = ({ onLocationSelect }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    
    const [map, setMap] = useState(null);
    const [marker, setMarker] = useState(null);
    const [address, setAddress] = useState('');
    const [setSearchBox] = useState(null);
    const [center, setCenter] = useState({
        lat: 45.2671,
        lng: 19.8335
    });
    const mapContainerStyle = {
        width: '100%',
        height: '400px'
    };


    const mapOptions = {
        scrollwheel: true,
        mapTypeControl: true,
        mapTypeControlOptions: {
            style: 'DEFAULT',
            mapTypeIds: ["roadmap", "satellite", "hybrid"]
        },
        streetViewControl: false,
        fullscreenControl: false,
    };

    const handleLocationSelect = (location) => {
        const getAddressComponent = (type) => {
            return location.address_components?.find(
                component => component.types.includes(type)
            )?.long_name || '';
        };

        const streetNumber = getAddressComponent('street_number');
        const route = getAddressComponent('route');

        const fullStreetAddress = route
            ? (streetNumber ? `${route}, ${streetNumber}` : route)
            : '';

        onLocationSelect({
            latitude: location.latitude,
            longitude: location.longitude,
            formatted_address: location.formatted_address,
            address_components: {
                street: fullStreetAddress || getAddressComponent('sublocality') || '',
                city: getAddressComponent('locality'),
                state: getAddressComponent('administrative_area_level_1'),
                country: getAddressComponent('country'),
                postal_code: getAddressComponent('postal_code')
            }
        });
    };

    const onMapLoad = useCallback((map) => {
        setMap(map);
        const searchInput = document.getElementById('location-search');
        if (searchInput && window.google) {
            const searchBoxInstance = new window.google.maps.places.SearchBox(searchInput);
            setSearchBox(searchBoxInstance);

            searchBoxInstance.addListener('places_changed', () => {
                const places = searchBoxInstance.getPlaces();
                if (places.length === 0) return;

                const place = places[0];
                if (!place.geometry) return;

                const newCenter = {
                    lat: place.geometry.location.lat(),
                    lng: place.geometry.location.lng()
                };

                setCenter(newCenter);
                map.setCenter(newCenter);
                map.setZoom(17);

                setMarker(newCenter);
                setAddress(place.formatted_address);

                handleLocationSelect({
                    latitude: newCenter.lat,
                    longitude: newCenter.lng,
                    formatted_address: place.formatted_address,
                    address_components: place.address_components
                });
            });
        }
    }, [onLocationSelect]);

    const handleMapClick = useCallback((e) => {
        const lat = e.latLng.lat();
        const lng = e.latLng.lng();

        setMarker({ lat, lng });

        if (window.google) {
            const geocoder = new window.google.maps.Geocoder();
            geocoder.geocode(
                { location: { lat, lng } },
                (results, status) => {
                    if (status === 'OK' && results[0]) {
                        const place = results[0];
                        setAddress(place.formatted_address);
                        const location = {
                            latitude: lat,
                            longitude: lng,
                            formatted_address: place.formatted_address,
                            address_components: place.address_components
                        };
                        handleLocationSelect(location);
                    }
                }
            );
        }
    }, []);

    const handleCurrentLocation = () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (position) => {
                    const newCenter = {
                        lat: position.coords.latitude,
                        lng: position.coords.longitude
                    };

                    setCenter(newCenter);
                    setMarker(newCenter);
                    
                    if (map) {
                        map.setCenter(newCenter);
                        map.setZoom(17);
                    }

                    if (window.google) {
                        const geocoder = new window.google.maps.Geocoder();
                        geocoder.geocode(
                            { location: newCenter },
                            (results, status) => {
                                if (status === 'OK' && results[0]) {
                                    const place = results[0];
                                    setAddress(place.formatted_address);
                                    handleLocationSelect({
                                        latitude: newCenter.lat,
                                        longitude: newCenter.lng,
                                        formatted_address: place.formatted_address,
                                        address_components: place.address_components
                                    });
                                }
                            }
                        );
                    }
                },
                (error) => {
                    console.error("Error getting current location:", error);
                    alert("Не удалось получить текущее местоположение");
                }
            );
        } else {
            alert("Геолокация не поддерживается вашим браузером");
        }
    };

    return (
        <Paper sx={{ p: isMobile ? 0 : 2, bgcolor: isMobile ? 'transparent' : 'background.paper', elevation: isMobile ? 0 : 1 }}>
            <Typography variant={isMobile ? "subtitle1" : "h6"} gutterBottom>
                Местоположение объекта
            </Typography>
            <Box sx={{ mb: isMobile ? 1 : 2 }}>
                <TextField
                    id="location-search"
                    fullWidth
                    placeholder="Поиск по адресу..."
                    value={address}
                    onChange={(e) => setAddress(e.target.value)}
                    size={isMobile ? "small" : "medium"}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon fontSize={isMobile ? "small" : "medium"} />
                            </InputAdornment>
                        ),
                        endAdornment: (
                            <InputAdornment position="end">
                                <IconButton
                                    onClick={handleCurrentLocation}
                                    title="Мое местоположение"
                                    size={isMobile ? "small" : "medium"}
                                >
                                    <MyLocationIcon fontSize={isMobile ? "small" : "medium"} />
                                </IconButton>
                            </InputAdornment>
                        )
                    }}
                />
            </Box>
            <GoogleMap
                mapContainerStyle={mapContainerStyle}
                center={center}  // Используем центр из state вместо defaultCenter
                zoom={13}
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
            <Typography variant={isMobile ? "caption" : "body2"} color="text.secondary" sx={{ mt: 1 }}>
                Кликните по карте или введите адрес для выбора местоположения
            </Typography>
        </Paper>
    );
};

export default LocationPicker;
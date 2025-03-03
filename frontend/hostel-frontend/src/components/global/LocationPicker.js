import React, { useState, useEffect, useRef } from 'react';
import L from 'leaflet';
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
import 'leaflet/dist/leaflet.css';
import '../maps/leaflet-icons'; // Импортируем файл с фиксом иконок

// Исправляем проблему с маркерами Leaflet в React
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png',
    iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
});

// Эта функция больше не нужна при использовании прямого API Leaflet
// удаляем MapEventHandler, так как он был для react-leaflet

// Компонент для геокодирования (поиск адреса по координатам)
const reverseGeocode = async (lat, lng) => {
    try {
        const response = await fetch(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&addressdetails=1`);
        const data = await response.json();
        
        if (data && data.display_name) {
            return {
                formatted_address: data.display_name,
                address_components: {
                    street: data.address.road || data.address.pedestrian || '',
                    city: data.address.city || data.address.town || data.address.village || '',
                    state: data.address.state || '',
                    country: data.address.country || '',
                    postal_code: data.address.postcode || ''
                }
            };
        }
        return null;
    } catch (error) {
        console.error("Error in reverse geocoding:", error);
        return null;
    }
};

// Компонент для поиска адреса
const searchAddress = async (query) => {
    try {
        const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&limit=1`);
        const data = await response.json();
        
        if (data && data.length > 0) {
            const place = data[0];
            return {
                latitude: parseFloat(place.lat),
                longitude: parseFloat(place.lon),
                formatted_address: place.display_name,
                address_components: {
                    street: '',
                    city: '',
                    state: '',
                    country: '',
                    postal_code: ''
                }
            };
        }
        return null;
    } catch (error) {
        console.error("Error searching address:", error);
        return null;
    }
};

const LocationPicker = ({ onLocationSelect, initialLocation }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const [marker, setMarker] = useState(null);
    const [address, setAddress] = useState('');
    const [center, setCenter] = useState({
        lat: initialLocation?.latitude || 45.2671,
        lng: initialLocation?.longitude || 19.8335
    });

    // Реф для доступа к DOM-элементу карты
    const mapContainerRef = useRef(null);
    // Реф для хранения экземпляра карты Leaflet
    const leafletMapRef = useRef(null);
    // Реф для хранения маркера
    const markerRef = useRef(null);

    // Устанавливаем начальный маркер, если есть initialLocation
    useEffect(() => {
        if (initialLocation && initialLocation.latitude && initialLocation.longitude) {
            setMarker({
                lat: initialLocation.latitude,
                lng: initialLocation.longitude
            });
            
            if (initialLocation.formatted_address) {
                setAddress(initialLocation.formatted_address);
            }
        }
    }, [initialLocation]);
    
    // Инициализация карты Leaflet
    useEffect(() => {
        if (!mapContainerRef.current) return;
        
        // Инициализируем карту только если её еще нет
        if (!leafletMapRef.current) {
            leafletMapRef.current = L.map(mapContainerRef.current).setView([center.lat, center.lng], 13);
            
            // Добавляем слой тайлов OpenStreetMap
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
                maxZoom: 19
            }).addTo(leafletMapRef.current);
            
            // Добавляем обработчик события клика
            leafletMapRef.current.on('click', async (e) => {
                const { lat, lng } = e.latlng;
                await handleMapClick({ latlng: { lat, lng } });
            });
        } else {
            // Если карта уже существует, просто обновляем её центр
            leafletMapRef.current.setView([center.lat, center.lng], 13);
        }
        
        // Очистка при размонтировании компонента
        return () => {
            if (leafletMapRef.current) {
                leafletMapRef.current.remove();
                leafletMapRef.current = null;
            }
        };
    }, []);

    // Обработчик клика по карте
    const handleMapClick = async (e) => {
        const { lat, lng } = e.latlng;
        setMarker({ lat, lng });
        
        // Обновляем маркер на карте
        if (leafletMapRef.current) {
            // Удаляем предыдущий маркер, если он существует
            if (markerRef.current) {
                leafletMapRef.current.removeLayer(markerRef.current);
            }
            
            // Создаем новый маркер
            markerRef.current = L.marker([lat, lng], { draggable: true })
                .addTo(leafletMapRef.current)
                .on('dragend', function(event) {
                    const marker = event.target;
                    const position = marker.getLatLng();
                    handleMarkerDragEnd({ target: marker });
                });
        }

        // Получаем адрес по координатам
        const locationInfo = await reverseGeocode(lat, lng);
        if (locationInfo) {
            setAddress(locationInfo.formatted_address);
            
            // Передаем информацию о местоположении родительскому компоненту
            onLocationSelect({
                latitude: lat,
                longitude: lng,
                formatted_address: locationInfo.formatted_address,
                address_components: locationInfo.address_components
            });
        }
    };

    // Обработчик поиска адреса
    const handleSearch = async () => {
        if (!address) return;
        
        const result = await searchAddress(address);
        if (result) {
            setMarker({
                lat: result.latitude,
                lng: result.longitude
            });
            
            setCenter({
                lat: result.latitude,
                lng: result.longitude
            });
            
            if (leafletMapRef.current) {
                leafletMapRef.current.setView([result.latitude, result.longitude], 17);
                
                // Удаляем предыдущий маркер, если он существует
                if (markerRef.current) {
                    leafletMapRef.current.removeLayer(markerRef.current);
                }
                
                // Создаем новый маркер
                markerRef.current = L.marker([result.latitude, result.longitude], { draggable: true })
                    .addTo(leafletMapRef.current)
                    .on('dragend', function(event) {
                        const marker = event.target;
                        handleMarkerDragEnd({ target: marker });
                    });
            }
            
            setAddress(result.formatted_address);
            
            // Передаем информацию о местоположении родительскому компоненту
            onLocationSelect({
                latitude: result.latitude,
                longitude: result.longitude,
                formatted_address: result.formatted_address,
                address_components: result.address_components
            });
        }
    };

    // Обработчик получения текущего местоположения
    const handleCurrentLocation = () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    const { latitude, longitude } = position.coords;
                    
                    setMarker({
                        lat: latitude,
                        lng: longitude
                    });
                    
                    setCenter({
                        lat: latitude,
                        lng: longitude
                    });
                    
                    if (leafletMapRef.current) {
                        leafletMapRef.current.setView([latitude, longitude], 17);
                        
                        // Удаляем предыдущий маркер, если он существует
                        if (markerRef.current) {
                            leafletMapRef.current.removeLayer(markerRef.current);
                        }
                        
                        // Создаем новый маркер
                        markerRef.current = L.marker([latitude, longitude], { draggable: true })
                            .addTo(leafletMapRef.current)
                            .on('dragend', function(event) {
                                const marker = event.target;
                                handleMarkerDragEnd({ target: marker });
                            });
                    }

                    // Получаем адрес по координатам
                    const locationInfo = await reverseGeocode(latitude, longitude);
                    if (locationInfo) {
                        setAddress(locationInfo.formatted_address);
                        
                        // Передаем информацию о местоположении родительскому компоненту
                        onLocationSelect({
                            latitude,
                            longitude,
                            formatted_address: locationInfo.formatted_address,
                            address_components: locationInfo.address_components
                        });
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

    // Обработчик перетаскивания маркера
    const handleMarkerDragEnd = async (e) => {
        const position = e.target.getLatLng();
        const lat = position.lat;
        const lng = position.lng;
        
        setMarker({
            lat,
            lng
        });

        // Получаем адрес по координатам
        const locationInfo = await reverseGeocode(lat, lng);
        if (locationInfo) {
            setAddress(locationInfo.formatted_address);
            
            // Передаем информацию о местоположении родительскому компоненту
            onLocationSelect({
                latitude: lat,
                longitude: lng,
                formatted_address: locationInfo.formatted_address,
                address_components: locationInfo.address_components
            });
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
                    onKeyPress={(e) => {
                        if (e.key === 'Enter') {
                            handleSearch();
                        }
                    }}
                    size={isMobile ? "small" : "medium"}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <IconButton onClick={handleSearch}>
                                    <SearchIcon fontSize={isMobile ? "small" : "medium"} />
                                </IconButton>
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
            
            <Box sx={{ height: '400px', width: '100%' }}>
                <div 
                    ref={mapContainerRef} 
                    style={{ height: '100%', width: '100%', borderRadius: '4px' }}
                />
            </Box>
            
            <Typography variant={isMobile ? "caption" : "body2"} color="text.secondary" sx={{ mt: 1 }}>
                Кликните по карте или введите адрес для выбора местоположения
            </Typography>
        </Paper>
    );
};

export default LocationPicker;
// frontend/hostel-frontend/src/components/global/LocationPicker.js
import React, { useState, useEffect, useRef, useCallback } from 'react';
import L from 'leaflet';
import {
    Box,
    TextField,
    Paper,
    Typography,
    InputAdornment,
    IconButton,
    useTheme,
    useMediaQuery,
    List,
    ListItem,
    ListItemText,
    Collapse,
    CircularProgress
} from '@mui/material';
import { Search as SearchIcon, MyLocation as MyLocationIcon, KeyboardArrowDown, KeyboardArrowUp } from '@mui/icons-material';
import axios from '../../api/axios';
import 'leaflet/dist/leaflet.css';
import '../maps/leaflet-icons'; // Импортируем файл с фиксом иконок

// Исправляем проблему с маркерами Leaflet в React
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png',
    iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
});

// Компонент для геокодирования (поиск адреса по координатам)
const reverseGeocode = async (lat, lng) => {
    try {
        // Сначала пробуем использовать наш собственный API
        try {
            const response = await axios.get(`/api/v1/geocode/reverse`, {
                params: { lat, lon: lng }
            });
            
            if (response.data && response.data.data) {
                const data = response.data.data;
                return {
                    formatted_address: `${data.city}, ${data.country}`,
                    address_components: {
                        city: data.city || '',
                        country: data.country || '',
                    },
                    latitude: lat,
                    longitude: lng
                };
            }
        } catch (err) {
            console.log('Ошибка нашего API геокодирования, используем OSM:', err);
        }
        
        // В случае ошибки используем Nominatim OSM
        const response = await fetch(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&addressdetails=1`, {
            headers: {
                'User-Agent': 'HostelBookingApp/1.0'
            }
        });
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
                },
                latitude: lat,
                longitude: lng
            };
        }
        return null;
    } catch (error) {
        console.error("Error in reverse geocoding:", error);
        return null;
    }
};

// Компонент для поиска адреса с автодополнением
const searchAddress = async (query) => {
    try {
        // Добавляем "Сербия" к запросу, если её нет, чтобы улучшить геолокацию
        if (!query.toLowerCase().includes('serbia') && !query.toLowerCase().includes('srbija')) {
            query = query + ', Serbia';
        }
        
        // Используем OSM Nominatim с явным указанием страны и полными деталями
        const response = await fetch(
            `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&countrycodes=rs&addressdetails=1&limit=5`, 
            {
                headers: {
                    'User-Agent': 'HostelBookingApp/1.0'
                }
            }
        );
        
        const data = await response.json();
        
        if (data && data.length > 0) {
            return data.map(place => {
                // Извлекаем максимально детальную информацию
                const district = place.address?.suburb || 
                                place.address?.neighbourhood || 
                                place.address?.district || 
                                place.address?.municipality || 
                                '';
                                
                const city = place.address?.city || 
                            place.address?.town || 
                            place.address?.village || 
                            place.address?.municipality || 
                            place.address?.county || 
                            '';
                
                // Формируем локацию для отображения
                let displayLocation = '';
                if (district && city && district !== city) {
                    displayLocation = `${district}, ${city}`;
                } else if (city) {
                    displayLocation = city;
                } else if (district) {
                    displayLocation = district;
                } else {
                    displayLocation = place.display_name;
                }
                
                return {
                    latitude: parseFloat(place.lat),
                    longitude: parseFloat(place.lon),
                    formatted_address: place.display_name,
                    address_components: {
                        city: displayLocation,
                        street: place.address?.road || '',
                        district: district,
                        municipality: place.address?.municipality || '',
                        county: place.address?.county || '',
                        state: place.address?.state || '',
                        country: place.address?.country || 'Serbia',
                        postal_code: place.address?.postcode || ''
                    }
                };
            });
        }
        return [];
    } catch (error) {
        console.error("Error searching address:", error);
        return [];
    }
};

const LocationPicker = ({ onLocationSelect, initialLocation }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const [marker, setMarker] = useState(null);
    const [address, setAddress] = useState('');
    const [addressSuggestions, setAddressSuggestions] = useState([]);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const [searching, setSearching] = useState(false);
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
    // Таймер для задержки поиска
    const searchTimer = useRef(null);

    // Устанавливаем начальный маркер, если есть initialLocation
    useEffect(() => {
        if (initialLocation && initialLocation.latitude && initialLocation.longitude) {
            setMarker({
                lat: initialLocation.latitude,
                lng: initialLocation.longitude
            });
            
            setCenter({
                lat: initialLocation.latitude,
                lng: initialLocation.longitude
            });
            
            if (initialLocation.formatted_address) {
                setAddress(initialLocation.formatted_address);
            }
            
            // Если карта уже инициализирована, обновляем её вид
            if (leafletMapRef.current) {
                leafletMapRef.current.setView([initialLocation.latitude, initialLocation.longitude], 13);
                
                // Удаляем предыдущий маркер, если он существует
                if (markerRef.current) {
                    leafletMapRef.current.removeLayer(markerRef.current);
                }
                
                // Создаем новый маркер
                markerRef.current = L.marker([initialLocation.latitude, initialLocation.longitude], { draggable: true })
                    .addTo(leafletMapRef.current)
                    .on('dragend', function(event) {
                        const marker = event.target;
                        handleMarkerDragEnd({ target: marker });
                    });
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

    // Обработчик изменения адреса с автодополнением
    const handleAddressChange = (e) => {
        const value = e.target.value;
        setAddress(value);
        
        // Очищаем предыдущий таймер
        if (searchTimer.current) {
            clearTimeout(searchTimer.current);
        }
        
        // Если строка поиска пустая, очищаем подсказки
        if (!value.trim()) {
            setAddressSuggestions([]);
            setShowSuggestions(false);
            return;
        }
        
        // Устанавливаем таймер для поиска
        searchTimer.current = setTimeout(async () => {
            setSearching(true);
            const suggestions = await searchAddress(value);
            setAddressSuggestions(suggestions);
            setShowSuggestions(suggestions.length > 0);
            setSearching(false);
        }, 500); // задержка 500мс
    };

    // Обработчик выбора адреса из подсказок
    const handleSelectSuggestion = (suggestion) => {
        setAddress(suggestion.formatted_address);
        setShowSuggestions(false);
        
        const { latitude, longitude } = suggestion;
        
        setMarker({
            lat: latitude,
            lng: longitude
        });
        
        setCenter({
            lat: latitude,
            lng: longitude
        });
        
        if (leafletMapRef.current) {
            leafletMapRef.current.setView([latitude, longitude], 15);
            
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
        
        // Передаем информацию о местоположении родительскому компоненту
        onLocationSelect({
            latitude,
            longitude,
            formatted_address: suggestion.formatted_address,
            address_components: suggestion.address_components
        });
    };

// Обработчик ручного поиска адреса (по нажатию на кнопку или Enter)
const handleSearch = async () => {
    if (!address) return;
    
    setSearching(true);
    try {
        // Пробуем найти точное совпадение
        const suggestions = await searchAddress(address);
        
        // Если есть результаты, используем их
        if (suggestions && suggestions.length > 0) {
            const suggestion = suggestions[0];
            
            // Формируем отображаемое местоположение из результатов геокодирования
            // Приоритет: district + city > city > municipality > district
            let displayLocation = suggestion.address_components.city;
            
            // Обновляем карту, если есть координаты
            if (suggestion.latitude && suggestion.longitude) {
                setMarker({
                    lat: suggestion.latitude,
                    lng: suggestion.longitude
                });
                
                setCenter({
                    lat: suggestion.latitude,
                    lng: suggestion.longitude
                });
                
                if (leafletMapRef.current) {
                    leafletMapRef.current.setView([suggestion.latitude, suggestion.longitude], 13);
                    
                    if (markerRef.current) {
                        leafletMapRef.current.removeLayer(markerRef.current);
                    }
                    
                    markerRef.current = L.marker([suggestion.latitude, suggestion.longitude], { draggable: true })
                        .addTo(leafletMapRef.current)
                        .on('dragend', function(event) {
                            const marker = event.target;
                            handleMarkerDragEnd({ target: marker });
                        });
                }
            }
            
            // Передаем данные родительскому компоненту
            onLocationSelect({
                latitude: suggestion.latitude,
                longitude: suggestion.longitude,
                formatted_address: address,
                address_components: suggestion.address_components
            });
            
            setSearching(false);
            return;
        }
        
        // Если не нашли по точному запросу, пробуем разбить адрес и искать по частям
        const addressParts = address.split(',').map(part => part.trim()).filter(Boolean);
        
        // Пробуем найти по последним частям (они обычно содержат район/город)
        if (addressParts.length > 2) {
            const lastParts = addressParts.slice(-2).join(', ') + ', Serbia';
            const partialResults = await searchAddress(lastParts);
            
            if (partialResults && partialResults.length > 0) {
                const result = partialResults[0];
                
                // Передаем данные родительскому компоненту
                onLocationSelect({
                    latitude: result.latitude,
                    longitude: result.longitude,
                    formatted_address: address,
                    address_components: result.address_components
                });
                
                // Обновляем карту
                if (result.latitude && result.longitude) {
                    setMarker({
                        lat: result.latitude,
                        lng: result.longitude
                    });
                    
                    setCenter({
                        lat: result.latitude,
                        lng: result.longitude
                    });
                    
                    if (leafletMapRef.current) {
                        leafletMapRef.current.setView([result.latitude, result.longitude], 13);
                        
                        if (markerRef.current) {
                            leafletMapRef.current.removeLayer(markerRef.current);
                        }
                        
                        markerRef.current = L.marker([result.latitude, result.longitude], { draggable: true })
                            .addTo(leafletMapRef.current)
                            .on('dragend', function(event) {
                                const marker = event.target;
                                handleMarkerDragEnd({ target: marker });
                            });
                    }
                }
                
                setSearching(false);
                return;
            }
        }
        
        // Если не нашли ничего, сохраняем адрес без координат
        onLocationSelect({
            formatted_address: address,
            address_components: {
                city: addressParts[0] || address,
                country: 'Serbia'
            }
        });
        
    } catch (error) {
        console.error("Error in geocoding address:", error);
        
        // При ошибке сохраняем адрес без координат
        const addressParts = address.split(',').map(part => part.trim()).filter(Boolean);
        onLocationSelect({
            formatted_address: address,
            address_components: {
                city: addressParts[0] || address,
                country: 'Serbia'
            }
        });
    } finally {
        setSearching(false);
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
                        leafletMapRef.current.setView([latitude, longitude], 15);
                        
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

            <Box sx={{ mb: isMobile ? 1 : 2, position: 'relative' }}>
                <TextField
                    id="location-search"
                    fullWidth
                    placeholder="Поиск по адресу..."
                    value={address}
                    onChange={handleAddressChange}
                    onKeyPress={(e) => {
                        if (e.key === 'Enter') {
                            handleSearch();
                        }
                    }}
                    onFocus={() => setShowSuggestions(addressSuggestions.length > 0)}
                    onBlur={() => {
                        // Задержка перед скрытием, чтобы успеть выбрать элемент при клике
                        setTimeout(() => setShowSuggestions(false), 200);
                    }}
                    size={isMobile ? "small" : "medium"}
                    autoComplete="off"
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <IconButton onClick={handleSearch}>
                                    {searching ? 
                                        <CircularProgress size={18} /> : 
                                        <SearchIcon fontSize={isMobile ? "small" : "medium"} />
                                    }
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
                
                {/* Выпадающий список подсказок */}
                <Collapse in={showSuggestions}>
                    <Paper 
                        sx={{ 
                            position: 'absolute', 
                            top: '100%', 
                            left: 0, 
                            right: 0, 
                            zIndex: 1000,
                            maxHeight: 200,
                            overflowY: 'auto',
                            mt: 0.5
                        }}
                        elevation={3}
                    >
                        <List dense>
                            {addressSuggestions.map((suggestion, index) => (
                                <ListItem 
                                    button 
                                    key={index} 
                                    onClick={() => handleSelectSuggestion(suggestion)}
                                >
                                    <ListItemText 
                                        primary={suggestion.formatted_address} 
                                    />
                                </ListItem>
                            ))}
                        </List>
                    </Paper>
                </Collapse>
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
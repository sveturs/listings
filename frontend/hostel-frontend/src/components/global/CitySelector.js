// frontend/hostel-frontend/src/components/global/CitySelector.js
import { MapPin, Navigation, Search } from 'lucide-react';
import React, { useState, useRef, useEffect } from 'react';
import {
    Box,
    Typography,
    Popover,
    List,
    ListItem,
    ListItemText,
    TextField,
    InputAdornment,
    Button,
    Divider,
    CircularProgress,
    Tooltip
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';

const CitySelector = ({ currentCity, onCityChange, isMobile = false }) => {
    const { t } = useTranslation('common');
    const [open, setOpen] = useState(false);
    const [searchValue, setSearchValue] = useState('');
    const [suggestions, setSuggestions] = useState([]);
    const [loading, setLoading] = useState(false);
    const anchorRef = useRef(null);

    const handleSearch = async (value) => {
        if (!value || value.length < 2) {
            setSuggestions([]);
            return;
        }

        setLoading(true);
        try {
            const response = await axios.get('/api/v1/cities/suggest', {
                params: { q: value, limit: 10 }
            });

            if (response.data?.data) {
                setSuggestions(response.data.data);
            } else {
                setSuggestions([]);
            }
        } catch (error) {
            console.error('Error fetching city suggestions:', error);
            setSuggestions([]);
        } finally {
            setLoading(false);
        }
    };

    // Обработка ввода с задержкой
    const handleSearchChange = (e) => {
        const value = e.target.value;
        setSearchValue(value);

        // Поиск с задержкой в 300мс
        const timerId = setTimeout(() => {
            handleSearch(value);
        }, 300);

        return () => clearTimeout(timerId);
    };

    // Популярные города в Сербии
    const popularCities = [
        { id: 1, city: 'Белград', country: 'Сербия', lat: 44.8178, lon: 20.4570 },
        { id: 2, city: 'Нови-Сад', country: 'Сербия', lat: 45.2671, lon: 19.8335 },
        { id: 3, city: 'Ниш', country: 'Сербия', lat: 43.3209, lon: 21.8958 },
        { id: 4, city: 'Крагуевац', country: 'Сербия', lat: 44.0128, lon: 20.9114 },
    ];

    const handleToggle = () => {
        setOpen((prevOpen) => !prevOpen);
        setSearchValue('');
        setSuggestions([]);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleSelectCity = (city) => {
        console.log(`Выбран город: ${city.city}, координаты: ${city.lat}, ${city.lon}`);
        onCityChange(city);
        handleClose();
    };

    const handleUseLocation = () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    const { latitude, longitude } = position.coords;
                    try {
                        const response = await axios.get('/api/v1/geocode/reverse', {
                            params: { lat: latitude, lon: longitude }
                        });

                        if (response.data?.data) {
                            onCityChange({
                                ...response.data.data,
                                lat: latitude,
                                lon: longitude
                            });
                        }
                    } catch (error) {
                        console.error('Error fetching city from coordinates:', error);
                    }
                    handleClose();
                },
                (error) => {
                    console.error('Error getting location:', error);
                }
            );
        }
    };

    useEffect(() => {
        if (currentCity) {
            // При изменении currentCity убедимся, что фильтры актуальны
            const checkFiltersSync = () => {
                // Получаем текущие параметры URL
                const urlParams = new URLSearchParams(window.location.search);
                const urlLat = urlParams.get('latitude');
                const urlLon = urlParams.get('longitude');

                // Если координаты в URL отличаются от currentCity, генерируем событие cityChanged
                if (urlLat !== String(currentCity.lat) || urlLon !== String(currentCity.lon)) {
                    const cityChangeEvent = new CustomEvent('cityChanged', {
                        detail: {
                            lat: currentCity.lat,
                            lon: currentCity.lon,
                            city: currentCity.city,
                            country: currentCity.country
                        }
                    });
                    window.dispatchEvent(cityChangeEvent);
                }
            };

            // Вызываем сразу и при изменении URL
            checkFiltersSync();
            window.addEventListener('popstate', checkFiltersSync);

            return () => {
                window.removeEventListener('popstate', checkFiltersSync);
            };
        }
    }, [currentCity]);
    return (
        <>
            <Tooltip title={t('location.changeCity', { defaultValue: 'Изменить город' })}>
                <Box
                    ref={anchorRef}
                    onClick={handleToggle}
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        cursor: 'pointer',
                        px: isMobile ? 0.5 : 1,
                        py: isMobile ? 0.3 : 0.5,
                        borderRadius: 1,
                        '&:hover': {
                            bgcolor: 'rgba(0,0,0,0.04)'
                        }
                    }}
                >
                    <MapPin size={18} />
                    {!isMobile && (
                        <Typography variant="body2" sx={{ mx: 0.5, maxWidth: 120, textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>
                            {currentCity?.city || t('location.selectCity', { defaultValue: 'Выбрать город' })}
                        </Typography>
                    )}
                </Box>
            </Tooltip>

            <Popover
                open={open}
                anchorEl={anchorRef.current}
                onClose={handleClose}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
                sx={{ mt: 1 }}
                PaperProps={{
                    sx: { width: 300, maxHeight: 500, overflow: 'hidden' }
                }}
            >
                <Box sx={{ p: 2 }}>
                    <Typography variant="subtitle1" gutterBottom>
                        {t('location.selectCity', { defaultValue: 'Выбрать город' })}
                    </Typography>

                    <TextField
                        fullWidth
                        size="small"
                        placeholder={t('location.searchCity', { defaultValue: 'Поиск города...' })}
                        value={searchValue}
                        onChange={handleSearchChange}
                        InputProps={{
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Search size={18} />
                                </InputAdornment>
                            ),
                            endAdornment: loading && (
                                <InputAdornment position="end">
                                    <CircularProgress size={20} />
                                </InputAdornment>
                            )
                        }}
                    />

                    <Button
                        fullWidth
                        variant="outlined"
                        sx={{ mt: 1 }}
                        onClick={handleUseLocation}
                        startIcon={<Navigation size={16} />}
                    >
                        {t('location.useCurrentLocation', { defaultValue: 'Использовать моё местоположение' })}
                    </Button>
                </Box>

                <Divider />

                <Box sx={{ maxHeight: 350, overflowY: 'auto' }}>
                    {searchValue.length > 1 ? (
                        <List dense>
                            {loading ? (
                                <ListItem>
                                    <ListItemText primary={t('common.loading', { defaultValue: 'Загрузка...' })} />
                                </ListItem>
                            ) : suggestions.length === 0 ? (
                                <ListItem>
                                    <ListItemText primary={t('location.noResults', { defaultValue: 'Ничего не найдено' })} />
                                </ListItem>
                            ) : (
                                suggestions.map((city) => (
                                    <ListItem
                                        button
                                        key={city.id || `${city.city}-${city.country}`}
                                        onClick={() => handleSelectCity(city)}
                                    >
                                        <ListItemText
                                            primary={city.city}
                                            secondary={city.country}
                                        />
                                    </ListItem>
                                ))
                            )}
                        </List>
                    ) : (
                        <>
                            <Typography variant="subtitle2" sx={{ px: 2, py: 1 }}>
                                {t('location.popularCities', { defaultValue: 'Популярные города' })}
                            </Typography>
                            <List dense>
                                {popularCities.map((city) => (
                                    <ListItem
                                        button
                                        key={city.id}
                                        onClick={() => handleSelectCity(city)}
                                    >
                                        <ListItemText
                                            primary={city.city}
                                            secondary={city.country}
                                        />
                                    </ListItem>
                                ))}
                            </List>
                        </>
                    )}
                </Box>
            </Popover>
        </>
    );
};

export default CitySelector;
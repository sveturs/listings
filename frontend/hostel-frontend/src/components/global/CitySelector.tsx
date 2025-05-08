// frontend/hostel-frontend/src/components/global/CitySelector.tsx
import { MapPin, Navigation, Search, Save } from 'lucide-react';
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
    Tooltip,
    Checkbox,
    FormControlLabel,
    Snackbar,
    Alert
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { useLocation } from '../../contexts/LocationContext';
import { useAuth } from '../../contexts/AuthContext';
// Import custom type declaration
import './tooltip-wrapper.d.ts';

interface CitySelectorProps {
    isMobile?: boolean;
}

interface City {
    id?: number | string;
    city: string;
    country: string;
    lat?: number;
    lon?: number;
}

const CitySelector: React.FC<CitySelectorProps> = ({ isMobile = false }) => {
    const { t } = useTranslation(['common', 'marketplace']);
    const { userLocation, setCity, detectUserLocation, isGeolocating } = useLocation();
    const { user } = useAuth();
    const isAuthenticated = !!user; // Вычисляем isAuthenticated на основе наличия пользователя
    const [open, setOpen] = useState<boolean>(false);
    const [searchValue, setSearchValue] = useState<string>('');
    const [suggestions, setSuggestions] = useState<City[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [saveToProfile, setSaveToProfile] = useState<boolean>(true);
    const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
    const [snackbarMessage, setSnackbarMessage] = useState<string>('');
    const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error' | 'warning' | 'info'>('success');
    const anchorRef = useRef<HTMLDivElement | null>(null);

    const handleSearch = async (value: string) => {
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
    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setSearchValue(value);

        // Поиск с задержкой в 300мс
        const timerId = setTimeout(() => {
            handleSearch(value);
        }, 300);

        return () => clearTimeout(timerId);
    };

    // Популярные города в Сербии
    const popularCities: City[] = [
        { id: 1, city: t('cities.belgrade', { defaultValue: 'Белград', ns: 'marketplace' }),
          country: t('countries.serbia', { defaultValue: 'Сербия', ns: 'marketplace' }), lat: 44.8178, lon: 20.4570 },
        { id: 2, city: t('cities.noviSad', { defaultValue: 'Нови-Сад', ns: 'marketplace' }),
          country: t('countries.serbia', { defaultValue: 'Сербия', ns: 'marketplace' }), lat: 45.2671, lon: 19.8335 },
        { id: 3, city: t('cities.nis', { defaultValue: 'Ниш', ns: 'marketplace' }),
          country: t('countries.serbia', { defaultValue: 'Сербия', ns: 'marketplace' }), lat: 43.3209, lon: 21.8958 },
        { id: 4, city: t('cities.kragujevac', { defaultValue: 'Крагуевац', ns: 'marketplace' }),
          country: t('countries.serbia', { defaultValue: 'Сербия', ns: 'marketplace' }), lat: 44.0128, lon: 20.9114 },
    ];

    const handleToggle = () => {
        setOpen((prevOpen) => !prevOpen);
        setSearchValue('');
        setSuggestions([]);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleSelectCity = async (city: City) => {
        // Явно устанавливаем город через контекст
        setCity({
            city: city.city,
            country: city.country,
            lat: city.lat,
            lon: city.lon
        });

        // Если пользователь авторизован и выбрал опцию сохранения в профиль
        if (isAuthenticated && saveToProfile) {
            try {
                await updateUserProfile(city);
                // Показываем уведомление об успешном сохранении
                setSnackbarMessage(t('location.savedToProfile', { defaultValue: 'Город сохранен в профиле', ns: 'marketplace' }));
                setSnackbarSeverity('success');
                setSnackbarOpen(true);
            } catch (error) {
                console.error('Error updating user profile:', error);
                // Показываем уведомление об ошибке
                setSnackbarMessage(t('location.errorSavingToProfile', { defaultValue: 'Ошибка сохранения города в профиле', ns: 'marketplace' }));
                setSnackbarSeverity('error');
                setSnackbarOpen(true);
            }
        }

        handleClose();
    };

    // Функция для обновления профиля пользователя
    const updateUserProfile = async (city: City) => {
        await axios.put('/api/v1/users/profile', {
            city: city.city,
            country: city.country
        });
    };

    const handleUseLocation = async () => {
        try {
            await detectUserLocation();

            // Если пользователь авторизован и выбрал опцию сохранения в профиль
            if (isAuthenticated && saveToProfile && userLocation) {
                // Создаем небольшую задержку, чтобы userLocation успел обновиться
                setTimeout(async () => {
                    try {
                        await updateUserProfile({
                            city: userLocation.city || '',
                            country: userLocation.country || ''
                        });
                        // Показываем уведомление об успешном сохранении
                        setSnackbarMessage(t('location.savedToProfile', { defaultValue: 'Город сохранен в профиле', ns: 'marketplace' }));
                        setSnackbarSeverity('success');
                        setSnackbarOpen(true);
                    } catch (error) {
                        console.error('Error updating user profile:', error);
                        // Показываем уведомление об ошибке
                        setSnackbarMessage(t('location.errorSavingToProfile', { defaultValue: 'Ошибка сохранения города в профиле', ns: 'marketplace' }));
                        setSnackbarSeverity('error');
                        setSnackbarOpen(true);
                    }
                }, 500);
            }

            handleClose();
        } catch (error) {
            console.error('Error getting location:', error);
        }
    };

    const handleSnackbarClose = (_event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }
        setSnackbarOpen(false);
    };

    return (
        <>
            {/* @ts-ignore */}
            <Tooltip
                title={t('location.changeCity', { defaultValue: 'Изменить город', ns: 'marketplace' })}
                arrow
                placement="bottom"
            >
                <Box component="span"
                    ref={anchorRef}
                    onClick={handleToggle}
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        cursor: 'pointer',
                        width: isMobile ? '40px' : 'auto',
                        height: isMobile ? '40px' : 'auto',
                        minWidth: isMobile ? '40px' : 'auto',
                        px: isMobile ? 0 : 1,
                        py: isMobile ? 0 : 0.5,
                        borderRadius: 1,
                        border: isMobile ? '1px solid rgba(0, 68, 148, 0.5)' : 'none',
                        '&:hover': {
                            bgcolor: 'rgba(0,0,0,0.04)',
                            borderColor: isMobile ? 'rgba(0, 68, 148, 0.8)' : 'none'
                        }
                    }}
                >
                    <MapPin size={18} color="#004494" />
                    {!isMobile && (
                        <Typography variant="body2" sx={{ mx: 0.5, maxWidth: 120, textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>
                            {userLocation?.city || t('location.selectCity', { defaultValue: 'Выбрать город', ns: 'marketplace' })}
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
                        {t('location.selectCity', { defaultValue: 'Выбрать город', ns: 'marketplace' })}
                    </Typography>

                    <TextField
                        fullWidth
                        size="small"
                        placeholder={t('location.searchCity', { defaultValue: 'Поиск города...', ns: 'marketplace' })}
                        value={searchValue}
                        onChange={handleSearchChange}
                        InputProps={{
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Search size={18} />
                                </InputAdornment>
                            ),
                            endAdornment: loading ? (
                                <InputAdornment position="end">
                                    <CircularProgress size={20} />
                                </InputAdornment>
                            ) : null
                        }}
                    />

                    <Button
                        fullWidth
                        variant="outlined"
                        sx={{ mt: 1 }}
                        onClick={handleUseLocation}
                        startIcon={<Navigation size={16} />}
                        disabled={isGeolocating}
                    >
                        {isGeolocating
                            ? t('location.detectingLocation', { defaultValue: 'Определение местоположения...', ns: 'marketplace' })
                            : t('listings.map.useMyLocation', { defaultValue: 'Использовать моё местоположение', ns: 'marketplace' })}
                    </Button>

                    {isAuthenticated && (
                        <FormControlLabel
                            control={
                                <Checkbox
                                    checked={saveToProfile}
                                    onChange={(e) => setSaveToProfile(e.target.checked)}
                                    size="small"
                                />
                            }
                            label={
                                <Typography variant="body2">
                                    {t('location.saveToProfile', { defaultValue: 'Сохранить в профиле', ns: 'marketplace' })}
                                </Typography>
                            }
                            sx={{ mt: 1 }}
                        />
                    )}
                </Box>

                <Divider />

                <Box sx={{ maxHeight: 350, overflowY: 'auto' }}>
                    {searchValue.length > 1 ? (
                        <List dense>
                            {loading ? (
                                <ListItem>
                                    <ListItemText primary={t('common.loading', { defaultValue: 'Загрузка...', ns: 'marketplace' })} />
                                </ListItem>
                            ) : suggestions.length === 0 ? (
                                <ListItem>
                                    <ListItemText primary={t('location.noResults', { defaultValue: 'Ничего не найдено', ns: 'marketplace' })} />
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
                                {t('location.popularCities', { defaultValue: 'Популярные города', ns: 'marketplace' })}
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

            <Snackbar
                open={snackbarOpen}
                autoHideDuration={4000}
                onClose={handleSnackbarClose}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                <Alert onClose={handleSnackbarClose} severity={snackbarSeverity} sx={{ width: '100%' }}>
                    {snackbarMessage}
                </Alert>
            </Snackbar>
        </>
    );
};

export default CitySelector;
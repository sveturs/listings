//frontend/hostel-frontend/src/pages/marketplace/MarketplacePage.js
import { useTranslation } from 'react-i18next';
import { Map as MapIcon } from 'lucide-react';
import { useEffect, useState, useCallback, useRef } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useLocation } from '../../contexts/LocationContext';
import { debounce } from 'lodash';

import {
    Container,
    Grid,
    Box,
    CircularProgress,
    Button,
    useTheme,
    useMediaQuery,
    IconButton,
    Alert,
    Paper,
    Chip,
    Fab,
    Tooltip,
}
    from '@mui/material';
import { Plus, Search, X, List } from 'lucide-react';
import ListingCard from '../../components/marketplace/ListingCard';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import {
    MobileFilters,
    MobileListingCard,
    MobileHeader,
} from '../../components/marketplace/MobileComponents';
import CompactMarketplaceFilters from '../../components/marketplace/MarketplaceFilters';
import MapView from '../../components/marketplace/MapView';
import axios from '../../api/axios';

const MobileListingGrid = ({ listings }) => {
    const navigate = useNavigate();

    return (
        <Box sx={{ px: 1 }}>
            <Grid container spacing={1}>
                {listings.map((listing, index) => {
                    // Создаем уникальный идентификатор из других полей, если ID = 0
                    const effectiveId = listing.id || `temp-${listing.category_id}-${listing.user_id}-${index}`;
                    return (
                        <Grid item xs={12} sm={6} md={4} key={effectiveId}>
                            <div onClick={() => {
                                if (listing.id) {
                                    navigate(`/marketplace/listings/${listing.id}`);
                                } else {
                                    // Используем стандартный поиск для получения объявления по другим параметрам
                                    const url = `/api/v1/marketplace/listings?category_id=${listing.category_id}&title=${encodeURIComponent(listing.title)}`;
                                    console.log("Переход к объявлению с временным URL:", url);
                                    // Можно показать уведомление пользователю, что фильтр применен
                                }
                            }}>
                                <ListingCard listing={listing} />
                            </div>
                        </Grid>
                    );
                })}
            </Grid>
        </Box>
    );
};

const MarketplacePage = () => {
    const { t } = useTranslation('marketplace');
    const { userLocation } = useLocation();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const navigate = useNavigate();
    const [searchParams, setSearchParams] = useSearchParams();

    const [listings, setListings] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isFilterOpen, setIsFilterOpen] = useState(false);
    const [categoryPath, setCategoryPath] = useState([]);
    const [mapViewActive, setMapViewActive] = useState(false);
    const [userLocationState, setUserLocationState] = useState(null);
    const [spellingSuggestion, setSpellingSuggestion] = useState(null);

    const [filters, setFilters] = useState({
        query: searchParams.get('query') || '',
        category_id: searchParams.get('category_id') || '',
        min_price: searchParams.get('min_price') || '',
        max_price: searchParams.get('max_price') || '',
        city: searchParams.get('city') || '',
        country: searchParams.get('country') || '',
        condition: searchParams.get('condition') || '',
        sort_by: searchParams.get('sort_by') || 'date_desc',
        distance: searchParams.get('distance') || '',
    });
    const fetchListings = useCallback(async (currentFilters = {}) => {
        try {
            // Предотвращаем повторные запросы с теми же параметрами
            const queryString = JSON.stringify(currentFilters);
            if (queryString === lastQueryRef.current) {
                return;
            }
            lastQueryRef.current = queryString;

            setLoading(true);
            setError(null);
            setSpellingSuggestion(null);

            const params = {};
            Object.entries(currentFilters).forEach(([key, value]) => {
                if (value !== '' && key !== 'city' && key !== 'country') {
                    if (key === 'query') {
                        params['q'] = value;
                    } else {
                        params[key] = value;
                    }
                }
            });

            console.log('Отправляем запрос:', params);
            const response = await axios.get('/api/v1/marketplace/search', { params });
            console.log('Получен ответ API:', response.data);

            // Улучшенная обработка данных с дополнительными проверками
            if (response.data && response.data.data) {
                if (Array.isArray(response.data.data)) {
                    console.log('Найдено объявлений:', response.data.data.length);
                    setListings(response.data.data);
                } else if (response.data.data.data && Array.isArray(response.data.data.data)) {
                    console.log('Найдено объявлений (вложенная структура):', response.data.data.data.length);
                    setListings(response.data.data.data);
                } else {
                    console.error('Данные не являются массивом:', response.data.data);
                    setListings([]);
                }
            } else {
                console.error('Ответ API не содержит ожидаемую структуру данных:', response.data);
                setListings([]);
            }

            if (response.data.meta && response.data.meta.spelling_suggestion) {
                setSpellingSuggestion(response.data.meta.spelling_suggestion);
            }
        } catch (err) {
            console.error('Ошибка при получении объявлений:', err);
            setError('Не удалось загрузить объявления');
            setListings([]);
        } finally {
            setLoading(false);
        }
    }, []); // Пустой массив зависимостей, так как функция не зависит от внешних переменных

    // Добавляем ref для отслеживания последнего запроса
    const lastQueryRef = useRef('');

    const debouncedFetchListings = useRef(
        debounce((filters) => {
            fetchListings(filters);
        }, 300)
    ).current;

    useEffect(() => {
        if (userLocation) {
            // Обновляем состояние фильтров с учетом нового местоположения
            setFilters(prevFilters => ({
                ...prevFilters,
                city: userLocation.city || '',
                country: userLocation.country || '',
                latitude: userLocation.lat,
                longitude: userLocation.lon
            }));

            // Обновляем параметры URL с новыми координатами
            setSearchParams(prevParams => {
                const nextParams = new URLSearchParams(prevParams);
                if (userLocation.city) nextParams.set('city', userLocation.city);
                if (userLocation.country) nextParams.set('country', userLocation.country);
                nextParams.set('latitude', userLocation.lat);
                nextParams.set('longitude', userLocation.lon);
                return nextParams;
            });

            // Выполняем поиск с учетом нового местоположения
            fetchListings({
                ...filters,
                city: userLocation.city || '',
                country: userLocation.country || '',
                latitude: userLocation.lat,
                longitude: userLocation.lon
            });
        }
    }, [userLocation]);

    // Проверяем, запрошен ли режим карты в URL
    useEffect(() => {
        const viewMode = searchParams.get('viewMode');
        if (viewMode === 'map') {
            setMapViewActive(true);
        } else if (viewMode === 'list') {
            setMapViewActive(false);
        }
    }, [searchParams]);

    useEffect(() => {
        const handleCityChange = (event) => {
            const { lat, lon, city, country } = event.detail;

            // Обновляем фильтры с новыми координатами
            setFilters(prev => ({
                ...prev,
                latitude: lat,
                longitude: lon,
                city: city || '',
                country: country || '',
                // Если distance не установлен, устанавливаем по умолчанию
                distance: prev.distance || '30km'
            }));

            // Обновляем параметры URL
            setSearchParams(prev => {
                const nextParams = new URLSearchParams(prev);
                nextParams.set('latitude', lat);
                nextParams.set('longitude', lon);
                if (city) nextParams.set('city', city);
                if (country) nextParams.set('country', country);
                if (!nextParams.has('distance')) nextParams.set('distance', '30km');
                return nextParams;
            });

            // Обновляем состояние для карты (если используется)
            if (typeof setUserLocationState === 'function') {
                setUserLocationState({
                    latitude: lat,
                    longitude: lon
                });
            }

            // Выполняем поиск с новыми координатами
            fetchListings({
                ...filters,
                latitude: lat,
                longitude: lon,
                city: city || '',
                country: country || '',
                distance: filters.distance || '30km'
            });
        };

        window.addEventListener('cityChanged', handleCityChange);
        return () => {
            window.removeEventListener('cityChanged', handleCityChange);
        };
    }, [filters, setSearchParams, fetchListings]);

    useEffect(() => {
        const fetchInitialData = async () => {
            try {
                // Проверяем наличие distance без координат
                const distanceParam = searchParams.get('distance');
                const latParam = searchParams.get('latitude');
                const lonParam = searchParams.get('longitude');

                // Удаляем этот код, так как он сбрасывает параметр distance, и вместо этого просто логируем
                if (distanceParam && (!latParam || !lonParam)) {
                    console.log("Параметр distance присутствует без координат, будет запрошена текущая геолокация");

                    // Проверим, поддерживается ли геолокация в браузере
                    if (navigator.geolocation) {
                        try {
                            // Запросим текущее местоположение пользователя
                            navigator.geolocation.getCurrentPosition(
                                (position) => {
                                    const latitude = position.coords.latitude;
                                    const longitude = position.coords.longitude;
                                    console.log(`Получены координаты: lat=${latitude}, lon=${longitude}`);

                                    // Обновим параметры с полученными координатами
                                    const updatedParams = new URLSearchParams(searchParams);
                                    updatedParams.set('latitude', latitude);
                                    updatedParams.set('longitude', longitude);
                                    setSearchParams(updatedParams);

                                    // Обновим состояние фильтров и выполним поиск с координатами
                                    setFilters(prev => ({
                                        ...prev,
                                        latitude: latitude,
                                        longitude: longitude
                                    }));

                                    // Сразу же обновим местоположение пользователя
                                    setUserLocationState({
                                        latitude: latitude,
                                        longitude: longitude
                                    });
                                },
                                (error) => {
                                    console.error("Ошибка получения геолокации:", error);
                                }
                            );
                        } catch (geoError) {
                            console.error("Ошибка при запросе геолокации:", geoError);
                        }
                    }
                }

                // Остальной код функции без изменений
                const categoriesResponse = await axios.get('/api/v1/marketplace/category-tree');

                if (categoriesResponse.data?.data) {
                    setCategories(categoriesResponse.data.data);
                }

                const initialFilters = {
                    query: searchParams.get('query') || '',
                    category_id: searchParams.get('category_id') || '',
                    min_price: searchParams.get('min_price') || '',
                    max_price: searchParams.get('max_price') || '',
                    city: searchParams.get('city') || '',
                    country: searchParams.get('country') || '',
                    condition: searchParams.get('condition') || '',
                    sort_by: searchParams.get('sort_by') || 'date_desc',
                    distance: searchParams.get('distance') || '',
                    latitude: searchParams.get('latitude') ? parseFloat(searchParams.get('latitude')) : null,
                    longitude: searchParams.get('longitude') ? parseFloat(searchParams.get('longitude')) : null,
                };

                setFilters(initialFilters);
                await fetchListings(initialFilters);
            } catch (err) {
                console.error('Error fetching initial data:', err);
                setError('Произошла ошибка при загрузке данных');
            }
        };

        fetchInitialData();
    }, [searchParams, fetchListings, setSearchParams]);

    useEffect(() => {
        if (!window.location.pathname.includes('/marketplace')) {
            navigate({
                pathname: '/marketplace',
                search: window.location.search
            }, { replace: true });
        }
    }, [navigate]);
    // Добавьте временную кнопку для тестирования






    const findCategoryPath = (categoryId, categoriesTree) => {
        if (!categoryId || !categoriesTree || categoriesTree.length === 0) {
            return [];
        }

        // Создаем плоскую карту всех категорий для быстрого поиска
        const categoryMap = new Map(); // Используем нативный Map, а не компонент из lucide-react

        const flattenCategories = (categories) => {
            for (const category of categories) {
                categoryMap.set(String(category.id), category);
                if (category.children && category.children.length > 0) {
                    flattenCategories(category.children);
                }
            }
        };

        // Заполняем карту всеми категориями
        flattenCategories(categoriesTree);

        // Строим путь от выбранной категории до корня
        const path = [];
        let currentId = String(categoryId);

        while (currentId) {
            const category = categoryMap.get(currentId);
            if (!category) break;

            // Добавляем категорию в начало пути
            path.unshift({
                id: category.id,
                name: category.name,
                slug: category.slug,
                translations: category.translations
            });

            // Переходим к родителю
            currentId = category.parent_id ? String(category.parent_id) : null;
        }

        return path;
    };

    useEffect(() => {
        if (filters.category_id && categories.length > 0) {
            const path = findCategoryPath(filters.category_id, categories);
            setCategoryPath(path);
        } else {
            setCategoryPath([]);
        }
    }, [filters.category_id, categories]);


    const handleFilterChange = useCallback((newFilters) => {
        setFilters(prev => {
            // Создаем обновленные фильтры
            const updated = { ...prev, ...newFilters };

            // Обновляем URL
            const nextParams = new URLSearchParams();
            Object.entries(updated).forEach(([key, value]) => {
                if (value !== null && value !== undefined && value !== '') {
                    nextParams.set(key, value);
                }
            });

            setSearchParams(nextParams);

            // Выполняем поиск с обновленными фильтрами
            fetchListings(updated);

            return updated;
        });
    }, [setSearchParams, fetchListings]);

    // Обработчик переключения режима просмотра (список/карта)
    const handleToggleMapView = useCallback(() => {
        const nextParams = new URLSearchParams(searchParams);
      
        if (mapViewActive) {
          // Переключаемся на список
          nextParams.set('viewMode', 'list');
          setMapViewActive(false);
        } else {
          // Переключаемся на карту
          nextParams.set('viewMode', 'map');
          setMapViewActive(true);
      
          // Используем координаты из userLocation, если они есть
          if (userLocation) {
            nextParams.set('latitude', userLocation.lat);
            nextParams.set('longitude', userLocation.lon);
            nextParams.set('distance', filters.distance || '5km');
            
            // Обновляем состояние для MapView
            setUserLocationState({
              latitude: userLocation.lat,
              longitude: userLocation.lon
            });
          }
          // Если нет данных в userLocation, но есть в фильтрах
          else if (filters.latitude && filters.longitude) {
            nextParams.set('latitude', filters.latitude);
            nextParams.set('longitude', filters.longitude);
            nextParams.set('distance', filters.distance || '5km');
          }
        }
      
        setSearchParams(nextParams);
      }, [mapViewActive, searchParams, setSearchParams, filters, userLocation]);

    const getActiveFiltersCount = () => {
        return Object.entries(filters).reduce((count, [key, value]) => {
            if (key !== 'sort_by' && value !== '') {
                return count + 1;
            }
            return count;
        }, 0);
    };

    // Добавьте эту функцию в MarketplacePage.js
    const resetAllFilters = () => {
        const nextParams = new URLSearchParams();
        if (searchParams.get('viewMode')) {
            nextParams.set('viewMode', searchParams.get('viewMode'));
        }
        setSearchParams(nextParams);

        const defaultFilters = {
            query: "",
            category_id: "",
            min_price: "",
            max_price: "",
            city: "",
            country: "",
            condition: "",
            sort_by: "date_desc",
            distance: "",
            latitude: null,
            longitude: null
        };

        setFilters(defaultFilters);
        fetchListings({});
    };

    const renderContent = () => {
        if (loading) {
            return (
                <Box display="flex" justifyContent="center" p={4}>
                    <CircularProgress />
                </Box>
            );
        }

        if (error) {
            return (
                <Alert severity="error" sx={{ m: 2 }}>
                    {error}
                </Alert>
            );
        }

        // Если активен режим карты
        if (mapViewActive) {
            return (
                <MapView
                    listings={listings}
                    userLocation={userLocationState}
                    filters={filters}
                    onFilterChange={handleFilterChange}
                    onMapClose={handleToggleMapView}
                />
            );
        }

        // Проверка, что listings - это массив
        if (!listings || !Array.isArray(listings) || listings.length === 0) {
            return (
                <Alert severity="info" sx={{ m: 2 }}>
                    No results found for your search
                </Alert>
            );
        }

        if (spellingSuggestion && listings.length === 0) {
            return (
                <>
                    <Alert
                        severity="info"
                        sx={{ m: 2 }}
                        action={
                            <Button
                                color="inherit"
                                size="small"
                                onClick={() => handleFilterChange({ query: spellingSuggestion })}
                            >
                                {t('search.usesuggestion')}
                            </Button>
                        }
                    >
                        {t('search.didyoumean')} <strong>{spellingSuggestion}</strong>?
                    </Alert>
                    <Box sx={{ m: 2, textAlign: 'center' }}>
                        {t('search.noresults')}
                    </Box>
                </>
            );
        }

        return isMobile ? (
            <MobileListingGrid listings={listings} />
        ) : (
            <Grid container spacing={3}>
                {listings.map((listing, index) => {
                    // Создаем уникальный идентификатор из других полей, если ID = 0
                    const effectiveId = listing.id || `temp-${listing.category_id}-${listing.user_id}-${index}`;
                    return (
                        <Grid item xs={12} sm={6} md={4} key={effectiveId}>
                            <div onClick={() => {
                                if (listing.id) {
                                    navigate(`/marketplace/listings/${listing.id}`);
                                } else {
                                    // Используем стандартный поиск для получения объявления по другим параметрам
                                    const url = `/api/v1/marketplace/listings?category_id=${listing.category_id}&title=${encodeURIComponent(listing.title)}`;
                                    console.log("Переход к объявлению с временным URL:", url);
                                    // Можно показать уведомление пользователю, что фильтр применен
                                }
                            }}>
                                <ListingCard listing={listing} />
                            </div>
                        </Grid>
                    );
                })}
            </Grid>
        );
    };

    if (isMobile) {
        return (
            <Box sx={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
                <MobileHeader
                    onOpenFilters={() => setIsFilterOpen(true)}
                    filtersCount={getActiveFiltersCount()}
                    onSearch={(query) => handleFilterChange({ query })}
                    searchValue={filters.query}
                />

                <Box sx={{
                    position: 'sticky',
                    top: 104,
                    zIndex: 1,
                    bgcolor: 'background.paper',
                    borderBottom: 1,
                    borderColor: 'divider'
                }}>
                    {/* Активные фильтры */}
                    {Object.entries(filters).some(([key, value]) => value && key !== 'sort_by') && (
                        <Box sx={{ px: 2, py: 1, display: 'flex', gap: 1, overflowX: 'auto', alignItems: 'center' }}>
                            {Object.entries(filters).map(([key, value]) => {
                                if (!value || key === 'sort_by' || key === 'latitude' || key === 'longitude') return null;
                                let label = value;
                                if (key === 'category_id') {
                                    const category = categories.find(c => String(c.id) === String(value));
                                    label = category ? category.name : value;
                                } else if (key === 'distance') {
                                    label = `Радиус: ${value}`;
                                }
                                return (
                                    <Chip
                                        key={key}
                                        label={label}
                                        size="small"
                                        onDelete={() => handleFilterChange({ [key]: '' })}
                                    />
                                );
                            })}
                            <Button
                                variant="outlined"
                                color="error"
                                size="small"
                                onClick={resetAllFilters}
                                sx={{ ml: 'auto', whiteSpace: 'nowrap' }}
                            >
                                {t('listings.filters.resetAll', { defaultValue: 'Сбросить всё' })}
                            </Button>
                        </Box>
                    )}
                </Box>

                <Box sx={{ flex: 1, bgcolor: mapViewActive ? 'transparent' : 'grey.50' }}>
                    {renderContent()}
                </Box>

                {/* Плавающая кнопка переключения режима просмотра */}
                <Fab
                    color={mapViewActive ? "default" : "primary"}
                    sx={{ position: 'fixed', bottom: 16, right: 16, zIndex: 1050 }}
                    onClick={handleToggleMapView}
                >
                    {mapViewActive ? <List /> : <MapIcon />}
                </Fab>

                <MobileFilters
                    open={isFilterOpen}
                    onClose={() => setIsFilterOpen(false)}
                    filters={filters}
                    onFilterChange={handleFilterChange}
                    categories={categories}
                    onToggleMapView={handleToggleMapView}
                />
            </Box>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    mb: 0
                }}
            >

                <Breadcrumbs paths={categoryPath} categories={categories} />


            </Box>


            <Grid container spacing={3}>
                <Grid item xs={12} md={3}>
                    <CompactMarketplaceFilters
                        filters={filters}
                        onFilterChange={handleFilterChange}
                        categories={categories}
                        selectedCategoryId={filters.category_id}
                        isLoading={loading}
                        onToggleMapView={handleToggleMapView}
                    />
                </Grid>
                <Grid item xs={12} md={9}>
                    {renderContent()}
                </Grid>
            </Grid>
        </Container>
    );
};

export default MarketplacePage;
// frontend/hostel-frontend/src/pages/marketplace/MarketplacePage.js

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
    ToggleButtonGroup,
    ToggleButton,
    Typography
}
    from '@mui/material';
import { Plus, Search, X, List, Grid as GridIcon } from 'lucide-react';
import ListingCard from '../../components/marketplace/ListingCard';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import {
    MobileFilters,
    MobileListingCard,
    MobileHeader,
} from '../../components/marketplace/MobileComponents';
import CompactMarketplaceFilters from '../../components/marketplace/MarketplaceFilters';
import CentralAttributeFilters from '../../components/marketplace/CentralAttributeFilters';
import MapView from '../../components/marketplace/MapView';
import axios from '../../api/axios';
import InfiniteScroll from '../../components/marketplace/InfiniteScroll';
import MarketplaceListingsList from '../../components/marketplace/MarketplaceListingsList';

const MobileListingGrid = ({ listings, viewMode = 'grid' }) => {
    const navigate = useNavigate();
    const { t } = useTranslation('marketplace');

    // Проверяем, какой режим отображения выбран
    if (viewMode === 'list') {
        return (
            <Box sx={{ px: 1 }}>
                {listings.map((listing, index) => {
                    const effectiveId = listing.id || `temp-${listing.category_id}-${listing.user_id}-${index}`;
                    return (
                        <Box
                            key={effectiveId}
                            onClick={() => {
                                if (listing.id) {
                                    navigate(`/marketplace/listings/${listing.id}`);
                                } else {
                                    const url = `/api/v1/marketplace/listings?category_id=${listing.category_id}&title=${encodeURIComponent(listing.title)}`;
                                    console.log("Переход к объявлению с временным URL:", url);
                                }
                            }}
                            sx={{ mb: 1 }}
                        >
                            <MobileListingCard listing={listing} viewMode="list" />
                        </Box>
                    );
                })}

                {listings.length === 0 && (
                    <Box sx={{ py: 4, textAlign: 'center' }}>
                        <Typography variant="body2" color="text.secondary">
                            {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
                        </Typography>
                    </Box>
                )}
            </Box>
        );
    }

    // Режим сетки (grid) - стандартное отображение
    return (
        <Grid container spacing={1}>
            {listings.map((listing, index) => {
                const effectiveId = listing.id || `temp-${listing.category_id}-${listing.user_id}-${index}`;
                return (
                    <Grid item xs={6} key={effectiveId}>
                        <Box
                            onClick={() => {
                                if (listing.id) {
                                    navigate(`/marketplace/listings/${listing.id}`);
                                } else {
                                    const url = `/api/v1/marketplace/listings?category_id=${listing.category_id}&title=${encodeURIComponent(listing.title)}`;
                                    console.log("Переход к объявлению с временным URL:", url);
                                }
                            }}
                        >
                            <ListingCard listing={listing} isMobile={true} />
                        </Box>
                    </Grid>
                );
            })}

            {listings.length === 0 && (
                <Grid item xs={12}>
                    <Box sx={{ py: 4, textAlign: 'center' }}>
                        <Typography variant="body2" color="text.secondary">
                            {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
                        </Typography>
                    </Box>
                </Grid>
            )}
        </Grid>
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
    const [viewMode, setViewMode] = useState('grid');

    // Состояния для пагинации и загрузки дополнительных объявлений
    const [page, setPage] = useState(1);
    const [hasMoreListings, setHasMoreListings] = useState(true);
    const [totalListings, setTotalListings] = useState(0);
    const [loadingMore, setLoadingMore] = useState(false);

    // Состояния для сортировки
    const [sortField, setSortField] = useState('created_at');
    const [sortDirection, setSortDirection] = useState('desc');

    // Добавляем ref для отслеживания последнего запроса
    const lastQueryRef = useRef('');

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
        attributeFilters: {}
    });

    const fetchListings = useCallback(async (currentFilters = {}, isLoadMore = false) => {
        try {
            // Предотвращаем повторные запросы с теми же параметрами
            const queryString = JSON.stringify(currentFilters);
            if (queryString === lastQueryRef.current && !isLoadMore) {
                return;
            }

            if (!isLoadMore) {
                lastQueryRef.current = queryString;
                setLoading(true);
                setError(null);
                setSpellingSuggestion(null);
            } else {
                setLoadingMore(true);
            }

            const params = {};

            // Обрабатываем основные фильтры
            Object.entries(currentFilters).forEach(([key, value]) => {
                if (value !== '' && key !== 'city' && key !== 'country' && key !== 'attributeFilters') {
                    if (key === 'query') {
                        params['q'] = value;
                    } else {
                        params[key] = value;
                    }
                }
            });

            // Явно передаем параметры сортировки
            if (currentFilters.sort_by) {
                params.sort_by = currentFilters.sort_by;
            }

            // Добавляем атрибуты, если они есть
            if (currentFilters.attributeFilters && typeof currentFilters.attributeFilters === 'object') {
                Object.entries(currentFilters.attributeFilters).forEach(([attrKey, attrValue]) => {
                    if (attrValue) {
                        params[`attr_${attrKey}`] = attrValue;
                        console.log(`Добавлен атрибут в запрос: attr_${attrKey}=${attrValue}`);
                    }
                });
            }

            // Добавляем параметры пагинации
            params.page = isLoadMore ? page + 1 : 1;
            params.size = 20; // Размер страницы

            console.log('Отправляем запрос поиска объявлений с параметрами:', params);
            const response = await axios.get('/api/v1/marketplace/search', { params });
            console.log('Получен ответ API:', response.data);

            // Проверка на наличие данных
            if (!response.data) {
                console.error('Ответ API не содержит данных');
                setListings([]);
                setHasMoreListings(false);
                return;
            }

            // Получаем данные о результатах и пагинации
            let newListings = [];
            let totalCount = 0;

            if (response.data.data && Array.isArray(response.data.data)) {
                newListings = response.data.data;
                totalCount = response.data.meta?.total || newListings.length;
            } else if (response.data.data?.data && Array.isArray(response.data.data.data)) {
                newListings = response.data.data.data;
                totalCount = response.data.data.meta?.total || newListings.length;
            } else if (response.data.data?.items && Array.isArray(response.data.data.items)) {
                newListings = response.data.data.items;
                totalCount = response.data.data?.total || newListings.length;
            }

            // Обновляем состояние списка объявлений
            if (isLoadMore) {
                setListings(prev => [...prev, ...newListings]);
                setPage(prev => prev + 1);
            } else {
                setListings(newListings);
                setPage(1);
            }

            setTotalListings(totalCount);

            // Проверяем, есть ли еще объявления для загрузки
            setHasMoreListings(listings.length + newListings.length < totalCount);

            // Обработка подсказок для исправления опечаток
            if (response.data.meta && response.data.meta.spelling_suggestion) {
                setSpellingSuggestion(response.data.meta.spelling_suggestion);
            }

        } catch (err) {
            console.error('Ошибка при получении объявлений:', err);
            setError('Не удалось загрузить объявления');
            setListings(isLoadMore ? [...listings] : []);
            setHasMoreListings(false);
        } finally {
            setLoading(false);
            setLoadingMore(false);
        }
    }, [page, listings]);

    const debouncedFetchListings = useRef(
        debounce((filters) => {
            fetchListings(filters);
        }, 300)
    ).current;

    // Обработчик изменения фильтров
    const handleFilterChange = useCallback((newFilters) => {
        console.log(`MarketplaceFilters: Выбраны фильтры:`, newFilters);

        setFilters(prev => {
            // Создаем обновленные фильтры
            let updated = { ...prev };

            // Обрабатываем attributeFilters особым образом
            if (newFilters.attributeFilters) {
                updated.attributeFilters = {
                    ...(prev.attributeFilters || {}),
                    ...newFilters.attributeFilters
                };

                // Удаляем attributeFilters из newFilters, чтобы избежать дублирования
                const { attributeFilters, ...restFilters } = newFilters;
                Object.assign(updated, restFilters);
            } else {
                // Если это обычное обновление фильтров, просто обновляем
                Object.assign(updated, newFilters);
            }

            // Обновляем URL параметры
            const nextParams = new URLSearchParams();

            // Добавляем стандартные фильтры
            Object.entries(updated).forEach(([key, value]) => {
                if (value !== null && value !== undefined && value !== '' && key !== 'attributeFilters') {
                    nextParams.set(key, value);
                }
            });

            // Добавляем атрибуты с префиксом attr_
            if (updated.attributeFilters) {
                Object.entries(updated.attributeFilters).forEach(([attrKey, attrValue]) => {
                    if (attrValue) {
                        nextParams.set(`attr_${attrKey}`, attrValue);
                    }
                });
            }

            setSearchParams(nextParams);

            // Вызываем запрос с обновленными фильтрами
            debouncedFetchListings(updated);

            return updated;
        });
    }, [debouncedFetchListings, setSearchParams]);

    // Обработчик сортировки
    const handleSortChange = (field, direction) => {
        console.log(`Sorting changed: ${field} ${direction}`);
        setSortField(field);
        setSortDirection(direction);

        // Преобразуем поле сортировки в формат, понятный API
        let apiSortField;
        switch (field) {
            case 'created_at':
                apiSortField = 'date';
                break;
            case 'title':
                apiSortField = 'title';
                break;
            case 'price':
                apiSortField = 'price';
                break;
            case 'location':
                apiSortField = 'location';
                break;
            default:
                apiSortField = 'date';
        }

        // Создаем новый объект фильтров с обновленной сортировкой
        const updatedFilters = {
            ...filters,
            sort_by: `${apiSortField}_${direction}`
        };

        // Обновляем URL-параметры
        const nextParams = new URLSearchParams(searchParams);
        nextParams.set('sort_by', `${apiSortField}_${direction}`);
        setSearchParams(nextParams);

        // Обновляем состояние фильтров
        setFilters(updatedFilters);

        // Сбрасываем данные для загрузки с первой страницы
        setListings([]);
        setPage(1);
        setHasMoreListings(true);

        // Выполняем запрос с новыми фильтрами
        fetchListings(updatedFilters);
    };

    // Обработчик сброса атрибутных фильтров
    const resetAttributeFilters = useCallback(() => {
        handleFilterChange({ attributeFilters: {} });
    }, [handleFilterChange]);

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

    // Обработчик переключения режима отображения (сетка/список)
    const handleViewModeChange = (event, newMode) => {
        if (newMode !== null) {
            setViewMode(newMode);

            // Опционально: сохраняем предпочтение пользователя в localStorage
            localStorage.setItem('marketplace-view-mode', newMode);
        }
    };

    // Обработчик загрузки дополнительных объявлений
    const handleLoadMore = () => {
        if (loadingMore || !hasMoreListings) return;
        fetchListings(filters, true);
    };

    const getActiveFiltersCount = () => {
        return Object.entries(filters).reduce((count, [key, value]) => {
            if (key !== 'sort_by' && value !== '') {
                return count + 1;
            }
            return count;
        }, 0);
    };

    // Функция сброса всех фильтров
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

    useEffect(() => {
        // Загружаем сохраненный режим просмотра из localStorage
        const savedViewMode = localStorage.getItem('marketplace-view-mode');
        if (savedViewMode && (savedViewMode === 'grid' || savedViewMode === 'list')) {
            setViewMode(savedViewMode);
        }
    }, []);

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
                // Загрузка категорий
                const categoriesResponse = await axios.get('/api/v1/marketplace/category-tree');
                if (categoriesResponse.data?.data) {
                    setCategories(categoriesResponse.data.data);
                }

                // Извлекаем атрибуты из URL
                const attributeFilters = {};
                searchParams.forEach((value, key) => {
                    if (key.startsWith('attr_')) {
                        const attrName = key.substring(5); // Удаляем префикс 'attr_'
                        attributeFilters[attrName] = value;
                    }
                });

                // Создаем объект с начальными фильтрами
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
                    longitude: searchParams.get('longitude') ? parseFloat(searchParams.get('longitude')) : null
                };

                // Добавляем атрибуты, только если они есть
                if (Object.keys(attributeFilters).length > 0) {
                    initialFilters.attributeFilters = attributeFilters;
                }

                console.log('Начальные фильтры с атрибутами:', initialFilters);

                // Устанавливаем начальные фильтры
                setFilters(initialFilters);

                // Выполняем первоначальный запрос данных
                await fetchListings(initialFilters);
            } catch (err) {
                console.error('Error fetching initial data:', err);
                setError('Произошла ошибка при загрузке данных');
            }
        };

        fetchInitialData();

        // Важно! Не добавляем fetchListings в зависимости
    }, [searchParams, setSearchParams]);

    useEffect(() => {
        if (!window.location.pathname.includes('/marketplace')) {
            navigate({
                pathname: '/marketplace',
                search: window.location.search
            }, { replace: true });
        }
    }, [navigate]);

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

    const renderContent = () => {
        // Отображаем фильтры атрибутов до проверки наличия объявлений
        const categoryFilters = filters.category_id ? (
            <CentralAttributeFilters
                categoryId={filters.category_id}
                onFilterChange={(newAttrFilters) => {
                    console.log("CentralAttributeFilters вызвал onFilterChange с:", newAttrFilters);
                    handleFilterChange({ attributeFilters: newAttrFilters });
                }}
                filters={filters.attributeFilters || {}}
                resetAttributeFilters={resetAttributeFilters}
            />
        ) : null;

        if (loading) {
            return (
                <>
                    {categoryFilters}
                    <Box display="flex" justifyContent="center" p={4}>
                        <CircularProgress />
                    </Box>
                </>
            );
        }

        if (error) {
            return (
                <>
                    {categoryFilters}
                    <Alert severity="error" sx={{ m: 2 }}>
                        {error}
                    </Alert>
                </>
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
            if (spellingSuggestion) {
                return (
                    <>
                        {categoryFilters}
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

            return (
                <>
                    {categoryFilters}
                    <Alert severity="info" sx={{ m: 2 }}>
                        {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
                    </Alert>
                </>
            );
        }

        return (
            <>
                {categoryFilters}

                                {!mapViewActive && (
                    <InfiniteScroll
                        hasMore={hasMoreListings}
                        loading={loadingMore}
                        onLoadMore={handleLoadMore}
                        autoLoad={!isMobile} // Автозагрузка только на десктопе
                        loadingMessage={t('listings.loading', { defaultValue: 'Загрузка...' })}
                        loadMoreButtonText={t('listings.loadMore', { defaultValue: 'Показать ещё' })}
                        noMoreItemsText={t('listings.noMoreListings', { defaultValue: 'Больше нет объявлений' })}
                    >
                        {isMobile ? (
                            <MobileListingGrid listings={listings} viewMode={viewMode} />
                        ) : viewMode === 'grid' ? (
                            <Grid container spacing={3}>
                                {listings.map((listing, index) => {
                                    const effectiveId = listing.id || `temp-${listing.category_id}-${listing.user_id}-${index}`;
                                    return (
                                        <Grid item xs={12} sm={6} md={4} key={effectiveId}>
                                            <Box onClick={() => {
                                                if (listing.id) {
                                                    navigate(`/marketplace/listings/${listing.id}`);
                                                } else {
                                                    const url = `/api/v1/marketplace/listings?category_id=${listing.category_id}&title=${encodeURIComponent(listing.title)}`;
                                                    console.log("Переход к объявлению с временным URL:", url);
                                                }
                                            }}>
                                                <ListingCard listing={listing} />
                                            </Box>
                                        </Grid>
                                    );
                                })}
                            </Grid>
                        ) : (
                            <MarketplaceListingsList
                                listings={listings}
                                showSelection={false}
                                onSortChange={handleSortChange}
                                filters={filters}
                            />
                        )}
                    </InfiniteScroll>
                )}
            </>
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
                    // Добавляем два новых пропса:
                    viewMode={viewMode}
                    onViewModeChange={(e, newMode) => {
                        if (newMode !== null) {
                            setViewMode(newMode);
                            // Сохраняем предпочтение пользователя в localStorage
                            localStorage.setItem('marketplace-view-mode', newMode);
                        }
                    }}
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
                                if (!value || key === 'sort_by' || key === 'latitude' || key === 'longitude' || key === 'attributeFilters') return null;

                                let label = '';

                                // Преобразуем значение в строку с учетом специфики каждого поля
                                if (key === 'category_id') {
                                    const category = categories.find(c => String(c.id) === String(value));
                                    label = category ? category.name : String(value);
                                } else if (key === 'distance') {
                                    label = `Радиус: ${value}`;
                                } else {
                                    // Убедимся, что любое другое значение преобразуется в строку
                                    label = String(value);
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
                    {/* Добавляем CentralAttributeFilters здесь для мобильной версии */}
                    {filters.category_id && !mapViewActive && (
                        <Box sx={{ px: 2, pt: 2 }}>
                            <CentralAttributeFilters
                                categoryId={filters.category_id}
                                onFilterChange={(newAttrFilters) => {
                                    console.log("CentralAttributeFilters вызвал onFilterChange с:", newAttrFilters);
                                    handleFilterChange({ attributeFilters: newAttrFilters });
                                }}
                                filters={filters.attributeFilters || {}}  // Гарантированно передаем объект
                                resetAttributeFilters={resetAttributeFilters}
                            />
                        </Box>
                    )}

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
                        selectedCategoryId={filters.category_id}
                        onToggleMapView={handleToggleMapView}
                        setSearchParams={setSearchParams}
                        fetchListings={fetchListings}
                        viewMode={viewMode}
                        handleViewModeChange={handleViewModeChange}
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
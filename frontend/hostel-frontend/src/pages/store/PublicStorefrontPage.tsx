// frontend/hostel-frontend/src/pages/store/PublicStorefrontPage.tsx
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, Link, useLocation } from 'react-router-dom';
import StorefrontRating from '../../components/store/StorefrontRating';

import {
    Container,
    Typography,
    Grid,
    Box,
    CircularProgress,
    Alert,
    Card,
    CardContent,
    Divider,
    Paper,
    Avatar,
    Button,
    Chip,
    Stack,
    ToggleButtonGroup,
    ToggleButton,
    useTheme,
    useMediaQuery
} from '@mui/material';
import { Store, Phone, Mail, Globe, MapPin, Grid as GridIcon, List } from 'lucide-react';
import axios from '../../api/axios';
import ListingCard from '../../components/marketplace/ListingCard';
import MarketplaceListingsList from '../../components/marketplace/MarketplaceListingsList';
import MiniMap from '../../components/maps/MiniMap';
import InfiniteScroll from '../../components/marketplace/InfiniteScroll';

interface Storefront {
    id: number;
    name: string;
    description: string;
    slug: string;
    phone: string;
    email: string;
    website: string;
    address: string;
    city: string;
    country: string;
    latitude: number | null;
    longitude: number | null;
    logo_path: string | null;
    created_at: string;
    [key: string]: any;
}

interface Listing {
    id: number;
    title: string;
    description: string;
    price: number;
    images: any[];
    status: string;
    created_at: string;
    category_id: number;
    user_id: number;
    city?: string;
    country?: string;
    [key: string]: any;
}

interface ResponseData<T> {
    data: T;
    meta?: {
        total: number;
        per_page: number;
        current_page: number;
        last_page: number;
    };
}

type RouteParams = {
    id: string;
}

type SortField = 'created_at' | 'title' | 'price' | 'reviews';
type SortDirection = 'asc' | 'desc';
type ViewMode = 'grid' | 'list';

interface SortParams {
    field?: SortField;
    direction?: SortDirection;
}

const PublicStorefrontPage: React.FC = () => {
    const { t } = useTranslation(['common', 'marketplace']);
    const { id } = useParams<keyof RouteParams>();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const location = useLocation();

    // Получаем параметр highlightedListingId из URL, если он есть
    const queryParams = new URLSearchParams(location.search);
    const highlightedListingId = queryParams.get('highlightedListingId');

    const [loading, setLoading] = useState<boolean>(true);
    const [storefront, setStorefront] = useState<Storefront | null>(null);
    const [storeListings, setStoreListings] = useState<Listing[]>([]);
    const [highlightedListing, setHighlightedListing] = useState<Listing | null>(null);
    const [error, setError] = useState<string | null>(null);
    
    // Добавляем состояния для пагинации и бесконечной прокрутки
    const [page, setPage] = useState<number>(1);
    const [hasMoreListings, setHasMoreListings] = useState<boolean>(true);
    const [loadingMore, setLoadingMore] = useState<boolean>(false);
    const [totalListings, setTotalListings] = useState<number>(0);
    
    // Состояния для сортировки
    const [sortField, setSortField] = useState<SortField>('created_at');
    const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
    
    // Состояние для режима отображения (grid/list)
    const [viewMode, setViewMode] = useState<ViewMode>(() => {
        // Загружаем сохраненный режим просмотра из localStorage
        const savedViewMode = localStorage.getItem('storefront-view-mode');
        return savedViewMode === 'list' ? 'list' : 'grid';
    });

    // Функция для загрузки начальных данных с учетом сортировки
    const fetchInitialData = async (sortParams: SortParams = {}): Promise<void> => {
        try {
            setLoading(true);
            setError(null);
            
            // Формируем параметры API сортировки
            let apiSort: string = '';
            const currentSortField = sortParams.field || sortField;
            const currentSortDirection = sortParams.direction || sortDirection;
            
            switch (currentSortField) {
                case 'created_at':
                    apiSort = `date_${currentSortDirection}`;
                    break;
                case 'title':
                    apiSort = `title_${currentSortDirection}`;
                    break;
                case 'price':
                    apiSort = `price_${currentSortDirection}`;
                    break;
                case 'reviews':
                    apiSort = `rating_${currentSortDirection}`;
                    break;
                default:
                    apiSort = `date_${currentSortDirection}`;
            }
            
            console.log(`Запрос с сортировкой: ${apiSort}, поле=${currentSortField}, направление=${currentSortDirection}`);
            
            // Определяем параметры промисов для загрузки данных
            const promises = [
                axios.get(`/api/v1/public/storefronts/${id}`),
                axios.get('/api/v1/marketplace/listings', {
                    params: { 
                        storefront_id: id,
                        page: 1,
                        size: 1000, // Увеличиваем размер, чтобы загружать все объявления для отображения на карте
                        sort_by: apiSort, // Добавляем параметр сортировки
                        view_mode: 'map' // Указываем, что данные нужны для карты
                    }
                })
            ];

            // Если есть ID выделенного объявления, добавляем промис для его загрузки
            if (highlightedListingId) {
                promises.push(axios.get(`/api/v1/marketplace/listings/${highlightedListingId}`));
            }

            // Выполняем все промисы параллельно
            const responses = await Promise.all(promises);
            
            const storefrontResponse = responses[0];
            const listingsResponse = responses[1];
            
            setStorefront(storefrontResponse.data.data);
            
            // Обрабатываем результаты первой страницы
            const listings: Listing[] = listingsResponse.data.data?.data || [];
            
            // Если был запрос на выделенное объявление, получаем его
            if (highlightedListingId && responses.length > 2) {
                const highlightedListingResponse = responses[2];
                console.log('Получены данные выделенного объявления:', highlightedListingResponse.data);
                
                if (highlightedListingResponse.data) {
                    // API может возвращать данные в разных форматах
                    const listingData = highlightedListingResponse.data.data || highlightedListingResponse.data;
                    setHighlightedListing(listingData);
                    
                    // Удаляем выделенное объявление из обычного списка, если оно там есть
                    // (чтобы избежать дублирования)
                    const filteredListings = listings.filter(
                        listing => listing.id !== parseInt(highlightedListingId)
                    );
                    setStoreListings(filteredListings);
                    
                    console.log('Выделенное объявление установлено:', listingData.id);
                } else {
                    console.log('Данных выделенного объявления нет в ответе API');
                    setStoreListings(listings);
                }
            } else {
                console.log('Нет запроса на выделенное объявление или нет ответа');
                setStoreListings(listings);
            }
            
            // Обновляем информацию о пагинации
            if (listingsResponse.data.data?.meta) {
                setTotalListings(listingsResponse.data.data.meta.total || 0);
                setHasMoreListings(listings.length < listingsResponse.data.data.meta.total);
            } else {
                setHasMoreListings(listings.length >= 20); // Если размер страницы - 20 элементов
            }
            
            // Обновляем состояние сортировки если переданы новые параметры
            if (sortParams.field) setSortField(sortParams.field);
            if (sortParams.direction) setSortDirection(sortParams.direction);
            
            setPage(1);
        } catch (err) {
            console.error('Error fetching storefront data:', err);
            setError(t('marketplace:listings.errors.loadFailed'));
        } finally {
            setLoading(false);
        }
    };

    // Функция для загрузки дополнительных объявлений
    const fetchMoreListings = async (): Promise<void> => {
        if (!hasMoreListings || loadingMore) return;
        
        try {
            setLoadingMore(true);
            const nextPage = page + 1;
            
            // Формируем параметры API сортировки
            let apiSort: string = '';
            switch (sortField) {
                case 'created_at':
                    apiSort = `date_${sortDirection}`;
                    break;
                case 'title':
                    apiSort = `title_${sortDirection}`;
                    break;
                case 'price':
                    apiSort = `price_${sortDirection}`;
                    break;
                case 'reviews':
                    apiSort = `rating_${sortDirection}`;
                    break;
                default:
                    apiSort = `date_${sortDirection}`;
            }
            
            console.log(`Загрузка дополнительных объявлений с сортировкой: ${apiSort}`);
            
            const listingsResponse = await axios.get('/api/v1/marketplace/listings', {
                params: { 
                    storefront_id: id,
                    page: nextPage,
                    size: 1000, // Увеличиваем размер для загрузки большего количества объявлений
                    sort_by: apiSort, // Добавляем параметр сортировки
                    view_mode: 'map' // Указываем, что данные нужны для карты
                }
            });
            
            const newListings: Listing[] = listingsResponse.data.data?.data || [];
            
            if (newListings.length > 0) {
                setStoreListings(prev => [...prev, ...newListings]);
                setPage(nextPage);
                
                // Проверяем, есть ли еще объявления для загрузки
                if (listingsResponse.data.data?.meta) {
                    const total = listingsResponse.data.data.meta.total || 0;
                    setHasMoreListings(storeListings.length + newListings.length < total);
                } else {
                    setHasMoreListings(newListings.length >= 20);
                }
            } else {
                setHasMoreListings(false);
            }
        } catch (err) {
            console.error('Error fetching more listings:', err);
        } finally {
            setLoadingMore(false);
        }
    };

    useEffect(() => {
        if (id) {
            fetchInitialData();
        }
    }, [id, t]);

    // Обработчик для бесконечной прокрутки
    const handleLoadMore = (): void => {
        fetchMoreListings();
    };
    
    // Обработчик переключения режима отображения (сетка/список)
    const handleViewModeChange = (_event: React.MouseEvent<HTMLElement>, newMode: ViewMode | null): void => {
        if (newMode !== null) {
            setViewMode(newMode);
            // Сохраняем предпочтение пользователя в localStorage
            localStorage.setItem('storefront-view-mode', newMode);
        }
    };
    
    // Обработчик изменения сортировки
    const handleSortChange = (field: SortField, direction: SortDirection): void => {
        console.log(`Обработчик сортировки получил: поле=${field}, направление=${direction}`);
        
        // Сбрасываем предыдущие объявления и состояние пагинации
        setStoreListings([]);
        setPage(1);
        setHasMoreListings(true);
        
        // Запрашиваем данные с новыми параметрами сортировки
        fetchInitialData({ field, direction });
    };

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
                    <CircularProgress />
                </Box>
            </Container>
        );
    }

    if (error || !storefront) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error || t('common:common.notFound')}
                </Alert>
            </Container>
        );
    }

    const hasLocation = storefront.latitude && storefront.longitude;
    const hasContactInfo = storefront.phone || storefront.email || storefront.website;

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Paper elevation={2} sx={{ p: 3, mb: 4 }}>
                <Box display="flex" alignItems="center" gap={2} mb={2}>
                    <Avatar sx={{ bgcolor: 'primary.main', width: 64, height: 64 }}>
                        {storefront.logo_path ? (
                            <img
                                src={`${process.env.REACT_APP_BACKEND_URL}/uploads/${storefront.logo_path}`}
                                alt={storefront.name}
                                style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                            />
                        ) : (
                            <Store size={32} />
                        )}
                    </Avatar>
                    <Box>
                        <Typography variant="h4" component="h1">
                            {storefront.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            {t('marketplace:store.openSince', { date: new Date(storefront.created_at).toLocaleDateString() })}
                        </Typography>
                    </Box>
                </Box>

                {storefront.description && (
                    <Typography variant="body1" paragraph>
                        {storefront.description}
                    </Typography>
                )}

                <Grid container spacing={3} sx={{ mt: 1 }}>
                    {hasContactInfo && (
                        <Grid item xs={12} md={hasLocation ? 6 : 12}>
                            <Card variant="outlined">
                                <CardContent>
                                    <Typography variant="h6" gutterBottom>
                                        {t('marketplace:store.contactInfo')}
                                    </Typography>
                                    <Stack spacing={1}>
                                        {storefront.phone && (
                                            <Box display="flex" alignItems="center" gap={1}>
                                                <Phone size={18} />
                                                <Typography>
                                                    <a href={`tel:${storefront.phone}`}>{storefront.phone}</a>
                                                </Typography>
                                            </Box>
                                        )}
                                        {storefront.email && (
                                            <Box display="flex" alignItems="center" gap={1}>
                                                <Mail size={18} />
                                                <Typography>
                                                    <a href={`mailto:${storefront.email}`}>{storefront.email}</a>
                                                </Typography>
                                            </Box>
                                        )}
                                        {storefront.website && (
                                            <Box display="flex" alignItems="center" gap={1}>
                                                <Globe size={18} />
                                                <Typography>
                                                    <a href={storefront.website} target="_blank" rel="noopener noreferrer">
                                                        {storefront.website.replace(/^https?:\/\//, '')}
                                                    </a>
                                                </Typography>
                                            </Box>
                                        )}
                                        {storefront.address && (
                                            <Box display="flex" alignItems="center" gap={1}>
                                                <MapPin size={18} />
                                                <Typography>
                                                    {storefront.address}
                                                </Typography>
                                            </Box>
                                        )}
                                    </Stack>
                                </CardContent>
                            </Card>
                        </Grid>
                    )}

                    {hasLocation && (
                        <Grid item xs={12} md={hasContactInfo ? 6 : 12}>
                            <Card variant="outlined" sx={{ height: '100%' }}>
                                <CardContent>
                                    <Typography variant="h6" gutterBottom>
                                        {t('marketplace:store.location')}
                                    </Typography>
                                    <MiniMap
                                        latitude={storefront.latitude}
                                        longitude={storefront.longitude}
                                        address={storefront.address || `${storefront.city}, ${storefront.country}`}
                                    />
                                    {storefront.city && storefront.country && (
                                        <Chip
                                            icon={<MapPin size={14} />}
                                            label={`${storefront.city}, ${storefront.country}`}
                                            size="small"
                                            sx={{ mt: 1 }}
                                        />
                                    )}
                                </CardContent>
                            </Card>
                        </Grid>
                    )}
                </Grid>

                <Divider sx={{ my: 3 }} />

                <StorefrontRating storefrontId={id} />
                <Box display="flex" justifyContent="center" mt={2}>
                    <Button
                        component={Link}
                        to={`/shop/${id}/reviews`}
                        variant="outlined"
                    >
                        {t('store.seeAllReviews')}
                    </Button>
                </Box>
            </Paper>

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h5" component="h2">
                    {t('marketplace:store.storeProducts')}
                </Typography>
                
                {/* Добавляем переключатель режима отображения */}
                <ToggleButtonGroup
                    value={viewMode}
                    exclusive
                    onChange={handleViewModeChange}
                    aria-label="view mode"
                    size="small"
                >
                    <ToggleButton value="grid" aria-label="grid view">
                        <GridIcon size={18} />
                    </ToggleButton>
                    <ToggleButton value="list" aria-label="list view">
                        <List size={18} />
                    </ToggleButton>
                </ToggleButtonGroup>
            </Box>

            {storeListings.length === 0 && !highlightedListing ? (
                <Alert severity="info">
                    {t('marketplace:store.noProductsMessage')}
                </Alert>
            ) : (
                <>
                    {/* Блок для выделенного объявления, если оно есть */}
                    {highlightedListing && (
                        <Box sx={{ mb: 4 }}>
                            <Typography variant="h6" component="h3" sx={{ mb: 2 }}>
                                {t('marketplace:store.highlightedListing', { defaultValue: 'Выбранное объявление' })}
                            </Typography>
                            <Paper elevation={1} sx={{ p: 2, mb: 2, bgcolor: theme.palette.primary?.light + '22' }}>
                                {viewMode === 'grid' ? (
                                    <Grid container>
                                        <Grid item xs={12} sm={6} md={4}>
                                            <Link
                                                to={`/marketplace/listings/${highlightedListing.id}`}
                                                style={{ textDecoration: 'none' }}
                                            >
                                                <ListingCard listing={highlightedListing} />
                                            </Link>
                                        </Grid>
                                    </Grid>
                                ) : (
                                    <MarketplaceListingsList
                                        listings={[highlightedListing]}
                                        showSelection={false}
                                        initialSortField={sortField}
                                        initialSortOrder={sortDirection}
                                        onSortChange={handleSortChange}
                                    />
                                )}
                            </Paper>
                            
                            {storeListings.length > 0 && (
                                <Typography variant="h6" component="h3" sx={{ mt: 4, mb: 2 }}>
                                    {t('marketplace:store.otherListings', { defaultValue: 'Другие объявления' })}
                                </Typography>
                            )}
                        </Box>
                    )}
                    
                    {/* Блок для остальных объявлений */}
                    {storeListings.length > 0 && (
                        <InfiniteScroll
                            hasMore={hasMoreListings}
                            loading={loadingMore}
                            onLoadMore={handleLoadMore}
                            autoLoad={!isMobile} // Автозагрузка только на десктопе
                            loadingMessage={t('marketplace:listings.loading', { defaultValue: 'Загрузка...' })}
                            loadMoreButtonText={t('marketplace:listings.loadMore', { defaultValue: 'Показать ещё' })}
                            noMoreItemsText={t('marketplace:store.noMoreProducts', { defaultValue: 'Больше нет товаров' })}
                        >
                            {viewMode === 'grid' ? (
                                <Grid container spacing={3}>
                                    {storeListings.map((listing) => (
                                        <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                            <Link
                                                to={`/marketplace/listings/${listing.id}`}
                                                style={{ textDecoration: 'none' }}
                                            >
                                                <ListingCard listing={listing} />
                                            </Link>
                                        </Grid>
                                    ))}
                                </Grid>
                            ) : (
                                <MarketplaceListingsList
                                    listings={storeListings}
                                    showSelection={false}
                                    initialSortField={sortField}
                                    initialSortOrder={sortDirection}
                                    onSortChange={handleSortChange}
                                />
                            )}
                        </InfiniteScroll>
                    )}
                </>
            )}
        </Container>
    );
};

export default PublicStorefrontPage;
//frontend/hostel-frontend/src/pages/MarketplacePage.js
import { useEffect, useState, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import {
    Container,
    Grid,
    Box,
    Typography,
    CircularProgress,
    Button,
    useTheme,
    useMediaQuery,
    Drawer,
    IconButton,
    Fab,
    Alert,
    Paper,
    Chip,
} from '@mui/material';
import { Plus, Filter, X } from 'lucide-react';
import ListingCard from '../../components/marketplace/ListingCard';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import {
    MobileFilters,
    MobileHeader,
    MobileListingCard
} from '../../components/marketplace/MobileComponents';
import CompactMarketplaceFilters from '../../components/marketplace/MarketplaceFilters';
import axios from '../../api/axios';
import { debounce } from 'lodash';
import { useSearchParams } from 'react-router-dom';

const MarketplacePage = () => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const navigate = useNavigate();

    const [listings, setListings] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isFilterOpen, setIsFilterOpen] = useState(false);
    const [searchParams, setSearchParams] = useSearchParams();
    const [categoryPath, setCategoryPath] = useState([]);
    const [filters, setFilters] = useState({
        query: '',
        category_id: searchParams.get('category_id') || '', 
        min_price: '',
        max_price: '',
        city: '',
        country: '',
        condition: '',
        sort_by: 'date_desc'
    });
    




    const fetchListings = useCallback(async (currentFilters = {}) => {
        try {
            setLoading(true);
            setError(null);
    
            // Формируем параметры запроса только из непустых значений
            const params = {};
            Object.entries(currentFilters).forEach(([key, value]) => {
                if (value !== '') {
                    params[key] = value;
                }
            });
    
            const response = await axios.get('/api/v1/marketplace/listings', { 
                params
            });
    
            if (response.data?.data?.data) {
                setListings(response.data.data.data);
            } else {
                setListings([]);
            }
            
        } catch (err) {
            console.error('Error fetching listings:', err);
            setError('Не удалось загрузить объявления');
        } finally {
            setLoading(false);
        }
    }, []);
    

    // Эффект для загрузки категорий
    useEffect(() => {
        const fetchInitialData = async () => {
            try {
                // Загружаем категории
                const categoriesResponse = await axios.get('/api/v1/marketplace/category-tree');
                if (categoriesResponse.data?.data) {
                    setCategories(categoriesResponse.data.data);
                }
    
                // Получаем начальный categoryId из URL
                const categoryId = searchParams.get('category_id');
                
                // Загружаем листинги с учётом начальных фильтров
                await fetchListings({
                    ...filters,
                    category_id: categoryId || '',
                    sort_by: 'date_desc'
                });
                    
            } catch (err) {
                console.error('Error fetching initial data:', err);
                setError('Произошла ошибка при загрузке данных');
            }
        };
    
        fetchInitialData();
    }, [searchParams]); // Добавляем зависимость от searchParams
    const findCategoryPath = (categoryId, categoriesTree) => {
        const path = [];
        
        const findPath = (id, categories) => {
            for (const category of categories) {
                if (String(category.id) === String(id)) {
                    path.unshift({ id: category.id, name: category.name, slug: category.slug });
                    return true;
                }
                
                if (category.children && findPath(id, category.children)) {
                    path.unshift({ id: category.id, name: category.name, slug: category.slug });
                    return true;
                }
            }
            return false;
        };
        
        findPath(categoryId, categoriesTree);
        return path;
    };
    // Заменяем второй эффект
    useEffect(() => {
        // Проверяем путь и делаем редирект если нужно
        if (!window.location.pathname.includes('/marketplace')) {
            navigate({
                pathname: '/marketplace',
                search: window.location.search
            }, { replace: true }); // Используем replace чтобы не создавать новую запись в истории
        }
    }, []); // Срабатывает только при монтировании
    
    // Третий эффект оставляем как есть
    useEffect(() => {
        const categoryId = searchParams.get('category_id');
        
        if (categoryId !== filters.category_id) {
            setFilters(prev => ({
                ...prev,
                category_id: categoryId || ''
            }));
            
            fetchListings({
                ...filters,
                category_id: categoryId || ''
            });
        }
    }, [searchParams]);
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
            const updated = { ...prev, ...newFilters };
            
            // Обновляем URL при изменении категории
            if (newFilters.category_id !== undefined) {
                const nextParams = new URLSearchParams(searchParams);
                if (newFilters.category_id) {
                    nextParams.set('category_id', newFilters.category_id);
                    // Убедимся, что мы находимся на правильном пути
                    if (!window.location.pathname.includes('/marketplace')) {
                        navigate({
                            pathname: '/marketplace',
                            search: nextParams.toString()
                        });
                    } else {
                        setSearchParams(nextParams);
                    }
                } else {
                    nextParams.delete('category_id');
                    setSearchParams(nextParams);
                }
            }
    
            const cleanFilters = {};
            Object.entries(updated).forEach(([key, value]) => {
                if (value !== '') {
                    cleanFilters[key] = value;
                }
            });
            fetchListings(cleanFilters);
            return updated;
        });
    }, [searchParams, setSearchParams, navigate]);

    const getActiveFiltersCount = () => {
        return Object.entries(filters).reduce((count, [key, value]) => {
            if (key !== 'sort_by' && value !== '') {
                return count + 1;
            }
            return count;
        }, 0);
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
                <Alert
                    severity="error"
                    sx={{ m: 2 }}
                    action={
                        <IconButton size="small" onClick={() => setError(null)}>
                            <X size={16} />
                        </IconButton>
                    }
                >
                    {error}
                </Alert>
            );
        }

        if (listings.length === 0) {
            return (
                <Alert severity="info" sx={{ m: 2 }}>
                    По вашему запросу ничего не найдено
                </Alert>
            );
        }

        return (
            <Grid container spacing={isMobile ? 1 : 3}>
                {listings.map((listing) => (
                    <Grid item xs={isMobile ? 6 : 12} sm={6} md={4} key={listing.id}>
                        <Link
                            to={`/marketplace/listings/${listing.id}`}
                            style={{ textDecoration: 'none' }}
                        >
                            {isMobile ? (
                                <MobileListingCard listing={listing} />
                            ) : (
                                <ListingCard listing={listing} />
                            )}
                        </Link>
                    </Grid>
                ))}
            </Grid>
        );
    };

    if (isMobile) {
        return (
            <Box sx={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
                <MobileHeader
                    onOpenFilters={() => setIsFilterOpen(true)}
                    filtersCount={getActiveFiltersCount()}
                />
    
                {/* Breadcrumbs и кнопка в одной строке */}
                <Box
                    sx={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        px: 2,
                        py: 1,
                        bgcolor: 'grey.100'
                    }}
                >
                    {categoryPath.length > 0 && <Breadcrumbs paths={categoryPath} />}
                    <Button
                        variant="contained"
                        size="small"
                        onClick={() => navigate('/marketplace/create')}
                        startIcon={<Plus size={20} />}
                    >
                        Создать
                    </Button>
                </Box>
    
                <Box sx={{ flex: 1, p: 1, bgcolor: 'grey.50' }}>
                    {filters.category_id && (
                        <Box sx={{ px: 1, mb: 1 }}>
                            <Chip
                                label={categories.find(c => c.id === filters.category_id)?.name}
                                onDelete={() => handleFilterChange({ category_id: '' })}
                                size="small"
                            />
                        </Box>
                    )}
                    {renderContent()}
                </Box>
    
                <MobileFilters
                    open={isFilterOpen}
                    onClose={() => setIsFilterOpen(false)}
                    filters={filters}
                    onFilterChange={handleFilterChange}
                    categories={categories}
                />
            </Box>
        );
    }
    
    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            {/* Breadcrumbs и кнопка в одной строке */}
            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 4
                }}
            >
                {categoryPath.length > 0 ? (
                    <Breadcrumbs paths={categoryPath} />
                ) : (
                    // Пустое место, если хлебных крошек нет
                    <Box sx={{ flex: 1 }} />
                )}
                <Button
                    variant="contained"
                    onClick={() => navigate('/marketplace/create')}
                    startIcon={<Plus />}
                >
                    Создать объявление
                </Button>
            </Box>
    
            <Grid container spacing={3}>
                <Grid item xs={12} md={3}>
                    <CompactMarketplaceFilters
                        filters={filters}
                        onFilterChange={handleFilterChange}
                        categories={categories}
                        selectedCategoryId={filters.category_id}
                        isLoading={loading}
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
//frontend/hostel-frontend/src/pages/marketplace/MarketplacePage.js
import { useTranslation } from 'react-i18next';

import { useEffect, useState, useCallback } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
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
}
    from '@mui/material';
import { Plus, Search as SearchIcon, X, Store } from 'lucide-react';
import ListingCard from '../../components/marketplace/ListingCard';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import {
    MobileFilters,
    MobileListingCard,
    MobileHeader,
} from '../../components/marketplace/MobileComponents';
import CompactMarketplaceFilters from '../../components/marketplace/MarketplaceFilters';
import axios from '../../api/axios';



const MobileListingGrid = ({ listings }) => {
    const navigate = useNavigate();

    return (
        <Box sx={{ px: 1 }}>
            <Grid container spacing={1}>
                {listings.map((listing) => (
                    <Grid item xs={6} key={listing.id}>
                        <Box
                            component={Paper}
                            variant="outlined"
                            sx={{
                                height: '100%',
                                overflow: 'hidden',
                                transition: 'transform 0.2s, box-shadow 0.2s',
                                '&:active': {
                                    transform: 'scale(0.98)'
                                }
                            }}
                        >
                            <Link
                                to={`/marketplace/listings/${listing.id}`}
                                style={{ textDecoration: 'none', color: 'inherit' }}
                                onClick={(e) => {
                                    const shopButton = e.target.closest('[data-shop-button="true"]');
                                    if (shopButton) {
                                        e.preventDefault();
                                        return;
                                    }
                                }}
                            >
                                <MobileListingCard listing={listing} />
                            </Link>
                        </Box>
                    </Grid>
                ))}
            </Grid>
        </Box>
    );
};

const MarketplacePage = () => {
    const { t } = useTranslation('marketplace');
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
    const [filters, setFilters] = useState({
        query: searchParams.get('query') || '',
        category_id: searchParams.get('category_id') || '',
        min_price: searchParams.get('min_price') || '',
        max_price: searchParams.get('max_price') || '',
        city: searchParams.get('city') || '',
        country: searchParams.get('country') || '',
        condition: searchParams.get('condition') || '',
        sort_by: searchParams.get('sort_by') || 'date_desc'
    });

    const fetchListings = useCallback(async (currentFilters = {}) => {
        try {
          setLoading(true);
          setError(null);
      
          const params = {};
          Object.entries(currentFilters).forEach(([key, value]) => {
            if (value !== '') {
              params[key] = value;
            }
          });
      
          // Используем новый endpoint для поиска
          const response = await axios.get('/api/v1/marketplace/search', { params });
          
          console.log('API response:', response.data);
          
          // Проверяем различные варианты структуры данных
          if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
            // Случай с двойным вложением data
            setListings(response.data.data.data);
          } else if (response.data?.data && Array.isArray(response.data.data)) {
            // Случай с одинарным вложением data
            setListings(response.data.data);
          } else {
            console.error('Ответ API не содержит массив данных:', response.data);
            setListings([]);
          }
        } catch (err) {
          console.error('Error fetching listings:', err);
          setError('Не удалось загрузить объявления');
          setListings([]);
        } finally {
          setLoading(false);
        }
      }, []);

    useEffect(() => {
        const fetchInitialData = async () => {
            try {
                const categoriesResponse = await axios.get('/api/v1/marketplace/category-tree');
                console.log('Fetched categories:', categoriesResponse.data?.data); // Добавьте этот лог

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
                    sort_by: searchParams.get('sort_by') || 'date_desc'
                };

                setFilters(initialFilters);
                await fetchListings(initialFilters);
            } catch (err) {
                console.error('Error fetching initial data:', err);
                setError('Произошла ошибка при загрузке данных');
            }
        };

        fetchInitialData();
    }, [searchParams]);

    useEffect(() => {
        if (!window.location.pathname.includes('/marketplace')) {
            navigate({
                pathname: '/marketplace',
                search: window.location.search
            }, { replace: true });
        }
    }, []);

    const findCategoryPath = (categoryId, categoriesTree) => {
        if (!categoryId || !categoriesTree || categoriesTree.length === 0) {
            return [];
        }

        // Создаем плоскую карту всех категорий для быстрого поиска
        const categoryMap = new Map();

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
            const updated = { ...prev, ...newFilters };

            const nextParams = new URLSearchParams(searchParams);
            Object.entries(updated).forEach(([key, value]) => {
                if (value) {
                    nextParams.set(key, value);
                } else {
                    nextParams.delete(key);
                }
            });

            if (!window.location.pathname.includes('/marketplace')) {
                navigate({
                    pathname: '/marketplace',
                    search: nextParams.toString()
                });
            } else {
                setSearchParams(nextParams);
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
    }, [searchParams, setSearchParams, navigate, fetchListings]);

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
    
        // Проверяем, что listings - это массив
        if (!listings || !Array.isArray(listings) || listings.length === 0) {
            return (
                <Alert severity="info" sx={{ m: 2 }}>
                    {t('listings.filters.noResults')}
                </Alert>
            );
        }
    
        return isMobile ? (
            <MobileListingGrid listings={listings} />
        ) : (
            <Grid container spacing={3}>
                {listings.map((listing) => (
                    <Grid item xs={12} sm={6} md={4} key={listing.id}>
                        <Link
                            to={`/marketplace/listings/${listing.id}`}
                            style={{ textDecoration: 'none' }}
                            onClick={(e) => {
                                if (listing.storefront_id && e.target.closest('[data-shop-button="true"]')) {
                                    e.preventDefault();
                                    navigate(`/shop/${listing.storefront_id}`);
                                }
                            }}
                        >
                            <ListingCard listing={listing} />
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
                        <Box sx={{ px: 2, py: 1, display: 'flex', gap: 1, overflowX: 'auto' }}>
                            {Object.entries(filters).map(([key, value]) => {
                                if (!value || key === 'sort_by') return null;
                                let label = value;
                                if (key === 'category_id') {
                                    const category = categories.find(c => String(c.id) === String(value));
                                    label = category ? category.name : value;
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
                        </Box>
                    )}
                </Box>

                <Box sx={{ flex: 1, bgcolor: 'grey.50' }}>
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
            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 4
                }}
            >
                {categoryPath.length > 0 ? (
                    <Breadcrumbs paths={categoryPath} categories={categories} />
                ) : (
                    <Box sx={{ flex: 1 }} />
                )}
                <Button
                    id="createAnnouncementButton"
                    variant="contained"
                    onClick={() => navigate('/marketplace/create')}
                    startIcon={<Plus />}
                >

                    {t('listings.create.title')}


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
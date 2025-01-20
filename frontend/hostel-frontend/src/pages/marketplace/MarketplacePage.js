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
    IconButton,
    Alert,
    Paper,
    Chip,
    InputBase,
    Toolbar,
} from '@mui/material';
import { Plus, Search as SearchIcon, Filter, X } from 'lucide-react';
import ListingCard from '../../components/marketplace/ListingCard';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import {
    MobileFilters,
    MobileListingCard,
    MobileHeader,
} from '../../components/marketplace/MobileComponents';
import CompactMarketplaceFilters from '../../components/marketplace/MarketplaceFilters';
import axios from '../../api/axios';
import { debounce } from 'lodash';
import { useSearchParams } from 'react-router-dom';


const MobileListingGrid = ({ listings }) => (
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
                        >
                            <MobileListingCard listing={listing} />
                        </Link>
                    </Box>
                </Grid>
            ))}
        </Grid>
    </Box>
);

const MarketplacePage = () => {
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

            const response = await axios.get('/api/v1/marketplace/listings', { params });

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

        if (listings.length === 0) {
            return (
                <Alert severity="info" sx={{ m: 2 }}>
                    По вашему запросу ничего не найдено
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
                    <Box sx={{ 
                        px: 2, 
                        py: 1,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        overflowX: 'auto',
                        WebkitOverflowScrolling: 'touch'
                    }}>
                        <Box sx={{ flex: 1, mr: 2 }}>
                            {categoryPath.length > 0 ? (
                                <Breadcrumbs paths={categoryPath} />
                            ) : (
                                <Typography variant="body2" color="text.secondary">
                                    Все категории
                                </Typography>
                            )}
                        </Box>
                        <Box sx={{ flexShrink: 0 }}>
                            <Button
                                variant="contained"
                                size="small"
                                onClick={() => navigate('/marketplace/create')}
                                startIcon={<Plus size={16} />}
                            >
                                Создать
                            </Button>
                        </Box>
                    </Box>
     
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
                    <Breadcrumbs paths={categoryPath} />
                ) : (
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
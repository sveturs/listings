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
import ListingCard from '../components/marketplace/ListingCard';
import { 
    MobileFilters, 
    MobileHeader, 
    MobileListingCard 
  } from '../components/marketplace/mobile/MobileComponents';
  import CompactMarketplaceFilters from '../components/marketplace/MarketplaceFilters';
  import axios from '../api/axios';
import { debounce } from 'lodash';

const MarketplacePage = () => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const navigate = useNavigate();

    const [listings, setListings] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [isFilterOpen, setIsFilterOpen] = useState(false);
    const [filters, setFilters] = useState({
        query: '',
        category_id: '',
        min_price: '',
        max_price: '',
        city: '',
        country: '',
        condition: '',
        sort_by: 'date_desc'
    });

    // Загрузка категорий
    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const response = await axios.get('/api/v1/marketplace/category-tree');
                setCategories(response.data.data || []);
            } catch (err) {
                console.error('Error fetching categories:', err);
                setError('Не удалось загрузить категории');
            }
        };
        fetchCategories();
    }, []);

    // Загрузка объявлений
    const fetchListings = useCallback(async (currentFilters) => {
        try {
            setLoading(true);
            setError(null);

            const params = Object.entries(currentFilters).reduce((acc, [key, value]) => {
                if (value !== '' && value !== null && value !== undefined) {
                    acc[key] = value;
                }
                return acc;
            }, {});

            const response = await axios.get('/api/v1/marketplace/listings', { params });
            const listingsData = response.data?.data?.data || [];
            setListings(listingsData);
        } catch (error) {
            console.error('Error fetching listings:', error);
            setError('Не удалось загрузить объявления');
        } finally {
            setLoading(false);
        }
    }, []);

    const handleFilterChange = useCallback((newFilters) => {
        setFilters(prev => ({
            ...prev,
            ...newFilters
        }));
    }, []);

    useEffect(() => {
        const debouncedFetch = debounce(() => fetchListings(filters), 500);
        debouncedFetch();
        return () => debouncedFetch.cancel();
    }, [fetchListings, filters]);
    const getActiveFiltersCount = () => {
        return Object.entries(filters).reduce((count, [key, value]) => {
            if (key !== 'sort_by' && value !== '') {
                return count + 1;
            }
            return count;
        }, 0);
    };
    if (isMobile) {
        return (
            <Box sx={{ 
                minHeight: '100vh',
                display: 'flex',
                flexDirection: 'column'
            }}>
                <MobileHeader
                    onOpenFilters={() => setIsFilterOpen(true)}
                    filtersCount={getActiveFiltersCount()}
                />

                {error && (
                    <Alert 
                        severity="error" 
                        sx={{ mx: 2, my: 1 }}
                        action={
                            <IconButton
                                size="small"
                                onClick={() => setError(null)}
                            >
                                <X size={16} />
                            </IconButton>
                        }
                    >
                        {error}
                    </Alert>
                )}

                <Box sx={{ 
                    flex: 1,
                    p: 1,
                    bgcolor: 'grey.50'
                }}>
                    {loading ? (
                        <Box display="flex" justifyContent="center" p={4}>
                            <CircularProgress />
                        </Box>
                    ) : (
                        <>
                            {filters.category_id && (
                                <Box sx={{ px: 1, mb: 1 }}>
                                    <Chip
                                        label={categories.find(c => c.id === filters.category_id)?.name}
                                        onDelete={() => handleFilterChange({ category_id: '' })}
                                        size="small"
                                    />
                                </Box>
                            )}
                            
                            <Grid container spacing={1}>
                                {listings.map((listing) => (
                                    <Grid item xs={6} key={listing.id}>
                                        <Link
                                            to={`/marketplace/listings/${listing.id}`}
                                            style={{ textDecoration: 'none' }}
                                        >
                                            <MobileListingCard listing={listing} />
                                        </Link>
                                    </Grid>
                                ))}
                                {listings.length === 0 && !loading && (
                                    <Grid item xs={12}>
                                        <Box 
                                            sx={{ 
                                                textAlign: 'center',
                                                py: 4,
                                                color: 'text.secondary'
                                            }}
                                        >
                                            <Typography variant="body2">
                                                По вашему запросу ничего не найдено
                                            </Typography>
                                        </Box>
                                    </Grid>
                                )}
                            </Grid>
                        </>
                    )}
                </Box>

                <MobileFilters
                    open={isFilterOpen}
                    onClose={() => setIsFilterOpen(false)}
                    filters={filters}
                    onFilterChange={handleFilterChange}
                    categories={categories}
                />

                <Box
                    component={Paper}
                    elevation={3}
                    sx={{
                        position: 'sticky',
                        bottom: 0,
                        p: 2,
                        borderRadius: 0
                    }}
                >
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={() => navigate('/marketplace/create')}
                        startIcon={<Plus size={20} />}
                    >
                        Разместить объявление
                    </Button>
                </Box>
            </Box>
        );
    }

// Десктопная версия
return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h4">
                Объявления
            </Typography>
            <Button
                variant="contained"
                onClick={() => navigate('/marketplace/create')}
                startIcon={<Plus />}
            >
                Создать объявление
            </Button>
        </Box>

        {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
                {error}
            </Alert>
        )}

        {/* Оборачиваем всё в контейнер Grid */}
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
                {loading ? (
                    <Box display="flex" justifyContent="center" p={4}>
                        <CircularProgress />
                    </Box>
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
                        {listings.length === 0 && !loading && (
                            <Grid item xs={12}>
                                <Alert severity="info">
                                    По вашему запросу ничего не найдено
                                </Alert>
                            </Grid>
                        )}
                    </Grid>
                )}
            </Grid>
        </Grid>
    </Container>
);
};

export default MarketplacePage;
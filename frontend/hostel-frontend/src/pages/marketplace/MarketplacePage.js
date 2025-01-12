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
import {
    MobileFilters,
    MobileHeader,
    MobileListingCard
} from '../../components/marketplace/MobileComponents';
import CompactMarketplaceFilters from '../../components/marketplace/MarketplaceFilters';
import axios from '../../api/axios';
import { debounce } from 'lodash';

const MarketplacePage = () => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const navigate = useNavigate();

    const [listings, setListings] = useState([]); 
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
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


    const handleFilterChange = useCallback((newFilters) => {
        setFilters(prev => ({
            ...prev,
            ...newFilters
        }));
    }, []);

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

    // Эффект для загрузки категорий
    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                // Параллельная загрузка категорий и листингов
                const [categoriesResponse, listingsResponse] = await Promise.all([
                    axios.get('/api/v1/marketplace/category-tree'),
                    axios.get('/api/v1/marketplace/listings', {
                        params: Object.entries(filters).reduce((acc, [key, value]) => {
                            if (value !== '' && value !== null && value !== undefined) {
                                acc[key] = value;
                            }
                            return acc;
                        }, {})
                    })
                ]);

                if (categoriesResponse.data?.data) {
                    setCategories(categoriesResponse.data.data);
                }
                setListings(listingsResponse.data?.data?.data || []);
                setError(null);
            } catch (err) {
                console.error('Error fetching data:', err);
                setError('Произошла ошибка при загрузке данных');
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [filters]); // Убираем лишний useEffect и объединяем загрузку данных

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
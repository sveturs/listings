import React, { useState, useEffect, useCallback } from 'react';
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
    Fab
} from '@mui/material';
import { Plus, Filter, X } from 'lucide-react';
import ListingCard from '../components/marketplace/ListingCard';
import MarketplaceFilters from '../components/marketplace/MarketplaceFilters';
import axios from '../api/axios';
import { debounce } from 'lodash';

const MarketplacePage = () => {
    const theme = useTheme();
    // Изменяем точку перелома на md
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const navigate = useNavigate();
    const [listings, setListings] = useState([]);
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
        with_photos: false,
        sort_by: 'date_desc'
    });

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
            console.log('API Response:', response.data);

            // Проверяем структуру ответа и извлекаем данные
            const listingsData = response.data?.data?.data || [];
            if (Array.isArray(listingsData)) {
                setListings(listingsData);
            } else {
                console.error('Invalid listings data format:', listingsData);
                setListings([]);
            }
        } catch (error) {
            console.error('Error fetching listings:', error);
            setError('Не удалось загрузить объявления');
            setListings([]);
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
    if (isMobile) {
        return (
            <Box sx={{ pb: 7, position: 'relative', minHeight: '100vh' }}>
                <Box sx={{
                    p: 2,
                    position: 'sticky',
                    top: 0,
                    bgcolor: 'background.paper',
                    borderBottom: '1px solid',
                    borderColor: 'divider',
                    zIndex: 10,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between'
                }}>
                    <Typography variant="h6">
                        Объявления
                    </Typography>
                    <IconButton onClick={() => setIsFilterOpen(true)}>
                        <Filter />
                    </IconButton>
                </Box>

                <Box sx={{ p: 1 }}>
                    {loading ? (
                        <Box display="flex" justifyContent="center" p={4}>
                            <CircularProgress />
                        </Box>
                    ) : error ? (
                        <Box display="flex" justifyContent="center" p={4}>
                            <Typography color="error">{error}</Typography>
                        </Box>
                    ) : (
                        <Grid container spacing={1}>
                            {listings.map((listing) => (
                                <Grid item xs={4} key={listing.id}>
                                    <Link
                                        to={`/marketplace/listings/${listing.id}`}
                                        style={{ textDecoration: 'none' }}
                                        onClick={(e) => {
                                            // Stop event propagation if needed
                                            e.stopPropagation();
                                            navigate(`/marketplace/listings/${listing.id}`);
                                        }}
                                    >
                                        <ListingCard listing={listing} isMobile={true} />
                                    </Link>
                                </Grid>
                            ))}
                        </Grid>
                    )}
                </Box>

                <Fab
                    color="primary"
                    sx={{
                        position: 'fixed',
                        bottom: 16,
                        right: 16,
                    }}
                    onClick={() => navigate('/marketplace/create')}
                >
                    <Plus />
                </Fab>

                <Drawer
                    anchor="bottom"
                    open={isFilterOpen}
                    onClose={() => setIsFilterOpen(false)}
                    PaperProps={{
                        sx: {
                            maxHeight: '90vh',
                            borderTopLeftRadius: 16,
                            borderTopRightRadius: 16,
                            pb: 4
                        }
                    }}
                >
                    <Box sx={{ p: 2 }}>
                        <Box sx={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            mb: 2
                        }}>
                            <Typography variant="h6">Фильтры</Typography>
                            <IconButton onClick={() => setIsFilterOpen(false)}>
                                <X />
                            </IconButton>
                        </Box>
                        <MarketplaceFilters
                            filters={filters}
                            onFilterChange={(newFilters) => {
                                handleFilterChange(newFilters);
                                setIsFilterOpen(false);
                            }}
                            isMobile={true}
                        />
                    </Box>
                </Drawer>
            </Box>
        );
    }
    return (
        <Container maxWidth="lg" sx={{ mt: 4 }}>
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
                    <MarketplaceFilters
                        filters={filters}
                        onFilterChange={handleFilterChange}
                    />
                </Grid>
                <Grid item xs={12} md={9}>
                    {loading ? (
                        <Box display="flex" justifyContent="center" p={4}>
                            <CircularProgress />
                        </Box>
                    ) : error ? (
                        <Box display="flex" justifyContent="center" p={4}>
                            <Typography color="error">{error}</Typography>
                        </Box>
                    ) : (
                        <Grid container spacing={2}>

                            {listings.map((listing) => (
                                <Grid item xs={4} key={listing.id}>
                                    <Link
                                        to={`/marketplace/listings/${listing.id}`}
                                        style={{ textDecoration: 'none' }}
                                    >
                                        <ListingCard listing={listing} isMobile={true} />
                                    </Link>
                                </Grid>
                            ))}
                        </Grid>
                    )}
                </Grid>
            </Grid>
        </Container>
    );
};

export default MarketplacePage;
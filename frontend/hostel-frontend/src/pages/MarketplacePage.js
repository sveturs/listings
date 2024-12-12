import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Container,
    Grid,
    Box,
    Typography,
    CircularProgress,
    Button,
    useTheme,
    useMediaQuery
} from '@mui/material';
import { Plus } from 'lucide-react';
import ListingCard from '../components/marketplace/ListingCard';
import MarketplaceFilters from '../components/marketplace/MarketplaceFilters';
import axios from '../api/axios';
import { debounce } from 'lodash';

const MarketplacePage = () => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    const navigate = useNavigate();
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
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
                    ) : listings.length > 0 ? (
                        <Grid container spacing={2}>
                            {listings.map((listing) => (
                                <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                    <ListingCard listing={listing} />
                                </Grid>
                            ))}
                        </Grid>
                    ) : (
                        <Box display="flex" justifyContent="center" p={4}>
                            <Typography>Объявления не найдены</Typography>
                        </Box>
                    )}
                </Grid>
            </Grid>
        </Container>
    );
};

export default MarketplacePage;
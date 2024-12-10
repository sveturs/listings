import React, { useState, useEffect, useCallback, useMemo } from 'react';
import {
    Container,
    Grid,
    Box,
    Typography,
    CircularProgress
} from '@mui/material';
import ListingCard from '../components/marketplace/ListingCard';
import MarketplaceFilters from '../components/marketplace/MarketplaceFilters';
import axios from '../api/axios';
import { debounce } from 'lodash';

const MarketplacePage = () => {
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(false);
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
            const params = Object.entries(currentFilters).reduce((acc, [key, value]) => {
                if (value !== '' && value !== null && value !== undefined) {
                    acc[key] = value;
                }
                return acc;
            }, {});

            const response = await axios.get('/api/v1/marketplace/listings', { params });
            console.log('Server response:', response.data);
            
            // Извлекаем данные из правильного уровня вложенности
            const listingsData = response.data.data.data || [];
            console.log('Extracted listings:', listingsData);
            
            setListings(listingsData);

            // Обновляем максимальную цену только при первой загрузке
            if (listingsData.length > 0 && !currentFilters.max_price) {
                const maxPrice = Math.max(...listingsData.map(listing => listing.price));
                setFilters(prev => ({
                    ...prev,
                    max_price: maxPrice
                }));
            }
        } catch (error) {
            console.error('Error fetching listings:', error);
            setListings([]);
        } finally {
            setLoading(false);
        }
    }, []);

    // Создаем дебаунсированную версию функции fetchListings
    const debouncedFetch = useMemo(
        () => debounce((filters) => fetchListings(filters), 500),
        [fetchListings]
    );

    // Обработчик изменения фильтров с дебаунсингом
    const handleFilterChange = useCallback((newFilters) => {
        setFilters(prevFilters => {
            const updatedFilters = {
                ...prevFilters,
                ...newFilters
            };
            debouncedFetch(updatedFilters);
            return updatedFilters;
        });
    }, [debouncedFetch]);

    useEffect(() => {
        fetchListings(filters);
        return () => debouncedFetch.cancel();
    }, []);

    return (
        <Container maxWidth="lg" sx={{ mt: 4 }}>
            <Box sx={{ mb: 4 }}>
                <Typography variant="h4" gutterBottom>
                    Объявления
                </Typography>
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
                    ) : (
                        <Grid container spacing={2}>
                            {listings && listings.length > 0 ? (
                                listings.map((listing) => (
                                    <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                        <ListingCard listing={listing} />
                                    </Grid>
                                ))
                            ) : (
                                <Grid item xs={12}>
                                    <Typography variant="body1" align="center">
                                        Объявления не найдены
                                    </Typography>
                                </Grid>
                            )}
                        </Grid>
                    )}
                </Grid>
            </Grid>
        </Container>
    );
};

export default React.memo(MarketplacePage);
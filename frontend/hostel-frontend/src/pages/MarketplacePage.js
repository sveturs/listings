import React, { useState, useEffect } from 'react';
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
        sort_by: 'date_desc'
    });

    const fetchListings = async () => {
        try {
            setLoading(true);
            const response = await axios.get('/api/v1/marketplace/listings', { params: filters });
            console.log('Response data:', response.data);
            setListings(response.data.data); // Предполагается, что массив находится в `data.data`
        } catch (error) {
            console.error('Error fetching listings:', error);
            setListings([]); // Устанавливаем пустой массив в случае ошибки
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchListings();
    }, [filters]);

    const handleFilterChange = (newFilters) => {
        setFilters(prev => ({ ...prev, ...newFilters }));
    };

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
                            {Array.isArray(listings) && listings.length > 0 ? (
                                listings.map((listing) => (
                                    <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                        <ListingCard listing={listing} />
                                    </Grid>
                                ))
                            ) : (
                                <Typography variant="body1">No listings available</Typography>
                            )}
                        </Grid>
                    )}
                </Grid>
            </Grid>
        </Container>
    );
};

export default MarketplacePage;
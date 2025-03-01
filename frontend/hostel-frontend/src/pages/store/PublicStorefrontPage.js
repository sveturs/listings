// frontend/hostel-frontend/src/pages/store/PublicStorefrontPage.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, Link } from 'react-router-dom';
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
    Button
} from '@mui/material';
import { Store } from 'lucide-react';
import axios from '../../api/axios';
import ListingCard from '../../components/marketplace/ListingCard';

const PublicStorefrontPage = () => {
    const { t } = useTranslation(['common', 'marketplace']);
    const { id } = useParams();

    const [loading, setLoading] = useState(true);
    const [storefront, setStorefront] = useState(null);
    const [storeListings, setStoreListings] = useState([]);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                const [storefrontResponse, listingsResponse] = await Promise.all([
                    axios.get(`/api/v1/public/storefronts/${id}`),
                    axios.get('/api/v1/marketplace/listings', {
                        params: { storefront_id: id }
                    })
                ]);

                setStorefront(storefrontResponse.data.data);
                setStoreListings(listingsResponse.data.data?.data || []);
            } catch (err) {
                console.error('Error fetching storefront data:', err);
                setError(t('marketplace:listings.errors.loadFailed'));
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [id, t]);

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
                <StorefrontRating storefrontId={id} />
                {/* Кнопка для просмотра всех отзывов */}
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

            <Typography variant="h5" component="h2" gutterBottom sx={{ mb: 3 }}>
                {t('marketplace:store.storeProducts')}
            </Typography>

            {storeListings.length === 0 ? (
                <Alert severity="info">
                    {t('marketplace:store.noProductsMessage')}
                </Alert>
            ) : (
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
            )}
        </Container>
    );
};

export default PublicStorefrontPage;
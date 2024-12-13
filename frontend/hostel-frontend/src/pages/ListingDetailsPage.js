import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
    Container,
    Grid,
    Box,
    Typography,
    Tabs,
    Tab,
    CircularProgress,
    Rating,
    Stack
} from '@mui/material';
import ReviewsSection from '../components/reviews/ReviewsSection';
import axios from '../api/axios';

const ListingDetailsPage = () => {
    const { id } = useParams();
    const [listing, setListing] = useState(null);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState(0);

    // Определяем fetchListing внутри useEffect с использованием useCallback
    useEffect(() => {
        const fetchListing = async () => {
            if (!id) return;
            
            try {
                const response = await axios.get(`/api/v1/marketplace/listings/${id}`);
                setListing(response.data.data);
            } catch (error) {
                console.error('Error fetching listing:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchListing();
    }, [id]);

    if (loading) {
        return (
            <Container sx={{ py: 4, textAlign: 'center' }}>
                <CircularProgress />
            </Container>
        );
    }

    if (!listing) {
        return (
            <Container sx={{ py: 4, textAlign: 'center' }}>
                <Typography>Объявление не найдено</Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Grid container spacing={4}>
                <Grid item xs={12} md={7}>
                    {/* Здесь компонент галереи */}
                </Grid>
                <Grid item xs={12} md={5}>
                    <Box>
                        <Typography variant="h4" gutterBottom>
                            {listing.title}
                        </Typography>
                        {listing.rating > 0 && (
                            <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 2 }}>
                                <Rating value={listing.rating} readOnly precision={0.1} />
                                <Typography variant="body2" color="text.secondary">
                                    {listing.rating?.toFixed(1)} ({listing.reviews_count || 0} отзывов)
                                </Typography>
                            </Stack>
                        )}
                    </Box>
                </Grid>
            </Grid>

            <Box sx={{ mt: 4 }}>
                <Tabs
                    value={activeTab}
                    onChange={(e, newValue) => setActiveTab(newValue)}
                    sx={{ borderBottom: 1, borderColor: 'divider' }}
                >
                    <Tab label="Описание" />
                    <Tab
                        label={`Отзывы (${listing.reviews_count || 0})`}
                        id="reviews-tab"
                    />
                </Tabs>

                {activeTab === 0 && (
                    <Box sx={{ py: 3 }}>
                        <Typography>{listing.description}</Typography>
                    </Box>
                )}

                {activeTab === 1 && (
                    <Box sx={{ py: 3 }}>
                        <ReviewsSection
                            entityType="listing"
                            entityId={listing.id}
                            entityTitle={listing.title}
                            canAddReview={true}
                        />
                    </Box>
                )}
            </Box>
        </Container>
    );
};

export default ListingDetailsPage;
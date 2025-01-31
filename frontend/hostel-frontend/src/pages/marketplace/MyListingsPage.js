// frontend/hostel-frontend/src/pages/marketplace/MyListingsPage.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

import {
    Container,
    Typography,
    Grid,
    Box,
    CircularProgress,
    Alert,
    Button
} from '@mui/material';
import { Plus } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import ListingCard from '../../components/marketplace/ListingCard';
import axios from '../../api/axios';

const MyListingsPage = () => {
        const { t } = useTranslation('marketplace');
    
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user } = useAuth();

    useEffect(() => {
        const fetchMyListings = async () => {
            try {
                setLoading(true);
                const response = await axios.get('/api/v1/marketplace/listings', {
                    withCredentials: true // Added this to ensure authentication
                });

                console.log('Listings API Response:', response);

                if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                    // Filter listings by user ID
                    const userListings = response.data.data.data.filter(listing =>
                        String(listing.user_id) === String(user?.id)
                    );
                    setListings(userListings);
                } else {
                    console.log('Unexpected listings data structure:', response.data);
                    setListings([]);
                }
            } catch (err) {
                console.error('Error fetching listings:', err);
                setError('Не удалось загрузить объявления');
            } finally {
                setLoading(false);
            }
        };

        if (user?.id) {
            fetchMyListings();
        } else {
            setLoading(false); // Останавливаем загрузку, если нет пользователя
        }
    }, [user]);

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" p={4}>
                <CircularProgress />
            </Box>
        );
    }

    if (error) {
        return (
            <Container>
                <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>
            </Container>
        );
    }

    return (
        <Container sx={{ py: 4 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Typography variant="h4" component="h1">
                {t('listings.myListings')}
                </Typography>
                <Button
                    id="createAnnouncementButton" // Добавлено id
                    component={Link}
                    to="/marketplace/create"
                    variant="contained"
                    startIcon={<Plus />}
                >
                    {t('listings.create.title')}



                </Button>
            </Box>

            <Grid container spacing={3}>
                {listings.length === 0 ? (
                    <Grid item xs={12}>
                        <Alert severity="info">
                        {t('listings.Youdonthave')}
                        </Alert>
                    </Grid>
                ) : (
                    listings.map((listing) => (
                        <Grid item xs={12} sm={6} md={4} key={listing.id}>
                            <Link
                                to={`/marketplace/listings/${listing.id}`}
                                style={{ textDecoration: 'none' }}
                            >
                                <ListingCard listing={listing} />
                            </Link>
                        </Grid>
                    ))
                )}
            </Grid>
        </Container>
    );
};

export default MyListingsPage;
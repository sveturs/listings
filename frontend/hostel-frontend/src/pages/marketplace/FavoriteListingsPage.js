import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import {
    Container,
    Typography,
    Grid,
    Box,
    CircularProgress,
    Alert,
} from '@mui/material';
import ListingCard from '../../components/marketplace/ListingCard';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';

const FavoriteListingsPage = () => {
    const [listings, setListings] = useState(null); // Изменили на null
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user } = useAuth();

    useEffect(() => {
        const fetchFavorites = async () => {
            try {
                setLoading(true);
                const response = await axios.get('/api/v1/marketplace/favorites');
                setListings(response.data.data);
            } catch (err) {
                console.error('Error fetching favorites:', err);
                setError('Не удалось загрузить избранные объявления');
                setListings([]); // Устанавливаем пустой массив в случае ошибки
            } finally {
                setLoading(false);
            }
        };

        if (user) {
            fetchFavorites();
        } else {
            setListings([]); // Если пользователь не авторизован, устанавливаем пустой массив
            setLoading(false);
        }
    }, [user]);

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" p={4}>
                <CircularProgress />
            </Box>
        );
    }

    if (!user) {
        return (
            <Container sx={{ py: 4 }}>
                <Alert severity="info">
                    Для просмотра избранных объявлений необходимо авторизоваться
                </Alert>
            </Container>
        );
    }

    return (
        <Container sx={{ py: 4 }}>
            <Typography variant="h4" gutterBottom>
                Избранные объявления
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            {listings && listings.length === 0 ? (
                <Alert severity="info">
                    У вас пока нет избранных объявлений
                </Alert>
            ) : (
                <Grid container spacing={3}>
                    {listings && listings.map((listing) => (
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

export default FavoriteListingsPage;
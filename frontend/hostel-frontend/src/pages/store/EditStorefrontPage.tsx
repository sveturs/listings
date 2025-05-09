// frontend/hostel-frontend/src/pages/store/EditStorefrontPage.tsx
import React, { useState, useEffect, ChangeEvent, FormEvent } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import {
    Container,
    Typography,
    TextField,
    Button,
    Paper,
    Box,
    Grid,
    CircularProgress,
    Alert,
    Divider
} from '@mui/material';
import axios from '../../api/axios';
import LocationPicker from '../../components/global/LocationPicker';

interface StorefrontData {
    name: string;
    description: string;
    slug: string;
    phone: string;
    email: string;
    website: string;
    address: string;
    city: string;
    country: string;
    latitude: number | null;
    longitude: number | null;
    [key: string]: string | number | null;
}

interface LocationData {
    latitude: number | null;
    longitude: number | null;
    formatted_address: string;
    address_components?: {
        city: string;
        country: string;
    };
}

type RouteParams = {
    id: string;
}

const EditStorefrontPage: React.FC = () => {
    const { t } = useTranslation(['common', 'marketplace']);
    const { id } = useParams<keyof RouteParams>();
    const navigate = useNavigate();

    const [loading, setLoading] = useState<boolean>(true);
    const [saving, setSaving] = useState<boolean>(false);
    const [storefront, setStorefront] = useState<StorefrontData>({
        name: '',
        description: '',
        slug: '',
        phone: '',
        email: '',
        website: '',
        address: '',
        city: '',
        country: '',
        latitude: null,
        longitude: null
    });
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<boolean>(false);

    useEffect(() => {
        const fetchStorefront = async (): Promise<void> => {
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/storefronts/${id}`);

                if (response.data?.data) {
                    setStorefront(response.data.data);
                }
            } catch (err) {
                console.error('Error fetching storefront:', err);
                setError(t('marketplace:store.errors.loadFailed'));
            } finally {
                setLoading(false);
            }
        };

        if (id) {
            fetchStorefront();
        }
    }, [id, t]);

    const handleSubmit = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
        e.preventDefault();

        try {
            setSaving(true);

            await axios.put(`/api/v1/storefronts/${id}`, storefront);

            setSuccess(true);
            setTimeout(() => {
                navigate(`/storefronts/${id}`);
            }, 1500);
        } catch (err) {
            console.error('Error updating storefront:', err);
            setError(t('marketplace:store.errors.updateFailed'));
        } finally {
            setSaving(false);
        }
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>): void => {
        const { name, value } = e.target;
        setStorefront(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleLocationSelect = (location: LocationData): void => {
        setStorefront(prev => ({
            ...prev,
            address: location.formatted_address,
            city: location.address_components?.city || '',
            country: location.address_components?.country || '',
            latitude: location.latitude,
            longitude: location.longitude
        }));
    };

    if (loading) {
        return (
            <Container maxWidth="md" sx={{ py: 4, textAlign: 'center' }}>
                <CircularProgress />
                <Typography sx={{ mt: 2 }}>{t('common:common.loading')}</Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="md" sx={{ py: 4 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                {t('marketplace:store.edit.title')}
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error}
                </Alert>
            )}

            {success && (
                <Alert severity="success" sx={{ mb: 3 }}>
                    {t('marketplace:store.edit.success')}
                </Alert>
            )}

            <Paper sx={{ p: 3 }}>
                <form onSubmit={handleSubmit}>
                    <Grid container spacing={3}>
                        <Grid item xs={12}>
                            <Typography variant="h6" gutterBottom>
                                {t('marketplace:store.settings.basic')}
                            </Typography>
                        </Grid>

                        <Grid item xs={12}>
                            <TextField
                                name="name"
                                label={t('marketplace:store.settings.name')}
                                value={storefront.name || ''}
                                onChange={handleChange}
                                fullWidth
                                required
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <TextField
                                name="slug"
                                label={t('marketplace:store.settings.slug')}
                                value={storefront.slug || ''}
                                onChange={handleChange}
                                fullWidth
                                helperText={t('marketplace:store.edit.slugHelp')}
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <TextField
                                name="description"
                                label={t('common:common.description')}
                                value={storefront.description || ''}
                                onChange={handleChange}
                                fullWidth
                                multiline
                                rows={4}
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Divider sx={{ my: 2 }} />
                            <Typography variant="h6" gutterBottom>
                                {t('marketplace:store.settings.contactInfo')}
                            </Typography>
                        </Grid>

                        <Grid item xs={12} md={6}>
                            <TextField
                                name="phone"
                                label={t('marketplace:store.settings.phone')}
                                value={storefront.phone || ''}
                                onChange={handleChange}
                                fullWidth
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12} md={6}>
                            <TextField
                                name="email"
                                label={t('marketplace:store.settings.email')}
                                value={storefront.email || ''}
                                onChange={handleChange}
                                fullWidth
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <TextField
                                name="website"
                                label={t('marketplace:store.settings.website')}
                                value={storefront.website || ''}
                                onChange={handleChange}
                                fullWidth
                                disabled={saving}
                                placeholder="https://"
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Divider sx={{ my: 2 }} />
                            <Typography variant="h6" gutterBottom>
                                {t('marketplace:store.settings.location')}
                            </Typography>
                        </Grid>

                        <Grid item xs={12}>
                            <LocationPicker
                                onLocationSelect={handleLocationSelect}
                                initialLocation={{
                                    latitude: storefront.latitude,
                                    longitude: storefront.longitude,
                                    formatted_address: storefront.address,
                                    address_components: {
                                        city: storefront.city,
                                        country: storefront.country
                                    }
                                }}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 3 }}>
                                <Button
                                    variant="outlined"
                                    onClick={() => navigate(`/storefronts/${id}`)}
                                    disabled={saving}
                                >
                                    {t('common:buttons.cancel')}
                                </Button>
                                <Button
                                    type="submit"
                                    variant="contained"
                                    disabled={saving}
                                >
                                    {saving ? t('marketplace:store.edit.saving') : t('marketplace:store.edit.saveChanges')}
                                </Button>
                            </Box>
                        </Grid>
                    </Grid>
                </form>
            </Paper>
        </Container>
    );
};

export default EditStorefrontPage;
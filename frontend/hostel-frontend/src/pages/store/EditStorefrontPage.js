// frontend/hostel-frontend/src/pages/store/EditStorefrontPage.js
import React, { useState, useEffect } from 'react';
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
    Alert
} from '@mui/material';
import axios from '../../api/axios';

const EditStorefrontPage = () => {
    const { t } = useTranslation(['common', 'marketplace']);
    const { id } = useParams();
    const navigate = useNavigate();
    
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [storefront, setStorefront] = useState({
        name: '',
        description: '',
        slug: ''
    });
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(false);

    useEffect(() => {
        const fetchStorefront = async () => {
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/storefronts/${id}`);
                
                if (response.data?.data) {
                    setStorefront(response.data.data);
                }
            } catch (err) {
                console.error('Error fetching storefront:', err);
                setError('Не удалось загрузить данные витрины');
            } finally {
                setLoading(false);
            }
        };

        fetchStorefront();
    }, [id]);

    const handleSubmit = async (e) => {
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
            setError('Не удалось обновить данные витрины');
        } finally {
            setSaving(false);
        }
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setStorefront(prev => ({
            ...prev,
            [name]: value
        }));
    };

    if (loading) {
        return (
            <Container maxWidth="md" sx={{ py: 4, textAlign: 'center' }}>
                <CircularProgress />
                <Typography sx={{ mt: 2 }}>Загрузка данных витрины...</Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="md" sx={{ py: 4 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                Редактирование витрины
            </Typography>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error}
                </Alert>
            )}

            {success && (
                <Alert severity="success" sx={{ mb: 3 }}>
                    Витрина успешно обновлена
                </Alert>
            )}

            <Paper sx={{ p: 3 }}>
                <form onSubmit={handleSubmit}>
                    <Grid container spacing={3}>
                        <Grid item xs={12}>
                            <TextField
                                name="name"
                                label="Название магазина"
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
                                label="URL (slug)"
                                value={storefront.slug || ''}
                                onChange={handleChange}
                                fullWidth
                                helperText="Уникальный идентификатор для URL адреса витрины"
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <TextField
                                name="description"
                                label="Описание"
                                value={storefront.description || ''}
                                onChange={handleChange}
                                fullWidth
                                multiline
                                rows={4}
                                disabled={saving}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
                                <Button
                                    variant="outlined"
                                    onClick={() => navigate(`/storefronts/${id}`)}
                                    disabled={saving}
                                >
                                    Отмена
                                </Button>
                                <Button
                                    type="submit"
                                    variant="contained"
                                    disabled={saving}
                                >
                                    {saving ? 'Сохранение...' : 'Сохранить изменения'}
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
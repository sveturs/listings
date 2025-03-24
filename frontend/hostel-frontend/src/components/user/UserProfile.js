
import React, { useState, useEffect } from 'react';

import {
    Box,
    TextField,
    Button,
    Switch,
    FormControlLabel,
    Paper,
    Typography,
    Alert,
    Stack,
    Avatar,
    IconButton,
    Divider
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import { useAuth } from '../../contexts/AuthContext';
import { useLocation } from '../../contexts/LocationContext';
import axios from '../../api/axios';

const UserProfile = ({ onClose }) => {
    const { user } = useAuth();
    const { userLocation } = useLocation();
    const [profile, setProfile] = useState(null);
    const [isEditing, setIsEditing] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [formData, setFormData] = useState({
        phone: '',
        bio: '',
        notification_email: true,
        notification_push: true,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        city: '',
        country: ''
    });

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await axios.get('/api/v1/users/profile');
                setProfile(response.data.data);
                setFormData({
                    phone: response.data.data.phone || '',
                    bio: response.data.data.bio || '',
                    notification_email: response.data.data.notification_email,
                    notification_push: response.data.data.notification_push,
                    timezone: response.data.data.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone,
                    city: response.data.data.city || '',
                    country: response.data.data.country || ''
                });
            } catch (err) {
                setError('Ошибка загрузки профиля');
            }
        };
        fetchProfile();
    }, []);

    // Используем данные местоположения, если они доступны и нет данных в профиле
    useEffect(() => {
        if (userLocation && (!formData.city || !formData.country)) {
            setFormData(prev => ({
                ...prev,
                city: prev.city || userLocation.city || '',
                country: prev.country || userLocation.country || ''
            }));
        }
    }, [userLocation, formData.city, formData.country]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        try {
            await axios.put('/api/v1/users/profile', formData);
            setSuccess('Профиль успешно обновлен');
            setIsEditing(false);
            
            // Обновляем данные профиля после сохранения
            const response = await axios.get('/api/v1/users/profile');
            setProfile(response.data.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Ошибка обновления профиля');
        }
    };

    if (!profile) {
        return <Box sx={{ p: 3 }}><Typography>Загрузка...</Typography></Box>;
    }

    return (
        <Box sx={{ maxWidth: 600, mx: 'auto', position: 'relative' }}>
            {/* Close button */}
            <IconButton
                onClick={onClose}
                sx={{
                    position: 'absolute',
                    right: 8,
                    top: 8,
                    zIndex: 1
                }}
            >
                <CloseIcon />
            </IconButton>

            <Paper sx={{ p: 3 }}>
                <Stack spacing={3}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Avatar
                            src={profile.picture_url}
                            alt={profile.name}
                            sx={{ width: 80, height: 80 }}
                        />
                        <Box>
                            <Typography variant="h5">{profile.name}</Typography>
                            <Typography variant="body2" color="text.secondary">
                                {profile.email}
                            </Typography>
                            {profile.city && (
                                <Typography variant="body2" color="text.secondary">
                                    {profile.city}{profile.country ? `, ${profile.country}` : ''}
                                </Typography>
                            )}
                        </Box>
                    </Box>

                    {error && <Alert severity="error">{error}</Alert>}
                    {success && <Alert severity="success">{success}</Alert>}

                    <form onSubmit={handleSubmit}>
                        <Stack spacing={2}>
                            <TextField
                                label="Телефон"
                                value={formData.phone}
                                onChange={(e) => setFormData({
                                    ...formData,
                                    phone: e.target.value
                                })}
                                disabled={!isEditing}
                                fullWidth
                            />

                            <TextField
                                label="О себе"
                                value={formData.bio}
                                onChange={(e) => setFormData({
                                    ...formData,
                                    bio: e.target.value
                                })}
                                disabled={!isEditing}
                                multiline
                                rows={4}
                                fullWidth
                            />

                            <Divider sx={{ my: 1 }} />
                            <Typography variant="subtitle1">Местоположение</Typography>

                            <TextField
                                label="Город"
                                value={formData.city}
                                onChange={(e) => setFormData({
                                    ...formData,
                                    city: e.target.value
                                })}
                                disabled={!isEditing}
                                fullWidth
                            />

                            <TextField
                                label="Страна"
                                value={formData.country}
                                onChange={(e) => setFormData({
                                    ...formData,
                                    country: e.target.value
                                })}
                                disabled={!isEditing}
                                fullWidth
                            />

                            {isEditing ? (
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                    <Button
                                        type="submit"
                                        variant="contained"
                                        fullWidth
                                    >
                                        Сохранить
                                    </Button>
                                    <Button
                                        onClick={() => setIsEditing(false)}
                                        variant="outlined"
                                        fullWidth
                                    >
                                        Отмена
                                    </Button>
                                </Box>
                            ) : (
                                <Button
                                    onClick={() => setIsEditing(true)}
                                    variant="contained"
                                    fullWidth
                                >
                                    Редактировать
                                </Button>
                            )}
                        </Stack>
                    </form>
                </Stack>
            </Paper>
        </Box>
    );
};

export default UserProfile;
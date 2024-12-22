// src/components/RoomDetailsDialog.js
import React from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    Box,
    Grid,
    Typography,
    Chip,
    Divider,
    IconButton
} from '@mui/material';
import {
    Close as CloseIcon,
    SingleBed as SingleBedIcon,
    Hotel as HotelIcon,
    Apartment as ApartmentIcon,
    LocationOn as LocationIcon
} from '@mui/icons-material';
import ReviewsSection from '../reviews/ReviewsSection';
import { useAuth } from '../../contexts/AuthContext';
import ImageGallery from './ImageGallery';
import GalleryViewer from '../shared/GalleryViewer';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const RoomDetailsDialog = ({ open, onClose, room, onBook }) => {
    const { user } = useAuth();

    const getAccommodationIcon = () => {
        switch (room?.accommodation_type) {
            case 'bed':
                return <SingleBedIcon />;
            case 'apartment':
                return <ApartmentIcon />;
            default:
                return <HotelIcon />;
        }
    };

    const getAccommodationName = () => {
        switch (room?.accommodation_type) {
            case 'bed':
                return 'Койко-место';
            case 'apartment':
                return 'Апартаменты';
            default:
                return 'Комната';
        }
    };

    return (
        <Dialog
            open={open}
            onClose={onClose}
            maxWidth="md"
            fullWidth
        >
            <DialogTitle sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                pb: 1
            }}>
                <Box>
                    <Typography variant="h5" component="div">
                        {room?.name}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
                        <LocationIcon sx={{ mr: 1, fontSize: 20 }} />
                        <Typography variant="body2">
                            {room?.address_street}, {room?.address_city}
                        </Typography>
                    </Box>
                </Box>
                <IconButton onClick={onClose}>
                    <CloseIcon />
                </IconButton>
            </DialogTitle>

            <DialogContent>
                <Grid container spacing={3}>
                    {/* Галерея изображений */}
                    <Grid item xs={12}>
                        {room.images?.length > 0 && (
                            <GalleryViewer
                                images={room.images}
                                galleryMode="thumbnails"
                                thumbnailSize={{ width: '100%', height: '300px' }}
                                gridColumns={{ xs: 12 }}
                            />
                        )}
                    </Grid>

                    {/* Основная информация */}
                    <Grid item xs={12} md={8}>
                        <Box sx={{ mb: 3 }}>
                            <Chip
                                icon={getAccommodationIcon()}
                                label={getAccommodationName()}
                                sx={{ mr: 1 }}
                            />
                            <Chip
                                label={`${room?.price_per_night} ₽/ночь`}
                                color="primary"
                            />
                            {room?.accommodation_type === 'bed' && (
                                <Chip
                                    label={`${room?.available_beds}/${room?.total_beds} мест`}
                                    color="secondary"
                                    sx={{ ml: 1 }}
                                />
                            )}
                        </Box>

                        <Typography variant="h6" gutterBottom>
                            Характеристики
                        </Typography>
                        <Grid container spacing={2} sx={{ mb: 3 }}>
                            <Grid item xs={6}>
                                <Typography variant="body2" color="text.secondary">
                                    Вместимость: {room?.capacity} чел.
                                </Typography>
                            </Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" color="text.secondary">
                                    {room?.has_private_bathroom ? 'Отдельный санузел' : 'Общий санузел'}
                                </Typography>
                            </Grid>
                        </Grid>

                        {/* Отзывы */}
                        <Box sx={{ mt: 4 }}>
                            <ReviewsSection
                                entityType="room"
                                entityId={room?.id}
                                entityTitle={room?.name}
                                canReview={Boolean(user)}
                            />
                        </Box>
                    </Grid>

                    {/* Правая панель с кнопкой бронирования */}
                    <Grid item xs={12} md={4}>
                        <Box sx={{ position: 'sticky', top: 24 }}>
                            <Button
                                variant="contained"
                                fullWidth
                                onClick={() => {
                                    onClose();
                                    onBook(room);
                                }}
                            >
                                Забронировать
                            </Button>
                        </Box>
                    </Grid>
                </Grid>
            </DialogContent>
        </Dialog>
    );
};

export default RoomDetailsDialog;
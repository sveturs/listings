// frontend/hostel-frontend/src/components/accommodation/RoomDetailsDialog.js
import React, { useState } from 'react';

import {
    Dialog,
    DialogTitle,
    DialogContent,
    Grid,
    Typography,
    Chip,
    Box,
    IconButton,
    Button
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
import GalleryViewer from '../shared/GalleryViewer';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const RoomDetailsDialog = ({ open, onClose, room, onBook }) => {
    const { user } = useAuth();
    const [galleryOpen, setGalleryOpen] = useState(false);
    const [currentImageIndex, setCurrentImageIndex] = useState(0);

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

    const handleThumbnailClick = (index) => {
        setCurrentImageIndex(index); // Меняет текущее изображение
    };

    const handleImageClick = () => {
        setGalleryOpen(true); // Открывает полноэкранную галерею
    };

    return (
        <>
            <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
                <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', pb: 1 }}>
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
                        <Grid item xs={12}>
                            {room?.images?.length > 0 && (
                                <Box>
                                    {/* Большое изображение */}
                                    <Box
                                        component="img"
                                        src={`${BACKEND_URL}/uploads/${room.images[currentImageIndex].file_path}`}
                                        alt={`Image ${currentImageIndex + 1}`}
                                        sx={{
                                            width: '100%',
                                            height: '400px',
                                            objectFit: 'cover',
                                            borderRadius: 1,
                                            cursor: 'pointer',
                                        }}
                                        onClick={handleImageClick}
                                    />
                                    {/* Эскизы */}
                                    <Box
                                        sx={{
                                            mt: 2,
                                            display: 'flex',
                                            gap: 1,
                                            overflowX: 'auto',
                                            pb: 1,
                                        }}
                                    >
                                        {room.images.map((image, index) => (
                                            <Box
                                                key={index}
                                                component="img"
                                                src={`${BACKEND_URL}/uploads/${image.file_path}`}
                                                alt={`Thumbnail ${index + 1}`}
                                                sx={{
                                                    height: 80,
                                                    width: 120,
                                                    objectFit: 'cover',
                                                    borderRadius: 1,
                                                    cursor: 'pointer',
                                                    opacity: currentImageIndex === index ? 1 : 0.7,
                                                    transition: 'opacity 0.2s',
                                                    '&:hover': {
                                                        opacity: 1,
                                                    },
                                                }}
                                                onClick={() => handleThumbnailClick(index)}
                                            />
                                        ))}
                                    </Box>
                                </Box>
                            )}
                        </Grid>

                        <Grid item xs={12} md={8}>
                            <Box sx={{ mb: 3 }}>
                                <Chip
                                    icon={getAccommodationIcon()}
                                    label={getAccommodationName()}
                                    sx={{ mr: 1 }}
                                />
                                <Chip label={`${room?.price_per_night} ₽/ночь`} color="primary" />
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
                            <Box sx={{ mt: 4 }}>
                                <ReviewsSection
                                    entityType="room"
                                    entityId={room?.id}
                                    entityTitle={room?.name}
                                    canReview={Boolean(user)}
                                />
                            </Box>
                        </Grid>

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
            <GalleryViewer
                images={room?.images || []}
                open={galleryOpen}
                onClose={() => setGalleryOpen(false)}
                initialIndex={currentImageIndex}
                galleryMode="fullscreen"
            />
        </>
    );
};

export default RoomDetailsDialog;

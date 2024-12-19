// frontend/hostel-frontend/src/components/car/CarDetailsDialog.js
import React from 'react';
import {
    Dialog,
    DialogContent,
    DialogTitle,
    IconButton,
    Typography,
    Box,
    Grid,
    Chip,
    Button,
    Rating,
    Divider,
} from '@mui/material';
import {
    Close as CloseIcon,
    LocalGasStation as FuelIcon,
    Speed as TransmissionIcon,
    AirlineSeatReclineNormal as SeatsIcon,
    LocationOn as LocationIcon,
} from '@mui/icons-material';
import ReviewsSection from '../reviews/ReviewsSection';
import { useAuth } from '../../contexts/AuthContext';

const CarDetailsDialog = ({ open, onClose, car, onBook }) => {
    const { user } = useAuth();
    const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;
    console.log('User:', user);
    console.log('Car:', car);
    console.log('Can Review:', Boolean(user));
    if (!car) return null;

    return (
        <Dialog
            open={open}
            onClose={onClose}
            maxWidth="md"
            fullWidth
        >
            <DialogTitle
                sx={{
                    m: 0,
                    p: 2,
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center'
                }}
            >
                <Box component="div"> {/* Оборачиваем в Box чтобы избежать вложенности h-тегов */}
                    <Typography component="span" variant="h6">
                        {car.make} {car.model} {car.year}
                    </Typography>
                </Box>
                <IconButton
                    onClick={onClose}
                    size="small"
                >
                    <CloseIcon />
                </IconButton>
            </DialogTitle>

            <DialogContent dividers>
                <Grid container spacing={3}>
                    {/* Галерея изображений */}
                    <Grid item xs={12}>
                        <Box sx={{ display: 'flex', gap: 1, overflowX: 'auto', pb: 2 }}>
                            {car.images?.map((image, index) => (
                                <Box
                                    key={index}
                                    component="img"
                                    src={`${BACKEND_URL}/uploads/${image.file_path}`}
                                    alt={`${car.make} ${car.model}`}
                                    sx={{
                                        height: 300,
                                        minWidth: 400,
                                        objectFit: 'cover',
                                        borderRadius: 1,
                                    }}
                                />
                            ))}
                        </Box>
                    </Grid>

                    {/* Основная информация */}
                    <Grid item xs={12} md={8}>
                        <Typography variant="h6" gutterBottom>
                            Характеристики
                        </Typography>

                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                            <Chip
                                icon={<FuelIcon />}
                                label={car.fuel_type === 'petrol' ? 'Бензин' :
                                    car.fuel_type === 'diesel' ? 'Дизель' :
                                        car.fuel_type === 'electric' ? 'Электро' : 'Гибрид'}
                            />
                            <Chip
                                icon={<TransmissionIcon />}
                                label={car.transmission === 'automatic' ? 'Автомат' : 'Механика'}
                            />
                            <Chip
                                icon={<SeatsIcon />}
                                label={`${car.seats} мест`}
                            />
                        </Box>

                        <Typography variant="body1" paragraph>
                            {car.description}
                        </Typography>

                        {car.features?.length > 0 && (
                            <>
                                <Typography variant="subtitle1" gutterBottom>
                                    Комплектация
                                </Typography>
                                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                                    {car.features.map((feature, index) => (
                                        <Chip key={index} label={feature} size="small" />
                                    ))}
                                </Box>
                            </>
                        )}
                    </Grid>

                    {/* Боковая информация */}
                    <Grid item xs={12} md={4}>
                        <Box sx={{ p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
                            <Typography variant="h6" color="primary" gutterBottom>
                                {car.price_per_day} ₽ в день
                            </Typography>

                            <Box sx={{ mb: 2 }}>
                                <Typography variant="body2" color="text.secondary" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <LocationIcon fontSize="small" />
                                    {car.location}
                                </Typography>
                            </Box>

                            <Button
                                variant="contained"
                                fullWidth
                                onClick={() => {
                                    onClose();
                                    onBook(car);
                                }}
                            >
                                Забронировать
                            </Button>
                        </Box>
                    </Grid>

                    {/* Отзывы */}
                    <Grid item xs={12}>
                        <Divider sx={{ my: 3 }} />
                        <ReviewsSection
                            entityType="car"
                            entityId={car.id}
                            entityTitle={`${car.make} ${car.model}`}
                            canReview={Boolean(user)}
                        />
                    </Grid>
                </Grid>
            </DialogContent>
        </Dialog>
    );
};

export default CarDetailsDialog;
import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import ReviewsSection from './reviews/ReviewsSection';

import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Box,
    Typography,
    Alert,
    Grid,
    Paper,
    Divider,
    Chip,
    FormControlLabel,
    Checkbox,
} from '@mui/material';
import {
    DirectionsCar as CarIcon,
    LocalGasStation as FuelIcon,
    Speed as TransmissionIcon,
    EventAvailable as DateIcon,
    LocationOn as LocationIcon,
    CreditCard as PaymentIcon,
    Close as CloseIcon,
    Assignment as InsuranceIcon,
} from '@mui/icons-material';
import axios from '../api/axios';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const CarBookingDialog = ({ open, onClose, car, startDate, endDate }) => {
    const { user } = useAuth();
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const [loading, setLoading] = useState(false);
    const [bookingData, setBookingData] = useState({
        pickup_location: car?.location || '',
        dropoff_location: car?.location || '',
        include_insurance: true,
        additional_notes: ''
    });

    // Расчет общей стоимости
    const calculateTotalPrice = () => {
        if (!startDate || !endDate) return 0;

        const start = new Date(startDate);
        const end = new Date(endDate);
        const days = Math.ceil((end - start) / (1000 * 60 * 60 * 24));

        let total = car.price_per_day * days;

        // Добавляем стоимость страховки если выбрана
        if (bookingData.include_insurance) {
            total += days * 500; // Предположим, страховка стоит 500р в день
        }

        return total;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess(false);
        setLoading(true);

        if (!user) {
            setError('Необходимо войти в систему для бронирования');
            setLoading(false);
            return;
        }

        try {
            const bookingPayload = {
                car_id: car.id,
                start_date: startDate,
                end_date: endDate,
                pickup_location: bookingData.pickup_location,
                dropoff_location: bookingData.dropoff_location,
                include_insurance: bookingData.include_insurance,
                additional_notes: bookingData.additional_notes,
                total_price: calculateTotalPrice()
            };

            await axios.post('/api/v1/car-bookings', bookingPayload, {
                withCredentials: true
            });

            setSuccess(true);
            setTimeout(() => {
                onClose();
            }, 2000);
        } catch (error) {
            if (error.response?.status === 401) {
                setError('Необходимо войти в систему');
            } else {
                setError(error.response?.data?.error || 'Произошла ошибка при бронировании');
            }
        } finally {
            setLoading(false);
        }
    };

    const renderCarDetails = () => (
        <Box sx={{ mb: 3 }}>
            <Typography variant="h6" gutterBottom>
                {car.make} {car.model}
                <Typography component="span" color="text.secondary" sx={{ ml: 1 }}>
                    {car.year}
                </Typography>
            </Typography>

            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                <Chip
                    icon={<FuelIcon />}
                    label={car.fuel_type === 'petrol' ? 'Бензин' :
                        car.fuel_type === 'diesel' ? 'Дизель' :
                            car.fuel_type === 'electric' ? 'Электро' : 'Гибрид'}
                    size="small"
                />
                <Chip
                    icon={<TransmissionIcon />}
                    label={car.transmission === 'automatic' ? 'Автомат' : 'Механика'}
                    size="small"
                />
                <Chip
                    label={`${car.seats} мест`}
                    size="small"
                />
            </Box>

            {car.features?.length > 0 && (
                <Typography variant="body2" color="text.secondary">
                    {car.features.join(' • ')}
                </Typography>
            )}
        </Box>
    );

    const renderPriceDetails = () => {
        const totalPrice = calculateTotalPrice();
        const days = Math.ceil((new Date(endDate) - new Date(startDate)) / (1000 * 60 * 60 * 24));

        return (
            <Paper variant="outlined" sx={{ p: 2, mt: 3 }}>
                <Typography variant="h6" gutterBottom>
                    Детали оплаты
                </Typography>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="body2">
                        Аренда ({days} {days === 1 ? 'день' : days < 5 ? 'дня' : 'дней'})
                    </Typography>
                    <Typography variant="body2">
                        {car.price_per_day * days} ₽
                    </Typography>
                </Box>
                {bookingData.include_insurance && (
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                        <Typography variant="body2">
                            Страховка ({days} {days === 1 ? 'день' : days < 5 ? 'дня' : 'дней'})
                        </Typography>
                        <Typography variant="body2">
                            {500 * days} ₽
                        </Typography>
                    </Box>
                )}
                <Divider sx={{ my: 1 }} />
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="subtitle1" fontWeight="bold">
                        Итого
                    </Typography>
                    <Typography variant="subtitle1" fontWeight="bold" color="primary">
                        {totalPrice} ₽
                    </Typography>
                </Box>
            </Paper>
        );
    };


    return (
        <Dialog
            open={open}
            onClose={onClose}
            maxWidth="md"
            fullWidth
            PaperProps={{
                sx: { p: 2 }
            }}
        >
            <DialogTitle sx={{ p: 0, mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="h5">
                        Бронирование автомобиля
                    </Typography>
                    <Button
                        onClick={onClose}
                        color="inherit"
                        startIcon={<CloseIcon />}
                    >
                        Закрыть
                    </Button>
                </Box>
            </DialogTitle>

            {!user ? (
                <Alert
                    severity="warning"
                    action={
                        <Button color="inherit" size="small">
                            Войти
                        </Button>
                    }
                >
                    Для бронирования необходимо войти в систему
                </Alert>
            ) : (
                <>
                    <form onSubmit={handleSubmit}>
                        <DialogContent sx={{ p: 0 }}>
                            {error && (
                                <Alert severity="error" sx={{ mb: 2 }}>
                                    {error}
                                </Alert>
                            )}
                            {success && (
                                <Alert severity="success" sx={{ mb: 2 }}>
                                    Бронирование успешно создано!
                                </Alert>
                            )}

                            {renderCarDetails()}

                            <Divider sx={{ my: 3 }} />

                            <Grid container spacing={3}>
                                <Grid item xs={12} sm={6}>
                                    <Typography variant="subtitle2" sx={{ mb: 1 }}>
                                        <DateIcon sx={{ mr: 1, verticalAlign: 'bottom' }} />
                                        Период аренды
                                    </Typography>
                                    <Box sx={{ display: 'flex', gap: 2 }}>
                                        <Typography variant="body2" color="text.secondary">
                                            {new Date(startDate).toLocaleDateString()} - {new Date(endDate).toLocaleDateString()}
                                        </Typography>
                                    </Box>
                                </Grid>

                                <Grid item xs={12}>
                                    <Typography variant="subtitle2" sx={{ mb: 1 }}>
                                        <LocationIcon sx={{ mr: 1, verticalAlign: 'bottom' }} />
                                        Место получения и возврата
                                    </Typography>
                                    <Grid container spacing={2}>
                                        <Grid item xs={12} sm={6}>
                                            <TextField
                                                label="Место получения"
                                                fullWidth
                                                size="small"
                                                value={bookingData.pickup_location}
                                                onChange={(e) => setBookingData({
                                                    ...bookingData,
                                                    pickup_location: e.target.value
                                                })}
                                                required
                                            />
                                        </Grid>
                                        <Grid item xs={12} sm={6}>
                                            <TextField
                                                label="Место возврата"
                                                fullWidth
                                                size="small"
                                                value={bookingData.dropoff_location}
                                                onChange={(e) => setBookingData({
                                                    ...bookingData,
                                                    dropoff_location: e.target.value
                                                })}
                                                required
                                            />
                                        </Grid>
                                    </Grid>
                                </Grid>

                                <Grid item xs={12}>
                                    <FormControlLabel
                                        control={
                                            <Checkbox
                                                checked={bookingData.include_insurance}
                                                onChange={(e) => setBookingData({
                                                    ...bookingData,
                                                    include_insurance: e.target.checked
                                                })}
                                                icon={<InsuranceIcon />}
                                                checkedIcon={<InsuranceIcon />}
                                            />
                                        }
                                        label="Включить страховку (500₽ в день)"
                                    />
                                </Grid>

                                <Grid item xs={12}>
                                    <TextField
                                        label="Дополнительные пожелания"
                                        fullWidth
                                        multiline
                                        rows={3}
                                        value={bookingData.additional_notes}
                                        onChange={(e) => setBookingData({
                                            ...bookingData,
                                            additional_notes: e.target.value
                                        })}
                                    />
                                </Grid>
                            </Grid>

                            {renderPriceDetails()}
                        </DialogContent>

                        <DialogActions sx={{ p: 0, mt: 3 }}>
                            <Button onClick={onClose} color="inherit">
                                Отмена
                            </Button>
                            <Button
                                type="submit"
                                variant="contained"
                                disabled={loading || !startDate || !endDate}
                                startIcon={<PaymentIcon />}
                            >
                                Забронировать
                            </Button>
                        </DialogActions>
                    </form>


                </>
            )}
        </Dialog>
    );
};

export default CarBookingDialog;
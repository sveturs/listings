import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
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
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Grid,
    //    Tooltip,
    Popover,
} from '@mui/material';
import axios from "../api/axios";

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const BookingDialog = ({ open, onClose, room, startDate, endDate }) => {
    const { user } = useAuth(); // Получаем данные пользователя из контекста
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const [selectedBed, setSelectedBed] = useState('');
    const [availableBeds, setAvailableBeds] = useState([]);
    const [bedImages, setBedImages] = useState({});
    const [anchorEl, setAnchorEl] = useState(null);
    const [activeBed, setActiveBed] = useState(null);
    const [bookingStartDate, setBookingStartDate] = useState(startDate);
    const [bookingEndDate, setBookingEndDate] = useState(endDate);

    useEffect(() => {
        if (open && room && room.accommodation_type === 'bed' && bookingStartDate && bookingEndDate) {
            // Сбрасываем предыдущий выбор койки при изменении дат
            setSelectedBed('');
            setError('');

            axios.get(`/rooms/${room.id}/available-beds`, {
                params: {
                    start_date: bookingStartDate,
                    end_date: bookingEndDate
                }
            })
                .then(response => {
                    setAvailableBeds(response.data);
                    response.data.forEach(bed => {
                        axios.get(`/beds/${bed.id}/images`)
                            .then(imgResponse => {
                                setBedImages(prev => ({
                                    ...prev,
                                    [bed.id]: imgResponse.data
                                }));
                            })
                            .catch(console.error);
                    });
                })
                .catch(err => {
                    console.error('Ошибка загрузки доступных койко-мест:', err);
                    setError('Не удалось загрузить список доступных койко-мест');
                });
        }
    }, [open, room, bookingStartDate, bookingEndDate]);
    const handleMouseEnter = (event, bed) => {
        if (bedImages[bed.id]?.length > 0) {
            setActiveBed(bed);
            setAnchorEl(event.currentTarget);
        }
    };

    const handleMouseLeave = () => {
        setActiveBed(null);
        setAnchorEl(null);
    };

    const renderBedImage = () => {
        if (!activeBed || !bedImages[activeBed.id]?.length) return null;

        const image = bedImages[activeBed.id][0];
        return (
            <Box sx={{ p: 1 }}>
                <img
                    src={`${BACKEND_URL}/uploads/${image.file_path}`}
                    alt={`Койко-место ${activeBed.bed_number}`}
                    style={{
                        width: '200px',
                        height: '150px',
                        objectFit: 'cover',
                        borderRadius: '4px'
                    }}
                />
            </Box>
        );
    };

    const calculateTotalPrice = () => {
        if (!bookingStartDate || !bookingEndDate) return 0;
    
        const start = new Date(bookingStartDate);
        const end = new Date(bookingEndDate);
        const daysCount = Math.ceil((end - start) / (1000 * 60 * 60 * 24));

        let pricePerNight;
        if (room.accommodation_type === 'bed' && selectedBed) {
            const selectedBedData = availableBeds.find(bed => bed.id === selectedBed);
            pricePerNight = selectedBedData ? selectedBedData.price_per_night : 0;
        } else {
            pricePerNight = room.price_per_night;
        }

        return pricePerNight * daysCount;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess(false);

        if (!user) {
            setError('Необходимо войти в систему для бронирования');
            return;
        }

        if (!bookingStartDate || !bookingEndDate) {
            setError('Выберите даты проживания');
            return;
        }

        if (room.accommodation_type === 'bed' && !selectedBed) {
            setError('Выберите койко-место');
            return;
        }

        if (bookingStartDate === bookingEndDate) {
            setError('Дата выезда должна быть позже даты заезда');
            return;
        }

        try {
            const bookingData = {
                room_id: room.id,
                start_date: bookingStartDate,
                end_date: bookingEndDate
            };

            if (room.accommodation_type === 'bed') {
                bookingData.bed_id = selectedBed;
            }

            // Отправляем запрос с куками для авторизации
            const response = await axios.post('/bookings', bookingData, {
                withCredentials: true
            });

            setSuccess(true);
            setTimeout(() => {
                onClose();
                setSelectedBed('');
                setError('');
                setSuccess(false);
            }, 2000);
        } catch (error) {
            if (error.response?.status === 401) {
                setError('Необходимо войти в систему');
            } else {
                setError(error.response?.data || 'Произошла ошибка при бронировании');
            }
        }
    };

    const getDialogTitle = () => {
        switch (room?.accommodation_type) {
            case 'bed':
                return 'Бронирование койко-места';
            case 'apartment':
                return 'Бронирование квартиры';
            default:
                return 'Бронирование комнаты';
        }
    };

    const today = new Date().toISOString().split('T')[0];

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>{getDialogTitle()}</DialogTitle>
            <DialogContent>
                {!user ? (
                    <Alert severity="warning" sx={{ mt: 2 }}>
                        Для бронирования необходимо войти в систему
                    </Alert>
                ) : (
                    <>
                        {success && (
                            <Alert severity="success" sx={{ mt: 2 }}>
                                Бронирование успешно создано!
                            </Alert>
                        )}
                        {error && (
                            <Alert severity="error" sx={{ mt: 2 }}>
                                {error}
                            </Alert>
                        )}
                        {room && (
                            <Box sx={{ mt: 2 }}>
                                <Typography variant="h6">{room.name}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {room.address_street}, {room.address_city}
                                </Typography>
    
                                <Grid container spacing={2} sx={{ mt: 1, mb: 2 }}>
                                    <Grid item xs={12} sm={6}>
                                        <TextField
                                            label="Дата заезда"
                                            type="date"
                                            fullWidth
                                            value={bookingStartDate}
                                            onChange={(e) => setBookingStartDate(e.target.value)}
                                            inputProps={{ min: today }}
                                            InputLabelProps={{ shrink: true }}
                                        />
                                    </Grid>
                                    <Grid item xs={12} sm={6}>
                                        <TextField
                                            label="Дата выезда"
                                            type="date"
                                            fullWidth
                                            value={bookingEndDate}
                                            onChange={(e) => setBookingEndDate(e.target.value)}
                                            inputProps={{ min: bookingStartDate || today }}
                                            InputLabelProps={{ shrink: true }}
                                        />
                                    </Grid>
                                </Grid>
    
                                {room.accommodation_type === 'bed' && (
                                    <FormControl fullWidth sx={{ mt: 2 }}>
                                        <InputLabel>Выберите койко-место</InputLabel>
                                        <Select
                                            value={selectedBed}
                                            onChange={(e) => setSelectedBed(e.target.value)}
                                            label="Выберите койко-место"
                                        >
                                            {availableBeds.map(bed => (
                                                <MenuItem
                                                    key={bed.id}
                                                    value={bed.id}
                                                    onMouseEnter={(e) => handleMouseEnter(e, bed)}
                                                    onMouseLeave={handleMouseLeave}
                                                >
                                                    Место {bed.bed_number} - {bed.price_per_night} руб./ночь
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                )}
    
                                <Typography variant="h6" sx={{ mt: 2 }}>
                                    Итого к оплате: {calculateTotalPrice()} руб.
                                </Typography>
                            </Box>
                        )}
                    </>
                )}
    
                <Popover
                    open={Boolean(anchorEl)}
                    anchorEl={anchorEl}
                    onClose={handleMouseLeave}
                    anchorOrigin={{
                        vertical: 'center',
                        horizontal: 'right',
                    }}
                    transformOrigin={{
                        vertical: 'center',
                        horizontal: 'left',
                    }}
                    sx={{
                        pointerEvents: 'none',
                    }}
                >
                    {renderBedImage()}
                </Popover>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="inherit">
                    Отмена
                </Button>
                <Button
                    onClick={handleSubmit}
                    color="primary"
                    variant="contained"
                    disabled={!user ||
                        !bookingStartDate ||
                        !bookingEndDate ||
                        bookingStartDate === bookingEndDate ||
                        (room?.accommodation_type === 'bed' && !selectedBed)}
                >
                    Забронировать
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default BookingDialog;
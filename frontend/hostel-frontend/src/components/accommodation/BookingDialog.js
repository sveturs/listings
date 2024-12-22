//hostel-booking-system/frontend/hostel-frontend/src/components/accommodation/BookingDialog.js
import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { ReviewsSection } from '../reviews'; 
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
  IconButton,
  Chip,
} from '@mui/material';
import {
  BedOutlined as BedIcon,
  CalendarMonth as CalendarIcon,
  LocationOn as LocationIcon,
  CreditCard as PaymentIcon,
  Close as CloseIcon,
  HotelOutlined as HotelIcon,
  PowerOutlined as PowerIcon,
  LightbulbOutlined as LightIcon,
  Storage as ShelfIcon,
} from '@mui/icons-material';
import axios from "../../api/axios";

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

// Компонент карточки койко-места
// Компонент карточки койко-места
const BedCard = ({ bed, selected, onClick, bedImages, disabled }) => (
  <Paper
    elevation={selected ? 3 : 1}
    sx={{
      p: 2,
      border: '1px solid',
      borderColor: selected ? 'primary.main' : 'divider',
      borderRadius: 2,
      cursor: disabled ? 'default' : 'pointer',
      transition: 'all 0.2s',
      bgcolor: disabled ? 'action.disabledBackground' : selected ? 'action.selected' : 'background.paper',
      '&:hover': !disabled && {
        borderColor: 'primary.main',
        transform: 'translateY(-2px)',
        boxShadow: 3,
      },
      opacity: disabled ? 0.7 : 1,
      height: '100%',
      display: 'flex',
      flexDirection: 'column',
    }}
    onClick={!disabled ? onClick : undefined}
  >
    {/* Изображение кровати */}
    <Box
      sx={{
        position: 'relative',
        mb: 2
      }}
    >
      <Box
        sx={{
          width: '100%',
          height: 160,
          borderRadius: 1,
          overflow: 'hidden',
          bgcolor: 'grey.100',
          position: 'relative'
        }}
      >
        {bedImages?.[bed.id]?.length > 0 ? (
          <img
            src={`${BACKEND_URL}/uploads/${bedImages[bed.id][0].file_path}`}
            alt={`Койко-место ${bed.bed_number}`}
            style={{
              width: '100%',
              height: '100%',
              objectFit: 'cover',
            }}
          />
        ) : (
          <Box
            sx={{
              height: '100%',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              color: 'text.secondary',
            }}
          >
            <HotelIcon sx={{ fontSize: 40, mb: 1 }} />
            <Typography variant="body2">Фото не доступно</Typography>
          </Box>
        )}
      </Box>

      {selected && (
        <Chip
          label="Выбрано"
          color="primary"
          size="small"
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
          }}
        />
      )}
    </Box>

    {/* Информация о кровати */}
    <Box sx={{ flex: 1 }}>
      <Typography variant="h6" gutterBottom>
        Место {bed.bed_number}
      </Typography>

      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
        {bed.bed_type === 'single' ? 'Отдельностоящая кровать' :
          bed.bed_type === 'top' ? 'Верхний ярус' : 'Нижний ярус'}
      </Typography>

      <Grid container spacing={1} sx={{ mb: 2 }}>
        {bed.has_outlet && (
          <Grid item xs={6}>
            <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
              <PowerIcon sx={{ mr: 1, fontSize: 20 }} />
              <Typography variant="body2">Розетка</Typography>
            </Box>
          </Grid>
        )}
        {bed.has_light && (
          <Grid item xs={6}>
            <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
              <LightIcon sx={{ mr: 1, fontSize: 20 }} />
              <Typography variant="body2">Светильник</Typography>
            </Box>
          </Grid>
        )}
        {bed.has_shelf && (
          <Grid item xs={6}>
            <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
              <ShelfIcon sx={{ mr: 1, fontSize: 20 }} />
              <Typography variant="body2">Полка</Typography>
            </Box>
          </Grid>
        )}
      </Grid>
    </Box>

    {/* Цена */}
    <Box sx={{ pt: 2, borderTop: '1px solid', borderColor: 'divider' }}>
      <Typography variant="h6" color="primary.main" align="center">
        {bed.price_per_night} ₽
        <Typography component="span" variant="body2" color="text.secondary">
          /ночь
        </Typography>
      </Typography>
    </Box>
  </Paper>
);

const BookingDialog = ({ open, onClose, room, startDate, endDate }) => {
  const { user } = useAuth();
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [selectedBed, setSelectedBed] = useState('');
  const [availableBeds, setAvailableBeds] = useState([]);
  const [bedImages, setBedImages] = useState({});
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
          setAvailableBeds(response.data.data || []);
          response.data.data.forEach(bed => {
            axios.get(`/beds/${bed.id}/images`)
              .then(imgResponse => {
                setBedImages(prev => ({
                  ...prev,
                  [bed.id]: imgResponse.data.data || []
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

      await axios.post('/api/v1/bookings', bookingData, {
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
        setError(error.response?.data?.error || 'Произошла ошибка при бронировании');
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
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 2,
        }
      }}
    >
      {/* Заголовок */}
      <DialogTitle sx={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'flex-start',
        pb: 1
      }}>
        <Box>
          <Typography variant="h5" component="div" gutterBottom>
            {getDialogTitle()}
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary' }}>
            <LocationIcon sx={{ mr: 1, fontSize: 20 }} />
            <Typography variant="body2">
              {room?.address_street}, {room?.address_city}
            </Typography>
          </Box>
        </Box>
        <IconButton onClick={onClose} size="small">
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <Divider />

      <DialogContent>
        {!user ? (
          <Alert
            severity="warning"
            sx={{ mt: 2 }}
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

            {/* Секция выбора дат */}
            <Box sx={{ mt: 2 }}>
              <Typography variant="subtitle1" sx={{
                display: 'flex',
                alignItems: 'center',
                mb: 2
              }}>
                <CalendarIcon sx={{ mr: 1 }} />
                Даты проживания
              </Typography>

              <Grid container spacing={2}>
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
            </Box>

            {/* Секция выбора койко-места */}
            {room?.accommodation_type === 'bed' && bookingStartDate && bookingEndDate && (
              <Box sx={{ mt: 4 }}>
                <Typography variant="subtitle1" sx={{
                  display: 'flex',
                  alignItems: 'center',
                  mb: 2
                }}>
                  <BedIcon sx={{ mr: 1 }} />
                  Доступные койко-места
                </Typography>

                <Grid container spacing={2}>
                  {availableBeds.map(bed => (
                    <Grid item xs={12} sm={6} md={4} key={bed.id}>
                      <BedCard
                        bed={bed}
                        selected={selectedBed === bed.id}
                        onClick={() => setSelectedBed(bed.id)}
                        bedImages={bedImages}
                      />
                    </Grid>
                  ))}
                </Grid>
              </Box>
            )}

            {/* Секция с итоговой стоимостью */}
            {(bookingStartDate && bookingEndDate) && (
              <Paper
                variant="outlined"
                sx={{
                  mt: 4,
                  p: 2,
                  bgcolor: 'background.default'
                }}
              >
                <Box sx={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between'
                }}>
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <PaymentIcon sx={{ mr: 1 }} />
                    <Typography variant="subtitle1">
                      Итого к оплате:
                    </Typography>
                  </Box>
                  <Typography variant="h5" color="primary.main">
                    {calculateTotalPrice()} ₽
                  </Typography>
                </Box>
              </Paper>
            )}
          </>
        )}
        <Box sx={{ mt: 4 }}>
          <ReviewsSection
            entityType="room"
            entityId={room?.id}
            entityTitle={room?.name}
            canReview={Boolean(user)}
          />
        </Box>
      </DialogContent>

      <Divider />

      <DialogActions sx={{ p: 2 }}>
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
          startIcon={<CalendarIcon />}
        >
          Забронировать
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default BookingDialog;
import React, { useState, useEffect } from 'react';
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
  MenuItem
} from '@mui/material';
import axios from "../api/axios";

const BookingDialog = ({ open, onClose, room, startDate, endDate }) => {
  const [userId, setUserId] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [selectedBed, setSelectedBed] = useState('');
  const [availableBeds, setAvailableBeds] = useState([]);

  useEffect(() => {
    if (open && room && room.accommodation_type === 'bed' && startDate && endDate) {
      axios.get(`/rooms/${room.id}/available-beds`, {
        params: {
          start_date: startDate,
          end_date: endDate
        }
      })
      .then(response => {
        setAvailableBeds(response.data);
        setSelectedBed('');
        setError('');
      })
      .catch(err => {
        console.error('Ошибка загрузки доступных койко-мест:', err);
        setError('Не удалось загрузить список доступных койко-мест');
      });
    }
  }, [open, room, startDate, endDate]);

  const calculateTotalPrice = () => {
    if (!startDate || !endDate) return 0;
    
    const start = new Date(startDate);
    const end = new Date(endDate);
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

    if (!startDate || !endDate) {
      setError('Выберите даты проживания');
      return;
    }

    if (!userId) {
      setError('Введите ID пользователя');
      return;
    }

    if (room.accommodation_type === 'bed' && !selectedBed) {
      setError('Выберите койко-место');
      return;
    }

    try {
      const bookingData = {
        user_id: parseInt(userId),
        room_id: room.id,
        start_date: startDate,
        end_date: endDate
      };

      if (room.accommodation_type === 'bed') {
        bookingData.bed_id = selectedBed;
      }

      await axios.post('/bookings', bookingData);
      setSuccess(true);
      
      setTimeout(() => {
        onClose();
        setUserId('');
        setSelectedBed('');
        setError('');
        setSuccess(false);
      }, 2000);
    } catch (error) {
      setError(error.response?.data || 'Произошла ошибка при бронировании');
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

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{getDialogTitle()}</DialogTitle>
      <DialogContent>
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
            <Typography sx={{ mt: 2 }}>
              Период проживания: {startDate || '-'} - {endDate || '-'}
            </Typography>

            {room.accommodation_type === 'bed' && (
              <FormControl fullWidth sx={{ mt: 2 }}>
                <InputLabel>Выберите койко-место</InputLabel>
                <Select
                  value={selectedBed}
                  onChange={(e) => setSelectedBed(e.target.value)}
                  label="Выберите койко-место"
                >
                  {availableBeds.map(bed => (
                    <MenuItem key={bed.id} value={bed.id}>
                      Место {bed.bed_number} - {bed.price_per_night} руб./ночь
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}

            <TextField
              margin="dense"
              label="ID пользователя"
              type="number"
              fullWidth
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              sx={{ mt: 2 }}
            />

            <Typography variant="h6" sx={{ mt: 2 }}>
              Итого к оплате: {calculateTotalPrice()} руб.
            </Typography>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="inherit">
          Отмена
        </Button>
        <Button 
          onClick={handleSubmit} 
          color="primary" 
          variant="contained"
          disabled={
            !userId || 
            !startDate || 
            !endDate || 
            (room?.accommodation_type === 'bed' && !selectedBed)
          }
        >
          Забронировать
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default BookingDialog;
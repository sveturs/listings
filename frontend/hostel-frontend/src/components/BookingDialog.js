import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Typography,
  Alert
} from '@mui/material';
import axios from "../api/axios";

const BookingDialog = ({ open, onClose, room, startDate, endDate }) => {
  const [userId, setUserId] = React.useState('');
  const [error, setError] = React.useState('');
  const [success, setSuccess] = React.useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess(false);

    try {
      await axios.post('/bookings', {
        user_id: parseInt(userId),
        room_id: room.id,
        start_date: startDate,
        end_date: endDate
      });
      setSuccess(true);
      setTimeout(() => {
        onClose();
        setUserId('');
        setSuccess(false);
      }, 2000);
    } catch (error) {
      setError(error.response?.data || 'Произошла ошибка при бронировании');
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Бронирование комнаты</DialogTitle>
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
            <Typography sx={{ mt: 1 }}>
              Даты: {startDate} - {endDate}
            </Typography>
            <Typography sx={{ mt: 1 }}>
              Общая стоимость: {room.price_per_night * 
                (Math.ceil((new Date(endDate) - new Date(startDate)) / (1000 * 60 * 60 * 24)))} руб.
            </Typography>
          </Box>
        )}
        <TextField
          autoFocus
          margin="dense"
          label="ID пользователя"
          type="number"
          fullWidth
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          sx={{ mt: 2 }}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="inherit">
          Отмена
        </Button>
        <Button onClick={handleSubmit} color="primary" variant="contained">
          Забронировать
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default BookingDialog;
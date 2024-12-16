import React, { useState, useEffect } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Chip,
  Box,
  CircularProgress,
  Alert
} from '@mui/material';
import {
  Apartment as ApartmentIcon,
  Hotel as HotelIcon,
  SingleBed as SingleBedIcon
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import axios from '../api/axios';

const BookingsList = () => {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { user } = useAuth();

  useEffect(() => {
    const fetchBookings = async () => {
        try {
            setLoading(true);
            
            const response = await axios.get('/api/v1/bookings', {
                withCredentials: true
            });
            
            console.log('User ID:', user?.id);
            console.log('API Response:', response.data);

            if (response.data?.data) {
                // Преобразуем ID в строки для сравнения
                const userBookings = response.data.data.filter(booking => 
                    String(booking.user_id) === String(user?.id)
                );
                setBookings(userBookings);
            }
            
        } catch (error) {
            console.error('Error:', error);
            setError('Не удалось загрузить бронирования');
        } finally {
            setLoading(false);
        }
    };

    if (user?.id) {
        fetchBookings();
    } else {
        setLoading(false);
    }
}, [user]);

  const getAccommodationIcon = (type) => {
    switch (type) {
      case 'bed':
        return <SingleBedIcon color="primary" />;
      case 'room':
        return <HotelIcon color="primary" />;
      case 'apartment':
        return <ApartmentIcon color="primary" />;
      default:
        return null;
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" p={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>
    );
  }

  if (!bookings || bookings.length === 0) {
    return (
      <Box sx={{ mt: 2 }}>
        <Typography variant="h5" gutterBottom>
          Мои бронирования
        </Typography>
        <Alert severity="info">
          У вас пока нет бронирований
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        Мои бронирования
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Тип размещения</TableCell>
              <TableCell>Комната</TableCell>
              <TableCell>Даты проживания</TableCell>
              <TableCell>Статус</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {bookings.map((booking) => (
              <TableRow key={booking.id} hover>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {getAccommodationIcon(booking.type)}
                    <Typography variant="body2">
                      {booking.type === 'bed'
                        ? 'Койко-место'
                        : booking.type === 'room'
                          ? 'Комната'
                          : 'Квартира'}
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Typography variant="body2">
                    {booking.room_name}
                    {booking.bed_id && (
                      <Typography variant="caption" display="block" color="text.secondary">
                        Место {booking.bed_id}
                      </Typography>
                    )}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Box>
                    <Typography variant="body2">
                      {`${formatDate(booking.start_date)} - ${formatDate(booking.end_date)}`}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {Math.ceil(
                        (new Date(booking.end_date) - new Date(booking.start_date)) /
                        (1000 * 60 * 60 * 24)
                      )} дней
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip
                    size="small"
                    label={booking.status === 'confirmed' ? 'Подтверждено' : 'В обработке'}
                    color={booking.status === 'confirmed' ? 'success' : 'warning'}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export default BookingsList;
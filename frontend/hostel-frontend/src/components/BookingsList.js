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
  Box
} from '@mui/material';
import {
  Apartment as ApartmentIcon,
  Hotel as HotelIcon,
  SingleBed as SingleBedIcon
} from '@mui/icons-material';
import axios from '../api/axios';

const BookingsList = () => {
  const [bookings, setBookings] = useState([]);

  useEffect(() => {
    const fetchBookings = async () => {
        try {
            const response = await axios.get('/api/v1/bookings'); // Используем /api/v1
            if (Array.isArray(response.data.data)) { // Проверяем структуру ответа
                setBookings(response.data.data);
            } else {
                console.error('Unexpected data format:', response.data);
                setBookings([]);
            }
        } catch (error) {
            console.error('Error fetching bookings:', error);
            setBookings([]);
        }
    };

    fetchBookings();
}, []);


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

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        Список бронирований
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Тип размещения</TableCell>
              <TableCell>Комната</TableCell>
              <TableCell>Клиент</TableCell>
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
                  <Typography variant="body2">
                    {booking.user_name}
                    <Typography variant="caption" display="block" color="text.secondary">
                      {booking.user_email}
                    </Typography>
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
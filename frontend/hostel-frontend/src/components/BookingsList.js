// frontend/hostel-frontend/src/components/BookingsList.js
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
  Box
} from '@mui/material';
import axios from "../api/axios";

const BookingsList = () => {
  const [bookings, setBookings] = useState([]);
  const [rooms, setRooms] = useState({});
  const [users, setUsers] = useState({});

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Получаем все бронирования
        const bookingsResponse = await axios.get('/bookings');
        setBookings(bookingsResponse.data);

        // Получаем информацию о комнатах
        const roomsResponse = await axios.get('/rooms');
        const roomsMap = {};
        roomsResponse.data.forEach(room => {
          roomsMap[room.id] = room;
        });
        setRooms(roomsMap);

        // Получаем информацию о пользователях
        const usersResponse = await axios.get('/users');
        const usersMap = {};
        usersResponse.data.forEach(user => {
          usersMap[user.id] = user;
        });
        setUsers(usersMap);
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h5" gutterBottom>
        Список бронирований
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Пользователь</TableCell>
              <TableCell>Комната</TableCell>
              <TableCell>Тип размещения</TableCell>
              <TableCell>Дата заезда</TableCell>
              <TableCell>Дата выезда</TableCell>
              <TableCell>Статус</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {bookings.map((booking) => (
              <TableRow key={booking.id}>
                <TableCell>{booking.id}</TableCell>
                <TableCell>
                  {users[booking.user_id]?.name || `Пользователь ${booking.user_id}`}
                </TableCell>
                <TableCell>
                  {rooms[booking.room_id]?.name || `Комната ${booking.room_id}`}
                </TableCell>
                <TableCell>
                  {booking.bed_id ? 'Койко-место' : 'Комната целиком'}
                </TableCell>
                <TableCell>{booking.start_date}</TableCell>
                <TableCell>{booking.end_date}</TableCell>
                <TableCell>{booking.status || 'Подтверждено'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export default BookingsList;
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
  Box,
  IconButton,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle
} from '@mui/material';
import { Delete as DeleteIcon } from '@mui/icons-material';
import axios from "../api/axios";

const BookingsList = () => {
  const [bookings, setBookings] = useState([]);
  const [rooms, setRooms] = useState({});
  const [users, setUsers] = useState({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [bookingToDelete, setBookingToDelete] = useState(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const bookingsResponse = await axios.get('/bookings');
      setBookings(bookingsResponse.data);

      const roomsResponse = await axios.get('/rooms');
      const roomsMap = {};
      roomsResponse.data.forEach(room => {
        roomsMap[room.id] = room;
      });
      setRooms(roomsMap);

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

  const handleDeleteClick = (booking) => {
    setBookingToDelete(booking);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    try {
      if (bookingToDelete.bed_id) {
        await axios.delete(`/beds/${bookingToDelete.bed_id}/bookings/${bookingToDelete.id}`);
      } else {
        await axios.delete(`/rooms/${bookingToDelete.room_id}/bookings/${bookingToDelete.id}`);
      }
      setDeleteDialogOpen(false);
      setBookingToDelete(null);
      fetchData(); // Обновляем список после удаления
    } catch (error) {
      console.error('Ошибка удаления бронирования:', error);
      alert('Ошибка удаления бронирования');
    }
  };

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
              <TableCell>Действия</TableCell>
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
                <TableCell>
                  <IconButton 
                    color="error" 
                    onClick={() => handleDeleteClick(booking)}
                    title="Удалить бронирование"
                  >
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Диалог подтверждения удаления */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
      >
        <DialogTitle>Подтверждение удаления</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Вы действительно хотите удалить это бронирование?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleDeleteConfirm} color="error" autoFocus>
            Удалить
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default BookingsList;
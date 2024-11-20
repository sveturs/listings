import React, { useState, useEffect } from "react";
import {
  Container,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  MenuItem,
  Grid
} from "@mui/material";
import axios from "../api/axios";

const AddBooking = () => {
  const [booking, setBooking] = useState({
    user_id: "",
    room_id: "",
    start_date: "",
    end_date: ""
  });
  
  const [rooms, setRooms] = useState([]); // Список всех комнат
  const [users, setUsers] = useState([]); // Список всех пользователей
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  // Загрузка списка комнат и пользователей
  useEffect(() => {
    const fetchData = async () => {
      try {
        const [roomsResponse, usersResponse] = await Promise.all([
          axios.get("/rooms"),
          axios.get("/users")
        ]);
        setRooms(roomsResponse.data);
        setUsers(usersResponse.data);
      } catch (error) {
        console.error("Ошибка загрузки данных:", error);
      }
    };
    fetchData();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    try {
      await axios.post("/bookings", {
        ...booking,
        user_id: parseInt(booking.user_id),
        room_id: parseInt(booking.room_id)
      });
      setSuccess(true);
      setBooking({
        user_id: "",
        room_id: "",
        start_date: "",
        end_date: ""
      });
    } catch (error) {
      setError(error.response?.data || "Ошибка добавления бронирования");
      console.error("Ошибка добавления бронирования:", error);
    }
  };

  // Получаем текущую дату для ограничения выбора дат
  const today = new Date().toISOString().split('T')[0];

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Добавить бронирование (Админ)
        </Typography>

        {success && (
          <Alert severity="success" sx={{ mb: 2 }}>
            Бронирование добавлено успешно!
          </Alert>
        )}

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <form onSubmit={handleSubmit}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                select
                label="Пользователь"
                fullWidth
                value={booking.user_id}
                onChange={(e) => setBooking({ ...booking, user_id: e.target.value })}
              >
                {users.map((user) => (
                  <MenuItem key={user.id} value={user.id}>
                    {user.name} ({user.email})
                  </MenuItem>
                ))}
              </TextField>
            </Grid>

            <Grid item xs={12}>
              <TextField
                select
                label="Комната"
                fullWidth
                value={booking.room_id}
                onChange={(e) => setBooking({ ...booking, room_id: e.target.value })}
              >
                {rooms.map((room) => (
                  <MenuItem key={room.id} value={room.id}>
                    {room.name} ({room.address_city}, {room.price_per_night} евро/сутки)
                  </MenuItem>
                ))}
              </TextField>
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Дата заезда"
                type="date"
                fullWidth
                InputLabelProps={{ shrink: true }}
                value={booking.start_date}
                onChange={(e) => setBooking({ ...booking, start_date: e.target.value })}
                inputProps={{ min: today }}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Дата выезда"
                type="date"
                fullWidth
                InputLabelProps={{ shrink: true }}
                value={booking.end_date}
                onChange={(e) => setBooking({ ...booking, end_date: e.target.value })}
                inputProps={{ min: booking.start_date || today }}
              />
            </Grid>

            <Grid item xs={12}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                fullWidth
                size="large"
              >
                Добавить бронирование
              </Button>
            </Grid>
          </Grid>
        </form>
      </Box>
    </Container>
  );
};

export default AddBooking;
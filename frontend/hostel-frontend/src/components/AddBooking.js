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

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const AddBooking = () => {
  const [booking, setBooking] = useState({
    user_id: "",
    room_id: "",
    start_date: "",
    end_date: ""
  });
  
  const [selectedRoom, setSelectedRoom] = useState(null);
  const [availableBeds, setAvailableBeds] = useState([]);
  const [bedImages, setBedImages] = useState({});
  const [selectedBed, setSelectedBed] = useState(null);
  const [rooms, setRooms] = useState([]);
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  // Загрузка списка комнат и пользователей
  useEffect(() => {
    const fetchData = async () => {
      try {
        const [roomsResponse, usersResponse] = await Promise.all([
          axios.get("/api/v1/rooms"),
          axios.get("/api/v1/users")
        ]);
        setRooms(roomsResponse.data.data || []);
        setUsers(usersResponse.data.data || []);
      } catch (error) {
        console.error("Ошибка загрузки данных:", error);
        setError("Не удалось загрузить данные");
      }
    };
    fetchData();
  }, []);

  // Загрузка доступных кроватей при выборе дат
  useEffect(() => {
    if (selectedRoom?.accommodation_type === 'bed' && booking.start_date && booking.end_date) {
      axios.get(`/rooms/${selectedRoom.id}/available-beds`, {
        params: {
          start_date: booking.start_date,
          end_date: booking.end_date
        }
      })
      .then(response => {
        setAvailableBeds(response.data.data || []);
        // Загружаем изображения для кроватей
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
      .catch(error => {
        console.error('Ошибка загрузки койко-мест:', error);
        setError('Не удалось загрузить доступные койко-места');
      });
    }
  }, [selectedRoom, booking.start_date, booking.end_date]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (selectedRoom?.accommodation_type === 'bed' && !selectedBed) {
      setError("Выберите койко-место");
      return;
    }

    try {
      const bookingData = {
        ...booking,
        user_id: parseInt(booking.user_id),
        room_id: parseInt(booking.room_id)
      };

      if (selectedRoom?.accommodation_type === 'bed') {
        bookingData.bed_id = selectedBed;
      }

      await axios.post("/api/v1/bookings", bookingData);
      setSuccess(true);
      // Сброс формы
      setBooking({
        user_id: "",
        room_id: "",
        start_date: "",
        end_date: ""
      });
      setSelectedRoom(null);
      setSelectedBed(null);
    } catch (error) {
      setError(error.response?.data?.error || "Ошибка добавления бронирования");
      console.error("Ошибка добавления бронирования:", error);
    }
  };

  // Получаем текущую дату для ограничения выбора дат
  const today = new Date().toISOString().split('T')[0];

  // Расчет общей стоимости
  const calculateTotalPrice = () => {
    if (!booking.start_date || !booking.end_date) return 0;

    const start = new Date(booking.start_date);
    const end = new Date(booking.end_date);
    const daysCount = Math.ceil((end - start) / (1000 * 60 * 60 * 24));

    let pricePerNight;
    if (selectedRoom?.accommodation_type === 'bed' && selectedBed) {
      const selectedBedData = availableBeds.find(bed => bed.id === selectedBed);
      pricePerNight = selectedBedData ? selectedBedData.price_per_night : 0;
    } else {
      pricePerNight = selectedRoom?.price_per_night || 0;
    }

    return pricePerNight * daysCount;
  };

  return (
    <Container maxWidth="md">
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
                required
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
                required
                value={booking.room_id}
                onChange={(e) => {
                  const room = rooms.find(r => r.id === parseInt(e.target.value));
                  setSelectedRoom(room);
                  setBooking({ ...booking, room_id: e.target.value });
                  setSelectedBed(null);
                }}
              >
                {rooms.map((room) => (
                  <MenuItem key={room.id} value={room.id}>
                    {room.name} ({room.address_city}, {room.price_per_night} ₽/сутки)
                  </MenuItem>
                ))}
              </TextField>
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Дата заезда"
                type="date"
                fullWidth
                required
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
                required
                InputLabelProps={{ shrink: true }}
                value={booking.end_date}
                onChange={(e) => setBooking({ ...booking, end_date: e.target.value })}
                inputProps={{ min: booking.start_date || today }}
              />
            </Grid>

            {selectedRoom?.accommodation_type === 'bed' && booking.start_date && booking.end_date && (
              <Grid item xs={12}>
                <Typography variant="subtitle1" sx={{ mb: 2 }}>
                  Выберите койко-место:
                </Typography>
                <Grid container spacing={2}>
                  {availableBeds.map(bed => (
                    <Grid item xs={12} sm={6} md={4} key={bed.id}>
                      <Box
                        sx={{
                          p: 2,
                          border: '1px solid',
                          borderColor: selectedBed === bed.id ? 'primary.main' : 'divider',
                          borderRadius: 1,
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                          bgcolor: selectedBed === bed.id ? 'action.selected' : 'background.paper',
                          '&:hover': {
                            borderColor: 'primary.main',
                            boxShadow: 1
                          },
                          height: '100%',
                          display: 'flex',
                          flexDirection: 'column'
                        }}
                        onClick={() => setSelectedBed(bed.id)}
                      >
                        {bedImages[bed.id]?.length > 0 && (
                          <Box sx={{ mb: 1, width: '100%', height: 140 }}>
                            <img
                              src={`${BACKEND_URL}/uploads/${bedImages[bed.id][0].file_path}`}
                              alt={`Койко-место ${bed.bed_number}`}
                              style={{
                                width: '100%',
                                height: '100%',
                                objectFit: 'cover',
                                borderRadius: '4px'
                              }}
                            />
                          </Box>
                        )}
                        <Typography variant="h6" component="div">
                          Место {bed.bed_number}
                        </Typography>
                        <Typography color="primary.main" variant="h6" sx={{ mt: 'auto' }}>
                          {bed.price_per_night} ₽/ночь
                        </Typography>
                      </Box>
                    </Grid>
                  ))}
                </Grid>
              </Grid>
            )}

            {(booking.start_date && booking.end_date && selectedRoom) && (
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mt: 2 }}>
                  Итого к оплате: {calculateTotalPrice()} ₽
                </Typography>
              </Grid>
            )}

            <Grid item xs={12}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                fullWidth
                size="large"
                disabled={!booking.user_id || !booking.room_id || !booking.start_date || !booking.end_date || 
                         (selectedRoom?.accommodation_type === 'bed' && !selectedBed)}
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
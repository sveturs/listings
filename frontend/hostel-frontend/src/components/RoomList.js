import React, { useState, useEffect, useCallback } from "react";
import axios from "../api/axios";

const RoomList = () => {
  const [rooms, setRooms] = useState([]);
  const [filters, setFilters] = useState({ capacity: "", minPrice: "", maxPrice: "" });
  const [bookings, setBookings] = useState([]);

  const fetchRooms = useCallback(async () => {
    try {
      const params = new URLSearchParams();
      if (filters.capacity) params.append("capacity", filters.capacity);
      if (filters.minPrice) params.append("min_price", filters.minPrice);
      if (filters.maxPrice) params.append("max_price", filters.maxPrice);
  
      const response = await axios.get(`/rooms?${params.toString()}`);
      setRooms(response.data || []);
    } catch (error) {
      console.error("Ошибка при получении списка комнат:", error);
      alert("Ошибка при загрузке комнат.");
    }
  }, [filters]);

  useEffect(() => {
    fetchRooms();
  }, [fetchRooms]);

  const fetchAllRooms = async () => {
    try {
      const response = await axios.get("/rooms");
      setRooms(response.data || []);
    } catch (error) {
      console.error("Ошибка при получении всех комнат:", error);
      alert("Не удалось получить список комнат.");
    }
  };

  const showBookings = async () => {
    try {
      const response = await axios.get("/bookings");
      setBookings(response.data || []);
      if (response.data.length === 0) {
        alert("Брони не найдены.");
      } else {
        alert(
          response.data
            .map(
              (booking) =>
                `Бронь ID: ${booking.id}, Комната ID: ${booking.room_id}, Пользователь ID: ${booking.user_id}, Даты: ${booking.start_date} - ${booking.end_date}`
            )
            .join("\n")
        );
      }
    } catch (error) {
      console.error("Ошибка при получении списка бронирований:", error);
      alert("Не удалось получить список бронирований.");
    }
  };

  const createRoom = async () => {
    try {
      const newRoom = {
        name: "New Room",
        capacity: 2,
        price_per_night: 200,
      };
      await axios.post("/rooms", newRoom);
      fetchAllRooms(); // Обновить список комнат
    } catch (error) {
      console.error("Ошибка при создании комнаты:", error);
      alert("Не удалось создать новую комнату.");
    }
  };

  const bookRoom = async (roomId) => {
    const userId = prompt("Введите ID пользователя:");
    const startDate = prompt("Введите дату начала (YYYY-MM-DD):");
    const endDate = prompt("Введите дату окончания (YYYY-MM-DD):");

    if (!userId || !startDate || !endDate) {
      alert("Все поля обязательны для заполнения!");
      return;
    }

    try {
      await axios.post("/bookings", {
        user_id: parseInt(userId, 10),
        room_id: roomId,
        start_date: startDate,
        end_date: endDate,
      });
      alert("Бронирование успешно добавлено!");
      showBookings(); // Обновить список бронирований
    } catch (error) {
      console.error("Ошибка при бронировании:", error);
      alert("Не удалось забронировать комнату. Проверьте данные.");
    }
  };

  return (
    <div>
      <h1>Список комнат</h1>
      <div>
        <label>Вместимость: </label>
        <input
          type="number"
          value={filters.capacity}
          onChange={(e) => setFilters({ ...filters, capacity: e.target.value })}
        />
        <label>Мин. цена: </label>
        <input
          type="number"
          value={filters.minPrice}
          onChange={(e) => setFilters({ ...filters, minPrice: e.target.value })}
        />
        <label>Макс. цена: </label>
        <input
          type="number"
          value={filters.maxPrice}
          onChange={(e) => setFilters({ ...filters, maxPrice: e.target.value })}
        />
        <button onClick={fetchRooms}>Фильтровать</button>
        <button onClick={fetchAllRooms}>Показать все комнаты</button>
        <button onClick={showBookings}>Показать все брони</button>
        <button onClick={createRoom}>Создать новую комнату</button>
      </div>
      <ul>
        {rooms.length > 0 ? (
          rooms.map((room) => (
            <li key={room.id}>
              {room.name} - {room.capacity} мест - {room.price_per_night} руб./ночь
              <button onClick={() => bookRoom(room.id)}>Забронировать</button>
            </li>
          ))
        ) : (
          <p>Комнаты не найдены.</p>
        )}
      </ul>
      <h2>Список бронирований</h2>
      <ul>
        {bookings.length > 0 ? (
          bookings.map((booking) => (
            <li key={booking.id}>
              Бронь ID: {booking.id}, Комната ID: {booking.room_id}, Пользователь ID: {booking.user_id}, Даты: {booking.start_date} - {booking.end_date}
            </li>
          ))
        ) : (
          <p>Брони отсутствуют.</p>
        )}
      </ul>
    </div>
  );
};

export default RoomList;

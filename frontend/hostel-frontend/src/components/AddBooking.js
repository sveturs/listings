import React, { useState } from "react";
import axios from "../api/axios";

const AddBooking = () => {
  const [booking, setBooking] = useState({
    user_id: "",
    room_id: "",
    start_date: "",
    end_date: "",
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post("/bookings", booking);
      alert("Бронирование добавлено успешно!");
    } catch (error) {
      console.error("Ошибка добавления бронирования:", error);
      alert("Не удалось добавить бронирование. Проверьте данные и повторите попытку.");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h1>Добавить бронирование</h1>
      <label>Пользователь ID:</label>
      <input
        type="number"
        value={booking.user_id}
        onChange={(e) => setBooking({ ...booking, user_id: e.target.value })}
      />
      <label>Комната ID:</label>
      <input
        type="number"
        value={booking.room_id}
        onChange={(e) => setBooking({ ...booking, room_id: e.target.value })}
      />
      <label>Дата начала:</label>
      <input
        type="date"
        value={booking.start_date}
        onChange={(e) => setBooking({ ...booking, start_date: e.target.value })}
      />
      <label>Дата окончания:</label>
      <input
        type="date"
        value={booking.end_date}
        onChange={(e) => setBooking({ ...booking, end_date: e.target.value })}
      />
      <button type="submit">Добавить бронирование</button>
    </form>
  );
};

export default AddBooking;

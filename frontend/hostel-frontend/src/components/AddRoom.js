import React, { useState } from "react";
import axios from "axios";

const AddRoom = () => {
  const [room, setRoom] = useState({ name: "", capacity: 0, price_per_night: 0 });

  const handleSubmit = async (e) => {
    e.preventDefault();
    await axios.post("http://localhost:3000/rooms", room);
    alert("Комната добавлена!");
  };

  return (
    <form onSubmit={handleSubmit}>
      <h1>Добавить комнату</h1>
      <label>Название:</label>
      <input
        type="text"
        value={room.name}
        onChange={(e) => setRoom({ ...room, name: e.target.value })}
      />
      <label>Вместимость:</label>
      <input
        type="number"
        value={room.capacity}
        onChange={(e) => setRoom({ ...room, capacity: +e.target.value })}
      />
      <label>Цена за ночь:</label>
      <input
        type="number"
        value={room.price_per_night}
        onChange={(e) => setRoom({ ...room, price_per_night: +e.target.value })}
      />
      <button type="submit">Добавить</button>
    </form>
  );
};

export default AddRoom;

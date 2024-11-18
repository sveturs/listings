import React, { useState, useEffect, useCallback } from "react";
import axios from "../api/axios";

const RoomList = () => {
  const [rooms, setRooms] = useState([]);
  const [filters, setFilters] = useState({ capacity: "", minPrice: "", maxPrice: "" });

  // Используем useCallback для fetchRooms
  const fetchRooms = useCallback(async () => {
    try {
      const params = new URLSearchParams(filters).toString();
      const response = await axios.get(`/rooms?${params}`);
      setRooms(response.data);
    } catch (error) {
      console.error("Ошибка при получении списка комнат:", error);
    }
  }, [filters]); // Завиимость: filters

  // useEffect вызывает fetchRooms
  useEffect(() => {
    fetchRooms();
  }, [fetchRooms]);

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
      </div>
      <ul>
        {rooms.map((room) => (
          <li key={room.id}>
            {room.name} - {room.capacity} мест - {room.price_per_night} руб./ночь
          </li>
        ))}
      </ul>
    </div>
  );
};

export default RoomList;

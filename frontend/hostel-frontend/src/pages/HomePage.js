import React from "react";
import { Link } from "react-router-dom";
import RoomList from "../components/RoomList";

const HomePage = () => (
  <div>
    <h1>Главная</h1>
    <nav>
      <ul>
        <li><Link to="/add-room">Добавить комнату</Link></li>
        <li><Link to="/add-user">Добавить пользователя</Link></li>
        <li><Link to="/admin">Админская панель</Link></li>
      </ul>
    </nav>
    <RoomList />
  </div>
);

export default HomePage;

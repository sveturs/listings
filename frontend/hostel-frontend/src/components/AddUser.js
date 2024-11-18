import React, { useState } from "react";
import axios from "axios";

const AddUser = () => {
  const [user, setUser] = useState({ name: "", email: "" });

  const handleSubmit = async (e) => {
    e.preventDefault();
    await axios.post("http://localhost:3000/users", user);
    alert("Пользователь добавлен!");
  };

  return (
    <form onSubmit={handleSubmit}>
      <h1>Добавить пользователя</h1>
      <label>Имя:</label>
      <input
        type="text"
        value={user.name}
        onChange={(e) => setUser({ ...user, name: e.target.value })}
      />
      <label>Email:</label>
      <input
        type="email"
        value={user.email}
        onChange={(e) => setUser({ ...user, email: e.target.value })}
      />
      <button type="submit">Добавить</button>
    </form>
  );
};

export default AddUser;

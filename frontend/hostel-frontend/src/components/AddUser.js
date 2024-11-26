import React, { useState } from "react";
import { TextField, Button, Container, Typography } from "@mui/material";
import axios from "../api/axios";

const AddUser = () => {
  const [user, setUser] = useState({ name: "", email: "" });

  const handleSubmit = async (e) => {
    e.preventDefault();
    await axios.post('/users', user);
    alert("Пользователь добавлен!");
  };

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        Добавить пользователя
      </Typography>
      <form onSubmit={handleSubmit}>
        <TextField
          label="Имя"
          fullWidth
          margin="normal"
          value={user.name}
          onChange={(e) => setUser({ ...user, name: e.target.value })}
        />
        <TextField
          label="Email"
          type="email"
          fullWidth
          margin="normal"
          value={user.email}
          onChange={(e) => setUser({ ...user, email: e.target.value })}
        />
        <Button type="submit" variant="contained" color="primary">
          Добавить
        </Button>
      </form>
    </Container>
  );
};

export default AddUser;

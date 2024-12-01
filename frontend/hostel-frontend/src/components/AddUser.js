import React, { useState } from "react";
import { TextField, Button, Container, Typography, Alert } from "@mui/material";
import axios from "../api/axios";
import { useAuth } from "../contexts/AuthContext";

const AddUser = () => {
  const { user } = useAuth();
  const [userForm, setUserForm] = useState({ name: "", email: "" });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (!user) {
      setError("Необходимо войти в систему");
      return;
    }

    try {
      await axios.post('/users', userForm, {
        withCredentials: true
      });
      setSuccess(true);
      setUserForm({ name: "", email: "" });
    } catch (err) {
      setError(err.response?.data || "Ошибка при добавлении пользователя");
    }
  };

  if (!user) {
    return (
      <Container>
        <Alert severity="warning">
          Для добавления пользователей необходимо войти в систему
        </Alert>
      </Container>
    );
  }

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        Добавить пользователя
      </Typography>
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      {success && <Alert severity="success" sx={{ mb: 2 }}>Пользователь успешно добавлен!</Alert>}
      <form onSubmit={handleSubmit}>
        <TextField
          label="Имя"
          fullWidth
          margin="normal"
          value={userForm.name}
          onChange={(e) => setUserForm({ ...userForm, name: e.target.value })}
          required
        />
        <TextField
          label="Email"
          type="email"
          fullWidth
          margin="normal"
          value={userForm.email}
          onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
          required
        />
        <Button 
          type="submit" 
          variant="contained" 
          color="primary"
          sx={{ mt: 2 }}
        >
          Добавить
        </Button>
      </form>
    </Container>
  );
};

export default AddUser;
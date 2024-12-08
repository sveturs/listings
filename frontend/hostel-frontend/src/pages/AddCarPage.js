import React, { useState } from "react";
import { Container, TextField, Button, Box, Alert } from "@mui/material";

const AddCarPage = () => {
  const [formData, setFormData] = useState({
    make: "",
    model: "",
    year: "",
    pricePerDay: "",
    location: "",
  });

  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch("/api/cars", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...formData, availability: true }),
    })
      .then((response) => {
        if (!response.ok) {
          // Обработка ошибок сервера
          return response.json().then((data) => {
            throw new Error(data.error || "Ошибка сервера");
          });
        }
        return response.json(); // Парсим ответ как JSON
      })
      .then((data) => {
        setSuccessMessage("Автомобиль успешно добавлен!");
        setFormData({ make: "", model: "", year: "", pricePerDay: "", location: "" });
      })
      .catch((error) => {
        setErrorMessage(error.message);
      });
  };
  

  return (
    <Container sx={{ marginTop: 4 }}>
      {successMessage && <Alert severity="success">{successMessage}</Alert>}
      {errorMessage && <Alert severity="error">{errorMessage}</Alert>}
      <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
        <TextField
          label="Марка"
          name="make"
          value={formData.make}
          onChange={handleChange}
          fullWidth
          margin="normal"
          required
        />
        <TextField
          label="Модель"
          name="model"
          value={formData.model}
          onChange={handleChange}
          fullWidth
          margin="normal"
          required
        />
        <TextField
          label="Год"
          name="year"
          value={formData.year}
          onChange={handleChange}
          fullWidth
          margin="normal"
          required
          type="number"
        />
        <TextField
          label="Цена в день"
          name="pricePerDay"
          value={formData.pricePerDay}
          onChange={handleChange}
          fullWidth
          margin="normal"
          required
          type="number"
        />
        <TextField
          label="Локация"
          name="location"
          value={formData.location}
          onChange={handleChange}
          fullWidth
          margin="normal"
          required
        />
        <Button type="submit" variant="contained" color="primary" sx={{ mt: 2 }}>
          Добавить автомобиль
        </Button>
      </Box>
    </Container>
  );
};

export default AddCarPage;

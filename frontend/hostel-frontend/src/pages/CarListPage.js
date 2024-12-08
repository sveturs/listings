import React, { useState, useEffect } from "react";
import { Container, TextField, Grid, Card, CardContent, Typography } from "@mui/material";

const CarListPage = () => {
  const [cars, setCars] = useState([]);
  const [filter, setFilter] = useState("");

  useEffect(() => {
    // Загрузка списка автомобилей
    fetch("/api/cars/available")
      .then((response) => response.json())
      .then((data) => setCars(data))
      .catch((error) => console.error("Error fetching cars:", error));
  }, []);

  const filteredCars = cars.filter((car) =>
    car.make.toLowerCase().includes(filter.toLowerCase()) ||
    car.model.toLowerCase().includes(filter.toLowerCase()) ||
    car.location.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <Container sx={{ marginTop: 4 }}>
      <TextField
        label="Фильтр по марке, модели или локации"
        variant="outlined"
        fullWidth
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
        sx={{ marginBottom: 4 }}
      />
      <Grid container spacing={3}>
        {filteredCars.map((car) => (
          <Grid item xs={12} sm={6} md={4} key={car.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{`${car.make} ${car.model}`}</Typography>
                <Typography variant="body2">Год: {car.year}</Typography>
                <Typography variant="body2">Цена в день: ${car.pricePerDay}</Typography>
                <Typography variant="body2">Локация: {car.location}</Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
};

export default CarListPage;

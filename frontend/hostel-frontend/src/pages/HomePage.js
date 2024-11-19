import React from "react";
import { Link } from "react-router-dom";
import { AppBar, Toolbar, Typography, Button, Container } from "@mui/material";
import RoomList from "../components/RoomList";

const HomePage = () => (
  <div>
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          Hostel Booking System
        </Typography>
        <Button color="inherit" component={Link} to="/add-room">
          Добавить комнату
        </Button>
        <Button color="inherit" component={Link} to="/add-user">
          Добавить пользователя
        </Button>
        <Button color="inherit" component={Link} to="/admin">
          Админская панель
        </Button>
      </Toolbar>
    </AppBar>
    <Container sx={{ marginTop: 4 }}>
      <Typography variant="h4" gutterBottom>
        Список комнат
      </Typography>
      <RoomList />
    </Container>
  </div>
);

export default HomePage;

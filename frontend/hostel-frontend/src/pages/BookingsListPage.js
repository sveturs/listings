// frontend/hostel-frontend/src/pages/BookingsListPage.js
import React from 'react';
import { Link } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button, Container } from '@mui/material';
import BookingsList from '../components/BookingsList';

const BookingsListPage = () => {
  return (
    <div>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            Hostel Booking System
          </Typography>
          <Button color="inherit" component={Link} to="/">
            На главную
          </Button>
        </Toolbar>
      </AppBar>
      <Container sx={{ mt: 4 }}>
        <BookingsList />
      </Container>
    </div>
  );
};

export default BookingsListPage;
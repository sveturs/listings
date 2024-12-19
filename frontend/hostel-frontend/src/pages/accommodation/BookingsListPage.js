// frontend/hostel-frontend/src/pages/accommodation/BookingsListPage.js
import React from 'react';
import { Container } from '@mui/material';
import BookingsList from '../../components/accommodation/BookingsList';

const BookingsListPage = () => {
  return (
    <Container sx={{ mt: 4 }}>
      <BookingsList />
    </Container>
  );
};

export default BookingsListPage;
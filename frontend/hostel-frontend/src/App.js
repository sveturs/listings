// frontend/hostel-frontend/src/App.js
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Container } from '@mui/material';
import HomePage from "./pages/HomePage";
import AddRoomPage from "./pages/AddRoomPage";
import AddUserPage from "./pages/AddUserPage";
import BookingsListPage from "./pages/BookingsListPage";
import AdminPanelPage from "./pages/AdminPanelPage";

function App() {
  return (
    <Router>
      <Container>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/bookings" element={<BookingsListPage />} />
          <Route path="/add-room" element={<AddRoomPage />} />
          <Route path="/add-user" element={<AddUserPage />} />
          <Route path="/admin" element={<AdminPanelPage />} />
        </Routes>
      </Container>
    </Router>
  );
}

export default App;
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import AddRoomPage from "./pages/AddRoomPage";
import AddUserPage from "./pages/AddUserPage";
import BookRoomPage from "./pages/BookRoomPage";
import AdminPanelPage from "./pages/AdminPanelPage";
import RoomList from "./components/RoomList";
import AddBookingPage from "./pages/AddBookingPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/add-room" element={<AddRoomPage />} />
        <Route path="/add-user" element={<AddUserPage />} />
        <Route path="/book-room" element={<BookRoomPage />} />
        <Route path="/admin" element={<AdminPanelPage />} />
        <Route path="/add-booking" element={<AddBookingPage />} />
      </Routes>
    </Router>
  );
}

export default App;

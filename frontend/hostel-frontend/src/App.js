import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./contexts/AuthContext";
import Layout from "./components/Layout";
import HomePage from "./pages/HomePage";
import AddRoomPage from "./pages/AddRoomPage";
import AddUserPage from "./pages/AddUserPage";
import BookingsListPage from "./pages/BookingsListPage";
import AdminPanelPage from "./pages/AdminPanelPage";
import PrivacyPolicy from "./pages/PrivacyPolicy";

function App() {
  return (
    <AuthProvider>
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/bookings" element={<BookingsListPage />} />
            <Route path="/add-room" element={<AddRoomPage />} />
            <Route path="/add-user" element={<AddUserPage />} />
            <Route path="/admin" element={<AdminPanelPage />} />
            <Route path="/privacy-policy" element={<PrivacyPolicy />} />
          </Routes>
        </Layout>
      </Router>
    </AuthProvider>
  );
}

export default App;
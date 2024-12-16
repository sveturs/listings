//frontend/hostel-frontend/src/App.js
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { future } from "@remix-run/router";


import { AuthProvider } from "./contexts/AuthContext";
import Layout from "./components/Layout";
import HomePage from "./pages/HomePage";
import AddRoomPage from "./pages/AddRoomPage";
import AddUserPage from "./pages/AddUserPage";
import BookingsListPage from "./pages/BookingsListPage";
import AdminPanelPage from "./pages/AdminPanelPage";
import PrivacyPolicy from "./pages/PrivacyPolicy";
import CarListPage from "./pages/CarListPage";
import AddCarPage from "./pages/AddCarPage";
import MarketplacePage from "./pages/MarketplacePage";
import CreateListingPage from "./pages/CreateListingPage";
import ListingDetailsPage from './pages/ListingDetailsPage';
import UserProfile from './components/UserProfile';
import MyListingsPage from './pages/MyListingsPage';


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
            <Route path="/cars" element={<CarListPage />} /> {/* Список автомобилей */}
            <Route path="/add-car" element={<AddCarPage />} /> {/* Добавление автомобиля */}
            <Route path="/marketplace" element={<MarketplacePage />} />
            <Route path="/marketplace/create" element={<CreateListingPage />} />
            <Route path="/marketplace/listings/:id" element={<ListingDetailsPage />} />
            <Route path="/profile" element={<UserProfile />} />
            <Route path="/marketplace" element={<MarketplacePage />} />
            
            
            <Route path="/my-listings" element={<MyListingsPage />} />
            

          
          </Routes>
        </Layout>
      </Router>
    </AuthProvider>
  );
}

export default App;
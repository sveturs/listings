//frontend/hostel-frontend/src/App.js
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { future } from "@remix-run/router";


import { AuthProvider } from "./contexts/AuthContext";
import Layout from "./components/global/Layout";
import HomePage from "./pages/global/HomePage";
import AddRoomPage from "./pages/accommodation/AddRoomPage";
import AddUserPage from "./pages/user/AddUserPage";
import BookingsListPage from "./pages/accommodation/BookingsListPage";
import AdminPanelPage from "./pages/global/AdminPanelPage";
import PrivacyPolicy from "./pages/accommodation/PrivacyPolicy";
import CarListPage from "./pages/car/CarListPage";
import AddCarPage from "./pages/car/AddCarPage";
import MarketplacePage from "./pages/marketplace/MarketplacePage";
import CreateListingPage from "./pages/marketplace/CreateListingPage";
import ListingDetailsPage from './pages/marketplace/ListingDetailsPage';
import UserProfile from './components/user/UserProfile';
import MyListingsPage from './pages/marketplace/MyListingsPage';
import FavoriteListingsPage from './pages/marketplace/FavoriteListingsPage';
import { MapProvider } from './components/maps/MapProvider';
 

function App() {
  return (
    <MapProvider>
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
              <Route path="/favorites" element={<FavoriteListingsPage />} />

            </Routes>
          </Layout>
        </Router>
      </AuthProvider>
    </MapProvider>
  );
}

export default App;
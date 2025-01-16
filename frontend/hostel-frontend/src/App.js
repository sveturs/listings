import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { AuthProvider } from "./contexts/AuthContext";
import Layout from "./components/global/Layout";
import HomePage from "./pages/global/HomePage";
import AddUserPage from "./pages/user/AddUserPage";
import AdminPanelPage from "./pages/global/AdminPanelPage";
import PrivacyPolicy from "./pages/accommodation/PrivacyPolicy";
import MarketplacePage from "./pages/marketplace/MarketplacePage";
import CreateListingPage from "./pages/marketplace/CreateListingPage";
import ListingDetailsPage from './pages/marketplace/ListingDetailsPage';
import UserProfile from './components/user/UserProfile';
import MyListingsPage from './pages/marketplace/MyListingsPage';
import FavoriteListingsPage from './pages/marketplace/FavoriteListingsPage';
import { MapProvider } from './components/maps/MapProvider';
import ChatPage from "./pages/marketplace/ChatPage";

function App() {
  return (
    <BrowserRouter
      future={{
        v7_startTransition: true,
        v7_relativeSplatPath: true
      }}
    >
      <MapProvider>
        <AuthProvider>
          <Layout>
            <Routes>
            <Route path="/" element={<MarketplacePage />} />
              <Route path="/add-user" element={<AddUserPage />} />
              <Route path="/admin" element={<AdminPanelPage />} />
              <Route path="/privacy-policy" element={<PrivacyPolicy />} />
              <Route path="/marketplace" element={<MarketplacePage />} />
              <Route path="/marketplace/create" element={<CreateListingPage />} />
              <Route path="/marketplace/listings/:id" element={<ListingDetailsPage />} />
              <Route path="/profile" element={<UserProfile />} />
              <Route path="/marketplace/chat" element={<ChatPage />} />
              <Route path="/my-listings" element={<MyListingsPage />} />
              <Route path="/favorites" element={<FavoriteListingsPage />} />
            </Routes>
          </Layout>
        </AuthProvider>
      </MapProvider>
    </BrowserRouter>
  );
}

export default App;
import React, { Suspense } from 'react';
import 'leaflet/dist/leaflet.css';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from 'react-query';
import { AuthProvider } from "./contexts/AuthContext";
import Layout from "./components/global/Layout";
import HomePage from "./pages/global/HomePage";
import PrivacyPolicy from "./pages/accommodation/PrivacyPolicy";
import MarketplacePage from "./pages/marketplace/MarketplacePage";
import CreateListingPage from "./pages/marketplace/CreateListingPage";
import ListingDetailsPage from './pages/marketplace/ListingDetailsPage';
import UserProfile from './components/user/UserProfile.tsx';
import MyListingsPage from './pages/marketplace/MyListingsPage';
import FavoriteListingsPage from './pages/marketplace/FavoriteListingsPage';
import { MapProvider } from './components/maps/MapProvider';
import ChatPage from "./pages/marketplace/ChatPage";
import { LanguageProvider } from "./contexts/LanguageContext";
import EditListingPage from './pages/marketplace/EditListingPage';
import { ChatProvider } from './contexts/ChatContext';
import PrivateRoute from "./components/global/PrivateRoute.tsx";
import { NotificationProvider } from './contexts/NotificationContext';
import NotificationSettings from './components/notifications/NotificationSettings';
import i18n from './i18n/config';
import './i18n/config';
import { CircularProgress, Box } from '@mui/material';
import TransactionsPage from './pages/balance/TransactionsPage';
import StorefrontPage from "./pages/store/StorefrontPage";
import StorefrontDetailPage from "./pages/store/StorefrontDetailPage";
import PublicStorefrontPage from "./pages/store/PublicStorefrontPage";
import UserReviewsPage from './pages/user/UserReviewsPage';
import StorefrontReviewsPage from './pages/store/StorefrontReviewsPage';
import EditStorefrontPage from "./pages/store/EditStorefrontPage";
import { LocationProvider } from './contexts/LocationContext';
import AdminRoute from "./components/global/AdminRoute.tsx";
import AdminPage from "./pages/admin/AdminPage";
import UsersManagementPage from "./pages/admin/UsersManagementPage";
import { GISMapPage } from './pages/gis';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      refetchOnMount: false,
      refetchOnReconnect: false,
      retry: 1,
      staleTime: 5 * 60 * 1000,
    },
  },
});

function App() {
  return (
    <BrowserRouter
      future={{
        v7_startTransition: true,
        v7_relativeSplatPath: true
      }}
    >
      <QueryClientProvider client={queryClient}>
        <MapProvider>
          <LanguageProvider>
            <AuthProvider>
              <ChatProvider>
                <LocationProvider>
                  <NotificationProvider>
                    <Layout>
                      <Routes>
                        <Route path="/" element={<MarketplacePage />} />
                        <Route path="/privacy-policy" element={<PrivacyPolicy />} />
                        <Route path="/marketplace" element={<MarketplacePage />} />
                        {/*  <Route path="/marketplace/create" element={<CreateListingPage />} />*/}
                        <Route path="/storefronts" element={
                          <PrivateRoute>
                            <StorefrontPage />
                          </PrivateRoute>
                        } />
                        <Route
                          path="/storefronts/:id/edit"
                          element={
                            <PrivateRoute>
                              <EditStorefrontPage />
                            </PrivateRoute>
                          }
                        />
                        {/*<Route path="/storefronts/:id" element={<StorefrontDetailPage />} />*/}
                        <Route path="/storefronts/:id" element={<PrivateRoute><StorefrontDetailPage /></PrivateRoute>} />
                        <Route path="/shop/:id" element={<PublicStorefrontPage />} />
                        <Route path="/user/:id/reviews" element={<UserReviewsPage />} />
                        <Route path="/shop/:id/reviews" element={<StorefrontReviewsPage />} />
                        <Route
                          path="/marketplace/create"
                          element={
                            <PrivateRoute>
                              <CreateListingPage />
                            </PrivateRoute>
                          }
                        />
                        <Route
                          path="/notifications/settings"
                          element={
                            <PrivateRoute>
                              <NotificationSettings />
                            </PrivateRoute>
                          }
                        />
                        <Route path="/marketplace/listings/:id" element={<ListingDetailsPage />} />
                        <Route path="/profile" element={<UserProfile />} />
                        <Route path="/marketplace/chat" element={<ChatPage />} />
                        <Route path="/my-listings" element={<MyListingsPage />} />
                        <Route path="/favorites" element={<FavoriteListingsPage />} />
                        <Route path="/marketplace/listings/:id/edit" element={<EditListingPage />} />
                        <Route path="/balance" element={<PrivateRoute>      <TransactionsPage />    </PrivateRoute>} />
                        <Route path="/gis" element={<GISMapPage />} />
                        <Route path="/admin" element={<AdminRoute>      <AdminPage />    </AdminRoute>} />
                        <Route path="/admin/users" element={<AdminRoute>      <UsersManagementPage />    </AdminRoute>} />
                      </Routes>
                    </Layout>
                  </NotificationProvider>
                </LocationProvider>
              </ChatProvider>
            </AuthProvider>

          </LanguageProvider>
        </MapProvider>
      </QueryClientProvider>
    </BrowserRouter >
  );
}

export default App;
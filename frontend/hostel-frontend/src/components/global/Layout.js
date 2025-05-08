// frontend/hostel-frontend/src/components/global/Layout.js
import React, { useState, useEffect } from "react";
import { Storefront } from '@mui/icons-material';
import NewMessageIndicator from './NewMessageIndicator.tsx';
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from '../../contexts/AuthContext';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { Bookmark } from '@mui/icons-material';
import UserProfile from "../user/UserProfile.tsx";
import { useChat } from '../../contexts/ChatContext';
import LanguageSwitcher from '../shared/LanguageSwitcher.tsx';
import NotificationDrawer from '../notifications/NotificationDrawer';
import { Settings } from '@mui/icons-material';
import { AccountBalanceWallet } from '@mui/icons-material';
import SveTuLogo from '../icons/SveTuLogo.tsx';
import { ShoppingBag, Store } from '@mui/icons-material';
import { useLocation as useGeoLocation } from '../../contexts/LocationContext'; // Добавляем импорт
import CitySelector from './CitySelector.tsx';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import CategoryMenu from '../marketplace/CategoryMenu'; // Импортируем компонент меню категорий
import { Plus } from 'lucide-react';
import {
  AppBar,
  Toolbar,
  Box,
  Container,
  Typography,
  IconButton,
  Avatar,
  Tooltip,
  Menu,
  MenuItem,
  Divider,
  useMediaQuery,
  useTheme,
  Modal,
  Alert,
  Slide,
  Button,
} from "@mui/material";
import {
  HomeWork,
  DirectionsCar,
  Key,
  Logout,
  ListAlt,
  AddHome,
  AccountCircle,
  Chat,
  Map,
} from "@mui/icons-material";

const Layout = ({ children }) => {
  const { userLocation, setCity, locationDismissed, dismissLocationSuggestion } = useGeoLocation();
  const [showLocationAlert, setShowLocationAlert] = useState(false);
  const { t, i18n } = useTranslation(['common', 'marketplace', 'gis']);
  const theme = useTheme();
  const navigate = useNavigate();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));
  const location = useLocation();
  const currentPath = location.pathname;
  const { user, login, logout } = useAuth();
  const [notificationDrawerOpen, setNotificationDrawerOpen] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);
  const [chatService, setChatService] = useState(null);
  const { getChatService } = useChat();
  const [anchorEl, setAnchorEl] = useState(null);
  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const [categoryPath] = useState([]);

  const handleOpenProfile = () => {
    setIsProfileOpen(true);
    handleCloseMenu();
  };

  const handleCloseProfile = () => {
    setIsProfileOpen(false);
  };

  const menuItems = [
    {
      path: "/",
      label: "Sve Tu",
      icon: <SveTuLogo width={isMobile ? 50 : 60} height={isMobile ? 50 : 60} />,
    }
  ];



  useEffect(() => {
    let unsubscribe;

    if (user?.id) {
      const chatService = getChatService(user.id);

      const messageHandler = (message) => {
        if (message.receiver_id === user.id && !message.is_read) {
          setUnreadCount(prev => prev + 1);
          const audio = new Audio('/notification.mp3');
          audio.play().catch(e => console.log('Notification sound error:', e));
        }
      };

      const fetchUnreadCount = async () => {
        try {
          const response = await axios.get('/api/v1/marketplace/chat/unread-count');
          setUnreadCount(response.data.data.count);
        } catch (error) {
          console.error('Error fetching unread count:', error);
        }
      };

      chatService.connect();
      unsubscribe = chatService.onMessage(messageHandler);
      fetchUnreadCount();
    }

    return () => {
      if (unsubscribe) {
        unsubscribe();
      }
    };
  }, [user?.id, getChatService]);

  const handleOpenMenu = (e) => {
    setAnchorEl(e.currentTarget);
  };

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  const renderMessagesMenuItem = () => (
    <MenuItem
      component={Link}
      to="/marketplace/chat"
      onClick={() => {
        handleCloseMenu();
        setUnreadCount(0);
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <NewMessageIndicator unreadCount={unreadCount} />
        <Chat size={20} />
        <Typography>{t('navigation.messages')}</Typography>
      </Box>
    </MenuItem>
  );
  useEffect(() => {
    if (userLocation && !locationDismissed) {
      setShowLocationAlert(true);
    }
  }, [userLocation, locationDismissed]);

  const handleCloseLocationAlert = () => {
    setShowLocationAlert(false);
    dismissLocationSuggestion();
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar
        position="static"
        sx={{
          bgcolor: "background.default",
          color: "text.primary",
          borderBottom: "1px solid #e0e0e0",
          boxShadow: "none",
        }}
      >
        <Container maxWidth="lg">
          <Toolbar
            disableGutters
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              minHeight: "56px",
              px: isMobile ? '4px' : 2,
            }}
          >
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                gap: isMobile ? 0.5 : 3,
              }}
            >
              {menuItems.map((item) => (
                <Box
                  key={item.path}
                  component={Link}
                  to={item.path}
                  sx={{
                    textDecoration: "none",
                    display: "flex",
                    alignItems: "center",
                    gap: 1.5,
                    color: currentPath === item.path ? "primary.main" : "text.secondary",
                    fontWeight: currentPath === item.path ? 600 : 400,
                    fontSize: "1rem",
                    transition: "color 0.3s ease, transform 0.3s ease",
                    "&:hover": {
                      color: "primary.main",
                      transform: "scale(1.05)",
                    },
                  }}
                >
                  {item.icon}
                  {!isMobile && (
                    <Typography
                      variant="h6"
                      sx={{
                        fontSize: "1.1rem",
                        fontWeight: "bold",
                      }}
                    >
                      {item.label}
                    </Typography>
                  )}
                </Box>
              ))}

              {/* Добавляем меню категорий */}
              {!isMobile && <CategoryMenu />}
            </Box>

            <Box sx={{ display: "flex", alignItems: "center", gap: isMobile ? 0.5 : 2 }}>
              <Button
                id="createAnnouncementButton"
                variant="contained"
                onClick={() => navigate('/marketplace/create')}
                startIcon={isMobile ? null : <Plus />}
                sx={{ mr: isMobile ? 0.5 : 1, minWidth: isMobile ? '40px' : 'auto', px: isMobile ? 1 : 2 }}
              >
                {isMobile ? <Plus /> : t('listings.create.title', { ns: 'marketplace' })}
              </Button>
              
              <Button
                variant="outlined"
                color="primary"
                onClick={() => navigate('/gis')}
                startIcon={<Map />}
                sx={{ mr: isMobile ? 0.5 : 2, display: { xs: 'none', sm: 'flex' } }}
              >
                {t('title', { ns: 'gis' })}
              </Button>

              {/* Добавляем меню категорий для мобильной версии */}
              {isMobile && <CategoryMenu />}
              
              {/* Кнопка GIS для мобильной версии */}
              {isMobile && (
                <Tooltip title={t('title', { ns: 'gis' })}>
                  <IconButton 
                    color="primary" 
                    onClick={() => navigate('/gis')}
                    sx={{ mr: 0.5 }}
                  >
                    <Map />
                  </IconButton>
                </Tooltip>
              )}

              {/* Используем обновленный CitySelector без передачи onCityChange */}
              <CitySelector isMobile={isMobile} />

              <LanguageSwitcher />
              {!user ? (
                <Tooltip title={t('auth.signIn')}>
                  <IconButton onClick={() => {
                    const returnUrl = window.location.pathname + window.location.search;
                    const encodedReturnUrl = encodeURIComponent(returnUrl);
                    login(`?returnTo=${encodedReturnUrl}`);
                  }} color="primary">
                    <Key />
                  </IconButton>
                </Tooltip>
              ) : (
                <>
                  {unreadCount > 0 && (
                    <IconButton
                      component={Link}
                      to="/marketplace/chat"
                      onClick={() => setUnreadCount(0)}
                    >
                      <NewMessageIndicator unreadCount={unreadCount} />
                    </IconButton>
                  )}

                  <Tooltip title={t('auth.profile')}>
                    <IconButton onClick={handleOpenMenu}>
                      <Avatar
                        src={user.pictureUrl}
                        alt={user.name}
                        sx={{ width: 32, height: 32 }}
                      />
                    </IconButton>
                  </Tooltip>

                  <Menu
                    anchorEl={anchorEl}
                    open={Boolean(anchorEl)}
                    onClose={handleCloseMenu}
                    PaperProps={{
                      sx: { mt: 1.5, width: 220 },
                    }}
                    transformOrigin={{ horizontal: "right", vertical: "top" }}
                    anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
                  >
                    <MenuItem
                      onClick={handleOpenProfile}
                      sx={{ display: "flex", flexDirection: "column", alignItems: "flex-start", gap: 0.5 }}
                    >
                      <Typography variant="subtitle2" noWrap sx={{ fontWeight: 600 }}>
                        {user.name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary" noWrap>
                        {user.email}
                      </Typography>
                    </MenuItem>
                    <Divider />

                    {renderMessagesMenuItem()}
                    <MenuItem
                      onClick={() => {
                        handleCloseMenu();
                        navigate('/notifications/settings');
                      }}
                    >
                      <Settings sx={{ mr: 1 }} />
                      {t('navigation.notifications')}
                    </MenuItem>
                    <MenuItem
                      component={Link}
                      to="/my-listings"
                      onClick={handleCloseMenu}
                    >
                      <ShoppingBag fontSize="small" sx={{ mr: 1 }} />
                      {t('navigation.myListings')}
                    </MenuItem>

                    <MenuItem component={Link} to="/storefronts" onClick={handleCloseMenu}>
                      <Store fontSize="small" sx={{ mr: 1 }} />
                      {t('navigation.storefronts', { defaultValue: 'Мои витрины' })}
                    </MenuItem>
                    <MenuItem component={Link} to="/favorites">
                      <Bookmark fontSize="small" sx={{ mr: 1 }} />
                      {t('navigation.favorites')}
                    </MenuItem>
                    <MenuItem component={Link} to="/balance" onClick={handleCloseMenu}>
                      <AccountBalanceWallet fontSize="small" sx={{ mr: 1 }} />
                      {t('navigation.balance')}
                    </MenuItem>
                    <MenuItem component={Link} to="/gis" onClick={handleCloseMenu}>
                      <Map fontSize="small" sx={{ mr: 1 }} />
                      {t('gis:title', { defaultValue: 'Map Search' })}
                    </MenuItem>

                    {user.email === 'voroshilovdo@gmail.com' && (
                      <MenuItem component={Link} to="/admin" onClick={handleCloseMenu}>
                        <Settings fontSize="small" sx={{ mr: 1 }} />
                        Администрирование
                      </MenuItem>
                    )}

                    <Divider />
                    <MenuItem onClick={logout}>
                      <Logout fontSize="small" sx={{ mr: 1 }} />
                      {t('auth.signOut')}
                    </MenuItem>
                  </Menu>
                </>
              )}
            </Box>
          </Toolbar>
        </Container>
      </AppBar>

      <Slide direction="down" in={showLocationAlert} mountOnEnter unmountOnExit>
        <Alert
          severity="info"
          sx={{
            mb: 2,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between'
          }}
          onClose={handleCloseLocationAlert}
          action={
            <Button
              color="inherit"
              size="small"
              onClick={handleCloseLocationAlert}
              sx={{ ml: 2 }}
            >
              {t('location.useThisCity', { defaultValue: 'Использовать этот город' })}
            </Button>
          }
        >
          {t('location.detectedCity', {
            defaultValue: 'Мы определили, что вы находитесь в городе {{city}}',
            city: userLocation?.city
          })}
        </Alert>
      </Slide>

      <Container maxWidth="lg" sx={{ py: 0, px: { xs: 1, sm: 2, md: 3 }, width: '100%', boxSizing: 'border-box' }}>
        {currentPath !== '/' && (
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              mb: 2,
              mt: 2
            }}
          >
            <Breadcrumbs categoryPath={categoryPath} />
          </Box>
        )}
        {children}
      </Container>

      <Modal
        open={isProfileOpen}
        onClose={handleCloseProfile}
      >
        <Box
          sx={{
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            width: "90%",
            maxWidth: 600,
            bgcolor: "background.paper",
            borderRadius: 2,
            boxShadow: 24,
            p: 4,
          }}
        >
          <UserProfile onClose={handleCloseProfile} />
        </Box>
      </Modal>

      <NotificationDrawer
        open={notificationDrawerOpen}
        onClose={() => setNotificationDrawerOpen(false)}
      />
      
      {/* Footer */}
      <Box 
        component="footer" 
        sx={{
          mt: 4,
          py: 3,
          bgcolor: 'background.paper',
          borderTop: '1px solid',
          borderColor: 'divider',
          color: 'text.secondary'
        }}
      >
        <Container maxWidth="lg">
          <Box sx={{ 
            display: 'flex', 
            flexDirection: {xs: 'column', md: 'row'}, 
            justifyContent: 'space-between',
            alignItems: {xs: 'center', md: 'flex-start'},
            textAlign: {xs: 'center', md: 'left'},
            gap: 2
          }}>
            <Box>
              <SveTuLogo width={50} height={50} />
              <Typography variant="body2" sx={{ mt: 1 }}>
                &copy; {new Date().getFullYear()} Sve Tu Platforma DOO
              </Typography>
              <Typography variant="body2">
                {t('footer.allRightsReserved', { defaultValue: 'All rights reserved' })}
              </Typography>
            </Box>
            
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 'bold', mb: 1 }}>
                {t('footer.links', { defaultValue: 'Links' })}
              </Typography>
              <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
                <Typography variant="body2" sx={{ '&:hover': { color: 'primary.main' } }}>
                  {t('navigation.home', { defaultValue: 'Home' })}
                </Typography>
              </Link>
              <Link to="/marketplace" style={{ textDecoration: 'none', color: 'inherit' }}>
                <Typography variant="body2" sx={{ '&:hover': { color: 'primary.main' } }}>
                  {t('navigation.marketplace', { defaultValue: 'Marketplace' })}
                </Typography>
              </Link>
              <Link to="/gis" style={{ textDecoration: 'none', color: 'inherit' }}>
                <Typography variant="body2" sx={{ '&:hover': { color: 'primary.main' } }}>
                  {t('gis:title', { defaultValue: 'Map Search' })}
                </Typography>
              </Link>
            </Box>
            
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 'bold', mb: 1 }}>
                {t('footer.legal', { defaultValue: 'Legal' })}
              </Typography>
              <Link to="/privacy-policy" style={{ textDecoration: 'none', color: 'inherit' }}>
                <Typography variant="body2" sx={{ '&:hover': { color: 'primary.main' } }}>
                  {t('footer.privacyPolicy', { defaultValue: 'Privacy Policy' })}
                </Typography>
              </Link>
              <Link to="/terms" style={{ textDecoration: 'none', color: 'inherit' }}>
                <Typography variant="body2" sx={{ '&:hover': { color: 'primary.main' } }}>
                  {t('footer.terms', { defaultValue: 'Terms of Service' })}
                </Typography>
              </Link>
            </Box>
            
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 'bold', mb: 1 }}>
                {t('footer.contact', { defaultValue: 'Contact' })}
              </Typography>
              <Typography variant="body2">
                Novi Sad, Serbia
              </Typography>
              <Typography variant="body2">
                <a href="mailto:info@svetu.rs" style={{ textDecoration: 'none', color: 'inherit', '&:hover': { color: 'primary.main' } }}>
                  info@svetu.rs
                </a>
              </Typography>
            </Box>
          </Box>
        </Container>
      </Box>
    </Box>
  );
};

export default Layout;
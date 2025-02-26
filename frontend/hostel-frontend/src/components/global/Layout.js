import React, { useState, useEffect } from "react";
import { Storefront } from '@mui/icons-material';
import NewMessageIndicator from './NewMessageIndicator';
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from '../../contexts/AuthContext';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { Bookmark } from '@mui/icons-material';
import UserProfile from "../user/UserProfile";
import { useChat } from '../../contexts/ChatContext';
import LanguageSwitcher from '../shared/LanguageSwitcher';
import NotificationDrawer from '../notifications/NotificationDrawer';
import { Settings } from '@mui/icons-material';
import { AccountBalanceWallet } from '@mui/icons-material';
import SveTuLogo from '../icons/SveTuLogo'; 
import { ShoppingBag, Store } from '@mui/icons-material';

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
} from "@mui/icons-material";

const Layout = ({ children }) => {
  const { t, i18n } = useTranslation('common');
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
      icon: <SveTuLogo width={60} height={60} />, 
    }
  ];
  
  useEffect(() => {
  //  console.log('Current language:', i18n.language);
  //  console.log('Available languages:', i18n.languages);
   // console.log('Translations loaded:', i18n.store.data);
  }, [i18n.language]);
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
              px: 2,
            }}
          >
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                gap: isMobile ? 1.5 : 3,
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
                  <Typography
                    variant="h6"
                    sx={{
                      fontSize: isMobile ? "0.85rem" : "1.1rem",
                      fontWeight: "bold",
                    }}
                  >
                    {item.label}
                  </Typography>
                </Box>
              ))}
            </Box>

            <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
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
  {t('navigation.storefronts')}
</MenuItem>
                    <MenuItem component={Link} to="/favorites">
                      <Bookmark fontSize="small" sx={{ mr: 1 }} />
                      {t('navigation.favorites')}
                    </MenuItem>
                    <MenuItem component={Link} to="/balance" onClick={handleCloseMenu}>
  <AccountBalanceWallet fontSize="small" sx={{ mr: 1 }} />
  {t('navigation.balance')}
</MenuItem>
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

      <Container maxWidth="lg" sx={{ py: 0 }}>
        {children}
      </Container>

      <NotificationDrawer
        open={notificationDrawerOpen}
        onClose={() => setNotificationDrawerOpen(false)}
      />
    </Box>
  );
};

export default Layout;
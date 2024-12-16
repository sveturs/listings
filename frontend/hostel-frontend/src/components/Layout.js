//frontend/hostel-frontend/src/components/Layout.js
import React, { useState } from "react";
import { ShoppingBag } from '@mui/icons-material';
import { Link, useLocation } from "react-router-dom";
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
  AttachMoney,
  Key,
  Logout,
  ListAlt,
  AddHome,
  AccountCircle,
} from "@mui/icons-material";
import { useAuth } from "../contexts/AuthContext";
import UserProfile from "../components/UserProfile";

const Layout = ({ children }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));
  const location = useLocation();

  // Определение текущей страницы
  const currentPath = location.pathname;

  const { user, login, logout } = useAuth();

  const [anchorEl, setAnchorEl] = useState(null);
  const [isProfileOpen, setIsProfileOpen] = useState(false);

  const handleOpenMenu = (e) => {
    setAnchorEl(e.currentTarget);
  };

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  const handleOpenProfile = () => {
    setIsProfileOpen(true);
    handleCloseMenu();
  };

  const handleCloseProfile = () => {
    setIsProfileOpen(false);
  };

  const menuItems = [
    { path: "/", label: "Hostel", icon: <HomeWork fontSize="medium" /> },
    { path: "/cars", label: "Auto", icon: <DirectionsCar fontSize="medium" /> },
    { path: "/marketplace", label: "Market", icon: <AttachMoney fontSize="medium" /> },
  ];

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
            {/* Левый блок (меню) */}
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
                    flexDirection: "column",
                    alignItems: "center",
                    gap: 0.3,
                    color: currentPath === item.path ? "primary.main" : "text.secondary",
                    fontWeight: currentPath === item.path ? 600 : 400,
                    fontSize: "0.9rem",
                    transition: "color 0.3s ease, transform 0.3s ease",
                    "&:hover": {
                      color: "primary.main",
                      transform: "scale(1.05)",
                    },
                  }}
                >
                  {item.icon}
                  <Typography
                    variant="body2"
                    sx={{
                      fontSize: isMobile ? "0.75rem" : "0.85rem",
                      textAlign: "center",
                    }}
                  >
                    {item.label}
                  </Typography>
                </Box>
              ))}
            </Box>

            {/* Правый блок (авторизация) */}
            <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
              {!user ? (
                <Tooltip title="Войти">
                  <IconButton onClick={login} color="primary">
                    <Key />
                  </IconButton>
                </Tooltip>
              ) : (
                <>
                  <Tooltip title="Мой профиль">
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
                    <MenuItem component={Link} to="/bookings">
                      <ListAlt fontSize="small" sx={{ mr: 1 }} />
                      Мои бронирования
                    </MenuItem>
                    <MenuItem component={Link} to="/my-listings">
                      <ShoppingBag fontSize="small" sx={{ mr: 1 }} />
                      Мои объявления
                    </MenuItem>
                    <MenuItem component={Link} to="/add-room">
                      <AddHome fontSize="small" sx={{ mr: 1 }} />
                      Добавить жильё
                    </MenuItem>
                    <MenuItem component={Link} to="/add-car">
                      <DirectionsCar fontSize="small" sx={{ mr: 1 }} />
                      Добавить автомобиль
                    </MenuItem>
                    <Divider />
                    <MenuItem onClick={logout}>
                      <Logout fontSize="small" sx={{ mr: 1 }} />
                      Выйти
                    </MenuItem>
                  </Menu>
                </>
              )}
            </Box>
          </Toolbar>
        </Container>
      </AppBar>

      {/* Модальное окно для редактирования профиля */}
      <Modal open={isProfileOpen} onClose={handleCloseProfile}>
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

      {/* Основной контент */}
      <Container maxWidth="lg" sx={{ py: 3 }}>
        {children}
      </Container>
    </Box>
  );
};

export default Layout;

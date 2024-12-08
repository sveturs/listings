import React from "react";
import { Link, useLocation } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Avatar,
  Menu,
  MenuItem,
  Box,
  Container,
  Divider,
  IconButton,
  useTheme,
  useMediaQuery
} from "@mui/material";
import { DirectionsCar, HomeWork, Menu as MenuIcon } from '@mui/icons-material';
import { useAuth } from "../contexts/AuthContext";

const Layout = ({ children }) => {
  const { user, login, logout } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const location = useLocation();
  
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [mobileMenuAnchor, setMobileMenuAnchor] = React.useState(null);

  // Группируем пункты меню по категориям
  const accommodationItems = [
    { to: "/bookings", label: "Все бронирования" },
    { to: "/add-room", label: "Добавить объявление" },
  ];

  const carItems = [
    { to: "/cars", label: "Список автомобилей" },
    { to: "/add-car", label: "Добавить автомобиль" },
  ];

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static" elevation={0} sx={{ bgcolor: 'background.paper', borderBottom: 1, borderColor: 'divider' }}>
        <Container maxWidth="xl">
          <Toolbar disableGutters>
            {/* Логотип */}
            <Typography
              variant="h6"
              component={Link}
              to="/"
              sx={{
                flexGrow: { xs: 1, md: 0 },
                mr: 4,
                textDecoration: 'none',
                color: 'text.primary',
                fontWeight: 700,
                letterSpacing: '-0.5px'
              }}
            >
              Hostel Booking System
            </Typography>

            {/* Десктопное меню */}
            {!isMobile && user && (
              <>
                <Box sx={{ display: 'flex', gap: 3, alignItems: 'center' }}>
                  {/* Секция жилья */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <HomeWork sx={{ color: 'primary.main' }} />
                    {accommodationItems.map((item) => (
                      <Button
                        key={item.to}
                        component={Link}
                        to={item.to}
                        sx={{
                          color: 'text.primary',
                          '&.active': {
                            color: 'primary.main',
                          }
                        }}
                      >
                        {item.label}
                      </Button>
                    ))}
                  </Box>

                  <Divider orientation="vertical" flexItem />

                  {/* Секция автомобилей */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <DirectionsCar sx={{ color: 'primary.main' }} />
                    {carItems.map((item) => (
                      <Button
                        key={item.to}
                        component={Link}
                        to={item.to}
                        sx={{
                          color: 'text.primary',
                          '&.active': {
                            color: 'primary.main',
                          }
                        }}
                      >
                        {item.label}
                      </Button>
                    ))}
                  </Box>
                </Box>

                <Box sx={{ flexGrow: 1 }} />
              </>
            )}

            {/* Профиль пользователя */}
            {user ? (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <Avatar
                  onClick={(e) => setAnchorEl(e.currentTarget)}
                  sx={{
                    width: 40,
                    height: 40,
                    cursor: 'pointer',
                    bgcolor: 'primary.main'
                  }}
                  src={user.pictureUrl}
                >
                  {user.name?.charAt(0)}
                </Avatar>
              </Box>
            ) : (
              <Button
                variant="contained"
                onClick={login}
                sx={{
                  borderRadius: 2,
                  textTransform: 'none',
                  px: 3
                }}
              >
                Войти через Google
              </Button>
            )}
          </Toolbar>
        </Container>
      </AppBar>

      {/* Основной контент */}
      <Container maxWidth="xl" sx={{ py: 4 }}>
        {children}
      </Container>
    </Box>
  );
};

export default Layout;
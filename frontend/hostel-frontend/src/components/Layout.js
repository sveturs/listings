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
  useTheme,
  useMediaQuery,
  ListItemIcon,
  ListItemText
} from "@mui/material";
import { 
  DirectionsCar, 
  HomeWork, 
  ListAlt,
  AddHome,
  Logout,
  Person
} from '@mui/icons-material';
import { useAuth } from "../contexts/AuthContext";

const Layout = ({ children }) => {
  const { user, login, logout } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const location = useLocation();
  
  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static" elevation={0} sx={{ bgcolor: 'background.paper', borderBottom: 1, borderColor: 'divider' }}>
        <Container maxWidth="xl">
        <Toolbar disableGutters sx={{ minHeight: '64px', display: 'flex', justifyContent: 'space-between' }}>
  {/* Логотипы */}
  <Box sx={{ display: 'flex', alignItems: 'center', flexGrow: 1, justifyContent: 'center' }}>
    <Typography
      variant="h6"
      component={Link}
      to="/"
      sx={{
        textDecoration: 'none',
        color: 'text.primary',
        fontWeight: 700,
        fontSize: '1.2rem',
        display: 'flex',
        alignItems: 'center',
        mr: 1,
      }}
    >
      Hostel Booking
      <HomeWork sx={{ ml: 1, color: 'primary.main' }} />
    </Typography>
    <Box sx={{ mx: 2, borderLeft: 1, height: '24px', borderColor: 'divider' }} />
    <Typography
      variant="h6"
      component={Link}
      to="/cars"
      sx={{
        textDecoration: 'none',
        color: 'text.primary',
        fontWeight: 700,
        fontSize: '1.2rem',
        display: 'flex',
        alignItems: 'center',
        ml: 1,
      }}
    >
      <DirectionsCar sx={{ mr: 1, color: 'primary.main' }} />
      Auto Booking
    </Typography>
  </Box>

  {/* Кнопка Войти */}
  <Box sx={{ position: 'absolute', right: '25%', transform: 'translateX(50%)' }}>
    {!user && (
      <Button
        variant="contained"
        onClick={login}
        sx={{ borderRadius: 2, px: 3 }}
      >
        Войти
      </Button>
    )}
  </Box>

  {/* Профиль */}
  {user && (
    <>
      <Avatar
        src={user.pictureUrl}
        sx={{ cursor: 'pointer', width: 36, height: 36 }}
        onClick={(e) => setAnchorEl(e.currentTarget)}
      />
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleCloseMenu}
        onClick={handleCloseMenu}
        PaperProps={{
          sx: { width: 220, mt: 1.5 },
        }}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
      >
        <MenuItem disabled>
          <ListItemText
            primary={user.name}
            secondary={user.email}
            primaryTypographyProps={{
              variant: 'subtitle2',
              noWrap: true,
            }}
            secondaryTypographyProps={{
              variant: 'caption',
              noWrap: true,
            }}
          />
        </MenuItem>
        <Divider />
        <MenuItem component={Link} to="/bookings">
          <ListItemIcon>
            <ListAlt fontSize="small" />
          </ListItemIcon>
          <ListItemText>Мои бронирования</ListItemText>
        </MenuItem>
        <Divider />
        <MenuItem component={Link} to="/add-room">
          <ListItemIcon>
            <AddHome fontSize="small" />
          </ListItemIcon>
          <ListItemText>Добавить жильё</ListItemText>
        </MenuItem>
        <MenuItem component={Link} to="/add-car">
          <ListItemIcon>
            <DirectionsCar fontSize="small" />
          </ListItemIcon>
          <ListItemText>Добавить автомобиль</ListItemText>
        </MenuItem>
        <Divider />
        <MenuItem onClick={logout}>
          <ListItemIcon>
            <Logout fontSize="small" />
          </ListItemIcon>
          <ListItemText>Выйти</ListItemText>
        </MenuItem>
      </Menu>
    </>
  )}
</Toolbar>

        </Container>
      </AppBar>
      {/* Основной контент */}
      <Container maxWidth="xl" sx={{ py: 0 }}>
        {children}
      </Container>
    </Box>
  );
};

export default Layout;
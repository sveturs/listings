import React from "react";
import { Link } from "react-router-dom";
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  Button, 
  Avatar,
  Menu,
  MenuItem,
  Box,
} from "@mui/material";
import { useAuth } from "../contexts/AuthContext";

const Layout = ({ children }) => {
  const { user, login, logout } = useAuth();
  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <Box>
      <AppBar position="static">
        <Toolbar>
          <Typography
            variant="h6"
            component={Link}
            to="/"
            sx={{
              flexGrow: 1,
              textDecoration: 'none',
              color: 'inherit'
            }}
          >
            Hostel Booking System
          </Typography>

          {user ? (
            <>
              <Button
                color="inherit"
                component={Link}
                to="/bookings"
              >
                Все бронирования
              </Button>
              <Button
                color="inherit"
                component={Link}
                to="/add-room"
              >
                Добавить объявление
              </Button>

              <Box 
                onClick={handleMenu}
                sx={{ 
                  display: 'flex',
                  alignItems: 'center',
                  ml: 2,
                  cursor: 'pointer'
                }}
              >
                <Avatar
                  sx={{ 
                    width: 32,
                    height: 32,
                    bgcolor: 'primary.dark'
                  }}
                >
                  {user.name.charAt(0)}
                </Avatar>
              </Box>
              <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                open={Boolean(anchorEl)}
                onClose={handleClose}
              >
                <MenuItem disabled>
                  <Typography variant="body2">{user.email}</Typography>
                </MenuItem>
                <MenuItem onClick={logout}>Выйти</MenuItem>
              </Menu>
            </>
          ) : (
            <Button 
              color="inherit"
              onClick={login}
            >
              ВОЙТИ ЧЕРЕЗ GOOGLE
            </Button>
          )}
        </Toolbar>
      </AppBar>
      {children}
    </Box>
  );
};

export default Layout;
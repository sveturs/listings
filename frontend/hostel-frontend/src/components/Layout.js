import React from "react";
import { Link } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Container,
  IconButton,
  Tooltip,
  Avatar,
  Menu,
  MenuItem,
  Divider,
  useMediaQuery,
  useTheme,
  ListItemIcon,
  ListItemText
} from "@mui/material";
import {
  DirectionsCar,
  HomeWork,
  Key,
  Logout,
  ListAlt,
  AddHome
} from "@mui/icons-material";
import { useAuth } from "../contexts/AuthContext";

const Layout = ({ children }) => {
  const { user, login, logout } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar
        position="static"
        elevation={0}
        sx={{
          bgcolor: "background.paper",
          borderBottom: 1,
          borderColor: "divider"
        }}
      >
        <Container maxWidth="xl">
          <Toolbar
            disableGutters
            sx={{
              minHeight: "64px",
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center"
            }}
          >
            {/* Левый блок */}
            <Box sx={{ display: "flex", alignItems: "center" }}>
              <Typography
                variant="h6"
                component={Link}
                to="/"
                sx={{
                  textDecoration: "none",
                  color: "text.primary",
                  fontWeight: 700,
                  fontSize: "1.2rem",
                  display: "flex",
                  alignItems: "center",
                  mr: 0.4
                }}
              >
                Hostel Booking
              </Typography>
              <HomeWork sx={{ color: "primary.main", mr: 0.4 }} />
              <Divider
                orientation="vertical"
                flexItem
                sx={{
                  bgcolor: "text.secondary",
                  mx: 0,
                  height: "1.5rem"
                }}
              />
              <DirectionsCar sx={{ color: "primary.main", mr: 0.2 }} />
              <Typography
                variant="h6"
                component={Link}
                to="/cars"
                sx={{
                  textDecoration: "none",
                  color: "text.primary",
                  fontWeight: 700,
                  fontSize: "1.2rem",
                  display: "flex",
                  alignItems: "center"
                }}
              >
                Auto Booking
              </Typography>
            </Box>
            <Box sx={{ display: "flex", alignItems: "center" }}>
    {/* существующие ссылки */}
    <Divider orientation="vertical" flexItem sx={{ mx: 1 }} />
    <Typography
        variant="h6"
        component={Link}
        to="/marketplace"
        sx={{
            textDecoration: "none",
            color: "text.primary",
            fontWeight: 700,
            fontSize: "1.2rem",
            display: "flex",
            alignItems: "center"
        }}
    >
        Marketplace
    </Typography>
</Box>
            {/* Правый блок */}
            <Box sx={{ display: "flex", alignItems: "center" }}>
              {!user ? (
                <Tooltip title={!isMobile ? "Залогиниться" : ""}>
                  <IconButton onClick={login}>
                    <Key fontSize="large" sx={{ color: "primary.main" }} />
                  </IconButton>
                </Tooltip>
              ) : (
                <>
                  <Avatar
                    src={user.pictureUrl}
                    sx={{ cursor: "pointer", width: 36, height: 36, ml: 2 }}
                    onClick={(e) => setAnchorEl(e.currentTarget)}
                  />
                  <Menu
                    anchorEl={anchorEl}
                    open={Boolean(anchorEl)}
                    onClose={handleCloseMenu}
                    onClick={handleCloseMenu}
                    PaperProps={{
                      sx: { width: 220, mt: 1.5 }
                    }}
                    transformOrigin={{ horizontal: "right", vertical: "top" }}
                    anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
                  >
                    <MenuItem disabled>
                      <ListItemText
                        primary={user.name}
                        secondary={user.email}
                        primaryTypographyProps={{
                          variant: "subtitle2",
                          noWrap: true
                        }}
                        secondaryTypographyProps={{
                          variant: "caption",
                          noWrap: true
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
            </Box>
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

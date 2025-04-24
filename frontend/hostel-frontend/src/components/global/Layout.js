// frontend/hostel-frontend/src/components/global/Layout.js
import React, { useState, useEffect } from "react";
import NewMessageIndicator from './NewMessageIndicator';
import { Link, useLocation, useNavigate, useSearchParams } from "react-router-dom";
import { useAuth } from '../../contexts/AuthContext';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import UserProfile from "../user/UserProfile";
import { useChat } from '../../contexts/ChatContext';
import LanguageSwitcher from '../shared/LanguageSwitcher';
import NotificationDrawer from '../notifications/NotificationDrawer';
import SveTuLogo from '../icons/SveTuLogo';
import { useLocation as useGeoLocation } from '../../contexts/LocationContext';
import CitySelector from './CitySelector';
import Breadcrumbs from '../../components/marketplace/Breadcrumbs';
import { Plus, MapPin, X, Menu as MenuIcon, Check } from 'lucide-react';
import AutocompleteInput from '../shared/AutocompleteInput';
import { isAdmin } from '../../utils/adminUtils';
import {
  HomeWork,
  DirectionsCar,
  Key,
  Logout,
  ListAlt,
  AddHome,
  AccountCircle,
  Chat,
  Favorite,
  Store,
  AccountBalanceWallet,
  Search,
  Settings,
  ShoppingBag,
  Notifications as NotificationsIcon,
  KeyboardArrowDown,
  HelpOutline,
  Person,
  Business,
  Login
} from "@mui/icons-material";
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
  ListItemIcon,
  ListItemText,
  Divider,
  useTheme,
  useMediaQuery,
  Modal,
  Alert,
  Slide,
  Button,
  TextField,
  InputAdornment,
  Grid,
  Drawer,
  List,
  ListItem
} from "@mui/material";

const Layout = ({ children }) => {
  const { userLocation, setCity, locationDismissed, dismissLocationSuggestion } = useGeoLocation();
  const [showLocationAlert, setShowLocationAlert] = useState(false);
  const [showLocationPicker, setShowLocationPicker] = useState(false);
  const { t, i18n } = useTranslation(['common', 'marketplace']);
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
  const [categoryMenuAnchor, setCategoryMenuAnchor] = useState(null);
  const [mobileCategoryDrawerOpen, setMobileCategoryDrawerOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [languageAnchorEl, setLanguageAnchorEl] = useState(null);

  // Добавляем состояние для поискового запроса
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchQuery, setSearchQuery] = useState('');

  // Добавляем данные категорий (или загрузите их из API)
  const categoryGroups = [
    {
      title: t('categories.realestate', { defaultValue: 'Недвижимость', ns: 'marketplace' }),
      icon: <HomeWork />,
      path: "/marketplace?category_id=1000",
      subcategories: [
        { name: t('categories.apartments', { defaultValue: 'Квартиры', ns: 'marketplace' }), path: "/marketplace?category_id=1100" },
        { name: t('categories.houses', { defaultValue: 'Дома', ns: 'marketplace' }), path: "/marketplace?category_id=1300" },
        { name: t('categories.rooms', { defaultValue: 'Комнаты', ns: 'marketplace' }), path: "/marketplace?category_id=1200" },
        { name: t('categories.commercial', { defaultValue: 'Коммерческая', ns: 'marketplace' }), path: "/marketplace?category_id=1600" },
        { name: t('categories.land', { defaultValue: 'Земельные участки', ns: 'marketplace' }), path: "/marketplace?category_id=1400" },
        { name: t('categories.garages', { defaultValue: 'Гаражи', ns: 'marketplace' }), path: "/marketplace?category_id=1500" },
      ]
    },
    {
      title: t('categories.transport', { defaultValue: 'Транспорт', ns: 'marketplace' }),
      icon: <DirectionsCar />,
      path: "/marketplace?category_id=2000",
      subcategories: [
        { name: t('categories.cars', { defaultValue: 'Автомобили', ns: 'marketplace' }), path: "/marketplace?category_id=2000" },
        { name: t('categories.motorcycles', { defaultValue: 'Мотоциклы', ns: 'marketplace' }), path: "/marketplace?category_id=2600" },
        { name: t('categories.parts', { defaultValue: 'Запчасти', ns: 'marketplace' }), path: "/marketplace?category_id=2800" },
        { name: t('categories.trucks', { defaultValue: 'Грузовики', ns: 'marketplace' }), path: "/marketplace?category_id=2200" },
        { name: t('categories.specialVehicles', { defaultValue: 'Спецтехника', ns: 'marketplace' }), path: "/marketplace?category_id=2300" },
      ]
    },
    {
      title: t('categories.electronics', { defaultValue: 'Электроника', ns: 'marketplace' }),
      icon: <ShoppingBag />,
      path: "/marketplace?category_id=3000",
      subcategories: [
        { name: t('categories.phones', { defaultValue: 'Телефоны', ns: 'marketplace' }), path: "/marketplace?category_id=3100" },
        { name: t('categories.computers', { defaultValue: 'Компьютеры', ns: 'marketplace' }), path: "/marketplace?category_id=3300" },
        { name: t('categories.tvAudio', { defaultValue: 'ТВ и аудио', ns: 'marketplace' }), path: "/marketplace?category_id=3200" },
        { name: t('categories.photoVideo', { defaultValue: 'Фото и видео', ns: 'marketplace' }), path: "/marketplace?category_id=3700" },
      ]
    },
    {
      title: t('categories.forHome', { defaultValue: 'Для дома', ns: 'marketplace' }),
      icon: <HomeWork />,
      path: "/marketplace?category_id=5000",
      subcategories: [
        { name: t('categories.furniture', { defaultValue: 'Мебель', ns: 'marketplace' }), path: "/marketplace?category_id=5200" },
        { name: t('categories.appliances', { defaultValue: 'Бытовая техника', ns: 'marketplace' }), path: "/marketplace?category_id=4100" },
        { name: t('categories.kitchenware', { defaultValue: 'Посуда', ns: 'marketplace' }), path: "/marketplace?category_id=5400" },
        { name: t('categories.renovation', { defaultValue: 'Ремонт', ns: 'marketplace' }), path: "/marketplace?category_id=5100" },
      ]
    },
    {
      title: t('categories.forGarden', { defaultValue: 'Для сада', ns: 'marketplace' }),
      icon: <HomeWork />,
      path: "/marketplace?category_id=6000",
      subcategories: [
        { name: t('categories.gardenFurniture', { defaultValue: 'Садовая мебель', ns: 'marketplace' }), path: "/marketplace?category_id=6050" },
        { name: t('categories.plants', { defaultValue: 'Растения', ns: 'marketplace' }), path: "/marketplace?category_id=6750" },
        { name: t('categories.tools', { defaultValue: 'Инструменты', ns: 'marketplace' }), path: "/marketplace?category_id=6100" },
      ]
    }
  ];

  // Обработчики для языкового меню
  const handleOpenLanguageMenu = (event) => {
    setLanguageAnchorEl(event.currentTarget);
  };

  const handleCloseLanguageMenu = () => {
    setLanguageAnchorEl(null);
  };

  const handleLanguageChange = (lang) => {
    i18n.changeLanguage(lang);
    localStorage.setItem('preferredLanguage', lang);
    document.documentElement.lang = lang;
    handleCloseLanguageMenu();
    if (mobileMenuOpen) {
      handleCloseMobileMenu();
    }
  };

  // Функция для получения текущего отображаемого языка
  const getCurrentLanguageDisplay = () => {
    switch (i18n.language) {
      case 'en':
        return 'EN';
      case 'sr':
        return 'SR';
      case 'ru':
        return 'RU';
      default:
        return 'RU';
    }
  };

  // Для отображения флага в мобильной версии
  const getLanguageFlag = () => {
    switch (i18n.language) {
      case 'en':
        return '🇬🇧';
      case 'sr':
        return '🇷🇸';
      case 'ru':
        return '🇷🇺';
      default:
        return '🇷🇺';
    }
  };

  // Обработчик поиска
  const handleSearch = (event) => {
    // Если нажата клавиша Enter или клик по кнопке поиска
    if (event.key === 'Enter' || event.type === 'click') {
      // Сохраняем текущие параметры поиска
      const currentParams = Object.fromEntries(searchParams.entries());

      // Обновляем параметр запроса
      setSearchParams({
        ...currentParams,
        query: searchQuery
      });

      // Перенаправляем на страницу маркетплейса, если мы на другой странице
      if (!currentPath.includes('/marketplace')) {
        navigate('/marketplace?query=' + encodeURIComponent(searchQuery));
      }
    }
  };

  // При изменении searchParams обновляем поле ввода
  useEffect(() => {
    const query = searchParams.get('query') || '';
    setSearchQuery(query);
  }, [searchParams]);

  const handleOpenCategoryMenu = (event) => {
    setCategoryMenuAnchor(event.currentTarget);
  };

  const handleCloseCategoryMenu = () => {
    setCategoryMenuAnchor(null);
  };

  const handleOpenMobileCategoryDrawer = () => {
    setMobileCategoryDrawerOpen(true);
  };

  const handleCloseMobileCategoryDrawer = () => {
    setMobileCategoryDrawerOpen(false);
  };

  // Функции для управления мобильным меню
  const handleOpenMobileMenu = () => {
    setMobileMenuOpen(true);
  };

  const handleCloseMobileMenu = () => {
    setMobileMenuOpen(false);
  };

  // Обновленная функция handleCategoryClick
  const handleCategoryClick = (path) => {
    navigate(path);
    handleCloseCategoryMenu();
    handleCloseMobileCategoryDrawer(); // Закрываем мобильное меню тоже
  };

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
        <Typography>{t('navigation.messages', { defaultValue: 'Сообщения', ns: 'marketplace' })}</Typography>
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
          bgcolor: "#104054", // Темно-синий цвет верхней части
          color: "white",
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
              px: 5,
            }}
          >
            <Box
              component={Link}
              to="/"
              sx={{
                textDecoration: "none",
                display: "flex",
                alignItems: "center",
                color: "white",
                fontWeight: 600,
                fontSize: "1rem",
              }}
            >
              {/* компонент SveTuLogo */}
              <SveTuLogo width={60} height={60} />
            </Box>

            <Box sx={{ display: "flex", alignItems: "center", gap: isMobile ? 0.5 : 2 }}>
              {/* Добавляем иконку меню на мобильных устройствах */}
              {isMobile && (
                <IconButton
                  onClick={handleOpenMobileMenu}
                  sx={{ color: 'white' }}
                >
                  <MenuIcon size={24} />
                </IconButton>
              )}

              {!isMobile && (
                <>
                  <Button
                    variant="text"
                    component={Link}
                    to="/business"
                    sx={{ color: 'white', textTransform: 'none' }}
                  >
                    {t('navigation.forBusiness', { defaultValue: 'Для бизнеса', ns: 'marketplace' })}
                  </Button>

                  <Button
                    variant="text"
                    component={Link}
                    to="/help"
                    sx={{ color: 'white', textTransform: 'none' }}
                  >
                    {t('navigation.help', { defaultValue: 'Помощь', ns: 'marketplace' })}
                  </Button>

                  <Button
                    variant="text"
                    component={Link}
                    to="/profile"
                    sx={{ color: 'white', textTransform: 'none' }}
                  >
                    {t('navigation.cabinet', { defaultValue: 'Кабинет', ns: 'marketplace' })}
                  </Button>

                  {/* Добавляем компонент переключения языка в десктопную версию */}
                  <LanguageSwitcher />
                </>
              )}

              <Button
                variant="contained"
                onClick={() => navigate('/marketplace/create')}
                sx={{
                  bgcolor: '#FF5000',
                  '&:hover': { bgcolor: '#FF6A00' },
                  borderRadius: '4px',
                  textTransform: 'none',
                  fontSize: isMobile ? '0.75rem' : 'inherit',
                  px: isMobile ? 1 : 2
                }}
              >
                {isMobile ? 
                  t('navigation.post', { defaultValue: 'Подать', ns: 'marketplace' }) : 
                  t('navigation.postListing', { defaultValue: 'Подать объявление', ns: 'marketplace' })
                }
              </Button>

              {!isMobile && (
                <Button
                  variant="outlined"
                  component={Link}
                  to="/my-listings"
                  sx={{
                    color: 'white',
                    borderColor: 'white',
                    textTransform: 'none',
                    '&:hover': { borderColor: 'white', bgcolor: 'rgba(255,255,255,0.1)' }
                  }}
                >
                  {t('navigation.myListings', { defaultValue: 'Мои объявления', ns: 'marketplace' })}
                </Button>
              )}

              {!user ? (
                <IconButton
                  sx={{
                    width: isMobile ? 32 : 38,
                    height: isMobile ? 32 : 38
                  }}
                  onClick={handleOpenLanguageMenu}
                >
                  <Avatar sx={{
                    bgcolor: '#FFFFFF',
                    color: '#004494',
                    width: isMobile ? 30 : 36,
                    height: isMobile ? 30 : 36
                  }}>
                    <Typography variant="button" sx={{ fontWeight: 'bold', fontSize: isMobile ? '0.6rem' : '0.75rem' }}>
                      {getCurrentLanguageDisplay()}
                    </Typography>
                  </Avatar>
                </IconButton>
              ) : (
                <IconButton onClick={handleOpenMenu} sx={{ p: 0 }}>
                  <Avatar
                    src={user.pictureUrl}
                    alt={user.name}
                    sx={{ width: isMobile ? 32 : 38, height: isMobile ? 32 : 38 }}
                  />
                </IconButton>
              )}
            </Box>
          </Toolbar>
        </Container>
      </AppBar>

      {/* Добавляем вторую строку с логотипом и поиском */}
      <Box
        sx={{
          bgcolor: "#FFF5F0", // Светло-оранжевый фон как в образце
          py: 2
        }}
      >
        <Container maxWidth="lg">
          <Box sx={{
            display: "flex",
            flexDirection: isMobile ? 'column' : 'row',
            justifyContent: "space-between",
            alignItems: isMobile ? "stretch" : "center",
            gap: isMobile ? 2 : 0
          }}>
            <Box sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: isMobile ? 'space-between' : 'flex-start',
              gap: 2
            }}>
              <Typography
                component={Link}
                to="/"
                variant={isMobile ? "h5" : "h4"}
                sx={{
                  fontWeight: 'bold',
                  textDecoration: 'none',
                  color: '#FF5000',
                  '& span': { color: '#004494' }
                }}
              >
                Sve <span>Tu</span>
              </Typography>

              <Button
                sx={{
                  color: '#004494',
                  textTransform: 'none',
                  fontWeight: 'normal',
                  fontSize: isMobile ? '0.875rem' : 'inherit'
                }}
                onClick={isMobile ? handleOpenMobileCategoryDrawer : handleOpenCategoryMenu}
                endIcon={<KeyboardArrowDown />}
              >
                {t('navigation.allCategories', { defaultValue: 'ВСЕ КАТЕГОРИИ', ns: 'marketplace' })}
              </Button>

              {isMobile && (
                <Box
                  sx={{
                    display: "flex",
                    alignItems: "center",
                    gap: 0.5,
                    cursor: 'pointer'
                  }}
                  onClick={() => setShowLocationPicker(true)}
                >
                  <MapPin size={16} color="#004494" />
                  <Typography
                    variant="caption"
                    sx={{
                      color: "#004494",
                      fontWeight: "medium"
                    }}
                  >
                    {userLocation?.city || t('cities.noviSad', { defaultValue: 'Нови-Сад', ns: 'marketplace' })}
                  </Typography>
                </Box>
              )}
            </Box>

            <Box sx={{ display: "flex", width: isMobile ? "100%" : "40%" }}>
              <AutocompleteInput
                placeholder={t('search.find', { defaultValue: 'Найти', ns: 'marketplace' })}
                value={searchQuery}
                onChange={(value) => setSearchQuery(value)}
                onSearch={(value, categoryId) => {
                  // Сохраняем текущие параметры поиска
                  const currentParams = Object.fromEntries(searchParams.entries());

                  // Обновляем параметры запроса
                  const updatedParams = {
                    ...currentParams,
                    query: value
                  };

                  // Если передан ID категории, добавляем его в параметры
                  if (categoryId) {
                    updatedParams.category_id = categoryId;
                  }

                  setSearchParams(updatedParams);

                  // Перенаправляем на страницу маркетплейса, если мы на другой странице
                  if (!currentPath.includes('/marketplace')) {
                    const queryString = new URLSearchParams(updatedParams).toString();
                    navigate(`/marketplace?${queryString}`);
                  }
                }}
                sx={{
                  width: '100%',
                  '& .MuiOutlinedInput-root': {
                    borderRadius: '4px',
                    bgcolor: 'white',
                  }
                }}
              />
            </Box>


            {!isMobile && (
              <Box
                sx={{
                  display: "flex",
                  alignItems: "center",
                  gap: 1,
                  cursor: 'pointer'
                }}
                onClick={() => setShowLocationPicker(true)}
              >
                <MapPin size={18} color="#004494" />
                <Typography
                  variant="body2"
                  sx={{
                    color: "#004494",
                    fontWeight: "medium",
                    '&:hover': { textDecoration: 'underline' }
                  }}
                >
                  {userLocation?.city || t('cities.noviSad', { defaultValue: 'Нови-Сад', ns: 'marketplace' })}
                </Typography>
              </Box>
            )}
          </Box>
        </Container>
      </Box>

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
              {t('location.useThisCity', { defaultValue: 'Использовать этот город', ns: 'marketplace' })}
            </Button>
          }
        >
          {t('location.detectedCity', {
            defaultValue: 'Мы определили, что вы находитесь в городе {{city}}',
            city: userLocation?.city,
            ns: 'marketplace'
          })}
        </Alert>
      </Slide>

      <Container maxWidth="lg" sx={{ py: 0 }}>
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

      {/* Модальное окно для выбора города */}
      <Modal
        open={showLocationPicker}
        onClose={() => setShowLocationPicker(false)}
      >
        <Box
          sx={{
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            width: "90%",
            maxWidth: 500,
            bgcolor: "background.paper",
            borderRadius: 2,
            boxShadow: 24,
            p: 4,
          }}
        >
          <Typography variant="h6" component="h2" gutterBottom>
            {t('location.citySelection', { defaultValue: 'Выбор города', ns: 'marketplace' })}
          </Typography>
          <CitySelector onClose={() => setShowLocationPicker(false)} />
        </Box>
      </Modal>

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

      {/* Выпадающее меню для категорий (для десктопов) */}
      <Menu
        anchorEl={categoryMenuAnchor}
        open={Boolean(categoryMenuAnchor)}
        onClose={handleCloseCategoryMenu}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'left',
        }}
        PaperProps={{
          style: {
            maxHeight: '80vh',
            width: '650px', // Увеличиваем ширину для мультиколоночного отображения
            padding: '12px'
          },
        }}
      >
        <Grid container spacing={2}>
          {categoryGroups.map((group, index) => (
            <Grid item xs={4} key={index}>
              <Box
                sx={{
                  mb: 2,
                  pb: 2,
                  borderBottom: index < categoryGroups.length - 3 ? '1px solid' : 'none',
                  borderColor: 'divider'
                }}
              >
                {/* Заголовок категории */}
                <Box
                  component={Button}
                  onClick={() => handleCategoryClick(group.path)}
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    width: '100%',
                    textAlign: 'left',
                    justifyContent: 'flex-start',
                    color: 'primary.main',
                    fontWeight: 'bold',
                    textTransform: 'none',
                    mb: 1,
                    '&:hover': { backgroundColor: 'transparent' }
                  }}
                >
                  <Box sx={{ color: 'primary.main', mr: 1 }}>
                    {group.icon}
                  </Box>
                  <Typography variant="subtitle1">
                    {group.title}
                  </Typography>
                </Box>

                {/* Подкатегории */}
                <Box sx={{ pl: 2 }}>
                  {group.subcategories.map((subcat, subIdx) => (
                    <Button
                      key={subIdx}
                      onClick={() => handleCategoryClick(subcat.path)}
                      sx={{
                        display: 'block',
                        width: '100%',
                        textAlign: 'left',
                        justifyContent: 'flex-start',
                        py: 0.5,
                        px: 1,
                        color: 'text.primary',
                        textTransform: 'none',
                        '&:hover': { backgroundColor: 'action.hover' }
                      }}
                    >
                      <Typography variant="body2">
                        {subcat.name}
                      </Typography>
                    </Button>
                  ))}
                </Box>
              </Box>
            </Grid>
          ))}
        </Grid>
      </Menu>

      {/* Drawer для мобильного меню категорий */}
      <Drawer
        anchor="right"
        open={mobileCategoryDrawerOpen}
        onClose={handleCloseMobileCategoryDrawer}
        sx={{
          '& .MuiDrawer-paper': {
            width: { xs: '80%', sm: '60%' },
            maxWidth: '350px',
          }
        }}
      >
        <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">{t('categories.title', { defaultValue: 'Категории', ns: 'marketplace' })}</Typography>
            <IconButton onClick={handleCloseMobileCategoryDrawer}>
              <X size={20} />
            </IconButton>
          </Box>
        </Box>

        <Box sx={{ overflow: 'auto', flex: 1 }}>
          {categoryGroups.map((group, index) => (
            <Box key={index} sx={{ mb: 2 }}>
              <Box
                component={Link}
                to={group.path}
                onClick={() => handleCategoryClick(group.path)}
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  px: 2,
                  py: 1,
                  color: 'primary.main',
                  fontWeight: 'bold',
                  textDecoration: 'none',
                  borderBottom: '1px solid',
                  borderColor: 'divider',
                  '&:hover': {
                    bgcolor: 'action.hover'
                  }
                }}
              >
                <Box sx={{ mr: 2, color: 'primary.main' }}>
                  {group.icon}
                </Box>
                <Typography variant="subtitle1">{group.title}</Typography>
                <Box sx={{ ml: 'auto' }}>
                  <KeyboardArrowDown />
                </Box>
              </Box>

              <Box sx={{ bgcolor: 'background.default' }}>
                {group.subcategories.map((subcat, idx) => (
                  <Box
                    key={idx}
                    component={Link}
                    to={subcat.path}
                    onClick={() => handleCategoryClick(subcat.path)}
                    sx={{
                      display: 'block',
                      px: 4,
                      py: 1,
                      color: 'text.primary',
                      textDecoration: 'none',
                      borderBottom: '1px solid',
                      borderColor: 'divider',
                      '&:hover': {
                        bgcolor: 'action.hover'
                      }
                    }}
                  >
                    <Typography variant="body2">{subcat.name}</Typography>
                  </Box>
                ))}
              </Box>
            </Box>
          ))}
        </Box>
      </Drawer>

      {/* Добавляем мобильное меню пользователя (гамбургер-меню) */}
      <Drawer
        anchor="left"
        open={mobileMenuOpen}
        onClose={handleCloseMobileMenu}
        sx={{
          '& .MuiDrawer-paper': {
            width: '80%',
            maxWidth: '300px',
          }
        }}
      >
        <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">{t('navigation.menu', { defaultValue: 'Меню', ns: 'marketplace' })}</Typography>
            <IconButton onClick={handleCloseMobileMenu}>
              <X size={20} />
            </IconButton>
          </Box>
        </Box>

        <List>
          <ListItem button component={Link} to="/business" onClick={handleCloseMobileMenu}>
            <ListItemIcon><Business size={20} /></ListItemIcon>
            <ListItemText primary={t('navigation.forBusiness', { defaultValue: 'Для бизнеса', ns: 'marketplace' })} />
          </ListItem>

          <ListItem button component={Link} to="/help" onClick={handleCloseMobileMenu}>
            <ListItemIcon><HelpOutline size={20} /></ListItemIcon>
            <ListItemText primary={t('navigation.help', { defaultValue: 'Помощь', ns: 'marketplace' })} />
          </ListItem>

          <ListItem button component={Link} to="/profile" onClick={handleCloseMobileMenu}>
            <ListItemIcon><Person size={20} /></ListItemIcon>
            <ListItemText primary={t('navigation.cabinet', { defaultValue: 'Кабинет', ns: 'marketplace' })} />
          </ListItem>

          <ListItem button component={Link} to="/my-listings" onClick={handleCloseMobileMenu}>
            <ListItemIcon><ListAlt size={20} /></ListItemIcon>
            <ListItemText primary={t('navigation.myListings', { defaultValue: 'Мои объявления', ns: 'marketplace' })} />
          </ListItem>

          <Divider />

          {user ? (
            <>
              <ListItem button component={Link} to="/favorites" onClick={handleCloseMobileMenu}>
                <ListItemIcon><Favorite size={20} /></ListItemIcon>
                <ListItemText primary={t('navigation.favorites', { defaultValue: 'Избранное', ns: 'marketplace' })} />
              </ListItem>

              <ListItem button component={Link} to="/marketplace/chat" onClick={handleCloseMobileMenu}>
                <ListItemIcon><Chat size={20} /></ListItemIcon>
                <ListItemText primary={t('navigation.messages', { defaultValue: 'Сообщения', ns: 'marketplace' })} />
              </ListItem>

              <ListItem button component={Link} to="/notifications/settings" onClick={handleCloseMobileMenu}>
                <ListItemIcon><NotificationsIcon size={20} /></ListItemIcon>
                <ListItemText primary={t('navigation.notifications', { defaultValue: 'Оповещения', ns: 'marketplace' })} />
              </ListItem>

              <ListItem button component={Link} to="/storefronts" onClick={handleCloseMobileMenu}>
                <ListItemIcon><Store size={20} /></ListItemIcon>
                <ListItemText primary={t('navigation.storefronts', { defaultValue: 'Мои витрины', ns: 'marketplace' })} />
              </ListItem>

              <ListItem button component={Link} to="/balance" onClick={handleCloseMobileMenu}>
                <ListItemIcon><AccountBalanceWallet size={20} /></ListItemIcon>
                <ListItemText primary={t('navigation.balance', { defaultValue: 'Мой баланс', ns: 'marketplace' })} />
              </ListItem>

              <Divider />

              <ListItem>
                <ListItemText primary={t('navigation.language', { defaultValue: 'Язык', ns: 'marketplace' })} />
              </ListItem>
              <ListItem button onClick={() => handleLanguageChange('en')}>
                <ListItemText primary="English" secondary="English" />
                {i18n.language === 'en' && <ListItemIcon sx={{ minWidth: 'auto' }}><Check size={16} /></ListItemIcon>}
              </ListItem>
              <ListItem button onClick={() => handleLanguageChange('sr')}>
                <ListItemText primary="Srpski" secondary={t('languages.serbian', { defaultValue: 'Сербский', ns: 'marketplace' })} />
                {i18n.language === 'sr' && <ListItemIcon sx={{ minWidth: 'auto' }}><Check size={16} /></ListItemIcon>}
              </ListItem>
              <ListItem button onClick={() => handleLanguageChange('ru')}>
                <ListItemText primary={t('languages.russian', { defaultValue: 'Русский', ns: 'marketplace' })} />
                {i18n.language === 'ru' && <ListItemIcon sx={{ minWidth: 'auto' }}><Check size={16} /></ListItemIcon>}
              </ListItem>

              <Divider />

              <ListItem button onClick={() => { logout(); handleCloseMobileMenu(); }}>
                <ListItemIcon><Logout size={20} /></ListItemIcon>
                <ListItemText primary={t('auth.signout', { defaultValue: 'Выйти', ns: 'marketplace' })} />
              </ListItem>
            </>
          ) : (
            <>
              <ListItem button onClick={() => { login(); handleCloseMobileMenu(); }}>
                <ListItemIcon><Login size={20} /></ListItemIcon>
                <ListItemText primary={t('auth.signin', { defaultValue: 'Войти', ns: 'marketplace' })} />
              </ListItem>

              <Divider />

              <ListItem>
                <ListItemText primary={t('navigation.language', { defaultValue: 'Язык', ns: 'marketplace' })} />
              </ListItem>
              <ListItem button onClick={() => handleLanguageChange('en')}>
                <ListItemText primary="English" secondary="English" />
                {i18n.language === 'en' && <ListItemIcon sx={{ minWidth: 'auto' }}><Check size={16} /></ListItemIcon>}
              </ListItem>
              <ListItem button onClick={() => handleLanguageChange('sr')}>
                <ListItemText primary="Srpski" secondary={t('languages.serbian', { defaultValue: 'Сербский', ns: 'marketplace' })} />
                {i18n.language === 'sr' && <ListItemIcon sx={{ minWidth: 'auto' }}><Check size={16} /></ListItemIcon>}
              </ListItem>
              <ListItem button onClick={() => handleLanguageChange('ru')}>
                <ListItemText primary={t('languages.russian', { defaultValue: 'Русский', ns: 'marketplace' })} />
                {i18n.language === 'ru' && <ListItemIcon sx={{ minWidth: 'auto' }}><Check size={16} /></ListItemIcon>}
              </ListItem>
            </>
          )}
        </List>
      </Drawer>

      {/* Меню для пользователя */}
      <Menu
        anchorEl={anchorEl}
        id="user-menu"
        open={Boolean(anchorEl)}
        onClose={handleCloseMenu}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
      >
        {user && (
          <>
            <MenuItem onClick={handleOpenProfile}>
              <ListItemIcon>
                <AccountCircle />
              </ListItemIcon>
              <Typography variant="inherit">{user.name}</Typography>
            </MenuItem>

            <MenuItem component={Link} to="/favorites" onClick={handleCloseMenu}>
              <ListItemIcon>
                <Favorite />
              </ListItemIcon>
              <Typography variant="inherit">{t('navigation.favorites', { defaultValue: 'Избранное', ns: 'marketplace' })}</Typography>
            </MenuItem>

            {renderMessagesMenuItem()}

            <MenuItem onClick={() => { setNotificationDrawerOpen(true); handleCloseMenu(); }}>
              <ListItemIcon>
                <NotificationsIcon />
              </ListItemIcon>
              <Typography variant="inherit">{t('navigation.notifications', { defaultValue: 'Оповещения', ns: 'marketplace' })}</Typography>
            </MenuItem>

            <MenuItem component={Link} to="/storefronts" onClick={handleCloseMenu}>
              <ListItemIcon>
                <Store />
              </ListItemIcon>
              <Typography variant="inherit">{t('navigation.storefronts', { defaultValue: 'Мои витрины', ns: 'marketplace' })}</Typography>
            </MenuItem>

            <MenuItem component={Link} to="/balance" onClick={handleCloseMenu}>
              <ListItemIcon>
                <AccountBalanceWallet />
              </ListItemIcon>
              <Typography variant="inherit">{t('navigation.balance', { defaultValue: 'Мой баланс', ns: 'marketplace' })}</Typography>
            </MenuItem>


            {isAdmin(user.email) && (
                <MenuItem component={Link} to="/admin" onClick={handleCloseMenu}>
                  <Settings fontSize="small" sx={{ mr: 1 }} />
                  Администрирование
                </MenuItem>
            )}

            <Divider />

            <MenuItem onClick={() => { logout(); handleCloseMenu(); }}>
              <ListItemIcon>
                <Logout />
              </ListItemIcon>
              <Typography variant="inherit">{t('auth.signout', { defaultValue: 'Выйти', ns: 'marketplace' })}</Typography>
            </MenuItem>
          </>
        )}

        {!user && (
          <MenuItem onClick={() => { login(); handleCloseMenu(); }}>
            <ListItemIcon>
              <Login />
            </ListItemIcon>
            <Typography variant="inherit">{t('auth.signin', { defaultValue: 'Войти', ns: 'marketplace' })}</Typography>
          </MenuItem>
        )}
      </Menu>

      {/* Меню выбора языка */}
      <Menu
        anchorEl={languageAnchorEl}
        open={Boolean(languageAnchorEl)}
        onClose={handleCloseLanguageMenu}
        PaperProps={{
          elevation: 0,
          sx: {
            overflow: 'visible',
            filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
            mt: 1.5,
          },
        }}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
      >
        <MenuItem onClick={() => handleLanguageChange('en')} selected={i18n.language === 'en'}>
          <Typography>English</Typography>
        </MenuItem>
        <MenuItem onClick={() => handleLanguageChange('sr')} selected={i18n.language === 'sr'}>
          <Typography>Srpski</Typography>
        </MenuItem>
        <MenuItem onClick={() => handleLanguageChange('ru')} selected={i18n.language === 'ru'}>
          <Typography>{t('languages.russian', { defaultValue: 'Русский', ns: 'marketplace' })}</Typography>
        </MenuItem>
      </Menu>
    </Box>
  );
};

export default Layout;
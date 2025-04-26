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
  Login,
  Devices,
  Home,
  Grass,
  SportsSoccer,
  Pets,
  ChildFriendly,
  ExpandMore,
  ExpandLess,
  KeyboardArrowRight,
  Work, // Добавьте этот импорт
  Checkroom,
  Security,
  MoreHoriz
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
  ListItem,
  Collapse
} from "@mui/material";

// Компонент меню категорий для десктопной версии
// Обновленный компонент CategoryMenu для многоуровневой навигации
const CategoryMenu = ({ categories, onSelect, anchorEl, open, onClose }) => {
  const [expandedGroups, setExpandedGroups] = useState({});
  const [expandedSubItems, setExpandedSubItems] = useState({});

  const handleToggleGroup = (key) => {
    setExpandedGroups(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
  };

  const handleToggleSubItem = (path) => {
    setExpandedSubItems(prev => ({
      ...prev,
      [path]: !prev[path]
    }));
  };

  // Количество категорий в первой строке
  const firstRowCount = 5;
  const topCategories = categories.slice(0, firstRowCount);
  const moreCategories = categories.slice(firstRowCount);

  return (
      <Menu
          anchorEl={anchorEl}
          open={open}
          onClose={onClose}
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
              width: '900px', // Увеличиваем ширину
              padding: '16px',
              overflow: 'auto'
            },
          }}
      >
        <Grid container spacing={2}>
          {/* Первая строка с основными категориями */}
          {topCategories.map((group, index) => (
              <Grid item xs={12 / firstRowCount} key={group.key}>
                <Box sx={{ mb: 2 }}>
                  {/* Заголовок категории */}
                  <Button
                      onClick={() => onSelect(group.path)}
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
                        '&:hover': { backgroundColor: 'rgba(25, 118, 210, 0.04)' }
                      }}
                  >
                    <Box sx={{ color: 'primary.main', mr: 1 }}>
                      {group.icon}
                    </Box>
                    <Typography variant="subtitle1" noWrap>
                      {group.title}
                    </Typography>
                  </Button>

                  {/* Подкатегории */}
                  <Box sx={{ pl: 2 }}>
                    {group.subcategories.slice(0, 6).map((subcat, subIdx) => (
                        <React.Fragment key={subIdx}>
                          <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                            <Button
                                onClick={() => onSelect(subcat.path)}
                                sx={{
                                  display: 'flex',
                                  justifyContent: 'space-between',
                                  alignItems: 'center',
                                  width: subcat.subItems ? 'calc(100% - 30px)' : '100%',
                                  textAlign: 'left',
                                  py: 0.5,
                                  px: 1,
                                  color: 'text.primary',
                                  textTransform: 'none',
                                  fontSize: '0.85rem',
                                  '&:hover': { backgroundColor: 'action.hover' }
                                }}
                            >
                              <Typography variant="body2" noWrap>
                                {subcat.name}
                              </Typography>
                            </Button>
                            {subcat.subItems && (
                                <IconButton
                                    size="small"
                                    onClick={(e) => {
                                      e.stopPropagation();
                                      handleToggleSubItem(subcat.path);
                                    }}
                                    sx={{ p: 0, ml: 0.5, minWidth: 24 }}
                                >
                                  {expandedSubItems[subcat.path] ? <ExpandLess fontSize="small" /> : <ExpandMore fontSize="small" />}
                                </IconButton>
                            )}
                          </Box>
                          {/* Третий уровень категорий */}
                          {subcat.subItems && (
                              <Collapse in={expandedSubItems[subcat.path]} timeout="auto" unmountOnExit>
                                <Box sx={{ pl: 2 }}>
                                  {subcat.subItems.map((subItem, subItemIdx) => (
                                      <Button
                                          key={subItemIdx}
                                          onClick={() => onSelect(subItem.path)}
                                          sx={{
                                            display: 'block',
                                            width: '100%',
                                            textAlign: 'left',
                                            justifyContent: 'flex-start',
                                            py: 0.5,
                                            px: 1,
                                            color: 'text.secondary',
                                            textTransform: 'none',
                                            fontSize: '0.8rem',
                                            '&:hover': { backgroundColor: 'action.hover' }
                                          }}
                                      >
                                        <Typography variant="body2" noWrap>
                                          {subItem.name}
                                        </Typography>
                                      </Button>
                                  ))}
                                </Box>
                              </Collapse>
                          )}
                        </React.Fragment>
                    ))}

                    {/* Кнопка "Показать еще" - отображается только если есть дополнительные подкатегории и они не развернуты */}
                    {group.subcategories.length > 6 && !expandedGroups[group.key] && (
                        <Button
                            onClick={() => handleToggleGroup(group.key)}
                            sx={{
                              display: 'flex',
                              width: '100%',
                              textAlign: 'left',
                              justifyContent: 'space-between',
                              alignItems: 'center',
                              py: 0.5,
                              px: 1,
                              color: 'primary.main',
                              textTransform: 'none',
                              fontSize: '0.85rem',
                              '&:hover': { backgroundColor: 'action.hover' }
                            }}
                        >
                          <Typography variant="body2">
                            Показать еще
                          </Typography>
                          <ExpandMore fontSize="small" />
                        </Button>
                    )}


                    <Collapse in={!!expandedGroups[group.key]} timeout="auto" unmountOnExit>
                      <Box sx={{ pl: 0 }}>
                        {group.subcategories.slice(6).map((subcat, subIdx) => (
                            <React.Fragment key={subIdx}>
                              <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                                <Button
                                    onClick={() => onSelect(subcat.path)}
                                    sx={{
                                      display: 'flex',
                                      justifyContent: 'space-between',
                                      alignItems: 'center',
                                      width: subcat.subItems ? 'calc(100% - 30px)' : '100%',
                                      textAlign: 'left',
                                      py: 0.5,
                                      px: 1,
                                      color: 'text.primary',
                                      textTransform: 'none',
                                      fontSize: '0.85rem',
                                      '&:hover': { backgroundColor: 'action.hover' }
                                    }}
                                >
                                  <Typography variant="body2" noWrap>
                                    {subcat.name}
                                  </Typography>
                                </Button>
                                {subcat.subItems && (
                                    <IconButton
                                        size="small"
                                        onClick={(e) => {
                                          e.stopPropagation();
                                          handleToggleSubItem(subcat.path);
                                        }}
                                        sx={{ p: 0, ml: 0.5, minWidth: 24 }}
                                    >
                                      {expandedSubItems[subcat.path] ? <ExpandLess fontSize="small" /> : <ExpandMore fontSize="small" />}
                                    </IconButton>
                                )}
                              </Box>
                              {/* Третий уровень категорий */}
                              {subcat.subItems && (
                                  <Collapse in={expandedSubItems[subcat.path]} timeout="auto" unmountOnExit>
                                    <Box sx={{ pl: 2 }}>
                                      {subcat.subItems.map((subItem, subItemIdx) => (
                                          <Button
                                              key={subItemIdx}
                                              onClick={() => onSelect(subItem.path)}
                                              sx={{
                                                display: 'block',
                                                width: '100%',
                                                textAlign: 'left',
                                                justifyContent: 'flex-start',
                                                py: 0.5,
                                                px: 1,
                                                color: 'text.secondary',
                                                textTransform: 'none',
                                                fontSize: '0.8rem',
                                                '&:hover': { backgroundColor: 'action.hover' }
                                              }}
                                          >
                                            <Typography variant="body2" noWrap>
                                              {subItem.name}
                                            </Typography>
                                          </Button>
                                      ))}
                                    </Box>
                                  </Collapse>
                              )}
                            </React.Fragment>
                        ))}

                        {/* Кнопка "Свернуть" - отображается только когда список развернут, и находится в конце развернутого списка */}
                        {expandedGroups[group.key] && (
                            <Button
                                onClick={() => handleToggleGroup(group.key)}
                                sx={{
                                  display: 'flex',
                                  width: '100%',
                                  textAlign: 'left',
                                  justifyContent: 'space-between',
                                  alignItems: 'center',
                                  py: 0.5,
                                  px: 1,
                                  color: 'primary.main',
                                  textTransform: 'none',
                                  fontSize: '0.85rem',
                                  '&:hover': { backgroundColor: 'action.hover' }
                                }}
                            >
                              <Typography variant="body2">
                                Свернуть
                              </Typography>
                              <ExpandLess fontSize="small" />
                            </Button>
                        )}
                      </Box>
                    </Collapse>

                  </Box>
                </Box>
              </Grid>
          ))}
        </Grid>

        {/* Остальные категории в компактном виде */}
        {moreCategories.length > 0 && (
            <>
              <Divider sx={{ my: 2 }} />
              <Typography variant="subtitle1" sx={{ fontWeight: 'bold', mb: 1, px: 1 }}>
                Больше категорий
              </Typography>
              <Grid container spacing={2}>
                {moreCategories.map((group, index) => (
                    <Grid item xs={12/3} key={group.key}>
                      <Button
                          onClick={() => onSelect(group.path)}
                          sx={{
                            display: 'flex',
                            alignItems: 'center',
                            width: '100%',
                            textAlign: 'left',
                            justifyContent: 'flex-start',
                            color: 'text.primary',
                            textTransform: 'none',
                            py: 0.5,
                            '&:hover': { backgroundColor: 'action.hover' }
                          }}
                      >
                        <Box sx={{ color: 'primary.main', mr: 1 }}>
                          {group.icon}
                        </Box>
                        <Typography variant="body2" noWrap>
                          {group.title}
                        </Typography>
                        <KeyboardArrowRight fontSize="small" sx={{ ml: 'auto' }} />
                      </Button>
                    </Grid>
                ))}
              </Grid>
            </>
        )}
      </Menu>
  );
};

// Обновленный компонент для мобильной версии с поддержкой многоуровневой вложенности
const MobileCategoryDrawer = ({ categories, onSelect, open, onClose, t }) => {
  const [expandedCategory, setExpandedCategory] = useState(null);
  const [expandedSubcategory, setExpandedSubcategory] = useState(null);

  const handleToggleCategory = (key) => {
    if (expandedCategory === key) {
      setExpandedCategory(null);
      setExpandedSubcategory(null); // Сбрасываем раскрытые подкатегории при закрытии категории
    } else {
      setExpandedCategory(key);
      setExpandedSubcategory(null);
    }
  };

  const handleToggleSubcategory = (path, event) => {
    event.stopPropagation();
    if (expandedSubcategory === path) {
      setExpandedSubcategory(null);
    } else {
      setExpandedSubcategory(path);
    }
  };

  return (
      <Drawer
          anchor="right"
          open={open}
          onClose={onClose}
          sx={{
            '& .MuiDrawer-paper': {
              width: { xs: '85%', sm: '70%' },
              maxWidth: '400px',
            }
          }}
      >
        <Box sx={{ p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid', borderColor: 'divider' }}>
          <Typography variant="h6">{t('categories.title', { defaultValue: 'Категории', ns: 'marketplace' })}</Typography>
          <IconButton onClick={onClose}>
            <X size={20} />
          </IconButton>
        </Box>

        <Box sx={{ overflowY: 'auto', flex: 1 }}>
          {categories.map((group) => (
              <React.Fragment key={group.key}>
                <Box
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      px: 2,
                      py: 1.5,
                      borderBottom: '1px solid',
                      borderColor: 'divider',
                      backgroundColor: expandedCategory === group.key ? 'action.selected' : 'transparent'
                    }}
                >
                  <IconButton
                      size="small"
                      onClick={() => handleToggleCategory(group.key)}
                      sx={{ mr: 1 }}
                  >
                    {expandedCategory === group.key ? <ExpandLess /> : <ExpandMore />}
                  </IconButton>
                  <Box
                      component={Button}
                      onClick={() => {
                        onSelect(group.path);
                        onClose();
                      }}
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        flex: 1,
                        color: 'text.primary',
                        fontWeight: expandedCategory === group.key ? 'bold' : 'normal',
                        textTransform: 'none',
                        justifyContent: 'flex-start',
                        px: 0,
                        '&:hover': { backgroundColor: 'transparent' }
                      }}
                  >
                    <Box sx={{ color: 'primary.main', mr: 2 }}>
                      {group.icon}
                    </Box>
                    <Typography variant="subtitle1">{group.title}</Typography>
                  </Box>
                </Box>

                <Collapse in={expandedCategory === group.key} timeout="auto" unmountOnExit>
                  <Box sx={{ bgcolor: 'background.default' }}>
                    {group.subcategories.map((subcat, idx) => (
                        <React.Fragment key={idx}>
                          <Box
                              sx={{
                                display: 'flex',
                                alignItems: 'center',
                                borderBottom: '1px solid',
                                borderColor: 'divider',
                                backgroundColor: expandedSubcategory === subcat.path ? 'action.hover' : 'transparent'
                              }}
                          >
                            {/* Индикатор наличия подкатегорий */}
                            {subcat.subItems ? (
                                <IconButton
                                    size="small"
                                    onClick={(e) => handleToggleSubcategory(subcat.path, e)}
                                    sx={{ ml: 3, mr: 1 }}
                                >
                                  {expandedSubcategory === subcat.path ? <ExpandLess fontSize="small" /> : <ExpandMore fontSize="small" />}
                                </IconButton>
                            ) : (
                                <Box sx={{ ml: 5, width: 24 }} /> // Отступ для выравнивания
                            )}

                            <Button
                                onClick={() => {
                                  onSelect(subcat.path);
                                  onClose();
                                }}
                                sx={{
                                  display: 'flex',
                                  width: '100%',
                                  textAlign: 'left',
                                  justifyContent: 'flex-start',
                                  py: 1.5,
                                  px: 1,
                                  color: 'text.primary',
                                  textTransform: 'none',
                                  fontSize: '0.95rem',
                                  '&:hover': { backgroundColor: 'action.hover' }
                                }}
                            >
                              <Typography variant="body2">
                                {subcat.name}
                              </Typography>
                            </Button>
                          </Box>

                          {/* Третий уровень категорий */}
                          {subcat.subItems && (
                              <Collapse in={expandedSubcategory === subcat.path} timeout="auto" unmountOnExit>
                                <Box sx={{ bgcolor: 'background.paper' }}>
                                  {subcat.subItems.map((subItem, subItemIdx) => (
                                      <Button
                                          key={subItemIdx}
                                          onClick={() => {
                                            onSelect(subItem.path);
                                            onClose();
                                          }}
                                          sx={{
                                            display: 'block',
                                            width: '100%',
                                            textAlign: 'left',
                                            justifyContent: 'flex-start',
                                            py: 1.2,
                                            pl: 9, // Увеличенный отступ для визуальной иерархии
                                            pr: 2,
                                            color: 'text.secondary',
                                            textTransform: 'none',
                                            borderBottom: '1px solid',
                                            borderColor: 'divider',
                                            '&:hover': { backgroundColor: 'action.hover' }
                                          }}
                                      >
                                        <Typography variant="body2">
                                          {subItem.name}
                                        </Typography>
                                      </Button>
                                  ))}
                                </Box>
                              </Collapse>
                          )}
                        </React.Fragment>
                    ))}
                  </Box>
                </Collapse>
              </React.Fragment>
          ))}
        </Box>
      </Drawer>
  );
};

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

// Обновляем массив categoryGroups, добавляя все категории и учитывая глубокую вложенность
  const categoryGroups = [
    {
      title: t('categories.realestate', { defaultValue: 'Недвижимость', ns: 'marketplace' }),
      icon: <HomeWork />,
      path: "/marketplace?category_id=1000",
      key: "realestate",
      subcategories: [
        { name: t('categories.apartments', { defaultValue: 'Квартиры', ns: 'marketplace' }), path: "/marketplace?category_id=1100" },
        { name: t('categories.rooms', { defaultValue: 'Комнаты', ns: 'marketplace' }), path: "/marketplace?category_id=1200" },
        { name: t('categories.houses', { defaultValue: 'Дома, дачи, коттеджи', ns: 'marketplace' }), path: "/marketplace?category_id=1300" },
        { name: t('categories.land', { defaultValue: 'Земельные участки', ns: 'marketplace' }), path: "/marketplace?category_id=1400" },
        { name: t('categories.garages', { defaultValue: 'Гаражи', ns: 'marketplace' }), path: "/marketplace?category_id=1500" },
        { name: t('categories.commercial', { defaultValue: 'Коммерческая', ns: 'marketplace' }), path: "/marketplace?category_id=1600" },
        { name: t('categories.foreign', { defaultValue: 'Недвижимость за рубежом', ns: 'marketplace' }), path: "/marketplace?category_id=1700" },
        { name: t('categories.hotel', { defaultValue: 'Отель', ns: 'marketplace' }), path: "/marketplace?category_id=1800" },
        { name: t('categories.apartment_property', { defaultValue: 'Апартаменты', ns: 'marketplace' }), path: "/marketplace?category_id=1900" },
      ]
    },
    {
      title: t('categories.transport', { defaultValue: 'Автомобили', ns: 'marketplace' }),
      icon: <DirectionsCar />,
      path: "/marketplace?category_id=2000",
      key: "transport",
      subcategories: [
        { name: t('categories.passenger', { defaultValue: 'Легковые автомобили', ns: 'marketplace' }), path: "/marketplace?category_id=2100" },
        {
          name: t('categories.commercial', { defaultValue: 'Грузовые автомобили', ns: 'marketplace' }),
          path: "/marketplace?category_id=2200",
          subItems: [
            { name: t('categories.trucks', { defaultValue: 'Грузовики', ns: 'marketplace' }), path: "/marketplace?category_id=2210" },
            { name: t('categories.semitrailers', { defaultValue: 'Полуприцепы', ns: 'marketplace' }), path: "/marketplace?category_id=2220" },
            { name: t('categories.light_commercial', { defaultValue: 'Лёгкий коммерческий транспорт', ns: 'marketplace' }), path: "/marketplace?category_id=2230" },
            { name: t('categories.buses', { defaultValue: 'Автобусы', ns: 'marketplace' }), path: "/marketplace?category_id=2240" },
          ]
        },
        {
          name: t('categories.special', { defaultValue: 'Спецтехника', ns: 'marketplace' }),
          path: "/marketplace?category_id=2300",
          subItems: [
            { name: t('categories.excavators', { defaultValue: 'Экскаваторы', ns: 'marketplace' }), path: "/marketplace?category_id=2310" },
            { name: t('categories.loaders', { defaultValue: 'Погрузчики', ns: 'marketplace' }), path: "/marketplace?category_id=2315" },
            { name: t('categories.backhoe_loaders', { defaultValue: 'Экскаваторы-погрузчики', ns: 'marketplace' }), path: "/marketplace?category_id=2320" },
            // Остальные подкатегории спецтехники по необходимости
          ]
        },
        {
          name: t('categories.agricultural', { defaultValue: 'Сельхозтехника', ns: 'marketplace' }),
          path: "/marketplace?category_id=2400",
          subItems: [
            { name: t('categories.tractors', { defaultValue: 'Тракторы', ns: 'marketplace' }), path: "/marketplace?category_id=2410" },
            { name: t('categories.mini_tractors', { defaultValue: 'Мини-тракторы', ns: 'marketplace' }), path: "/marketplace?category_id=2415" },
            // Остальные подкатегории сельхозтехники по необходимости
          ]
        },
        { name: t('categories.rent', { defaultValue: 'Аренда авто и спецтехники', ns: 'marketplace' }), path: "/marketplace?category_id=2500" },
        { name: t('categories.motorcycles', { defaultValue: 'Мотоциклы', ns: 'marketplace' }), path: "/marketplace?category_id=2600" },
        { name: t('categories.water', { defaultValue: 'Водный транспорт', ns: 'marketplace' }), path: "/marketplace?category_id=2700" },
        { name: t('categories.parts', { defaultValue: 'Запчасти', ns: 'marketplace' }), path: "/marketplace?category_id=2800" },
      ]
    },
    {
      title: t('categories.electronics', { defaultValue: 'Электроника', ns: 'marketplace' }),
      icon: <Devices />,
      path: "/marketplace?category_id=3000",
      key: "electronics",
      subcategories: [
        {
          name: t('categories.phones', { defaultValue: 'Телефоны', ns: 'marketplace' }),
          path: "/marketplace?category_id=3100",
          subItems: [
            { name: t('categories.mobile_phones', { defaultValue: 'Мобильные телефоны', ns: 'marketplace' }), path: "/marketplace?category_id=3110" },
            { name: t('categories.accessories', { defaultValue: 'Аксессуары', ns: 'marketplace' }), path: "/marketplace?category_id=3120" },
            { name: t('categories.radio', { defaultValue: 'Рации', ns: 'marketplace' }), path: "/marketplace?category_id=3130" },
            { name: t('categories.landline', { defaultValue: 'Стационарные телефоны', ns: 'marketplace' }), path: "/marketplace?category_id=3140" },
          ]
        },
        { name: t('categories.audio_video', { defaultValue: 'Аудио и видео', ns: 'marketplace' }), path: "/marketplace?category_id=3200" },
        {
          name: t('categories.computers', { defaultValue: 'Товары для компьютера', ns: 'marketplace' }),
          path: "/marketplace?category_id=3300",
          subItems: [
            { name: t('categories.desktops', { defaultValue: 'Системные блоки', ns: 'marketplace' }), path: "/marketplace?category_id=3310" },
            { name: t('categories.all_in_one', { defaultValue: 'Моноблоки', ns: 'marketplace' }), path: "/marketplace?category_id=3320" },
            { name: t('categories.components', { defaultValue: 'Комплектующие', ns: 'marketplace' }), path: "/marketplace?category_id=3330" },
            // Остальные подкатегории по необходимости
          ]
        },
        { name: t('categories.games', { defaultValue: 'Игры и приставки', ns: 'marketplace' }), path: "/marketplace?category_id=3500" },
        { name: t('categories.laptops', { defaultValue: 'Ноутбуки', ns: 'marketplace' }), path: "/marketplace?category_id=3600" },
        { name: t('categories.photo', { defaultValue: 'Фототехника', ns: 'marketplace' }), path: "/marketplace?category_id=3700" },
        { name: t('categories.tablets', { defaultValue: 'Планшеты и электронные книги', ns: 'marketplace' }), path: "/marketplace?category_id=3800" },
        { name: t('categories.office', { defaultValue: 'Оргтехника и расходники', ns: 'marketplace' }), path: "/marketplace?category_id=3900" },
        { name: t('categories.appliances', { defaultValue: 'Бытовая техника', ns: 'marketplace' }), path: "/marketplace?category_id=4100" },
      ]
    },
    {
      title: t('categories.forHome', { defaultValue: 'Для дома', ns: 'marketplace' }),
      icon: <Home />,
      path: "/marketplace?category_id=5000",
      key: "home",
      subcategories: [
        { name: t('categories.repair', { defaultValue: 'Ремонт и строительство', ns: 'marketplace' }), path: "/marketplace?category_id=5100" },
        { name: t('categories.furniture', { defaultValue: 'Мебель и интерьер', ns: 'marketplace' }), path: "/marketplace?category_id=5200" },
        { name: t('categories.food', { defaultValue: 'Продукты питания', ns: 'marketplace' }), path: "/marketplace?category_id=5300" },
        { name: t('categories.kitchenware', { defaultValue: 'Посуда и товары для кухни', ns: 'marketplace' }), path: "/marketplace?category_id=5400" },
      ]
    },
    {
      title: t('categories.forGarden', { defaultValue: 'Для сада', ns: 'marketplace' }),
      icon: <Grass />,
      path: "/marketplace?category_id=6000",
      key: "garden",
      subcategories: [
        { name: t('categories.gardenFurniture', { defaultValue: 'Садовая мебель', ns: 'marketplace' }), path: "/marketplace?category_id=6050" },
        { name: t('categories.tools', { defaultValue: 'Садовые инструменты', ns: 'marketplace' }), path: "/marketplace?category_id=6100" },
        { name: t('categories.seeds', { defaultValue: 'Семена и рассада', ns: 'marketplace' }), path: "/marketplace?category_id=6200" },
        { name: t('categories.bbq', { defaultValue: 'Барбекю и аксессуары', ns: 'marketplace' }), path: "/marketplace?category_id=6250" },
        { name: t('categories.plants', { defaultValue: 'Садовые растения', ns: 'marketplace' }), path: "/marketplace?category_id=6750" },
        // Добавьте остальные подкатегории для сада
      ]
    },
    {
      title: t('categories.hobby', { defaultValue: 'Хобби и отдых', ns: 'marketplace' }),
      icon: <SportsSoccer />,
      path: "/marketplace?category_id=7000",
      key: "hobby",
      subcategories: [
        { name: t('categories.music', { defaultValue: 'Музыкальные инструменты', ns: 'marketplace' }), path: "/marketplace?category_id=7050" },
        { name: t('categories.books', { defaultValue: 'Книги и журналы', ns: 'marketplace' }), path: "/marketplace?category_id=7100" },
        { name: t('categories.sport_equipment', { defaultValue: 'Спортивный инвентарь', ns: 'marketplace' }), path: "/marketplace?category_id=7150" },
        { name: t('categories.collecting', { defaultValue: 'Коллекционирование', ns: 'marketplace' }), path: "/marketplace?category_id=7250" },
        { name: t('categories.art', { defaultValue: 'Предметы искусства', ns: 'marketplace' }), path: "/marketplace?category_id=7300" },
        { name: t('categories.bicycles', { defaultValue: 'Велосипеды', ns: 'marketplace' }), path: "/marketplace?category_id=7400" },
        { name: t('categories.hunting', { defaultValue: 'Охота и рыбалка', ns: 'marketplace' }), path: "/marketplace?category_id=7500" },
        { name: t('categories.camping', { defaultValue: 'Кемпинг', ns: 'marketplace' }), path: "/marketplace?category_id=7650" },
        { name: t('categories.antiques', { defaultValue: 'Антиквариат', ns: 'marketplace' }), path: "/marketplace?category_id=7700" },
        { name: t('categories.tickets', { defaultValue: 'Билеты и путешествия', ns: 'marketplace' }), path: "/marketplace?category_id=7750" },
        { name: t('categories.sport', { defaultValue: 'Спорт', ns: 'marketplace' }), path: "/marketplace?category_id=7800" },
        { name: t('categories.crafts', { defaultValue: 'Народное ремесло', ns: 'marketplace' }), path: "/marketplace?category_id=7865" },
        { name: t('categories.beekeeping', { defaultValue: 'Пчеловодство', ns: 'marketplace' }), path: "/marketplace?category_id=7900" },
        { name: t('categories.rural_tourism', { defaultValue: 'Сельский туризм', ns: 'marketplace' }), path: "/marketplace?category_id=7950" },
      ]
    },
    {
      title: t('categories.animals', { defaultValue: 'Животные', ns: 'marketplace' }),
      icon: <Pets />,
      path: "/marketplace?category_id=8000",
      key: "animals",
      subcategories: [
        { name: t('categories.dogs', { defaultValue: 'Собаки', ns: 'marketplace' }), path: "/marketplace?category_id=8050" },
        { name: t('categories.cats', { defaultValue: 'Кошки', ns: 'marketplace' }), path: "/marketplace?category_id=8100" },
        { name: t('categories.birds', { defaultValue: 'Птицы', ns: 'marketplace' }), path: "/marketplace?category_id=8150" },
        { name: t('categories.aquarium', { defaultValue: 'Аквариум', ns: 'marketplace' }), path: "/marketplace?category_id=8200" },
        { name: t('categories.otherAnimals', { defaultValue: 'Другие животные', ns: 'marketplace' }), path: "/marketplace?category_id=8250" },
      ]
    },
    {
      title: t('categories.business', { defaultValue: 'Готовый бизнес и оборудование', ns: 'marketplace' }),
      icon: <Business />,
      path: "/marketplace?category_id=8500",
      key: "business",
      subcategories: []
    },
    {
      title: t('categories.jobs', { defaultValue: 'Работа', ns: 'marketplace' }),
      icon: <Work />,
      path: "/marketplace?category_id=9000",
      key: "jobs",
      subcategories: [
        { name: t('categories.vacancies', { defaultValue: 'Вакансии', ns: 'marketplace' }), path: "/marketplace?category_id=9050" },
        { name: t('categories.resumes', { defaultValue: 'Резюме', ns: 'marketplace' }), path: "/marketplace?category_id=9100" },
        { name: t('categories.remote_work', { defaultValue: 'Удаленная работа', ns: 'marketplace' }), path: "/marketplace?category_id=9150" },
        { name: t('categories.partnership', { defaultValue: 'Партнерство', ns: 'marketplace' }), path: "/marketplace?category_id=9200" },
        { name: t('categories.training', { defaultValue: 'Обучение и стажировка', ns: 'marketplace' }), path: "/marketplace?category_id=9250" },
        {
          name: t('categories.seasonal', { defaultValue: 'Сезонные работы', ns: 'marketplace' }),
          path: "/marketplace?category_id=9300",
          subItems: [
            { name: t('categories.harvesting', { defaultValue: 'Сбор урожая', ns: 'marketplace' }), path: "/marketplace?category_id=9310" },
            { name: t('categories.vineyard', { defaultValue: 'Работа на винограднике', ns: 'marketplace' }), path: "/marketplace?category_id=9315" },
            { name: t('categories.construction_seasonal', { defaultValue: 'Сезонные строительные работы', ns: 'marketplace' }), path: "/marketplace?category_id=9320" },
          ]
        },
      ]
    },
    {
      title: t('categories.fashion', { defaultValue: 'Одежда, обувь, аксессуары', ns: 'marketplace' }),
      icon: <Checkroom />,
      path: "/marketplace?category_id=9500",
      key: "fashion",
      subcategories: []
    },
    {
      title: t('categories.forChildren', { defaultValue: 'Товары для детей', ns: 'marketplace' }),
      icon: <ChildFriendly />,
      path: "/marketplace?category_id=9700",
      key: "children",
      subcategories: [
        { name: t('categories.strollers', { defaultValue: 'Детские коляски', ns: 'marketplace' }), path: "/marketplace?category_id=9705" },
        { name: t('categories.childrenFurniture', { defaultValue: 'Детская мебель', ns: 'marketplace' }), path: "/marketplace?category_id=9710" },
        { name: t('categories.bikes_scooters', { defaultValue: 'Велосипеды и самокаты', ns: 'marketplace' }), path: "/marketplace?category_id=9715" },
        { name: t('categories.feeding', { defaultValue: 'Товары для кормления', ns: 'marketplace' }), path: "/marketplace?category_id=9720" },
        { name: t('categories.carSeats', { defaultValue: 'Автомобильные кресла', ns: 'marketplace' }), path: "/marketplace?category_id=9725" },
        { name: t('categories.toys', { defaultValue: 'Игрушки', ns: 'marketplace' }), path: "/marketplace?category_id=9730" },
        { name: t('categories.childrenClothes', { defaultValue: 'Детская одежда и обувь', ns: 'marketplace' }), path: "/marketplace?category_id=9750" },
      ]
    },
    {
      title: t('categories.security', { defaultValue: 'Безопасность', ns: 'marketplace' }),
      icon: <Security />,
      path: "/marketplace?category_id=10000",
      key: "security",
      subcategories: [
        { name: t('categories.video_surveillance', { defaultValue: 'Видеонаблюдение', ns: 'marketplace' }), path: "/marketplace?category_id=10100" },
      ]
    },
    {
      title: t('categories.other', { defaultValue: 'Прочее', ns: 'marketplace' }),
      icon: <MoreHoriz />,
      path: "/marketplace?category_id=9999",
      key: "other",
      subcategories: []
    },
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

    // Обновляем поисковые параметры, если необходимо
    if (path.includes('?')) {
      const searchParamsString = path.split('?')[1];
      const newSearchParams = new URLSearchParams(searchParamsString);
      setSearchParams(newSearchParams);
    }
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
                  minHeight: "48px",
                  px: 1, py: 0,
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
{/*               <Typography
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
*/}              
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
              py: 0
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

        {/* Новый компонент меню категорий для десктопа */}
        <CategoryMenu
            categories={categoryGroups}
            onSelect={handleCategoryClick}
            anchorEl={categoryMenuAnchor}
            open={Boolean(categoryMenuAnchor)}
            onClose={handleCloseCategoryMenu}
        />

        {/* Обновленный мобильный drawer для категорий */}
        <MobileCategoryDrawer
            categories={categoryGroups}
            onSelect={handleCategoryClick}
            open={mobileCategoryDrawerOpen}
            onClose={handleCloseMobileCategoryDrawer}
            t={t}
        />

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
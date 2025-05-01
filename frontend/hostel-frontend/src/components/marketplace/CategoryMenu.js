// frontend/hostel-frontend/src/components/marketplace/CategoryMenu.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import axios from '../../api/axios';
import {
  Box,
  Button,
  Menu,
  MenuItem,
  Grid,
  Typography,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  useTheme,
  useMediaQuery,
  Collapse,
  Divider,
  Link
} from '@mui/material';
import {
  KeyboardArrowDown,
  ChevronRight,
  ChevronDown,
  Close as CloseIcon,
  ExpandMore,
  ExpandLess,
  ArrowBack,
  HomeWork,
  DirectionsCar,
  ShoppingBag,
  Apartment,
  Weekend,
  Devices,
  Category as CategoryIcon
} from '@mui/icons-material';

// Функция для получения переведенного имени категории
const getTranslatedName = (category, language) => {
  if (category.translations && category.translations[language]) {
    return category.translations[language];
  }
  return category.name;
};

// Иконки для категорий верхнего уровня
const CATEGORY_ICONS = {
  1000: <HomeWork />,          // Недвижимость
  2000: <DirectionsCar />,     // Транспорт
  3000: <Devices />,           // Электроника
  4000: <ShoppingBag />,       // Личные вещи
  5000: <Weekend />,           // Для дома
  6000: <Apartment />          // Для сада
};

// Получение иконки для категории
const getCategoryIcon = (categoryId) => {
  // Преобразуем ID в число и округляем до ближайшей тысячи вниз
  const baseId = Math.floor(parseInt(categoryId) / 1000) * 1000;
  return CATEGORY_ICONS[baseId] || <ShoppingBag />;
};

const CategoryMenu = () => {
  const { t, i18n } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const navigate = useNavigate();

  const [categoryMenuAnchor, setCategoryMenuAnchor] = useState(null);
  const [mobileCategoryDrawerOpen, setMobileCategoryDrawerOpen] = useState(false);
  const [categories, setCategories] = useState([]);
  const [categoryGroups, setCategoryGroups] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [expandedGroups, setExpandedGroups] = useState({});
  const [expandedMobileCategories, setExpandedMobileCategories] = useState({});
  const [mobileNavigationStack, setMobileNavigationStack] = useState([]);
  const [currentMobileCategory, setCurrentMobileCategory] = useState(null);

  // Константа для ограничения количества подкатегорий
  const MAX_VISIBLE_SUBCATEGORIES = 5;

  // Порядок категорий слева направо
  const categoryOrder = [
    2000, 1000, 9000,
    3000, 9500, 7000,
    8000, 8500, 5000,
    6000, 9700, 10000,
    9999
  ];

  // Функция для сортировки категорий по заданному порядку
  const sortCategoriesByOrder = (categories) => {
    return [...categories].sort((a, b) => {
      const indexA = categoryOrder.indexOf(a.id);
      const indexB = categoryOrder.indexOf(b.id);

      // Если обе категории есть в массиве порядка, сортируем по их позиции
      if (indexA !== -1 && indexB !== -1) {
        return indexA - indexB;
      }

      // Если только одна категория есть в массиве порядка, она идет первой
      if (indexA !== -1) return -1;
      if (indexB !== -1) return 1;

      // Если ни одной категории нет в массиве порядка, сортируем по ID
      return a.id - b.id;
    });
  };

  // Загрузка категорий
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setIsLoading(true);
        const response = await axios.get('/api/v1/marketplace/category-tree');
        if (response.data?.data) {
          const allCategories = response.data.data;
          setCategories(allCategories);

          // Группируем категории верхнего уровня
          const topLevelCategories = allCategories.filter(cat => !cat.parent_id);

          // Сортируем категории по заданному порядку
          const sortedTopLevelCategories = sortCategoriesByOrder(topLevelCategories);

          // Создаем группы категорий для меню
          const groups = sortedTopLevelCategories.map(category => {
            // Находим подкатегории для текущей категории верхнего уровня
            const subcategories = allCategories
              .filter(cat => cat.parent_id === category.id)
              .map(subcat => ({
                id: subcat.id,
                name: getTranslatedName(subcat, i18n.language),
                path: `/marketplace?category_id=${subcat.id}`,
                hasChildren: allCategories.some(cat => cat.parent_id === subcat.id)
              }));

            return {
              id: category.id,
              title: getTranslatedName(category, i18n.language),
              icon: getCategoryIcon(category.id),
              path: `/marketplace?category_id=${category.id}`,
              subcategories: subcategories,
              hasMoreSubcategories: subcategories.length > MAX_VISIBLE_SUBCATEGORIES
            };
          });

          setCategoryGroups(groups);
        }
      } catch (error) {
        console.error('Error fetching categories:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchCategories();
  }, [i18n.language]);

  // Обработчики для меню категорий
  const handleOpenCategoryMenu = (event) => {
    setCategoryMenuAnchor(event.currentTarget);
  };

  const handleCloseCategoryMenu = () => {
    setCategoryMenuAnchor(null);
    // Сбрасываем состояние развернутых групп и подкатегорий при закрытии меню
    setExpandedGroups({});
    setExpandedDesktopSubcategories({});
  };

  const handleOpenMobileCategoryDrawer = () => {
    setMobileCategoryDrawerOpen(true);
    setMobileNavigationStack([]);
    setCurrentMobileCategory(null);
  };

  const handleCloseMobileCategoryDrawer = () => {
    setMobileCategoryDrawerOpen(false);
    setMobileNavigationStack([]);
    setCurrentMobileCategory(null);
    setExpandedMobileCategories({});
  };

  // Обработчик клика по категории
  const handleCategoryClick = (path) => {
    navigate(path);
    handleCloseCategoryMenu();
    handleCloseMobileCategoryDrawer();
  };

  // Обработчик для разворачивания/сворачивания группы подкатегорий
  const handleToggleGroup = (groupId, event) => {
    if (event) {
      event.stopPropagation();
      event.preventDefault();
    }
    setExpandedGroups(prev => ({
      ...prev,
      [groupId]: !prev[groupId]
    }));
  };

  // Обработчик для разворачивания/сворачивания мобильных категорий
  const handleToggleMobileCategory = (categoryId, event) => {
    if (event) {
      event.stopPropagation();
      event.preventDefault();
    }
    // Предотвращаем закрытие панели категорий
    setTimeout(() => {
      setExpandedMobileCategories(prev => ({
        ...prev,
        [categoryId]: !prev[categoryId]
      }));
    }, 0);
  };

  // Функция для навигации по категориям в мобильном режиме
  const handleMobileSubcategoryClick = (category, event) => {
    if (event) {
      event.stopPropagation();
      event.preventDefault();
    }

    // Предотвращаем закрытие панели категорий
    setTimeout(() => {
      // Проверяем, есть ли у категории подкатегории
      const hasSubcategories = categories.some(cat => cat.parent_id === category.id);

      // Если у категории есть дочерние элементы, открываем подкатегории
      if (hasSubcategories) {
        setMobileNavigationStack(prev => [...prev, currentMobileCategory]);
        setCurrentMobileCategory(category);
      } else {
        // Если нет дочерних элементов, переходим на страницу категории
        handleCategoryClick(category.path);
      }
    }, 0);
  };

  // Функция для навигации по подкатегориям в десктопной версии
  const [expandedDesktopSubcategories, setExpandedDesktopSubcategories] = useState({});
  const [desktopSubcategoriesCache, setDesktopSubcategoriesCache] = useState({});

  // Функция для получения подкатегорий для десктопной версии
  const getDesktopSubcategories = (parentId) => {
    if (desktopSubcategoriesCache[parentId]) {
      return desktopSubcategoriesCache[parentId];
    }

    const subcats = categories
      .filter(cat => cat.parent_id === parentId)
      .map(cat => ({
        id: cat.id,
        name: getTranslatedName(cat, i18n.language),
        path: `/marketplace?category_id=${cat.id}`,
        hasChildren: categories.some(c => c.parent_id === cat.id)
      }));

    setDesktopSubcategoriesCache(prev => ({
      ...prev,
      [parentId]: subcats
    }));

    return subcats;
  };

  const handleDesktopSubcategoryClick = (subcat, event) => {
    if (event) {
      event.stopPropagation();
      event.preventDefault();
    }

    // Переключаем состояние развернутости подкатегории
    setExpandedDesktopSubcategories(prev => ({
      ...prev,
      [subcat.id]: !prev[subcat.id]
    }));
  };

  // Функция для возврата на предыдущий уровень в мобильном режиме
  const handleMobileBack = () => {
    if (mobileNavigationStack.length > 0) {
      const previousCategory = mobileNavigationStack[mobileNavigationStack.length - 1];
      setCurrentMobileCategory(previousCategory);
      setMobileNavigationStack(prev => prev.slice(0, -1));
    } else {
      setCurrentMobileCategory(null);
    }
  };

  // Функция для получения подкатегорий текущей категории в мобильном режиме
  const getCurrentMobileSubcategories = () => {
    if (!currentMobileCategory) {
      // Для главного экрана используем тот же порядок, что и в десктопной версии
      return categoryGroups;
    }

    // Для подкатегорий получаем и сортируем по ID
    const subcategories = categories
      .filter(cat => cat.parent_id === currentMobileCategory.id)
      .map(cat => ({
        id: cat.id,
        title: getTranslatedName(cat, i18n.language),
        icon: getCategoryIcon(cat.id),
        path: `/marketplace?category_id=${cat.id}`,
        hasChildren: categories.some(c => c.parent_id === cat.id)
      }));

    // Сортируем подкатегории по ID
    return subcategories.sort((a, b) => a.id - b.id);
  };

  return (
    <>
      {/* Кнопка для открытия меню категорий */}
      <Button
        variant={isMobile ? "outlined" : "text"}
        sx={{
          color: '#004494',
          textTransform: 'none',
          fontWeight: 'normal',
          fontSize: isMobile ? '0.875rem' : 'inherit',
          minWidth: isMobile ? '40px' : 'auto',
          width: isMobile ? '40px' : 'auto',
          height: isMobile ? '40px' : 'auto',
          p: isMobile ? 0 : 2,
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          borderColor: isMobile ? 'rgba(0, 68, 148, 0.5)' : 'transparent',
          '&:hover': {
            borderColor: isMobile ? 'rgba(0, 68, 148, 0.8)' : 'transparent',
            backgroundColor: isMobile ? 'rgba(0, 68, 148, 0.04)' : 'rgba(0, 68, 148, 0.04)'
          }
        }}
        onClick={isMobile ? handleOpenMobileCategoryDrawer : handleOpenCategoryMenu}
        startIcon={!isMobile && <CategoryIcon fontSize="small" />}
        endIcon={!isMobile && <KeyboardArrowDown />}
        disabled={isLoading}
      >
        {isMobile ? <CategoryIcon fontSize="small" /> : t('navigation.allCategories', { defaultValue: 'ВСЕ КАТЕГОРИИ' })}
      </Button>

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
          {/* Разбиваем категории на ряды по 3 в каждом */}
          {Array.from({ length: Math.ceil(categoryGroups.length / 3) }).map((_, rowIndex) => (
            <React.Fragment key={`row-${rowIndex}`}>
              {categoryGroups.slice(rowIndex * 3, rowIndex * 3 + 3).map((group, colIndex) => (
                <Grid item xs={4} key={group.id}>
                  <Box
                    sx={{
                      mb: 2,
                      pb: 2,
                      borderBottom: rowIndex < Math.ceil(categoryGroups.length / 3) - 1 ? '1px solid' : 'none',
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
                  {/* Показываем ограниченное количество подкатегорий или все, если группа развернута */}
                  {(expandedGroups[group.id] ? group.subcategories : group.subcategories.slice(0, MAX_VISIBLE_SUBCATEGORIES)).map((subcat, subIdx) => (
                    <Box key={subcat.id} sx={{ width: '100%' }}>
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Button
                          onClick={() => handleCategoryClick(subcat.path)}
                          sx={{
                            display: 'flex',
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
                          <Typography variant="body2" noWrap sx={{ flex: 1 }}>
                            {subcat.name}
                          </Typography>
                        </Button>
                        {subcat.hasChildren && (
                          <IconButton
                            size="small"
                            onClick={(e) => {
                              e.stopPropagation();
                              e.preventDefault();
                              handleDesktopSubcategoryClick(subcat, e);
                            }}
                            sx={{ ml: 'auto', p: 0.5 }}
                          >
                            {expandedDesktopSubcategories[subcat.id] ?
                              <ExpandLess fontSize="small" /> :
                              <ChevronRight fontSize="small" />}
                          </IconButton>
                        )}
                      </Box>

                      {/* Подкатегории третьего уровня */}
                      {subcat.hasChildren && expandedDesktopSubcategories[subcat.id] && (
                        <Collapse in={expandedDesktopSubcategories[subcat.id]} timeout="auto" unmountOnExit>
                          <Box sx={{ pl: 2 }}>
                            {getDesktopSubcategories(subcat.id).map(childCat => (
                              <Button
                                key={childCat.id}
                                onClick={() => handleCategoryClick(childCat.path)}
                                sx={{
                                  display: 'block',
                                  width: '100%',
                                  textAlign: 'left',
                                  justifyContent: 'flex-start',
                                  py: 0.5,
                                  px: 1,
                                  color: 'text.secondary',
                                  fontSize: '0.8rem',
                                  textTransform: 'none',
                                  '&:hover': { backgroundColor: 'action.hover' }
                                }}
                              >
                                <Typography variant="body2" sx={{ fontSize: '0.8rem' }}>
                                  {childCat.name}
                                </Typography>
                              </Button>
                            ))}
                          </Box>
                        </Collapse>
                      )}
                    </Box>
                  ))}

                  {/* Кнопка "Показать еще" для групп с большим количеством подкатегорий */}
                  {group.hasMoreSubcategories && (
                    <Button
                      onClick={(e) => handleToggleGroup(group.id, e)}
                      sx={{
                        display: 'flex',
                        width: '100%',
                        justifyContent: 'flex-start',
                        py: 0.5,
                        px: 1,
                        color: 'primary.main',
                        textTransform: 'none',
                        fontSize: '0.75rem',
                        '&:hover': { backgroundColor: 'transparent' }
                      }}
                    >
                      {expandedGroups[group.id] ? (
                        <>
                          <ExpandLess fontSize="small" sx={{ mr: 0.5 }} />
                          {t('categories.showLess', { defaultValue: 'Свернуть' })}
                        </>
                      ) : (
                        <>
                          <ExpandMore fontSize="small" sx={{ mr: 0.5 }} />
                          {t('categories.showMore', { defaultValue: 'Показать еще' })}
                          ({group.subcategories.length - MAX_VISIBLE_SUBCATEGORIES})
                        </>
                      )}
                    </Button>
                  )}
                </Box>
              </Box>
                </Grid>
              ))}
            </React.Fragment>
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
            width: { xs: '90%', sm: '60%' },
            maxWidth: '350px',
            overflowX: 'hidden', // Предотвращаем горизонтальную прокрутку
          }
        }}
      >
        <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            {currentMobileCategory ? (
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <IconButton onClick={handleMobileBack} edge="start">
                  <ArrowBack />
                </IconButton>
                <Typography variant="h6" sx={{ ml: 1 }}>
                  {currentMobileCategory.title || currentMobileCategory.name}
                </Typography>
              </Box>
            ) : (
              <Typography variant="h6">{t('categories.title', { defaultValue: 'Категории' })}</Typography>
            )}
            <IconButton onClick={handleCloseMobileCategoryDrawer}>
              <CloseIcon />
            </IconButton>
          </Box>
        </Box>

        <Box sx={{ overflow: 'auto', flex: 1, width: '100%', maxWidth: '100%', boxSizing: 'border-box' }}>
          {getCurrentMobileSubcategories().map((category) => (
            <Box key={category.id}>
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  px: 2,
                  py: 1.5,
                  width: '100%',
                  maxWidth: '100%',
                  boxSizing: 'border-box',
                  color: 'text.primary',
                  fontWeight: category.hasChildren ? 'bold' : 'normal',
                  textDecoration: 'none',
                  borderBottom: '1px solid',
                  borderColor: 'divider',
                  justifyContent: 'space-between',
                  pr: 0, // Убираем отступ справа
                  flexWrap: 'nowrap', // Запрещаем перенос элементов на новую строку
                }}
              >
                <Button
                  onClick={(e) => {
                    e.stopPropagation();
                    e.preventDefault();

                    // Добавляем отладочную информацию
                    console.log('Category clicked:', category.id, category.title || category.name);
                    console.log('hasChildren:', category.hasChildren);
                    console.log('currentMobileCategory:', currentMobileCategory);

                    // Проверяем, есть ли у категории подкатегории
                    const hasSubcategories = categories.some(cat => cat.parent_id === category.id);
                    console.log('Actual subcategories check:', hasSubcategories);

                    // Если мы находимся в подкатегориях и у категории есть дочерние элементы
                    if (currentMobileCategory && hasSubcategories) {
                      console.log('Toggle subcategories');
                      // Разворачиваем/сворачиваем подкатегории
                      handleToggleMobileCategory(category.id, e);
                    }
                    // Если мы на главном экране и у категории есть дочерние элементы
                    else if (!currentMobileCategory && hasSubcategories) {
                      console.log('Navigate to subcategory');
                      // Переходим в подкатегорию
                      handleMobileSubcategoryClick(category, e);
                    }
                    // Если нет дочерних элементов, переходим на страницу категории
                    else {
                      console.log('Navigate to category page');
                      handleCategoryClick(category.path);
                    }
                  }}
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    textAlign: 'left',
                    justifyContent: 'flex-start',
                    textTransform: 'none',
                    color: 'text.primary',
                    px: 0,
                    flex: 1,
                    minWidth: 0, // Позволяем кнопке сжиматься меньше минимальной ширины
                    maxWidth: 'calc(100% - 40px)', // Оставляем место для кнопки разворачивания
                    overflow: 'hidden', // Скрываем выходящий за пределы текст
                    '&:hover': {
                      bgcolor: 'transparent'
                    }
                  }}
                >
                  {category.icon && (
                    <Box sx={{ mr: 2, color: 'primary.main' }}>
                      {category.icon}
                    </Box>
                  )}
                  <Typography variant={category.hasChildren ? "subtitle1" : "body2"}>
                    {category.title || category.name}
                  </Typography>
                </Button>

                {/* Добавляем более заметную кнопку для разворачивания только для категорий с подкатегориями */}
                {categories.some(cat => cat.parent_id === category.id) && <Button
                  variant="outlined"
                  size="small"
                  onClick={(e) => {
                    e.stopPropagation();
                    e.preventDefault();

                    // Добавляем отладочную информацию
                    console.log('Expand button clicked:', category.id, category.title || category.name);

                    // Проверяем, есть ли у категории подкатегории
                    const hasSubcategories = categories.some(cat => cat.parent_id === category.id);
                    console.log('Actual subcategories check (button):', hasSubcategories);

                    // Если мы находимся в подкатегориях и у категории есть дочерние элементы
                    if (currentMobileCategory && hasSubcategories) {
                      console.log('Toggle subcategories (button)');
                      // Разворачиваем/сворачиваем подкатегории
                      handleToggleMobileCategory(category.id, e);
                    }
                    // Если мы на главном экране и у категории есть дочерние элементы
                    else if (!currentMobileCategory && hasSubcategories) {
                      console.log('Navigate to subcategory (button)');
                      // Переходим в подкатегорию
                      handleMobileSubcategoryClick(category, e);
                    }
                    // Если нет дочерних элементов, переходим на страницу категории
                    else {
                      console.log('Navigate to category page (button)');
                      handleCategoryClick(category.path);
                    }
                  }}
                  sx={{
                    minWidth: '32px',
                    width: '32px',
                    maxWidth: '32px',
                    ml: 0,
                    p: 0.5,
                    borderColor: 'divider',
                    flexShrink: 0, // Запрещаем сжиматься кнопке
                    '&:hover': { borderColor: 'primary.main' }
                  }}
                >
                  {currentMobileCategory ?
                    (expandedMobileCategories[category.id] ?
                      <ExpandLess fontSize="small" />
                    :
                      <ExpandMore fontSize="small" />
                    ) :
                    <ChevronRight fontSize="small" />
                  }
                </Button>}
              </Box>

              {/* Развернутые подкатегории в мобильном режиме */}
              {currentMobileCategory && categories.some(cat => cat.parent_id === category.id) && (
                <Collapse in={expandedMobileCategories[category.id]} timeout="auto" unmountOnExit>
                  <List disablePadding>
                    {categories
                      .filter(cat => cat.parent_id === category.id)
                      .sort((a, b) => a.id - b.id) // Сортируем по ID
                      .map(subcat => (
                        <ListItem
                          key={subcat.id}
                          button
                          onClick={() => handleCategoryClick(`/marketplace?category_id=${subcat.id}`)}
                          sx={{ pl: 4 }}
                        >
                          <ListItemText primary={getTranslatedName(subcat, i18n.language)} />
                          {categories.some(c => c.parent_id === subcat.id) && (
                            <IconButton
                              size="small"
                              onClick={(e) => {
                                e.stopPropagation();
                                e.preventDefault();
                                handleMobileSubcategoryClick({
                                  id: subcat.id,
                                  title: getTranslatedName(subcat, i18n.language),
                                  hasChildren: true,
                                  path: `/marketplace?category_id=${subcat.id}`
                                }, e);
                              }}
                            >
                              <ChevronRight fontSize="small" />
                            </IconButton>
                          )}
                        </ListItem>
                      ))
                    }
                  </List>
                </Collapse>
              )}
            </Box>
          ))}
        </Box>
      </Drawer>
    </>
  );
};

export default CategoryMenu;

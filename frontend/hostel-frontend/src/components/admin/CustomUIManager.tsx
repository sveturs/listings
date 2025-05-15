import React, { useState, useEffect } from 'react';
import {
  Typography,
  Paper,
  Box,
  FormControlLabel,
  Checkbox,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  Divider,
  Alert,
  Tooltip,
  CircularProgress,
  alpha,
  useTheme
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';

interface CustomUIManagerProps {
  categoryId: number;
  onCategoryUpdate: () => void;
}

// Доступные UI-компоненты для категорий
const availableComponents = [
  { value: "AutoCategoryUI", label: "Автомобили" },
  { value: "RealEstateCategoryUI", label: "Недвижимость" },
  { value: "ElectronicsCategoryUI", label: "Электроника" },
  { value: "FurnitureCategoryUI", label: "Мебель" },
  { value: "ClothingCategoryUI", label: "Одежда" },
  { value: "ServicesCategoryUI", label: "Услуги" },
];

const CustomUIManager: React.FC<CustomUIManagerProps> = ({ categoryId, onCategoryUpdate }) => {
  const { t } = useTranslation();
  const theme = useTheme();
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasCustomUI, setHasCustomUI] = useState<boolean>(false);
  const [customUIComponent, setCustomUIComponent] = useState<string>('');
  const [changed, setChanged] = useState<boolean>(false);

  // Загрузка информации о текущем состоянии пользовательского UI для категории
  useEffect(() => {
    const fetchCategoryData = async () => {
      if (!categoryId) return;

      try {
        setLoading(true);
        setError(null);
        
        // Для мок-данных
        const mockCategories = [
          { id: 1, name: 'Электроника', slug: 'electronics', has_custom_ui: true, custom_ui_component: 'ElectronicsCategoryUI' },
          { id: 2, name: 'Смартфоны', slug: 'smartphones', has_custom_ui: false },
          { id: 3, name: 'Ноутбуки', slug: 'laptops', has_custom_ui: true, custom_ui_component: 'ElectronicsCategoryUI' },
          { id: 4, name: 'Одежда', slug: 'clothing', has_custom_ui: false },
        ];
        
        // Находим категорию по ID
        const categoryData = mockCategories.find(cat => cat.id === categoryId) || { has_custom_ui: false, custom_ui_component: '' };
        
        // Устанавливаем состояние из мок-данных
        setHasCustomUI(categoryData.has_custom_ui || false);
        setCustomUIComponent(categoryData.custom_ui_component || '');
        setChanged(false);
        
        // Симулируем задержку загрузки
        await new Promise(resolve => setTimeout(resolve, 500));
        
      } catch (err) {
        console.error('Error fetching category UI data:', err);
        setError(t('admin.customUI.fetchError'));
      } finally {
        setLoading(false);
      }
    };

    fetchCategoryData();
  }, [categoryId, t]);

  // Обработчик изменения флага кастомного UI
  const handleCustomUIChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setHasCustomUI(e.target.checked);
    if (!e.target.checked) {
      setCustomUIComponent('');
    }
    setChanged(true);
  };

  // Обработчик изменения компонента UI
  const handleComponentChange = (e: any) => {
    setCustomUIComponent(e.target.value);
    setChanged(true);
  };

  // Сохранение изменений
  const handleSave = async () => {
    if (!categoryId) return;

    try {
      setLoading(true);
      setError(null);
      
      // Симулируем запрос к API
      console.log('Сохранение настроек UI для категории:', {
        categoryId,
        has_custom_ui: hasCustomUI,
        custom_ui_component: customUIComponent
      });
      
      // Симулируем задержку запроса
      await new Promise(resolve => setTimeout(resolve, 700));
      
      setChanged(false);
      
      // Вызываем callback для обновления данных родительского компонента
      onCategoryUpdate();
      
    } catch (err) {
      console.error('Error updating category UI settings:', err);
      setError(t('admin.customUI.updateError'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper 
      variant="outlined" 
      sx={{ 
        p: 2, 
        mb: 3,
        borderRadius: 1,
        borderColor: alpha(theme.palette.primary.main, 0.3),
        borderWidth: '1px',
        position: 'relative',
        overflow: 'hidden',
        '&:before': {
          content: '""',
          position: 'absolute',
          top: 0,
          left: 0,
          width: '4px',
          height: '100%',
          bgcolor: 'primary.main'
        }
      }}
    >
      <Typography 
        variant="h6" 
        gutterBottom 
        sx={{ 
          display: 'flex', 
          alignItems: 'center', 
          gap: 0.5,
          color: 'primary.main',
          pl: 1
        }}
      >
        <InfoOutlinedIcon fontSize="small" />
        {t('admin.customUI.title')}
        <Tooltip title={t('admin.customUI.tooltip')}>
          <span>
            <HelpOutlineIcon fontSize="small" sx={{ ml: 1, color: 'text.secondary', cursor: 'help' }} />
          </span>
        </Tooltip>
      </Typography>
      
      <Divider sx={{ mb: 2 }} />
      
      {error && (
        <Alert 
          severity="error" 
          sx={{ mb: 2 }} 
          onClose={() => setError(null)}
        >
          {error}
        </Alert>
      )}
      
      <Box sx={{ mt: 2, position: 'relative' }}>
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress size={30} />
          </Box>
        ) : (
          <>
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                {t('admin.customUI.description')}
              </Typography>
              
              <FormControlLabel
                control={
                  <Checkbox
                    checked={hasCustomUI}
                    onChange={handleCustomUIChange}
                    name="has_custom_ui"
                    color="primary"
                  />
                }
                label={t('admin.customUI.enableCustomUI')}
              />
            </Box>
            
            {hasCustomUI && (
              <Box sx={{ mt: 2 }}>
                <FormControl fullWidth margin="normal">
                  <InputLabel id="custom-component-label">
                    {t('admin.customUI.selectComponent')}
                  </InputLabel>
                  <Select
                    labelId="custom-component-label"
                    value={customUIComponent}
                    onChange={handleComponentChange}
                    label={t('admin.customUI.selectComponent')}
                    required={hasCustomUI}
                  >
                    <MenuItem value="">
                      <em>{t('admin.customUI.noComponent')}</em>
                    </MenuItem>
                    {availableComponents.map(component => (
                      <MenuItem key={component.value} value={component.value}>
                        {component.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                  {t('admin.customUI.componentDescription')}
                </Typography>
              </Box>
            )}
            
            <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
              <Button 
                variant="contained" 
                color="primary" 
                onClick={handleSave}
                disabled={!changed}
              >
                {t('common.save')}
              </Button>
            </Box>
          </>
        )}
      </Box>
    </Paper>
  );
};

export default CustomUIManager;
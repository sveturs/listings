// frontend/hostel-frontend/src/components/marketplace/MarketplaceFilters.tsx
import React, { useMemo, useEffect, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import AutocompleteInput from '../shared/AutocompleteInput';
import { Search, X, Map } from 'lucide-react';
import {
    Paper,
    Box,
    TextField,
    InputAdornment,
    IconButton,
    Typography,
    Stack,
    Divider,
    Button,
    FormControl,
    FormControlLabel,
    Select,
    MenuItem,
    InputLabel,
    RadioGroup,
    Radio,
    ToggleButton,
    ToggleButtonGroup,
    SelectChangeEvent
} from '@mui/material';
import {
    List,
    Grid as GridIcon
} from 'lucide-react';
import { useLocation } from '../../contexts/LocationContext';
import VirtualizedCategoryTree from './VirtualizedCategoryTree';
import AttributeFilters from './AttributeFilters';

// TypeScript interfaces
interface FilterOptions {
  query?: string;
  category_id?: number | string;
  min_price?: string | number;
  max_price?: string | number;
  distance?: string; // Values like "1km", "3km", etc.
  condition?: '' | 'new' | 'used';
  latitude?: number;
  longitude?: number;
  [key: string]: any; // For dynamic attribute filters
}

interface AttributeFiltersType {
  [key: string]: any; // Dynamic nature of attribute filters
}

interface MarketplaceFiltersProps {
  filters: FilterOptions;
  onFilterChange: (newFilters: FilterOptions | ((prevFilters: FilterOptions) => FilterOptions)) => void;
  selectedCategoryId: number | string | null;
  onToggleMapView: () => void;
  setSearchParams: (params: Record<string, string>) => void;
  fetchListings: () => void;
  viewMode: 'grid' | 'list';
  handleViewModeChange: (event: React.MouseEvent<HTMLElement>, newViewMode: 'grid' | 'list' | null) => void;
}

const MarketplaceFilters: React.FC<MarketplaceFiltersProps> = ({
    filters,
    onFilterChange,
    selectedCategoryId,
    onToggleMapView,
    setSearchParams,
    fetchListings,
    viewMode,
    handleViewModeChange
}) => {
    const { t } = useTranslation(['marketplace', 'common']);
    const { userLocation, detectUserLocation } = useLocation();
    
    // Обработчик для выбора категории
    const handleCategorySelect = useCallback((id: number | string): void => {
        console.log(`MarketplaceFilters: Выбрана категория с ID: ${id}`);

        // Вызываем onFilterChange, но с пустым запросом, чтобы показать все товары категории
        onFilterChange({
            ...filters,
            category_id: id,
            query: '' // Очищаем поисковый запрос при выборе категории
        });
    }, [filters, onFilterChange]);

    // Обновленная функция handleFilterChange
    const handleFilterChange = useCallback((attributeFilters: AttributeFiltersType): void => {
        console.log(`MarketplaceFilters: Выбраны фильтры:`, attributeFilters);

        // ВАЖНО: Теперь мы полностью игнорируем атрибуты недвижимости в боковой панели
        // Они будут обрабатываться только через CentralAttributeFilters
        const nonRealEstateAttrs = { ...attributeFilters };

        // Вызываем родительский обработчик только с базовыми фильтрами
        onFilterChange(prev => {
            const updated = { ...prev };
            Object.assign(updated, nonRealEstateAttrs);
            return updated;
        });
    }, [onFilterChange]);

    const handleDistanceChange = async (value: string): Promise<void> => {
        // Если выбрано расстояние, но нет координат, используем геолокацию
        if (value && (!filters.latitude || !filters.longitude)) {
            try {
                // Используем функцию из контекста местоположения
                await detectUserLocation();

                // userLocation должен обновиться автоматически через контекст
                // Затем будет вызвано событие cityChanged, которое обновит фильтры

                // Здесь просто обновляем distance, так как координаты обновятся через событие
                onFilterChange({ ...filters, distance: value });
            } catch (error) {
                console.error("Ошибка получения геолокации:", error);
                // Показываем уведомление пользователю
                alert(t('listings.filters.distance.locationError', {
                    defaultValue: 'Для использования фильтра по расстоянию необходимо разрешить доступ к вашему местоположению'
                }));
            }
        } else {
            // Если координаты уже есть, просто обновляем фильтр расстояния
            onFilterChange({ ...filters, distance: value });
        }
    };

    const isMapAvailable = useMemo((): boolean => {
        // Карта доступна, если либо нет фильтра по расстоянию, либо есть координаты
        return !filters.distance || (userLocation?.lat !== undefined && userLocation?.lon !== undefined);
    }, [filters.distance, userLocation]);

    const isDistanceWithoutCoordinates = filters.distance && (!userLocation?.lat || !userLocation?.lon);

    return (
        <Paper variant="elevation" elevation={3} sx={{ height: '100%', display: 'flex', flexDirection: 'column', width: '100%', minWidth: '250px', maxWidth: '100%', boxSizing: 'border-box', overflow: 'hidden' }}>
            {/* Поиск с автодополнением */}
            <Box sx={{ p: 2, backgroundColor: 'background.default', boxShadow: '0px 1px 2px rgba(0, 0, 0, 0.1)', zIndex: 1 }}>
                <AutocompleteInput
                    value={filters.query || ''}
                    onChange={(value) => onFilterChange({ ...filters, query: value })}
                    onSearch={(value, categoryId) => {
                        // Если предоставлен categoryId, обновляем и категорию
                        if (categoryId) {
                            onFilterChange({ query: value, category_id: categoryId });
                        } else {
                            onFilterChange({ query: value });
                        }
                    }}
                    placeholder={t('buttons.search', { ns: 'common' })}
                />
            </Box>

            {/* Панель переключения режимов просмотра */}
            <Box sx={{ px: 2, py: 1, display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 1 }}>
                <Button
                    variant="outlined"
                    startIcon={<Map />}
                    onClick={onToggleMapView}
                    disabled={!isMapAvailable && Boolean(filters.distance)}
                    sx={{ flex: 1, mr: 1 }}
                    size="medium"
                >
                    {t('listings.map.map')}
                </Button>

                <ToggleButtonGroup
                    value={viewMode}
                    exclusive
                    onChange={handleViewModeChange}
                    aria-label="view mode"
                    size="small"
                >
                    <ToggleButton value="grid" aria-label="grid view">
                        <GridIcon size={18} />
                    </ToggleButton>
                    <ToggleButton value="list" aria-label="list view">
                        <List size={18} />
                    </ToggleButton>
                </ToggleButtonGroup>
            </Box>

            {/* Предупреждение о необходимости выбрать местоположение */}
            {!isMapAvailable && filters.distance && (
                <Box sx={{ px: 2, py: 1, color: 'warning.main' }}>
                    <Typography variant="caption">
                        {t('listings.map.needLocation')}
                    </Typography>
                </Box>
            )}

            <Divider sx={{ my: 1 }} />

            {/* Основные фильтры */}
            <Box sx={{ p: 2, overflowY: 'auto' }}>
                <Typography variant="subtitle1" gutterBottom>{t('listings.filters.title')}</Typography>
                <Stack spacing={2}>
                    <Box>
                        <Typography gutterBottom>{t('listings.filters.price.label')}</Typography>
                        <Stack direction="row" spacing={1}>
                            <TextField
                                size="small"
                                type="number"
                                placeholder={t('listings.filters.price.min')}
                                value={filters.min_price || ''}
                                onChange={(e) => onFilterChange({ ...filters, min_price: e.target.value })}
                            />
                            <TextField
                                size="small"
                                type="number"
                                placeholder={t('listings.filters.price.max')}
                                value={filters.max_price || ''}
                                onChange={(e) => onFilterChange({ ...filters, max_price: e.target.value })}
                            />
                        </Stack>
                    </Box>

                    {/* Фильтр по расстоянию (для карты) */}
                    <Box>
                        <Typography gutterBottom>{t('listings.filters.distance.label')}</Typography>
                        <FormControl fullWidth size="small">
                            <Select
                                value={filters.distance || ''}
                                onChange={(e: SelectChangeEvent<string>) => handleDistanceChange(e.target.value)}
                                displayEmpty
                            >
                                <MenuItem value="">{t('listings.filters.distance.any')}</MenuItem>
                                <MenuItem value="1km">1 км</MenuItem>
                                <MenuItem value="3km">3 км</MenuItem>
                                <MenuItem value="5km">5 км</MenuItem>
                                <MenuItem value="10km">10 км</MenuItem>
                                <MenuItem value="15km">15 км</MenuItem>
                                <MenuItem value="30km">30 км</MenuItem>
                            </Select>
                        </FormControl>
                        <Typography variant="caption" color={isDistanceWithoutCoordinates ? "error.main" : "text.secondary"}>
                            {userLocation?.city
                                ? t('listings.filters.distance.fromCity', { city: userLocation.city })
                                : t('listings.filters.distance.needLocation')}
                        </Typography>
                        {isDistanceWithoutCoordinates && (
                            <Typography variant="caption" color="error.main" sx={{ display: 'block', mt: 0.5 }}>
                                {t('listings.filters.distance.warningNoLocation', { defaultValue: 'Укажите местоположение, чтобы использовать фильтр по расстоянию' })}
                            </Typography>
                        )}
                    </Box>

                    {/* Фильтр по состоянию */}
                    <Box>
                        <Typography gutterBottom>{t('listings.filters.condition.label')}</Typography>
                        <RadioGroup
                            value={filters.condition || ''}
                            onChange={(e) => onFilterChange({ ...filters, condition: e.target.value })}
                        >
                            <FormControlLabel
                                value=""
                                control={<Radio size="small" />}
                                label={t('listings.filters.condition.any')}
                            />
                            <FormControlLabel
                                value="new"
                                control={<Radio size="small" />}
                                label={t('listings.conditions.new')}
                            />
                            <FormControlLabel
                                value="used"
                                control={<Radio size="small" />}
                                label={t('listings.conditions.used')}
                            />
                        </RadioGroup>
                    </Box>

                    {/* Добавляем компонент AttributeFilters если выбрана категория */}

                </Stack>
            </Box>

            {/* Категории */}
            <Box sx={{
                flex: 1,
                overflow: 'auto',
                p: 2,
                backgroundColor: 'background.paper',
                borderTop: 1,
                borderColor: 'divider'
            }}>
                <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                    {t('listings.create.categories')}
                </Typography>
                <VirtualizedCategoryTree
                    selectedId={selectedCategoryId}
                    onSelectCategory={handleCategorySelect}
                />
            </Box>
        </Paper>
    );
};

export default React.memo(MarketplaceFilters);
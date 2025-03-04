// frontend/hostel-frontend/src/components/marketplace/MarketplaceFilters.js
import React, { useMemo, useEffect, useCallback } from 'react';
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
    Radio
} from '@mui/material';

import VirtualizedCategoryTree from './VirtualizedCategoryTree';

const CompactMarketplaceFilters = ({ filters, onFilterChange, selectedCategoryId, onToggleMapView }) => {
    const { t } = useTranslation('marketplace', 'common');

    const handleCategorySelect = useCallback((id) => {
        console.log(`MarketplaceFilters: Выбрана категория с ID: ${id}`);

        // Вызываем onFilterChange, но с пустым запросом, чтобы показать все товары категории
        onFilterChange({
            ...filters,
            category_id: id,
            query: '' // Очищаем поисковый запрос при выборе категории
        });
    }, [filters, onFilterChange]);


    // Проверяем правильность передачи поискового запроса
    const handleSearchChange = useCallback((value) => {
        console.log("Поисковый запрос изменен на:", value);
        onFilterChange({ ...filters, query: value });
    }, [filters, onFilterChange]);

    // Определяем, доступна ли карта (не используется в запросе distance без координат)
    const isMapAvailable = useMemo(() => {
        return !filters.distance || (filters.latitude && filters.longitude);
    }, [filters.distance, filters.latitude, filters.longitude]);
    const isDistanceWithoutCoordinates = filters.distance && (!filters.latitude || !filters.longitude);
    return (
        <Paper variant="elevation" elevation={3} sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* Поиск с автодополнением */}
            <Box sx={{ p: 2, backgroundColor: 'background.default', boxShadow: '0px 1px 2px rgba(0, 0, 0, 0.1)', zIndex: 1 }}>
                <AutocompleteInput
                    value={filters.query || ''}
                    onChange={handleSearchChange} // Используем функцию, определенную в этом компоненте
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

            {/* Кнопка просмотра на карте */}
            <Box sx={{ px: 2, py: 1, display: 'flex', justifyContent: 'center' }}>
                <Button
                    variant="outlined"
                    startIcon={<Map />}
                    fullWidth
                    onClick={onToggleMapView}
                    disabled={!isMapAvailable && filters.distance}
                >
                    {t('listings.map.showOnMap')}
                </Button>
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
                                onChange={(e) => onFilterChange({ ...filters, distance: e.target.value })}
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
                            {filters.latitude && filters.longitude
                                ? t('listings.filters.distance.fromYourLocation')
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
                    {t('listings.create.сategories')}
                </Typography>
                <VirtualizedCategoryTree
                    selectedId={selectedCategoryId}
                    onSelectCategory={handleCategorySelect}
                />
            </Box>
        </Paper>
    );
};

export default React.memo(CompactMarketplaceFilters);
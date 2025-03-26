// frontend/hostel-frontend/src/components/marketplace/CentralAttributeFilters.js
import React, { useState, useEffect, useCallback, useRef } from 'react';
import { 
    Paper, 
    Typography, 
    Box, 
    Collapse, 
    Button, 
    Divider,
    useTheme,
    useMediaQuery,
    IconButton,
    Alert
} from '@mui/material';
import { ChevronDown, ChevronUp, Filter } from 'lucide-react';
import AttributeFilters from './AttributeFilters';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';

const CentralAttributeFilters = ({ 
    categoryId, 
    onFilterChange, 
    filters = {},
    resetAttributeFilters 
}) => {
    const { t } = useTranslation('marketplace');
    const [expanded, setExpanded] = useState(true);
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const [hasAttributes, setHasAttributes] = useState(false);
    const [debugData, setDebugData] = useState(null);
    const [debugError, setDebugError] = useState(null);
    const [debugLoading, setDebugLoading] = useState(false);

    // Преобразуем categoryId в строку для безопасного сравнения
    const catId = String(categoryId);
    
    console.log(`DEBUG-CAF: Инициализация компонента, categoryId=${categoryId}, тип=${typeof categoryId}`);

    const hasActiveFilters = filters && 
                           typeof filters === 'object' && 
                           Object.keys(filters).length > 0;
    
    // Прямой отладочный запрос к API
    useEffect(() => {
        if (!categoryId) return;
        
        console.log(`DEBUG-CAF: Прямой запрос к API для категории ${categoryId}`);
        setDebugLoading(true);
        setDebugError(null);
        
        const fetchDebugData = async () => {
            try {
                const response = await axios.get(`/api/v1/marketplace/categories/${categoryId}/attributes`);
                console.log(`DEBUG-CAF: Ответ API для категории ${categoryId}:`, response);
                
                if (response.data?.data) {
                    const filterableAttrs = response.data.data.filter(attr => attr.is_filterable);
                    console.log(`DEBUG-CAF: Получено ${response.data.data.length} атрибутов, фильтруемых: ${filterableAttrs.length}`);
                    setDebugData({
                        total: response.data.data.length,
                        filterable: filterableAttrs.length,
                        attributes: filterableAttrs.map(a => a.name)
                    });
                    setHasAttributes(filterableAttrs.length > 0);
                } else {
                    console.log(`DEBUG-CAF: Нет данных в ответе API`);
                    setDebugData({ total: 0, filterable: 0, attributes: [] });
                    setHasAttributes(false);
                }
            } catch (error) {
                console.error(`DEBUG-CAF: Ошибка при запросе к API:`, error);
                setDebugError(error.message || 'Ошибка при загрузке атрибутов');
                setHasAttributes(false);
            } finally {
                setDebugLoading(false);
            }
        };
        
        fetchDebugData();
    }, [categoryId]);
    
    useEffect(() => {
        if (hasActiveFilters) {
            setExpanded(true);
        }
    }, [hasActiveFilters]);

    const toggleExpanded = useCallback(() => {
        setExpanded(prev => !prev);
    }, []);

    const handleAttributesLoaded = useCallback((hasAny) => {
        console.log(`DEBUG-CAF: handleAttributesLoaded(${hasAny}) вызван для категории ${categoryId}`);
        setHasAttributes(hasAny);
    }, [categoryId]);

    // Функция getCategoryTitle всегда возвращает правильный заголовок по categoryId
    const getCategoryTitle = () => {
        // В соответствии с данными API и базы данных:
        switch (catId) {
            case '1100': return t('listings.filters.property_attributes', { defaultValue: 'Параметры недвижимости' });
            case '2000': return t('listings.filters.auto_attributes', { defaultValue: 'Параметры автомобиля' });
            case '3110': return t('listings.filters.phone_attributes', { defaultValue: 'Параметры телефона' });
            case '3810': return t('listings.filters.tablet_attributes', { defaultValue: 'Параметры планшета' });
            case '3310': return t('listings.filters.desktop_attributes', { defaultValue: 'Параметры компьютера' });
            case '3320': return t('listings.filters.laptop_attributes', { defaultValue: 'Параметры ноутбука' });
            case '3600': return t('listings.filters.tech_attributes', { defaultValue: 'Параметры техники' });
            default: return t('listings.filters.specific_attributes', { defaultValue: 'Параметры' });
        }
    };

    return (
        <Paper 
            sx={{ 
                p: 2, 
                mb: 3, 
                width: '100%', 
                position: 'relative',
                overflow: 'hidden'
            }}
        >
            <Box 
                sx={{ 
                    display: 'flex', 
                    justifyContent: 'space-between', 
                    alignItems: 'center',
                    cursor: 'pointer',
                    pb: 1,
                    mb: expanded ? 1 : 0
                }}
                onClick={toggleExpanded}
            >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Filter size={18} />
                    <Typography variant="h6" sx={{ fontSize: '1.1rem' }}>
                        {getCategoryTitle()} [{categoryId}]
                    </Typography>
                    
                    {hasActiveFilters && (
                        <Box 
                            component="span" 
                            sx={{ 
                                ml: 1, 
                                bgcolor: 'primary.main', 
                                color: 'white', 
                                width: 20, 
                                height: 20, 
                                borderRadius: '50%', 
                                display: 'flex', 
                                alignItems: 'center', 
                                justifyContent: 'center', 
                                fontSize: '0.75rem' 
                            }}
                        >
                            {Object.keys(filters).length}
                        </Box>
                    )}
                </Box>
                
                <IconButton 
                    size="small"
                    onClick={(e) => {
                        e.stopPropagation();
                        toggleExpanded();
                    }}
                >
                    {expanded ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
                </IconButton>
            </Box>
            
            {/* Добавляем кнопку сброса фильтров, которая всегда видна если есть активные фильтры */}
            {hasActiveFilters && expanded && (
                <Box sx={{ mb: 2 }}>
                    <Button 
                        size="small" 
                        color="error" 
                        variant="outlined"
                        onClick={(e) => {
                            e.stopPropagation();
                            if (resetAttributeFilters) {
                                resetAttributeFilters();
                            }
                        }}
                    >
                        {t('listings.filters.resetAttributes', { defaultValue: 'Сбросить параметры' })}
                    </Button>
                </Box>
            )}

            <Collapse in={expanded}>
                <Divider sx={{ mb: 2 }} />
                
                <AttributeFilters
                    categoryId={categoryId}
                    onFilterChange={onFilterChange}
                    filters={filters}
                    onAttributesLoaded={handleAttributesLoaded}
                />
            </Collapse>
            
            {/* Сообщение если нет атрибутов показываем только после загрузки */}
            {expanded && !hasAttributes && !debugLoading && (
                <Box sx={{ py: 2, textAlign: 'center' }}>
                    <Typography color="text.secondary">
                        {t('listings.filters.noAttributes', { defaultValue: 'Нет доступных параметров для фильтрации' })}
                    </Typography>
                </Box>
            )}
        </Paper>
    );
};

export default CentralAttributeFilters;
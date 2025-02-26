import React from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import {
    Breadcrumbs as MuiBreadcrumbs,
    Typography,
    Box,
    useTheme
} from '@mui/material';
import { ChevronRight } from 'lucide-react';

const Breadcrumbs = ({ paths, categories }) => {
    const { t, i18n } = useTranslation(['common', 'marketplace']);
    const theme = useTheme();
    const navigate = useNavigate();

    const getTranslatedName = (pathCategory) => {
        if (!pathCategory) return '';
        
        // Если у категории есть переводы, пробуем их использовать
        if (pathCategory.translations && Object.keys(pathCategory.translations).length > 0) {
            if (i18n.language && pathCategory.translations[i18n.language]) {
                return pathCategory.translations[i18n.language];
            }
            
            // Приоритеты языков
            const langPriority = [i18n.language, 'sr', 'ru', 'en'];
            for (const lang of langPriority) {
                if (pathCategory.translations[lang]) {
                    return pathCategory.translations[lang];
                }
            }
        }
        
        // Если переводов нет или они не подходят, ищем категорию в общем списке
        if (categories && Array.isArray(categories)) {
            const categoryFromList = categories.find(c => String(c.id) === String(pathCategory.id));
            if (categoryFromList && categoryFromList.translations) {
                const langPriority = [i18n.language, 'sr', 'ru', 'en'];
                for (const lang of langPriority) {
                    if (categoryFromList.translations[lang]) {
                        return categoryFromList.translations[lang];
                    }
                }
            }
        }
        
        // Если ничего не найдено, возвращаем исходное имя
        return pathCategory.name;
    };
    if (!paths || paths.length === 0) {
        return null;
    }
// Добавьте эту функцию в компонент Breadcrumbs, если её нет
const truncateLongPath = (paths) => {
    // Если путь слишком длинный, показываем только первые 2 и последнюю категории
    if (paths.length > 3) {
        const truncatedPath = [
            paths[0], // Первая категория
            { 
                id: 'truncated', 
                name: '...', 
                slug: '', 
                translations: {} 
            }, // Отображение многоточия
            ...paths.slice(paths.length - 1) // Последняя категория
        ];
        return truncatedPath;
    }
    return paths;
};
    const handleCategoryClick = (categoryId, event) => {
        event.preventDefault();
        navigate(`/marketplace?category_id=${categoryId}`);
    };

    return (
        <Box sx={{
            mb: 2,
            py: 0,
            overflowX: 'auto',
            whiteSpace: 'nowrap',
            '&::-webkit-scrollbar': {
                display: 'none'
            },
            scrollbarWidth: 'none'
        }}>
            <MuiBreadcrumbs
                separator={<ChevronRight size={16} />}
                aria-label="breadcrumb"
                sx={{
                    '& .MuiBreadcrumbs-ol': {
                        flexWrap: 'nowrap'
                    },
                    '& .MuiBreadcrumbs-li': {
                        display: 'flex',
                        alignItems: 'center'
                    }
                }}
            >
                <Link
                    to="/marketplace"
                    style={{
                        color: theme.palette.text.secondary,
                        textDecoration: 'none',
                        padding: '4px 8px',
                        borderRadius: '4px',
                        transition: 'all 0.2s'
                    }}
                >
                    {t('navigation.home', { ns: 'common' })}
                </Link>

                {paths.map((path, index) => {
                    const isLast = index === paths.length - 1;
                    const translatedName = getTranslatedName(path);

                    if (isLast) {
                        return (
                            <Typography
                                key={path.id}
                                color="text.primary"
                                sx={{
                                    fontSize: '0.875rem',
                                    padding: '4px 8px',
                                    fontWeight: 500
                                }}
                            >
                                {translatedName}
                            </Typography>
                        );
                    }

                    return (
                        <Link
                            key={path.id}
                            to={`/marketplace?category_id=${path.id}`}
                            onClick={(e) => handleCategoryClick(path.id, e)}
                            style={{
                                color: theme.palette.text.secondary,
                                textDecoration: 'none',
                                padding: '4px 8px',
                                borderRadius: '4px',
                                fontSize: '0.875rem',
                                transition: 'all 0.2s'
                            }}
                        >
                            {translatedName}
                        </Link>
                    );
                })}
            </MuiBreadcrumbs>
        </Box>
    );
};

export default Breadcrumbs;
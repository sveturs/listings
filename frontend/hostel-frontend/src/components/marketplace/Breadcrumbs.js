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
     //   console.log('Current language:', i18n.language);
     //   console.log('All translations:', pathCategory.translations);
      //  console.log('Serbian translation:', pathCategory.translations?.sr);
      //  console.log('Russian translation:', pathCategory.translations?.ru);
        
        // Сначала пробуем найти перевод для текущего языка
        if (i18n.language === 'sr' && pathCategory.translations?.sr) {
            return pathCategory.translations.sr;
        }
    
        // Если не нашли, возвращаем в порядке sr -> ru -> en -> name
        const translations = pathCategory.translations || {};
        return translations[i18n.language] || translations.sr || translations.ru || translations.en || pathCategory.name;
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
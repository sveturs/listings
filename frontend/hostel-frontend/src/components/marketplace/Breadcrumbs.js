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
    const { t } = useTranslation('common');
    const { i18n } = useTranslation('marketplace');
    const theme = useTheme();
    const navigate = useNavigate();

    const getTranslatedName = (pathCategory) => {
        // Ищем полную информацию о категории из общего списка категорий
        const fullCategory = categories?.find(c => c.id === pathCategory.id);
        
        if (fullCategory?.translations?.[i18n.language]) {
            return fullCategory.translations[i18n.language];
        }
        
        return pathCategory.name;
    };

    if (!paths || paths.length === 0) {
        return null;
    }

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
                    {t('navigation.home')}
                </Link>

                {paths.map((path, index) => {
                    const isLast = index === paths.length - 1;

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
                                {getTranslatedName(path)}
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
                                transition: 'all 0.2s',
                                '&:hover': {
                                    backgroundColor: theme.palette.action.hover,
                                    color: theme.palette.primary.main
                                }
                            }}
                        >
                            {getTranslatedName(path)}
                        </Link>
                    );
                })}
            </MuiBreadcrumbs>
        </Box>
    );
};

export default Breadcrumbs;